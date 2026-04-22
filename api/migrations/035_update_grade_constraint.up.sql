-- Update grade constraint to allow only grades 1-4
-- Начальная школа: 1-4 класс

-- Drop old constraint
ALTER TABLE child_profiles
DROP CONSTRAINT IF EXISTS child_profiles_grade_check;

-- Add new constraint for grades 1-4
ALTER TABLE child_profiles
ADD CONSTRAINT child_profiles_grade_check CHECK (grade >= 1 AND grade <= 4);

-- Comment
COMMENT ON COLUMN child_profiles.grade IS 'Класс ученика (начальная школа: 1-4)';
