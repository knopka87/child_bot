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

// ListVillains получает список всех злодеев из БД
func (s *VillainService) ListVillains(ctx context.Context, childProfileID string) ([]Villain, error) {
	// Получаем всех злодеев из справочника
	villains := []Villain{}

	for i := 1; i <= 20; i++ { // Загружаем до 20 злодеев
		v, err := s.store.Villains.GetVillainByOrder(ctx, i)
		if err != nil {
			log.Printf("[VillainService] Failed to get villain %d: %v", i, err)
			continue
		}
		if v == nil {
			break // Больше злодеев нет
		}

		// Получаем статус битвы для этого профиля
		battle, _, _ := s.store.Villains.GetVillainBattleByVillainID(ctx, childProfileID, v.ID)

		hp := v.MaxHP
		isActive := false
		isDefeated := false
		var unlockedAt *time.Time

		if battle != nil {
			hp = battle.CurrentHP
			isActive = battle.Status == "active"
			isDefeated = battle.Status == "defeated"
			unlockedAt = &battle.StartedAt
		}

		villains = append(villains, Villain{
			ID:          v.ID,
			Name:        v.Name,
			Description: v.Description,
			ImageURL:    v.ImageURL,
			HP:          hp,
			MaxHP:       v.MaxHP,
			Level:       v.Level,
			IsActive:    isActive,
			IsDefeated:  isDefeated,
			UnlockedAt:  unlockedAt,
		})
	}

	return villains, nil
}

// GetActiveVillain получает активного злодея
func (s *VillainService) GetActiveVillain(ctx context.Context, childProfileID string) (*Villain, error) {
	// Получаем последнюю битву (активную или побеждённую сегодня)
	battle, villainRow, err := s.store.Villains.GetActiveVillainBattle(ctx, childProfileID)
	if err != nil {
		return nil, err
	}

	// Определяем какой злодей должен быть сегодня
	today := time.Now().UTC()
	dayOfWeek := int(today.Weekday()) // Sunday=0, Monday=1, ..., Saturday=6
	if dayOfWeek == 0 {
		dayOfWeek = 7 // Sunday = 7
	}

	// Если нет активной битвы ИЛИ злодей не соответствует дню недели
	shouldCreateNew := false
	if battle == nil || villainRow == nil {
		shouldCreateNew = true
	} else {
		// Проверяем соответствует ли злодей сегодняшнему дню
		correctVillain, err := s.store.Villains.GetVillainByOrder(ctx, dayOfWeek)
		if err != nil {
			log.Printf("[VillainService] Failed to get villain for day %d: %v", dayOfWeek, err)
		} else if correctVillain != nil && correctVillain.ID != villainRow.ID {
			// Злодей не соответствует дню - помечаем старую битву как заброшенную
			log.Printf("[VillainService] Villain mismatch: current=%s, should be=%s for day %d. Replacing...",
				villainRow.ID, correctVillain.ID, dayOfWeek)

			_ = s.store.Villains.MarkBattleAbandoned(ctx, battle.ID)
			shouldCreateNew = true
		}
	}

	if shouldCreateNew {
		// Проверяем, был ли злодей побеждён сегодня (тогда не создаём нового)
		defeatedToday, err := s.wasDefeatedToday(ctx, childProfileID)
		if err != nil {
			log.Printf("[VillainService] Failed to check if defeated today: %v", err)
		}

		if defeatedToday {
			log.Printf("[VillainService] Villain already defeated today for %s", childProfileID)
			return nil, nil
		}

		// Создаём нового злодея на сегодня
		log.Printf("[VillainService] Creating new villain for day %d (%s) for %s",
			dayOfWeek, today.Weekday().String(), childProfileID)
		err = s.ensureDailyVillain(ctx, childProfileID)
		if err != nil {
			return nil, fmt.Errorf("failed to create daily villain: %w", err)
		}
		// Повторно получаем битву
		battle, villainRow, err = s.store.Villains.GetActiveVillainBattle(ctx, childProfileID)
		if err != nil || battle == nil {
			return nil, fmt.Errorf("failed to get battle after creation: %w", err)
		}
	}

	// Проверяем, что битва начата сегодня (для сброса HP)
	battleDate := battle.StartedAt.Truncate(24 * time.Hour)
	todayDate := today.Truncate(24 * time.Hour)

	if battleDate.Before(todayDate) {
		// Битва начата не сегодня - сбрасываем HP
		log.Printf("[VillainService] Battle from %s, resetting HP for today", battle.StartedAt.Format("2006-01-02"))

		if battle.Status == "active" {
			err = s.store.Villains.ResetBattleHP(ctx, battle.ID, villainRow.MaxHP)
			if err != nil {
				log.Printf("[VillainService] Failed to reset HP: %v", err)
			}
			battle.CurrentHP = villainRow.MaxHP
		}
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

// wasDefeatedToday проверяет, был ли злодей побеждён сегодня
func (s *VillainService) wasDefeatedToday(ctx context.Context, childProfileID string) (bool, error) {
	defeatedAt, err := s.store.Villains.GetLastDefeatedAt(ctx, childProfileID)
	if err != nil {
		return false, err
	}

	if defeatedAt == nil {
		return false, nil
	}

	// Проверяем что победа была сегодня
	today := time.Now().UTC().Truncate(24 * time.Hour)
	battleDate := defeatedAt.Truncate(24 * time.Hour)

	return battleDate.Equal(today), nil
}

// ensureDailyVillain создаёт злодея на сегодня если ещё нет
func (s *VillainService) ensureDailyVillain(ctx context.Context, childProfileID string) error {
	// Проверяем, нет ли уже активного злодея
	battle, _, err := s.store.Villains.GetActiveVillainBattle(ctx, childProfileID)
	if err != nil {
		return fmt.Errorf("failed to check active battle: %w", err)
	}

	if battle != nil && battle.Status == "active" {
		// Уже есть активный злодей
		return nil
	}

	// Определяем день недели (1=понедельник, 7=воскресенье)
	today := time.Now().UTC()
	dayOfWeek := int(today.Weekday()) // Sunday=0, Monday=1, ..., Saturday=6
	if dayOfWeek == 0 {
		dayOfWeek = 7 // Sunday = 7
	}

	// Получаем злодея для этого дня недели (unlock_order = day_of_week)
	villain, err := s.store.Villains.GetVillainByOrder(ctx, dayOfWeek)
	if err != nil {
		return fmt.Errorf("failed to get villain for day %d: %w", dayOfWeek, err)
	}

	if villain == nil {
		// Если нет злодея для этого дня, пробуем найти любого доступного
		log.Printf("[VillainService] No villain for day %d, trying to find any available", dayOfWeek)
		for i := 1; i <= 7; i++ {
			v, err := s.store.Villains.GetVillainByOrder(ctx, i)
			if err == nil && v != nil {
				villain = v
				break
			}
		}
	}

	if villain == nil {
		return fmt.Errorf("no villains available for any day")
	}

	// Пропускаем боссов если нужно
	if villain.IsBoss {
		lastBossDate, err := s.store.Villains.GetLastBossDefeatedAt(ctx, childProfileID)
		if err == nil && lastBossDate != nil {
			daysSinceLastBoss := time.Since(*lastBossDate).Hours() / 24
			if daysSinceLastBoss < 7 {
				// Берём любого не-босса
				for i := 1; i <= 7; i++ {
					v, err := s.store.Villains.GetVillainByOrder(ctx, i)
					if err == nil && v != nil && !v.IsBoss {
						villain = v
						break
					}
				}
			}
		}
	}

	// Создаём битву
	err = s.store.Villains.CreateBattle(ctx, childProfileID, villain.ID, villain.MaxHP)
	if err != nil {
		return fmt.Errorf("failed to create battle: %w", err)
	}

	log.Printf("[VillainService] Created daily villain battle: %s (day %d, weekday %s) for child %s",
		villain.ID, dayOfWeek, today.Weekday().String(), childProfileID)
	return nil
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
		// Нет активной битвы - проверяем, был ли побеждён сегодня
		defeatedToday, checkErr := s.wasDefeatedToday(ctx, childProfileID)
		if checkErr != nil {
			log.Printf("[VillainService] Failed to check defeat status: %v", checkErr)
		}

		if defeatedToday {
			// Уже побеждён сегодня - не наносим урон
			log.Printf("[VillainService] Villain already defeated today for %s, skipping damage", childProfileID)
			return false, 0, nil
		}

		// Создаём нового злодея
		log.Printf("[VillainService] No active battle for %s, creating new villain", childProfileID)
		err = s.ensureDailyVillain(ctx, childProfileID)
		if err != nil {
			return false, 0, fmt.Errorf("failed to create villain: %w", err)
		}
		// Повторно получаем битву
		battle, villainRow, err = s.store.Villains.GetActiveVillainBattle(ctx, childProfileID)
		if err != nil || battle == nil {
			return false, 0, fmt.Errorf("failed to get battle after creation: %w", err)
		}
	}

	// Проверяем что битва активна
	if battle.Status != "active" {
		return false, 0, fmt.Errorf("battle is not active")
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

		// Не создаём следующего злодея до завтра
		log.Printf("[VillainService] Villain defeated for today, next villain will spawn tomorrow for %s", childProfileID)

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

// GetVillainByID получает злодея по ID из БД
func (s *VillainService) GetVillainByID(ctx context.Context, childProfileID, villainID string) (*Villain, error) {
	// Получаем злодея из справочника
	villainRow, err := s.store.Villains.GetVillainByID(ctx, villainID)
	if err != nil {
		return nil, fmt.Errorf("failed to get villain: %w", err)
	}

	if villainRow == nil {
		return nil, nil
	}

	// Получаем статус битвы
	battle, _, _ := s.store.Villains.GetVillainBattleByVillainID(ctx, childProfileID, villainID)

	hp := villainRow.MaxHP
	isActive := false
	isDefeated := false
	var unlockedAt *time.Time
	var defeatedAt *time.Time

	if battle != nil {
		hp = battle.CurrentHP
		isActive = battle.Status == "active"
		isDefeated = battle.Status == "defeated"
		unlockedAt = &battle.StartedAt
		if battle.DefeatedAt.Valid {
			defeatedAt = &battle.DefeatedAt.Time
		}
	}

	return &Villain{
		ID:          villainRow.ID,
		Name:        villainRow.Name,
		Description: villainRow.Description,
		ImageURL:    villainRow.ImageURL,
		HP:          hp,
		MaxHP:       villainRow.MaxHP,
		Level:       villainRow.Level,
		IsActive:    isActive,
		IsDefeated:  isDefeated,
		UnlockedAt:  unlockedAt,
		DefeatedAt:  defeatedAt,
	}, nil
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

// DealDamageToVillain наносит урон злодею и проверяет победу
func (s *VillainService) DealDamage(ctx context.Context, childProfileID, villainID string, attemptID string, damage int) (*DamageResult, error) {
	// Получаем битву
	battle, villainRow, err := s.store.Villains.GetActiveVillainBattle(ctx, childProfileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active battle: %w", err)
	}

	if battle == nil || villainRow == nil {
		return nil, fmt.Errorf("no active battle found")
	}

	// Проверяем что битва активна
	if battle.Status != "active" {
		return nil, fmt.Errorf("battle is not active")
	}

	// Обновляем HP
	newHP := battle.CurrentHP - damage
	if newHP < 0 {
		newHP = 0
	}

	// Записываем событие урона
	attemptUUID, err := uuid.Parse(attemptID)
	if err != nil {
		return nil, fmt.Errorf("invalid attempt ID: %w", err)
	}

	err = s.store.Villains.RecordDamageEvent(ctx, battle.ID, attemptUUID, damage, "manual")
	if err != nil {
		log.Printf("[VillainService] Failed to record damage event: %v", err)
	}

	// Обновляем прогресс
	err = s.store.Villains.UpdateBattleProgress(ctx, battle.ID, newHP, damage)
	if err != nil {
		return nil, fmt.Errorf("failed to update battle progress: %w", err)
	}

	// Проверяем победу
	defeated := newHP <= 0

	if defeated {
		log.Printf("[VillainService] Villain %s defeated! Awarding coins: %d", villainRow.ID, villainRow.RewardCoins)

		// Помечаем битву как побеждённую
		err = s.store.Villains.MarkBattleDefeated(ctx, battle.ID)
		if err != nil {
			log.Printf("[VillainService] Failed to mark battle as defeated: %v", err)
		}

		// Создаём следующего монстра
		err = s.createNextVillain(ctx, childProfileID, villainRow.UnlockOrder)
		if err != nil {
			log.Printf("[VillainService] Failed to create next villain: %v", err)
		}

		// Проверяем достижения
		if s.achievementService != nil {
			err = s.achievementService.CheckVillainAchievements(ctx, childProfileID)
			if err != nil {
				log.Printf("[VillainService] Failed to check villain achievements: %v", err)
			}
		}
	}

	// Формируем результат
	result := &DamageResult{
		DamageDealt:  damage,
		VillainHP:    newHP,
		VillainMaxHP: villainRow.MaxHP,
		IsDefeated:   defeated,
		Rewards:      []VictoryReward{},
	}

	if defeated {
		result.Rewards = append(result.Rewards, VictoryReward{
			Type:   "coins",
			ID:     fmt.Sprintf("coins_%d", villainRow.RewardCoins),
			Name:   fmt.Sprintf("%d монет", villainRow.RewardCoins),
			Amount: villainRow.RewardCoins,
		})
	}

	return result, nil
}

// GetVillainVictory получает информацию о победе
func (s *VillainService) GetVillainVictory(ctx context.Context, childProfileID, villainID string) (*VictoryData, error) {
	// Получаем битву
	battle, villainRow, err := s.store.Villains.GetVillainBattleByVillainID(ctx, childProfileID, villainID)
	if err != nil {
		return nil, fmt.Errorf("failed to get villain battle: %w", err)
	}

	if battle == nil || villainRow == nil {
		return nil, fmt.Errorf("villain battle not found")
	}

	// Формируем награды
	rewards := []VictoryReward{
		{
			Type:   "coins",
			ID:     fmt.Sprintf("coins_%d", villainRow.RewardCoins),
			Name:   fmt.Sprintf("%d монет", villainRow.RewardCoins),
			Amount: villainRow.RewardCoins,
		},
	}

	// Получаем следующего злодея
	var nextVillain *Villain
	nextVillainRow, err := s.store.Villains.GetVillainByOrder(ctx, villainRow.UnlockOrder+1)
	if err != nil {
		log.Printf("[VillainService] Failed to get next villain: %v", err)
	}

	if nextVillainRow != nil {
		nextVillain = &Villain{
			ID:          nextVillainRow.ID,
			Name:        nextVillainRow.Name,
			Description: nextVillainRow.Description,
			ImageURL:    nextVillainRow.ImageURL,
			HP:          nextVillainRow.MaxHP,
			MaxHP:       nextVillainRow.MaxHP,
			Level:       nextVillainRow.Level,
		}
	}

	victory := &VictoryData{
		VillainID:      villainRow.ID,
		VillainName:    villainRow.Name,
		DefeatedAt:     battle.DefeatedAt.Time,
		TotalDamage:    villainRow.MaxHP,
		TasksCompleted: battle.CorrectTasksCount,
		Rewards:        rewards,
		NextVillain:    nextVillain,
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
