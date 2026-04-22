-- Parent Consent History table for audit trail
CREATE TABLE IF NOT EXISTS parent_consent_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Связь с parent_consents
    consent_id UUID NOT NULL,
    parent_user_id VARCHAR(255) NOT NULL,
    platform_id VARCHAR(20) NOT NULL,

    -- Тип изменения
    action VARCHAR(20) NOT NULL, -- 'created', 'updated', 'revoked'

    -- Snapshot согласий на момент изменения
    privacy_policy_version VARCHAR(20) NOT NULL,
    privacy_policy_accepted BOOLEAN NOT NULL,

    terms_version VARCHAR(20) NOT NULL,
    terms_accepted BOOLEAN NOT NULL,

    adult_consent BOOLEAN NOT NULL,

    -- Что изменилось (для action='updated')
    changed_fields TEXT[], -- ['privacy_policy_version', 'terms_version']
    previous_values JSONB, -- {"privacy_policy_version": "1.0", "terms_version": "1.0"}

    -- Audit metadata
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Индексы для быстрого поиска истории
CREATE INDEX IF NOT EXISTS idx_parent_consent_history_consent
    ON parent_consent_history (consent_id);

CREATE INDEX IF NOT EXISTS idx_parent_consent_history_parent
    ON parent_consent_history (platform_id, parent_user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_parent_consent_history_action
    ON parent_consent_history (action, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_parent_consent_history_created
    ON parent_consent_history (created_at DESC);

-- Комментарии
COMMENT ON TABLE parent_consent_history IS 'История изменений согласий родителей (audit trail)';
COMMENT ON COLUMN parent_consent_history.consent_id IS 'ID записи в parent_consents';
COMMENT ON COLUMN parent_consent_history.action IS 'Тип действия: created, updated, revoked';
COMMENT ON COLUMN parent_consent_history.changed_fields IS 'Массив полей которые изменились';
COMMENT ON COLUMN parent_consent_history.previous_values IS 'Предыдущие значения измененных полей (JSON)';
