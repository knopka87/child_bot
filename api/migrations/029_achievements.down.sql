DROP TRIGGER IF EXISTS child_achievements_updated_at ON child_achievements;
DROP INDEX IF EXISTS idx_child_achievements_unlocked;
DROP INDEX IF EXISTS idx_child_achievements_profile;
DROP INDEX IF EXISTS idx_achievements_type;
DROP TABLE IF EXISTS child_achievements;
DROP TABLE IF EXISTS achievements;
