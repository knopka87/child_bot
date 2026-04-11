-- Rollback achievements viewed tracking

ALTER TABLE child_profiles
DROP COLUMN IF EXISTS achievements_last_viewed_at;
