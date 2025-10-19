-- 001_create_timeline.sql
CREATE TABLE IF NOT EXISTS timeline_events (
                                               id               BIGSERIAL PRIMARY KEY,
                                               chat_id          BIGINT NOT NULL,
                                               task_session_id  TEXT,                -- UUID как строка
                                               direction        TEXT NOT NULL,       -- 'in' | 'out' | 'api'
                                               event_type       TEXT NOT NULL,       -- tg_text|tg_photo|tg_callback|tg_out_text|api_detect|api_parse|api_hint|api_normalize|api_check|api_analogue|...
                                               provider         TEXT,                -- gemini|gpt для api_*
                                               ok               BOOLEAN,
                                               latency_ms       BIGINT,
                                               tg_message_id    BIGINT,              -- для входящих/исходящих сообщений
                                               text             TEXT,
                                               input_payload    JSONB,               -- вход (redacted)
                                               output_payload   JSONB,               -- выход (redacted)
                                               error            TEXT,
                                               created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_timeline_chat_time   ON timeline_events (chat_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_timeline_session_time ON timeline_events (task_session_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_timeline_event_type   ON timeline_events (chat_id, event_type);