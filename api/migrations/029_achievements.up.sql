-- Achievements - справочник достижений
CREATE TABLE IF NOT EXISTS achievements (
    id VARCHAR(100) PRIMARY KEY, -- например: 'streak_7', 'tasks_10', 'villain_1_defeated'

    -- Основная информация
    type VARCHAR(50) NOT NULL, -- 'streak', 'tasks', 'fixes', 'villain_defeater', etc.
    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    icon VARCHAR(200) NOT NULL, -- emoji или URL изображения

    -- Условия разблокировки
    requirement_type VARCHAR(50) NOT NULL, -- 'streak_days', 'tasks_count', 'villain_defeated', etc.
    requirement_value INTEGER NOT NULL, -- целевое значение

    -- Награда
    reward_type VARCHAR(50) NOT NULL, -- 'coins', 'sticker', 'avatar', 'badge'
    reward_id VARCHAR(100),
    reward_name VARCHAR(200),
    reward_amount INTEGER, -- для coins

    -- UI
    shelf_order INTEGER DEFAULT 1, -- полка (1, 2, 3)
    position_in_shelf INTEGER DEFAULT 0, -- позиция на полке (0-3)

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Child Achievements - прогресс достижений для каждого профиля
CREATE TABLE IF NOT EXISTS child_achievements (
    id BIGSERIAL PRIMARY KEY,
    child_profile_id UUID NOT NULL REFERENCES child_profiles(id) ON DELETE CASCADE,
    achievement_id VARCHAR(100) NOT NULL REFERENCES achievements(id) ON DELETE CASCADE,

    -- Прогресс
    current_progress INTEGER DEFAULT 0 CHECK (current_progress >= 0),
    is_unlocked BOOLEAN DEFAULT FALSE,
    is_claimed BOOLEAN DEFAULT FALSE, -- забрана ли награда

    -- Timestamps
    unlocked_at TIMESTAMPTZ,
    claimed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (child_profile_id, achievement_id)
);

-- Индексы
CREATE INDEX IF NOT EXISTS idx_achievements_type
    ON achievements (type);

CREATE INDEX IF NOT EXISTS idx_child_achievements_profile
    ON child_achievements (child_profile_id, is_unlocked);

CREATE INDEX IF NOT EXISTS idx_child_achievements_unlocked
    ON child_achievements (child_profile_id, unlocked_at DESC)
    WHERE is_unlocked = TRUE;

-- Триггер для обновления updated_at
CREATE TRIGGER child_achievements_updated_at
    BEFORE UPDATE ON child_achievements
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Комментарии
COMMENT ON TABLE achievements IS 'Справочник всех достижений';
COMMENT ON TABLE child_achievements IS 'Прогресс достижений для каждого ребенка';
COMMENT ON COLUMN achievements.requirement_type IS 'Тип условия: streak_days, tasks_count, villain_defeated';
COMMENT ON COLUMN achievements.shelf_order IS 'Номер полки для отображения (1-3)';
COMMENT ON COLUMN child_achievements.is_claimed IS 'Забрана ли награда пользователем';

-- Вставим примеры достижений
INSERT INTO achievements (id, type, title, description, icon, requirement_type, requirement_value,
                          reward_type, reward_id, reward_name, reward_amount, shelf_order, position_in_shelf)
VALUES
    ('streak_3', 'streak', 'Три дня подряд', 'Решай задачи 3 дня подряд', '🔥', 'streak_days', 3,
     'coins', 'coins_30', '30 монет', 30, 1, 0),

    ('streak_7', 'streak', 'Неделя успеха', 'Решай задачи 7 дней подряд', '🔥', 'streak_days', 7,
     'coins', 'coins_50', '50 монет', 50, 1, 1),

    ('tasks_10', 'tasks', 'Начинающий', 'Реши 10 задач', '🎯', 'tasks_count', 10,
     'coins', 'coins_20', '20 монет', 20, 1, 2),

    ('tasks_50', 'tasks', 'Знаток', 'Реши 50 задач', '🎯', 'tasks_count', 50,
     'coins', 'coins_100', '100 монет', 100, 1, 3),

    ('villain_1_defeated', 'villain_defeater', 'Победитель Графа', 'Победи Графа Ошибок', '⚔️',
     'villain_defeated', 1, 'sticker', 'sticker_villain_1', 'Стикер Героя', NULL, 2, 0)
ON CONFLICT (id) DO NOTHING;
