package handler

import (
	"encoding/json"
	"net/http"

	"child-bot/api/internal/api/response"
)

// AnalyticsHandler обрабатывает analytics events
type AnalyticsHandler struct{}

// NewAnalyticsHandler создает новый AnalyticsHandler
func NewAnalyticsHandler() *AnalyticsHandler {
	return &AnalyticsHandler{}
}

// AnalyticsEvent представляет событие аналитики
type AnalyticsEvent struct {
	Name      string                 `json:"name"`
	Timestamp int64                  `json:"timestamp"`
	SessionID string                 `json:"sessionId"` // camelCase для соответствия frontend
	Params    map[string]interface{} `json:"params"`
}

// AnalyticsEventsRequest обертка для batch событий
type AnalyticsEventsRequest struct {
	Events []AnalyticsEvent `json:"events"`
}

// SendEvents обрабатывает POST /analytics/events
// Принимает batch событий от frontend и отправляет их в платформенную аналитику
func (h *AnalyticsHandler) SendEvents(w http.ResponseWriter, r *http.Request) {
	var req AnalyticsEventsRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body: "+err.Error())
		return
	}

	// TODO: Реализовать отправку событий в VK/MAX analytics
	// Сейчас просто логируем и возвращаем успех
	// В production здесь будет:
	// - Валидация событий
	// - Обогащение метаданными (platform, user_id, etc.)
	// - Отправка в VK Mini Apps Analytics / MAX Analytics
	// - Сохранение в БД для отчетов (опционально)

	// Логируем количество событий
	// log.Printf("[Analytics] Received %d events", len(req.Events))

	response.OK(w, map[string]interface{}{
		"success":      true,
		"events_count": len(req.Events),
		"message":      "Events received successfully",
	})
}
