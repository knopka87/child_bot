-- Payments - транзакции платежей
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Связь с подпиской
    subscription_id BIGINT REFERENCES subscriptions(id) ON DELETE SET NULL,
    child_profile_id UUID NOT NULL REFERENCES child_profiles(id) ON DELETE CASCADE,
    plan_id VARCHAR(100) NOT NULL REFERENCES subscription_plans(id),

    -- Информация о платеже
    amount_cents INTEGER NOT NULL CHECK (amount_cents >= 0),
    currency VARCHAR(10) NOT NULL DEFAULT 'RUB',

    -- VK Pay данные
    vk_order_id VARCHAR(255) UNIQUE, -- ID заказа в VK Pay
    vk_transaction_id VARCHAR(255), -- ID транзакции от VK
    vk_user_id VARCHAR(255), -- VK User ID плательщика

    -- Статус платежа
    status VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'refunded', 'cancelled')),

    -- Метаданные
    payment_method VARCHAR(50) DEFAULT 'vk_pay',
    description TEXT,

    -- IP и User Agent для безопасности
    ip_address VARCHAR(45),
    user_agent TEXT,

    -- Дополнительные данные
    metadata JSONB DEFAULT '{}'::jsonb,

    -- Даты
    paid_at TIMESTAMPTZ,
    refunded_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ, -- Для pending платежей

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Payment Events - события платежей (для аудита)
CREATE TABLE IF NOT EXISTS payment_events (
    id BIGSERIAL PRIMARY KEY,
    payment_id UUID NOT NULL REFERENCES payments(id) ON DELETE CASCADE,

    -- Тип события
    event_type VARCHAR(50) NOT NULL, -- created, processing, completed, failed, refunded, webhook_received

    -- Данные события
    old_status VARCHAR(20),
    new_status VARCHAR(20),

    -- VK webhook данные
    vk_event_type VARCHAR(100),
    vk_event_data JSONB,

    -- Ошибки
    error_code VARCHAR(100),
    error_message TEXT,

    -- Timestamp
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Индексы для payments
CREATE INDEX IF NOT EXISTS idx_payments_subscription
    ON payments (subscription_id);

CREATE INDEX IF NOT EXISTS idx_payments_profile
    ON payments (child_profile_id, status);

CREATE INDEX IF NOT EXISTS idx_payments_vk_order
    ON payments (vk_order_id)
    WHERE vk_order_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_payments_status
    ON payments (status, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_payments_expires
    ON payments (expires_at)
    WHERE status = 'pending' AND expires_at IS NOT NULL;

-- Индексы для payment_events
CREATE INDEX IF NOT EXISTS idx_payment_events_payment
    ON payment_events (payment_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_payment_events_type
    ON payment_events (event_type, created_at DESC);

-- Триггер для обновления updated_at
CREATE TRIGGER payments_updated_at
    BEFORE UPDATE ON payments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Комментарии
COMMENT ON TABLE payments IS 'Транзакции платежей';
COMMENT ON TABLE payment_events IS 'События и аудит платежей';
COMMENT ON COLUMN payments.vk_order_id IS 'ID заказа в VK Pay (уникальный)';
COMMENT ON COLUMN payments.vk_transaction_id IS 'ID транзакции от VK Pay';
COMMENT ON COLUMN payments.status IS 'Статус: pending, processing, completed, failed, refunded, cancelled';
COMMENT ON COLUMN payments.expires_at IS 'Время истечения для pending платежей (обычно +30 минут)';
