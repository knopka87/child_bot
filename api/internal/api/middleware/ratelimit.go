package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// RateLimitConfig определяет конфигурацию rate limiting
type RateLimitConfig struct {
	RequestsPerWindow int           // Количество запросов
	WindowDuration    time.Duration // Размер временного окна
}

// Preset конфигурации для разных типов endpoints
var (
	// RateLimitDefault - стандартный лимит для большинства endpoints
	RateLimitDefault = RateLimitConfig{
		RequestsPerWindow: 100,
		WindowDuration:    time.Minute,
	}

	// RateLimitStrict - строгий лимит для критических операций
	RateLimitStrict = RateLimitConfig{
		RequestsPerWindow: 10,
		WindowDuration:    time.Minute,
	}

	// RateLimitRelaxed - мягкий лимит для read-only операций
	RateLimitRelaxed = RateLimitConfig{
		RequestsPerWindow: 300,
		WindowDuration:    time.Minute,
	}
)

// rateLimitEntry хранит информацию о запросах клиента
type rateLimitEntry struct {
	mu         sync.Mutex
	requests   []time.Time
	lastAccess time.Time
}

// RateLimiter управляет rate limiting
type RateLimiter struct {
	mu      sync.RWMutex
	entries map[string]*rateLimitEntry
	config  RateLimitConfig
}

// newRateLimiter создает новый RateLimiter
func newRateLimiter(config RateLimitConfig) *RateLimiter {
	rl := &RateLimiter{
		entries: make(map[string]*rateLimitEntry),
		config:  config,
	}
	// Запускаем cleanup горутину
	go rl.cleanupLoop()
	return rl
}

// RateLimit создает middleware с заданной конфигурацией
func RateLimit(config RateLimitConfig) func(http.Handler) http.Handler {
	// Создаем отдельный limiter для этого middleware
	limiter := newRateLimiter(config)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// В development режиме можно отключить rate limiting
			env := os.Getenv("ENV")
			isDev := env == "development" || env == "dev" || env == ""

			if isDev {
				// В dev режиме пропускаем все запросы
				next.ServeHTTP(w, r)
				return
			}

			// Получаем идентификатор клиента (IP адрес)
			clientID := getClientIP(r)

			// Проверяем rate limit
			allowed := limiter.Allow(clientID, config)

			if !allowed {
				// Устанавливаем Retry-After header
				retryAfter := int(config.WindowDuration.Seconds())
				w.Header().Set("Retry-After", fmt.Sprintf("%d", retryAfter))
				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", config.RequestsPerWindow))
				w.Header().Set("X-RateLimit-Remaining", "0")

				log.Printf("[RateLimit] Rate limit exceeded for %s on %s %s", clientID, r.Method, r.URL.Path)

				http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
				return
			}

			// Добавляем rate limit headers
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", config.RequestsPerWindow))

			next.ServeHTTP(w, r)
		})
	}
}

// Allow проверяет можно ли пропустить запрос от клиента
func (rl *RateLimiter) Allow(clientID string, config RateLimitConfig) bool {
	now := time.Now()

	// Получаем или создаем entry для клиента
	entry := rl.getOrCreateEntry(clientID)

	entry.mu.Lock()
	defer entry.mu.Unlock()

	// Удаляем запросы вне текущего окна
	windowStart := now.Add(-config.WindowDuration)
	validRequests := []time.Time{}
	for _, reqTime := range entry.requests {
		if reqTime.After(windowStart) {
			validRequests = append(validRequests, reqTime)
		}
	}
	entry.requests = validRequests

	// Проверяем лимит
	if len(entry.requests) >= config.RequestsPerWindow {
		return false
	}

	// Добавляем текущий запрос
	entry.requests = append(entry.requests, now)
	entry.lastAccess = now

	return true
}

// getOrCreateEntry получает или создает entry для клиента
func (rl *RateLimiter) getOrCreateEntry(clientID string) *rateLimitEntry {
	rl.mu.RLock()
	entry, exists := rl.entries[clientID]
	rl.mu.RUnlock()

	if exists {
		return entry
	}

	// Создаем новый entry
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Double-check после получения write lock
	if entry, exists := rl.entries[clientID]; exists {
		return entry
	}

	entry = &rateLimitEntry{
		requests:   make([]time.Time, 0),
		lastAccess: time.Now(),
	}
	rl.entries[clientID] = entry

	return entry
}

// cleanupLoop периодически удаляет устаревшие entries
func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.cleanup()
	}
}

// cleanup удаляет entries которые не использовались последние 10 минут
func (rl *RateLimiter) cleanup() {
	now := time.Now()
	threshold := now.Add(-10 * time.Minute)

	rl.mu.Lock()
	defer rl.mu.Unlock()

	for clientID, entry := range rl.entries {
		entry.mu.Lock()
		if entry.lastAccess.Before(threshold) {
			delete(rl.entries, clientID)
		}
		entry.mu.Unlock()
	}
}

// getClientIP извлекает IP адрес клиента
func getClientIP(r *http.Request) string {
	// Проверяем X-Forwarded-For (если за proxy/load balancer)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Берем первый IP из списка
		// X-Forwarded-For: client, proxy1, proxy2
		return forwarded
	}

	// Проверяем X-Real-IP
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Используем RemoteAddr как fallback
	return r.RemoteAddr
}

// GetRemainingRequests возвращает количество оставшихся запросов для клиента
func (rl *RateLimiter) GetRemainingRequests(clientID string, config RateLimitConfig) int {
	rl.mu.RLock()
	entry, exists := rl.entries[clientID]
	rl.mu.RUnlock()

	if !exists {
		return config.RequestsPerWindow
	}

	entry.mu.Lock()
	defer entry.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-config.WindowDuration)

	// Считаем запросы в текущем окне
	count := 0
	for _, reqTime := range entry.requests {
		if reqTime.After(windowStart) {
			count++
		}
	}

	remaining := config.RequestsPerWindow - count
	if remaining < 0 {
		return 0
	}

	return remaining
}
