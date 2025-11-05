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
func (r *Router) lastParseMeta(ctx context.Context, chatID int64) (subject string, taskType string, grade int, ctxParse json.RawMessage) {
	pt, err := r.ParseRepo.FindByChatID(ctx, chatID)
	r.sendDebug(chatID, "err", err)
	r.sendDebug(chatID, "pt", pt)
	if err == nil && pt.Accepted {
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

	_, _, _, parseCtx := r.lastParseMeta(ctx, chatID)

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
	sid, _ := r.getSession(chatID)
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

// suggestSolutionShape — простая эвристика: если по парсингу известна форма — берём её, иначе number
func (r *Router) suggestSolutionShape(chatID int64) string {
	// Попробуем вывести форму ответа на основе последнего подтверждённого парсинга.
	// Если данных нет — вернём дефолт: number.
	if r.ParseRepo != nil {
		if pr, ok := r.ParseRepo.FindLastConfirmed(context.Background(), chatID); ok {
			util.PrintInfo("suggestSolutionShape", r.LlmManager.Get(chatID), chatID, fmt.Sprintf("parsed_raw: %+v", pr))

			subj := strings.ToLower(strings.TrimSpace(pr.Subject))
			tt := strings.ToLower(strings.TrimSpace(pr.TaskType))

			// Простая эвристика по предмету/типу задания
			// Русский язык — чаще всего ожидаем строку (слово/фразу)
			if subj == "russian" || subj == "ru" || subj == "русский" {
				if strings.Contains(tt, "list") || strings.Contains(tt, "спис") {
					return "list"
				}
				if strings.Contains(tt, "steps") || strings.Contains(tt, "шаг") {
					return "steps"
				}
				return "string"
			}

			// Математика/прочее
			if strings.Contains(tt, "list") || strings.Contains(tt, "спис") || strings.Contains(tt, "перечис") {
				return "list"
			}
			if strings.Contains(tt, "steps") || strings.Contains(tt, "шаг") || strings.Contains(tt, "пошаг") {
				return "steps"
			}
			if strings.Contains(tt, "word") || strings.Contains(tt, "слово") {
				return "string"
			}
			return "number"
		}
	}
	return "number"
}

// sendNormalizePreview — короткий текст для пользователя по NormalizeResult
func (r *Router) sendNormalizePreview(chatID int64) {
	b := &strings.Builder{}
	b.WriteString("✅ Принял ответ.")
	r.send(chatID, b.String(), nil)
}
