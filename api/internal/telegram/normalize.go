package telegram

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"child-bot/api/internal/ocr/types"
	"child-bot/api/internal/store"
	"child-bot/api/internal/util"
)

// lastParseMeta — извлекает метаданные последнего подтверждённого парсинга
func (r *Router) lastParseMeta(chatID int64) (subject string, taskType string, grade int, ctx json.RawMessage) {
	if r.ParseRepo != nil {
		if pr, ok := r.ParseRepo.FindLastConfirmed(context.Background(), chatID); ok {
			subject = pr.Subject
			taskType = pr.TaskType
			grade = pr.Grade
			// Если вы храните сырое JSON парсинга — присвойте в ctx:
			// ctx = pr.RawJSON
		}
	}
	return
}

// normalizeText — отправляет текст ученика на нормализацию в LLM-прокси
func (r *Router) normalizeText(ctx context.Context, chatID int64, userID *int64, text string) {
	setState(chatID, Normalize)
	llmName := r.EngManager.Get(chatID)
	shape := r.suggestSolutionShape(chatID)

	subject, taskType, grade, parseCtx := r.lastParseMeta(chatID)

	var userIDAnon string
	if userID != nil {
		userIDAnon = fmt.Sprint(*userID)
	}

	sid, _ := r.getSession(chatID)

	in := types.NormalizeInput{
		TaskID:        sid,
		UserIDAnon:    userIDAnon,
		Grade:         grade,
		Subject:       subject,
		TaskType:      taskType,
		SolutionShape: shape,
		Answer:        types.NormalizeAnswer{Source: "text", Text: strings.TrimSpace(text)},
		ParseContext:  parseCtx,
		Provider:      llmName,
	}
	util.PrintInfo("normalizeText", r.EngManager.Get(chatID), chatID, fmt.Sprintf("normalize_input: %+v", in))
	r.sendDebug(chatID, "normalize_input", in)

	start := time.Now()
	res, err := r.LLM.Normalize(ctx, llmName, in)
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

		b := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report"))
		r.send(chatID, fmt.Sprintf("Не удалось нормализовать ответ: %v", err), b)
		return
	}

	_ = r.Metrics.InsertEvent(ctx, store.MetricEvent{
		Stage:      "normalize",
		Provider:   llmName,
		OK:         true,
		DurationMS: latency,
		ChatID:     &chatID,
		UserIDAnon: userID,
		Details: map[string]any{
			"source":          "text",
			"shape":           res.Shape,
			"needs_clarify":   res.NeedsClarification,
			"uncertain_count": len(res.UncertainReasons),
		},
	})

	r.sendNormalizePreview(chatID, res)

	// Попробуем сразу проверить решение, если в системе есть ожидаемое решение
	r.maybeCheckSolution(ctx, chatID, userID, res)
	clearMode(chatID)
}

// normalizePhoto — скачивает фото из Telegram и отправляет на нормализацию
func (r *Router) normalizePhoto(ctx context.Context, msg tgbotapi.Message) {
	llmName := r.EngManager.Get(util.GetChatIDFromTgMessage(msg))
	chatID := util.GetChatIDFromTgMessage(msg)
	subject, taskType, grade, parseCtx := r.lastParseMeta(chatID)

	if len(msg.Photo) == 0 {
		util.PrintInfo("normalizePhoto", llmName, chatID, "not found photo")
		return
	}

	ph := msg.Photo[len(msg.Photo)-1] // последнее
	data, mime, err := r.downloadFileBytes(ph.FileID)
	if err != nil {
		util.PrintError("normalizePhoto", llmName, chatID, "не удалось получить фото", err)
		b := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report"))
		r.send(chatID, fmt.Sprintf("Не удалось получить фото: %v", err), b)
		return
	}
	if mime == "application/octet-stream" {
		// Попробуем руками распознать распространённые форматы и HEIC/AVIF
		if len(data) >= 2 && data[0] == 0xFF && data[1] == 0xD8 {
			mime = "image/jpeg"
		}
		if len(data) >= 8 &&
			data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 &&
			data[4] == 0x0D && data[5] == 0x0A && data[6] == 0x1A && data[7] == 0x0A {
			mime = "image/png"
		}
		if heicAvif := util.SniffHEICorAVIF(data); heicAvif != "" {
			mime = heicAvif
		}
	}

	shape := r.suggestSolutionShape(chatID)
	sid, _ := r.getSession(chatID)

	in := types.NormalizeInput{
		TaskID: sid,
		UserIDAnon: func() string {
			u := util.GetUserIDFromTgMessage(msg)
			if u != nil {
				return fmt.Sprint(*u)
			}
			return ""
		}(),
		Grade:         grade,
		Subject:       subject,
		TaskType:      taskType,
		SolutionShape: shape,
		Answer: types.NormalizeAnswer{
			Source:   "photo",
			PhotoB64: base64.StdEncoding.EncodeToString(data),
			Mime:     mime,
		},
		ParseContext: parseCtx,
		Provider:     llmName,
	}
	// util.PrintInfo("normalizePhoto", llmName, chatID, fmt.Sprintf("normalize_input: %v", in))
	userID := util.GetUserIDFromTgMessage(msg)
	start := time.Now()
	res, err := r.LLM.Normalize(ctx, llmName, in)
	latency := time.Since(start).Milliseconds()
	_ = r.History.Insert(ctx, store.TimelineEvent{
		ChatID:        chatID,
		TaskSessionID: sid,
		Direction:     "api",
		EventType:     string(Normalize),
		Provider:      llmName,
		OK:            err == nil,
		LatencyMS:     &latency,
		TgMessageID:   &msg.MessageID,
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
			TaskID:     sid,
			Details: map[string]any{
				"source": "photo",
				"mime":   mime,
				"bytes":  len(data),
			},
		})

		util.PrintError("normalizePhoto", llmName, chatID, "Не удалось нормализовать ответ (фото)", err)
		b := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report"))
		r.send(chatID, "Не удалось нормализовать ответ (фото)", b)
		return
	}

	_ = r.Metrics.InsertEvent(ctx, store.MetricEvent{
		Stage:      "normalize",
		Provider:   llmName,
		OK:         true,
		DurationMS: latency,
		ChatID:     &chatID,
		UserIDAnon: userID,
		TaskID:     sid,
		Details: map[string]any{
			"source":          "photo",
			"mime":            mime,
			"bytes":           len(data),
			"shape":           res.Shape,
			"shape_detected":  res.ShapeDetected,
			"needs_clarify":   res.NeedsClarification,
			"uncertain_count": len(res.UncertainReasons),
		},
	})

	util.PrintInfo("normalizePhoto", llmName, chatID, fmt.Sprintf("normalize_photo: %+v", res))
	r.sendNormalizePreview(chatID, res)
	if res.Success {
		// Попробуем сразу проверить решение, если в системе есть ожидаемое решение
		r.maybeCheckSolution(ctx, chatID, userID, res)
	}
}

// suggestSolutionShape — простая эвристика: если по парсингу известна форма — берём её, иначе number
func (r *Router) suggestSolutionShape(chatID int64) string {
	// Попробуем вывести форму ответа на основе последнего подтверждённого парсинга.
	// Если данных нет — вернём дефолт: number.
	if r.ParseRepo != nil {
		if pr, ok := r.ParseRepo.FindLastConfirmed(context.Background(), chatID); ok {
			util.PrintInfo("suggestSolutionShape", r.EngManager.Get(chatID), chatID, fmt.Sprintf("parsed_raw: %+v", pr))

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
func (r *Router) sendNormalizePreview(chatID int64, nr types.NormalizeResult) {
	shape := strings.ToLower(strings.TrimSpace(nr.Shape))
	val := ""
	util.PrintInfo("sendNormalizePreview", r.EngManager.Get(chatID), chatID, fmt.Sprintf("normalize_result_value: %v", nr.Value))
	switch v := nr.Value.(type) {
	case string:
		val = v
	case float64:
		val = strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		val = strconv.Itoa(v)
	case int32:
		val = strconv.FormatInt(int64(v), 10)
	case int64:
		val = strconv.FormatInt(v, 10)
	case []string:
		val = strings.Join(v, "; ")
	default:
		val = "(не удалось отобразить значение)"
	}
	b := &strings.Builder{}
	b.WriteString("✅ Принял ответ.")
	if shape != "" {
		b.WriteString("\nФорма: ")
		b.WriteString(shape)
	}
	if val != "" {
		b.WriteString("\nЗначение: ")
		b.WriteString(val)
	}
	if nr.UncertainReasons != nil && len(nr.UncertainReasons) > 0 {
		b.WriteString("\nПредупреждения: ")
		b.WriteString(strings.Join(nr.UncertainReasons, ", "))
	}
	if nr.NeedsClarification != nil && *nr.NeedsClarification &&
		nr.NeedsUserActionMessage != nil && *nr.NeedsUserActionMessage != "" {
		b.WriteString("\nНужно уточнение: ")
		b.WriteString(*nr.NeedsUserActionMessage)
	}
	r.send(chatID, b.String(), nil)
}

// downloadFileBytes — скачивает файл Telegram по fileID и возвращает bytes и mime
func (r *Router) downloadFileBytes(fileID string) ([]byte, string, error) {
	url, err := r.Bot.GetFileDirectURL(fileID)
	if err != nil {
		return nil, "", err
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	mime := resp.Header.Get("Content-Type")
	if mime == "" {
		mime = "image/jpeg"
	}
	return b, mime, nil
}
