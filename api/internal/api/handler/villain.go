package handler

import (
	"context"
	"log"
	"net/http"

	"child-bot/api/internal/api/middleware"
	"child-bot/api/internal/api/response"
	"child-bot/api/internal/api/validation"
	"child-bot/api/internal/service"
)

// VillainServiceInterface интерфейс для VillainService
type VillainServiceInterface interface {
	GetActiveVillain(ctx context.Context, childProfileID string) (*service.Villain, error)
	GetVillainBattle(ctx context.Context, childProfileID, villainID string) (*service.VillainBattle, error)
	GetVillainVictory(ctx context.Context, childProfileID, villainID string) (*service.VictoryData, error)
}

// VillainHandler обрабатывает запросы злодеев
type VillainHandler struct {
	service VillainServiceInterface
}

// NewVillainHandler создает новый VillainHandler
func NewVillainHandler(villainService VillainServiceInterface) *VillainHandler {
	return &VillainHandler{service: villainService}
}

// Villain структура злодея
type Villain struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	HP          int    `json:"hp"`
	MaxHP       int    `json:"max_hp"`
	Level       int    `json:"level"`
	Taunt       string `json:"taunt,omitempty"`
	IsActive    bool   `json:"is_active"`
	IsDefeated  bool   `json:"is_defeated"`
	UnlockedAt  string `json:"unlocked_at,omitempty"`
	DefeatedAt  string `json:"defeated_at,omitempty"`
}

type VillainBattle struct {
	VillainID    string        `json:"villain_id"`
	BattleStats  BattleStats   `json:"battle_stats"`
	RecentDamage []DamageEvent `json:"recent_damage"`
	NextDamageAt string        `json:"next_damage_at,omitempty"`
	CanDamageNow bool          `json:"can_damage_now"`
}

type BattleStats struct {
	TotalDamageDealt  int     `json:"total_damage_dealt"`
	CorrectTasksCount int     `json:"correct_tasks_count"`
	DamagePerTask     int     `json:"damage_per_task"`
	ProgressPercent   float64 `json:"progress_percent"`
}

type DamageEvent struct {
	ID        string `json:"id"`
	Damage    int    `json:"damage"`
	TaskType  string `json:"task_type"` // help, check
	CreatedAt string `json:"created_at"`
}

type VictoryData struct {
	VillainID      string          `json:"villain_id"`
	VillainName    string          `json:"villain_name"`
	DefeatedAt     string          `json:"defeated_at"`
	TotalDamage    int             `json:"total_damage"`
	TasksCompleted int             `json:"tasks_completed"`
	Rewards        []VictoryReward `json:"rewards"`
	NextVillain    *Villain        `json:"next_villain,omitempty"`
}

type VictoryReward struct {
	Type     string `json:"type"` // coins, sticker, avatar, achievement
	ID       string `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url,omitempty"`
	Amount   int    `json:"amount,omitempty"`
}

type DamageRequest struct {
	AttemptID string `json:"attempt_id"`
	Damage    int    `json:"damage"`
}

// List получает список всех злодеев
// GET /villains
func (h *VillainHandler) List(w http.ResponseWriter, r *http.Request) {
	childProfileID := middleware.GetChildProfileID(r.Context())
	if childProfileID == "" {
		response.Unauthorized(w, "Missing child_profile_id")
		return
	}

	// TODO: Phase 4 - получение списка через service layer
	// villains, err := h.service.ListVillains(r.Context(), childProfileID)

	// Placeholder
	villains := []Villain{
		{
			ID:          "villain_1",
			Name:        "Граф Ошибок",
			Description: "Злодей, который распространяет ошибки в задачах",
			ImageURL:    "/assets/villains/count_error.png",
			HP:          100,
			MaxHP:       100,
			Level:       1,
			Taunt:       "Ха-ха! Попробуй-ка реши задачки!",
			IsActive:    true,
			IsDefeated:  false,
		},
	}

	response.OK(w, villains)
}

// GetActive получает активного злодея
// GET /villains/active
func (h *VillainHandler) GetActive(w http.ResponseWriter, r *http.Request) {
	childProfileID := middleware.GetChildProfileID(r.Context())
	if childProfileID == "" {
		response.Unauthorized(w, "Missing child_profile_id")
		return
	}

	// Получаем активного злодея через service layer
	villainData, err := h.service.GetActiveVillain(r.Context(), childProfileID)
	if err != nil {
		log.Printf("[VillainHandler] GetActiveVillain error: %v", err)
		response.InternalError(w, "Failed to get active villain")
		return
	}

	log.Printf("[VillainHandler] GetActiveVillain result: %+v", villainData)

	// Если нет активного злодея
	if villainData == nil {
		response.OK(w, nil)
		return
	}

	// Список реплик для случайного выбора на основе HP
	var taunt string
	healthPercent := float64(villainData.HP) / float64(villainData.MaxHP) * 100.0

	if healthPercent > 75 {
		taunt = "Ха-ха! Попробуй-ка реши задачки!"
	} else if healthPercent > 50 {
		taunt = "Думаешь, справишься?"
	} else if healthPercent > 25 {
		taunt = "Ещё немного, и ты сдашься!"
	} else {
		taunt = "Я непобедим!"
	}

	// Конвертируем в API response
	villain := Villain{
		ID:          villainData.ID,
		Name:        villainData.Name,
		Description: villainData.Description,
		ImageURL:    villainData.ImageURL,
		HP:          villainData.HP,
		MaxHP:       villainData.MaxHP,
		Level:       villainData.Level,
		Taunt:       taunt,
		IsActive:    villainData.IsActive,
		IsDefeated:  villainData.IsDefeated,
	}

	if villainData.UnlockedAt != nil {
		villain.UnlockedAt = villainData.UnlockedAt.Format("2006-01-02T15:04:05Z07:00")
	}
	if villainData.DefeatedAt != nil {
		villain.DefeatedAt = villainData.DefeatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	response.OK(w, villain)
}

// GetByID получает информацию о конкретном злодее
// GET /villains/{id}
func (h *VillainHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	villainID := r.PathValue("id")
	if err := validation.ValidateRequired(villainID, "villain_id"); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	childProfileID := middleware.GetChildProfileID(r.Context())
	if childProfileID == "" {
		response.Unauthorized(w, "Missing child_profile_id")
		return
	}

	// TODO: Phase 4 - получение через service layer
	// villain, err := h.service.GetVillainByID(r.Context(), childProfileID, villainID)

	// Placeholder
	villain := Villain{
		ID:          villainID,
		Name:        "Граф Ошибок",
		Description: "Злодей, который распространяет ошибки в задачах",
		ImageURL:    "/assets/villains/count_error.png",
		HP:          75,
		MaxHP:       100,
		Level:       1,
		Taunt:       "Ха-ха! Попробуй-ка реши задачки!",
		IsActive:    true,
		IsDefeated:  false,
	}

	response.OK(w, villain)
}

// GetBattle получает информацию о битве со злодеем
// GET /villains/{id}/battle
func (h *VillainHandler) GetBattle(w http.ResponseWriter, r *http.Request) {
	villainID := r.PathValue("id")
	if err := validation.ValidateRequired(villainID, "villain_id"); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	childProfileID := middleware.GetChildProfileID(r.Context())
	if childProfileID == "" {
		response.Unauthorized(w, "Missing child_profile_id")
		return
	}

	// Получаем битву через service layer
	battleData, err := h.service.GetVillainBattle(r.Context(), childProfileID, villainID)
	if err != nil {
		response.InternalError(w, "Failed to get villain battle")
		return
	}

	if battleData == nil {
		response.NotFound(w, "Battle not found")
		return
	}

	// Конвертируем события урона
	recentDamage := make([]DamageEvent, 0, len(battleData.RecentDamage))
	for _, event := range battleData.RecentDamage {
		recentDamage = append(recentDamage, DamageEvent{
			ID:        event.ID,
			Damage:    event.Damage,
			TaskType:  event.TaskType,
			CreatedAt: event.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	battle := VillainBattle{
		VillainID: villainID,
		BattleStats: BattleStats{
			TotalDamageDealt:  battleData.BattleStats.TotalDamageDealt,
			CorrectTasksCount: battleData.BattleStats.CorrectTasksCount,
			DamagePerTask:     battleData.BattleStats.DamagePerTask,
			ProgressPercent:   battleData.BattleStats.ProgressPercent,
		},
		RecentDamage: recentDamage,
		CanDamageNow: battleData.CanDamageNow,
	}

	if battleData.NextDamageAt != nil {
		battle.NextDamageAt = battleData.NextDamageAt.Format("2006-01-02T15:04:05Z07:00")
	}

	response.OK(w, battle)
}

// GetVictory получает информацию о победе над злодеем
// GET /villains/{id}/victory
func (h *VillainHandler) GetVictory(w http.ResponseWriter, r *http.Request) {
	villainID := r.PathValue("id")
	if err := validation.ValidateRequired(villainID, "villain_id"); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	childProfileID := middleware.GetChildProfileID(r.Context())
	if childProfileID == "" {
		response.Unauthorized(w, "Missing child_profile_id")
		return
	}

	// TODO: Phase 4 - получение победы через service layer
	// victory, err := h.service.GetVillainVictory(r.Context(), childProfileID, villainID)

	// Placeholder
	victory := VictoryData{
		VillainID:      villainID,
		VillainName:    "Граф Ошибок",
		DefeatedAt:     "2024-03-31T10:00:00Z",
		TotalDamage:    100,
		TasksCompleted: 20,
		Rewards: []VictoryReward{
			{
				Type:   "coins",
				ID:     "coins_100",
				Name:   "100 монет",
				Amount: 100,
			},
			{
				Type:     "achievement",
				ID:       "achievement_villain_1",
				Name:     "Победитель Графа",
				ImageURL: "/assets/achievements/villain_1.png",
			},
		},
	}

	response.OK(w, victory)
}

// DealDamage наносит урон злодею
// POST /villains/{id}/damage
func (h *VillainHandler) DealDamage(w http.ResponseWriter, r *http.Request) {
	villainID := r.PathValue("id")
	if err := validation.ValidateRequired(villainID, "villain_id"); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	childProfileID := middleware.GetChildProfileID(r.Context())
	if childProfileID == "" {
		response.Unauthorized(w, "Missing child_profile_id")
		return
	}

	var req DamageRequest
	if err := validation.DecodeJSON(r, &req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	// TODO: Phase 4 - нанести урон через service layer
	// result, err := h.service.DealDamageToVillain(r.Context(), childProfileID, villainID, req)

	// Placeholder
	result := map[string]interface{}{
		"damage_dealt": 5,
		"villain_hp":   70,
		"is_defeated":  false,
		"message":      "Нанесено 5 урона!",
	}

	response.OK(w, result)
}
