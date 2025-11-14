package telegram

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

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
				"subject_hint": dres.SubjectHint,
				"grade_hint":   dres.GradeHint,
				"confidence":   dres.Confidence,
				"debug_reason": dres.DebugReason,
			},
		})
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
	r.runParseAndMaybeConfirm(ctx, chatID, userID, sc, dres.SubjectHint, dres.GradeHint)
	util.PrintInfo("runDetectThenParse", llmName, chatID, fmt.Sprintf("Total time: %d", time.Since(start).Milliseconds()))
}
