package telegram

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"child-bot/api/internal/store"
	"child-bot/api/internal/util"
	"child-bot/api/internal/v1/types"
)

// TaskChoice хранит отображаемый номер варианта и краткое описание + индекс задачи
type TaskChoice struct {
	Number      string
	Description string
	TaskIndex   int
}

func (r *Router) hasPendingCorrection(chatID int64) bool { _, ok := parseWait.Load(chatID); return ok }
func (r *Router) clearPendingCorrection(chatID int64)    { parseWait.Delete(chatID) }

func (r *Router) runDetectThenParse(ctx context.Context, chatID int64, userID *int64, image []byte, mediaGroupID string) {
	setState(chatID, Detect)
	mime := util.SniffMimeHTTP(image)
	llmName := r.LlmManager.Get(chatID)

	// DETECT через llmproxy
	var dres types.DetectResponse
	in := types.DetectRequest{
		Image:    base64.StdEncoding.EncodeToString(image),
		Locale:   "ru-RU",
		MaxTasks: 1,
	}
	start := time.Now()
	dr, err := r.GetLLMClient().Detect(ctx, llmName, in)
	latency := time.Since(start).Milliseconds()
	if err == nil {
		dres = dr
		r.sendDebug(chatID, "detect_req", in)
		r.sendDebug(chatID, "detect_res", dres)

		errM := r.Metrics.InsertEvent(ctx, store.MetricEvent{
			Stage:      "detect",
			Provider:   llmName,
			OK:         true,
			DurationMS: latency,
			ChatID:     &chatID,
			UserIDAnon: userID,
			Details: map[string]any{
				"subject_hint": dres.SubjectHint,
				"grade_hint":   dres.GradeHint,
				"confidence":   dres.Confidence,
				"debug_reason": dres.DebugReason,
			},
		})
		if errM != nil {
			util.PrintError("runDetectThenParse", llmName, chatID, "error insert metrics", errM)
		}
	} else {
		if r.Metrics != nil {
			_ = r.Metrics.InsertEvent(ctx, store.MetricEvent{
				Stage:      "detect",
				Provider:   llmName,
				OK:         false,
				DurationMS: latency,
				ChatID:     &chatID,
				UserIDAnon: userID,
				Error:      err.Error(),
			})
		}
		log.Printf("detect failed (chat=%d): %v; fallback to parse without detect", chatID, err)
		b := make([][]tgbotapi.InlineKeyboardButton, 0, 1)
		b = append(b, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report")))
		r.send(chatID, "ℹ️ Не удалось выделить области на фото, попробую распознать задание целиком.", b)
	}
	util.PrintInfo("runDetectThenParse", llmName, chatID, fmt.Sprintf("Received a response from LLMClient: %d", time.Since(start).Milliseconds()))

	sid, _ := r.getSession(chatID)
	_ = r.History.Insert(ctx, store.TimelineEvent{
		ChatID:        chatID,
		TaskSessionID: sid,
		Provider:      llmName,
		Direction:     "api",
		EventType:     string(Detect),
		InputPayload:  in,
		OutputPayload: dres,
		Error:         err,
		OK:            err == nil,
		LatencyMS:     &latency,
	})

	// без выбора — сразу PARSE
	setState(chatID, DecideTasks)
	r.send(chatID, "Изображение распознано, перехожу к парсингу.", nil)
	sc := &selectionContext{Image: image, Mime: mime, MediaGroupID: mediaGroupID, Detect: dres}
	r.runParseAndMaybeConfirm(ctx, chatID, userID, sc, dres.SubjectHint, dres.GradeHint)
	util.PrintInfo("runDetectThenParse", llmName, chatID, fmt.Sprintf("Total time: %d", time.Since(start).Milliseconds()))
}
