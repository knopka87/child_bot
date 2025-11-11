package telegram

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"child-bot/api/internal/store"
	"child-bot/api/internal/util"
	"child-bot/api/internal/v1/types"
)

// OCR — скачивает фото из Telegram и отправляет на нормализацию
func (r *Router) OCR(ctx context.Context, msg tgbotapi.Message) {
	llmName := r.LlmManager.Get(util.GetChatIDFromTgMessage(msg))
	chatID := util.GetChatIDFromTgMessage(msg)

	if len(msg.Photo) == 0 {
		util.PrintInfo("OCR", llmName, chatID, "not found photo")
		return
	}

	ph := msg.Photo[len(msg.Photo)-1] // последнее
	data, mime, err := r.downloadFileBytes(ph.FileID)
	if err != nil {
		util.PrintError("OCR", llmName, chatID, "не удалось получить фото", err)
		b := make([][]tgbotapi.InlineKeyboardButton, 0, 1)
		b = append(b,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📝 Сообщить об ошибке", "report")),
		)
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

	sid, _ := r.getSession(chatID)

	in := types.OCRRequest{
		Image:  base64.StdEncoding.EncodeToString(data),
		Locale: "ru_RU",
	}
	// util.PrintInfo("OCR", llmName, chatID, fmt.Sprintf("ocr_input: %v", in))
	userID := util.GetUserIDFromTgMessage(msg)
	start := time.Now()
	res, err := r.GetLLMClient().OCR(ctx, llmName, in)
	latency := time.Since(start).Milliseconds()
	_ = r.Store.InsertHistory(ctx, store.TimelineEvent{
		ChatID:        chatID,
		TaskSessionID: sid,
		Direction:     "api",
		EventType:     string(OCR),
		Provider:      llmName,
		OK:            err == nil,
		LatencyMS:     &latency,
		TgMessageID:   &msg.MessageID,
		InputPayload:  in,
		OutputPayload: res,
		Error:         err,
	})
	if err != nil {
		_ = r.Store.InsertEvent(ctx, store.MetricEvent{
			Stage:      "ocr",
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

		util.PrintError("OCR", llmName, chatID, "Не удалось нормализовать ответ (фото)", err)
		b := make([][]tgbotapi.InlineKeyboardButton, 0, 1)
		b = append(b,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📝 Сообщить об ошибке", "report")),
		)
		r.send(chatID, "Не удалось нормализовать ответ (фото)", b)
		return
	}

	_ = r.Store.InsertEvent(ctx, store.MetricEvent{
		Stage:      "ocr",
		Provider:   llmName,
		OK:         true,
		DurationMS: latency,
		ChatID:     &chatID,
		UserIDAnon: userID,
		TaskID:     sid,
		Details: map[string]any{
			"source":     "photo",
			"mime":       mime,
			"bytes":      len(data),
			"confidence": res.Confidence,
		},
	})

	util.PrintInfo("OCR", llmName, chatID, fmt.Sprintf("ocr_photo: %+v", res))
	r.normalizeText(ctx, chatID, userID, res.RawAnswerText)
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
