package telegram

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"child-bot/api/internal/store"
	"child-bot/api/internal/v1/types"
)

// --- ANALOGUE SOLUTION (v1.1) ----------------------------------------------
// По кнопке «Похожее задание» генерируем аналог по тем же приёмам, но с другими данными.
// Основано на инструкции ANALOGUE_SOLUTION v1.1.

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
// Вызовите его из вашего обработчика, когда callback.Data == "ANALOGUE".
func (r *Router) HandleAnalogueCallback(chatID int64, userID *int64) {
	ctx := context.Background()
	if err := r.runAnalogue(ctx, chatID, userID); err != nil {
		b := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report"))
		r.send(chatID, "Не удалось подготовить аналогичное задание: "+err.Error(), b)
	}
}

// runAnalogue — собирает вход из последнего подтверждённого парсинга и вызывает LLMClient-прокси
func (r *Router) runAnalogue(ctx context.Context, chatID int64, userID *int64) error {
	in, err := r.buildAnalogueInput(ctx, chatID)
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
			"has_minichecks": len(ar.MiniChecks) > 0,
		},
	})
	r.sendAnalogueResult(chatID, ar)

	return nil
}

// buildAnalogueInput — конструирует вход для ANALOGUE из данных последнего парсинга
func (r *Router) buildAnalogueInput(ctx context.Context, chatID int64) (types.AnalogueSolutionInput, error) {
	if r.ParseRepo == nil {
		return types.AnalogueSolutionInput{}, errors.New("ParseRepo is not configured")
	}
	sid, _ := r.getSession(chatID)
	tasks, ok := r.ParseRepo.FindLastConfirmed(ctx, sid)
	if !ok {
		return types.AnalogueSolutionInput{}, errors.New("нет подтверждённого задания — пришлите фото и подтвердите распознавание")
	}

	var p types.ParseResult
	_ = json.Unmarshal(tasks.ResultJSON, &p)

	// Берём краткую суть, либо строим её из вопроса/сырого текста, удаляя числа/единицы
	base := strings.TrimSpace(tasks.Question)
	if base == "" {
		base = strings.TrimSpace(tasks.RawTaskText)
	}
	norm := stripNumbersUnits(base)
	if norm == "" {
		return types.AnalogueSolutionInput{}, errors.New("не удалось получить краткую суть задания")
	}

	in := types.AnalogueSolutionInput{
		TaskID:              sid,
		UserIDAnon:          fmt.Sprint(chatID),
		Grade:               tasks.Grade,
		Subject:             tasks.Subject,  // "math"|"russian"|...
		TaskType:            tasks.TaskType, // если классификатор есть
		OriginalTaskEssence: norm,           // без чисел/единиц исходника
		Locale:              "ru",
	}
	return in, nil
}

var reNums = regexp.MustCompile(`(?i)(\d+[\d\s./,:-]*\d*|\d+)`)
var reUnits = regexp.MustCompile(`(?i)(см|мм|м|кг|г|л|мл|ч|мин|сек|%|грн|руб|р\.|км)\.?`)

// stripNumbersUnits — удаляет из текста числа и типичные единицы/знаки, чтобы
// получить краткую суть без утечки исходных данных (см. anti‑leak в v1.1)
func stripNumbersUnits(s string) string {
	out := reNums.ReplaceAllString(s, "N")
	out = reUnits.ReplaceAllString(out, "U")
	out = strings.TrimSpace(strings.Join(strings.Fields(out), " "))
	return out
}

// sendAnalogueResult — формирует человекочитаемый вывод без раскрытия ответа исходника
func (r *Router) sendAnalogueResult(chatID int64, ar types.AnalogueSolutionResult) {
	var b strings.Builder
	if t := strings.TrimSpace(ar.AnalogyTitle); t != "" {
		b.WriteString("📘 ")
		b.WriteString(t)
		b.WriteString("\n\n")
	}
	if t := strings.TrimSpace(ar.AnalogyTask); t != "" {
		b.WriteString("Похожее задание:\n")
		b.WriteString(t)
		b.WriteString("\n\n")
	}
	if len(ar.SolutionSteps) > 0 {
		b.WriteString("Как решать (тот же приём):\n")
		for i, s := range ar.SolutionSteps {
			b.WriteString(strconv.Itoa(i + 1))
			b.WriteString(". ")
			b.WriteString(strings.TrimSpace(s))
			b.WriteString("\n")
		}
	}
	if len(ar.TransferBridge) > 0 {
		b.WriteString("\nМостик переноса:\n")
		b.WriteString(ar.TransferBridge)
	}
	if s := strings.TrimSpace(ar.TransferCheck); s != "" {
		b.WriteString("\n\nПроверь себя: ")
		b.WriteString(s)
	}
	// Мини‑проверки: показываем без ответов
	if len(ar.MiniChecks) > 0 {
		b.WriteString("\n\nМини‑проверки:\n")
		for _, mc := range ar.MiniChecks {
			p := strings.TrimSpace(mc.Prompt)
			if p == "" && mc.Raw != "" {
				p = mc.Raw
			}
			if p != "" {
				b.WriteString("— ")
				b.WriteString(p)
				b.WriteString("\n")
			}
		}
	}
	// Короткая подсказка по безопасности/антилику
	if !ar.LeakGuardPassed || !ar.Safety.NoOriginalAnswerLeak {
		b.WriteString("\n(Замечание: аналог без ссылок на исходные данные, ответы не раскрываются.)")
	}
	r.send(chatID, b.String(), nil)
}
