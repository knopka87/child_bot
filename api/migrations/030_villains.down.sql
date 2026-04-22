DROP TRIGGER IF EXISTS villain_battles_updated_at ON villain_battles;
DROP INDEX IF EXISTS idx_damage_events_battle;
DROP INDEX IF EXISTS idx_villain_battles_active;
DROP INDEX IF EXISTS idx_villain_battles_profile_status;
DROP INDEX IF EXISTS idx_villains_unlock_order;
DROP TABLE IF EXISTS damage_events;
DROP TABLE IF EXISTS villain_battles;
DROP TABLE IF EXISTS villains;
