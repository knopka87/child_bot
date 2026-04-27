package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"child-bot/api/internal/api/response"
)

type contextKey string

const (
	// ContextKeyPlatformID ключ для platformID в context
	ContextKeyPlatformID contextKey = "platformID"
	// ContextKeyChildProfileID ключ для childProfileID в context
	ContextKeyChildProfileID contextKey = "childProfileID"
	// ContextKeyVKUserID ключ для VK user ID в context
	ContextKeyVKUserID contextKey = "vkUserID"
)

// Auth middleware для проверки платформы и профиля
// Ожидает заголовки:
// - X-Platform-ID: vk|telegram|max|web
// - X-Child-Profile-ID: uuid профиля ребенка
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		platformID := strings.TrimSpace(r.Header.Get("X-Platform-ID"))
		childProfileID := strings.TrimSpace(r.Header.Get("X-Child-Profile-ID"))

		log.Printf("[Auth] Request: %s %s, Platform-ID: '%s', Child-Profile-ID: '%s'",
			r.Method, r.URL.Path, platformID, childProfileID)
		// Также выведем в stdout для docker logs
		fmt.Printf("[Auth] Request: %s %s, Platform-ID: '%s', Child-Profile-ID: '%s'\n",
			r.Method, r.URL.Path, platformID, childProfileID)

		// Для некоторых endpoints (health, onboarding) auth не требуется
		// Проверим, нужна ли аутентификация для этого пути
		if !requiresAuth(r.URL.Path) {
			log.Printf("[Auth] Skipping auth for path: %s", r.URL.Path)
			next.ServeHTTP(w, r)
			return
		}

		// Проверка platformID
		if platformID == "" {
			log.Printf("[Auth] Missing X-Platform-ID header for path: %s", r.URL.Path)
			response.Unauthorized(w, "Missing X-Platform-ID header")
			return
		}

		if !isValidPlatform(platformID) {
			log.Printf("[Auth] Invalid platform ID: %s for path: %s", platformID, r.URL.Path)
			response.BadRequest(w, "Invalid platform ID")
			return
		}

		// Проверка childProfileID (для большинства endpoints)
		if requiresChildProfile(r.URL.Path) && childProfileID == "" {
			response.Unauthorized(w, "Missing X-Child-Profile-ID header")
			return
		}

		// Добавляем данные в context
		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextKeyPlatformID, platformID)
		if childProfileID != "" {
			ctx = context.WithValue(ctx, ContextKeyChildProfileID, childProfileID)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetPlatformID извлекает platformID из context
func GetPlatformID(ctx context.Context) string {
	if id, ok := ctx.Value(ContextKeyPlatformID).(string); ok {
		return id
	}
	return ""
}

// GetChildProfileID извлекает childProfileID из context
func GetChildProfileID(ctx context.Context) string {
	if id, ok := ctx.Value(ContextKeyChildProfileID).(string); ok {
		return id
	}
	return ""
}

// GetVKUserID извлекает VK user ID из context
func GetVKUserID(ctx context.Context) string {
	if id, ok := ctx.Value(ContextKeyVKUserID).(string); ok {
		return id
	}
	return ""
}

// requiresAuth проверяет, требует ли путь аутентификации
func requiresAuth(path string) bool {
	publicPaths := []string{
		"/health",
		"/api/health",
		"/onboarding/start",
		"/api/onboarding/start",
		"/onboarding/complete",
		"/api/onboarding/complete",
		"/avatars", // Временно: для отладки VK Mini App
		"/api/avatars",
		"/analytics/events", // Analytics events (могут отправляться без childProfileID)
		"/api/analytics/events",
		"/legal/", // Legal documents доступны всем
		"/api/legal/",
		"/webhooks/", // Webhooks (VK Pay, etc.) - используют собственную валидацию
	}

	for _, pp := range publicPaths {
		if strings.HasPrefix(path, pp) {
			return false
		}
	}
	return true
}

// requiresChildProfile проверяет, требует ли путь childProfileID
func requiresChildProfile(path string) bool {
	// Paths, которые НЕ требуют childProfileID (например, onboarding)
	noProfilePaths := []string{
		"/onboarding/",
		"/avatars",
		"/api/analytics/events",
		"/profiles/child", // Создание профиля - не требует ID, так как он ещё не создан
		"/api/profiles/child",
		"/profiles/by-platform", // Получение профиля по platform credentials - используется для auth
		"/api/profiles/by-platform",
		"/consent", // Сохранение согласия - часть onboarding
		"/api/consent",
		"/email/", // Email verification - часть onboarding, до создания профиля
	}

	for _, npp := range noProfilePaths {
		if strings.HasPrefix(path, npp) {
			return false
		}
	}
	return true
}

// isValidPlatform проверяет валидность platformID
func isValidPlatform(platform string) bool {
	validPlatforms := map[string]bool{
		"vk":       true,
		"telegram": true,
		"max":      true,
		"web":      true,
	}
	return validPlatforms[platform]
}
