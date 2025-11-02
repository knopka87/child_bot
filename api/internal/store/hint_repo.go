package store

import (
	"context"
	"database/sql"
	"time"
)

type HintCache struct {
	ImageHash string
	Engine    string
	Level     string
	HintJson  []byte
	CreatedAt time.Time
}

type HintRepo struct{ DB *sql.DB }

func NewHintRepo(db *sql.DB) *HintRepo { return &HintRepo{DB: db} }

// Find возвращает кэш подсказки указанного уровня (1..3) для (imageHash, engine).
// Если maxAge > 0 и запись старше, вернёт sql.ErrNoRows (чтобы вызвать LLMClient заново).
func (r *HintRepo) Find(ctx context.Context, imageHash, engine string, level int) (HintCache, error) {
	const q = `select hint_json, created_at
	           from hints_cache
	           where image_hash=$1 and engine=$2 and level=$3`
	var hc HintCache
	if err := r.DB.QueryRowContext(ctx, q, imageHash, engine, level).Scan(&hc.HintJson, &hc.CreatedAt); err != nil {
		return HintCache{}, err
	}
	return hc, nil
}

// Upsert сохраняет/обновляет подсказку указанного уровня.
// PK: (image_hash, engine, model, level).
func (r *HintRepo) Upsert(ctx context.Context, hc HintCache) error {
	const q = `
insert into hints_cache(image_hash, engine, level, hint_json)
values ($1,$2,$3,$4)
on conflict (image_hash, engine, level)
do update set hint_json=excluded.hint_json, created_at=now()`
	_, err := r.DB.ExecContext(ctx, q, hc.ImageHash, hc.Engine, hc.Level, hc.HintJson)
	return err
}
