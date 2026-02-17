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
	"child-bot/api/internal/v2/types"
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

func (r *Router) runParseAndMaybeConfirm(ctx context.Context, chatID int64, userID *int64, sc *selectionContext, subjectCandidate types.Subject) {
	setState(chatID, Parse)
	imgHash := util.SHA256Hex(sc.Image)
	llmName := r.LlmManager.Get(chatID)
	sid, _ := r.getSession(chatID)

	var grade int64
	if user, err := r.Store.FindUserByChatID(ctx, chatID); err == nil && user.Grade != nil {
		grade = *user.Grade
	}

	confidence := "high"
	if sc.Detect.Classification.Confidence < 0.7 {
		confidence = "low"
	} else if sc.Detect.Classification.Confidence < 0.9 {
		confidence = "medium"
	}

	in := types.ParseRequest{
		Image:             base64.StdEncoding.EncodeToString(sc.Image),
		TaskId:            sid,
		Locale:            "ru_RU",
		SubjectCandidate:  string(subjectCandidate),
		SubjectConfidence: confidence,
		Grade:             grade,
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
	taskType := ""
	if len(pr.Items) > 0 {
		taskType = pr.Items[0].PedKeys.TaskType
	}
	_ = r.Store.InsertEvent(ctx, store.MetricEvent{
		Stage:      "parse",
		Provider:   llmName,
		OK:         true,
		DurationMS: time.Since(start).Milliseconds(),
		ChatID:     &chatID,
		UserIDAnon: userID,
		Details: map[string]any{
			"subject":     pr.Task.Subject,
			"task_type":   taskType,
			"items_count": len(pr.Items),
		},
	})
	util.PrintInfo("runParseAndMaybeConfirm", llmName, chatID, fmt.Sprintf("Received a response from LLMClient: %d", time.Since(start).Milliseconds()))

	// 4) Сохраняем черновик PARSE в БД
	js, _ := json.Marshal(pr)
	data := store.ParsedTasks{
		CreatedAt:             time.Now(),
		ChatID:                chatID,
		SessionID:             sid,
		MediaGroupID:          sc.MediaGroupID,
		ImageHash:             imgHash,
		Engine:                llmName,
		Subject:               string(pr.Task.Subject),
		Grade:                 pr.Task.Grade,
		RawTaskText:           pr.Task.TaskTextClean,
		ResultJSON:            js,
		NeedsUserConfirmation: false,
		TaskType:              taskType,
		CombinedSubpoints:     false,
		Accepted:              true,
		AcceptReason:          "",
		TaskID:                pr.Task.TaskId,
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
	if s := strings.TrimSpace(pr.Task.TaskTextClean); s != "" {
		b.WriteString("```\n")
		b.WriteString(s)
		b.WriteString("\n```\n")
	}

	// Добавляем информацию о педагогическом шаблоне
	if templateID := getTemplateID(pr.Task, pr.Items); templateID != "" {
		b.WriteString(fmt.Sprintf("\n🎓 Шаблон: `%s`", templateID))
	} else {
		b.WriteString("\n🎓 Шаблон: не удалось подобрать")
	}

	text := fmt.Sprintf(TaskViewText, b.String())
	r.sendMarkdown(chatID, text, makeParseConfirmButtons())
}
