-- Drop obsolete non-unique index on (engine, image_hash)
DROP INDEX IF EXISTS idx_parsed_tasks_engine_img;

-- Add unique index on session_id (required for upsert / lookups)
CREATE UNIQUE INDEX IF NOT EXISTS idx_parsed_tasks_session_id
  ON parsed_tasks (session_id);
