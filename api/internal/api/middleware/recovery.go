package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"child-bot/api/internal/api/response"
)

// Recovery middleware для перехвата panic
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Логируем stack trace
				log.Printf("PANIC: %v\n%s", err, debug.Stack())

				// Отправляем generic ошибку клиенту (не раскрываем детали)
				response.InternalError(w, "Internal server error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}
