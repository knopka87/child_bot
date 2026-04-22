package handler

import (
	"log"
	"net/http"
	"time"

	"child-bot/api/internal/api/middleware"
	"child-bot/api/internal/api/response"
	"child-bot/api/internal/api/validation"
	"child-bot/api/internal/service"
	"child-bot/api/internal/store"
)

// SubscriptionHandler обрабатывает запросы подписок
type SubscriptionHandler struct {
	store        *store.Store
	vkPayService *service.VKPayService
}

// NewSubscriptionHandler создает новый SubscriptionHandler
func NewSubscriptionHandler(store *store.Store, vkPayService *service.VKPayService) *SubscriptionHandler {
	return &SubscriptionHandler{
		store:        store,
		vkPayService: vkPayService,
	}
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

	// Получаем активную подписку
	subscription, err := h.store.GetActiveSubscription(r.Context(), childProfileID)
	if err != nil && err.Error() != "subscription not found" {
		log.Printf("[SubscriptionHandler] Failed to get subscription: %v", err)
		response.InternalError(w, "Failed to get subscription status")
		return
	}

	// Базовые фичи для всех
	features := []string{"unlimited_tasks", "hints", "achievements"}

	// Если нет подписки - возвращаем trial статус
	if subscription == nil {
		status := SubscriptionStatus{
			Status:             "trial",
			Features:           features,
			TrialDaysRemaining: 7,
			CanCancel:          false,
			CanResume:          false,
		}
		response.OK(w, status)
		return
	}

	// Загружаем информацию о плане
	plan, err := h.store.GetSubscriptionPlan(r.Context(), subscription.PlanID)
	if err != nil {
		log.Printf("[SubscriptionHandler] Failed to get plan: %v", err)
	}

	// Формируем статус
	status := SubscriptionStatus{
		Status:    subscription.Status,
		PlanID:    subscription.PlanID,
		Features:  features,
		ExpiresAt: subscription.ExpiresAt.Format(time.RFC3339),
		CanCancel: subscription.Status == "active" && subscription.CancelledAt == nil,
		CanResume: subscription.CancelledAt != nil,
	}

	if plan != nil {
		status.PlanName = plan.Name
	}

	if subscription.CancelledAt != nil {
		status.CancelledAt = subscription.CancelledAt.Format(time.RFC3339)
	}

	if subscription.Status == "trial" && subscription.TrialEndsAt != nil {
		daysRemaining := int(time.Until(*subscription.TrialEndsAt).Hours() / 24)
		if daysRemaining < 0 {
			daysRemaining = 0
		}
		status.TrialDaysRemaining = daysRemaining
	}

	if subscription.Status == "active" && subscription.AutoRenew {
		status.RenewsAt = subscription.ExpiresAt.Format(time.RFC3339)
	}

	response.OK(w, status)
}

// GetPlans получает список доступных планов
// GET /subscription/plans
func (h *SubscriptionHandler) GetPlans(w http.ResponseWriter, r *http.Request) {
	// Загружаем активные планы из БД
	dbPlans, err := h.store.GetActivePlans(r.Context())
	if err != nil {
		log.Printf("[SubscriptionHandler] Failed to get plans: %v", err)
		response.InternalError(w, "Failed to get subscription plans")
		return
	}

	// Базовые features для всех планов
	baseFeatures := []string{
		"Неограниченное количество задач",
		"Умные подсказки",
		"Проверка решений",
		"Достижения и награды",
	}

	// Преобразуем в API формат
	plans := make([]SubscriptionPlan, 0, len(dbPlans))
	for _, dbPlan := range dbPlans {
		features := make([]string, len(baseFeatures))
		copy(features, baseFeatures)

		// Добавляем дополнительные features для годовой подписки
		if dbPlan.ID == "yearly" {
			features = append(features, "Приоритетная поддержка")
		}

		// Определяем duration string
		duration := "month"
		if dbPlan.DurationDays >= 365 {
			duration = "year"
		}

		plan := SubscriptionPlan{
			ID:              dbPlan.ID,
			Name:            dbPlan.Name,
			Description:     dbPlan.Description,
			Price:           dbPlan.PriceCents,
			Currency:        dbPlan.Currency,
			Duration:        duration,
			Features:        features,
			IsPopular:       dbPlan.IsPopular,
			TrialDays:       dbPlan.TrialDays,
			DiscountPercent: dbPlan.DiscountPercent,
		}
		plans = append(plans, plan)
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

	// Получаем VK User ID из контекста
	vkUserID := middleware.GetVKUserID(r.Context())
	if vkUserID == "" {
		response.Unauthorized(w, "Missing vk_user_id")
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

	// Поддерживаем только VK Pay
	if req.PaymentMethod != "vk_pay" {
		response.BadRequest(w, "Only vk_pay payment method is supported")
		return
	}

	// Создаем платеж через VK Pay Service
	ipAddress := r.RemoteAddr
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ipAddress = xff
	}

	userAgent := r.Header.Get("User-Agent")

	paymentReq := service.CreatePaymentRequest{
		ChildProfileID: childProfileID,
		PlanID:         req.PlanID,
		IPAddress:      ipAddress,
		UserAgent:      userAgent,
		VKUserID:       vkUserID,
	}

	paymentOrder, err := h.vkPayService.CreatePayment(r.Context(), paymentReq)
	if err != nil {
		log.Printf("[SubscriptionHandler] Failed to create payment: %v", err)
		response.InternalError(w, "Failed to create payment")
		return
	}

	// Возвращаем данные для открытия VK Pay формы
	result := map[string]interface{}{
		"payment_id": paymentOrder.PaymentID,
		"order_id":   paymentOrder.OrderID,
		"vk_pay_url": paymentOrder.VKPayURL,
		"amount":     paymentOrder.Amount,
		"currency":   paymentOrder.Currency,
		"status":     "pending",
		"expires_at": paymentOrder.ExpiresAt.Format(time.RFC3339),
		"metadata":   paymentOrder.Metadata,
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

	// Отменяем подписку
	err := h.store.CancelSubscription(r.Context(), childProfileID)
	if err != nil {
		if err.Error() == "no active subscription to cancel" {
			response.BadRequest(w, "No active subscription to cancel")
			return
		}
		log.Printf("[SubscriptionHandler] Failed to cancel subscription: %v", err)
		response.InternalError(w, "Failed to cancel subscription")
		return
	}

	// Получаем обновленную подписку
	subscription, err := h.store.GetActiveSubscription(r.Context(), childProfileID)
	if err != nil {
		log.Printf("[SubscriptionHandler] Failed to get subscription after cancel: %v", err)
		// Не критичная ошибка, можем вернуть базовый ответ
		result := map[string]interface{}{
			"status":  "cancelled",
			"message": "Подписка отменена. Доступ сохраняется до конца оплаченного периода.",
		}
		response.OK(w, result)
		return
	}

	result := map[string]interface{}{
		"status":       subscription.Status,
		"cancelled_at": subscription.CancelledAt.Format(time.RFC3339),
		"expires_at":   subscription.ExpiresAt.Format(time.RFC3339),
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

	// Возобновляем подписку
	err := h.store.ResumeSubscription(r.Context(), childProfileID)
	if err != nil {
		if err.Error() == "no cancelled subscription to resume" {
			response.BadRequest(w, "No cancelled subscription to resume")
			return
		}
		log.Printf("[SubscriptionHandler] Failed to resume subscription: %v", err)
		response.InternalError(w, "Failed to resume subscription")
		return
	}

	// Получаем обновленную подписку
	subscription, err := h.store.GetActiveSubscription(r.Context(), childProfileID)
	if err != nil {
		log.Printf("[SubscriptionHandler] Failed to get subscription after resume: %v", err)
		result := map[string]interface{}{
			"status":  "active",
			"message": "Подписка возобновлена",
		}
		response.OK(w, result)
		return
	}

	result := map[string]interface{}{
		"status":    subscription.Status,
		"renews_at": subscription.ExpiresAt.Format(time.RFC3339),
		"message":   "Подписка возобновлена. Автопродление включено.",
	}

	response.OK(w, result)
}
