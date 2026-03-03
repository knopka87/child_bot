package telegram

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"child-bot/api/internal/store"
	"child-bot/api/internal/util"
	"child-bot/api/internal/v2/types"
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
	r.setStateWithPersist(chatID, Detect)
	mime := util.SniffMimeHTTP(image)
	llmName := r.LlmManager.Get(chatID)

	// DETECT через llmproxy
	var dres types.DetectResponse
	in := types.DetectRequest{
		Image:  base64.StdEncoding.EncodeToString(image),
		Locale: "ru_RU",
	}
	start := time.Now()
	dr, err := r.GetLLMClient().Detect(ctx, llmName, in)
	latency := time.Since(start).Milliseconds()
	sid, _ := r.getSession(chatID)
	_ = r.Store.InsertHistory(ctx, store.TimelineEvent{
		ChatID:        chatID,
		TaskSessionID: sid,
		Provider:      llmName,
		Direction:     "api",
		EventType:     string(Detect),
		InputPayload:  in,
		OutputPayload: dr,
		Error:         err,
		OK:            err == nil,
		LatencyMS:     &latency,
	})
	if err == nil {
		dres = dr
		r.sendDebug(chatID, "detect_req", in)
		r.sendDebug(chatID, "detect_res", dres)

		_ = r.Store.InsertEvent(ctx, store.MetricEvent{
			Stage:      "detect",
			Provider:   llmName,
			OK:         true,
			DurationMS: latency,
			ChatID:     &chatID,
			UserIDAnon: userID,
			Details: map[string]any{
				"subject_candidate": dres.Classification.SubjectCandidate,
				"confidence":        dres.Classification.Confidence,
				"recommend_retake":  dres.Quality.RecommendRetake,
			},
		})

		// Проверка предмета на этапе DETECT (audit 3.1)
		// Если предмет определён как НЕ математика с высокой уверенностью — останавливаемся
		const highConfidenceThreshold = 0.7
		if dres.Classification.SubjectCandidate != types.SubjectMath &&
			dres.Classification.Confidence >= highConfidenceThreshold {
			util.PrintInfo("runDetectThenParse", llmName, chatID,
				fmt.Sprintf("Subject not supported: %s (confidence=%.2f) - stopping",
					dres.Classification.SubjectCandidate, dres.Classification.Confidence))
			r.send(chatID, SubjectNotSupportedText, makeErrorButtons())
			r.setStateWithPersist(chatID, AwaitingTask)
			return
		}
	} else {
		_ = r.Store.InsertEvent(ctx, store.MetricEvent{
			Stage:      "detect",
			Provider:   llmName,
			OK:         false,
			DurationMS: latency,
			ChatID:     &chatID,
			UserIDAnon: userID,
			Error:      err.Error(),
		})

		r._sendWithError(chatID, DetectErrorText, "", makeErrorButtons(), fmt.Errorf("detect failed (chat=%d): %v; fallback to parse without detect", chatID, err))
	}

	// без выбора — сразу PARSE
	r.send(chatID, ReadTaskText, nil)
	sc := &selectionContext{Image: image, Mime: mime, MediaGroupID: mediaGroupID, Detect: dres}
	r.runParseAndMaybeConfirm(ctx, chatID, userID, sc, dres.Classification.SubjectCandidate)
	util.PrintInfo("runDetectThenParse", llmName, chatID, fmt.Sprintf("Total time: %d", time.Since(start).Milliseconds()))
}
