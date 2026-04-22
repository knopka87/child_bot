package handler

import (
	"net/http"

	"child-bot/api/internal/api/response"
	"child-bot/api/internal/store"
)

// LegalHandler обрабатывает запросы юридических документов
type LegalHandler struct {
	store *store.Store
}

// NewLegalHandler создает новый LegalHandler
func NewLegalHandler(store *store.Store) *LegalHandler {
	return &LegalHandler{store: store}
}

// GetPrivacyPolicy возвращает текст политики конфиденциальности
// GET /legal/privacy
func (h *LegalHandler) GetPrivacyPolicy(w http.ResponseWriter, r *http.Request) {
	// Получаем язык из query параметра, по умолчанию 'ru'
	language := r.URL.Query().Get("lang")
	if language == "" {
		language = "ru"
	}

	// Получаем активную версию документа
	doc, err := h.store.GetActiveLegalDocument(r.Context(), "privacy_policy", language)
	if err != nil {
		response.InternalError(w, "Failed to get privacy policy")
		return
	}

	if doc == nil {
		response.NotFound(w, "Privacy policy not found")
		return
	}

	response.OK(w, map[string]interface{}{
		"id":             doc.ID,
		"version":        doc.Version,
		"title":          doc.Title,
		"content":        doc.Content,
		"language":       doc.Language,
		"effectiveDate":  doc.EffectiveDate,
		"lastUpdated":    doc.UpdatedAt,
	})
}

// GetTermsOfService возвращает текст условий использования
// GET /legal/terms
func (h *LegalHandler) GetTermsOfService(w http.ResponseWriter, r *http.Request) {
	// Получаем язык из query параметра, по умолчанию 'ru'
	language := r.URL.Query().Get("lang")
	if language == "" {
		language = "ru"
	}

	// Получаем активную версию документа
	doc, err := h.store.GetActiveLegalDocument(r.Context(), "terms_of_service", language)
	if err != nil {
		response.InternalError(w, "Failed to get terms of service")
		return
	}

	if doc == nil {
		response.NotFound(w, "Terms of service not found")
		return
	}

	response.OK(w, map[string]interface{}{
		"id":             doc.ID,
		"version":        doc.Version,
		"title":          doc.Title,
		"content":        doc.Content,
		"language":       doc.Language,
		"effectiveDate":  doc.EffectiveDate,
		"lastUpdated":    doc.UpdatedAt,
	})
}
