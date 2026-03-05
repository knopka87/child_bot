package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"child-bot/api/internal/store"
	"child-bot/api/internal/v2/types"
)

type hintSession struct {
	mu           sync.Mutex // защита от concurrent access
	Image        []byte
	Mime         string
	MediaGroupID string
	Parse        types.ParseResponse
	Detect       types.DetectResponse
	EngineName   string
	NextLevel    int
	MaxHints     int
}

func (r *Router) sendHint(_ context.Context, chatID int64, msgID int, hs *hintSession) {
	// Копируем данные под защитой мьютекса чтобы избежать race condition
	hs.mu.Lock()
	level := hs.NextLevel
	parseData := hs.Parse
	detectData := hs.Detect
	hs.mu.Unlock()

	sid, _ := r.getSession(chatID)

	// Определяем режим подсказки
	mode := "learn"
	if level > 1 {
		mode = "rescue"
	}

	// Определяем политику подсказок из первого item, если есть
	var appliedPolicy types.HintPolicy
	if len(parseData.Items) > 0 {
		appliedPolicy = parseData.Items[0].HintPolicy
	} else {
		appliedPolicy = types.HintPolicy{
			MaxHints:       3,
			DefaultVisible: 1,
			H3Reason:       types.H3ReasonNone,
		}
	}

	in := types.HintRequest{
		Task:          parseData.Task,
		Mode:          mode,
		Items:         parseData.Items,
		AppliedPolicy: appliedPolicy,
		Template:      getTemplate(parseData.Task, parseData.Items, detectData.Classification.SubjectCandidate),
	}

	llmName := r.LlmManager.Get(chatID)
	stopProgress := r.startHintProgress(chatID)
	start := time.Now()
	hrNew, err := r.GetLLMClient().Hint(context.Background(), llmName, in)
	latency := time.Since(start).Milliseconds()
	stopProgress()

	// Получаем template_id и task_type для метрик
	templateID := getTemplateID(parseData.Task, parseData.Items, detectData.Classification.SubjectCandidate)
	taskType := ""
	if len(parseData.Items) > 0 {
		taskType = parseData.Items[0].PedKeys.TaskType
	}

	// Получаем класс пользователя
	var grade int64
	if user, userErr := r.Store.FindUserByChatID(context.Background(), chatID); userErr == nil && user.Grade != nil {
		grade = *user.Grade
	}

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
		_ = r.Store.InsertEvent(context.Background(), store.MetricEvent{
			Stage:      "hint",
			Provider:   llmName,
			OK:         false,
			Error:      err.Error(),
			DurationMS: latency,
			ChatID:     &chatID,
			Details: map[string]any{
				"hint_level":  level,
				"template_id": templateID,
				"task_type":   taskType,
				"grade":       grade,
			},
		})
		r.sendError(chatID, fmt.Errorf("не удалось получить подсказку L%d: %s", level, err.Error()))
		return
	}

	// Метрика успешной подсказки
	_ = r.Store.InsertEvent(context.Background(), store.MetricEvent{
		Stage:      "hint",
		Provider:   llmName,
		OK:         true,
		DurationMS: latency,
		ChatID:     &chatID,
		Details: map[string]any{
			"hint_level":  level,
			"max_hints":   appliedPolicy.MaxHints,
			"template_id": templateID,
			"task_type":   taskType,
			"grade":       grade,
			"mode":        mode,
		},
	})
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
	r.send(chatID, formatHint(hrNew, level), makeHintButtons(level, appliedPolicy.MaxHints, true))
}

func formatHint(hr types.HintResponse, level int) string {
	// Собираем подсказки из всех items для указанного уровня
	var hints []string
	targetLevel := lvlToConst(level)

	for _, item := range hr.Items {
		for _, hint := range item.Hints {
			if hint.Level == targetLevel {
				hints = append(hints, hint.HintText)
			}
		}
	}

	if len(hints) == 0 {
		return "Подсказка не найдена"
	}

	hintText := strings.Join(hints, "\n\n")

	switch targetLevel {
	case types.HintL1:
		return fmt.Sprintf(HINT1Text, hintText)
	case types.HintL2:
		return fmt.Sprintf(HINT2Text, hintText)
	case types.HintL3:
		return fmt.Sprintf(HINT3Text, hintText)
	default:
		return hintText
	}
}
