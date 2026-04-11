-- Auto-claim rewards when achievement is unlocked (Variant A)

-- Функция для автоматического начисления награды при разблокировке достижения
CREATE OR REPLACE FUNCTION auto_claim_achievement_reward()
RETURNS TRIGGER AS $$
DECLARE
    achievement_reward_type VARCHAR(50);
    achievement_reward_amount INTEGER;
BEGIN
    -- Проверяем что достижение только что разблокировалось
    IF NEW.is_unlocked = TRUE AND (OLD.is_unlocked IS NULL OR OLD.is_unlocked = FALSE) THEN
        -- Получаем информацию о награде
        SELECT reward_type, reward_amount
        INTO achievement_reward_type, achievement_reward_amount
        FROM achievements
        WHERE id = NEW.achievement_id;

        -- Если награда - монеты, начисляем их
        IF achievement_reward_type = 'coins' AND achievement_reward_amount IS NOT NULL THEN
            UPDATE child_profiles
            SET coins_balance = coins_balance + achievement_reward_amount,
                updated_at = NOW()
            WHERE id = NEW.child_profile_id;
        END IF;

        -- Автоматически помечаем награду как полученную
        NEW.is_claimed := TRUE;
        NEW.claimed_at := NOW();
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Создаём триггер
CREATE TRIGGER trigger_auto_claim_achievement
    BEFORE UPDATE ON child_achievements
    FOR EACH ROW
    EXECUTE FUNCTION auto_claim_achievement_reward();

-- Также устанавливаем is_claimed = TRUE при INSERT если уже unlocked
CREATE OR REPLACE FUNCTION auto_claim_achievement_on_insert()
RETURNS TRIGGER AS $$
DECLARE
    achievement_reward_type VARCHAR(50);
    achievement_reward_amount INTEGER;
BEGIN
    -- Если создаётся уже разблокированное достижение
    IF NEW.is_unlocked = TRUE THEN
        -- Получаем информацию о награде
        SELECT reward_type, reward_amount
        INTO achievement_reward_type, achievement_reward_amount
        FROM achievements
        WHERE id = NEW.achievement_id;

        -- Если награда - монеты, начисляем их
        IF achievement_reward_type = 'coins' AND achievement_reward_amount IS NOT NULL THEN
            UPDATE child_profiles
            SET coins_balance = coins_balance + achievement_reward_amount,
                updated_at = NOW()
            WHERE id = NEW.child_profile_id;
        END IF;

        -- Автоматически помечаем награду как полученную
        NEW.is_claimed := TRUE;
        NEW.claimed_at := NOW();
        NEW.unlocked_at := NOW();
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_auto_claim_achievement_insert
    BEFORE INSERT ON child_achievements
    FOR EACH ROW
    EXECUTE FUNCTION auto_claim_achievement_on_insert();

-- Аналогично для referrals - автоматически начисляем награду при активации
CREATE OR REPLACE FUNCTION auto_claim_referral_reward()
RETURNS TRIGGER AS $$
BEGIN
    -- Проверяем что реферал только что активировался
    IF NEW.is_active = TRUE AND (OLD.is_active IS NULL OR OLD.is_active = FALSE) THEN
        -- Начисляем монеты рефереру
        IF NEW.reward_coins IS NOT NULL THEN
            UPDATE child_profiles
            SET coins_balance = coins_balance + NEW.reward_coins,
                updated_at = NOW()
            WHERE id = NEW.referrer_id;
        END IF;

        -- Автоматически помечаем награду как полученную
        NEW.reward_claimed := TRUE;
        NEW.reward_claimed_at := NOW();
        NEW.activated_at := NOW();
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_auto_claim_referral
    BEFORE UPDATE ON referrals
    FOR EACH ROW
    EXECUTE FUNCTION auto_claim_referral_reward();

-- Аналогично для child_referral_milestones
CREATE OR REPLACE FUNCTION auto_claim_milestone_reward()
RETURNS TRIGGER AS $$
DECLARE
    milestone_reward_coins INTEGER;
BEGIN
    -- При INSERT если is_claimed = FALSE, начисляем награду автоматически
    IF TG_OP = 'INSERT' THEN
        -- Получаем награду за milestone
        SELECT reward_coins
        INTO milestone_reward_coins
        FROM referral_milestones
        WHERE id = NEW.milestone_id;

        -- Начисляем монеты
        IF milestone_reward_coins IS NOT NULL THEN
            UPDATE child_profiles
            SET coins_balance = coins_balance + milestone_reward_coins,
                updated_at = NOW()
            WHERE id = NEW.child_profile_id;
        END IF;

        -- Автоматически помечаем как полученное
        NEW.is_claimed := TRUE;
        NEW.claimed_at := NOW();
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_auto_claim_milestone
    BEFORE INSERT ON child_referral_milestones
    FOR EACH ROW
    EXECUTE FUNCTION auto_claim_milestone_reward();

-- Комментарии
COMMENT ON FUNCTION auto_claim_achievement_reward() IS 'Автоматически начисляет награду при разблокировке достижения (Вариант А)';
COMMENT ON FUNCTION auto_claim_referral_reward() IS 'Автоматически начисляет награду за активацию реферала';
COMMENT ON FUNCTION auto_claim_milestone_reward() IS 'Автоматически начисляет награду за referral milestone';
