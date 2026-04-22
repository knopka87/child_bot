-- Improvements to villains system

-- 1. Добавляем поле is_boss (для задачи #7)
ALTER TABLE villains
ADD COLUMN is_boss BOOLEAN NOT NULL DEFAULT FALSE;

-- 2. reward_achievement_id уже nullable, но добавим комментарий
COMMENT ON COLUMN villains.reward_achievement_id IS 'ID достижения за победу (nullable - не всегда есть)';
COMMENT ON COLUMN villains.is_boss IS 'Является ли монстр боссом (появляется раз в неделю)';

-- 3. Увеличиваем награду за боссов и помечаем их
UPDATE villains SET is_boss = FALSE WHERE level <= 2;
-- Пока нет боссов, добавим позже

-- 4. Индекс для поиска активных и босс-монстров
CREATE INDEX IF NOT EXISTS idx_villains_is_boss
    ON villains (is_boss, unlock_order);

-- 5. Добавляем больше злодеев с разными наградами
INSERT INTO villains (id, name, description, image_url, max_hp, level, damage_per_correct_task,
                      unlock_order, reward_coins, reward_achievement_id, is_boss)
VALUES
    -- Обычные злодеи (4-6)
    ('sir_procrastination', 'Сэр Прокрастинация', 'Заставляет откладывать решение задач',
     '/assets/villains/sir_procrastination.png', 120, 2, 5, 4, 120, NULL, FALSE),

    ('madame_mistake', 'Мадам Ошибка', 'Подсовывает ошибки в вычисления',
     '/assets/villains/madame_mistake.png', 140, 3, 5, 5, 140, NULL, FALSE),

    ('lord_laziness', 'Лорд Лень', 'Внушает желание не делать домашку',
     '/assets/villains/lord_laziness.png', 160, 3, 5, 6, 160, NULL, FALSE),

    -- Первый босс (появляется раз в неделю)
    ('boss_week_chaos', 'БОСС: Хаос Недели', 'Могущественный босс, приходящий раз в неделю',
     '/assets/villains/boss_week_chaos.png', 300, 5, 10, 7, 500, NULL, TRUE)
ON CONFLICT (id) DO NOTHING;

-- Комментарий к таблице
COMMENT ON TABLE villains IS 'Справочник злодеев: обычные монстры и боссы (раз в неделю)';
