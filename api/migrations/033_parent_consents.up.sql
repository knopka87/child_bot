-- Parent Consents table for storing parent consent records
CREATE TABLE IF NOT EXISTS parent_consents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Родитель (платформенный ID)
    parent_user_id VARCHAR(255) NOT NULL,
    platform_id VARCHAR(20) NOT NULL, -- 'vk', 'telegram', 'max', 'web'

    -- Согласия
    privacy_policy_version VARCHAR(20) NOT NULL,
    privacy_policy_accepted BOOLEAN NOT NULL DEFAULT TRUE,
    privacy_policy_accepted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    terms_version VARCHAR(20) NOT NULL,
    terms_accepted BOOLEAN NOT NULL DEFAULT TRUE,
    terms_accepted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    adult_consent BOOLEAN NOT NULL DEFAULT FALSE, -- подтверждение что родитель совершеннолетний
    adult_consent_at TIMESTAMPTZ,

    -- IP для аудита
    ip_address VARCHAR(45), -- IPv4 или IPv6
    user_agent TEXT,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Уникальность по родителю и платформе (можно обновлять при изменении версий)
    UNIQUE (platform_id, parent_user_id)
);

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_parent_consents_parent
    ON parent_consents (platform_id, parent_user_id);

CREATE INDEX IF NOT EXISTS idx_parent_consents_created
    ON parent_consents (created_at DESC);

-- Триггер для обновления updated_at
CREATE TRIGGER parent_consents_updated_at
    BEFORE UPDATE ON parent_consents
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Комментарии
COMMENT ON TABLE parent_consents IS 'Согласия родителей на обработку данных и условия использования';
COMMENT ON COLUMN parent_consents.parent_user_id IS 'ID родителя на платформе (platform_user_id)';
COMMENT ON COLUMN parent_consents.platform_id IS 'Платформа: vk, telegram, max, web';
COMMENT ON COLUMN parent_consents.adult_consent IS 'Подтверждение совершеннолетия родителя';
COMMENT ON COLUMN parent_consents.privacy_policy_version IS 'Версия политики конфиденциальности (например: 1.0, 2.1)';
COMMENT ON COLUMN parent_consents.terms_version IS 'Версия условий использования (например: 1.0, 2.1)';
