package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type TaskSession struct {
	ChatID    int64
	SessionID string
	UpdatedAt time.Time
}

func (s *Store) UpsertSession(ctx context.Context, ts TaskSession) error {
	if ts.UpdatedAt.IsZero() {
		ts.UpdatedAt = time.Now()
	}

	_, err := s.DB.ExecContext(ctx, `
	INSERT INTO task_sessions
	(chat_id, session_id, updated_at)
	VALUES ($1,$2,$3)
	ON CONFLICT (chat_id) DO UPDATE
	SET session_id = EXCLUDED.session_id,
	    updated_at = EXCLUDED.updated_at
	`, ts.ChatID, ts.SessionID, ts.UpdatedAt)

	return err
}

func (s *Store) FindSession(ctx context.Context, chatID int64) (TaskSession, error) {
	const q = `SELECT session_id, updated_at
				FROM task_sessions
				WHERE chat_id=$1`
	var sessionID string
	var updatedAt time.Time
	err := s.DB.QueryRowContext(ctx, q, chatID).Scan(&sessionID, &updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Нет записи — возвращаем "пустую" сессиию и nil-ошибку.
			// Вызывающая сторона может проверить SessionID == "" как признак отсутствия.
			return TaskSession{}, nil
		}
		return TaskSession{}, err
	}

	return TaskSession{
		ChatID:    chatID,
		SessionID: sessionID,
		UpdatedAt: updatedAt,
	}, nil
}
