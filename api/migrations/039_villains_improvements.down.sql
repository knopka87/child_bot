-- Rollback villains improvements

-- Удаляем добавленных злодеев
DELETE FROM villains WHERE id IN (
    'sir_procrastination',
    'madame_mistake',
    'lord_laziness',
    'boss_week_chaos'
);

-- Удаляем индекс
DROP INDEX IF EXISTS idx_villains_is_boss;

-- Удаляем поле is_boss
ALTER TABLE villains
DROP COLUMN IF EXISTS is_boss;
