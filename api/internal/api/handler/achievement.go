package handler

import (
	"log"
	"net/http"

	"child-bot/api/internal/api/middleware"
	"child-bot/api/internal/api/response"
	"child-bot/api/internal/api/validation"
	"child-bot/api/internal/store"
)

// AchievementHandler обрабатывает запросы достижений
type AchievementHandler struct {
	store *store.Store
}

// NewAchievementHandler создает новый AchievementHandler
func NewAchievementHandler(store *store.Store) *AchievementHandler {
	return &AchievementHandler{store: store}
}

// toAchievement преобразует CombinedAchievement из store в Achievement для JSON
func toAchievement(ca store.CombinedAchievement) Achievement {
	// Вычисляем процент прогресса
	var percent float64
	if ca.RequirementValue > 0 {
		percent = float64(ca.CurrentProgress) / float64(ca.RequirementValue) * 100
		if percent > 100 {
			percent = 100
		}
	}

	// Формируем unlocked_at
	var unlockedAt string
	if ca.UnlockedAt.Valid {
		unlockedAt = ca.UnlockedAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}

	return Achievement{
		ID:          ca.ID,
		Type:        ca.Type,
		Title:       ca.Title,
		Description: ca.Description,
		Icon:        ca.Icon,
		IsUnlocked:  ca.IsUnlocked,
		UnlockedAt:  unlockedAt,
		Progress: AchievementProgress{
			Current: ca.CurrentProgress,
			Total:   ca.RequirementValue,
			Percent: percent,
		},
		Reward: AchievementReward{
			Type:   ca.RewardType,
			ID:     ca.RewardID.String,
			Name:   ca.RewardName.String,
			Amount: int(ca.RewardAmount.Int32),
		},
		Priority: ca.Priority,
	}
}

// Achievement структура достижения
type Achievement struct {
	ID          string                `json:"id"`
	Type        string                `json:"type"` // streak, tasks, fixes, etc.
	Title       string                `json:"title"`
	Description string                `json:"description"`
	Icon        string                `json:"icon"` // emoji или URL
	IsUnlocked  bool                  `json:"is_unlocked"`
	UnlockedAt  string                `json:"unlocked_at,omitempty"`
	Progress    AchievementProgress   `json:"progress"`
	Reward      AchievementReward     `json:"reward"`
	Priority    int                   `json:"priority"`             // Приоритет для сортировки (меньше = выше)
	NextLevel   *AchievementNextLevel `json:"next_level,omitempty"` // Информация о следующем уровне для серийных наград
}

// AchievementNextLevel информация о следующем уровне серийной награды
type AchievementNextLevel struct {
	Description      string `json:"description"`       // Описание следующего уровня
	RequirementValue int    `json:"requirement_value"` // Необходимое значение для следующего уровня
}

type AchievementProgress struct {
	Current int     `json:"current"`
	Total   int     `json:"total"`
	Percent float64 `json:"percent"`
}

type AchievementReward struct {
	Type     string `json:"type"` // sticker, coins, avatar, badge
	ID       string `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url,omitempty"`
	Amount   int    `json:"amount,omitempty"` // для coins
}

type AchievementsStats struct {
	UnlockedCount   int     `json:"unlocked_count"`
	TotalCount      int     `json:"total_count"`
	ProgressPercent float64 `json:"progress_percent"`
}

// List получает список всех достижений
// GET /achievements
func (h *AchievementHandler) List(w http.ResponseWriter, r *http.Request) {
	childProfileID := middleware.GetChildProfileID(r.Context())
	if childProfileID == "" {
		response.Unauthorized(w, "Missing child_profile_id")
		return
	}

	// Получаем достижения из базы данных
	combined, err := h.store.ListAchievementsWithProgress(r.Context(), childProfileID)
	if err != nil {
		response.InternalError(w, "Failed to load achievements")
		return
	}

	// Получаем информацию о следующих уровнях для серийных достижений
	nextLevels, err := h.store.GetNextLevelsForSeries(r.Context(), childProfileID)
	if err != nil {
		log.Printf("[AchievementHandler] Failed to get next levels: %v", err)
		// Не критично, продолжаем без next_level
	}

	// Преобразуем в формат для JSON
	achievements := make([]Achievement, 0, len(combined))
	for _, ca := range combined {
		ach := toAchievement(ca)

		// Добавляем информацию о следующем уровне если есть
		if nextLevel, ok := nextLevels[ca.ID]; ok {
			ach.NextLevel = &AchievementNextLevel{
				Description:      nextLevel.Description,
				RequirementValue: nextLevel.RequirementValue,
			}
		}

		achievements = append(achievements, ach)
	}

	response.OK(w, achievements)
}

// ListOld - старый код с hardcoded данными (удалить после проверки)
func (h *AchievementHandler) ListOld(w http.ResponseWriter, r *http.Request) {
	childProfileID := middleware.GetChildProfileID(r.Context())
	if childProfileID == "" {
		response.Unauthorized(w, "Missing child_profile_id")
		return
	}

	// Placeholder - мок данные из дизайна
	achievements := []Achievement{
		// Первая полка - разблокированные
		{
			ID:          "achievement_1",
			Type:        "streak",
			Title:       "5 дней подряд",
			Description: "Решай задачи 5 дней подряд",
			Icon:        "🔥",
			IsUnlocked:  true,
			UnlockedAt:  "2024-03-01T10:00:00Z",
			Progress: AchievementProgress{
				Current: 5,
				Total:   5,
				Percent: 100,
			},
			Reward: AchievementReward{
				Type:   "coins",
				ID:     "coins_50",
				Name:   "50 монет",
				Amount: 50,
			},
			Priority: 1,
		},
		{
			ID:          "achievement_2",
			Type:        "tasks",
			Title:       "10 проверок ДЗ",
			Description: "Проверь 10 домашних заданий",
			Icon:        "✅",
			IsUnlocked:  true,
			UnlockedAt:  "2024-03-02T10:00:00Z",
			Progress: AchievementProgress{
				Current: 10,
				Total:   10,
				Percent: 100,
			},
			Reward: AchievementReward{
				Type:   "coins",
				ID:     "coins_100",
				Name:   "100 монет",
				Amount: 100,
			},
			Priority: 1,
		},
		{
			ID:          "achievement_3",
			Type:        "fixes",
			Title:       "5 ошибок исправлено",
			Description: "Исправь 5 ошибок в заданиях",
			Icon:        "⭐",
			IsUnlocked:  true,
			UnlockedAt:  "2024-03-03T10:00:00Z",
			Progress: AchievementProgress{
				Current: 5,
				Total:   5,
				Percent: 100,
			},
			Reward: AchievementReward{
				Type:   "coins",
				ID:     "coins_75",
				Name:   "75 монет",
				Amount: 75,
			},
			Priority: 1,
		},
		{
			ID:          "achievement_4",
			Type:        "first",
			Title:       "Первое задание",
			Description: "Реши первое задание",
			Icon:        "🎯",
			IsUnlocked:  false,
			Progress: AchievementProgress{
				Current: 0,
				Total:   1,
				Percent: 0,
			},
			Reward: AchievementReward{
				Type:   "coins",
				ID:     "coins_25",
				Name:   "25 монет",
				Amount: 25,
			},
			Priority: 1,
		},
		// Вторая полка - заблокированные
		{
			ID:          "achievement_5",
			Type:        "speed",
			Title:       "Скоростной решатель",
			Description: "Реши задачу за 5 минут",
			Icon:        "⚡",
			IsUnlocked:  false,
			Progress: AchievementProgress{
				Current: 0,
				Total:   1,
				Percent: 0,
			},
			Reward: AchievementReward{
				Type:   "coins",
				ID:     "coins_150",
				Name:   "150 монет",
				Amount: 150,
			},
			Priority: 2,
		},
		{
			ID:          "achievement_6",
			Type:        "villain",
			Title:       "Победитель злодеев",
			Description: "Победи 3 злодеев",
			Icon:        "🏆",
			IsUnlocked:  false,
			Progress: AchievementProgress{
				Current: 0,
				Total:   3,
				Percent: 0,
			},
			Reward: AchievementReward{
				Type: "sticker",
				ID:   "sticker_trophy",
				Name: "Стикер Чемпиона",
			},
			Priority: 2,
		},
		{
			ID:          "achievement_7",
			Type:        "wisdom",
			Title:       "Мудрая сова",
			Description: "Получи 50 подсказок",
			Icon:        "🦉",
			IsUnlocked:  false,
			Progress: AchievementProgress{
				Current: 0,
				Total:   50,
				Percent: 0,
			},
			Reward: AchievementReward{
				Type:   "coins",
				ID:     "coins_200",
				Name:   "200 монет",
				Amount: 200,
			},
			Priority: 2,
		},
		{
			ID:          "achievement_8",
			Type:        "collection",
			Title:       "Коллекционер",
			Description: "Собери все стикеры",
			Icon:        "💎",
			IsUnlocked:  false,
			Progress: AchievementProgress{
				Current: 0,
				Total:   20,
				Percent: 0,
			},
			Reward: AchievementReward{
				Type: "badge",
				ID:   "badge_collector",
				Name: "Значок Коллекционера",
			},
			Priority: 2,
		},
		// Третья полка
		{
			ID:          "achievement_9",
			Type:        "rocket",
			Title:       "Ракета знаний",
			Description: "Реши 100 задач",
			Icon:        "🚀",
			IsUnlocked:  false,
			Progress: AchievementProgress{
				Current: 0,
				Total:   100,
				Percent: 0,
			},
			Reward: AchievementReward{
				Type:   "coins",
				ID:     "coins_500",
				Name:   "500 монет",
				Amount: 500,
			},
			Priority: 3,
		},
		{
			ID:          "achievement_10",
			Type:        "superhero",
			Title:       "Супергерой",
			Description: "Помоги друзьям 10 раз",
			Icon:        "🦸",
			IsUnlocked:  false,
			Progress: AchievementProgress{
				Current: 0,
				Total:   10,
				Percent: 0,
			},
			Reward: AchievementReward{
				Type: "avatar",
				ID:   "avatar_superhero",
				Name: "Аватар Супергероя",
			},
			Priority: 3,
		},
		{
			ID:          "achievement_11",
			Type:        "marathon",
			Title:       "Марафонец",
			Description: "Учись 30 дней подряд",
			Icon:        "🏛️",
			IsUnlocked:  false,
			Progress: AchievementProgress{
				Current: 0,
				Total:   30,
				Percent: 0,
			},
			Reward: AchievementReward{
				Type:   "coins",
				ID:     "coins_1000",
				Name:   "1000 монет",
				Amount: 1000,
			},
			Priority: 3,
		},
		{
			ID:          "achievement_12",
			Type:        "genius",
			Title:       "Гений",
			Description: "Реши все задачи без ошибок",
			Icon:        "🧠",
			IsUnlocked:  false,
			Progress: AchievementProgress{
				Current: 0,
				Total:   50,
				Percent: 0,
			},
			Reward: AchievementReward{
				Type: "badge",
				ID:   "badge_genius",
				Name: "Значок Гения",
			},
			Priority: 3,
		},
	}

	response.OK(w, achievements)
}

// GetUnlocked получает только разблокированные достижения
// GET /achievements/unlocked
func (h *AchievementHandler) GetUnlocked(w http.ResponseWriter, r *http.Request) {
	childProfileID := middleware.GetChildProfileID(r.Context())
	if childProfileID == "" {
		response.Unauthorized(w, "Missing child_profile_id")
		return
	}

	// TODO: Phase 4 - получение через service layer
	// unlocked, err := h.service.GetUnlockedAchievements(r.Context(), childProfileID)

	// Placeholder
	unlocked := []Achievement{}

	response.OK(w, unlocked)
}

// GetByID получает информацию о конкретном достижении
// GET /achievements/{id}
func (h *AchievementHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	achievementID := r.PathValue("id")
	if err := validation.ValidateRequired(achievementID, "achievement_id"); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	childProfileID := middleware.GetChildProfileID(r.Context())
	if childProfileID == "" {
		response.Unauthorized(w, "Missing child_profile_id")
		return
	}

	// Получаем достижение из базы данных
	combined, err := h.store.GetAchievementByID(r.Context(), childProfileID, achievementID)
	if err != nil {
		response.InternalError(w, "Failed to load achievement")
		return
	}
	if combined == nil {
		response.NotFound(w, "Achievement not found")
		return
	}

	achievement := toAchievement(*combined)
	response.OK(w, achievement)
}

// Claim забирает награду за достижение
// POST /achievements/{id}/claim
func (h *AchievementHandler) Claim(w http.ResponseWriter, r *http.Request) {
	achievementID := r.PathValue("id")
	if err := validation.ValidateRequired(achievementID, "achievement_id"); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	childProfileID := middleware.GetChildProfileID(r.Context())
	if childProfileID == "" {
		response.Unauthorized(w, "Missing child_profile_id")
		return
	}

	// TODO: Phase 4 - забрать награду через service layer
	// reward, err := h.service.ClaimAchievement(r.Context(), childProfileID, achievementID)

	// Placeholder
	reward := map[string]interface{}{
		"claimed":     true,
		"reward_type": "coins",
		"amount":      50,
		"message":     "Получено 50 монет!",
	}

	response.OK(w, reward)
}

// GetStats получает статистику достижений
// GET /achievements/stats
func (h *AchievementHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	childProfileID := middleware.GetChildProfileID(r.Context())
	if childProfileID == "" {
		response.Unauthorized(w, "Missing child_profile_id")
		return
	}

	// Получаем статистику из базы данных
	unlockedCount, totalCount, err := h.store.GetAchievementStats(r.Context(), childProfileID)
	if err != nil {
		response.InternalError(w, "Failed to load stats")
		return
	}

	// Вычисляем процент прогресса
	var progressPercent float64
	if totalCount > 0 {
		progressPercent = float64(unlockedCount) / float64(totalCount) * 100
	}

	stats := AchievementsStats{
		UnlockedCount:   unlockedCount,
		TotalCount:      totalCount,
		ProgressPercent: progressPercent,
	}

	response.OK(w, stats)
}

// HasNew проверяет есть ли новые (непросмотренные) достижения
// GET /achievements/has-new
func (h *AchievementHandler) HasNew(w http.ResponseWriter, r *http.Request) {
	childProfileID := middleware.GetChildProfileID(r.Context())
	if childProfileID == "" {
		response.Unauthorized(w, "Missing child_profile_id")
		return
	}

	hasNew, err := h.store.HasNewAchievements(r.Context(), childProfileID)
	if err != nil {
		response.InternalError(w, "Failed to check new achievements")
		return
	}

	response.OK(w, map[string]interface{}{
		"has_new": hasNew,
	})
}

// MarkViewed отмечает что пользователь просмотрел страницу достижений
// POST /achievements/mark-viewed
func (h *AchievementHandler) MarkViewed(w http.ResponseWriter, r *http.Request) {
	childProfileID := middleware.GetChildProfileID(r.Context())
	if childProfileID == "" {
		response.Unauthorized(w, "Missing child_profile_id")
		return
	}

	err := h.store.MarkAchievementsViewed(r.Context(), childProfileID)
	if err != nil {
		response.InternalError(w, "Failed to mark achievements viewed")
		return
	}

	response.OK(w, map[string]interface{}{
		"success": true,
	})
}
