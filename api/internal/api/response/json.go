package response

import (
	"encoding/json"
	"log"
	"net/http"
)

// Success отправляет успешный JSON ответ
func Success(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

// OK отправляет 200 OK с данными
func OK(w http.ResponseWriter, data interface{}) {
	Success(w, http.StatusOK, data)
}

// Created отправляет 201 Created с данными
func Created(w http.ResponseWriter, data interface{}) {
	Success(w, http.StatusCreated, data)
}

// NoContent отправляет 204 No Content
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// ErrorResponse структура для ошибок
type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}

// Error отправляет ошибку с указанным статусом
func Error(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(ErrorResponse{Error: message}); err != nil {
		log.Printf("failed to encode error response: %v", err)
	}
}

// ErrorWithCode отправляет ошибку с кодом
func ErrorWithCode(w http.ResponseWriter, status int, message, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(ErrorResponse{Error: message, Code: code}); err != nil {
		log.Printf("failed to encode error response: %v", err)
	}
}

// BadRequest отправляет 400 Bad Request
func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, message)
}

// Unauthorized отправляет 401 Unauthorized
func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, message)
}

// Forbidden отправляет 403 Forbidden
func Forbidden(w http.ResponseWriter, message string) {
	Error(w, http.StatusForbidden, message)
}

// NotFound отправляет 404 Not Found
func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, message)
}

// Conflict отправляет 409 Conflict
func Conflict(w http.ResponseWriter, message string) {
	Error(w, http.StatusConflict, message)
}

// InternalError отправляет 500 Internal Server Error
func InternalError(w http.ResponseWriter, message string) {
	Error(w, http.StatusInternalServerError, message)
}
