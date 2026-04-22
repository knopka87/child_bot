package handler

import (
	"net/http"

	"child-bot/api/internal/api/middleware"
	"child-bot/api/internal/api/response"
)

// CSRFHandler обрабатывает запросы связанные с CSRF
type CSRFHandler struct{}

// NewCSRFHandler создает новый CSRFHandler
func NewCSRFHandler() *CSRFHandler {
	return &CSRFHandler{}
}

// GetToken возвращает CSRF token для клиента
// GET /csrf-token
func (h *CSRFHandler) GetToken(w http.ResponseWriter, r *http.Request) {
	// Получаем существующий token или генерируем новый
	token := middleware.GetCSRFToken(r)

	if token == "" {
		// Генерируем новый token если нет
		var err error
		token, err = middleware.RegenerateCSRFToken(w)
		if err != nil {
			response.InternalError(w, "Failed to generate CSRF token")
			return
		}
	}

	response.OK(w, map[string]interface{}{
		"csrfToken":  token,
		"headerName": "X-CSRF-Token",
		"cookieName": "csrf_token",
	})
}
