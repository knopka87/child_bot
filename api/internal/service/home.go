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

	// Агрегируем данные из разных источников
	data := &HomeData{
		Profile: ProfileSummary{
			ID:                      childProfileID,
			DisplayName:             profile.DisplayName,
			Level:                   1, // TODO: Phase 5 - загружать level из БД
			LevelProgress:           0, // TODO: Phase 5 - вычислять прогресс уровня
			CoinsBalance:            coinsBalance,
			TasksSolvedCorrectCount: unlockedCount, // Используем количество разблокированных достижений как количество решённых задач
		},
		Mascot: MascotData{
			ID:       "owl_1",
			State:    "idle",
			ImageURL: "/assets/mascot/owl_idle.png",
			Message:  "Привет! Готов решать задачи?",
		},
		// ВРЕМЕННО: Мок-данные для villain
		Villain: &VillainSummary{
			ID:         "villain_1",
			Name:       "Кракозябра",
			ImageURL:   "/images/villain.png",
			HP:         2,
			MaxHP:      3,
			IsActive:   true,
			IsDefeated: false,
		},
		RecentAttempts: []RecentAttempt{},
		Achievements: AchievementsSummary{
			UnlockedCount: unlockedCount,
			TotalCount:    totalCount,
		},
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

	// Получить активного злодея
	// TODO: villain, _ := s.villainService.GetActiveVillain(ctx, childProfileID)

	return data, nil
}
