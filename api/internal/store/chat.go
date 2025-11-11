package store

import (
	"context"
	"database/sql"
)

type Chat struct {
	ID        int64
	Type      *string
	Username  *string
	FirstName *string
	LastName  *string
}

// FindChatByID returns a single chat by its primary key (id).
func (s *Store) FindChatByID(ctx context.Context, id int64) (Chat, error) {
	const q = `
select id, type, username, first_name, last_name
from chat
where id = $1`
	var (
		cid       int64
		t         sql.NullString
		username  sql.NullString
		firstname sql.NullString
		lastname  sql.NullString
	)
	if err := s.DB.QueryRowContext(ctx, q, id).Scan(
		&cid,
		&t,
		&username,
		&firstname,
		&lastname,
	); err != nil {
		return Chat{}, err
	}
	c := Chat{
		ID: cid,
	}
	if t.Valid {
		c.Type = &t.String
	}
	if username.Valid {
		c.Username = &username.String
	}
	if firstname.Valid {
		c.FirstName = &firstname.String
	}
	if lastname.Valid {
		c.LastName = &lastname.String
	}

	return c, nil
}

// UpsertChat inserts/updates a chat by id (PK: id).
func (s *Store) UpsertChat(ctx context.Context, c Chat) error {
	const q = `
insert into chat (id, type, username, first_name, last_name)
values ($1, $2, $3, $4, $5)
on conflict (id) do update
set
	type       = excluded.type,
	username   = excluded.username,
	first_name = excluded.first_name,
	last_name  = excluded.last_name`
	_, err := s.DB.ExecContext(ctx, q,
		c.ID,
		c.Type,
		c.Username,
		c.FirstName,
		c.LastName,
	)
	return err
}
