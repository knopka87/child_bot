DROP TRIGGER IF EXISTS child_profiles_updated_at ON child_profiles;
DROP FUNCTION IF EXISTS update_child_profiles_updated_at();
DROP INDEX IF EXISTS idx_child_profiles_created;
DROP INDEX IF EXISTS idx_child_profiles_activity;
DROP INDEX IF EXISTS idx_child_profiles_platform;
DROP TABLE IF EXISTS child_profiles;
