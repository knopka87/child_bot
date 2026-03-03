-- Добавляем колонки для сохранения состояния сессии между редеплоями
ALTER TABLE task_sessions
    ADD COLUMN IF NOT EXISTS current_state VARCHAR(50),
    ADD COLUMN IF NOT EXISTS chat_mode VARCHAR(50),
    ADD COLUMN IF NOT EXISTS hint_context JSONB;

-- Индекс для быстрого поиска активных сессий
CREATE INDEX IF NOT EXISTS idx_task_sessions_updated_at
    ON task_sessions (updated_at DESC)
    WHERE current_state IS NOT NULL;

COMMENT ON COLUMN task_sessions.current_state IS 'Текущее состояние машины состояний (awaiting_task, hints, etc.)';
COMMENT ON COLUMN task_sessions.chat_mode IS 'Режим чата (await_solution, await_new_task, etc.)';
COMMENT ON COLUMN task_sessions.hint_context IS 'Контекст для продолжения подсказок (Parse, Detect результаты)';
