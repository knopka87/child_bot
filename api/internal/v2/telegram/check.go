package telegram

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"child-bot/api/internal/store"
	"child-bot/api/internal/util"
	"child-bot/api/internal/v2/types"
)

// checkSolution — если есть ожидаемое решение для текущей задачи, проверяем ответ
func (r *Router) checkSolution(ctx context.Context, chatID int64, userID *int64, msg tgbotapi.Message) {
	llmName := r.LlmManager.Get(util.GetChatIDFromTgMessage(msg))

	if len(msg.Photo) == 0 {
		util.PrintInfo("Check", llmName, chatID, "not found photo")
		return
	}

	setState(chatID, Check)
	sid, _ := r.getSession(chatID)

	time1 := r.sendAlert(chatID, CheckAlert, 0, 15)

	ph := msg.Photo[len(msg.Photo)-1] // последнее
	data, mime, err := r.downloadFileBytes(ph.FileID)
	if err != nil {
		r.sendError(chatID, fmt.Errorf("не удалось получить фото: %v", err))
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

	// 0) Подтянем метаданные предмета/класса из последнего подтверждённого парсинга
	subj := "generic"
	grade := int64(0)
	rawTaskText := ""
	var taskStruct types.TaskStructCheck
	if pr, ok := r.Store.FindLastConfirmedParse(ctx, sid); ok {
		if s := strings.TrimSpace(pr.Subject); s != "" {
			subj = s
			grade = pr.Grade
			rawTaskText = pr.RawTaskText

			// Восстанавливаем ParseResponse из сохранённого JSON
			var parseResp types.ParseResponse
			if err := json.Unmarshal(pr.ResultJSON, &parseResp); err == nil {
				taskStruct = types.TaskStructCheck{
					TaskTextClean: parseResp.Task.TaskTextClean,
					VisualFacts:   parseResp.Task.VisualFacts,
					QualityFlags:  parseResp.Task.Quality,
					Items:         parseResp.Items,
				}
			}
		}
	}

	in := types.CheckRequest{
		Image:            base64.StdEncoding.EncodeToString(data),
		TaskStruct:       taskStruct,
		RawTaskText:      rawTaskText,
		PhotoQualityHint: "auto",
		Student: types.StudentCheck{
			Grade:   grade,
			Subject: subj,
			Locale:  "ru_RU",
		},
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

	if cr.IsCorrect != nil && *cr.IsCorrect {
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
