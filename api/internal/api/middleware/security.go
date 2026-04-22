package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"
)

// SecurityHeaders добавляет security headers ко всем ответам
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// HSTS (HTTP Strict Transport Security)
		// Принудительное использование HTTPS на 1 год
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// X-Frame-Options - защита от clickjacking
		// DENY полностью запрещает загрузку в iframe
		w.Header().Set("X-Frame-Options", "DENY")

		// X-Content-Type-Options - защита от MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// X-XSS-Protection - legacy XSS защита (для старых браузеров)
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Content-Security-Policy - защита от XSS и injection атак
		// Разрешаем только собственные скрипты и стили, VK Bridge API
		csp := strings.Join([]string{
			"default-src 'self'",
			"script-src 'self' 'unsafe-inline' 'unsafe-eval' https://vk.com https://*.vk.com https://*.vk.me",
			"style-src 'self' 'unsafe-inline'",
			"img-src 'self' data: https: blob:",
			"font-src 'self' data:",
			"connect-src 'self' https://vk.com https://*.vk.com https://*.vk.me wss://im-0*.vk.com",
			"frame-src https://vk.com https://*.vk.com",
			"frame-ancestors 'none'",
			"form-action 'self'",
			"base-uri 'self'",
			"object-src 'none'",
		}, "; ")
		w.Header().Set("Content-Security-Policy", csp)

		// Referrer-Policy - контроль referrer информации
		// strict-origin-when-cross-origin отправляет полный referrer только для same-origin
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions-Policy (ранее Feature-Policy)
		// Отключаем ненужные browser features
		permissions := strings.Join([]string{
			"geolocation=()",
			"microphone=()",
			"camera=()",
			"payment=()",
			"usb=()",
			"magnetometer=()",
			"gyroscope=()",
			"accelerometer=()",
		}, ", ")
		w.Header().Set("Permissions-Policy", permissions)

		next.ServeHTTP(w, r)
	})
}

// HTTPSRedirect редиректит HTTP на HTTPS в production режиме
func HTTPSRedirect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// В development режиме пропускаем редирект
		env := os.Getenv("ENV")
		if env == "development" || env == "dev" || env == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Проверяем протокол из заголовка X-Forwarded-Proto (для reverse proxy)
		proto := r.Header.Get("X-Forwarded-Proto")
		if proto == "" {
			// Если заголовка нет, проверяем схему запроса
			if r.TLS == nil {
				proto = "http"
			} else {
				proto = "https"
			}
		}

		// Если не HTTPS - редиректим
		if proto != "https" {
			// Строим HTTPS URL
			host := r.Host
			if host == "" {
				host = r.Header.Get("Host")
			}

			// Используем r.URL для правильного построения пути
			path := r.URL.Path
			if r.URL.RawQuery != "" {
				path += "?" + r.URL.RawQuery
			}

			httpsURL := "https://" + host + path

			log.Printf("[Security] Redirecting HTTP -> HTTPS: %s", path)

			// 301 Permanent Redirect
			http.Redirect(w, r, httpsURL, http.StatusMovedPermanently)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// SecureCookies настраивает безопасные параметры cookie
// Этот helper можно использовать при установке cookies в handlers
func SecureCookie(name, value string, maxAge int) *http.Cookie {
	env := os.Getenv("ENV")
	isProduction := env == "production"

	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,                    // Защита от XSS
		Secure:   isProduction,            // HTTPS only в production
		SameSite: http.SameSiteStrictMode, // CSRF защита
	}
}
