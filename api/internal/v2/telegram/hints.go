package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"child-bot/api/internal/store"
	"child-bot/api/internal/v2/types"
)

type hintSession struct {
	Image        []byte
	Mime         string
	MediaGroupID string
	Parse        types.ParseResponse
	Detect       types.DetectResponse
	EngineName   string
	NextLevel    int
}

func (r *Router) sendHint(ctx context.Context, chatID int64, msgID int, hs *hintSession) {
	level := hs.NextLevel
	sid, _ := r.getSession(chatID)

	grade := hs.Detect.GradeHint
	if user, err := r.Store.FindUserByChatID(ctx, chatID); err == nil && user.Grade != nil {
		grade = user.Grade
	}

	in := types.HintRequest{
		RawTaskText: hs.Parse.RawTaskText,
		Level:       lvlToConst(level),
		Grade:       grade,
		TaskStruct:  hs.Parse.TaskStruct,
		Locale:      "ru_RU",
	}
	hintLevel := level - 1
	for hintLevel > 0 {
		h, err := r.Store.FindHintBySID(ctx, sid, hintLevel)
		if err == nil {
			var hr types.HintResponse
			_ = json.Unmarshal(h.HintJson, &hr)
			in.PreviousHints = append(in.PreviousHints, hr.HintText)
		}
		hintLevel--
	}
	llmName := r.LlmManager.Get(chatID)
	start := time.Now()
	hrNew, err := r.GetLLMClient().Hint(context.Background(), llmName, in)
	latency := time.Since(start).Milliseconds()

	_ = r.Store.InsertHistory(context.Background(), store.TimelineEvent{
		ChatID:        chatID,
		TaskSessionID: sid,
		Direction:     "api",
		EventType:     string(Hints),
		Provider:      llmName,
		OK:            err == nil,
		LatencyMS:     &latency,
		TgMessageID:   &msgID,
		InputPayload:  in,
		OutputPayload: hrNew,
		Error:         err,
	})
	if err != nil {
		r.sendError(chatID, fmt.Errorf("не удалось получить подсказку L%d: %s", level, err.Error()))
		return
	}
	js, _ := json.Marshal(hrNew)
	data := store.HintCache{
		SessionID: sid,
		CreatedAt: time.Now(),
		Engine:    llmName,
		HintJson:  js,
		Level:     string(lvlToConst(level)),
	}
	err = r.Store.UpsertHint(context.Background(), data)
	if err != nil {
		_ = r.Store.InsertHistory(context.Background(), store.TimelineEvent{
			ChatID:        chatID,
			TaskSessionID: sid,
			Direction:     "db",
			EventType:     string(Hints),
			Provider:      llmName,
			OK:            false,
			Error:         err,
			CreatedAt:     time.Time{},
		})
	}
	r.send(chatID, formatHint(hrNew), makeHintButtons(level, true))

}

func formatHint(hr types.HintResponse) string {
	switch hr.Level {
	case types.HintL1:
		return fmt.Sprintf(HINT1Text, hr.HintText)
	case types.HintL2:
		return fmt.Sprintf(HINT2Text, hr.HintText)
	case types.HintL3:
		return fmt.Sprintf(HINT3Text, hr.HintText)
	default:
		return ""
	}
}
