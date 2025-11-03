package telegram

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"child-bot/api/internal/store"
	"child-bot/api/internal/v2/types"
)

// По кнопке «Похожее задание» генерируем аналог по тем же приёмам, но с другими данными.

// offerAnalogueButton — показывает кнопку для вызова аналога
func (r *Router) offerAnalogueButton(chatID int64) {
	kb := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("Похожее задание", "analogue_solution"),
		},
	)
	msg := tgbotapi.NewMessage(chatID, "Если нужно, покажу похожее задание тем же приёмом (без ответа исходной задачи).")
	msg.ReplyMarkup = kb
	_, _ = r.Bot.Send(msg)
}

// HandleAnalogueCallback — публичный помощник для существующего handleCallback
func (r *Router) HandleAnalogueCallback(chatID int64, userID *int64, reason types.AnalogueReason) {
	ctx := context.Background()
	if err := r.runAnalogue(ctx, chatID, userID, reason, "ru_RU"); err != nil {
		b := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report"))
		r.send(chatID, "Не удалось подготовить аналогичное задание: "+err.Error(), b)
	}
}

// runAnalogue — собирает вход из последнего подтверждённого парсинга и вызывает LLMClient-прокси
func (r *Router) runAnalogue(ctx context.Context, chatID int64, userID *int64, reason types.AnalogueReason, locale string) error {
	in, err := r.buildAnalogueInput(ctx, chatID, reason, locale)
	if err != nil {
		return err
	}
	llmName := r.LlmManager.Get(chatID)
	start := time.Now()
	ar, err := r.GetLLMClient().AnalogueSolution(ctx, llmName, in)
	latency := time.Since(start).Milliseconds()
	sid, _ := r.getSession(chatID)
	_ = r.History.Insert(ctx, store.TimelineEvent{
		ChatID:        chatID,
		TaskSessionID: sid,
		Direction:     "api",
		EventType:     string(Analogue),
		Provider:      llmName,
		OK:            err == nil,
		LatencyMS:     &latency,
		InputPayload:  in,
		OutputPayload: ar,
		Error:         err,
	})
	if err != nil {
		_ = r.Metrics.InsertEvent(ctx, store.MetricEvent{
			Stage:      "analogue",
			Provider:   llmName,
			OK:         false,
			Error:      err.Error(),
			DurationMS: latency,
			ChatID:     &chatID,
			UserIDAnon: userID,
		})

		return err
	}

	_ = r.Metrics.InsertEvent(ctx, store.MetricEvent{
		Stage:      "analogue",
		Provider:   llmName,
		OK:         true,
		DurationMS: time.Since(start).Milliseconds(),
		ChatID:     &chatID,
		UserIDAnon: userID,
		Details: map[string]any{
			"solution_steps": len(ar.SolutionSteps),
		},
	})
	r.sendAnalogueResult(chatID, ar, reason)

	return nil
}

// buildAnalogueInput — конструирует вход для ANALOGUE из данных последнего парсинга
func (r *Router) buildAnalogueInput(ctx context.Context, chatID int64, reason types.AnalogueReason, locale string) (types.AnalogueRequest, error) {
	if r.ParseRepo == nil {
		return types.AnalogueRequest{}, errors.New("ParseRepo is not configured")
	}
	pr, ok := r.ParseRepo.FindLastConfirmed(ctx, chatID)
	if !ok {
		return types.AnalogueRequest{}, errors.New("нет подтверждённого задания — пришлите фото и подтвердите распознавание")
	}

	in := types.AnalogueRequest{
		TaskStruct: types.TaskStruct{
			Subject:           pr.Subject,
			Type:              pr.TaskType,
			CombinedSubpoints: pr.CombinedSubpoints,
		},
		Reason: reason,
		Locale: locale,
	}
	return in, nil
}

// sendAnalogueResult — формирует человекочитаемый вывод без раскрытия ответа исходника
func (r *Router) sendAnalogueResult(chatID int64, ar types.AnalogueResponse, reason types.AnalogueReason) {
	var b strings.Builder

	b.WriteString("Аналогичная задача\n\n")
	b.WriteString(ar.ExampleTask)

	if len(ar.SolutionSteps) > 0 {
		b.WriteString("\n\n\n\n📘 Шаги решения\n\n")
	}
	for i, step := range ar.SolutionSteps {
		b.WriteString(strconv.Itoa(i+1) + "." + step + "\n\n")
	}

	button := makeActionsKeyboardRow(3, false)
	r.send(chatID, b.String(), button)
}
