-- Villains - справочник злодеев
CREATE TABLE IF NOT EXISTS villains (
    id VARCHAR(100) PRIMARY KEY, -- например: 'count_error', 'baron_confusion'

    -- Основная информация
    name VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    image_url TEXT NOT NULL,

    -- Характеристики
    max_hp INTEGER NOT NULL CHECK (max_hp > 0),
    level INTEGER NOT NULL DEFAULT 1 CHECK (level >= 1),
    damage_per_correct_task INTEGER NOT NULL DEFAULT 5,

    -- Порядок появления
    unlock_order INTEGER NOT NULL DEFAULT 1,

    -- Награды за победу
    reward_coins INTEGER DEFAULT 100,
    reward_achievement_id VARCHAR(100), -- связь с достижением

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Villain Battles - битвы пользователей со злодеями
CREATE TABLE IF NOT EXISTS villain_battles (
    id BIGSERIAL PRIMARY KEY,
    child_profile_id UUID NOT NULL REFERENCES child_profiles(id) ON DELETE CASCADE,
    villain_id VARCHAR(100) NOT NULL REFERENCES villains(id) ON DELETE CASCADE,

    -- Статус битвы
    status VARCHAR(20) NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'defeated', 'abandoned')),

    -- Прогресс
    current_hp INTEGER NOT NULL,
    total_damage_dealt INTEGER DEFAULT 0 CHECK (total_damage_dealt >= 0),
    correct_tasks_count INTEGER DEFAULT 0 CHECK (correct_tasks_count >= 0),

    -- Награды получены?
    rewards_claimed BOOLEAN DEFAULT FALSE,

    -- Timestamps
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    defeated_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (child_profile_id, villain_id, status)
);

-- Damage Events - история нанесения урона
CREATE TABLE IF NOT EXISTS damage_events (
    id BIGSERIAL PRIMARY KEY,
    battle_id BIGINT NOT NULL REFERENCES villain_battles(id) ON DELETE CASCADE,
    attempt_id UUID REFERENCES attempts(id) ON DELETE SET NULL,

    -- Урон
    damage INTEGER NOT NULL CHECK (damage >= 0),
    task_type VARCHAR(20) NOT NULL CHECK (task_type IN ('help', 'check')),

    -- Timestamp
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Индексы
CREATE INDEX IF NOT EXISTS idx_villains_unlock_order
    ON villains (unlock_order);

CREATE INDEX IF NOT EXISTS idx_villain_battles_profile_status
    ON villain_battles (child_profile_id, status);

CREATE INDEX IF NOT EXISTS idx_villain_battles_active
    ON villain_battles (child_profile_id)
    WHERE status = 'active';

CREATE INDEX IF NOT EXISTS idx_damage_events_battle
    ON damage_events (battle_id, created_at DESC);

-- Триггер для обновления updated_at
CREATE TRIGGER villain_battles_updated_at
    BEFORE UPDATE ON villain_battles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Комментарии
COMMENT ON TABLE villains IS 'Справочник злодеев';
COMMENT ON TABLE villain_battles IS 'Битвы пользователей со злодеями';
COMMENT ON TABLE damage_events IS 'История нанесения урона злодеям';
COMMENT ON COLUMN villains.damage_per_correct_task IS 'Урон за одну правильно решенную задачу';
COMMENT ON COLUMN villain_battles.status IS 'Статус: active, defeated, abandoned';

-- Вставим примеры злодеев
INSERT INTO villains (id, name, description, image_url, max_hp, level, damage_per_correct_task,
                      unlock_order, reward_coins, reward_achievement_id)
VALUES
    ('count_error', 'Граф Ошибок', 'Злодей, который распространяет ошибки в задачах',
     '/assets/villains/count_error.png', 100, 1, 5, 1, 100, 'villain_1_defeated'),

    ('baron_confusion', 'Барон Путаница', 'Мастер запутывания логических цепочек',
     '/assets/villains/baron_confusion.png', 150, 2, 5, 2, 150, NULL),

    ('duchess_distraction', 'Герцогиня Отвлечения', 'Заставляет терять концентрацию',
     '/assets/villains/duchess_distraction.png', 200, 3, 5, 3, 200, NULL)
ON CONFLICT (id) DO NOTHING;
