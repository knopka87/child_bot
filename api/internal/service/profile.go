package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"child-bot/api/internal/domain"
	"child-bot/api/internal/store"

	"github.com/google/uuid"
)

// ProfileService бизнес-логика для профиля
type ProfileService struct {
	store              *store.Store
	achievementService *AchievementService
}

// NewProfileService создает новый ProfileService
func NewProfileService(store *store.Store) *ProfileService {
	return &ProfileService{store: store}
}

// SetAchievementService устанавливает AchievementService (для избежания циклических зависимостей)
func (s *ProfileService) SetAchievementService(achievementService *AchievementService) {
	s.achievementService = achievementService
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
	Status             string // trial, active, expired, cancelled
	PlanID             string
	PlanName           string
	TrialDaysRemaining int
	ExpiresAt          *time.Time
	RenewsAt           *time.Time
	CancelledAt        *time.Time
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
	// Парсим UUID
	profileUUID, err := uuid.Parse(childProfileID)
	if err != nil {
		return nil, fmt.Errorf("invalid child_profile_id: %w", err)
	}

	// Получаем попытки из БД (последние 100)
	attempts, err := s.store.Attempts.GetRecentAttempts(ctx, profileUUID, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to get attempts: %w", err)
	}

	// Преобразуем в формат HistoryAttempt
	history := make([]HistoryAttempt, 0, len(attempts))
	for _, attempt := range attempts {
		// Определяем статус для фронтенда
		var status string
		switch attempt.Status {
		case "completed":
			// Для режима check - проверяем результат
			if attempt.AttemptType == "check" {
				if attempt.IsCorrect.Valid && attempt.IsCorrect.Bool {
					status = "success"
				} else if attempt.HasErrors.Valid && attempt.HasErrors.Bool {
					status = "error"
				} else {
					status = "completed"
				}
			} else {
				// Для режима help - просто completed
				status = "completed"
			}
		case "processing", "created":
			status = "in_progress"
		case "failed":
			status = "error"
		default:
			status = "in_progress"
		}

		// Собираем изображения
		images := []HistoryImage{}
		if attempt.TaskImageURL.Valid {
			images = append(images, HistoryImage{
				ID:           uuid.New().String(),
				Role:         "task",
				URL:          attempt.TaskImageURL.String,
				ThumbnailURL: attempt.TaskImageURL.String,
			})
		}
		if attempt.AnswerImageURL.Valid {
			images = append(images, HistoryImage{
				ID:           uuid.New().String(),
				Role:         "answer",
				URL:          attempt.AnswerImageURL.String,
				ThumbnailURL: attempt.AnswerImageURL.String,
			})
		}

		// Определяем scenario_type
		scenarioType := ""
		if attempt.TaskImageURL.Valid && attempt.AnswerImageURL.Valid {
			scenarioType = "two_photo"
		} else if attempt.TaskImageURL.Valid {
			scenarioType = "single_photo"
		}

		historyAttempt := HistoryAttempt{
			ID:           attempt.ID.String(),
			Mode:         attempt.AttemptType,
			Status:       status,
			ScenarioType: scenarioType,
			CreatedAt:    attempt.CreatedAt,
			HintsUsed:    attempt.HintsUsed,
			Images:       images,
		}

		if attempt.CompletedAt.Valid {
			historyAttempt.CompletedAt = &attempt.CompletedAt.Time
		}

		history = append(history, historyAttempt)
	}

	return history, nil
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
	// Проверяем, не активирован ли уже, и получаем referrer_id
	checkQuery := `
		SELECT is_active, referrer_id
		FROM referrals
		WHERE referred_id = $1
	`
	var isActive bool
	var referrerID string
	err := s.store.DB.QueryRowContext(ctx, checkQuery, childProfileID).Scan(&isActive, &referrerID)
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

	// Проверяем достижения за приглашённых друзей у реферера
	if s.achievementService != nil && referrerID != "" {
		err = s.achievementService.CheckFriendsInvitedAchievements(ctx, referrerID)
		if err != nil {
			log.Printf("[ActivateReferral] Failed to check friends invited achievements for referrer %s: %v", referrerID, err)
			// Не блокируем, продолжаем
		}
	}

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

// IncrementHintsUsed увеличивает счётчик использованных подсказок в профиле
func (s *ProfileService) IncrementHintsUsed(ctx context.Context, childProfileID string) error {
	query := `
		UPDATE child_profiles
		SET hints_used_total = hints_used_total + 1,
		    updated_at = NOW()
		WHERE id = $1
	`

	result, err := s.store.DB.ExecContext(ctx, query, childProfileID)
	if err != nil {
		log.Printf("[IncrementHintsUsed] Error incrementing hints_used_total for %s: %v", childProfileID, err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	log.Printf("[IncrementHintsUsed] Successfully incremented hints_used_total for child %s", childProfileID)
	return nil
}

// UpdateStreakAndActivity обновляет серию дней (streak) и last_activity_at
// Вызывается ОДИН РАЗ В ДЕНЬ при первом действии пользователя
func (s *ProfileService) UpdateStreakAndActivity(ctx context.Context, childProfileID string) error {
	// Получаем профиль
	query := `
		SELECT last_activity_at, streak_days
		FROM child_profiles
		WHERE id = $1
	`

	var lastActivityAt *time.Time
	var currentStreak int

	err := s.store.DB.QueryRowContext(ctx, query, childProfileID).Scan(&lastActivityAt, &currentStreak)
	if err != nil {
		log.Printf("[UpdateStreakAndActivity] Error getting profile %s: %v", childProfileID, err)
		return err
	}

	now := time.Now().UTC()
	newStreak := currentStreak

	// Если last_activity_at пустой (первый заход), устанавливаем streak = 1
	if lastActivityAt == nil {
		newStreak = 1
		log.Printf("[UpdateStreakAndActivity] First login for %s, setting streak = 1", childProfileID)
	} else {
		// Вычисляем разницу в днях (только дата, без времени)
		lastDate := lastActivityAt.Truncate(24 * time.Hour)
		currentDate := now.Truncate(24 * time.Hour)
		daysDiff := int(currentDate.Sub(lastDate).Hours() / 24)

		log.Printf("[UpdateStreakAndActivity] Profile %s: last_activity=%s, now=%s, days_diff=%d",
			childProfileID, lastActivityAt.Format("2006-01-02 15:04"), now.Format("2006-01-02 15:04"), daysDiff)

		if daysDiff == 0 {
			// Тот же день - ничего не делаем со streak, только обновляем time
			log.Printf("[UpdateStreakAndActivity] Same day, keeping streak=%d", currentStreak)
			// Обновляем только last_activity_at
			updateQuery := `
				UPDATE child_profiles
				SET last_activity_at = $1,
				    updated_at = NOW()
				WHERE id = $2
			`
			_, err = s.store.DB.ExecContext(ctx, updateQuery, now, childProfileID)
			if err != nil {
				return err
			}
			return nil
		} else if daysDiff == 1 {
			// Следующий день - увеличиваем streak
			newStreak = currentStreak + 1
			log.Printf("[UpdateStreakAndActivity] Next day, incrementing streak: %d -> %d", currentStreak, newStreak)
		} else {
			// Пропущено больше 1 дня - сбрасываем streak
			newStreak = 1
			log.Printf("[UpdateStreakAndActivity] Missed days (%d), resetting streak: %d -> 1", daysDiff, currentStreak)
		}
	}

	// Обновляем streak_days и last_activity_at
	updateQuery := `
		UPDATE child_profiles
		SET streak_days = $1,
		    last_activity_at = $2,
		    updated_at = NOW()
		WHERE id = $3
	`

	result, err := s.store.DB.ExecContext(ctx, updateQuery, newStreak, now, childProfileID)
	if err != nil {
		log.Printf("[UpdateStreakAndActivity] Error updating profile %s: %v", childProfileID, err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	log.Printf("[UpdateStreakAndActivity] Successfully updated profile %s: streak=%d, last_activity=%s",
		childProfileID, newStreak, now.Format("2006-01-02 15:04:05"))

	// Проверяем достижения за streak (только если изменился)
	if s.achievementService != nil && newStreak != currentStreak {
		err = s.achievementService.CheckStreakAchievements(ctx, childProfileID)
		if err != nil {
			log.Printf("[UpdateStreakAndActivity] Failed to check streak achievements: %v", err)
			// Не блокируем, продолжаем
		}
	}

	return nil
}
