package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

type ParsedTasks struct {
	ID                    int64           `db:"id"`
	SessionID             string          `db:"session_id"`
	CreatedAt             time.Time       `db:"created_at"`
	UpdatedAt             time.Time       `db:"created_at"`
	ChatID                int64           `db:"chat_id"`
	MediaGroupID          string          `db:"media_group_id"`
	ImageHash             string          `db:"image_hash"`
	Engine                string          `db:"engine"`
	Subject               string          `db:"subject"`
	Grade                 int             `db:"grade_hint"`
	RawTaskText           string          `db:"raw_task_text"`
	Question              string          `db:"question"`
	ResultJSON            json.RawMessage `db:"result_json"`
	NeedsUserConfirmation bool            `db:"needs_user_confirmation"`
	TaskType              string          `db:"task_type"`
	CombinedSubpoints     bool            `db:"combined_subpoints"`
	Confidence            float64         `db:"confidence"`
	Accepted              bool            `db:"accepted"`
	AcceptReason          string          `db:"accept_reason"`
}

var ErrNotFound = sql.ErrNoRows

type ParseRepo struct{ DB *sql.DB }

func NewParseRepo(db *sql.DB) *ParseRepo { return &ParseRepo{DB: db} }

// FindByChatID достаёт самую свежую запись по ключу (image_hash + engine).
// Если maxAge > 0 — проверяет "свежесть", иначе игнорирует возраст.
func (r *ParseRepo) FindByChatID(ctx context.Context, chatID int64) (*ParsedTasks, error) {
	const q = `
select pt.id,
       pt.created_at,
       pt.updated_at,
       pt.session_id,
       coalesce(pt.media_group_id,'') as media_group_id,
       pt.image_hash,
       pt.engine,
       pt.subject,
       pt.grade_hint,
       pt.raw_task_text,
       pt.question,
       pt.result_json,
       pt.needs_user_confirmation,
       pt.task_type,
       pt.combined_subpoints,
       pt.accepted,
       pt."accept_reason",
       pt.confidence
from parsed_tasks pt 
    left join task_sessions ts 
        ON pt.session_id = ts.session_id and ts.chat_id = pt.chat_id
where pt.chat_id = $1
limit 1`

	row := r.DB.QueryRowContext(ctx, q, chatID)

	var (
		id           int64
		createdAt    time.Time
		updatedAt    time.Time
		sessionID    string
		mediaGroupID string
		imgHash      string
		engName      string
		subject      string
		grade        int
		rawText      string
		question     string
		jsonBlob     []byte
		needConf     bool
		taskType     string
		combined     bool
		accepted     bool
		accReason    sql.NullString
		confidence   float64
	)

	if err := row.Scan(&id, &createdAt, &updatedAt, &sessionID, &mediaGroupID, &imgHash, &engName,
		&subject, &grade, &rawText, &question, &jsonBlob, &needConf, &taskType, &combined, &accepted, &accReason, &confidence); err != nil {
		return nil, err
	}

	return &ParsedTasks{
		ID:                    id,
		CreatedAt:             createdAt,
		UpdatedAt:             updatedAt,
		SessionID:             sessionID,
		ChatID:                chatID,
		MediaGroupID:          mediaGroupID,
		ImageHash:             imgHash,
		Engine:                engName,
		Subject:               subject,
		Grade:                 grade,
		RawTaskText:           rawText,
		Question:              question,
		ResultJSON:            json.RawMessage(jsonBlob),
		NeedsUserConfirmation: needConf,
		TaskType:              taskType,
		CombinedSubpoints:     combined,
		Accepted:              accepted,
		AcceptReason:          accReason.String,
		Confidence:            confidence,
	}, nil
}

// Upsert сохраняет PARSE (черновик или принятый).
// Если запись по (image_hash, engine) существует — обновит все поля.
func (r *ParseRepo) Upsert(
	ctx context.Context,
	pr ParsedTasks,
) error {

	const q = `
insert into parsed_tasks (
  chat_id,
  session_id,
  media_group_id,
  image_hash,
  engine,
  subject,
  grade_hint,
  raw_task_text,
  question,
  result_json,
  needs_user_confirmation,
  task_type,
  combined_subpoints,
  accepted,
  accept_reason,
  confidence,
  created_at,
  updated_at
) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$17)
on conflict (image_hash, engine) do update set
  chat_id = excluded.chat_id,
  session_id = excluded.session_id,
  media_group_id = excluded.media_group_id,
  subject = excluded.subject,
  grade_hint = excluded.grade_hint,
  raw_task_text = excluded.raw_task_text,
  question = excluded.question,
  result_json = excluded.result_json,
  needs_user_confirmation = excluded.needs_user_confirmation,
  task_type = excluded.task_type,
  combined_subpoints = excluded.combined_subpoints,
  accepted = excluded.accepted,
  accept_reason = excluded.accept_reason,
  confidence = excluded.confidence,
  updated_at = NOW()
  `
	_, err := r.DB.ExecContext(ctx, q,
		pr.ChatID,
		pr.MediaGroupID,
		pr.ImageHash,
		pr.Engine,
		pr.Subject,
		pr.Grade,
		pr.RawTaskText,
		pr.Question,
		pr.ResultJSON,
		pr.NeedsUserConfirmation,
		pr.TaskType,
		pr.CombinedSubpoints,
		pr.Accepted,
		pr.AcceptReason,
		pr.Confidence,
		pr.CreatedAt,
	)
	return err
}

// MarkAccepted помечает существующую запись как принятую (без изменения JSON).
func (r *ParseRepo) MarkAccepted(ctx context.Context, imageHash, engine, reason string) error {
	const q = `update parsed_tasks set accepted=true, accept_reason=$3, updated_at=NOW() where image_hash=$1 and engine=$2`
	res, err := r.DB.ExecContext(ctx, q, imageHash, engine, reason)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return ErrNotFound
	}
	return nil
}

// PurgeOlderThan удаляет очень старые записи-кэши, чтобы не раздувать БД.
func (r *ParseRepo) PurgeOlderThan(ctx context.Context, olderThan time.Duration) (int64, error) {
	if olderThan <= 0 {
		return 0, errors.New("olderThan must be > 0")
	}
	cutoff := time.Now().Add(-olderThan)
	const q = `delete from parsed_tasks where created_at < $1`
	res, err := r.DB.ExecContext(ctx, q, cutoff)
	if err != nil {
		return 0, err
	}
	aff, _ := res.RowsAffected()
	return aff, nil
}

// FindLastConfirmed возвращает последнюю ПРИНЯТУЮ (accepted=true) запись по chat_id.
// Удобно для шагов hint/check/analogue, где нужна подтверждённая формулировка.
func (r *ParseRepo) FindLastConfirmed(ctx context.Context, chatID int64) (*ParsedTasks, bool) {
	pr, err := r.FindByChatID(ctx, chatID)
	if err != nil || !pr.Accepted {
		return nil, false
	}
	return pr, true
}
