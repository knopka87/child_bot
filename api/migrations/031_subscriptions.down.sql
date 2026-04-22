DROP TRIGGER IF EXISTS subscriptions_updated_at ON subscriptions;
DROP INDEX IF EXISTS idx_subscriptions_expires;
DROP INDEX IF EXISTS idx_subscriptions_profile;
DROP INDEX IF EXISTS idx_subscription_plans_active;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS subscription_plans;
