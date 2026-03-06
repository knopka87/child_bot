-- Добавляем колонку для сохранения контекста парсинга между редеплоями
ALTER TABLE task_sessions
    ADD COLUMN IF NOT EXISTS parse_context JSONB;

COMMENT ON COLUMN task_sessions.parse_context IS 'Контекст ожидания подтверждения парсинга (parseWait)';
