package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"child-bot/api/internal/ocr/types"
	"child-bot/api/internal/store"
)

// maybeCheckSolution — если есть ожидаемое решение для текущей задачи, проверяем ответ
func (r *Router) maybeCheckSolution(ctx context.Context, chatID int64, userID *int64, nr types.NormalizeResult) {
	setState(chatID, Check)
	// 0) Подтянем метаданные предмета/класса из последнего подтверждённого парсинга
	subj := "math"
	grade := 0
	var parseCtx json.RawMessage
	if r.ParseRepo != nil {
		if pr, ok := r.ParseRepo.FindLastConfirmed(ctx, chatID); ok {
			subj = strings.TrimSpace(pr.Subject)
			grade = pr.Grade
			parseCtx, _ = json.Marshal(pr.Parse)
		}
	}

	// 1) Пытаемся взять ожидаемое решение из БД
	exp, ok := r.getExpectedForChat(ctx, chatID)
	if !ok {
		// 2) Фолбэк: строим «policy-only» ожидание по данным нормализации ученика
		shape := strings.TrimSpace(nr.Shape)
		if shape == "" && nr.ShapeDetected != nil {
			shape = strings.TrimSpace(*nr.ShapeDetected)
		}
		if shape == "" {
			shape = "number"
		}

		var units *types.UnitsExpectedSpec
		if nr.Units != nil {
			policy := "optional"
			if nr.Units.Kept != nil && *nr.Units.Kept {
				policy = "required"
			}
			primary := ""
			if nr.Units.Canonical != nil {
				primary = strings.TrimSpace(*nr.Units.Canonical)
			}
			alts := []string{}
			if nr.Units.Detected != nil {
				det := strings.TrimSpace(*nr.Units.Detected)
				if det != "" && det != primary {
					alts = append(alts, det)
				}
			}
			units = &types.UnitsExpectedSpec{
				Policy:          policy,  // требуем/не требуем единицы
				ExpectedPrimary: primary, // если нормализация вывела канон. единицу
				Alternatives:    alts,    // допустимые альтернативы
			}
		}

		exp = types.ExpectedSolution{
			Shape: shape,
			Units: units,
			// Number/String/List/Steps — не задаём без эталона, чтобы не «подгонять» под ответ
		}
	}

	llmName := r.EngManager.Get(chatID)
	in := types.CheckSolutionInput{
		UserIDAnon:   fmt.Sprint(chatID),
		Subject:      subj,
		Grade:        grade,
		Student:      nr,
		Expected:     exp,
		ParseContext: parseCtx,
	}
	r.sendDebug(chatID, "check_solution_input", in)
	start := time.Now()
	res, err := r.LLM.CheckSolution(ctx, llmName, in)
	latency := time.Since(start).Milliseconds()
	sid, _ := r.getSession(chatID)
	_ = r.History.Insert(ctx, store.TimelineEvent{
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
		_ = r.Metrics.InsertEvent(ctx, store.MetricEvent{
			Stage:      "check",
			Provider:   llmName,
			OK:         false,
			Error:      err.Error(),
			DurationMS: latency,
			ChatID:     &chatID,
			UserIDAnon: userID,
			Details: map[string]any{
				"subject": subj,
				"grade":   grade,
			},
		})

		b := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report"))
		r.send(chatID, fmt.Sprintf("Не удалось проверить решение: %v", err), b)
		r.offerAnalogueButton(chatID)
		return
	}

	_ = r.Metrics.InsertEvent(ctx, store.MetricEvent{
		Stage:      "check",
		Provider:   llmName,
		OK:         true,
		DurationMS: latency,
		ChatID:     &chatID,
		UserIDAnon: userID,
		Details: map[string]any{
			"subject": subj,
			"grade":   grade,
			"verdict": res.Verdict,
		},
	})

	r.sendCheckResult(chatID, res)
}

// getExpectedForChat — извлекает ожидаемое решение из вашей БД для текущей задачи чата
func (r *Router) getExpectedForChat(ctx context.Context, chatID int64) (types.ExpectedSolution, bool) {
	// if r.ParseRepo != nil {
	// 	if pr, ok := r.ParseRepo.FindLastConfirmed(ctx, chatID); ok {
	// 		return pr.Expected, true
	// 	}
	// }
	var exp types.ExpectedSolution
	return exp, false
}

// sendCheckResult — вывод краткого результата проверки
func (r *Router) sendCheckResult(chatID int64, cr types.CheckSolutionResult) {
	var b strings.Builder

	// 1) Вердикт
	switch strings.ToLower(strings.TrimSpace(cr.Verdict)) {
	case "correct":
		setState(chatID, Correct)
		b.WriteString("✅ Задача решена верно\n")
	case "incorrect":
		setState(chatID, Incorrect)
		b.WriteString("⚠️ Похоже, есть ошибка\n")
	case "uncertain":
		setState(chatID, Uncertain)
		b.WriteString("🤔 Пока не уверен в оценке\n")
	default:
		setState(chatID, Uncertain)
		b.WriteString("Результат проверки получен\n")
	}

	// 2) Короткая подсказка от проверки (без раскрытия ответа)
	if s := strings.TrimSpace(cr.ShortHint); s != "" {
		b.WriteString("Подсказка: ")
		b.WriteString(s)
		b.WriteString("\n")
	}

	// 3) Коды причин (если есть) — компактно
	if len(cr.ReasonCodes) > 0 {
		b.WriteString("Причины: ")
		b.WriteString(strings.Join(cr.ReasonCodes, ", "))
		b.WriteString("\n")
	}

	// 4) Диагностическая сводка (без чисел/конкретных значений)
	c := cr.Comparison

	// Единицы измерения
	if u := c.Units; u != nil {
		b.WriteString("Единицы: ")
		if u.Applied != nil && strings.TrimSpace(*u.Applied) != "" {
			// конверсия применена, без конкретных значений
			b.WriteString("конверсия применена\n")
		} else {
			// просто констатируем наличие/отсутствие
			if u.Detected != nil && strings.TrimSpace(*u.Detected) != "" {
				b.WriteString("указаны\n")
			} else {
				b.WriteString("не указаны\n")
			}
		}
	}

	// Числовая проверка
	if nd := c.NumberDiff; nd != nil {
		if nd.WithinTolerance {
			b.WriteString("Числовая проверка: в допустимых пределах\n")
		} else {
			// не раскрываем значения/форматы
			b.WriteString("Числовая проверка: требуется пересмотр\n")
		}
	}

	// Словесная проверка (для русского языка)
	if sm := c.StringMatch; sm != nil {
		mode := strings.TrimSpace(sm.Mode)
		if mode == "" {
			mode = "по тексту"
		}
		b.WriteString("Словесная сверка: ")
		b.WriteString(mode)
		b.WriteString("\n")
	}

	// Списки
	if lm := c.ListMatch; lm != nil {
		if lm.Extra > 0 || len(lm.Missing) > 0 {
			b.WriteString("Список: проверь комплектность и лишние элементы\n")
		} else if lm.Total > 0 {
			b.WriteString("Список: ок\n")
		}
	}

	// Шаги решения
	if st := c.StepsMatch; st != nil {
		if !st.OrderOK || len(st.Missing) > 0 || len(st.ExtraSteps) > 0 {
			b.WriteString("Шаги решения: проверь порядок и полноту\n")
		} else {
			b.WriteString("Шаги решения: ок\n")
		}
	}

	// 5) Короткая «озвучиваемая» фраза (до 140 символов)
	if s := strings.TrimSpace(cr.SpeakableMessage); s != "" {
		b.WriteString("\n")
		b.WriteString(s)
	}

	// 6) Рекомендованное следующее действие от модели
	if code := strings.TrimSpace(cr.NextActionCode); code != "" {
		var tip string
		switch code {
		case "ask_retry":
			tip = "→ Попробуй ещё раз: перепроверь и пришли новое фото решения."
		case "ask_rephoto":
			tip = "→ Пересними фото решения: чётко, без теней и бликов."
		case "ask_clarify_units":
			tip = "→ Уточни единицы измерения рядом с ответом."
		}
		if tip != "" {
			b.WriteString("\n")
			b.WriteString(tip)
		}
	}

	if strings.ToLower(strings.TrimSpace(cr.Verdict)) == "correct" {
		b.WriteString("\nДавай перейдём к решению следующе задачи.")
		clearMode(chatID)
		r.clearSession(chatID)
	}

	r.send(chatID, b.String(), nil)

	// 7) При ошибке или неуверенности предлагаем «Похожее задание»
	if strings.EqualFold(cr.Verdict, "incorrect") || strings.EqualFold(cr.Verdict, "uncertain") {
		r.offerAnalogueButton(chatID)
	}
}
