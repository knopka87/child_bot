package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

type TaskSession struct {
	ChatID       int64
	SessionID    string
	UpdatedAt    time.Time
	CurrentState *string         // текущее состояние машины состояний
	ChatMode     *string         // режим чата (await_solution, etc.)
	HintContext  json.RawMessage // контекст для продолжения подсказок (JSON)
	ParseContext json.RawMessage // контекст ожидания подтверждения парсинга (JSON)
}

// HintContextData — структура для сериализации контекста подсказок
type HintContextData struct {
	ParseJSON       json.RawMessage `json:"parse_json"`                  // результат Parse
	DetectJSON      json.RawMessage `json:"detect_json"`                 // результат Detect
	EngineName      string          `json:"engine_name"`                 // имя LLM
	NextLevel       int             `json:"next_level"`                  // следующий уровень подсказки
	MaxHints        int             `json:"max_hints"`                   // максимум подсказок
	ImageBase64     string          `json:"image_base64"`                // изображение в base64 (опционально)
	CachedHintsJSON json.RawMessage `json:"cached_hints_json,omitempty"` // кэш ответа LLM со всеми подсказками
}

// ParseContextData — структура для сериализации контекста ожидания подтверждения парсинга
type ParseContextData struct {
	ImageBase64  string          `json:"image_base64"`            // изображение в base64
	Mime         string          `json:"mime"`                    // MIME тип изображения
	MediaGroupID string          `json:"media_group_id"`          // ID альбома (если есть)
	DetectJSON   json.RawMessage `json:"detect_json"`             // результат Detect
	ParseJSON    json.RawMessage `json:"parse_json"`              // результат Parse
	LLM          string          `json:"llm"`                     // имя LLM провайдера
}

func (s *Store) UpsertSession(ctx context.Context, ts TaskSession) error {
	if ts.UpdatedAt.IsZero() {
		ts.UpdatedAt = time.Now()
	}

	_, err := s.DB.ExecContext(ctx, `
	INSERT INTO task_sessions
	(chat_id, session_id, updated_at, current_state, chat_mode, hint_context)
	VALUES ($1,$2,$3,$4,$5,$6)
	ON CONFLICT (chat_id) DO UPDATE
	SET session_id = EXCLUDED.session_id,
	    updated_at = EXCLUDED.updated_at,
	    current_state = EXCLUDED.current_state,
	    chat_mode = EXCLUDED.chat_mode,
	    hint_context = EXCLUDED.hint_context
	`, ts.ChatID, ts.SessionID, ts.UpdatedAt, ts.CurrentState, ts.ChatMode, ts.HintContext)

	return err
}

// UpdateSessionState обновляет только состояние и режим (без полного upsert)
func (s *Store) UpdateSessionState(ctx context.Context, chatID int64, state, mode *string) error {
	_, err := s.DB.ExecContext(ctx, `
	UPDATE task_sessions
	SET current_state = $2,
	    chat_mode = $3,
	    updated_at = NOW()
	WHERE chat_id = $1
	`, chatID, state, mode)
	return err
}

// UpdateSessionHintContext обновляет контекст подсказок
func (s *Store) UpdateSessionHintContext(ctx context.Context, chatID int64, hintContext json.RawMessage) error {
	_, err := s.DB.ExecContext(ctx, `
	UPDATE task_sessions
	SET hint_context = $2,
	    updated_at = NOW()
	WHERE chat_id = $1
	`, chatID, hintContext)
	return err
}

// ClearSessionHintContext очищает контекст подсказок
func (s *Store) ClearSessionHintContext(ctx context.Context, chatID int64) error {
	_, err := s.DB.ExecContext(ctx, `
	UPDATE task_sessions
	SET hint_context = NULL,
	    updated_at = NOW()
	WHERE chat_id = $1
	`, chatID)
	return err
}

// UpdateSessionParseContext обновляет контекст ожидания подтверждения парсинга
func (s *Store) UpdateSessionParseContext(ctx context.Context, chatID int64, parseContext json.RawMessage) error {
	_, err := s.DB.ExecContext(ctx, `
	UPDATE task_sessions
	SET parse_context = $2,
	    updated_at = NOW()
	WHERE chat_id = $1
	`, chatID, parseContext)
	return err
}

// ClearSessionParseContext очищает контекст парсинга
func (s *Store) ClearSessionParseContext(ctx context.Context, chatID int64) error {
	_, err := s.DB.ExecContext(ctx, `
	UPDATE task_sessions
	SET parse_context = NULL,
	    updated_at = NOW()
	WHERE chat_id = $1
	`, chatID)
	return err
}

func (s *Store) FindSession(ctx context.Context, chatID int64) (TaskSession, error) {
	const q = `SELECT session_id, updated_at, current_state, chat_mode, hint_context, parse_context
				FROM task_sessions
				WHERE chat_id=$1`
	var sessionID string
	var updatedAt time.Time
	var currentState, chatMode sql.NullString
	var hintContext, parseContext []byte

	err := s.DB.QueryRowContext(ctx, q, chatID).Scan(
		&sessionID, &updatedAt, &currentState, &chatMode, &hintContext, &parseContext,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return TaskSession{}, nil
		}
		return TaskSession{}, err
	}

	ts := TaskSession{
		ChatID:       chatID,
		SessionID:    sessionID,
		UpdatedAt:    updatedAt,
		HintContext:  hintContext,
		ParseContext: parseContext,
	}
	if currentState.Valid {
		ts.CurrentState = &currentState.String
	}
	if chatMode.Valid {
		ts.ChatMode = &chatMode.String
	}

	return ts, nil
}
