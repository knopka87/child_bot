package store

import (
	"context"
	"database/sql"
	"fmt"
)

type User struct {
	ID    int64
	Grade *int64
}

// FindUserByChatID returns a single user by user_id.
func (s *Store) FindUserByChatID(ctx context.Context, chatID int64) (User, error) {
	const q = `
select grade
from "user"
where chat_id = $1`
	var (
		grade sql.NullInt64
	)
	if err := s.DB.QueryRowContext(ctx, q, chatID).Scan(
		&grade,
	); err != nil {
		return User{}, err
	}

	c := User{ID: chatID}
	if !grade.Valid {
		return User{}, fmt.Errorf("invalid grade: %v", grade)
	}

	c.Grade = &grade.Int64
	return c, nil
}

// UpsertUser inserts/updates a user by chat_id (PK: chat_id).
func (s *Store) UpsertUser(ctx context.Context, c User) error {
	const q = `
insert into "user" (chat_id, grade) values ($1, $2)
on conflict (chat_id) do update
set
	grade       = excluded.grade`
	_, err := s.DB.ExecContext(ctx, q,
		c.ID,
		c.Grade,
	)
	return err
}
