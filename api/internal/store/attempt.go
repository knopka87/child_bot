package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"child-bot/api/internal/llm/types"

	"github.com/google/uuid"
)

// AttemptStore работает с попытками в БД
type AttemptStore struct {
	db *sql.DB
}

// NewAttemptStore создаёт новый AttemptStore
func NewAttemptStore(db *sql.DB) *AttemptStore {
	return &AttemptStore{db: db}
}

// Attempt модель попытки в БД
type Attempt struct {
	ID                uuid.UUID
	ChildProfileID    uuid.UUID
	AttemptType       string // help или check
	Status            string // created, processing, completed, failed
	TaskImageURL      sql.NullString
	AnswerImageURL    sql.NullString
	DetectResult      []byte // JSONB - может быть NULL (будет пустой слайс)
	ParseResult       []byte // JSONB - может быть NULL
	HintsResult       []byte // JSONB - может быть NULL
	CheckResult       []byte // JSONB - может быть NULL
	CurrentHintIndex  int
	HintsUsed         int
	TimeSpentSeconds  sql.NullInt64
	IsCorrect         sql.NullBool
	HasErrors         sql.NullBool
	CreatedAt         time.Time
	UpdatedAt         time.Time
	CompletedAt       sql.NullTime
}

// CreateAttempt создаёт новую попытку
func (s *AttemptStore) CreateAttempt(ctx context.Context, childProfileID uuid.UUID, attemptType string) (uuid.UUID, error) {
	query := `
		INSERT INTO attempts (child_profile_id, attempt_type, status)
		VALUES ($1, $2, 'created')
		RETURNING id, created_at, updated_at
	`

	var attempt Attempt
	err := s.db.QueryRowContext(ctx, query, childProfileID, attemptType).Scan(
		&attempt.ID,
		&attempt.CreatedAt,
		&attempt.UpdatedAt,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create attempt: %w", err)
	}

	return attempt.ID, nil
}

// UpdateTaskImage обновляет изображение задания
func (s *AttemptStore) UpdateTaskImage(ctx context.Context, attemptID uuid.UUID, imageURL string) error {
	query := `
		UPDATE attempts
		SET task_image_url = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := s.db.ExecContext(ctx, query, imageURL, attemptID)
	if err != nil {
		return fmt.Errorf("failed to update task image: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("attempt not found: %s", attemptID)
	}

	return nil
}

// UpdateAnswerImage обновляет изображение ответа
func (s *AttemptStore) UpdateAnswerImage(ctx context.Context, attemptID uuid.UUID, imageURL string) error {
	query := `
		UPDATE attempts
		SET answer_image_url = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := s.db.ExecContext(ctx, query, imageURL, attemptID)
	if err != nil {
		return fmt.Errorf("failed to update answer image: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("attempt not found: %s", attemptID)
	}

	return nil
}

// UpdateStatus обновляет статус попытки
func (s *AttemptStore) UpdateStatus(ctx context.Context, attemptID uuid.UUID, status string) error {
	query := `
		UPDATE attempts
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := s.db.ExecContext(ctx, query, status, attemptID)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("attempt not found: %s", attemptID)
	}

	return nil
}

// SaveDetectResult сохраняет результат Detect
func (s *AttemptStore) SaveDetectResult(ctx context.Context, attemptID uuid.UUID, result *types.DetectResponse) error {
	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal detect result: %w", err)
	}

	query := `
		UPDATE attempts
		SET detect_result = $1, updated_at = NOW()
		WHERE id = $2
	`

	_, err = s.db.ExecContext(ctx, query, data, attemptID)
	if err != nil {
		return fmt.Errorf("failed to save detect result: %w", err)
	}

	return nil
}

// SaveParseResult сохраняет результат Parse
func (s *AttemptStore) SaveParseResult(ctx context.Context, attemptID uuid.UUID, result *types.ParseResponse) error {
	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal parse result: %w", err)
	}

	query := `
		UPDATE attempts
		SET parse_result = $1, updated_at = NOW()
		WHERE id = $2
	`

	_, err = s.db.ExecContext(ctx, query, data, attemptID)
	if err != nil {
		return fmt.Errorf("failed to save parse result: %w", err)
	}

	return nil
}

// SaveHintsResult сохраняет результат Hint
func (s *AttemptStore) SaveHintsResult(ctx context.Context, attemptID uuid.UUID, result *types.HintResponse) error {
	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal hints result: %w", err)
	}

	query := `
		UPDATE attempts
		SET hints_result = $1, status = 'completed', completed_at = NOW(), updated_at = NOW()
		WHERE id = $2
	`

	_, err = s.db.ExecContext(ctx, query, data, attemptID)
	if err != nil {
		return fmt.Errorf("failed to save hints result: %w", err)
	}

	return nil
}

// SaveCheckResult сохраняет результат Check
func (s *AttemptStore) SaveCheckResult(ctx context.Context, attemptID uuid.UUID, result *types.CheckResponse) error {
	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal check result: %w", err)
	}

	isCorrect := result.Decision == types.CheckDecisionCorrect
	hasErrors := result.Decision != types.CheckDecisionCorrect

	query := `
		UPDATE attempts
		SET check_result = $1, is_correct = $2, has_errors = $3,
		    status = 'completed', completed_at = NOW(), updated_at = NOW()
		WHERE id = $4
	`

	_, err = s.db.ExecContext(ctx, query, data, isCorrect, hasErrors, attemptID)
	if err != nil {
		return fmt.Errorf("failed to save check result: %w", err)
	}

	return nil
}

// GetAttempt получает попытку по ID
func (s *AttemptStore) GetAttempt(ctx context.Context, attemptID uuid.UUID) (*Attempt, error) {
	query := `
		SELECT id, child_profile_id, attempt_type, status,
		       task_image_url, answer_image_url,
		       detect_result, parse_result, hints_result, check_result,
		       current_hint_index, hints_used, time_spent_seconds,
		       is_correct, has_errors,
		       created_at, updated_at, completed_at
		FROM attempts
		WHERE id = $1
	`

	var attempt Attempt
	err := s.db.QueryRowContext(ctx, query, attemptID).Scan(
		&attempt.ID,
		&attempt.ChildProfileID,
		&attempt.AttemptType,
		&attempt.Status,
		&attempt.TaskImageURL,
		&attempt.AnswerImageURL,
		&attempt.DetectResult,
		&attempt.ParseResult,
		&attempt.HintsResult,
		&attempt.CheckResult,
		&attempt.CurrentHintIndex,
		&attempt.HintsUsed,
		&attempt.TimeSpentSeconds,
		&attempt.IsCorrect,
		&attempt.HasErrors,
		&attempt.CreatedAt,
		&attempt.UpdatedAt,
		&attempt.CompletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("attempt not found: %s", attemptID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get attempt: %w", err)
	}

	return &attempt, nil
}

// GetUnfinishedAttempt получает незавершённую попытку пользователя
func (s *AttemptStore) GetUnfinishedAttempt(ctx context.Context, childProfileID uuid.UUID) (*Attempt, error) {
	query := `
		SELECT id, child_profile_id, attempt_type, status,
		       task_image_url, answer_image_url,
		       detect_result, parse_result, hints_result, check_result,
		       current_hint_index, hints_used, time_spent_seconds,
		       is_correct, has_errors,
		       created_at, updated_at, completed_at
		FROM attempts
		WHERE child_profile_id = $1 AND status IN ('created', 'processing')
		ORDER BY created_at DESC
		LIMIT 1
	`

	var attempt Attempt
	err := s.db.QueryRowContext(ctx, query, childProfileID).Scan(
		&attempt.ID,
		&attempt.ChildProfileID,
		&attempt.AttemptType,
		&attempt.Status,
		&attempt.TaskImageURL,
		&attempt.AnswerImageURL,
		&attempt.DetectResult,
		&attempt.ParseResult,
		&attempt.HintsResult,
		&attempt.CheckResult,
		&attempt.CurrentHintIndex,
		&attempt.HintsUsed,
		&attempt.TimeSpentSeconds,
		&attempt.IsCorrect,
		&attempt.HasErrors,
		&attempt.CreatedAt,
		&attempt.UpdatedAt,
		&attempt.CompletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil // Нет незавершённой попытки - это нормально
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get unfinished attempt: %w", err)
	}

	return &attempt, nil
}

// GetRecentAttempts получает последние попытки пользователя
func (s *AttemptStore) GetRecentAttempts(ctx context.Context, childProfileID uuid.UUID, limit int) ([]*Attempt, error) {
	query := `
		SELECT id, child_profile_id, attempt_type, status,
		       task_image_url, answer_image_url,
		       detect_result, parse_result, hints_result, check_result,
		       current_hint_index, hints_used, time_spent_seconds,
		       is_correct, has_errors,
		       created_at, updated_at, completed_at
		FROM attempts
		WHERE child_profile_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := s.db.QueryContext(ctx, query, childProfileID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent attempts: %w", err)
	}
	defer rows.Close()

	var attempts []*Attempt
	for rows.Next() {
		var attempt Attempt
		err := rows.Scan(
			&attempt.ID,
			&attempt.ChildProfileID,
			&attempt.AttemptType,
			&attempt.Status,
			&attempt.TaskImageURL,
			&attempt.AnswerImageURL,
			&attempt.DetectResult,
			&attempt.ParseResult,
			&attempt.HintsResult,
			&attempt.CheckResult,
			&attempt.CurrentHintIndex,
			&attempt.HintsUsed,
			&attempt.TimeSpentSeconds,
			&attempt.IsCorrect,
			&attempt.HasErrors,
			&attempt.CreatedAt,
			&attempt.UpdatedAt,
			&attempt.CompletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan attempt: %w", err)
		}
		attempts = append(attempts, &attempt)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return attempts, nil
}

// DeleteAttempt удаляет попытку
func (s *AttemptStore) DeleteAttempt(ctx context.Context, attemptID uuid.UUID) error {
	query := `DELETE FROM attempts WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, attemptID)
	if err != nil {
		return fmt.Errorf("failed to delete attempt: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("attempt not found: %s", attemptID)
	}

	return nil
}
