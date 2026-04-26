package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"

	"child-bot/api/internal/api/response"
)

// VKAuthMiddleware проверяет валидность sign параметра от VK Mini Apps
// Документация: https://dev.vk.com/ru/mini-apps/development/launch-parameters
func VKAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Пропускаем валидацию для /analytics и других публичных эндпоинтов
		publicPaths := []string{
			"/analytics/",
			"/health",
			"/avatars",
			"/legal/",
		}
		for _, pp := range publicPaths {
			if strings.HasPrefix(r.URL.Path, pp) {
				next.ServeHTTP(w, r)
				return
			}
		}

		// Проверяем наличие VK параметров
		query := r.URL.Query()
		vkUserID := query.Get("vk_user_id")

		// Пропускаем валидацию для локальной разработки
		if isDevelopmentMode() {
			log.Println("[VK Auth] Development mode: skipping sign validation")
			// Добавляем VK user ID в контекст если есть
			if vkUserID != "" {
				ctx := r.Context()
				ctx = context.WithValue(ctx, ContextKeyVKUserID, vkUserID)
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
			return
		}

		if !hasVKParams(query) {
			// Если нет VK параметров - не VK запрос, пропускаем
			next.ServeHTTP(w, r)
			return
		}

		// Получаем sign из параметров
		sign := query.Get("sign")
		if sign == "" {
			log.Println("[VK Auth] Missing sign parameter")
			response.Unauthorized(w, "Missing VK signature")
			return
		}

		// Валидируем sign
		if !validateVKSign(query, sign) {
			log.Printf("[VK Auth] Invalid VK signature for user: %s", vkUserID)
			response.Unauthorized(w, "Invalid VK signature")
			return
		}

		log.Printf("[VK Auth] Valid VK signature for user: %s", vkUserID)

		// Добавляем VK user ID в контекст
		if vkUserID != "" {
			ctx := r.Context()
			ctx = context.WithValue(ctx, ContextKeyVKUserID, vkUserID)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

// validateVKSign проверяет подпись VK Mini Apps
func validateVKSign(query url.Values, receivedSign string) bool {
	vkAppSecret := os.Getenv("VK_APP_SECRET")
	if vkAppSecret == "" {
		log.Println("[VK Auth] VK_APP_SECRET not configured")
		return false
	}

	// Собираем все vk_* параметры (кроме sign)
	var params []string
	for key := range query {
		if strings.HasPrefix(key, "vk_") {
			value := query.Get(key)
			params = append(params, fmt.Sprintf("%s=%s", key, value))
		}
	}

	// Сортируем параметры по ключу
	sort.Strings(params)

	// Создаем строку для подписи
	queryString := strings.Join(params, "&")

	// Вычисляем HMAC-SHA256
	mac := hmac.New(sha256.New, []byte(vkAppSecret))
	mac.Write([]byte(queryString))
	expectedHash := mac.Sum(nil)

	// Кодируем в base64 и заменяем символы для URL-safe формата
	expectedSign := base64.StdEncoding.EncodeToString(expectedHash)
	expectedSign = strings.ReplaceAll(expectedSign, "+", "-")
	expectedSign = strings.ReplaceAll(expectedSign, "/", "_")
	expectedSign = strings.TrimRight(expectedSign, "=")

	// Сравниваем с полученным sign
	return hmac.Equal([]byte(expectedSign), []byte(receivedSign))
}

// hasVKParams проверяет наличие VK параметров в запросе
func hasVKParams(query url.Values) bool {
	return query.Get("vk_user_id") != "" || query.Get("vk_app_id") != ""
}

// isDevelopmentMode проверяет режим разработки
func isDevelopmentMode() bool {
	env := os.Getenv("ENV")
	return env == "development" || env == "dev" || env == ""
}
