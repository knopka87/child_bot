-- metrics events table (flexible via JSONB details)
CREATE TABLE IF NOT EXISTS metrics_events (
                                              id             BIGSERIAL PRIMARY KEY,
                                              created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),

                                              stage          TEXT NOT NULL,       -- detect|parse|hint|normalize|check|analogue
                                              provider       TEXT,                -- gemini|gpt|yandex|...

                                              ok             BOOLEAN NOT NULL DEFAULT true,
                                              error          TEXT,
                                              http_code      INTEGER,
                                              duration_ms    BIGINT,

                                              chat_id        BIGINT,
                                              user_id_anon   TEXT,
                                              task_id        TEXT,
                                              correlation_id TEXT,
                                              request_id     TEXT,

                                              details        JSONB NOT NULL DEFAULT '{}'::jsonb
);

-- Helpful indexes
-- CREATE INDEX IF NOT EXISTS metrics_events_created_at_idx ON metrics_events (created_at);
CREATE INDEX IF NOT EXISTS metrics_events_stage_idx      ON metrics_events (stage);
CREATE INDEX IF NOT EXISTS metrics_events_provider_idx   ON metrics_events (provider);
CREATE INDEX IF NOT EXISTS metrics_events_chat_idx       ON metrics_events (chat_id);
-- CREATE INDEX IF NOT EXISTS metrics_events_ok_idx         ON metrics_events (ok);
CREATE INDEX IF NOT EXISTS metrics_events_details_gin    ON metrics_events USING GIN (details);