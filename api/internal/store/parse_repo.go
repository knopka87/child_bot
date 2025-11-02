package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type ParsedTasks struct {
	ID                    int64           `db:"id"`
	CreatedAt             time.Time       `db:"created_at"`
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

// FindByHash достаёт самую свежую запись по ключу (image_hash + engine).
// Если maxAge > 0 — проверяет "свежесть", иначе игнорирует возраст.
func (r *ParseRepo) FindByHash(ctx context.Context, imageHash, engine string, maxAge time.Duration) (*ParsedTasks, error) {
	const q = `
select id,
       created_at,
       coalesce(chat_id,0) as chat_id,
       coalesce(media_group_id,'') as media_group_id,
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
       confidence
from parsed_tasks
where image_hash = $1 and engine = $2
order by created_at desc
limit 1`

	row := r.DB.QueryRowContext(ctx, q, imageHash, engine)

	var (
		id           int64
		createdAt    time.Time
		chatID       int64
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

	if err := row.Scan(&id, &createdAt, &chatID, &mediaGroupID, &imgHash, &engName,
		&subject, &grade, &rawText, &question, &jsonBlob, &needConf, &taskType, &combined, &accepted, &accReason, &confidence); err != nil {
		return nil, err
	}

	if maxAge > 0 && time.Since(createdAt) > maxAge {
		return nil, ErrNotFound
	}

	return &ParsedTasks{
		ID:                    id,
		CreatedAt:             createdAt,
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
  confidence
) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
on conflict (image_hash, engine) do update set
  chat_id = excluded.chat_id,
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
  confidence = excluded.confidence
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
	)
	return err
}

// MarkAccepted помечает существующую запись как принятую (без изменения JSON).
func (r *ParseRepo) MarkAccepted(ctx context.Context, imageHash, engine, reason string) error {
	const q = `update parsed_tasks set accepted=true, accept_reason=$3 where image_hash=$1 and engine=$2`
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

// AcceptWithOverwrite — принять PARSE и одновременно переписать JSON/текст.
// Удобно при сценарии "Нет" + текстовая правка пользователя.
func (r *ParseRepo) AcceptWithOverwrite(
	ctx context.Context,
	pr ParsedTasks,
) error {
	const q = `
insert into parsed_tasks (
  chat_id,
  media_group_id,
  image_hash,
  engine,
  subject,
  grade_hint,
  raw_task_text,
  result_json,
  needs_user_confirmation,
  task_type,
  combined_subpoints,
  accepted,
  accept_reason
) values ($1,$2,$3,$4,$5,$6,$7,false,$8,$9,true,$10)
on conflict (image_hash, engine) do update set
  chat_id = excluded.chat_id,
  media_group_id = excluded.media_group_id,
  subject = excluded.subject,
  raw_task_text = excluded.raw_task_text,
  result_json = excluded.result_json,
  needs_user_confirmation = false,
  task_type = excluded.task_type,
  combined_subpoints = excluded.combined_subpoints,
  accepted = true,
  accept_reason = excluded.accept_reason`
	_, err := r.DB.ExecContext(ctx, q,
		pr.ChatID, pr.MediaGroupID, pr.ImageHash, pr.Engine,
		sql.NullString{String: strings.ToLower(pr.Subject), Valid: pr.Subject != ""},
		pr.RawTaskText,
		pr.ResultJSON,
		sql.NullString{String: pr.TaskType, Valid: pr.TaskType != ""},
		true,
		pr.AcceptReason,
	)
	return err
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
	const q = `
select id,
       created_at,
       coalesce(chat_id,0) as chat_id,
       coalesce(media_group_id,'') as media_group_id,
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
       confidence
from parsed_tasks
where chat_id = $1 and accepted = true
order by created_at desc
limit 1`

	row := r.DB.QueryRowContext(ctx, q, chatID)

	var (
		id         int64
		createdAt  time.Time
		cid        int64
		mgid       sql.NullString
		imgHash    string
		engine     string
		subject    string
		grade      int
		rawText    string
		question   string
		jsonBlob   []byte
		needConf   bool
		taskType   string
		combined   bool
		accepted   bool
		accReason  sql.NullString
		confidence float64
	)

	if err := row.Scan(&id, &createdAt, &cid, &mgid, &imgHash, &engine, &subject, &grade, &rawText, &question, &jsonBlob, &needConf, &taskType, &combined, &accepted, &accReason, &confidence); err != nil {
		return nil, false
	}

	return &ParsedTasks{
		ID:                    id,
		CreatedAt:             createdAt,
		ChatID:                cid,
		MediaGroupID:          mgid.String,
		ImageHash:             imgHash,
		Engine:                engine,
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
	}, true
}
