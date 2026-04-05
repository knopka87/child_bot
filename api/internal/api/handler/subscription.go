package handler

import (
	"net/http"

	"child-bot/api/internal/api/middleware"
	"child-bot/api/internal/api/response"
	"child-bot/api/internal/api/validation"
	"child-bot/api/internal/store"
)

// SubscriptionHandler обрабатывает запросы подписок
type SubscriptionHandler struct {
	store *store.Store
}

// NewSubscriptionHandler создает новый SubscriptionHandler
func NewSubscriptionHandler(store *store.Store) *SubscriptionHandler {
	return &SubscriptionHandler{store: store}
}

// SubscriptionStatus структура статуса подписки
type SubscriptionStatus struct {
	Status             string   `json:"status"` // trial, active, expired, cancelled
	PlanID             string   `json:"plan_id,omitempty"`
	PlanName           string   `json:"plan_name,omitempty"`
	Features           []string `json:"features"`
	TrialDaysRemaining int      `json:"trial_days_remaining,omitempty"`
	ExpiresAt          string   `json:"expires_at,omitempty"`
	RenewsAt           string   `json:"renews_at,omitempty"`
	CancelledAt        string   `json:"cancelled_at,omitempty"`
	CanCancel          bool     `json:"can_cancel"`
	CanResume          bool     `json:"can_resume"`
}

// SubscriptionPlan план подписки
type SubscriptionPlan struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Price           int      `json:"price"` // в копейках
	Currency        string   `json:"currency"`
	Duration        string   `json:"duration"` // month, year
	Features        []string `json:"features"`
	IsPopular       bool     `json:"is_popular"`
	TrialDays       int      `json:"trial_days"`
	DiscountPercent int      `json:"discount_percent,omitempty"`
}

// SubscribeRequest запрос на подписку
type SubscribeRequest struct {
	PlanID        string `json:"plan_id"`
	PaymentMethod string `json:"payment_method"` // card, yookassa, etc.
}

// GetStatus получает статус подписки
// GET /subscription/status
func (h *SubscriptionHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	childProfileID := middleware.GetChildProfileID(r.Context())
	if childProfileID == "" {
		response.Unauthorized(w, "Missing child_profile_id")
		return
	}

	// TODO: Phase 4 - получение статуса через service layer
	// status, err := h.service.GetSubscriptionStatus(r.Context(), childProfileID)

	// Placeholder
	status := SubscriptionStatus{
		Status:             "trial",
		Features:           []string{"unlimited_tasks", "hints", "achievements"},
		TrialDaysRemaining: 7,
		CanCancel:          false,
		CanResume:          false,
	}

	response.OK(w, status)
}

// GetPlans получает список доступных планов
// GET /subscription/plans
func (h *SubscriptionHandler) GetPlans(w http.ResponseWriter, r *http.Request) {
	// TODO: Phase 4 - получение планов через service layer
	// plans, err := h.service.GetSubscriptionPlans(r.Context())

	// Placeholder
	plans := []SubscriptionPlan{
		{
			ID:          "plan_monthly",
			Name:        "Месячная подписка",
			Description: "Полный доступ ко всем функциям на 1 месяц",
			Price:       49900, // 499 руб
			Currency:    "RUB",
			Duration:    "month",
			Features: []string{
				"Неограниченное количество задач",
				"Умные подсказки",
				"Проверка решений",
				"Достижения и награды",
			},
			IsPopular: true,
			TrialDays: 7,
		},
		{
			ID:          "plan_yearly",
			Name:        "Годовая подписка",
			Description: "Выгодная подписка на целый год",
			Price:       399900, // 3999 руб
			Currency:    "RUB",
			Duration:    "year",
			Features: []string{
				"Неограниченное количество задач",
				"Умные подсказки",
				"Проверка решений",
				"Достижения и награды",
				"Приоритетная поддержка",
			},
			IsPopular:       false,
			TrialDays:       14,
			DiscountPercent: 33,
		},
	}

	response.OK(w, plans)
}

// Subscribe оформляет подписку
// POST /subscription/subscribe
func (h *SubscriptionHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	childProfileID := middleware.GetChildProfileID(r.Context())
	if childProfileID == "" {
		response.Unauthorized(w, "Missing child_profile_id")
		return
	}

	var req SubscribeRequest
	if err := validation.DecodeJSON(r, &req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	// Валидация
	if err := validation.ValidateRequired(req.PlanID, "plan_id"); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	if err := validation.ValidateRequired(req.PaymentMethod, "payment_method"); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	// TODO: Phase 4 - оформление подписки через service layer
	// result, err := h.service.Subscribe(r.Context(), childProfileID, req)

	// Placeholder (возвращаем URL для оплаты)
	result := map[string]interface{}{
		"payment_url": "https://payment.example.com/checkout/123",
		"payment_id":  "payment_123",
		"status":      "pending",
		"expires_at":  "2024-03-31T11:00:00Z",
	}

	response.OK(w, result)
}

// Cancel отменяет подписку
// DELETE /subscription/cancel
func (h *SubscriptionHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	childProfileID := middleware.GetChildProfileID(r.Context())
	if childProfileID == "" {
		response.Unauthorized(w, "Missing child_profile_id")
		return
	}

	// TODO: Phase 4 - отмена подписки через service layer
	// err := h.service.CancelSubscription(r.Context(), childProfileID)

	// Placeholder
	result := map[string]interface{}{
		"status":       "cancelled",
		"cancelled_at": "2024-03-31T10:00:00Z",
		"expires_at":   "2024-04-30T23:59:59Z",
		"message":      "Подписка отменена. Доступ сохраняется до конца оплаченного периода.",
	}

	response.OK(w, result)
}

// Resume возобновляет подписку
// POST /subscription/resume
func (h *SubscriptionHandler) Resume(w http.ResponseWriter, r *http.Request) {
	childProfileID := middleware.GetChildProfileID(r.Context())
	if childProfileID == "" {
		response.Unauthorized(w, "Missing child_profile_id")
		return
	}

	// TODO: Phase 4 - возобновление подписки через service layer
	// err := h.service.ResumeSubscription(r.Context(), childProfileID)

	// Placeholder
	result := map[string]interface{}{
		"status":    "active",
		"renews_at": "2024-04-30T23:59:59Z",
		"message":   "Подписка возобновлена",
	}

	response.OK(w, result)
}
