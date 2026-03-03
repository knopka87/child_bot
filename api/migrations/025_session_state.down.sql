-- Откат миграции
DROP INDEX IF EXISTS idx_task_sessions_updated_at;

ALTER TABLE task_sessions
    DROP COLUMN IF EXISTS current_state,
    DROP COLUMN IF EXISTS chat_mode,
    DROP COLUMN IF EXISTS hint_context;
