-- Таблица для настроек отчётов родителям
CREATE TABLE IF NOT EXISTS report_settings (
    child_profile_id UUID PRIMARY KEY REFERENCES child_profiles(id) ON DELETE CASCADE,
    parent_email VARCHAR(255),
    weekly_report_enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Индекс для поиска по email
CREATE INDEX idx_report_settings_email ON report_settings(parent_email);

-- Триггер для обновления updated_at
CREATE TRIGGER report_settings_updated_at 
    BEFORE UPDATE ON report_settings 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
