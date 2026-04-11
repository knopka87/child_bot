-- Rollback sync email_verified

DROP TRIGGER IF EXISTS trigger_sync_email_verified ON email_verifications;
DROP FUNCTION IF EXISTS sync_email_verified();
