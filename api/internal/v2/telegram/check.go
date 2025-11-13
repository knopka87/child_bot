package telegram

import (
	"context"
	"fmt"
	"strings"
	"time"

	"child-bot/api/internal/store"
	"child-bot/api/internal/v2/types"
)

// checkSolution — если есть ожидаемое решение для текущей задачи, проверяем ответ
func (r *Router) checkSolution(ctx context.Context, chatID int64, userID *int64, nr types.NormalizeResponse) {
	setState(chatID, Check)
	sid, _ := r.getSession(chatID)

	time1 := r.sendAlert(chatID, CheckAlert, 0, 15)

	// 0) Подтянем метаданные предмета/класса из последнего подтверждённого парсинга
	subj := "generic"
	grade := int64(0)
	if pr, ok := r.Store.FindLastConfirmedParse(ctx, sid); ok {
		if s := strings.TrimSpace(pr.Subject); s != "" {
			subj = s
			grade = pr.Grade
		}
	}
	if user, err := r.Store.FindUserByChatID(ctx, chatID); err == nil && user.Grade != nil {
		grade = *user.Grade
	}

	llmName := r.LlmManager.Get(chatID)
	in := types.CheckRequest{
		NormAnswer: nr.NormAnswer,
		NormTask:   nr.NormTask,
		Grade:      grade,
	}

	start := time.Now()
	res, err := r.GetLLMClient().CheckSolution(ctx, llmName, in)
	latency := time.Since(start).Milliseconds()
	_ = r.Store.InsertHistory(ctx, store.TimelineEvent{
		ChatID:        chatID,
		TaskSessionID: sid,
		Direction:     "api",
		EventType:     string(Check),
		Provider:      llmName,
		OK:            err == nil,
		LatencyMS:     &latency,
		InputPayload:  in,
		OutputPayload: res,
		Error:         err,
	})
	if err != nil {
		_ = r.Store.InsertEvent(ctx, store.MetricEvent{
			Stage:      "check",
			Provider:   llmName,
			OK:         false,
			Error:      err.Error(),
			DurationMS: latency,
			ChatID:     &chatID,
			UserIDAnon: userID,
			Details: map[string]any{
				"subject": subj,
			},
		})

		r.sendError(chatID, fmt.Errorf("не удалось проверить решение: %v", err))
		return
	}

	r.sendDebug(chatID, "check input", in)
	r.sendDebug(chatID, "check res", res)

	_ = r.Store.InsertEvent(ctx, store.MetricEvent{
		Stage:      "check",
		Provider:   llmName,
		OK:         true,
		DurationMS: latency,
		ChatID:     &chatID,
		UserIDAnon: userID,
		Details: map[string]any{
			"subject":    subj,
			"confidence": res.Confidence,
			"is_correct": res.IsCorrect,
		},
	})

	time1.Stop()
	r.sendCheckResponse(chatID, res)
}

// sendCheckResponse — вывод краткого результата проверки (с учётом новой схемы v1.2)
func (r *Router) sendCheckResponse(chatID int64, cr types.CheckResponse) {
	var b strings.Builder

	if cr.IsCorrect {
		setState(chatID, Correct)
		clearMode(chatID)
		r.clearSession(chatID)
		r.send(chatID, AnswerCorrectText, makeCorrectAnswerButtons())
		return
	}

	setState(chatID, Incorrect)

	if cr.Feedback != "" {
		b.WriteString("\n" + cr.Feedback + "\n")
	}

	text := fmt.Sprintf(AnswerIncorrectText, cr.Feedback)
	r.send(chatID, text, makeIncorrectAnswerButtons())
}
