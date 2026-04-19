-- Добавляем XP и уровень в профиль пользователя
ALTER TABLE child_profiles
ADD COLUMN xp_total INTEGER NOT NULL DEFAULT 0,
ADD COLUMN level INTEGER NOT NULL DEFAULT 1;

-- Индекс для быстрого поиска по уровню
CREATE INDEX idx_child_profiles_level ON child_profiles(level DESC);
