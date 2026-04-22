-- Удаляем неиспользуемые таблицы referral milestones (заменены на achievements)
DROP TABLE IF EXISTS child_referral_milestones CASCADE;
DROP TABLE IF EXISTS referral_milestones CASCADE;

-- Удаляем deprecated таблицы Telegram бота (v1 и v2)
DROP TABLE IF EXISTS chat CASCADE;
DROP TABLE IF EXISTS task_sessions CASCADE;
DROP TABLE IF EXISTS timeline_events CASCADE;
DROP TABLE IF EXISTS metrics_events CASCADE;
DROP TABLE IF EXISTS hints_cache CASCADE;
DROP TABLE IF EXISTS parsed_tasks CASCADE;
DROP TABLE IF EXISTS "user" CASCADE;

-- Комментарий: Все эти таблицы были заменены на новую REST API архитектуру
-- referral_milestones → achievements (стикеры "Дружба" за 5, 10, 15... друзей)
-- Telegram таблицы → attempts, child_profiles (унифицированная система)
