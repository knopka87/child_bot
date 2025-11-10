package telegram

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

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

func (r *Router) runParseAndMaybeConfirm(ctx context.Context, chatID int64, userID *int64, sc *selectionContext, subjectHint types.Subject, gradeHint *int) {
	setState(chatID, Parse)
	imgHash := util.SHA256Hex(sc.Image)
	llmName := r.LlmManager.Get(chatID)
	sid, _ := r.getSession(chatID)

	// 1) Проверка кэша: если уже было подтверждено ранее — используем сразу
	if prRow, ok := r.ParseRepo.FindLastConfirmed(ctx, sid); ok {
		pr := types.ParseResponse{
			RawTaskText: prRow.RawTaskText,
			TaskStruct: types.TaskStruct{
				Subject:           prRow.Subject,
				Type:              prRow.TaskType,
				CombinedSubpoints: prRow.CombinedSubpoints,
			},
			NeedsUserConfirmation: prRow.NeedsUserConfirmation,
		}
		r.showTaskAndPrepareHints(chatID, sc, pr, llmName)
		return
	}

	// 2) Запрос к LLMClient.Parse
	in := types.ParseRequest{
		Image:       base64.StdEncoding.EncodeToString(sc.Image),
		Locale:      "ru_RU",
		SubjectHint: &subjectHint,
		GradeHint:   gradeHint,
	}
	start := time.Now()
	pr, err := r.GetLLMClient().Parse(ctx, llmName, in)
	latency := time.Since(start).Milliseconds()
	_ = r.History.Insert(ctx, store.TimelineEvent{
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
		_ = r.Metrics.InsertEvent(ctx, store.MetricEvent{
			Stage:      "parse",
			Provider:   llmName,
			OK:         false,
			Error:      err.Error(),
			DurationMS: latency,
			ChatID:     &chatID,
			UserIDAnon: userID,
		})

		util.PrintError("runParseAndMaybeConfirm", llmName, chatID, "parse", err)
		r.SendError(chatID, fmt.Errorf("parse: %w", err))
		return
	}

	r.sendDebug(chatID, "parse_req", in)
	r.sendDebug(chatID, "parse_res", pr)

	// 3) Метрики строго по новой структуре
	_ = r.Metrics.InsertEvent(ctx, store.MetricEvent{
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
	data := store.ParsedTasks{
		CreatedAt:             time.Now(),
		ChatID:                chatID,
		SessionID:             sid,
		MediaGroupID:          sc.MediaGroupID,
		ImageHash:             imgHash,
		Engine:                llmName,
		Subject:               pr.TaskStruct.Subject,
		RawTaskText:           pr.RawTaskText,
		ResultJSON:            js,
		NeedsUserConfirmation: pr.NeedsUserConfirmation,
		TaskType:              pr.TaskStruct.Type,
		CombinedSubpoints:     pr.TaskStruct.CombinedSubpoints,
		Accepted:              !pr.NeedsUserConfirmation,
		AcceptReason:          "",
	}
	if errP := r.ParseRepo.Upsert(ctx, data); errP != nil {
		util.PrintError("runParseAndMaybeConfirm", llmName, chatID, "error upsert parsed_tasks", errP)
	}

	// 6) Если требуется подтверждение — спрашиваем пользователя
	if pr.NeedsUserConfirmation {
		setState(chatID, Confirm)
		r.askParseConfirmation(chatID, pr)
		parseWait.Store(chatID, &parsePending{Sc: sc, PR: pr, LLM: llmName})
		return
	}

	// 7) Иначе — автоподтверждение и переход к подсказкам
	setState(chatID, AutoPick)
	_ = r.ParseRepo.MarkAcceptedBySession(ctx, sid, "auto")
	r.showTaskAndPrepareHints(chatID, sc, pr, llmName)
	util.PrintInfo("runParseAndMaybeConfirm", llmName, chatID, fmt.Sprintf("total time: %d", time.Since(start).Milliseconds()))
}

// Показ запроса подтверждения распознанного текста
func (r *Router) askParseConfirmation(chatID int64, pr types.ParseResponse) {
	var b strings.Builder
	b.WriteString("Я так прочитал задание. Всё верно?\n")
	if s := strings.TrimSpace(pr.RawTaskText); s != "" {
		b.WriteString("```\n")
		b.WriteString(s)
		b.WriteString("\n```\n")
	}

	msg := tgbotapi.NewMessage(chatID, b.String())
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = makeParseConfirmKeyboard()
	_, _ = r.Bot.Send(msg)
}
