-- Drop triggers
DROP TRIGGER IF EXISTS payments_updated_at ON payments;

-- Drop indexes
DROP INDEX IF EXISTS idx_payment_events_type;
DROP INDEX IF EXISTS idx_payment_events_payment;
DROP INDEX IF EXISTS idx_payments_expires;
DROP INDEX IF EXISTS idx_payments_status;
DROP INDEX IF EXISTS idx_payments_vk_order;
DROP INDEX IF EXISTS idx_payments_profile;
DROP INDEX IF EXISTS idx_payments_subscription;

-- Drop tables
DROP TABLE IF EXISTS payment_events;
DROP TABLE IF EXISTS payments;
