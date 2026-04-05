package domain

import "errors"

// Common domain errors
var (
	// ErrNotFound возвращается, когда ресурс не найден
	ErrNotFound = errors.New("resource not found")

	// ErrInvalidInput возвращается при невалидных входных данных
	ErrInvalidInput = errors.New("invalid input")

	// ErrUnauthorized возвращается при отсутствии авторизации
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden возвращается при недостаточных правах
	ErrForbidden = errors.New("forbidden")

	// ErrConflict возвращается при конфликте данных
	ErrConflict = errors.New("conflict")

	// ErrInternal возвращается при внутренней ошибке
	ErrInternal = errors.New("internal error")

	// ErrAttemptNotFound возвращается, когда попытка не найдена
	ErrAttemptNotFound = errors.New("attempt not found")

	// ErrAttemptAlreadyProcessed возвращается, когда попытка уже обработана
	ErrAttemptAlreadyProcessed = errors.New("attempt already processed")

	// ErrNoHintsAvailable возвращается, когда подсказки закончились
	ErrNoHintsAvailable = errors.New("no more hints available")
)
