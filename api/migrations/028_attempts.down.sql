DROP TRIGGER IF EXISTS attempts_updated_at ON attempts;
DROP INDEX IF EXISTS idx_attempts_unfinished;
DROP INDEX IF EXISTS idx_attempts_type;
DROP INDEX IF EXISTS idx_attempts_status;
DROP INDEX IF EXISTS idx_attempts_child_profile;
DROP TABLE IF EXISTS attempts;
