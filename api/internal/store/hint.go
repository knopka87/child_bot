package store

import (
	"context"
	"time"
)

type HintCache struct {
	SessionID string
	Engine    string
	Level     string
	HintJson  []byte
	CreatedAt time.Time
}

// FindHintBySID возвращает кэш подсказки указанного уровня (1..3) для session_id.
func (s *Store) FindHintBySID(ctx context.Context, sid string, level int) (HintCache, error) {
	const q = `select hint_json, created_at
	           from hints_cache
	           where session_id=$1 and level=$2`
	var hc HintCache
	if err := s.DB.QueryRowContext(ctx, q, sid, level).Scan(&hc.HintJson, &hc.CreatedAt); err != nil {
		return HintCache{}, err
	}
	return hc, nil
}

// UpsertHint сохраняет/обновляет подсказку указанного уровня.
// UK: (session_id, level).
func (s *Store) UpsertHint(ctx context.Context, hc HintCache) error {
	const q = `
insert into hints_cache(session_id, engine, level, hint_json)
values ($1,$2,$3,$4)
on conflict (session_id, level)
do update set 
    hint_json=excluded.hint_json,
    engine = excluded.engine,
    created_at=now()`
	_, err := s.DB.ExecContext(ctx, q, hc.SessionID, hc.Engine, hc.Level, hc.HintJson)
	return err
}
