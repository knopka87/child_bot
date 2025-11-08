package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"child-bot/api/internal/store"
	"child-bot/api/internal/util"
	"child-bot/api/internal/v2/types"
)

// lastParseMeta — извлекает метаданные последнего подтверждённого парсинга
func (r *Router) lastParseMeta(ctx context.Context, sid string) (subject string, taskType string, grade int, ctxParse json.RawMessage) {
	if pt, ok := r.ParseRepo.FindLastConfirmed(ctx, sid); ok {
		subject = pt.Subject
		taskType = pt.TaskType
		grade = pt.Grade
		ctxParse = pt.ResultJSON
	}

	return
}

// normalizeText — отправляет текст ученика на нормализацию в LLMClient-прокси
func (r *Router) normalizeText(ctx context.Context, chatID int64, userID *int64, text string) {
	setState(chatID, Normalize)
	llmName := r.LlmManager.Get(chatID)

	text = strings.TrimSpace(text)
	if text == "" {
		r.send(chatID, "Пожалуйста, пришлите текст ответа.", nil)
		return
	}

	sid, _ := r.getSession(chatID)
	_, _, _, parseCtx := r.lastParseMeta(ctx, sid)

	r.sendDebug(chatID, "parse context", parseCtx)

	var pr types.ParseResponse
	_ = json.Unmarshal(parseCtx, &pr)
	r.sendDebug(chatID, "parse response", pr)

	in := types.NormalizeRequest{
		TaskStruct:    pr.TaskStruct,
		RawTaskText:   pr.RawTaskText,
		RawAnswerText: text,
	}
	util.PrintInfo("normalizeText", r.LlmManager.Get(chatID), chatID, fmt.Sprintf("normalize_input: %+v", in))
	start := time.Now()
	res, err := r.GetLLMClient().Normalize(ctx, llmName, in)
	latency := time.Since(start).Milliseconds()
	_ = r.History.Insert(ctx, store.TimelineEvent{
		ChatID:        chatID,
		TaskSessionID: sid,
		Direction:     "api",
		EventType:     string(Normalize),
		Provider:      llmName,
		OK:            err == nil,
		LatencyMS:     &latency,
		InputPayload:  in,
		OutputPayload: res,
		Error:         err,
	})
	if err != nil {
		_ = r.Metrics.InsertEvent(ctx, store.MetricEvent{
			Stage:      "normalize",
			Provider:   llmName,
			OK:         false,
			Error:      err.Error(),
			DurationMS: latency,
			ChatID:     &chatID,
			UserIDAnon: userID,
			Details: map[string]any{
				"source":      "text",
				"input_chars": len(text),
			},
		})

		b := make([][]tgbotapi.InlineKeyboardButton, 0, 2)
		b = append(b,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Перейти к новой задаче", "new_task")),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report")),
		)
		r.send(chatID, fmt.Sprintf("Не удалось нормализовать ответ: %v", err), b)
		return
	}

	r.sendDebug(chatID, "normalize_input", in)
	r.sendDebug(chatID, "normalize_req", res)

	_ = r.Metrics.InsertEvent(ctx, store.MetricEvent{
		Stage:      "normalize",
		Provider:   llmName,
		OK:         true,
		DurationMS: latency,
		ChatID:     &chatID,
		UserIDAnon: userID,
		Details: map[string]any{
			"source": "text",
		},
	})

	r.sendNormalizePreview(chatID)

	// Попробуем сразу проверить решение, если в системе есть ожидаемое решение
	r.checkSolution(ctx, chatID, userID, res)
	clearMode(chatID)
}

// sendNormalizePreview — короткий текст для пользователя по NormalizeResult
func (r *Router) sendNormalizePreview(chatID int64) {
	b := &strings.Builder{}
	b.WriteString("✅ Принял ответ.")
	r.send(chatID, b.String(), nil)
}
