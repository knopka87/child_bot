package service

import (
	"context"
	"time"

	"child-bot/api/internal/store"
)

// VillainService бизнес-логика для злодеев
type VillainService struct {
	store *store.Store
}

// NewVillainService создает новый VillainService
func NewVillainService(store *store.Store) *VillainService {
	return &VillainService{store: store}
}

// Villain злодей
type Villain struct {
	ID          string
	Name        string
	Description string
	ImageURL    string
	HP          int
	MaxHP       int
	Level       int
	IsActive    bool
	IsDefeated  bool
	UnlockedAt  *time.Time
	DefeatedAt  *time.Time
}

// VillainBattle данные битвы
type VillainBattle struct {
	VillainID    string
	BattleStats  BattleStats
	RecentDamage []DamageEvent
	NextDamageAt *time.Time
	CanDamageNow bool
}

// BattleStats статистика битвы
type BattleStats struct {
	TotalDamageDealt  int
	CorrectTasksCount int
	DamagePerTask     int
	ProgressPercent   float64
}

// DamageEvent урон злодею
type DamageEvent struct {
	ID        string
	Damage    int
	TaskType  string // help, check
	CreatedAt time.Time
}

// VictoryData данные победы
type VictoryData struct {
	VillainID      string
	VillainName    string
	DefeatedAt     time.Time
	TotalDamage    int
	TasksCompleted int
	Rewards        []VictoryReward
	NextVillain    *Villain
}

// VictoryReward награда за победу
type VictoryReward struct {
	Type     string // coins, sticker, avatar, achievement
	ID       string
	Name     string
	ImageURL string
	Amount   int
}

// ListVillains получает список злодеев
func (s *VillainService) ListVillains(ctx context.Context, childProfileID string) ([]Villain, error) {
	// TODO: Phase 5 - загрузка из БД

	villains := []Villain{
		{
			ID:          "villain_1",
			Name:        "Граф Ошибок",
			Description: "Злодей, который распространяет ошибки в задачах",
			ImageURL:    "/assets/villains/count_error.png",
			HP:          100,
			MaxHP:       100,
			Level:       1,
			IsActive:    true,
			IsDefeated:  false,
		},
	}

	return villains, nil
}

// GetActiveVillain получает активного злодея
func (s *VillainService) GetActiveVillain(ctx context.Context, childProfileID string) (*Villain, error) {
	battle, villainRow, err := s.store.Villains.GetActiveVillainBattle(ctx, childProfileID)
	if err != nil {
		return nil, err
	}

	// Если нет активной битвы
	if battle == nil || villainRow == nil {
		return nil, nil
	}

	villain := &Villain{
		ID:          villainRow.ID,
		Name:        villainRow.Name,
		Description: villainRow.Description,
		ImageURL:    villainRow.ImageURL,
		HP:          battle.CurrentHP,
		MaxHP:       villainRow.MaxHP,
		Level:       villainRow.Level,
		IsActive:    battle.Status == "active",
		IsDefeated:  battle.Status == "defeated",
		UnlockedAt:  &battle.StartedAt,
		DefeatedAt:  nil,
	}

	if battle.DefeatedAt.Valid {
		villain.DefeatedAt = &battle.DefeatedAt.Time
	}

	return villain, nil
}

// GetVillainByID получает злодея по ID
func (s *VillainService) GetVillainByID(ctx context.Context, childProfileID, villainID string) (*Villain, error) {
	// TODO: Phase 5 - загрузка из БД

	return s.GetActiveVillain(ctx, childProfileID)
}

// GetVillainBattle получает информацию о битве
func (s *VillainService) GetVillainBattle(ctx context.Context, childProfileID, villainID string) (*VillainBattle, error) {
	battleRow, villainRow, err := s.store.Villains.GetActiveVillainBattle(ctx, childProfileID)
	if err != nil {
		return nil, err
	}

	if battleRow == nil || villainRow == nil {
		return nil, nil
	}

	// Загружаем последние события урона
	damageEvents, err := s.store.Villains.GetDamageEvents(ctx, battleRow.ID, 10)
	if err != nil {
		return nil, err
	}

	// Конвертируем в доменные объекты
	recentDamage := make([]DamageEvent, 0, len(damageEvents))
	for _, event := range damageEvents {
		recentDamage = append(recentDamage, DamageEvent{
			ID:        string(rune(event.ID)),
			Damage:    event.Damage,
			TaskType:  event.TaskType,
			CreatedAt: event.CreatedAt,
		})
	}

	// Вычисляем прогресс
	progressPercent := 0.0
	if villainRow.MaxHP > 0 {
		progressPercent = float64(battleRow.TotalDamageDealt) / float64(villainRow.MaxHP) * 100.0
	}

	battle := &VillainBattle{
		VillainID: villainID,
		BattleStats: BattleStats{
			TotalDamageDealt:  battleRow.TotalDamageDealt,
			CorrectTasksCount: battleRow.CorrectTasksCount,
			DamagePerTask:     villainRow.DamagePerCorrectTask,
			ProgressPercent:   progressPercent,
		},
		RecentDamage: recentDamage,
		CanDamageNow: battleRow.Status == "active",
	}

	return battle, nil
}

// DealDamage наносит урон злодею
func (s *VillainService) DealDamage(ctx context.Context, childProfileID, villainID string, attemptID string, damage int) (*DamageResult, error) {
	// TODO: Phase 5 - обновление HP злодея в БД

	result := &DamageResult{
		DamageDealt: damage,
		VillainHP:   70,
		VillainMaxHP: 100,
		IsDefeated:  false,
	}

	return result, nil
}

// GetVillainVictory получает информацию о победе
func (s *VillainService) GetVillainVictory(ctx context.Context, childProfileID, villainID string) (*VictoryData, error) {
	// TODO: Phase 5 - загрузка данных победы из БД

	victory := &VictoryData{
		VillainID:      villainID,
		VillainName:    "Граф Ошибок",
		DefeatedAt:     time.Now(),
		TotalDamage:    100,
		TasksCompleted: 20,
		Rewards: []VictoryReward{
			{
				Type:   "coins",
				ID:     "coins_100",
				Name:   "100 монет",
				Amount: 100,
			},
		},
	}

	return victory, nil
}

// DamageResult результат нанесения урона
type DamageResult struct {
	DamageDealt  int
	VillainHP    int
	VillainMaxHP int
	IsDefeated   bool
	Rewards      []VictoryReward
}
