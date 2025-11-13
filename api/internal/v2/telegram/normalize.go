package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"child-bot/api/internal/store"
	"child-bot/api/internal/v2/types"
)

// lastParseMeta — извлекает метаданные последнего подтверждённого парсинга
func (r *Router) lastParseMeta(ctx context.Context, sid string) (subject string, taskType string, grade int64, ctxParse json.RawMessage) {
	if pt, ok := r.Store.FindLastConfirmedParse(ctx, sid); ok {
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

	time1 := r.sendAlert(chatID, NormaliseAlert1, 0, 15)
	time2 := r.sendAlert(chatID, NormaliseAlert2, 15, 30)

	text = strings.TrimSpace(text)
	if text == "" {
		r.sendError(chatID, fmt.Errorf("распознан пустой ответ"))
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
	// util.PrintInfo("normalizeText", r.LlmManager.Get(chatID), chatID, fmt.Sprintf("normalize_input: %+v", in))
	start := time.Now()
	res, err := r.GetLLMClient().Normalize(ctx, llmName, in)
	latency := time.Since(start).Milliseconds()
	_ = r.Store.InsertHistory(ctx, store.TimelineEvent{
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
		_ = r.Store.InsertEvent(ctx, store.MetricEvent{
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

		r.sendError(chatID, fmt.Errorf("не удалось нормализовать ответ: %v", err))
		return
	}

	r.sendDebug(chatID, "normalize_input", in)
	r.sendDebug(chatID, "normalize_req", res)

	_ = r.Store.InsertEvent(ctx, store.MetricEvent{
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

	time2.Stop()
	time1.Stop()
	r.checkSolution(ctx, chatID, userID, res)
	clearMode(chatID)
}
