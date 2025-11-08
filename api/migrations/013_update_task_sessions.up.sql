-- Drop unique constraint first
ALTER TABLE parsed_tasks
    DROP CONSTRAINT IF EXISTS parsed_tasks_image_hash_engine_key;

-- Now we can safely drop the index if it still exists
DROP INDEX IF EXISTS parsed_tasks_image_hash_engine_idx;

CREATE INDEX IF NOT EXISTS ix_parsed_tasks_chat_updated
    ON parsed_tasks(chat_id, updated_at DESC);

ALTER TABLE parsed_tasks
    ALTER COLUMN session_id SET NOT NULL;