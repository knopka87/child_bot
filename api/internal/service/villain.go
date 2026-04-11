package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"child-bot/api/internal/store"

	"github.com/google/uuid"
)

// VillainService бизнес-логика для злодеев
type VillainService struct {
	store              *store.Store
	achievementService *AchievementService
}

// NewVillainService создает новый VillainService
func NewVillainService(store *store.Store) *VillainService {
	return &VillainService{store: store}
}

// SetAchievementService устанавливает AchievementService (для избежания циклических зависимостей)
func (s *VillainService) SetAchievementService(achievementService *AchievementService) {
	s.achievementService = achievementService
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

// DealDamageToVillain наносит урон активному злодею и проверяет победу
// Возвращает: (defeated bool, coinsEarned int, error)
func (s *VillainService) DealDamageToVillain(ctx context.Context, childProfileID string, attemptID uuid.UUID, taskType string) (bool, int, error) {
	// Получаем активную битву
	battle, villainRow, err := s.store.Villains.GetActiveVillainBattle(ctx, childProfileID)
	if err != nil {
		return false, 0, fmt.Errorf("failed to get active battle: %w", err)
	}

	if battle == nil || villainRow == nil {
		// Нет активной битвы - создаём первого монстра
		log.Printf("[VillainService] No active battle for %s, creating first villain", childProfileID)
		err = s.ensureActiveVillain(ctx, childProfileID)
		if err != nil {
			return false, 0, fmt.Errorf("failed to create first villain: %w", err)
		}
		// Повторно получаем битву
		battle, villainRow, err = s.store.Villains.GetActiveVillainBattle(ctx, childProfileID)
		if err != nil || battle == nil {
			return false, 0, fmt.Errorf("failed to get battle after creation: %w", err)
		}
	}

	// Вычисляем урон
	damage := villainRow.DamagePerCorrectTask
	log.Printf("[VillainService] Dealing %d damage to villain %s (current HP: %d/%d)",
		damage, villainRow.ID, battle.CurrentHP, villainRow.MaxHP)

	// Записываем событие урона
	err = s.store.Villains.RecordDamageEvent(ctx, battle.ID, attemptID, damage, taskType)
	if err != nil {
		log.Printf("[VillainService] Failed to record damage event: %v", err)
	}

	// Обновляем HP и счётчики битвы
	newHP := battle.CurrentHP - damage
	if newHP < 0 {
		newHP = 0
	}

	err = s.store.Villains.UpdateBattleProgress(ctx, battle.ID, newHP, damage)
	if err != nil {
		return false, 0, fmt.Errorf("failed to update battle progress: %w", err)
	}

	// Проверяем победу
	defeated := newHP <= 0
	coinsEarned := 0

	if defeated {
		log.Printf("[VillainService] Villain %s defeated! Awarding coins: %d", villainRow.ID, villainRow.RewardCoins)

		// Помечаем битву как побеждённую
		err = s.store.Villains.MarkBattleDefeated(ctx, battle.ID)
		if err != nil {
			log.Printf("[VillainService] Failed to mark battle as defeated: %v", err)
		}

		// Начисляем монеты за победу
		coinsEarned = villainRow.RewardCoins
		// Монеты начислит вызывающий код

		// Создаём следующего монстра
		err = s.createNextVillain(ctx, childProfileID, villainRow.UnlockOrder)
		if err != nil {
			log.Printf("[VillainService] Failed to create next villain: %v", err)
		}

		// Проверяем достижения за побеждённых монстров
		if s.achievementService != nil {
			err = s.achievementService.CheckVillainAchievements(ctx, childProfileID)
			if err != nil {
				log.Printf("[VillainService] Failed to check villain achievements: %v", err)
				// Не блокируем, продолжаем
			}
		}
	}

	return defeated, coinsEarned, nil
}

// ensureActiveVillain создаёт первого монстра если нет активного
func (s *VillainService) ensureActiveVillain(ctx context.Context, childProfileID string) error {
	// Получаем первого злодея (unlock_order = 1)
	villain, err := s.store.Villains.GetVillainByOrder(ctx, 1)
	if err != nil {
		return fmt.Errorf("failed to get first villain: %w", err)
	}

	if villain == nil {
		return fmt.Errorf("no villains found in database")
	}

	// Создаём битву
	err = s.store.Villains.CreateBattle(ctx, childProfileID, villain.ID, villain.MaxHP)
	if err != nil {
		return fmt.Errorf("failed to create battle: %w", err)
	}

	log.Printf("[VillainService] Created first villain battle: %s for child %s", villain.ID, childProfileID)
	return nil
}

// createNextVillain создаёт следующего монстра после победы
func (s *VillainService) createNextVillain(ctx context.Context, childProfileID string, currentOrder int) error {
	const maxAttempts = 20 // Защита от бесконечного цикла
	attempts := 0

	for attempts < maxAttempts {
		attempts++

		// Получаем следующего злодея
		nextOrder := currentOrder + 1
		villain, err := s.store.Villains.GetVillainByOrder(ctx, nextOrder)
		if err != nil {
			return fmt.Errorf("failed to get next villain: %w", err)
		}

		// Если нет следующего - начинаем цикл заново
		if villain == nil {
			log.Printf("[VillainService] No more villains, starting cycle again")
			villain, err = s.store.Villains.GetVillainByOrder(ctx, 1)
			if err != nil || villain == nil {
				return fmt.Errorf("failed to get first villain for cycle: %w", err)
			}
			currentOrder = 0 // Сброс для следующей итерации
		} else {
			currentOrder = villain.UnlockOrder
		}

		// Проверяем: если это босс, можем ли мы его создать?
		if villain.IsBoss {
			lastBossDate, err := s.store.Villains.GetLastBossDefeatedAt(ctx, childProfileID)
			if err != nil {
				log.Printf("[VillainService] Failed to check last boss date: %v", err)
				// Продолжаем, считаем что босса можно создать
			} else if lastBossDate != nil {
				// Проверяем прошло ли 7 дней
				daysSinceLastBoss := time.Since(*lastBossDate).Hours() / 24
				if daysSinceLastBoss < 7 {
					log.Printf("[VillainService] Boss %s skipped, last boss was %.1f days ago (need 7)",
						villain.ID, daysSinceLastBoss)
					// Пропускаем босса, пробуем следующего
					continue
				}
			}
		}

		// Создаём новую битву
		err = s.store.Villains.CreateBattle(ctx, childProfileID, villain.ID, villain.MaxHP)
		if err != nil {
			return fmt.Errorf("failed to create next battle: %w", err)
		}

		log.Printf("[VillainService] Created next villain battle: %s (order %d, is_boss: %v) for child %s",
			villain.ID, villain.UnlockOrder, villain.IsBoss, childProfileID)
		return nil
	}

	return fmt.Errorf("failed to find suitable villain after %d attempts", maxAttempts)
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
		DamageDealt:  damage,
		VillainHP:    70,
		VillainMaxHP: 100,
		IsDefeated:   false,
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
