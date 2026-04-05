package middleware

import (
	"net/http"
	"os"
	"strings"
)

// CORS middleware для поддержки cross-origin запросов
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Разрешенные origins (из env или default для development)
		allowedOrigins := getEnv("ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:3000")

		// Проверяем, разрешен ли origin
		if isOriginAllowed(origin, strings.Split(allowedOrigins, ",")) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Platform-ID, X-Child-Profile-ID")
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
