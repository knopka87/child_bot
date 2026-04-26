package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"
)

// CORS middleware для поддержки cross-origin запросов
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Разрешенные origins (из env или default для development)
		// Поддерживаем как HTTP так и HTTPS для гибридных запросов
		allowedOrigins := getEnv("ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:3000,https://localhost:5173,https://localhost:3000,http://127.0.0.1:5173,https://127.0.0.1:5173")

		// Проверяем, разрешен ли origin
		if isOriginAllowed(origin, strings.Split(allowedOrigins, ",")) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		} else if origin != "" {
			// Для production: строгая проверка
			// Для development/стейджинг: разрешаем запросы с предупреждением
			env := os.Getenv("ENV")
			if env == "development" || env == "staging" || env == "" {
				// В dev/staging разрешаем все origins для удобства разработки
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				log.Printf("[CORS] Development mode: allowing origin %s", origin)
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Platform-ID, X-Child-Profile-ID, X-VK-Sign")
		w.Header().Set("Access-Control-Expose-Headers", "X-Child-Profile-ID, X-Platform-ID")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		// Обработка preflight запросов
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isOriginAllowed(origin string, allowed []string) bool {
	if origin == "" {
		return false
	}

	for _, o := range allowed {
		if strings.TrimSpace(o) == origin {
			return true
		}
	}
	return false
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
