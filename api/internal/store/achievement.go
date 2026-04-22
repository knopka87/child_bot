package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// AchievementDB представляет достижение в базе данных
type AchievementDB struct {
	ID               string
	Type             string
	Title            string
	Description      string
	Icon             string
	RequirementType  string
	RequirementValue int
	RewardType       string
	RewardID         sql.NullString
	RewardName       sql.NullString
	RewardAmount     sql.NullInt32
	ShelfOrder       int
	PositionInShelf  int
	CreatedAt        time.Time
}

// ChildAchievementDB представляет прогресс достижения для ребенка
type ChildAchievementDB struct {
	ID              int64
	ChildProfileID  string
	AchievementID   string
	CurrentProgress int
	IsUnlocked      bool
	IsClaimed       bool
	UnlockedAt      sql.NullTime
	ClaimedAt       sql.NullTime
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// CombinedAchievement объединяет данные о достижении и прогрессе пользователя
type CombinedAchievement struct {
	// Achievement data
	ID               string
	Type             string
	Title            string
	Description      string
	Icon             string
	RequirementType  string
	RequirementValue int
	RewardType       string
	RewardID         sql.NullString
	RewardName       sql.NullString
	RewardAmount     sql.NullInt32
	Priority         int
	// Progress data
	CurrentProgress int
	IsUnlocked      bool
	UnlockedAt      sql.NullTime
}

// NextLevelInfo информация о следующем уровне серийной награды
type NextLevelInfo struct {
	Description      string
	RequirementValue int
}

// ListAchievementsWithProgress получает все достижения с прогрессом для конкретного пользователя
// Для серийных достижений (с одинаковым reward_name) показывает только максимальное разблокированное
func (s *Store) ListAchievementsWithProgress(ctx context.Context, childProfileID string) ([]CombinedAchievement, error) {
	query := `
		WITH achievements_with_progress AS (
			SELECT
				a.id,
				a.type,
				a.title,
				a.description,
				a.icon,
				a.requirement_type,
				a.requirement_value,
				a.reward_type,
				a.reward_id,
				a.reward_name,
				a.reward_amount,
				a.priority,
				COALESCE(ca.current_progress, 0) as current_progress,
				COALESCE(ca.is_unlocked, FALSE) as is_unlocked,
				ca.unlocked_at,
				-- Ранжируем внутри каждой серии (reward_name)
				-- Для разблокированных берем максимальное, для неразблокированных - минимальное
				ROW_NUMBER() OVER (
					PARTITION BY COALESCE(a.reward_name, a.id)
					ORDER BY
						CASE WHEN COALESCE(ca.is_unlocked, FALSE) = TRUE
							THEN a.requirement_value
							ELSE -a.requirement_value  -- для неразблокированных инвертируем сортировку
						END DESC
				) as rn
			FROM achievements a
			LEFT JOIN child_achievements ca
				ON a.id = ca.achievement_id
				AND ca.child_profile_id = $1
		)
		SELECT
			id, type, title, description, icon,
			requirement_type, requirement_value,
			reward_type, reward_id, reward_name, reward_amount,
			priority, current_progress, is_unlocked, unlocked_at
		FROM achievements_with_progress
		WHERE rn = 1  -- берем только первое из каждой серии
		ORDER BY
			is_unlocked DESC,                    -- активные первыми
			unlocked_at DESC NULLS LAST,          -- самые новые активные первыми
			priority ASC                          -- неактивные по приоритету
	`

	rows, err := s.DB.QueryContext(ctx, query, childProfileID)
	if err != nil {
		return nil, fmt.Errorf("query achievements: %w", err)
	}
	defer rows.Close()

	var achievements []CombinedAchievement
	for rows.Next() {
		var ach CombinedAchievement
		err := rows.Scan(
			&ach.ID,
			&ach.Type,
			&ach.Title,
			&ach.Description,
			&ach.Icon,
			&ach.RequirementType,
			&ach.RequirementValue,
			&ach.RewardType,
			&ach.RewardID,
			&ach.RewardName,
			&ach.RewardAmount,
			&ach.Priority,
			&ach.CurrentProgress,
			&ach.IsUnlocked,
			&ach.UnlockedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan achievement: %w", err)
		}
		achievements = append(achievements, ach)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate achievements: %w", err)
	}

	return achievements, nil
}

// GetNextLevelsForSeries получает информацию о следующих уровнях для серийных достижений
// Возвращает map[achievementID]NextLevelInfo
func (s *Store) GetNextLevelsForSeries(ctx context.Context, childProfileID string) (map[string]NextLevelInfo, error) {
	query := `
		WITH current_achievements AS (
			-- Получаем текущие достижения пользователя (те что показываются в списке)
			SELECT
				a.id,
				a.reward_type,
				a.reward_name,
				a.requirement_type,
				a.requirement_value,
				COALESCE(ca.current_progress, 0) as current_progress
			FROM achievements a
			LEFT JOIN child_achievements ca
				ON a.id = ca.achievement_id
				AND ca.child_profile_id = $1
		),
		series_achievements AS (
			-- Получаем все достижения в сериях (где reward_name не NULL)
			SELECT
				a.id,
				a.reward_name,
				a.description,
				a.requirement_type,
				a.requirement_value
			FROM achievements a
			WHERE a.reward_type = 'sticker'
			  AND a.reward_name IS NOT NULL
		)
		SELECT DISTINCT ON (ca.id)
			ca.id as current_id,
			sa.description as next_description,
			sa.requirement_value as next_requirement
		FROM current_achievements ca
		JOIN series_achievements sa
			ON ca.reward_name = sa.reward_name
			AND ca.requirement_type = sa.requirement_type
			AND sa.requirement_value > ca.requirement_value
		WHERE ca.reward_type = 'sticker'
		  AND ca.reward_name IS NOT NULL
		ORDER BY ca.id, sa.requirement_value ASC
	`

	rows, err := s.DB.QueryContext(ctx, query, childProfileID)
	if err != nil {
		return nil, fmt.Errorf("query next levels: %w", err)
	}
	defer rows.Close()

	result := make(map[string]NextLevelInfo)
	for rows.Next() {
		var currentID, nextDescription string
		var nextRequirement int

		err := rows.Scan(&currentID, &nextDescription, &nextRequirement)
		if err != nil {
			return nil, fmt.Errorf("scan next level: %w", err)
		}

		result[currentID] = NextLevelInfo{
			Description:      nextDescription,
			RequirementValue: nextRequirement,
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate next levels: %w", err)
	}

	return result, nil
}

// GetAchievementStats возвращает статистику достижений для пользователя
func (s *Store) GetAchievementStats(ctx context.Context, childProfileID string) (unlockedCount, totalCount int, err error) {
	query := `
		SELECT
			COUNT(CASE WHEN ca.is_unlocked THEN 1 END) as unlocked_count,
			COUNT(*) as total_count
		FROM achievements a
		LEFT JOIN child_achievements ca
			ON a.id = ca.achievement_id
			AND ca.child_profile_id = $1
	`

	err = s.DB.QueryRowContext(ctx, query, childProfileID).Scan(&unlockedCount, &totalCount)
	if err != nil {
		return 0, 0, fmt.Errorf("get achievement stats: %w", err)
	}

	return unlockedCount, totalCount, nil
}

// GetAchievementByID получает конкретное достижение с прогрессом пользователя
func (s *Store) GetAchievementByID(ctx context.Context, childProfileID, achievementID string) (*CombinedAchievement, error) {
	query := `
		SELECT
			a.id,
			a.type,
			a.title,
			a.description,
			a.icon,
			a.requirement_type,
			a.requirement_value,
			a.reward_type,
			a.reward_id,
			a.reward_name,
			a.reward_amount,
			a.priority,
			COALESCE(ca.current_progress, 0) as current_progress,
			COALESCE(ca.is_unlocked, FALSE) as is_unlocked,
			ca.unlocked_at
		FROM achievements a
		LEFT JOIN child_achievements ca
			ON a.id = ca.achievement_id
			AND ca.child_profile_id = $1
		WHERE a.id = $2
	`

	var ach CombinedAchievement
	err := s.DB.QueryRowContext(ctx, query, childProfileID, achievementID).Scan(
		&ach.ID,
		&ach.Type,
		&ach.Title,
		&ach.Description,
		&ach.Icon,
		&ach.RequirementType,
		&ach.RequirementValue,
		&ach.RewardType,
		&ach.RewardID,
		&ach.RewardName,
		&ach.RewardAmount,
		&ach.Priority,
		&ach.CurrentProgress,
		&ach.IsUnlocked,
		&ach.UnlockedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get achievement by id: %w", err)
	}

	return &ach, nil
}

// HasNewAchievements проверяет есть ли новые (непросмотренные) достижения
func (s *Store) HasNewAchievements(ctx context.Context, childProfileID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM child_achievements ca
			JOIN child_profiles cp ON ca.child_profile_id = cp.id
			WHERE ca.child_profile_id = $1
				AND ca.is_unlocked = TRUE
				AND ca.unlocked_at > COALESCE(cp.achievements_last_viewed_at, '1970-01-01'::timestamptz)
		)
	`

	var hasNew bool
	err := s.DB.QueryRowContext(ctx, query, childProfileID).Scan(&hasNew)
	if err != nil {
		return false, fmt.Errorf("check new achievements: %w", err)
	}

	return hasNew, nil
}

// MarkAchievementsViewed обновляет дату последнего просмотра достижений
func (s *Store) MarkAchievementsViewed(ctx context.Context, childProfileID string) error {
	query := `
		UPDATE child_profiles
		SET achievements_last_viewed_at = NOW(),
		    updated_at = NOW()
		WHERE id = $1
	`

	_, err := s.DB.ExecContext(ctx, query, childProfileID)
	if err != nil {
		return fmt.Errorf("mark achievements viewed: %w", err)
	}

	return nil
}
