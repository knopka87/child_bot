package store

import (
	"context"
	"fmt"
	"log"
)

// XPConfig конфигурация системы XP
type XPConfig struct {
	LevelUpCoinsReward int // Константные монеты за повышение уровня
}

// DefaultXPConfig стандартная конфигурация
var DefaultXPConfig = XPConfig{
	LevelUpCoinsReward: 100, // 100 монет за каждый уровень
}

// AddXP добавляет XP пользователю и проверяет повышение уровня
// Возвращает: (новый уровень, был ли повышен уровень, ошибка)
func (s *Store) AddXP(ctx context.Context, childProfileID string, xpAmount int, config XPConfig) (int, bool, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, false, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	// Получаем текущие XP и уровень
	var currentXP, currentLevel int
	query := `SELECT COALESCE(xp_total, 0), COALESCE(level, 1) FROM child_profiles WHERE id = $1`
	err = tx.QueryRowContext(ctx, query, childProfileID).Scan(&currentXP, &currentLevel)
	if err != nil {
		return 0, false, fmt.Errorf("get current xp and level: %w", err)
	}

	newXP := currentXP + xpAmount
	newLevel := currentLevel
	leveledUp := false

	// Проверяем повышение уровня
	for {
		xpNeeded := XPForLevel(newLevel)
		if newXP < xpNeeded {
			break
		}

		// Повышаем уровень
		newLevel++
		leveledUp = true
		log.Printf("[Store] 🎉 Level up! child=%s, level %d -> %d (XP: %d)",
			childProfileID, currentLevel, newLevel, newXP)
	}

	// Обновляем XP и уровень
	updateQuery := `
		UPDATE child_profiles
		SET xp_total = $1,
		    level = $2,
		    updated_at = NOW()
		WHERE id = $3
	`
	_, err = tx.ExecContext(ctx, updateQuery, newXP, newLevel, childProfileID)
	if err != nil {
		return 0, false, fmt.Errorf("update xp and level: %w", err)
	}

	// Если был повышен уровень, начисляем бонусные монеты
	if leveledUp {
		coinsReward := config.LevelUpCoinsReward
		coinsQuery := `
			UPDATE child_profiles
			SET coins_balance = coins_balance + $1,
			    updated_at = NOW()
			WHERE id = $2
		`
		_, err = tx.ExecContext(ctx, coinsQuery, coinsReward, childProfileID)
		if err != nil {
			return 0, false, fmt.Errorf("add level up coins: %w", err)
		}

		log.Printf("[Store] 🪙 Level up reward: child=%s, level=%d, coins=%d",
			childProfileID, newLevel, coinsReward)
	}

	if err := tx.Commit(); err != nil {
		return 0, false, fmt.Errorf("commit tx: %w", err)
	}

	if leveledUp {
		log.Printf("[Store] ✅ XP updated: child=%s, XP: %d -> %d, level: %d -> %d, coins_reward=%d",
			childProfileID, currentXP, newXP, currentLevel, newLevel, config.LevelUpCoinsReward)
	} else {
		log.Printf("[Store] ✅ XP updated: child=%s, XP: %d -> %d, level: %d",
			childProfileID, currentXP, newXP, currentLevel)
	}

	return newLevel, leveledUp, nil
}

// GetXPAndLevel получает текущие XP и уровень пользователя
func (s *Store) GetXPAndLevel(ctx context.Context, childProfileID string) (int, int, error) {
	var xpTotal, level int
	query := `SELECT COALESCE(xp_total, 0), COALESCE(level, 1) FROM child_profiles WHERE id = $1`
	err := s.DB.QueryRowContext(ctx, query, childProfileID).Scan(&xpTotal, &level)
	if err != nil {
		return 0, 0, fmt.Errorf("get xp and level: %w", err)
	}
	return xpTotal, level, nil
}

// XPForLevel рассчитывает сколько XP нужно для следующего уровня
// Формула: 50 × level² + 50 × level
func XPForLevel(level int) int {
	return 50*level*level + 50*level
}

// XPProgress рассчитывает прогресс до следующего уровня
// Возвращает: (текущий XP в уровне, XP нужно для уровня)
func XPProgress(xpTotal int, level int) (int, int) {
	xpNeeded := XPForLevel(level)
	currentLevelXP := xpTotal - XPForLevel(level-1)
	if currentLevelXP < 0 {
		currentLevelXP = 0
	}
	return currentLevelXP, xpNeeded
}
