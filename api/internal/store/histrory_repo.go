package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

type HistoryRepo struct{ DB *sql.DB }

func NewHistoryRepo(db *sql.DB) *HistoryRepo {
	return &HistoryRepo{DB: db}
}

type TimelineEvent struct {
	ChatID        int64
	TaskSessionID string
	Direction     string // in|out|api
	EventType     string
	Provider      string
	OK            bool
	LatencyMS     *int64
	TgMessageID   *int
	Text          string
	InputPayload  any
	OutputPayload any
	Error         error
	CreatedAt     time.Time
}

func (r *HistoryRepo) Insert(ctx context.Context, e TimelineEvent) error {
	var inb, outb []byte
	if e.InputPayload != nil {
		inb, _ = json.Marshal(e.InputPayload)
	}
	if e.OutputPayload != nil {
		outb, _ = json.Marshal(e.OutputPayload)
	}
	errStr := ""
	if e.Error != nil {
		errStr = e.Error.Error()
	}
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now()
	}

	_, err := r.DB.ExecContext(ctx, `
    INSERT INTO timeline_events
      (chat_id, task_session_id, direction, event_type, provider, ok, latency_ms,
       tg_message_id, text, input_payload, output_payload, error)
    VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
  `,
		e.ChatID, e.TaskSessionID, e.Direction, e.EventType, e.Provider, e.OK, e.LatencyMS,
		e.TgMessageID, e.Text, inb, outb, errStr,
	)
	return err
}

func (r *HistoryRepo) FindALLRecordsBySessionID(ctx context.Context, sid string) ([]TimelineEvent, error) {
	const q = `SELECT
	    chat_id,
	    task_session_id,
	    direction,
	    event_type,
	    provider,
	    ok,
	    latency_ms,
	    tg_message_id,
	    text,
	    input_payload,
	    output_payload,
	    error,
	    created_at
	  FROM timeline_events
	  WHERE task_session_id = $1
	  ORDER BY created_at`

	rows, err := r.DB.QueryContext(ctx, q, sid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []TimelineEvent
	for rows.Next() {
		var e TimelineEvent
		var (
			latency   sql.NullInt64
			tgMsgID   sql.NullInt64
			inb, outb []byte
			errStr    sql.NullString
			createdAt time.Time
		)

		if scanErr := rows.Scan(
			&e.ChatID,
			&e.TaskSessionID,
			&e.Direction,
			&e.EventType,
			&e.Provider,
			&e.OK,
			&latency,
			&tgMsgID,
			&e.Text,
			&inb,
			&outb,
			&errStr,
			&createdAt,
		); scanErr != nil {
			return nil, scanErr
		}

		if latency.Valid {
			e.LatencyMS = &latency.Int64
		}
		if tgMsgID.Valid {
			v := int(tgMsgID.Int64)
			e.TgMessageID = &v
		}
		if len(inb) > 0 {
			var anyIn any
			_ = json.Unmarshal(inb, &anyIn)
			e.InputPayload = anyIn
		}
		if len(outb) > 0 {
			var anyOut any
			_ = json.Unmarshal(outb, &anyOut)
			e.OutputPayload = anyOut
		}
		if errStr.Valid && errStr.String != "" {
			e.Error = errors.New(errStr.String)
		}
		e.CreatedAt = createdAt

		res = append(res, e)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return res, nil
}
