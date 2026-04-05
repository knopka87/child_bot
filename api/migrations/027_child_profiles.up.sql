-- Child Profiles table for REST API
-- Основная таблица профилей детей (не привязана к Telegram chat_id)
CREATE TABLE IF NOT EXISTS child_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Основная информация
    display_name VARCHAR(100) NOT NULL,
    avatar_id VARCHAR(50),
    grade INTEGER NOT NULL CHECK (grade >= 1 AND grade <= 11),

    -- Контактная информация (опционально)
    email VARCHAR(255),
    email_verified BOOLEAN DEFAULT FALSE,

    -- Платформа регистрации
    platform_id VARCHAR(20) NOT NULL, -- 'vk', 'telegram', 'max', 'web'
    platform_user_id VARCHAR(255), -- ID пользователя на платформе (опционально)

    -- Gamification
    level INTEGER DEFAULT 1 CHECK (level >= 1),
    experience_points INTEGER DEFAULT 0 CHECK (experience_points >= 0),
    coins_balance INTEGER DEFAULT 0 CHECK (coins_balance >= 0),

    -- Статистика
    tasks_solved_total INTEGER DEFAULT 0 CHECK (tasks_solved_total >= 0),
    tasks_solved_correct INTEGER DEFAULT 0 CHECK (tasks_solved_correct >= 0),
    hints_used_total INTEGER DEFAULT 0 CHECK (hints_used_total >= 0),
    streak_days INTEGER DEFAULT 0 CHECK (streak_days >= 0),
    last_activity_at TIMESTAMPTZ,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Unique constraint для платформ (один профиль на платформу+user)
    UNIQUE (platform_id, platform_user_id)
);

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_child_profiles_platform
    ON child_profiles (platform_id, platform_user_id);

CREATE INDEX IF NOT EXISTS idx_child_profiles_activity
    ON child_profiles (last_activity_at DESC);

CREATE INDEX IF NOT EXISTS idx_child_profiles_created
    ON child_profiles (created_at DESC);

-- Триггер для обновления updated_at
CREATE OR REPLACE FUNCTION update_child_profiles_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER child_profiles_updated_at
    BEFORE UPDATE ON child_profiles
    FOR EACH ROW
    EXECUTE FUNCTION update_child_profiles_updated_at();

-- Комментарии
COMMENT ON TABLE child_profiles IS 'Профили детей для REST API (независимо от Telegram)';
COMMENT ON COLUMN child_profiles.platform_id IS 'Платформа: vk, telegram, max, web';
COMMENT ON COLUMN child_profiles.platform_user_id IS 'ID пользователя на платформе (опционально)';
COMMENT ON COLUMN child_profiles.level IS 'Уровень пользователя в gamification';
COMMENT ON COLUMN child_profiles.experience_points IS 'Очки опыта для прогресса уровня';
