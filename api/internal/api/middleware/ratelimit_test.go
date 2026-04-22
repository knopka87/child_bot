package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestRateLimit_AllowWithinLimit(t *testing.T) {
	os.Setenv("ENV", "production")
	defer os.Unsetenv("ENV")

	config := RateLimitConfig{
		RequestsPerWindow: 5,
		WindowDuration:    time.Second,
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := RateLimit(config)(nextHandler)

	// Отправляем 5 запросов (в пределах лимита)
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Request %d: expected status 200, got %d", i+1, rec.Code)
		}
	}
}

func TestRateLimit_BlockExceedingLimit(t *testing.T) {
	os.Setenv("ENV", "production")
	defer os.Unsetenv("ENV")

	config := RateLimitConfig{
		RequestsPerWindow: 3,
		WindowDuration:    time.Second,
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := RateLimit(config)(nextHandler)

	// Отправляем 3 запроса (достигаем лимита)
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Request %d: expected status 200, got %d", i+1, rec.Code)
		}
	}

	// 4-й запрос должен быть заблокирован
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("Request over limit: expected status 429, got %d", rec.Code)
	}

	// Проверяем headers
	if rec.Header().Get("Retry-After") == "" {
		t.Error("Retry-After header should be set")
	}

	if rec.Header().Get("X-RateLimit-Limit") != "3" {
		t.Errorf("X-RateLimit-Limit = %s, expected 3", rec.Header().Get("X-RateLimit-Limit"))
	}

	if rec.Header().Get("X-RateLimit-Remaining") != "0" {
		t.Errorf("X-RateLimit-Remaining = %s, expected 0", rec.Header().Get("X-RateLimit-Remaining"))
	}
}

func TestRateLimit_SlidingWindow(t *testing.T) {
	os.Setenv("ENV", "production")
	defer os.Unsetenv("ENV")

	config := RateLimitConfig{
		RequestsPerWindow: 2,
		WindowDuration:    500 * time.Millisecond,
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := RateLimit(config)(nextHandler)

	// Отправляем 2 запроса (достигаем лимита)
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Request %d: expected status 200, got %d", i+1, rec.Code)
		}
	}

	// 3-й запрос сразу должен быть заблокирован
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("Request over limit: expected status 429, got %d", rec.Code)
	}

	// Ждем истечения окна
	time.Sleep(600 * time.Millisecond)

	// Теперь запрос должен пройти (старые запросы вне окна)
	req = httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	rec = httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("After window: expected status 200, got %d", rec.Code)
	}
}

func TestRateLimit_DifferentClients(t *testing.T) {
	os.Setenv("ENV", "production")
	defer os.Unsetenv("ENV")

	config := RateLimitConfig{
		RequestsPerWindow: 2,
		WindowDuration:    time.Second,
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := RateLimit(config)(nextHandler)

	// Client 1 отправляет 2 запроса
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Client1 request %d: expected status 200, got %d", i+1, rec.Code)
		}
	}

	// Client 1: 3-й запрос блокируется
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("Client1 over limit: expected status 429, got %d", rec.Code)
	}

	// Client 2 должен иметь свой лимит
	req = httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.2:12345"
	rec = httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Client2 first request: expected status 200, got %d", rec.Code)
	}
}

func TestRateLimit_DevelopmentMode(t *testing.T) {
	os.Setenv("ENV", "development")
	defer os.Unsetenv("ENV")

	config := RateLimitConfig{
		RequestsPerWindow: 1,
		WindowDuration:    time.Second,
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := RateLimit(config)(nextHandler)

	// В dev режиме должны проходить все запросы
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Dev mode request %d: expected status 200, got %d", i+1, rec.Code)
		}
	}
}

func TestRateLimit_XForwardedFor(t *testing.T) {
	os.Setenv("ENV", "production")
	defer os.Unsetenv("ENV")

	config := RateLimitConfig{
		RequestsPerWindow: 2,
		WindowDuration:    time.Second,
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := RateLimit(config)(nextHandler)

	// Запросы с X-Forwarded-For должны использовать этот IP
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "10.0.0.1:12345" // Proxy IP
		req.Header.Set("X-Forwarded-For", "203.0.113.1")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Request %d: expected status 200, got %d", i+1, rec.Code)
		}
	}

	// 3-й запрос с тем же X-Forwarded-For должен блокироваться
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "10.0.0.1:12345"
	req.Header.Set("X-Forwarded-For", "203.0.113.1")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("Over limit: expected status 429, got %d", rec.Code)
	}

	// Запрос с другим X-Forwarded-For должен проходить
	req = httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "10.0.0.1:12345"
	req.Header.Set("X-Forwarded-For", "203.0.113.2")
	rec = httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Different IP: expected status 200, got %d", rec.Code)
	}
}

func TestGetRemainingRequests(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerWindow: 5,
		WindowDuration:    time.Second,
	}

	testLimiter := newRateLimiter(config)
	clientID := "192.168.1.1:12345"

	// Изначально должно быть 5 доступных запросов
	remaining := testLimiter.GetRemainingRequests(clientID, config)
	if remaining != 5 {
		t.Errorf("Initial remaining = %d, expected 5", remaining)
	}

	// Делаем 2 запроса
	for i := 0; i < 2; i++ {
		testLimiter.Allow(clientID, config)
	}

	// Должно остаться 3
	remaining = testLimiter.GetRemainingRequests(clientID, config)
	if remaining != 3 {
		t.Errorf("After 2 requests remaining = %d, expected 3", remaining)
	}

	// Делаем еще 3 запроса (достигаем лимита)
	for i := 0; i < 3; i++ {
		testLimiter.Allow(clientID, config)
	}

	// Должно быть 0
	remaining = testLimiter.GetRemainingRequests(clientID, config)
	if remaining != 0 {
		t.Errorf("After 5 requests remaining = %d, expected 0", remaining)
	}
}

func TestGetClientIP(t *testing.T) {
	tests := []struct {
		name         string
		remoteAddr   string
		forwardedFor string
		realIP       string
		expectedIP   string
	}{
		{
			name:       "direct_connection",
			remoteAddr: "192.168.1.1:12345",
			expectedIP: "192.168.1.1:12345",
		},
		{
			name:         "x_forwarded_for",
			remoteAddr:   "10.0.0.1:12345",
			forwardedFor: "203.0.113.1",
			expectedIP:   "203.0.113.1",
		},
		{
			name:       "x_real_ip",
			remoteAddr: "10.0.0.1:12345",
			realIP:     "203.0.113.2",
			expectedIP: "203.0.113.2",
		},
		{
			name:         "x_forwarded_for_priority",
			remoteAddr:   "10.0.0.1:12345",
			forwardedFor: "203.0.113.1",
			realIP:       "203.0.113.2",
			expectedIP:   "203.0.113.1", // X-Forwarded-For имеет приоритет
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = tt.remoteAddr

			if tt.forwardedFor != "" {
				req.Header.Set("X-Forwarded-For", tt.forwardedFor)
			}

			if tt.realIP != "" {
				req.Header.Set("X-Real-IP", tt.realIP)
			}

			ip := getClientIP(req)
			if ip != tt.expectedIP {
				t.Errorf("getClientIP() = %s, expected %s", ip, tt.expectedIP)
			}
		})
	}
}
