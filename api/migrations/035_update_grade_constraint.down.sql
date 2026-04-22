-- Rollback: restore grade constraint to 1-11

-- Drop new constraint
ALTER TABLE child_profiles
DROP CONSTRAINT IF EXISTS child_profiles_grade_check;

-- Restore old constraint for grades 1-11
ALTER TABLE child_profiles
ADD CONSTRAINT child_profiles_grade_check CHECK (grade >= 1 AND grade <= 11);

-- Restore comment
COMMENT ON COLUMN child_profiles.grade IS 'Класс ученика (1-11)';
