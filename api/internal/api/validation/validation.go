package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

var (
	// ErrInvalidJSON возвращается при невалидном JSON
	ErrInvalidJSON = errors.New("invalid JSON")
	// ErrEmptyBody возвращается при пустом теле запроса
	ErrEmptyBody = errors.New("empty request body")
	// ErrInvalidUUID возвращается при невалидном UUID
	ErrInvalidUUID = errors.New("invalid UUID")
)

// DecodeJSON декодирует JSON из тела запроса
func DecodeJSON(r *http.Request, v interface{}) error {
	if r.Body == nil {
		return ErrEmptyBody
	}

	defer r.Body.Close()

	// Ограничиваем размер тела запроса (10MB)
	r.Body = http.MaxBytesReader(nil, r.Body, 10*1024*1024)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // Строгий режим: отклоняем неизвестные поля

	if err := dec.Decode(v); err != nil {
		if errors.Is(err, io.EOF) {
			return ErrEmptyBody
		}
		// Возвращаем более подробное сообщение об ошибке
		return fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}

	return nil
}

// ValidateUUID проверяет валидность UUID строки
func ValidateUUID(s string) error {
	if _, err := uuid.Parse(s); err != nil {
		return ErrInvalidUUID
	}
	return nil
}

// ValidateRequired проверяет, что строка не пустая
func ValidateRequired(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	return nil
}

// ValidateMaxLength проверяет максимальную длину строки
func ValidateMaxLength(value, fieldName string, maxLen int) error {
	if len(value) > maxLen {
		return fmt.Errorf("%s must be at most %d characters", fieldName, maxLen)
	}
	return nil
}

// ValidateMinLength проверяет минимальную длину строки
func ValidateMinLength(value, fieldName string, minLen int) error {
	if len(value) < minLen {
		return fmt.Errorf("%s must be at least %d characters", fieldName, minLen)
	}
	return nil
}

// ValidateEnum проверяет, что значение входит в список допустимых
func ValidateEnum(value, fieldName string, validValues []string) error {
	for _, v := range validValues {
		if value == v {
			return nil
		}
	}
	return fmt.Errorf("%s must be one of: %s", fieldName, strings.Join(validValues, ", "))
}

// ValidateAttemptType проверяет тип попытки (help или check)
func ValidateAttemptType(attemptType string) error {
	return ValidateEnum(attemptType, "attempt_type", []string{"help", "check"})
}

// ValidatePlatformID проверяет валидность platformID
func ValidatePlatformID(platformID string) error {
	return ValidateEnum(platformID, "platform_id", []string{"vk", "telegram", "max", "web"})
}

// ValidateBase64Image проверяет, что строка начинается с префикса base64
func ValidateBase64Image(imageData string) error {
	if !strings.HasPrefix(imageData, "data:image/") {
		return errors.New("image must be base64 encoded with data URI scheme")
	}
	return nil
}

// ValidateEmail проверяет валидность email адреса
func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)

	// Базовая проверка формата
	if email == "" {
		return errors.New("email is required")
	}

	// Проверяем наличие @
	if !strings.Contains(email, "@") {
		return errors.New("email must contain @")
	}

	// Разбиваем на части
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return errors.New("email must have exactly one @")
	}

	local := parts[0]
	domain := parts[1]

	// Проверка local part
	if local == "" {
		return errors.New("email local part cannot be empty")
	}

	if len(local) > 64 {
		return errors.New("email local part too long (max 64 characters)")
	}

	// Проверка domain
	if domain == "" {
		return errors.New("email domain cannot be empty")
	}

	if !strings.Contains(domain, ".") {
		return errors.New("email domain must contain a dot")
	}

	if len(domain) > 255 {
		return errors.New("email domain too long (max 255 characters)")
	}

	// Общая длина
	if len(email) > 320 {
		return errors.New("email too long (max 320 characters)")
	}

	return nil
}
