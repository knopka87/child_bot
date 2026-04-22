-- Track when user last viewed achievements page

ALTER TABLE child_profiles
ADD COLUMN achievements_last_viewed_at TIMESTAMPTZ DEFAULT NOW();

COMMENT ON COLUMN child_profiles.achievements_last_viewed_at IS 'Дата последнего просмотра страницы достижений (для badge новых достижений)';
