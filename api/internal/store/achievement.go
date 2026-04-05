package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// AchievementDB представляет достижение в базе данных
type AchievementDB struct {
	ID              string
	Type            string
	Title           string
	Description     string
	Icon            string
	RequirementType string
	RequirementValue int
	RewardType      string
	RewardID        sql.NullString
	RewardName      sql.NullString
	RewardAmount    sql.NullInt32
	ShelfOrder      int
	PositionInShelf int
	CreatedAt       time.Time
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
	ID              string
	Type            string
	Title           string
	Description     string
	Icon            string
	RequirementType string
	RequirementValue int
	RewardType      string
	RewardID        sql.NullString
	RewardName      sql.NullString
	RewardAmount    sql.NullInt32
	ShelfOrder      int
	PositionInShelf int
	// Progress data
	CurrentProgress int
	IsUnlocked      bool
	UnlockedAt      sql.NullTime
}

// ListAchievementsWithProgress получает все достижения с прогрессом для конкретного пользователя
func (s *Store) ListAchievementsWithProgress(ctx context.Context, childProfileID string) ([]CombinedAchievement, error) {
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
			a.shelf_order,
			a.position_in_shelf,
			COALESCE(ca.current_progress, 0) as current_progress,
			COALESCE(ca.is_unlocked, FALSE) as is_unlocked,
			ca.unlocked_at
		FROM achievements a
		LEFT JOIN child_achievements ca
			ON a.id = ca.achievement_id
			AND ca.child_profile_id = $1
		ORDER BY a.shelf_order, a.position_in_shelf
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
			&ach.ShelfOrder,
			&ach.PositionInShelf,
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
			a.shelf_order,
			a.position_in_shelf,
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
		&ach.ShelfOrder,
		&ach.PositionInShelf,
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
