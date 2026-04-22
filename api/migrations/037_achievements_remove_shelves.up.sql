-- Remove shelves from achievements, add priority instead

-- 1. Добавляем поле priority
ALTER TABLE achievements
ADD COLUMN priority INTEGER NOT NULL DEFAULT 100;

-- 2. Устанавливаем приоритеты для существующих достижений
-- Чем меньше число - тем выше приоритет (легче заработать)
UPDATE achievements SET priority = 10 WHERE id = 'streak_3';        -- Самое простое: 3 дня подряд
UPDATE achievements SET priority = 20 WHERE id = 'tasks_10';        -- Простое: 10 задач
UPDATE achievements SET priority = 30 WHERE id = 'streak_7';        -- Средне: 7 дней подряд
UPDATE achievements SET priority = 40 WHERE id = 'tasks_50';        -- Сложнее: 50 задач
UPDATE achievements SET priority = 50 WHERE id = 'villain_1_defeated'; -- Сложно: победить первого злодея

-- 3. Удаляем поля shelf_order и position_in_shelf
ALTER TABLE achievements
DROP COLUMN IF EXISTS shelf_order,
DROP COLUMN IF EXISTS position_in_shelf;

-- 4. Создаём индекс для сортировки неактивных достижений по приоритету
CREATE INDEX IF NOT EXISTS idx_achievements_priority
    ON achievements (priority ASC);

-- Комментарий
COMMENT ON COLUMN achievements.priority IS 'Приоритет отображения (меньше = выше приоритет, легче заработать)';
