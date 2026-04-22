package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestSecurityHeaders(t *testing.T) {
	// Создаем тестовый handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Применяем middleware
	handler := SecurityHeaders(nextHandler)

	// Создаем тестовый запрос
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	// Выполняем запрос
	handler.ServeHTTP(rec, req)

	// Проверяем что все security headers установлены
	tests := []struct {
		header   string
		expected string
		contains bool // если true, проверяем Contains вместо точного совпадения
	}{
		{
			header:   "Strict-Transport-Security",
			expected: "max-age=31536000; includeSubDomains",
			contains: false,
		},
		{
			header:   "X-Frame-Options",
			expected: "DENY",
			contains: false,
		},
		{
			header:   "X-Content-Type-Options",
			expected: "nosniff",
			contains: false,
		},
		{
			header:   "X-XSS-Protection",
			expected: "1; mode=block",
			contains: false,
		},
		{
			header:   "Content-Security-Policy",
			expected: "default-src 'self'",
			contains: true,
		},
		{
			header:   "Content-Security-Policy",
			expected: "frame-ancestors 'none'",
			contains: true,
		},
		{
			header:   "Referrer-Policy",
			expected: "strict-origin-when-cross-origin",
			contains: false,
		},
		{
			header:   "Permissions-Policy",
			expected: "geolocation=()",
			contains: true,
		},
	}

	for _, tt := range tests {
		actual := rec.Header().Get(tt.header)
		if actual == "" {
			t.Errorf("Header %s not set", tt.header)
			continue
		}

		if tt.contains {
			if !strings.Contains(actual, tt.expected) {
				t.Errorf("Header %s does not contain %q, got: %q", tt.header, tt.expected, actual)
			}
		} else {
			if actual != tt.expected {
				t.Errorf("Header %s = %q, expected %q", tt.header, actual, tt.expected)
			}
		}
	}
}

func TestHTTPSRedirect_Production(t *testing.T) {
	// Устанавливаем production режим
	os.Setenv("ENV", "production")
	defer os.Unsetenv("ENV")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := HTTPSRedirect(nextHandler)

	// Тест 1: HTTP запрос должен редиректиться
	t.Run("redirect_http_to_https", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/test?foo=bar", nil)
		req.Header.Set("X-Forwarded-Proto", "http")
		req.Host = "example.com"
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		// Проверяем редирект
		if rec.Code != http.StatusMovedPermanently {
			t.Errorf("Expected status 301, got %d", rec.Code)
		}

		location := rec.Header().Get("Location")
		expected := "https://example.com/test?foo=bar"
		if location != expected {
			t.Errorf("Location = %q, expected %q", location, expected)
		}
	})

	// Тест 2: HTTPS запрос не должен редиректиться
	t.Run("no_redirect_for_https", func(t *testing.T) {
		req := httptest.NewRequest("GET", "https://example.com/test", nil)
		req.Header.Set("X-Forwarded-Proto", "https")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		// Должен пропустить дальше без редиректа
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}

		location := rec.Header().Get("Location")
		if location != "" {
			t.Errorf("Should not redirect HTTPS requests, got Location: %q", location)
		}
	})
}

func TestHTTPSRedirect_Development(t *testing.T) {
	// Устанавливаем development режим
	os.Setenv("ENV", "development")
	defer os.Unsetenv("ENV")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := HTTPSRedirect(nextHandler)

	// В dev режиме HTTP запросы должны проходить без редиректа
	req := httptest.NewRequest("GET", "http://localhost:8080/test", nil)
	req.Header.Set("X-Forwarded-Proto", "http")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Не должно быть редиректа
	if rec.Code != http.StatusOK {
		t.Errorf("Dev mode should not redirect, got status %d", rec.Code)
	}

	location := rec.Header().Get("Location")
	if location != "" {
		t.Errorf("Dev mode should not redirect, got Location: %q", location)
	}
}

func TestHTTPSRedirect_NoEnv(t *testing.T) {
	// Удаляем ENV переменную (default behavior = development)
	os.Unsetenv("ENV")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := HTTPSRedirect(nextHandler)

	// Без ENV должно работать как development
	req := httptest.NewRequest("GET", "http://localhost:8080/test", nil)
	req.Header.Set("X-Forwarded-Proto", "http")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Не должно быть редиректа
	if rec.Code != http.StatusOK {
		t.Errorf("No ENV should behave as dev mode, got status %d", rec.Code)
	}
}

func TestSecureCookie(t *testing.T) {
	tests := []struct {
		name        string
		env         string
		expectHTTPS bool
	}{
		{
			name:        "production_mode",
			env:         "production",
			expectHTTPS: true,
		},
		{
			name:        "development_mode",
			env:         "development",
			expectHTTPS: false,
		},
		{
			name:        "no_env",
			env:         "",
			expectHTTPS: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.env != "" {
				os.Setenv("ENV", tt.env)
				defer os.Unsetenv("ENV")
			} else {
				os.Unsetenv("ENV")
			}

			cookie := SecureCookie("test_cookie", "test_value", 3600)

			// Проверяем параметры cookie
			if cookie.Name != "test_cookie" {
				t.Errorf("Name = %q, expected %q", cookie.Name, "test_cookie")
			}

			if cookie.Value != "test_value" {
				t.Errorf("Value = %q, expected %q", cookie.Value, "test_value")
			}

			if cookie.MaxAge != 3600 {
				t.Errorf("MaxAge = %d, expected %d", cookie.MaxAge, 3600)
			}

			if !cookie.HttpOnly {
				t.Error("HttpOnly should be true")
			}

			if cookie.Secure != tt.expectHTTPS {
				t.Errorf("Secure = %v, expected %v", cookie.Secure, tt.expectHTTPS)
			}

			if cookie.SameSite != http.SameSiteStrictMode {
				t.Errorf("SameSite = %v, expected %v", cookie.SameSite, http.SameSiteStrictMode)
			}

			if cookie.Path != "/" {
				t.Errorf("Path = %q, expected %q", cookie.Path, "/")
			}
		})
	}
}
