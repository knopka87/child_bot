-- Rollback: restore shelves for achievements

-- 1. Добавляем обратно поля shelf_order и position_in_shelf
ALTER TABLE achievements
ADD COLUMN shelf_order INTEGER DEFAULT 1,
ADD COLUMN position_in_shelf INTEGER DEFAULT 0;

-- 2. Восстанавливаем старые значения
UPDATE achievements SET shelf_order = 1, position_in_shelf = 0 WHERE id = 'streak_3';
UPDATE achievements SET shelf_order = 1, position_in_shelf = 1 WHERE id = 'streak_7';
UPDATE achievements SET shelf_order = 1, position_in_shelf = 2 WHERE id = 'tasks_10';
UPDATE achievements SET shelf_order = 1, position_in_shelf = 3 WHERE id = 'tasks_50';
UPDATE achievements SET shelf_order = 2, position_in_shelf = 0 WHERE id = 'villain_1_defeated';

-- 3. Удаляем индекс priority
DROP INDEX IF EXISTS idx_achievements_priority;

-- 4. Удаляем поле priority
ALTER TABLE achievements
DROP COLUMN IF EXISTS priority;
