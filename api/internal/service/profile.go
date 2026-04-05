package service

import (
	"context"
	"database/sql"
	"log"
	"time"

	"child-bot/api/internal/domain"
	"child-bot/api/internal/store"
)

// ProfileService бизнес-логика для профиля
type ProfileService struct {
	store *store.Store
}

// NewProfileService создает новый ProfileService
func NewProfileService(store *store.Store) *ProfileService {
	return &ProfileService{store: store}
}

// ProfileData полные данные профиля
type ProfileData struct {
	ID           string
	DisplayName  string
	AvatarID     string
	AvatarURL    string
	Grade        int
	Email        string
	Subscription SubscriptionData
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// SubscriptionData данные подписки
type SubscriptionData struct {
	Status              string // trial, active, expired, cancelled
	PlanID              string
	PlanName            string
	TrialDaysRemaining  int
	ExpiresAt           *time.Time
	RenewsAt            *time.Time
	CancelledAt         *time.Time
}

// CreateChildProfile создает профиль ребенка или обновляет существующий (UPSERT)
// Если профиль с таким platform_id + platform_user_id уже существует, обновляет его данные
func (s *ProfileService) CreateChildProfile(ctx context.Context, platformUserID, displayName, avatarID, platformID string, grade int) (string, error) {
	query := `
		INSERT INTO child_profiles (display_name, avatar_id, grade, platform_id, platform_user_id)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (platform_id, platform_user_id)
		DO UPDATE SET
			display_name = EXCLUDED.display_name,
			avatar_id = EXCLUDED.avatar_id,
			grade = EXCLUDED.grade,
			updated_at = NOW()
		RETURNING id
	`

	var childProfileID string
	err := s.store.DB.QueryRowContext(ctx, query, displayName, avatarID, grade, platformID, platformUserID).Scan(&childProfileID)
	if err != nil {
		return "", err
	}

	return childProfileID, nil
}

// Мапа avatar_id -> emoji
var avatarEmojiMap = map[string]string{
	"cat":     "🐱",
	"dog":     "🐶",
	"panda":   "🐼",
	"fox":     "🦊",
	"bear":    "🐻",
	"lion":    "🦁",
	"tiger":   "🐯",
	"unicorn": "🦄",
	"robot":   "🤖",
	"alien":   "👽",
}

// GetProfile получает профиль пользователя
func (s *ProfileService) GetProfile(ctx context.Context, childProfileID string) (*ProfileData, error) {
	// Загрузка базовых данных профиля
	query := `
		SELECT id, display_name, avatar_id, grade, created_at, updated_at
		FROM child_profiles
		WHERE id = $1
	`

	var profile ProfileData
	err := s.store.DB.QueryRowContext(ctx, query, childProfileID).Scan(
		&profile.ID,
		&profile.DisplayName,
		&profile.AvatarID,
		&profile.Grade,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Получаем emoji для аватара
	if emoji, ok := avatarEmojiMap[profile.AvatarID]; ok {
		profile.AvatarURL = emoji
	} else {
		// Fallback
		profile.AvatarURL = "🦊"
	}

	// Загрузка данных подписки
	subscriptionQuery := `
		SELECT
			status,
			plan_id,
			trial_ends_at,
			expires_at,
			cancelled_at
		FROM subscriptions
		WHERE child_profile_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	var status, planID string
	var trialEndsAt, expiresAt, cancelledAt *time.Time

	err = s.store.DB.QueryRowContext(ctx, subscriptionQuery, childProfileID).Scan(
		&status,
		&planID,
		&trialEndsAt,
		&expiresAt,
		&cancelledAt,
	)

	// Если подписки нет, создаем дефолтную trial
	if err == sql.ErrNoRows {
		profile.Subscription = SubscriptionData{
			Status:             "trial",
			TrialDaysRemaining: 7,
		}
	} else if err != nil {
		// Другая ошибка БД
		return nil, err
	} else {
		// Вычисляем оставшиеся дни trial
		trialDaysRemaining := 0
		if status == "trial" && trialEndsAt != nil {
			daysLeft := int(time.Until(*trialEndsAt).Hours() / 24)
			if daysLeft > 0 {
				trialDaysRemaining = daysLeft
			}
		}

		profile.Subscription = SubscriptionData{
			Status:             status,
			PlanID:             planID,
			TrialDaysRemaining: trialDaysRemaining,
			ExpiresAt:          expiresAt,
			CancelledAt:        cancelledAt,
		}
	}

	return &profile, nil
}

// UpdateProfile обновляет профиль
func (s *ProfileService) UpdateProfile(ctx context.Context, childProfileID string, updates map[string]interface{}) error {
	// TODO: Phase 5 - обновление в БД

	// Валидация updates
	if displayName, ok := updates["display_name"].(string); ok {
		if len(displayName) > 50 {
			return domain.ErrInvalidInput
		}
	}

	if grade, ok := updates["grade"].(int); ok {
		if grade < 1 || grade > 4 {
			return domain.ErrInvalidInput
		}
	}

	return nil
}

// GetHistory получает историю попыток
func (s *ProfileService) GetHistory(ctx context.Context, childProfileID string, filters map[string]string) ([]HistoryAttempt, error) {
	// TODO: Phase 5 - загрузка из БД с фильтрами

	return []HistoryAttempt{}, nil
}

// GetStats получает статистику профиля
func (s *ProfileService) GetStats(ctx context.Context, childProfileID string) (*ProfileStats, error) {
	// TODO: Phase 5 - агрегация статистики из БД

	stats := &ProfileStats{
		TotalAttempts:      50,
		SuccessfulAttempts: 42,
		ErrorsFixed:        35,
		StreakDays:         7,
		AverageAccuracy:    84.0,
		TotalHintsUsed:     15,
	}

	return stats, nil
}

// HistoryAttempt попытка в истории
type HistoryAttempt struct {
	ID           string
	Mode         string
	Status       string
	ScenarioType string
	CreatedAt    time.Time
	CompletedAt  *time.Time
	Images       []HistoryImage
	Result       *HistoryResult
	HintsUsed    int
}

// HistoryImage изображение в истории
type HistoryImage struct {
	ID           string
	Role         string // task, answer, single
	URL          string
	ThumbnailURL string
}

// HistoryResult результат в истории
type HistoryResult struct {
	Status     string
	ErrorCount int
	Feedback   []ErrorFeedback
	Summary    string
}

// ErrorFeedback ошибка в решении
type ErrorFeedback struct {
	ID            string
	StepNumber    int
	LineReference string
	Description   string
	LocationType  string
}

// ProfileStats статистика профиля
type ProfileStats struct {
	TotalAttempts      int
	SuccessfulAttempts int
	ErrorsFixed        int
	StreakDays         int
	AverageAccuracy    float64
	TotalHintsUsed     int
}

// ProcessReferral обрабатывает реферальный код и создаёт связь между пользователями
func (s *ProfileService) ProcessReferral(ctx context.Context, childProfileID, referralCode string) error {
	// 1. Найти владельца реферального кода
	var referrerID string
	query := `
		SELECT child_profile_id
		FROM referral_codes
		WHERE code = $1
	`
	err := s.store.DB.QueryRowContext(ctx, query, referralCode).Scan(&referrerID)
	if err == sql.ErrNoRows {
		return domain.ErrNotFound
	}
	if err != nil {
		return err
	}

	// Не разрешаем реферить самого себя
	if referrerID == childProfileID {
		return domain.ErrInvalidInput
	}

	// 2. Создать запись в таблице referrals
	insertQuery := `
		INSERT INTO referrals (referrer_id, referred_id, is_active, reward_coins)
		VALUES ($1, $2, false, 50)
		ON CONFLICT (referrer_id, referred_id) DO NOTHING
	`
	_, err = s.store.DB.ExecContext(ctx, insertQuery, referrerID, childProfileID)
	if err != nil {
		return err
	}

	// 3. Увеличить счётчик uses_count у реферального кода
	updateQuery := `
		UPDATE referral_codes
		SET uses_count = uses_count + 1
		WHERE code = $1
	`
	_, err = s.store.DB.ExecContext(ctx, updateQuery, referralCode)
	if err != nil {
		return err
	}

	return nil
}

// ActivateReferral активирует реферала после первого действия
func (s *ProfileService) ActivateReferral(ctx context.Context, childProfileID string) error {
	// Проверяем, не активирован ли уже
	checkQuery := `
		SELECT is_active
		FROM referrals
		WHERE referred_id = $1
	`
	var isActive bool
	err := s.store.DB.QueryRowContext(ctx, checkQuery, childProfileID).Scan(&isActive)
	if err == sql.ErrNoRows {
		// Пользователь не является рефералом - ничего не делаем
		log.Printf("[ActivateReferral] User %s is not a referral", childProfileID)
		return nil
	}
	if err != nil {
		log.Printf("[ActivateReferral] Error checking referral status for %s: %v", childProfileID, err)
		return err
	}

	// Если уже активирован - ничего не делаем
	if isActive {
		log.Printf("[ActivateReferral] User %s is already active", childProfileID)
		return nil
	}

	// Активируем реферала
	log.Printf("[ActivateReferral] Activating referral for %s", childProfileID)
	activateQuery := `
		UPDATE referrals
		SET is_active = true, activated_at = NOW()
		WHERE referred_id = $1 AND is_active = false
	`
	result, err := s.store.DB.ExecContext(ctx, activateQuery, childProfileID)
	if err != nil {
		log.Printf("[ActivateReferral] Error updating referral for %s: %v", childProfileID, err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("[ActivateReferral] Successfully activated referral for %s, rows affected: %d", childProfileID, rowsAffected)
	return nil
}

// AddCoins начисляет монеты пользователю
func (s *ProfileService) AddCoins(ctx context.Context, childProfileID string, amount int) error {
	if amount <= 0 {
		return domain.ErrInvalidInput
	}

	query := `
		UPDATE child_profiles
		SET coins_balance = coins_balance + $1,
		    updated_at = NOW()
		WHERE id = $2
	`

	result, err := s.store.DB.ExecContext(ctx, query, amount, childProfileID)
	if err != nil {
		log.Printf("[AddCoins] Error adding %d coins to %s: %v", amount, childProfileID, err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	log.Printf("[AddCoins] Added %d coins to child %s", amount, childProfileID)
	return nil
}
