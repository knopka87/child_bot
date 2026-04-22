package handler

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"child-bot/api/internal/api/middleware"
	"child-bot/api/internal/api/response"
	"child-bot/api/internal/api/validation"
	"child-bot/api/internal/service"
	"child-bot/api/internal/store"
)

// EmailHandler обрабатывает запросы верификации email
type EmailHandler struct {
	store        *store.Store
	emailService *service.EmailService
}

// NewEmailHandler создает новый EmailHandler
func NewEmailHandler(store *store.Store) *EmailHandler {
	return &EmailHandler{
		store:        store,
		emailService: service.NewEmailService(),
	}
}

type SendVerificationRequest struct {
	Email        string `json:"email"`
	ParentUserID string `json:"parentUserId"`
}

type VerifyCodeRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

// SendVerification отправляет код верификации на email
// POST /email/verify/send
func (h *EmailHandler) SendVerification(w http.ResponseWriter, r *http.Request) {
	// Получаем platformID из middleware
	platformID := middleware.GetPlatformID(r.Context())
	if platformID == "" {
		response.Unauthorized(w, "Missing platform ID")
		return
	}

	var req SendVerificationRequest
	if err := validation.DecodeJSON(r, &req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	// Валидация email
	if err := validation.ValidateEmail(req.Email); err != nil {
		response.BadRequest(w, "Invalid email: "+err.Error())
		return
	}

	if req.ParentUserID == "" {
		response.BadRequest(w, "parentUserId is required")
		return
	}

	// Генерируем 6-значный код
	code, err := generateVerificationCode()
	if err != nil {
		response.InternalError(w, "Failed to generate verification code")
		return
	}

	// Получаем IP для аудита
	ipAddress := r.Header.Get("X-Forwarded-For")
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}

	// Сохраняем в БД
	expiresAt := time.Now().Add(15 * time.Minute)
	err = h.store.CreateEmailVerification(
		r.Context(),
		req.Email,
		code,
		req.ParentUserID,
		platformID,
		ipAddress,
		expiresAt,
	)
	if err != nil {
		log.Printf("[EmailHandler] Failed to create verification: %v", err)
		response.InternalError(w, "Failed to create verification")
		return
	}

	// Отправляем email через email service
	if !h.emailService.IsDevelopmentMode() {
		err = h.emailService.SendVerificationCode(req.Email, code, expiresAt)
		if err != nil {
			log.Printf("[EmailHandler] Failed to send email: %v", err)
			// Не фейлим запрос если email не отправился, код всё равно в БД
		} else {
			log.Printf("[EmailHandler] Verification email sent to %s", req.Email)
		}
	} else {
		log.Printf("[EmailHandler] Dev mode: verification code for %s: %s (expires at %v)", req.Email, code, expiresAt)
	}

	// Формируем ответ
	responseData := map[string]interface{}{
		"message":   "Verification code sent to email",
		"expiresAt": expiresAt,
	}

	// Включаем код только в dev режиме для удобства тестирования
	if h.emailService.IsDevelopmentMode() {
		responseData["devCode"] = code
	}

	response.OK(w, responseData)
}

// VerifyCode проверяет введенный код верификации
// POST /email/verify/check
func (h *EmailHandler) VerifyCode(w http.ResponseWriter, r *http.Request) {
	var req VerifyCodeRequest
	if err := validation.DecodeJSON(r, &req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	// Валидация
	if err := validation.ValidateEmail(req.Email); err != nil {
		response.BadRequest(w, "Invalid email: "+err.Error())
		return
	}

	if len(req.Code) != 6 {
		response.BadRequest(w, "Code must be 6 digits")
		return
	}

	// Проверяем код
	isValid, err := h.store.VerifyEmailCode(r.Context(), req.Email, req.Code)
	if err != nil {
		log.Printf("[EmailHandler] Verification failed: %v", err)
		response.InternalError(w, "Verification failed")
		return
	}

	if !isValid {
		response.BadRequest(w, "Invalid or expired verification code")
		return
	}

	response.OK(w, map[string]interface{}{
		"verified": true,
		"message":  "Email verified successfully",
	})
}

// CheckVerification проверяет статус верификации email
// GET /email/verify/status?email=xxx
func (h *EmailHandler) CheckVerification(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		response.BadRequest(w, "email parameter is required")
		return
	}

	// Валидация email
	if err := validation.ValidateEmail(email); err != nil {
		response.BadRequest(w, "Invalid email: "+err.Error())
		return
	}

	// Проверяем статус
	isVerified, err := h.store.IsEmailVerified(r.Context(), email)
	if err != nil {
		response.InternalError(w, "Failed to check verification status")
		return
	}

	response.OK(w, map[string]interface{}{
		"email":    email,
		"verified": isVerified,
	})
}

// generateVerificationCode генерирует случайный 6-значный код
func generateVerificationCode() (string, error) {
	// Генерируем число от 0 до 999999
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	// Форматируем с ведущими нулями
	return fmt.Sprintf("%06d", n.Int64()), nil
}
