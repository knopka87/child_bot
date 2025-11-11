package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"child-bot/api/internal/store"
	"child-bot/api/internal/v1/types"
)

// checkSolution — если есть ожидаемое решение для текущей задачи, проверяем ответ
func (r *Router) checkSolution(ctx context.Context, chatID int64, userID *int64, nr types.NormalizeResponse) {
	setState(chatID, Check)
	sid, _ := r.getSession(chatID)

	// 0) Подтянем метаданные предмета/класса из последнего подтверждённого парсинга
	subj := "generic"
	if pr, ok := r.Store.FindLastConfirmedParse(ctx, sid); ok {
		if s := strings.TrimSpace(pr.Subject); s != "" {
			subj = s
		}
	}

	// 1) Определим ветку проверки из предмета/контекста
	branch := r.detectCheckBranch(subj)

	// 2) Пытаемся взять ожидаемое решение (внутренний JSON) из БД
	exp, ok := r.getExpectedForChat(ctx, chatID)
	if !ok || len(exp) == 0 {
		// 3) Фолбэк: передаём пустой объект — модель проведёт policy‑only проверку
		exp = json.RawMessage(`{}`)
	}

	llmName := r.LlmManager.Get(chatID)
	in := types.CheckRequest{
		NormAnswer: nr.NormAnswer,
		NormTask:   nr.NormTask,
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
				"branch":  branch,
			},
		})

		b := make([][]tgbotapi.InlineKeyboardButton, 0, 2)
		b = append(b,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Перейти к новой задаче", "new_task")),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📝 Сообщить об ошибке", "report")),
		)
		r.send(chatID, fmt.Sprintf("Не удалось проверить решение: %v", err), b)
		r.offerAnalogueButton(chatID)
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

	r.sendCheckResponse(chatID, res)
}

// getExpectedForChat — извлекает ожидаемое решение (внутренний эталон JSON) для текущей задачи чата
func (r *Router) getExpectedForChat(ctx context.Context, chatID int64) (json.RawMessage, bool) {
	// Если в ParseRepo хранится сырой JSON эталона, раскомментируйте:
	// if r.ParseRepo != nil {
	// 	if pr, ok := r.Store.FindLastConfirmedParseParse(ctx, chatID); ok {
	// 		// Возможные варианты поля в модели парсинга:
	// 		// 1) pr.ExpectedSolution []byte / json.RawMessage
	// 		// 2) pr.Expected json.RawMessage
	// 		// 3) pr.ExpectedObject (структура) — тогда нужно: b, _ := json.Marshal(pr.ExpectedObject); return json.RawMessage(b), true
	// 		if len(pr.ExpectedSolution) > 0 {
	// 			return pr.ExpectedSolution, true
	// 		}
	// 		if len(pr.Expected) > 0 {
	// 			return pr.Expected, true
	// 		}
	// 	}
	// }
	return nil, false
}

// detectCheckBranch — маппинг предмета/контекста к ветке проверки схемы
func (r *Router) detectCheckBranch(subject string) string {
	s := strings.ToLower(strings.TrimSpace(subject))
	if strings.Contains(s, "мат") || s == "math" || s == "математика" {
		return "math_branch"
	}
	if strings.Contains(s, "рус") || s == "russian" || s == "русский язык" {
		return "ru_branch"
	}
	return "generic_branch"
}

// sendCheckResponse — вывод краткого результата проверки (с учётом новой схемы v1.2)
func (r *Router) sendCheckResponse(chatID int64, cr types.CheckResponse) {
	var b strings.Builder

	if cr.IsCorrect {
		setState(chatID, Correct)
		b.WriteString("✅ Задача решена верно\n")
	} else {
		setState(chatID, Incorrect)
		b.WriteString("⚠️ Похоже, есть неточности в решении\n")
	}

	if cr.Feedback != "" {
		b.WriteString("\n" + cr.Feedback + "\n")
	}

	if getState(chatID) == Correct {
		b.WriteString("\n\nГотов двигаться дальше — присылай следующую задачу.")
		clearMode(chatID)
		r.clearSession(chatID)
	} else {
		b.WriteString("\nЕсли нужно — могу подобрать похожее задание для тренировки.")
	}

	r.send(chatID, b.String(), nil)

	// Предлагаем «Похожее задание», если решение не подтверждено
	if getState(chatID) != Correct {
		r.offerAnalogueButton(chatID)
	}
}
