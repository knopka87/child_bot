package telegram

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"child-bot/api/internal/store"
	"child-bot/api/internal/util"
	"child-bot/api/internal/v1/types"
)

type selectionContext struct {
	Image        []byte
	Mime         string
	MediaGroupID string
	Detect       types.DetectResponse
}

type parsePending struct {
	Sc  *selectionContext
	PR  types.ParseResponse
	LLM string // "gemini"|"gpt"
}

func (r *Router) runParseAndMaybeConfirm(ctx context.Context, chatID int64, userID *int64, sc *selectionContext, subjectHint types.Subject, gradeHint *int64) {
	setState(chatID, Parse)
	imgHash := util.SHA256Hex(sc.Image)
	llmName := r.LlmManager.Get(chatID)
	sid, _ := r.getSession(chatID)

	grade := gradeHint
	if user, err := r.Store.FindUserByChatID(ctx, chatID); err == nil && user.Grade != nil {
		grade = user.Grade
	}

	in := types.ParseRequest{
		Image:       base64.StdEncoding.EncodeToString(sc.Image),
		Locale:      "ru_RU",
		SubjectHint: &subjectHint,
		GradeHint:   grade,
	}
	start := time.Now()
	pr, err := r.GetLLMClient().Parse(ctx, llmName, in)
	latency := time.Since(start).Milliseconds()
	_ = r.Store.InsertHistory(ctx, store.TimelineEvent{
		ChatID:        chatID,
		TaskSessionID: sid,
		Direction:     "api",
		EventType:     string(Parse),
		Provider:      llmName,
		OK:            err == nil,
		LatencyMS:     &latency,
		InputPayload:  in,
		OutputPayload: pr,
		Error:         err,
	})
	if err != nil {
		_ = r.Store.InsertEvent(ctx, store.MetricEvent{
			Stage:      "parse",
			Provider:   llmName,
			OK:         false,
			Error:      err.Error(),
			DurationMS: latency,
			ChatID:     &chatID,
			UserIDAnon: userID,
		})

		util.PrintError("runParseAndMaybeConfirm", llmName, chatID, "parse", err)
		r.sendError(chatID, fmt.Errorf("parse: %w", err))
		return
	}

	r.sendDebug(chatID, "parse_req", in)
	r.sendDebug(chatID, "parse_res", pr)

	// 3) Метрики строго по новой структуре
	_ = r.Store.InsertEvent(ctx, store.MetricEvent{
		Stage:      "parse",
		Provider:   llmName,
		OK:         true,
		DurationMS: time.Since(start).Milliseconds(),
		ChatID:     &chatID,
		UserIDAnon: userID,
		Details: map[string]any{
			"need_user_confirmation": pr.NeedsUserConfirmation,
			"task_type":              pr.TaskStruct.Type,
		},
	})
	util.PrintInfo("runParseAndMaybeConfirm", llmName, chatID, fmt.Sprintf("Received a response from LLMClient: %d", time.Since(start).Milliseconds()))

	// 4) Сохраняем черновик PARSE в БД
	js, _ := json.Marshal(pr)
	gradeValue := int64(0)
	if grade != nil {
		gradeValue = *grade
	}
	data := store.ParsedTasks{
		CreatedAt:             time.Now(),
		ChatID:                chatID,
		SessionID:             sid,
		MediaGroupID:          sc.MediaGroupID,
		ImageHash:             imgHash,
		Engine:                llmName,
		Subject:               pr.TaskStruct.Subject,
		Grade:                 gradeValue,
		RawTaskText:           pr.RawTaskText,
		ResultJSON:            js,
		NeedsUserConfirmation: pr.NeedsUserConfirmation,
		TaskType:              pr.TaskStruct.Type,
		CombinedSubpoints:     pr.TaskStruct.CombinedSubpoints,
		Accepted:              !pr.NeedsUserConfirmation,
		AcceptReason:          "",
	}
	if errP := r.Store.UpsertParse(ctx, data); errP != nil {
		util.PrintError("runParseAndMaybeConfirm", llmName, chatID, "error upsert parsed_tasks", errP)
	}

	r.askParseConfirmation(chatID, pr)
	parseWait.Store(chatID, &parsePending{Sc: sc, PR: pr, LLM: llmName})
}

// Показ запроса подтверждения распознанного текста
func (r *Router) askParseConfirmation(chatID int64, pr types.ParseResponse) {
	var b strings.Builder
	if s := strings.TrimSpace(pr.RawTaskText); s != "" {
		b.WriteString("```\n")
		b.WriteString(s)
		b.WriteString("\n```\n")
	}

	text := fmt.Sprintf(TaskViewText, b.String())
	r.sendMarkdown(chatID, text, makeParseConfirmButtons())
}
