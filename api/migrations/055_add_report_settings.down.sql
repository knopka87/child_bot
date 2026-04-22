-- Удаление таблицы настроек отчётов
DROP TRIGGER IF EXISTS report_settings_updated_at ON report_settings;
DROP INDEX IF EXISTS idx_report_settings_email;
DROP TABLE IF EXISTS report_settings;
