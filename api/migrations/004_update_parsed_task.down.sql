-- Revert new index
DROP INDEX IF EXISTS idx_parsed_tasks_engine_img;

-- Remove acceptance fields
ALTER TABLE parsed_tasks
    DROP COLUMN IF EXISTS accepted_reason,
    DROP COLUMN IF EXISTS accepted;

-- Restore old index by chat/time
CREATE INDEX IF NOT EXISTS idx_parsed_tasks_chat_time ON parsed_tasks (chat_id, created_at DESC);