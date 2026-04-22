package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCSRFProtection_SafeMethod(t *testing.T) {
	os.Setenv("ENV", "production")
	defer os.Unsetenv("ENV")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := CSRFProtection(nextHandler)

	// GET запрос не требует CSRF token
	req := httptest.NewRequest("GET", "/api/v1/profile", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Должен пропустить
	if rec.Code != http.StatusOK {
		t.Errorf("Safe method should pass, got status %d", rec.Code)
	}

	// Должен установить CSRF cookie для будущих запросов
	cookies := rec.Result().Cookies()
	foundCSRF := false
	for _, cookie := range cookies {
		if cookie.Name == CSRFCookieName {
			foundCSRF = true
			if cookie.Value == "" {
				t.Error("CSRF cookie value should not be empty")
			}
			if cookie.HttpOnly {
				t.Error("CSRF cookie should not be HttpOnly (JS needs access)")
			}
			break
		}
	}

	if !foundCSRF {
		t.Error("CSRF cookie should be set for safe methods")
	}
}

func TestCSRFProtection_UnsafeMethodWithoutToken(t *testing.T) {
	os.Setenv("ENV", "production")
	defer os.Unsetenv("ENV")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := CSRFProtection(nextHandler)

	// POST без CSRF token
	req := httptest.NewRequest("POST", "/api/v1/profile", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Должен заблокировать
	if rec.Code != http.StatusForbidden {
		t.Errorf("Unsafe method without token should be forbidden, got %d", rec.Code)
	}
}

func TestCSRFProtection_UnsafeMethodWithValidToken(t *testing.T) {
	os.Setenv("ENV", "production")
	defer os.Unsetenv("ENV")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := CSRFProtection(nextHandler)

	// Генерируем валидный token
	token, err := generateCSRFToken()
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// POST с валидным CSRF token в cookie и header
	req := httptest.NewRequest("POST", "/api/v1/profile", nil)
	req.AddCookie(&http.Cookie{
		Name:  CSRFCookieName,
		Value: token,
	})
	req.Header.Set(CSRFHeaderName, token)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Должен пропустить
	if rec.Code != http.StatusOK {
		t.Errorf("Unsafe method with valid token should pass, got %d", rec.Code)
	}
}

func TestCSRFProtection_UnsafeMethodWithInvalidToken(t *testing.T) {
	os.Setenv("ENV", "production")
	defer os.Unsetenv("ENV")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := CSRFProtection(nextHandler)

	// Генерируем два разных token
	cookieToken, _ := generateCSRFToken()
	headerToken, _ := generateCSRFToken()

	// POST с разными tokens в cookie и header
	req := httptest.NewRequest("POST", "/api/v1/profile", nil)
	req.AddCookie(&http.Cookie{
		Name:  CSRFCookieName,
		Value: cookieToken,
	})
	req.Header.Set(CSRFHeaderName, headerToken)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Должен заблокировать
	if rec.Code != http.StatusForbidden {
		t.Errorf("Mismatched tokens should be forbidden, got %d", rec.Code)
	}
}

func TestCSRFProtection_ExemptPath(t *testing.T) {
	os.Setenv("ENV", "production")
	defer os.Unsetenv("ENV")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := CSRFProtection(nextHandler)

	// POST к /health должен быть exempt
	req := httptest.NewRequest("POST", "/health", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Должен пропустить без token
	if rec.Code != http.StatusOK {
		t.Errorf("Exempt path should pass without token, got %d", rec.Code)
	}
}

func TestCSRFProtection_DevelopmentMode(t *testing.T) {
	os.Setenv("ENV", "development")
	defer os.Unsetenv("ENV")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := CSRFProtection(nextHandler)

	// POST без token в dev режиме
	req := httptest.NewRequest("POST", "/api/v1/profile", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// В dev режиме должен пропустить с warning
	if rec.Code != http.StatusOK {
		t.Errorf("Dev mode should be lenient, got %d", rec.Code)
	}
}

func TestCSRFProtection_VariousMethods(t *testing.T) {
	os.Setenv("ENV", "production")
	defer os.Unsetenv("ENV")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := CSRFProtection(nextHandler)

	tests := []struct {
		method      string
		shouldPass  bool
		needsToken  bool
		description string
	}{
		{"GET", true, false, "GET is safe"},
		{"HEAD", true, false, "HEAD is safe"},
		{"OPTIONS", true, false, "OPTIONS is safe"},
		{"POST", false, true, "POST needs token"},
		{"PUT", false, true, "PUT needs token"},
		{"DELETE", false, true, "DELETE needs token"},
		{"PATCH", false, true, "PATCH needs token"},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/api/v1/test", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if tt.shouldPass {
				if rec.Code != http.StatusOK {
					t.Errorf("%s: %s, expected 200, got %d", tt.method, tt.description, rec.Code)
				}
			} else {
				if rec.Code != http.StatusForbidden {
					t.Errorf("%s: %s, expected 403, got %d", tt.method, tt.description, rec.Code)
				}
			}
		})
	}
}

func TestGenerateCSRFToken(t *testing.T) {
	// Генерируем несколько токенов
	tokens := make(map[string]bool)

	for i := 0; i < 100; i++ {
		token, err := generateCSRFToken()
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		// Проверяем что token не пустой
		if token == "" {
			t.Error("Token should not be empty")
		}

		// Проверяем что token уникальный
		if tokens[token] {
			t.Error("Generated duplicate token")
		}
		tokens[token] = true

		// Проверяем длину (base64 encoded 32 bytes)
		if len(token) < 40 {
			t.Errorf("Token too short: %d characters", len(token))
		}
	}
}

func TestTokensEqual(t *testing.T) {
	tests := []struct {
		a     string
		b     string
		equal bool
	}{
		{"same", "same", true},
		{"different", "tokens", false},
		{"", "", true},
		{"token", "", false},
		{"", "token", false},
		{"LoNgToKeN123456", "LoNgToKeN123456", true},
		{"LoNgToKeN123456", "LoNgToKeN123457", false},
	}

	for _, tt := range tests {
		result := tokensEqual(tt.a, tt.b)
		if result != tt.equal {
			t.Errorf("tokensEqual(%q, %q) = %v, expected %v", tt.a, tt.b, result, tt.equal)
		}
	}
}

func TestIsSafeMethod(t *testing.T) {
	tests := []struct {
		method string
		safe   bool
	}{
		{"GET", true},
		{"HEAD", true},
		{"OPTIONS", true},
		{"TRACE", true},
		{"POST", false},
		{"PUT", false},
		{"DELETE", false},
		{"PATCH", false},
		{"CONNECT", false},
	}

	for _, tt := range tests {
		result := isSafeMethod(tt.method)
		if result != tt.safe {
			t.Errorf("isSafeMethod(%s) = %v, expected %v", tt.method, result, tt.safe)
		}
	}
}

func TestIsCSRFExempt(t *testing.T) {
	tests := []struct {
		path   string
		exempt bool
	}{
		{"/health", true},
		{"/api/v1/analytics/events", true},
		{"/api/v1/profile", false},
		{"/api/v1/attempts", false},
		{"/healthcheck", false},
	}

	for _, tt := range tests {
		result := isCSRFExempt(tt.path)
		if result != tt.exempt {
			t.Errorf("isCSRFExempt(%s) = %v, expected %v", tt.path, result, tt.exempt)
		}
	}
}

func TestRegenerateCSRFToken(t *testing.T) {
	rec := httptest.NewRecorder()

	token, err := RegenerateCSRFToken(rec)
	if err != nil {
		t.Fatalf("Failed to regenerate token: %v", err)
	}

	if token == "" {
		t.Error("Regenerated token should not be empty")
	}

	// Проверяем что cookie установлена
	cookies := rec.Result().Cookies()
	foundCSRF := false
	for _, cookie := range cookies {
		if cookie.Name == CSRFCookieName {
			foundCSRF = true
			if cookie.Value != token {
				t.Errorf("Cookie value %q doesn't match returned token %q", cookie.Value, token)
			}
			break
		}
	}

	if !foundCSRF {
		t.Error("CSRF cookie should be set")
	}
}

func TestGetCSRFToken(t *testing.T) {
	// Создаем request с CSRF cookie
	req := httptest.NewRequest("GET", "/test", nil)
	expectedToken := "test_token_123"
	req.AddCookie(&http.Cookie{
		Name:  CSRFCookieName,
		Value: expectedToken,
	})

	token := GetCSRFToken(req)
	if token != expectedToken {
		t.Errorf("GetCSRFToken() = %q, expected %q", token, expectedToken)
	}

	// Request без cookie
	req2 := httptest.NewRequest("GET", "/test", nil)
	token2 := GetCSRFToken(req2)
	if token2 != "" {
		t.Errorf("GetCSRFToken() for request without cookie should return empty, got %q", token2)
	}
}
