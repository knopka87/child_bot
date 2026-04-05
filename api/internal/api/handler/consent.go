package handler

import (
	"net/http"

	"child-bot/api/internal/api/middleware"
	"child-bot/api/internal/api/response"
	"child-bot/api/internal/api/validation"
	"child-bot/api/internal/store"
)

// ConsentHandler обрабатывает запросы согласий
type ConsentHandler struct {
	store *store.Store
}

// NewConsentHandler создает новый ConsentHandler
func NewConsentHandler(store *store.Store) *ConsentHandler {
	return &ConsentHandler{store: store}
}

type SaveConsentRequest struct {
	ParentUserID         string `json:"parentUserId"`
	PrivacyPolicyVersion string `json:"privacyPolicyVersion"`
	TermsVersion         string `json:"termsVersion"`
	AdultConsent         bool   `json:"adultConsent"`
}

// SaveConsent сохраняет согласие пользователя
// POST /consent
func (h *ConsentHandler) SaveConsent(w http.ResponseWriter, r *http.Request) {
	// Получаем platformID из middleware
	platformID := middleware.GetPlatformID(r.Context())
	if platformID == "" {
		response.Unauthorized(w, "Missing platform ID")
		return
	}

	var req SaveConsentRequest
	if err := validation.DecodeJSON(r, &req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	// Валидация
	if req.ParentUserID == "" {
		response.BadRequest(w, "parentUserId is required")
		return
	}
	if req.PrivacyPolicyVersion == "" {
		response.BadRequest(w, "privacyPolicyVersion is required")
		return
	}
	if req.TermsVersion == "" {
		response.BadRequest(w, "termsVersion is required")
		return
	}
	if !req.AdultConsent {
		response.BadRequest(w, "adultConsent must be true")
		return
	}

	// Получаем IP и User-Agent для аудита
	ipAddress := r.Header.Get("X-Forwarded-For")
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}
	userAgent := r.Header.Get("User-Agent")

	// Сохранение в БД
	err := h.store.SaveParentConsent(
		r.Context(),
		platformID,
		req.ParentUserID,
		req.PrivacyPolicyVersion,
		req.TermsVersion,
		req.AdultConsent,
		ipAddress,
		userAgent,
	)
	if err != nil {
		response.InternalError(w, "Failed to save consent")
		return
	}

	response.OK(w, map[string]string{"message": "Consent saved successfully"})
}

// CheckConsent проверяет наличие действительного согласия
// GET /consent/check?parentUserId=xxx
func (h *ConsentHandler) CheckConsent(w http.ResponseWriter, r *http.Request) {
	// Получаем platformID из middleware
	platformID := middleware.GetPlatformID(r.Context())
	if platformID == "" {
		response.Unauthorized(w, "Missing platform ID")
		return
	}

	parentUserID := r.URL.Query().Get("parentUserId")
	if parentUserID == "" {
		response.BadRequest(w, "parentUserId is required")
		return
	}

	// Проверка наличия согласия
	hasConsent, err := h.store.HasValidConsent(r.Context(), platformID, parentUserID)
	if err != nil {
		response.InternalError(w, "Failed to check consent")
		return
	}

	response.OK(w, map[string]bool{"hasConsent": hasConsent})
}

// GetConsent получает информацию о согласии
// GET /consent?parentUserId=xxx
func (h *ConsentHandler) GetConsent(w http.ResponseWriter, r *http.Request) {
	// Получаем platformID из middleware
	platformID := middleware.GetPlatformID(r.Context())
	if platformID == "" {
		response.Unauthorized(w, "Missing platform ID")
		return
	}

	parentUserID := r.URL.Query().Get("parentUserId")
	if parentUserID == "" {
		response.BadRequest(w, "parentUserId is required")
		return
	}

	// Получение согласия
	consent, err := h.store.GetParentConsent(r.Context(), platformID, parentUserID)
	if err != nil {
		response.InternalError(w, "Failed to get consent")
		return
	}

	if consent == nil {
		response.NotFound(w, "Consent not found")
		return
	}

	// Формируем ответ
	consentData := map[string]interface{}{
		"id":                      consent.ID,
		"parentUserId":            consent.ParentUserID,
		"privacyPolicyVersion":    consent.PrivacyPolicyVersion,
		"privacyPolicyAccepted":   consent.PrivacyPolicyAccepted,
		"termsVersion":            consent.TermsVersion,
		"termsAccepted":           consent.TermsAccepted,
		"adultConsent":            consent.AdultConsent,
		"createdAt":               consent.CreatedAt,
		"updatedAt":               consent.UpdatedAt,
	}

	response.OK(w, consentData)
}
