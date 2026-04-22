-- Referrals - реферальная программа
CREATE TABLE IF NOT EXISTS referrals (
    id BIGSERIAL PRIMARY KEY,
    referrer_id UUID NOT NULL REFERENCES child_profiles(id) ON DELETE CASCADE, -- кто пригласил
    referred_id UUID NOT NULL REFERENCES child_profiles(id) ON DELETE CASCADE, -- кого пригласили

    -- Статус
    is_active BOOLEAN DEFAULT FALSE, -- активировался ли приглашенный (выполнил первое действие)

    -- Награды
    reward_coins INTEGER DEFAULT 50, -- награда за приглашение
    reward_claimed BOOLEAN DEFAULT FALSE, -- забрана ли награда

    -- Timestamps
    invited_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    activated_at TIMESTAMPTZ, -- когда приглашенный стал активным
    reward_claimed_at TIMESTAMPTZ,

    UNIQUE (referrer_id, referred_id)
);

-- Referral Codes - реферальные коды (опционально, если нужны уникальные коды)
CREATE TABLE IF NOT EXISTS referral_codes (
    id BIGSERIAL PRIMARY KEY,
    child_profile_id UUID NOT NULL REFERENCES child_profiles(id) ON DELETE CASCADE,

    -- Код
    code VARCHAR(20) NOT NULL UNIQUE, -- например: 'ABCD1234'

    -- Статистика
    uses_count INTEGER DEFAULT 0 CHECK (uses_count >= 0),
    max_uses INTEGER, -- NULL = unlimited

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ,

    UNIQUE (child_profile_id)
);

-- Referral Milestones - награды за количество приглашений
CREATE TABLE IF NOT EXISTS referral_milestones (
    id BIGSERIAL PRIMARY KEY,

    -- Условие
    friends_count INTEGER NOT NULL UNIQUE CHECK (friends_count > 0),

    -- Награда
    reward_coins INTEGER NOT NULL CHECK (reward_coins > 0),
    description TEXT NOT NULL,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Child Referral Milestones - прогресс milestone для каждого пользователя
CREATE TABLE IF NOT EXISTS child_referral_milestones (
    id BIGSERIAL PRIMARY KEY,
    child_profile_id UUID NOT NULL REFERENCES child_profiles(id) ON DELETE CASCADE,
    milestone_id BIGINT NOT NULL REFERENCES referral_milestones(id) ON DELETE CASCADE,

    -- Статус
    is_claimed BOOLEAN DEFAULT FALSE,
    claimed_at TIMESTAMPTZ,

    UNIQUE (child_profile_id, milestone_id)
);

-- Индексы
CREATE INDEX IF NOT EXISTS idx_referrals_referrer
    ON referrals (referrer_id, is_active);

CREATE INDEX IF NOT EXISTS idx_referrals_referred
    ON referrals (referred_id);

CREATE INDEX IF NOT EXISTS idx_referral_codes_code
    ON referral_codes (code);

CREATE INDEX IF NOT EXISTS idx_child_referral_milestones_profile
    ON child_referral_milestones (child_profile_id, is_claimed);

-- Комментарии
COMMENT ON TABLE referrals IS 'Реферальные связи между пользователями';
COMMENT ON TABLE referral_codes IS 'Уникальные реферальные коды пользователей';
COMMENT ON TABLE referral_milestones IS 'Награды за количество приглашенных друзей';
COMMENT ON TABLE child_referral_milestones IS 'Прогресс milestone для каждого пользователя';
COMMENT ON COLUMN referrals.is_active IS 'Активировался ли приглашенный пользователь';
COMMENT ON COLUMN referral_codes.code IS 'Уникальный код приглашения';

-- Вставим примеры milestone
INSERT INTO referral_milestones (friends_count, reward_coins, description)
VALUES
    (1, 50, 'Пригласи 1 друга'),
    (3, 100, 'Пригласи 3 друзей'),
    (5, 200, 'Пригласи 5 друзей'),
    (10, 500, 'Пригласи 10 друзей')
ON CONFLICT (friends_count) DO NOTHING;

-- Функция для генерации уникального реферального кода
CREATE OR REPLACE FUNCTION generate_referral_code()
RETURNS TEXT AS $$
DECLARE
    chars TEXT := 'ABCDEFGHJKLMNPQRSTUVWXYZ23456789'; -- без похожих символов
    result TEXT := '';
    i INTEGER;
BEGIN
    FOR i IN 1..8 LOOP
        result := result || substr(chars, floor(random() * length(chars) + 1)::int, 1);
    END LOOP;
    RETURN result;
END;
$$ LANGUAGE plpgsql;

-- Триггер для автоматического создания реферального кода при создании профиля
CREATE OR REPLACE FUNCTION create_referral_code_for_profile()
RETURNS TRIGGER AS $$
DECLARE
    new_code TEXT;
    max_attempts INTEGER := 10;
    attempt INTEGER := 0;
BEGIN
    LOOP
        new_code := generate_referral_code();
        attempt := attempt + 1;

        BEGIN
            INSERT INTO referral_codes (child_profile_id, code)
            VALUES (NEW.id, new_code);
            EXIT; -- успешно вставили, выходим
        EXCEPTION WHEN unique_violation THEN
            IF attempt >= max_attempts THEN
                RAISE EXCEPTION 'Failed to generate unique referral code after % attempts', max_attempts;
            END IF;
        END;
    END LOOP;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER child_profile_create_referral_code
    AFTER INSERT ON child_profiles
    FOR EACH ROW
    EXECUTE FUNCTION create_referral_code_for_profile();
