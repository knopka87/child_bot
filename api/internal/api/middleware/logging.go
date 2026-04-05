package middleware

import (
	"log"
	"net/http"
	"time"
)

// responseWriter wrapper для захвата статуса
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

// Logging middleware для логирования запросов
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := wrapResponseWriter(w)

		next.ServeHTTP(wrapped, r)

		// Логируем после выполнения запроса
		duration := time.Since(start)
		status := wrapped.status
		if status == 0 {
			status = http.StatusOK
		}

		log.Printf("%s %s %d %s", r.Method, r.URL.Path, status, duration)
	})
}
