package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"child-bot/api/internal/api/response"
	"child-bot/api/internal/service"
)

// VKPayWebhookHandler обрабатывает webhooks от VK Pay
type VKPayWebhookHandler struct {
	vkPayService *service.VKPayService
}

// NewVKPayWebhookHandler создает новый VKPayWebhookHandler
func NewVKPayWebhookHandler(vkPayService *service.VKPayService) *VKPayWebhookHandler {
	return &VKPayWebhookHandler{
		vkPayService: vkPayService,
	}
}

// VKPayNotification структура уведомления от VK Pay
type VKPayNotification struct {
	Type    string                 `json:"type"`
	Object  map[string]interface{} `json:"object"`
	GroupID int64                  `json:"group_id"`
	EventID string                 `json:"event_id"`
	AppID   int                    `json:"app_id,omitempty"`
	UserID  int                    `json:"user_id,omitempty"`
}

// HandleWebhook обрабатывает webhook от VK Pay
// POST /webhooks/vk-pay
func (h *VKPayWebhookHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	// Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("[VKPayWebhook] Failed to read request body: %v", err)
		response.BadRequest(w, "Invalid request body")
		return
	}
	defer r.Body.Close()

	log.Printf("[VKPayWebhook] Received webhook: %s", string(body))

	// Парсим JSON
	var notification VKPayNotification
	if err := json.Unmarshal(body, &notification); err != nil {
		log.Printf("[VKPayWebhook] Failed to parse JSON: %v", err)
		response.BadRequest(w, "Invalid JSON")
		return
	}

	// Обрабатываем разные типы уведомлений
	switch notification.Type {
	case "confirmation":
		// VK отправляет confirmation событие при первой настройке webhook
		// Нужно вернуть confirmation code из настроек приложения
		// TODO: Добавить в конфиг VK_CONFIRMATION_CODE
		confirmationCode := "your_confirmation_code_here"
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(confirmationCode))
		log.Printf("[VKPayWebhook] Sent confirmation code")
		return

	case "vkpay_transaction":
		// Обработка транзакции VK Pay
		if err := h.handleTransaction(r.Context(), notification.Object); err != nil {
			log.Printf("[VKPayWebhook] Failed to handle transaction: %v", err)
			response.InternalError(w, "Failed to process transaction")
			return
		}

	case "order_status_change", "order_status_change_test":
		// Обработка изменения статуса заказа
		if err := h.handleOrderStatusChange(r.Context(), notification.Type, notification.Object); err != nil {
			log.Printf("[VKPayWebhook] Failed to handle order status change: %v", err)
			response.InternalError(w, "Failed to process order status change")
			return
		}

	default:
		log.Printf("[VKPayWebhook] Unknown notification type: %s", notification.Type)
	}

	// VK требует ответ "ok" для всех событий кроме confirmation
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// handleTransaction обрабатывает транзакцию VK Pay
func (h *VKPayWebhookHandler) handleTransaction(ctx context.Context, object map[string]interface{}) error {
	// Извлекаем данные транзакции
	orderID, _ := object["order_id"].(string)
	userID, _ := object["from_id"].(float64)

	if orderID == "" {
		return fmt.Errorf("missing order_id in transaction")
	}

	log.Printf("[VKPayWebhook] Processing transaction for order: %s, user: %.0f", orderID, userID)

	// Передаем в VK Pay Service для обработки
	return h.vkPayService.ProcessWebhook(ctx, fmt.Sprintf("%.0f", userID), orderID, "vkpay_transaction", object)
}

// handleOrderStatusChange обрабатывает изменение статуса заказа
func (h *VKPayWebhookHandler) handleOrderStatusChange(ctx context.Context, notificationType string, object map[string]interface{}) error {
	// Извлекаем order_id
	orderID, _ := object["order_id"].(string)
	if orderID == "" {
		// Пробуем альтернативные поля
		if oid, ok := object["id"].(string); ok {
			orderID = oid
		}
	}

	if orderID == "" {
		return fmt.Errorf("missing order_id in status change")
	}

	// Извлекаем user_id если есть
	var userID string
	if uid, ok := object["user_id"].(float64); ok {
		userID = fmt.Sprintf("%.0f", uid)
	}

	log.Printf("[VKPayWebhook] Order %s status change, type: %s", orderID, notificationType)

	// Передаем в VK Pay Service
	return h.vkPayService.ProcessWebhook(ctx, userID, orderID, notificationType, object)
}
