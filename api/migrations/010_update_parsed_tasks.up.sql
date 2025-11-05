-- Добавляем колонку session_id
ALTER TABLE parsed_tasks
    ADD COLUMN IF NOT EXISTS session_id text;
