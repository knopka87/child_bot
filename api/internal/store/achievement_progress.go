package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

// UpdateAchievementProgress обновляет прогресс достижения и проверяет разблокировку
// Возвращает: (wasUnlocked bool, error)
func (s *Store) UpdateAchievementProgress(ctx context.Context, childProfileID, achievementID string, newProgress int) (bool, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	// Получаем requirement_value из achievements
	var requirementValue int
	achievementQuery := `SELECT requirement_value FROM achievements WHERE id = $1`
	err = tx.QueryRowContext(ctx, achievementQuery, achievementID).Scan(&requirementValue)
	if err != nil {
		return false, fmt.Errorf("get achievement requirement: %w", err)
	}

	// Проверяем есть ли уже запись в child_achievements
	var existingProgress int
	var isUnlocked bool
	checkQuery := `
		SELECT current_progress, is_unlocked
		FROM child_achievements
		WHERE child_profile_id = $1 AND achievement_id = $2
	`
	err = tx.QueryRowContext(ctx, checkQuery, childProfileID, achievementID).Scan(&existingProgress, &isUnlocked)

	if err == sql.ErrNoRows {
		// Создаём новую запись
		isUnlocked = newProgress >= requirementValue
		insertQuery := `
			INSERT INTO child_achievements (child_profile_id, achievement_id, current_progress, is_unlocked, unlocked_at)
			VALUES ($1, $2, $3, $4, CASE WHEN $4 THEN NOW() ELSE NULL END)
		`
		_, err = tx.ExecContext(ctx, insertQuery, childProfileID, achievementID, newProgress, isUnlocked)
		if err != nil {
			return false, fmt.Errorf("insert achievement progress: %w", err)
		}

		if err := tx.Commit(); err != nil {
			return false, fmt.Errorf("commit tx: %w", err)
		}

		log.Printf("[Store] Achievement progress created: child=%s, achievement=%s, progress=%d/%d, unlocked=%v",
			childProfileID, achievementID, newProgress, requirementValue, isUnlocked)

		return isUnlocked, nil
	} else if err != nil {
		return false, fmt.Errorf("check existing progress: %w", err)
	}

	// Обновляем существующую запись
	// Если уже разблокировано - ничего не делаем
	if isUnlocked {
		return false, nil
	}

	// Обновляем прогресс
	wasJustUnlocked := newProgress >= requirementValue

	updateQuery := `
		UPDATE child_achievements
		SET current_progress = $1,
		    is_unlocked = $2,
		    unlocked_at = CASE WHEN $2 AND unlocked_at IS NULL THEN NOW() ELSE unlocked_at END,
		    updated_at = NOW()
		WHERE child_profile_id = $3 AND achievement_id = $4
	`
	_, err = tx.ExecContext(ctx, updateQuery, newProgress, wasJustUnlocked, childProfileID, achievementID)
	if err != nil {
		return false, fmt.Errorf("update achievement progress: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("commit tx: %w", err)
	}

	log.Printf("[Store] Achievement progress updated: child=%s, achievement=%s, progress=%d/%d, unlocked=%v",
		childProfileID, achievementID, newProgress, requirementValue, wasJustUnlocked)

	return wasJustUnlocked, nil
}

// CheckAndUpdateAchievementsByType проверяет и обновляет прогресс всех достижений заданного типа
// Возвращает список ID разблокированных достижений
func (s *Store) CheckAndUpdateAchievementsByType(ctx context.Context, childProfileID, requirementType string, currentValue int) ([]string, error) {
	// Получаем все достижения данного типа
	query := `
		SELECT id, requirement_value
		FROM achievements
		WHERE requirement_type = $1
		ORDER BY requirement_value ASC
	`

	rows, err := s.DB.QueryContext(ctx, query, requirementType)
	if err != nil {
		return nil, fmt.Errorf("query achievements: %w", err)
	}
	defer rows.Close()

	var unlockedIDs []string

	for rows.Next() {
		var achievementID string
		var requirementValue int

		if err := rows.Scan(&achievementID, &requirementValue); err != nil {
			log.Printf("[Store] Failed to scan achievement: %v", err)
			continue
		}

		// Обновляем прогресс
		wasUnlocked, err := s.UpdateAchievementProgress(ctx, childProfileID, achievementID, currentValue)
		if err != nil {
			log.Printf("[Store] Failed to update achievement %s: %v", achievementID, err)
			continue
		}

		if wasUnlocked {
			unlockedIDs = append(unlockedIDs, achievementID)
			log.Printf("[Store] 🎉 Achievement unlocked! child=%s, achievement=%s, type=%s, value=%d",
				childProfileID, achievementID, requirementType, currentValue)
		}
	}

	if err := rows.Err(); err != nil {
		return unlockedIDs, fmt.Errorf("rows iteration: %w", err)
	}

	return unlockedIDs, nil
}

// GetCurrentStreakDays получает текущий streak пользователя
func (s *Store) GetCurrentStreakDays(ctx context.Context, childProfileID string) (int, error) {
	var streakDays int
	query := `SELECT COALESCE(streak_days, 0) FROM child_profiles WHERE id = $1`
	err := s.DB.QueryRowContext(ctx, query, childProfileID).Scan(&streakDays)
	if err != nil {
		return 0, fmt.Errorf("get streak days: %w", err)
	}
	return streakDays, nil
}

// GetVillainsDefeatedCount получает количество побеждённых монстров
func (s *Store) GetVillainsDefeatedCount(ctx context.Context, childProfileID string) (int, error) {
	var count int
	query := `
		SELECT COUNT(DISTINCT villain_id)
		FROM villain_battles
		WHERE child_profile_id = $1 AND status = 'defeated'
	`
	err := s.DB.QueryRowContext(ctx, query, childProfileID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("get villains defeated count: %w", err)
	}
	return count, nil
}

// GetTasksCorrectCount получает количество правильно решённых задач
func (s *Store) GetTasksCorrectCount(ctx context.Context, childProfileID string) (int, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM attempts
		WHERE child_profile_id = $1
		  AND attempt_type = 'check'
		  AND status = 'completed'
		  AND check_result->>'decision' = 'correct'
	`
	err := s.DB.QueryRowContext(ctx, query, childProfileID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("get tasks correct count: %w", err)
	}
	return count, nil
}

// GetTasksNoHintsCount получает количество задач решённых без подсказок
func (s *Store) GetTasksNoHintsCount(ctx context.Context, childProfileID string) (int, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM attempts
		WHERE child_profile_id = $1
		  AND attempt_type = 'check'
		  AND status = 'completed'
		  AND hints_used = 0
		  AND check_result->>'decision' = 'correct'
	`
	err := s.DB.QueryRowContext(ctx, query, childProfileID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("get tasks no hints count: %w", err)
	}
	return count, nil
}

// GetFriendsInvitedCount получает количество приглашённых друзей
func (s *Store) GetFriendsInvitedCount(ctx context.Context, childProfileID string) (int, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM referrals
		WHERE referrer_child_profile_id = $1
		  AND invited_child_profile_id IS NOT NULL
	`
	err := s.DB.QueryRowContext(ctx, query, childProfileID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("get friends invited count: %w", err)
	}
	return count, nil
}

// GetHintsUsedCount получает общее количество использованных подсказок
func (s *Store) GetHintsUsedCount(ctx context.Context, childProfileID string) (int, error) {
	var count int
	query := `SELECT COALESCE(hints_used_total, 0) FROM child_profiles WHERE id = $1`
	err := s.DB.QueryRowContext(ctx, query, childProfileID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("get hints used count: %w", err)
	}
	return count, nil
}

// GetStickersCollectedCount получает количество собранных стикеров
func (s *Store) GetStickersCollectedCount(ctx context.Context, childProfileID string) (int, error) {
	var count int
	query := `
		SELECT COUNT(DISTINCT ca.achievement_id)
		FROM child_achievements ca
		JOIN achievements a ON ca.achievement_id = a.id
		WHERE ca.child_profile_id = $1
		  AND ca.is_unlocked = TRUE
		  AND a.reward_type = 'sticker'
	`
	err := s.DB.QueryRowContext(ctx, query, childProfileID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("get stickers collected count: %w", err)
	}
	return count, nil
}

// GetErrorsFoundCount получает количество найденных ошибок во всех проверках
func (s *Store) GetErrorsFoundCount(ctx context.Context, childProfileID string) (int, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM attempts
		WHERE child_profile_id = $1
		  AND attempt_type = 'check'
		  AND status = 'completed'
		  AND check_result->>'decision' != 'correct'
	`
	err := s.DB.QueryRowContext(ctx, query, childProfileID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("get errors found count: %w", err)
	}
	return count, nil
}
