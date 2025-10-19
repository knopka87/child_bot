CREATE TABLE IF NOT EXISTS task_sessions (
                                             chat_id BIGINT PRIMARY KEY,
                                             task_session_id TEXT NOT NULL,
                                             updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
