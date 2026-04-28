package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"os"
)

const (
	// CSRFTokenLength - длина CSRF token в байтах
	CSRFTokenLength = 32

	// CSRFCookieName - имя cookie для CSRF token
	CSRFCookieName = "csrf_token"

	// CSRFHeaderName - имя header для CSRF token
	CSRFHeaderName = "X-CSRF-Token"

	// CSRFCookieMaxAge - время жизни CSRF cookie (24 часа)
	CSRFCookieMaxAge = 86400
)

// CSRFProtection защищает от Cross-Site Request Forgery атак
// Использует Double Submit Cookie pattern
func CSRFProtection(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Safe methods не требуют CSRF защиты
		if isSafeMethod(r.Method) {
			// Для safe methods генерируем новый token если его нет
			ensureCSRFToken(w, r)
			next.ServeHTTP(w, r)
			return
		}

		// Проверяем исключения (public endpoints)
		if isCSRFExempt(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Development режим - более гибкая валидация
		env := os.Getenv("ENV")
		isDev := env == "development" || env == "dev" || env == ""

		// Получаем token из cookie
		cookieToken, err := getCSRFTokenFromCookie(r)
		if err != nil {
			if isDev {
				log.Printf("[CSRF] Warning: No CSRF cookie found, generating new one")
				// В dev режиме создаем новый token и пропускаем запрос
				ensureCSRFToken(w, r)
				next.ServeHTTP(w, r)
				return
			}
			log.Printf("[CSRF] Missing CSRF cookie")
			http.Error(w, "CSRF token missing", http.StatusForbidden)
			return
		}

		// Получаем token из header
		headerToken := r.Header.Get(CSRFHeaderName)
		if headerToken == "" {
			if isDev {
				log.Printf("[CSRF] Warning: No CSRF header found for %s %s", r.Method, r.URL.Path)
				next.ServeHTTP(w, r)
				return
			}
			log.Printf("[CSRF] Missing CSRF header for %s %s", r.Method, r.URL.Path)
			http.Error(w, "CSRF token required in header", http.StatusForbidden)
			return
		}

		// Сравниваем tokens
		if !tokensEqual(cookieToken, headerToken) {
			log.Printf("[CSRF] Token mismatch for %s %s", r.Method, r.URL.Path)
			http.Error(w, "Invalid CSRF token", http.StatusForbidden)
			return
		}

		// Token валиден, пропускаем запрос
		next.ServeHTTP(w, r)
	})
}

// isSafeMethod проверяет что HTTP метод безопасный (не изменяет state)
func isSafeMethod(method string) bool {
	safeMethods := []string{"GET", "HEAD", "OPTIONS", "TRACE"}
	for _, safe := range safeMethods {
		if method == safe {
			return true
		}
	}
	return false
}

// isCSRFExempt проверяет что endpoint освобожден от CSRF проверки
func isCSRFExempt(path string) bool {
	exemptPaths := []string{
		"/health",
		"/api/health",
		"/api/analytics/events",         // Analytics events
		"/api/profiles/child",           // Child profile creation (onboarding)
		"/api/profiles/by-platform",     // Profile lookup by platform credentials
		"/api/onboarding/start",         // Onboarding start
		"/api/onboarding/complete",      // Onboarding complete
		"/api/consent",                  // Consent saving
		"/api/achievements/mark-viewed", // Mark achievements as viewed
	}

	for _, exempt := range exemptPaths {
		if path == exempt {
			return true
		}
	}
	return false
}

// ensureCSRFToken генерирует и устанавливает CSRF token если его нет
func ensureCSRFToken(w http.ResponseWriter, r *http.Request) {
	// Проверяем есть ли уже token
	_, err := getCSRFTokenFromCookie(r)
	if err == nil {
		// Token уже есть
		return
	}

	// Генерируем новый token
	token, err := generateCSRFToken()
	if err != nil {
		log.Printf("[CSRF] Failed to generate token: %v", err)
		return
	}

	// Устанавливаем cookie
	setCSRFCookie(w, token)
}

// generateCSRFToken генерирует cryptographically secure random token
func generateCSRFToken() (string, error) {
	bytes := make([]byte, CSRFTokenLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Кодируем в base64 URL-safe
	token := base64.URLEncoding.EncodeToString(bytes)
	return token, nil
}

// getCSRFTokenFromCookie извлекает CSRF token из cookie
func getCSRFTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(CSRFCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// setCSRFCookie устанавливает CSRF token cookie
func setCSRFCookie(w http.ResponseWriter, token string) {
	env := os.Getenv("ENV")
	isProduction := env == "production"

	// Для работы в iframe (VK Mini App) нужен SameSite=None с Secure=true
	sameSite := http.SameSiteLaxMode
	if isProduction {
		sameSite = http.SameSiteNoneMode // Разрешает cookie в cross-site iframe
	}

	cookie := &http.Cookie{
		Name:     CSRFCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   CSRFCookieMaxAge,
		HttpOnly: false, // JavaScript должен иметь доступ для отправки в header
		Secure:   isProduction, // Обязательно для SameSite=None
		SameSite: sameSite,
	}

	http.SetCookie(w, cookie)
}

// tokensEqual безопасно сравнивает два token (constant time)
func tokensEqual(a, b string) bool {
	// Используем constant-time comparison для защиты от timing attacks
	if len(a) != len(b) {
		return false
	}

	result := 0
	for i := 0; i < len(a); i++ {
		result |= int(a[i]) ^ int(b[i])
	}

	return result == 0
}

// GetCSRFToken - helper для получения CSRF token в handlers
// Используется для отправки token клиенту при первом запросе
func GetCSRFToken(r *http.Request) string {
	token, err := getCSRFTokenFromCookie(r)
	if err != nil {
		return ""
	}
	return token
}

// RegenerateCSRFToken - генерирует новый CSRF token
// Полезно после login/logout для дополнительной безопасности
func RegenerateCSRFToken(w http.ResponseWriter) (string, error) {
	token, err := generateCSRFToken()
	if err != nil {
		return "", err
	}

	setCSRFCookie(w, token)
	return token, nil
}
