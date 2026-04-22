package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sort"
	"strings"
	"time"

	"child-bot/api/internal/store"

	"github.com/google/uuid"
)

// VKPayService обрабатывает платежи через VK Pay
type VKPayService struct {
	store  *store.Store
	config VKPayConfig
}

// VKPayConfig конфигурация VK Pay
type VKPayConfig struct {
	AppID       string // VK App ID
	AppSecret   string // VK App Secret для валидации webhooks
	CallbackURL string // URL для webhooks
}

// Payment модель платежа
type Payment struct {
	ID              string
	SubscriptionID  *int64
	ChildProfileID  string
	PlanID          string
	AmountCents     int
	Currency        string
	VKOrderID       *string
	VKTransactionID *string
	VKUserID        *string
	Status          string
	PaymentMethod   string
	Description     string
	IPAddress       *string
	UserAgent       *string
	Metadata        map[string]interface{}
	PaidAt          *time.Time
	RefundedAt      *time.Time
	ExpiresAt       *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// PaymentEvent событие платежа
type PaymentEvent struct {
	PaymentID    string
	EventType    string
	OldStatus    *string
	NewStatus    *string
	VKEventType  *string
	VKEventData  map[string]interface{}
	ErrorCode    *string
	ErrorMessage *string
	CreatedAt    time.Time
}

// CreatePaymentRequest запрос на создание платежа
type CreatePaymentRequest struct {
	ChildProfileID string
	PlanID         string
	IPAddress      string
	UserAgent      string
	VKUserID       string
}

// PaymentOrderResponse данные для открытия платежной формы
type PaymentOrderResponse struct {
	PaymentID string                 `json:"payment_id"`
	OrderID   string                 `json:"order_id"`
	VKPayURL  string                 `json:"vk_pay_url"`
	Amount    int                    `json:"amount"`
	Currency  string                 `json:"currency"`
	ExpiresAt time.Time              `json:"expires_at"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// NewVKPayService создает новый VKPayService
func NewVKPayService(store *store.Store, config VKPayConfig) *VKPayService {
	return &VKPayService{
		store:  store,
		config: config,
	}
}

// CreatePayment создает новый платеж и возвращает данные для VK Pay формы
func (s *VKPayService) CreatePayment(ctx context.Context, req CreatePaymentRequest) (*PaymentOrderResponse, error) {
	// Загружаем план подписки
	plan, err := s.store.GetSubscriptionPlan(ctx, req.PlanID)
	if err != nil {
		return nil, fmt.Errorf("get subscription plan: %w", err)
	}

	if !plan.IsActive {
		return nil, fmt.Errorf("plan is not active")
	}

	// Генерируем уникальный order ID
	orderID := fmt.Sprintf("order_%s_%d", uuid.New().String()[:8], time.Now().Unix())

	// Создаем платеж в БД
	paymentID := uuid.New().String()
	expiresAt := time.Now().Add(30 * time.Minute) // Платеж действителен 30 минут

	query := `
		INSERT INTO payments (
			id, child_profile_id, plan_id, amount_cents, currency,
			vk_order_id, vk_user_id, status, payment_method,
			description, ip_address, user_agent, expires_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING created_at
	`

	description := fmt.Sprintf("Подписка: %s", plan.Name)

	var createdAt time.Time
	err = s.store.DB.QueryRowContext(ctx, query,
		paymentID, req.ChildProfileID, req.PlanID, plan.PriceCents, plan.Currency,
		orderID, req.VKUserID, "pending", "vk_pay",
		description, req.IPAddress, req.UserAgent, expiresAt,
	).Scan(&createdAt)

	if err != nil {
		return nil, fmt.Errorf("create payment: %w", err)
	}

	// Логируем событие создания
	if err := s.logEvent(ctx, paymentID, "created", nil, strPtr("pending"), nil); err != nil {
		log.Printf("[VKPayService] Failed to log payment created event: %v", err)
	}

	log.Printf("[VKPayService] Created payment %s for plan %s, amount: %d %s",
		paymentID, req.PlanID, plan.PriceCents, plan.Currency)

	// Формируем VK Pay URL
	vkPayURL := s.buildVKPayURL(orderID, plan.PriceCents, description, req.VKUserID)

	return &PaymentOrderResponse{
		PaymentID: paymentID,
		OrderID:   orderID,
		VKPayURL:  vkPayURL,
		Amount:    plan.PriceCents,
		Currency:  plan.Currency,
		ExpiresAt: expiresAt,
		Metadata: map[string]interface{}{
			"plan_name": plan.Name,
			"duration":  plan.DurationDays,
		},
	}, nil
}

// buildVKPayURL формирует URL для VK Pay
func (s *VKPayService) buildVKPayURL(orderID string, amountCents int, description, vkUserID string) string {
	// VK Pay использует копейки, но в параметре amount нужны рубли
	amountRubles := float64(amountCents) / 100.0

	params := url.Values{}
	params.Set("app_id", s.config.AppID)
	params.Set("order_id", orderID)
	params.Set("amount", fmt.Sprintf("%.2f", amountRubles))
	params.Set("description", description)
	params.Set("action", "pay-to-group") // Платеж в группу/приложение

	if s.config.CallbackURL != "" {
		params.Set("notification_url", s.config.CallbackURL)
	}

	// VK Bridge откроет форму оплаты
	return fmt.Sprintf("https://vk.com/app%s#order_id=%s", s.config.AppID, orderID)
}

// ProcessWebhook обрабатывает webhook от VK Pay
func (s *VKPayService) ProcessWebhook(ctx context.Context, vkUserID, orderID, notificationType string, payload map[string]interface{}) error {
	log.Printf("[VKPayService] Processing webhook: type=%s, order=%s, user=%s",
		notificationType, orderID, vkUserID)

	// Загружаем платеж по order_id
	payment, err := s.getPaymentByOrderID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("get payment by order_id: %w", err)
	}

	// Логируем webhook событие
	if err := s.logEvent(ctx, payment.ID, "webhook_received", nil, nil, &notificationType); err != nil {
		log.Printf("[VKPayService] Failed to log webhook event: %v", err)
	}

	// Обрабатываем разные типы уведомлений
	switch notificationType {
	case "order_status_change":
		return s.handleOrderStatusChange(ctx, payment, payload)
	case "order_status_change_test":
		// Тестовое уведомление - игнорируем
		log.Printf("[VKPayService] Received test webhook, ignoring")
		return nil
	default:
		log.Printf("[VKPayService] Unknown notification type: %s", notificationType)
		return fmt.Errorf("unknown notification type: %s", notificationType)
	}
}

// handleOrderStatusChange обрабатывает изменение статуса заказа
func (s *VKPayService) handleOrderStatusChange(ctx context.Context, payment *Payment, payload map[string]interface{}) error {
	// Извлекаем новый статус
	status, ok := payload["status"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid status in payload")
	}

	log.Printf("[VKPayService] Order %s status changed to: %s", *payment.VKOrderID, status)

	// Маппинг VK статусов на наши статусы
	var newStatus string
	switch status {
	case "chargeable":
		newStatus = "processing"
	case "charged":
		newStatus = "completed"
	case "refunded":
		newStatus = "refunded"
	case "declined":
		newStatus = "failed"
	case "cancelled":
		newStatus = "cancelled"
	default:
		log.Printf("[VKPayService] Unknown VK status: %s, setting to pending", status)
		newStatus = "pending"
	}

	// Обновляем платеж
	if err := s.updatePaymentStatus(ctx, payment.ID, newStatus, payload); err != nil {
		return fmt.Errorf("update payment status: %w", err)
	}

	// Если платеж успешен - создаем/активируем подписку
	if newStatus == "completed" {
		if err := s.activateSubscription(ctx, payment); err != nil {
			log.Printf("[VKPayService] Failed to activate subscription: %v", err)
			// Не возвращаем ошибку, чтобы не блокировать webhook
		}
	}

	return nil
}

// updatePaymentStatus обновляет статус платежа
func (s *VKPayService) updatePaymentStatus(ctx context.Context, paymentID, newStatus string, vkData map[string]interface{}) error {
	tx, err := s.store.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Получаем текущий статус
	var oldStatus string
	err = tx.QueryRowContext(ctx, `SELECT status FROM payments WHERE id = $1`, paymentID).Scan(&oldStatus)
	if err != nil {
		return fmt.Errorf("get current status: %w", err)
	}

	// Обновляем платеж
	query := `
		UPDATE payments
		SET status = $1,
		    paid_at = CASE WHEN $1 = 'completed' AND paid_at IS NULL THEN NOW() ELSE paid_at END,
		    refunded_at = CASE WHEN $1 = 'refunded' AND refunded_at IS NULL THEN NOW() ELSE refunded_at END,
		    updated_at = NOW()
		WHERE id = $2
	`
	_, err = tx.ExecContext(ctx, query, newStatus, paymentID)
	if err != nil {
		return fmt.Errorf("update payment: %w", err)
	}

	// Логируем событие
	eventQuery := `
		INSERT INTO payment_events (payment_id, event_type, old_status, new_status, vk_event_data)
		VALUES ($1, $2, $3, $4, $5)
	`
	vkDataJSON, _ := json.Marshal(vkData)
	_, err = tx.ExecContext(ctx, eventQuery, paymentID, "status_changed", oldStatus, newStatus, vkDataJSON)
	if err != nil {
		return fmt.Errorf("log event: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	log.Printf("[VKPayService] Payment %s status: %s -> %s", paymentID, oldStatus, newStatus)
	return nil
}

// activateSubscription создает или активирует подписку после успешного платежа
func (s *VKPayService) activateSubscription(ctx context.Context, payment *Payment) error {
	// Загружаем план
	plan, err := s.store.GetSubscriptionPlan(ctx, payment.PlanID)
	if err != nil {
		return fmt.Errorf("get subscription plan: %w", err)
	}

	// Проверяем есть ли уже активная подписка
	existingSubscription, err := s.store.GetActiveSubscription(ctx, payment.ChildProfileID)
	if err != nil && err.Error() != "subscription not found" {
		return fmt.Errorf("get active subscription: %w", err)
	}

	tx, err := s.store.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	var subscriptionID int64

	if existingSubscription != nil {
		// Продлеваем существующую подписку
		expiresAt := existingSubscription.ExpiresAt
		if time.Now().After(expiresAt) {
			// Если подписка истекла, начинаем с текущего момента
			expiresAt = time.Now()
		}
		newExpiresAt := expiresAt.Add(time.Duration(plan.DurationDays) * 24 * time.Hour)

		query := `
			UPDATE subscriptions
			SET plan_id = $1,
			    status = 'active',
			    expires_at = $2,
			    cancelled_at = NULL,
			    auto_renew = TRUE,
			    updated_at = NOW()
			WHERE id = $3
			RETURNING id
		`
		err = tx.QueryRowContext(ctx, query, plan.ID, newExpiresAt, existingSubscription.ID).Scan(&subscriptionID)
		if err != nil {
			return fmt.Errorf("extend subscription: %w", err)
		}

		log.Printf("[VKPayService] Extended subscription %d until %s", subscriptionID, newExpiresAt)
	} else {
		// Создаем новую подписку
		startedAt := time.Now()
		expiresAt := startedAt.Add(time.Duration(plan.DurationDays) * 24 * time.Hour)

		query := `
			INSERT INTO subscriptions (
				child_profile_id, plan_id, status, started_at, expires_at,
				auto_renew, payment_provider, payment_external_id
			) VALUES ($1, $2, 'active', $3, $4, TRUE, 'vk_pay', $5)
			RETURNING id
		`
		err = tx.QueryRowContext(ctx, query,
			payment.ChildProfileID, plan.ID, startedAt, expiresAt, payment.VKOrderID,
		).Scan(&subscriptionID)
		if err != nil {
			return fmt.Errorf("create subscription: %w", err)
		}

		log.Printf("[VKPayService] Created new subscription %d until %s", subscriptionID, expiresAt)
	}

	// Связываем платеж с подпиской
	_, err = tx.ExecContext(ctx, `UPDATE payments SET subscription_id = $1 WHERE id = $2`, subscriptionID, payment.ID)
	if err != nil {
		return fmt.Errorf("link payment to subscription: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// ValidateVKSignature валидирует подпись от VK
func (s *VKPayService) ValidateVKSignature(params url.Values) bool {
	vkSign := params.Get("sign")
	if vkSign == "" {
		return false
	}

	// Удаляем sign из параметров
	params.Del("sign")

	// Сортируем ключи
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Формируем строку для подписи
	var parts []string
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", key, params.Get(key)))
	}
	queryString := strings.Join(parts, "&")

	// Вычисляем HMAC-SHA256
	h := hmac.New(sha256.New, []byte(s.config.AppSecret))
	h.Write([]byte(queryString))
	expectedSign := hex.EncodeToString(h.Sum(nil))

	return hmac.Equal([]byte(vkSign), []byte(expectedSign))
}

// getPaymentByOrderID получает платеж по VK order_id
func (s *VKPayService) getPaymentByOrderID(ctx context.Context, orderID string) (*Payment, error) {
	query := `
		SELECT id, child_profile_id, plan_id, amount_cents, currency, status, created_at, updated_at
		FROM payments
		WHERE vk_order_id = $1
	`

	payment := &Payment{VKOrderID: &orderID}
	err := s.store.DB.QueryRowContext(ctx, query, orderID).Scan(
		&payment.ID, &payment.ChildProfileID, &payment.PlanID,
		&payment.AmountCents, &payment.Currency, &payment.Status,
		&payment.CreatedAt, &payment.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("query payment: %w", err)
	}

	return payment, nil
}

// GetPayment получает платеж по ID
func (s *VKPayService) GetPayment(ctx context.Context, paymentID string) (*Payment, error) {
	query := `
		SELECT id, child_profile_id, plan_id, amount_cents, currency,
		       vk_order_id, vk_transaction_id, vk_user_id, status,
		       payment_method, description, paid_at, created_at, updated_at
		FROM payments
		WHERE id = $1
	`

	payment := &Payment{}
	err := s.store.DB.QueryRowContext(ctx, query, paymentID).Scan(
		&payment.ID, &payment.ChildProfileID, &payment.PlanID,
		&payment.AmountCents, &payment.Currency,
		&payment.VKOrderID, &payment.VKTransactionID, &payment.VKUserID, &payment.Status,
		&payment.PaymentMethod, &payment.Description, &payment.PaidAt,
		&payment.CreatedAt, &payment.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("query payment: %w", err)
	}

	return payment, nil
}

// logEvent логирует событие платежа
func (s *VKPayService) logEvent(ctx context.Context, paymentID, eventType string, oldStatus, newStatus, vkEventType *string) error {
	query := `
		INSERT INTO payment_events (payment_id, event_type, old_status, new_status, vk_event_type)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := s.store.DB.ExecContext(ctx, query, paymentID, eventType, oldStatus, newStatus, vkEventType)
	return err
}

// Helper functions
func strPtr(s string) *string {
	return &s
}
