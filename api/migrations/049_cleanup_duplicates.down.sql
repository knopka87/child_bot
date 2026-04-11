-- Rollback cleanup duplicates

-- 1. Восстанавливаем "Стикер Дружба" → "Дружба"
UPDATE achievements
SET title = 'Стикер Дружба',
    reward_name = 'Стикер Дружба'
WHERE reward_name = 'Дружба';

-- 2. Восстанавливаем "Победитель злодеев" (Стикер чемпиона за 3 злодеев)
INSERT INTO achievements (id, type, title, description, icon, requirement_type, requirement_value, reward_type, reward_name, reward_amount, priority)
VALUES
    ('achievement_villain_defeater', 'villain', 'Победитель злодеев', 'За победу над злодеями', '🦹',
     'villains_defeated', 3, 'sticker', 'Стикер чемпиона', 0, 60)
ON CONFLICT (id) DO NOTHING;

-- 3. Восстанавливаем "5 дней подряд"
INSERT INTO achievements (id, type, title, description, icon, requirement_type, requirement_value, reward_type, reward_name, reward_amount, priority)
VALUES
    ('achievement_streak_5', 'streak', '5 дней подряд', 'Занимался 5 дней подряд', '🔥',
     'streak_days', 5, 'coins', '50 монет', 50, 10)
ON CONFLICT (id) DO NOTHING;

-- 4. Восстанавливаем Мудрую сову (75 монет за 5 подсказок)
INSERT INTO achievements (id, type, title, description, icon, requirement_type, requirement_value, reward_type, reward_name, reward_amount, priority)
VALUES
    ('achievement_wise_owl', 'mastery', 'Мудрая сова', 'За использование подсказок', '🦉',
     'hints_used', 5, 'coins', '75 монет', 75, 70)
ON CONFLICT (id) DO NOTHING;
