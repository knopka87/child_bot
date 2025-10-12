package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"child-bot/api/internal/ocr"
)

var ErrNotFound = sql.ErrNoRows

type ParseRepo struct{ DB *sql.DB }

func NewParseRepo(db *sql.DB) *ParseRepo { return &ParseRepo{DB: db} }

// ParsedRow — то, что чаще всего нужно наверх.
type ParsedRow struct {
	ID                 int64
	CreatedAt          time.Time
	ChatID             int64
	MediaGroupID       string
	ImageHash          string
	Engine             string
	Parse              ocr.ParseResult
	Accepted           bool
	AcceptReason       string
	ConfirmationNeeded bool
	Confidence         float64

	// Flattened поля из Parse (для удобства в handlers/router.go)
	RawText        string
	Question       string
	ShortEssence   string
	TaskID         string
	Grade          int
	Subject        string
	TaskType       string
	MethodTag      string
	DifficultyHint string
	Expected       ocr.ExpectedSolution
}

// FindByHash достаёт самую свежую запись по ключу (image_hash + engine).
// Если maxAge > 0 — проверяет "свежесть", иначе игнорирует возраст.
func (r *ParseRepo) FindByHash(ctx context.Context, imageHash, engine string, maxAge time.Duration) (*ParsedRow, error) {
	const q = `
select id, created_at,
       coalesce(chat_id,0) as chat_id,
       coalesce(media_group_id,'') as media_group_id,
       image_hash, engine,
       result_json,
       accepted, coalesce(accept_reason,'') as accept_reason,
       confirmation_needed,
       coalesce(confidence,0) as confidence
from parsed_tasks
where image_hash = $1 and engine = $2
order by created_at desc
limit 1`
	row := r.DB.QueryRowContext(ctx, q, imageHash, engine)

	var (
		id                 int64
		ts                 time.Time
		chatID             int64
		mediaGroupID       string
		imgHash            string
		engName            string
		js                 []byte
		accepted           bool
		acceptReason       string
		confirmationNeeded bool
		confidence         float64
	)
	if err := row.Scan(&id, &ts, &chatID, &mediaGroupID, &imgHash, &engName,
		&js, &accepted, &acceptReason, &confirmationNeeded, &confidence); err != nil {
		return nil, err
	}
	if maxAge > 0 && time.Since(ts) > maxAge {
		return nil, ErrNotFound
	}
	var pr ocr.ParseResult
	if err := json.Unmarshal(js, &pr); err != nil {
		// если JSON поломан — считаем, что не найдено
		return nil, ErrNotFound
	}
	return &ParsedRow{
		ID:                 id,
		CreatedAt:          ts,
		ChatID:             chatID,
		MediaGroupID:       mediaGroupID,
		ImageHash:          imgHash,
		Engine:             engName,
		Parse:              pr,
		Accepted:           accepted,
		AcceptReason:       acceptReason,
		ConfirmationNeeded: confirmationNeeded,
		Confidence:         confidence,

		RawText:        pr.RawText,
		Question:       pr.Question,
		ShortEssence:   "", // pr.ShortEssence,
		TaskID:         "", // pr.TaskID,
		Grade:          pr.Grade,
		Subject:        pr.Subject,
		TaskType:       pr.TaskType,
		MethodTag:      "", // pr.MethodTag,
		DifficultyHint: "", // pr.DifficultyHint,
		// Expected:       pr.Expected,
	}, nil
}

// Upsert сохраняет PARSE (черновик или принятый). Если запись по (image_hash, engine)
// существует — обновит все поля.
func (r *ParseRepo) Upsert(
	ctx context.Context,
	chatID int64,
	mediaGroupID, imageHash, engine string,
	pr ocr.ParseResult,
	accepted bool,
	reason string,
) error {
	js, _ := json.Marshal(pr)
	const q = `
insert into parsed_tasks (
  chat_id, media_group_id, image_hash, engine,
  raw_text, question, result_json, confidence, confirmation_needed,
  accepted, accept_reason
) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
on conflict (image_hash, engine) do update
set chat_id = excluded.chat_id,
    media_group_id = excluded.media_group_id,
    raw_text = excluded.raw_text,
    question = excluded.question,
    result_json = excluded.result_json,
    confidence = excluded.confidence,
    confirmation_needed = excluded.confirmation_needed,
    accepted = excluded.accepted,
    accept_reason = excluded.accept_reason`
	_, err := r.DB.ExecContext(ctx, q,
		chatID, mediaGroupID, imageHash, engine,
		pr.RawText, pr.Question, js, pr.Confidence, pr.ConfirmationNeeded,
		accepted, reason,
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
	chatID int64,
	mediaGroupID, imageHash, engine string,
	pr ocr.ParseResult,
	reason string,
) error {
	js, _ := json.Marshal(pr)
	const q = `
insert into parsed_tasks (
  chat_id, media_group_id, image_hash, engine,
  raw_text, question, result_json, confidence, confirmation_needed,
  accepted, accept_reason
) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,true,$10)
on conflict (image_hash, engine) do update
set chat_id = excluded.chat_id,
    media_group_id = excluded.media_group_id,
    raw_text = excluded.raw_text,
    question = excluded.question,
    result_json = excluded.result_json,
    confidence = excluded.confidence,
    confirmation_needed = excluded.confirmation_needed,
    accepted = true,
    accept_reason = excluded.accept_reason`
	_, err := r.DB.ExecContext(ctx, q,
		chatID, mediaGroupID, imageHash, engine,
		pr.RawText, pr.Question, js, pr.Confidence, pr.ConfirmationNeeded,
		reason,
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
func (r *ParseRepo) FindLastConfirmed(ctx context.Context, chatID int64) (*ParsedRow, bool) {
	const q = `
select id, created_at,
       coalesce(chat_id,0) as chat_id,
       coalesce(media_group_id,'') as media_group_id,
       coalesce(image_hash,'') as image_hash,
       coalesce(engine,'') as engine,
       coalesce(raw_text,'') as raw_text,
       coalesce(question,'') as question,
       result_json,
       accepted,
       coalesce(accept_reason,'') as accept_reason,
       confirmation_needed,
       coalesce(confidence,0) as confidence
from parsed_tasks
where chat_id = $1 and accepted = true
order by created_at desc
limit 1`
	row := r.DB.QueryRowContext(ctx, q, chatID)

	var (
		id           int64
		ts           time.Time
		cid          int64
		mgid         string
		imgHash      string
		engine       string
		rawTextField string
		questionFld  string
		jsonBlob     []byte
		accepted     bool
		reason       string
		needConfirm  bool
		conf         float64
	)
	if err := row.Scan(&id, &ts, &cid, &mgid, &imgHash, &engine,
		&rawTextField, &questionFld, &jsonBlob, &accepted, &reason, &needConfirm, &conf); err != nil {
		return nil, false
	}

	var pr ocr.ParseResult
	if err := json.Unmarshal(jsonBlob, &pr); err != nil {
		return nil, false
	}

	out := &ParsedRow{
		ID:                 id,
		CreatedAt:          ts,
		ChatID:             cid,
		MediaGroupID:       mgid,
		ImageHash:          imgHash,
		Engine:             engine,
		Parse:              pr,
		Accepted:           accepted,
		AcceptReason:       reason,
		ConfirmationNeeded: needConfirm,
		Confidence:         conf,

		RawText:  firstNonEmpty(rawTextField, pr.RawText),
		Question: firstNonEmpty(questionFld, pr.Question),
		// ShortEssence:   pr.ShortEssence,
		// TaskID:         pr.TaskID,
		Grade:    pr.Grade,
		Subject:  pr.Subject,
		TaskType: pr.TaskType,
		// MethodTag:      pr.MethodTag,
		// DifficultyHint: pr.DifficultyHint,
		// Expected:       pr.Expected,
	}
	return out, true
}

// firstNonEmpty returns a if not empty, otherwise b
func firstNonEmpty(a, b string) string {
	if strings.TrimSpace(a) != "" {
		return a
	}
	return b
}
