package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"child-bot/api/internal/util"
)

// MetricEvent — одна точка метрик.
// details — свободная структура для этап-специфичных полей.
type MetricEvent struct {
	Stage       string // detect|parse|hint|normalize|check|analogue
	Provider    string // gemini|gpt|yandex|...
	OK          bool
	Error       string // короткая причина (если ошибка)
	HTTPCode    *int   // статус код, если релевантно
	DurationMS  int64
	ChatID      *int64
	UserIDAnon  *int64
	TaskID      string
	Correlation string         // trace/correlation id (опц.)
	RequestID   string         // id от провайдера (опц.)
	Details     map[string]any // любые поля, напр. {"source":"photo","mime":"image/jpeg","bytes":1234}
	CreatedAt   time.Time
}

func (me *MetricEvent) ChatIDStr() string {
	if me.ChatID == nil {
		return ""
	}
	return strconv.FormatInt(*me.ChatID, 10)
}

type MetricsRepo struct{ db *sql.DB }

func NewMetricsRepo(db *sql.DB) *MetricsRepo { return &MetricsRepo{db: db} }

func (r *MetricsRepo) InsertEvent(ctx context.Context, ev MetricEvent) error {
	if ev.CreatedAt.IsZero() {
		ev.CreatedAt = time.Now()
	}
	var jb []byte
	if ev.Details == nil {
		jb = []byte("{}")
	} else {
		b, err := json.Marshal(ev.Details)
		if err != nil {
			jb = []byte(fmt.Sprintf("{}"))
			util.PrintError("InsertEvent", ev.Provider, *ev.ChatID, "error", err)
		} else {
			jb = b
		}
	}

	const q = `
	INSERT INTO metrics_events(
	    created_at, stage, provider, ok, error, http_code, duration_ms,
	    chat_id, user_id_anon, task_id, correlation_id, request_id, details
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13);
	`
	_, err := r.db.ExecContext(ctx, q,
		ev.CreatedAt,
		ev.Stage,
		ev.Provider,
		ev.OK,
		nullIfEmpty(ev.Error),
		ev.HTTPCode,
		ev.DurationMS,
		ev.ChatIDStr(),
		ev.UserIDAnon,
		nullIfEmpty(ev.TaskID),
		nullIfEmpty(ev.Correlation),
		nullIfEmpty(ev.RequestID),
		string(jb),
	)
	return err
}

func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}
