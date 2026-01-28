-- Откат миграции v2

DROP INDEX IF EXISTS idx_parsed_tasks_task_id;

-- Удаляем новые колонки
ALTER TABLE parsed_tasks
    DROP COLUMN IF EXISTS task_id;