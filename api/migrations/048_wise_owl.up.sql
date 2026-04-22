-- Добавляем достижение "Мудрая сова" за первую подсказку

INSERT INTO achievements (id, type, title, description, icon, requirement_type, requirement_value, reward_type, reward_name, reward_amount, priority)
VALUES
    -- Мудрая сова - за 1 подсказку (priority 600)
    ('wise_owl_1', 'mastery', 'Мудрая сова', 'Использовал первую подсказку', '🦉',
     'hints_used', 1, 'badge', 'Мудрая сова', 0, 600)
ON CONFLICT (id) DO NOTHING;

COMMENT ON COLUMN achievements.requirement_type IS
'Типы требований: streak_days, villains_defeated, tasks_correct, tasks_no_hints, hints_used, stickers_collected, friends_invited';
