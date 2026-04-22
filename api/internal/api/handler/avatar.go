package handler

import (
	"net/http"

	"child-bot/api/internal/api/response"
)

// AvatarHandler обрабатывает запросы аватаров
type AvatarHandler struct{}

// NewAvatarHandler создает новый AvatarHandler
func NewAvatarHandler() *AvatarHandler {
	return &AvatarHandler{}
}

// Avatar представляет доступный аватар
type Avatar struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ImageURL  string `json:"imageUrl"`
	IsPremium bool   `json:"isPremium"`
}

// GetAll возвращает список доступных аватаров
// GET /avatars
func (h *AvatarHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Список доступных аватаров (пока hardcoded)
	avatars := []Avatar{
		{ID: "cat", Name: "Кот", ImageURL: "🐱", IsPremium: false},
		{ID: "dog", Name: "Пёс", ImageURL: "🐶", IsPremium: false},
		{ID: "panda", Name: "Панда", ImageURL: "🐼", IsPremium: false},
		{ID: "fox", Name: "Лиса", ImageURL: "🦊", IsPremium: false},
		{ID: "bear", Name: "Медведь", ImageURL: "🐻", IsPremium: false},
		{ID: "lion", Name: "Лев", ImageURL: "🦁", IsPremium: false},
		{ID: "tiger", Name: "Тигр", ImageURL: "🐯", IsPremium: true},
		{ID: "unicorn", Name: "Единорог", ImageURL: "🦄", IsPremium: true},
		{ID: "robot", Name: "Робот", ImageURL: "🤖", IsPremium: false},
		{ID: "alien", Name: "Пришелец", ImageURL: "👽", IsPremium: true},
	}

	response.OK(w, avatars)
}
