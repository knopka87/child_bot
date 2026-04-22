package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"sort"
	"strings"
	"testing"
)

// Генерирует валидный VK sign для тестовых параметров
func generateValidVKSign(params url.Values, secret string) string {
	// Собираем все vk_* параметры
	var paramsList []string
	for key := range params {
		if strings.HasPrefix(key, "vk_") {
			value := params.Get(key)
			paramsList = append(paramsList, fmt.Sprintf("%s=%s", key, value))
		}
	}

	// ВАЖНО: Сортируем по алфавиту
	sort.Strings(paramsList)

	// Создаем строку для подписи
	queryString := strings.Join(paramsList, "&")

	// Вычисляем HMAC-SHA256
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(queryString))
	hash := mac.Sum(nil)

	// Кодируем в base64 URL-safe
	sign := base64.StdEncoding.EncodeToString(hash)
	sign = strings.ReplaceAll(sign, "+", "-")
	sign = strings.ReplaceAll(sign, "/", "_")
	sign = strings.TrimRight(sign, "=")

	return sign
}

func TestVKAuthMiddleware_ValidSign(t *testing.T) {
	// Устанавливаем тестовый secret
	testSecret := "test_secret_key_123"
	os.Setenv("VK_APP_SECRET", testSecret)
	os.Setenv("ENV", "production") // Не dev режим
	defer os.Unsetenv("VK_APP_SECRET")
	defer os.Unsetenv("ENV")

	// Создаем тестовые VK параметры
	params := url.Values{}
	params.Set("vk_user_id", "12345678")
	params.Set("vk_app_id", "54517931")

	// Генерируем валидный sign
	validSign := generateValidVKSign(params, testSecret)
	params.Set("sign", validSign)

	// Создаем тестовый запрос
	req := httptest.NewRequest("GET", "/test?"+params.Encode(), nil)
	rec := httptest.NewRecorder()

	// Создаем handler который просто вернет 200 OK
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Применяем middleware
	handler := VKAuthMiddleware(nextHandler)
	handler.ServeHTTP(rec, req)

	// Проверяем что запрос прошел (200 OK)
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

func TestVKAuthMiddleware_InvalidSign(t *testing.T) {
	// Устанавливаем тестовый secret
	testSecret := "test_secret_key_123"
	os.Setenv("VK_APP_SECRET", testSecret)
	os.Setenv("ENV", "production")
	defer os.Unsetenv("VK_APP_SECRET")
	defer os.Unsetenv("ENV")

	// Создаем тестовые VK параметры
	params := url.Values{}
	params.Set("vk_user_id", "12345678")
	params.Set("vk_app_id", "54517931")
	params.Set("sign", "invalid_signature_123") // Невалидная подпись

	// Создаем тестовый запрос
	req := httptest.NewRequest("GET", "/test?"+params.Encode(), nil)
	rec := httptest.NewRecorder()

	// Создаем handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Применяем middleware
	handler := VKAuthMiddleware(nextHandler)
	handler.ServeHTTP(rec, req)

	// Проверяем что запрос отклонен (401 Unauthorized)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rec.Code)
	}
}

func TestVKAuthMiddleware_MissingSign(t *testing.T) {
	// Устанавливаем тестовый secret
	os.Setenv("VK_APP_SECRET", "test_secret")
	os.Setenv("ENV", "production")
	defer os.Unsetenv("VK_APP_SECRET")
	defer os.Unsetenv("ENV")

	// Создаем параметры БЕЗ sign
	params := url.Values{}
	params.Set("vk_user_id", "12345678")
	params.Set("vk_app_id", "54517931")

	req := httptest.NewRequest("GET", "/test?"+params.Encode(), nil)
	rec := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := VKAuthMiddleware(nextHandler)
	handler.ServeHTTP(rec, req)

	// Должен вернуть 401
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rec.Code)
	}
}

func TestVKAuthMiddleware_DevelopmentMode(t *testing.T) {
	// В dev режиме должен пропускать все запросы
	os.Setenv("ENV", "development")
	defer os.Unsetenv("ENV")

	// Невалидные параметры
	params := url.Values{}
	params.Set("vk_user_id", "12345678")
	params.Set("sign", "totally_invalid")

	req := httptest.NewRequest("GET", "/test?"+params.Encode(), nil)
	rec := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := VKAuthMiddleware(nextHandler)
	handler.ServeHTTP(rec, req)

	// В dev режиме должен пропустить даже невалидный sign
	if rec.Code != http.StatusOK {
		t.Errorf("Dev mode should pass all requests, got %d", rec.Code)
	}
}

func TestVKAuthMiddleware_NoVKParams(t *testing.T) {
	// Запрос без VK параметров должен пройти
	os.Setenv("ENV", "production")
	defer os.Unsetenv("ENV")

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := VKAuthMiddleware(nextHandler)
	handler.ServeHTTP(rec, req)

	// Должен пропустить (не VK запрос)
	if rec.Code != http.StatusOK {
		t.Errorf("Non-VK request should pass, got %d", rec.Code)
	}
}
