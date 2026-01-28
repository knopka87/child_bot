-- Добавляем task_id из v2
ALTER TABLE parsed_tasks
    ADD COLUMN IF NOT EXISTS task_id TEXT;

-- Добавляем индекс по task_id для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_parsed_tasks_task_id ON parsed_tasks (task_id) WHERE task_id IS NOT NULL;
