package service

import (
	"context"
	"log"

	"child-bot/api/internal/store"
)

// HomeService бизнес-логика для главного экрана
type HomeService struct {
	store          *store.Store
	attemptService *AttemptService
	profileService *ProfileService
	villainService *VillainService
}

// NewHomeService создает новый HomeService
func NewHomeService(
	store *store.Store,
	attemptService *AttemptService,
	profileService *ProfileService,
	villainService *VillainService,
) *HomeService {
	return &HomeService{
		store:          store,
		attemptService: attemptService,
		profileService: profileService,
		villainService: villainService,
	}
}

// GetStore возвращает store для прямого доступа
func (s *HomeService) GetStore() *store.Store {
	return s.store
}

// HomeData данные для главного экрана
type HomeData struct {
	Profile           ProfileSummary
	Mascot            MascotData
	Villain           *VillainSummary
	UnfinishedAttempt *AttemptData
	RecentAttempts    []RecentAttempt
	Achievements      AchievementsSummary
}

// AchievementsSummary статистика достижений
type AchievementsSummary struct {
	UnlockedCount int
	TotalCount    int
}

// ProfileSummary краткие данные профиля для home
type ProfileSummary struct {
	ID                      string
	DisplayName             string
	Level                   int
	LevelProgress           int // 0-100
	CoinsBalance            int
	TasksSolvedCorrectCount int
}

// MascotData данные маскота
type MascotData struct {
	ID       string
	State    string // idle, happy, thinking, celebrating
	ImageURL string
	Message  string
}

// VillainSummary краткие данные злодея
type VillainSummary struct {
	ID         string
	Name       string
	ImageURL   string
	HP         int
	MaxHP      int
	IsActive   bool
	IsDefeated bool
}

// RecentAttempt последняя попытка
type RecentAttempt struct {
	ID            string
	Mode          string // help or check
	Status        string // success, error, in_progress
	CreatedAt     string
	Thumbnail     string
	ResultSummary string
}

// GetHomeData получает все данные для главного экрана
func (s *HomeService) GetHomeData(ctx context.Context, childProfileID string) (*HomeData, error) {
	// ОБНОВЛЯЕМ СЕРИЮ ДНЕЙ (streak) и активность при заходе на главную
	err := s.profileService.UpdateStreakAndActivity(ctx, childProfileID)
	if err != nil {
		// Логируем ошибку, но не блокируем загрузку home
		// (streak - не критичная функциональность)
		// log уже внутри UpdateStreakAndActivity
	}

	// Загружаем профиль из БД
	profile, err := s.profileService.GetProfile(ctx, childProfileID)
	if err != nil {
		return nil, err
	}

	// Подсчитываем статистику достижений
	unlockedCount, totalCount, err := s.store.GetAchievementStats(ctx, childProfileID)
	if err != nil {
		// Игнорируем ошибку, используем дефолтные значения
		unlockedCount = 0
		totalCount = 0
	}

	// Загружаем баланс монет из БД
	var coinsBalance int
	coinsQuery := `SELECT COALESCE(coins_balance, 0) FROM child_profiles WHERE id = $1`
	err = s.store.DB.QueryRowContext(ctx, coinsQuery, childProfileID).Scan(&coinsBalance)
	if err != nil {
		coinsBalance = 0 // Игнорируем ошибку, используем 0
	}

	// Загружаем XP и уровень из БД
	xpTotal, level, err := s.store.GetXPAndLevel(ctx, childProfileID)
	if err != nil {
		log.Printf("[HomeService] Failed to get XP and level: %v", err)
		xpTotal = 0
		level = 1
	}

	// Рассчитываем прогресс уровня
	xpForCurrentLevel := store.XPForLevel(level - 1)
	xpForNextLevel := store.XPForLevel(level)
	xpInCurrentLevel := xpTotal - xpForCurrentLevel
	xpNeeded := xpForNextLevel - xpForCurrentLevel

	levelProgress := 0
	if xpNeeded > 0 {
		levelProgress = (xpInCurrentLevel * 100) / xpNeeded
	}

	// Агрегируем данные из разных источников
	data := &HomeData{
		Profile: ProfileSummary{
			ID:                      childProfileID,
			DisplayName:             profile.DisplayName,
			Level:                   level,
			LevelProgress:           levelProgress,
			CoinsBalance:            coinsBalance,
			TasksSolvedCorrectCount: unlockedCount, // Используем количество разблокированных достижений как количество решённых задач
		},
		Mascot: MascotData{
			ID:       "owl_1",
			State:    "idle",
			ImageURL: "/assets/mascot/owl_idle.png",
			Message:  "Привет! Готов решать задачи?",
		},
		RecentAttempts: []RecentAttempt{},
		Achievements: AchievementsSummary{
			UnlockedCount: unlockedCount,
			TotalCount:    totalCount,
		},
	}

	// Загружаем активного злодея из БД
	activeVillain, err := s.villainService.GetActiveVillain(ctx, childProfileID)
	if err != nil {
		log.Printf("[HomeService] Failed to get active villain: %v", err)
		activeVillain = nil
	}

	if activeVillain != nil {
		data.Villain = &VillainSummary{
			ID:         activeVillain.ID,
			Name:       activeVillain.Name,
			ImageURL:   activeVillain.ImageURL,
			HP:         activeVillain.HP,
			MaxHP:      activeVillain.MaxHP,
			IsActive:   activeVillain.IsActive,
			IsDefeated: activeVillain.IsDefeated,
		}
		log.Printf("[HomeService] Loaded active villain: %s (HP: %d/%d)",
			activeVillain.ID, activeVillain.HP, activeVillain.MaxHP)
	} else {
		// Нет активного злодея - используем дефолтные значения
		data.Villain = &VillainSummary{
			ID:         "villain_1",
			Name:       "Граф Ошибок",
			ImageURL:   "/assets/villains/count_error.png",
			HP:         100,
			MaxHP:      100,
			IsActive:   true,
			IsDefeated: false,
		}
		log.Printf("[HomeService] No active villain, using default")
	}

	// Получить незавершенную попытку
	unfinished, err := s.attemptService.GetUnfinishedAttempt(ctx, childProfileID)
	if err != nil {
		log.Printf("[HomeService] Failed to get unfinished attempt: %v", err)
	}
	if unfinished != nil {
		log.Printf("[HomeService] Found unfinished attempt: %s (type=%s)", unfinished.ID, unfinished.Type)
		data.UnfinishedAttempt = unfinished
	} else {
		log.Printf("[HomeService] No unfinished attempt found for profile: %s", childProfileID)
	}

	return data, nil
}
