-- Sync email_verified between email_verifications and child_profiles

-- Функция для синхронизации email_verified
CREATE OR REPLACE FUNCTION sync_email_verified()
RETURNS TRIGGER AS $$
BEGIN
    -- Когда email становится verified, обновляем child_profiles
    IF NEW.is_verified = TRUE AND (OLD.is_verified IS NULL OR OLD.is_verified = FALSE) THEN
        UPDATE child_profiles
        SET email_verified = TRUE,
            updated_at = NOW()
        WHERE email = NEW.email;

        RAISE NOTICE 'Synced email_verified for email: %', NEW.email;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Триггер на UPDATE email_verifications
CREATE TRIGGER trigger_sync_email_verified
    AFTER UPDATE ON email_verifications
    FOR EACH ROW
    EXECUTE FUNCTION sync_email_verified();

-- Комментарий
COMMENT ON FUNCTION sync_email_verified() IS 'Синхронизирует email_verified из email_verifications в child_profiles';
