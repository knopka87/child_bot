package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// VillainStore работа с таблицами villains, villain_battles, damage_events
type VillainStore struct {
	db *sql.DB
}

// NewVillainStore создает новый VillainStore
func NewVillainStore(db *sql.DB) *VillainStore {
	return &VillainStore{db: db}
}

// VillainRow запись злодея из справочника
type VillainRow struct {
	ID                   string
	Name                 string
	Description          string
	ImageURL             string
	MaxHP                int
	Level                int
	DamagePerCorrectTask int
	UnlockOrder          int
	RewardCoins          int
	RewardAchievementID  sql.NullString
	CreatedAt            time.Time
}

// VillainBattleRow запись битвы со злодеем
type VillainBattleRow struct {
	ID                 int64
	ChildProfileID     string
	VillainID          string
	Status             string // active, defeated, abandoned
	CurrentHP          int
	TotalDamageDealt   int
	CorrectTasksCount  int
	RewardsClaimed     bool
	StartedAt          time.Time
	DefeatedAt         sql.NullTime
	UpdatedAt          time.Time
}

// DamageEventRow запись урона
type DamageEventRow struct {
	ID        int64
	BattleID  int64
	AttemptID sql.NullString
	Damage    int
	TaskType  string
	CreatedAt time.Time
}

// GetActiveVillainBattle получает активную битву со злодеем
func (s *VillainStore) GetActiveVillainBattle(ctx context.Context, childProfileID string) (*VillainBattleRow, *VillainRow, error) {
	query := `
		SELECT
			vb.id, vb.child_profile_id, vb.villain_id, vb.status,
			vb.current_hp, vb.total_damage_dealt, vb.correct_tasks_count,
			vb.rewards_claimed, vb.started_at, vb.defeated_at, vb.updated_at,
			v.id, v.name, v.description, v.image_url, v.max_hp, v.level,
			v.damage_per_correct_task, v.unlock_order, v.reward_coins,
			v.reward_achievement_id, v.created_at
		FROM villain_battles vb
		JOIN villains v ON vb.villain_id = v.id
		WHERE vb.child_profile_id = $1
			AND vb.status = 'active'
		ORDER BY vb.started_at DESC
		LIMIT 1
	`

	var battle VillainBattleRow
	var villain VillainRow

	err := s.db.QueryRowContext(ctx, query, childProfileID).Scan(
		&battle.ID, &battle.ChildProfileID, &battle.VillainID, &battle.Status,
		&battle.CurrentHP, &battle.TotalDamageDealt, &battle.CorrectTasksCount,
		&battle.RewardsClaimed, &battle.StartedAt, &battle.DefeatedAt, &battle.UpdatedAt,
		&villain.ID, &villain.Name, &villain.Description, &villain.ImageURL,
		&villain.MaxHP, &villain.Level, &villain.DamagePerCorrectTask,
		&villain.UnlockOrder, &villain.RewardCoins, &villain.RewardAchievementID,
		&villain.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil, nil // Нет активной битвы
	}
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get active villain battle: %w", err)
	}

	return &battle, &villain, nil
}

// GetDamageEvents получает последние события урона для битвы
func (s *VillainStore) GetDamageEvents(ctx context.Context, battleID int64, limit int) ([]DamageEventRow, error) {
	query := `
		SELECT id, battle_id, attempt_id, damage, task_type, created_at
		FROM damage_events
		WHERE battle_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := s.db.QueryContext(ctx, query, battleID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query damage events: %w", err)
	}
	defer rows.Close()

	var events []DamageEventRow
	for rows.Next() {
		var event DamageEventRow
		if err := rows.Scan(&event.ID, &event.BattleID, &event.AttemptID, &event.Damage, &event.TaskType, &event.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan damage event: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return events, nil
}

// DealDamage наносит урон злодею и записывает событие
func (s *VillainStore) DealDamage(ctx context.Context, battleID int64, attemptID string, damage int, taskType string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Обновляем битву
	updateQuery := `
		UPDATE villain_battles
		SET current_hp = GREATEST(current_hp - $1, 0),
			total_damage_dealt = total_damage_dealt + $1,
			correct_tasks_count = correct_tasks_count + 1,
			updated_at = NOW()
		WHERE id = $2
	`

	if _, err := tx.ExecContext(ctx, updateQuery, damage, battleID); err != nil {
		return fmt.Errorf("failed to update villain battle: %w", err)
	}

	// Записываем событие урона
	insertQuery := `
		INSERT INTO damage_events (battle_id, attempt_id, damage, task_type)
		VALUES ($1, $2, $3, $4)
	`

	if _, err := tx.ExecContext(ctx, insertQuery, battleID, attemptID, damage, taskType); err != nil {
		return fmt.Errorf("failed to insert damage event: %w", err)
	}

	// Проверяем победу (HP <= 0)
	var currentHP int
	checkQuery := `SELECT current_hp FROM villain_battles WHERE id = $1`
	if err := tx.QueryRowContext(ctx, checkQuery, battleID).Scan(&currentHP); err != nil {
		return fmt.Errorf("failed to check HP: %w", err)
	}

	if currentHP <= 0 {
		// Злодей побежден
		defeatQuery := `
			UPDATE villain_battles
			SET status = 'defeated',
				defeated_at = NOW(),
				updated_at = NOW()
			WHERE id = $1
		`
		if _, err := tx.ExecContext(ctx, defeatQuery, battleID); err != nil {
			return fmt.Errorf("failed to mark villain as defeated: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
