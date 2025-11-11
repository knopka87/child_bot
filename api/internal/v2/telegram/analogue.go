package telegram

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"child-bot/api/internal/store"
	"child-bot/api/internal/v2/types"
)

// HandleAnalogueCallback — публичный помощник для существующего handleCallback
func (r *Router) HandleAnalogueCallback(chatID int64, userID *int64, reason types.AnalogueReason) {
	ctx := context.Background()
	if err := r.runAnalogue(ctx, chatID, userID, reason, "ru_RU"); err != nil {
		r.sendError(chatID, err)
	}
}

// runAnalogue — собирает вход из последнего подтверждённого парсинга и вызывает LLMClient-прокси
func (r *Router) runAnalogue(ctx context.Context, chatID int64, userID *int64, reason types.AnalogueReason, locale string) error {
	sid, _ := r.getSession(chatID)
	in, err := r.buildAnalogueInput(ctx, sid, reason, locale)
	if err != nil {
		return err
	}
	llmName := r.LlmManager.Get(chatID)
	start := time.Now()
	ar, err := r.GetLLMClient().AnalogueSolution(ctx, llmName, in)
	latency := time.Since(start).Milliseconds()
	_ = r.Store.InsertHistory(ctx, store.TimelineEvent{
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
		_ = r.Store.InsertEvent(ctx, store.MetricEvent{
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

	_ = r.Store.InsertEvent(ctx, store.MetricEvent{
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
	r.sendAnalogueResult(chatID, ar)

	return nil
}

// buildAnalogueInput — конструирует вход для ANALOGUE из данных последнего парсинга
func (r *Router) buildAnalogueInput(ctx context.Context, sid string, reason types.AnalogueReason, locale string) (types.AnalogueRequest, error) {
	if r.Store == nil {
		return types.AnalogueRequest{}, errors.New("store is not configured")
	}
	pr, ok := r.Store.FindLastConfirmedParse(ctx, sid)
	if !ok {
		return types.AnalogueRequest{}, errors.New("нет подтверждённого задания — пришлите фото и подтвердите распознавание")
	}

	grade := pr.Grade
	if user, err := r.Store.FindUserByChatID(ctx, pr.ChatID); err != nil && user.Grade != nil {
		grade = *user.Grade
	}

	in := types.AnalogueRequest{
		TaskStruct: types.TaskStruct{
			Subject:           pr.Subject,
			Type:              pr.TaskType,
			CombinedSubpoints: pr.CombinedSubpoints,
		},
		Reason:      reason,
		Locale:      locale,
		Grade:       grade,
		RawTaskText: pr.RawTaskText,
	}
	return in, nil
}

// sendAnalogueResult — формирует человекочитаемый вывод без раскрытия ответа исходника
func (r *Router) sendAnalogueResult(chatID int64, ar types.AnalogueResponse) {
	var b strings.Builder

	b.WriteString(ar.ExampleTask)

	if len(ar.SolutionSteps) > 0 {
		b.WriteString(StepSolutionText)
	}
	for i, step := range ar.SolutionSteps {
		b.WriteString(strconv.Itoa(i+1) + ". " + step + "\n\n")
	}

	text := fmt.Sprintf(AnalogueTaskText, b.String())

	r.send(chatID, text, makeAnalogueButtons())
}
