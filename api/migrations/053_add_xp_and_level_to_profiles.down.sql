-- Откатываем XP и уровень
ALTER TABLE child_profiles
DROP COLUMN IF EXISTS xp_total,
DROP COLUMN IF EXISTS level;

DROP INDEX IF EXISTS idx_child_profiles_level;
