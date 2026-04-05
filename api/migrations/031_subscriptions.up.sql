-- Subscription Plans - планы подписок
CREATE TABLE IF NOT EXISTS subscription_plans (
    id VARCHAR(100) PRIMARY KEY, -- например: 'monthly', 'yearly'

    -- Основная информация
    name VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,

    -- Цена
    price_cents INTEGER NOT NULL CHECK (price_cents >= 0), -- цена в копейках
    currency VARCHAR(10) NOT NULL DEFAULT 'RUB',

    -- Параметры
    duration_days INTEGER NOT NULL CHECK (duration_days > 0), -- длительность в днях
    trial_days INTEGER DEFAULT 0 CHECK (trial_days >= 0),

    -- Скидка
    discount_percent INTEGER DEFAULT 0 CHECK (discount_percent >= 0 AND discount_percent <= 100),

    -- UI
    is_popular BOOLEAN DEFAULT FALSE,
    display_order INTEGER DEFAULT 0,

    -- Статус
    is_active BOOLEAN DEFAULT TRUE,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Subscriptions - подписки пользователей
CREATE TABLE IF NOT EXISTS subscriptions (
    id BIGSERIAL PRIMARY KEY,
    child_profile_id UUID NOT NULL REFERENCES child_profiles(id) ON DELETE CASCADE,
    plan_id VARCHAR(100) NOT NULL REFERENCES subscription_plans(id),

    -- Статус
    status VARCHAR(20) NOT NULL DEFAULT 'trial'
        CHECK (status IN ('trial', 'active', 'expired', 'cancelled')),

    -- Даты
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    trial_ends_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ NOT NULL,
    cancelled_at TIMESTAMPTZ,

    -- Автопродление
    auto_renew BOOLEAN DEFAULT TRUE,

    -- Payment info (опционально, для интеграции с платежными системами)
    payment_provider VARCHAR(50), -- 'yookassa', 'stripe', etc.
    payment_external_id VARCHAR(255), -- ID в платежной системе

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Только одна активная подписка на профиль (partial unique index)
CREATE UNIQUE INDEX IF NOT EXISTS idx_subscriptions_active_unique
    ON subscriptions (child_profile_id)
    WHERE status IN ('trial', 'active');

-- Индексы
CREATE INDEX IF NOT EXISTS idx_subscription_plans_active
    ON subscription_plans (is_active, display_order);

CREATE INDEX IF NOT EXISTS idx_subscriptions_profile
    ON subscriptions (child_profile_id, status);

CREATE INDEX IF NOT EXISTS idx_subscriptions_expires
    ON subscriptions (expires_at)
    WHERE status IN ('trial', 'active');

-- Триггер для обновления updated_at
CREATE TRIGGER subscriptions_updated_at
    BEFORE UPDATE ON subscriptions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Комментарии
COMMENT ON TABLE subscription_plans IS 'Планы подписок';
COMMENT ON TABLE subscriptions IS 'Подписки пользователей';
COMMENT ON COLUMN subscription_plans.price_cents IS 'Цена в копейках (499 руб = 49900 копеек)';
COMMENT ON COLUMN subscription_plans.duration_days IS 'Длительность подписки в днях';
COMMENT ON COLUMN subscriptions.status IS 'Статус: trial, active, expired, cancelled';

-- Вставим примеры планов
INSERT INTO subscription_plans (id, name, description, price_cents, currency, duration_days,
                                trial_days, discount_percent, is_popular, display_order)
VALUES
    ('monthly', 'Месячная подписка', 'Полный доступ ко всем функциям на 1 месяц',
     49900, 'RUB', 30, 7, 0, TRUE, 1),

    ('yearly', 'Годовая подписка', 'Выгодная подписка на целый год - экономия 33%',
     399900, 'RUB', 365, 14, 33, FALSE, 2)
ON CONFLICT (id) DO NOTHING;
