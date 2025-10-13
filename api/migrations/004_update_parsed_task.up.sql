-- Add acceptance fields
ALTER TABLE parsed_tasks
    ADD COLUMN IF NOT EXISTS accepted BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS accepted_reason TEXT;

-- Drop old index by chat/time
DROP INDEX IF EXISTS idx_parsed_tasks_chat_time;

-- New index by engine + image
CREATE INDEX IF NOT EXISTS idx_parsed_tasks_engine_img ON parsed_tasks (engine, image_hash DESC);
