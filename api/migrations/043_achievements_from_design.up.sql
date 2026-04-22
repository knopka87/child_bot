-- Replace achievements with design version
-- Удаляем старые примеры
DELETE FROM achievements WHERE id IN (
    'streak_3', 'streak_7', 'streak_14', 'streak_30',
    'villains_1', 'villains_3', 'villains_5', 'villains_10',
    'tasks_10', 'tasks_25', 'tasks_50', 'tasks_100',
    'no_hints_5', 'no_hints_10', 'no_hints_25', 'no_hints_50',
    'villain_1_defeated'
);

-- Вставляем достижения из дизайна (12 штук, 3 полки по 4)
INSERT INTO achievements (id, type, title, description, icon, requirement_type, requirement_value, reward_type, reward_name, reward_amount, priority)
VALUES
    -- Полка 1 (id 1-4, priority 10-40)
    ('achievement_streak_5', 'streak', '5 дней подряд', 'Занимался 5 дней подряд', '🔥',
     'streak_days', 5, 'coins', '50 монет', 50, 10),

    ('achievement_checks_10', 'tasks', '10 проверок ДЗ', 'Проверил 10 домашних заданий', '✅',
     'tasks_correct', 10, 'coins', '100 монет', 100, 20),

    ('achievement_fixes_5', 'fixes', '5 ошибок исправлено', 'Исправил 5 ошибок после проверки', '⭐',
     'tasks_correct', 5, 'coins', '50 монет', 50, 30),

    ('achievement_first_task', 'milestones', 'Первое задание', 'За первое решённое задание', '🎯',
     'tasks_correct', 1, 'coins', '25 монет', 25, 40),

    -- Полка 2 (id 5-8, priority 50-80)
    ('achievement_speed_solver', 'speed', 'Скоростной решатель', 'За быструю работу', '⚡',
     'tasks_correct', 5, 'coins', '150 монет', 150, 50),

    ('achievement_villain_defeater', 'villain', 'Победитель злодеев', 'За победу над злодеями', '🏆',
     'villains_defeated', 3, 'sticker', 'Стикер чемпиона', 0, 60),

    ('achievement_wise_owl', 'wisdom', 'Мудрая сова', 'За использование подсказок', '🦉',
     'tasks_correct', 5, 'coins', '75 монет', 75, 70),

    ('achievement_collector', 'collection', 'Коллекционер', 'За сбор стикеров', '💎',
     'stickers_collected', 10, 'badge', 'Значок коллекционера', 0, 80),

    -- Полка 3 (id 9-12, priority 90-120)
    ('achievement_knowledge_rocket', 'mastery', 'Ракета знаний', 'За 20 решённых заданий', '🚀',
     'tasks_correct', 20, 'coins', '200 монет', 200, 90),

    ('achievement_superstar', 'perfect', 'Суперзвезда', 'За 10 безошибочных проверок', '🌟',
     'tasks_no_hints', 10, 'coins', '300 монет', 300, 100),

    ('achievement_marathoner', 'streak', 'Марафонец', 'За 7 дней подряд', '🎪',
     'streak_days', 7, 'coins', '100 монет', 100, 110),

    ('achievement_genius', 'mastery', 'Гений', 'За 50 решённых заданий', '🧠',
     'tasks_correct', 50, 'coins', '500 монет', 500, 120)
ON CONFLICT (id) DO NOTHING;
