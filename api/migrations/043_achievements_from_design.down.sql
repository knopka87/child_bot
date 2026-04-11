-- Rollback achievements from design

DELETE FROM achievements WHERE id IN (
    'achievement_streak_5',
    'achievement_checks_10',
    'achievement_fixes_5',
    'achievement_first_task',
    'achievement_speed_solver',
    'achievement_villain_defeater',
    'achievement_wise_owl',
    'achievement_collector',
    'achievement_knowledge_rocket',
    'achievement_superstar',
    'achievement_marathoner',
    'achievement_genius'
);

-- Restore previous achievements from 042
INSERT INTO achievements (id, type, title, description, icon, requirement_type, requirement_value,
                          reward_type, reward_name, reward_amount, priority)
VALUES
    ('streak_3', 'streak', '3 дня подряд', 'Решай задачи 3 дня подряд', '🔥',
     'streak_days', 3, 'coins', '30 монет', 30, 10),

    ('streak_7', 'streak', 'Неделя подряд', 'Решай задачи 7 дней подряд', '🔥',
     'streak_days', 7, 'coins', '70 монет', 70, 20),

    ('streak_14', 'streak', '2 недели подряд', 'Решай задачи 14 дней подряд', '🔥',
     'streak_days', 14, 'coins', '140 монет', 140, 30),

    ('streak_30', 'streak', 'Месяц подряд', 'Решай задачи 30 дней подряд', '🔥',
     'streak_days', 30, 'coins', '300 монет', 300, 40),

    ('villains_1', 'villain', 'Первый монстр', 'Победи первого монстра', '⚔️',
     'villains_defeated', 1, 'coins', '50 монет', 50, 50),

    ('villains_3', 'villain', 'Охотник на монстров', 'Победи 3 монстров', '⚔️',
     'villains_defeated', 3, 'coins', '150 монет', 150, 60),

    ('villains_5', 'villain', 'Истребитель монстров', 'Победи 5 монстров', '⚔️',
     'villains_defeated', 5, 'coins', '250 монет', 250, 70),

    ('villains_10', 'villain', 'Гроза монстров', 'Победи 10 монстров', '⚔️',
     'villains_defeated', 10, 'coins', '500 монет', 500, 80),

    ('tasks_10', 'tasks', '10 задач', 'Реши правильно 10 задач', '✅',
     'tasks_correct', 10, 'coins', '100 монет', 100, 90),

    ('tasks_25', 'tasks', '25 задач', 'Реши правильно 25 задач', '✅',
     'tasks_correct', 25, 'coins', '250 монет', 250, 100),

    ('tasks_50', 'tasks', '50 задач', 'Реши правильно 50 задач', '✅',
     'tasks_correct', 50, 'coins', '500 монет', 500, 110),

    ('tasks_100', 'tasks', '100 задач', 'Реши правильно 100 задач', '✅',
     'tasks_correct', 100, 'coins', '1000 монет', 1000, 120),

    ('no_hints_5', 'perfect', '5 без подсказок', 'Реши 5 задач без подсказок', '🧠',
     'tasks_no_hints', 5, 'coins', '100 монет', 100, 130),

    ('no_hints_10', 'perfect', '10 без подсказок', 'Реши 10 задач без подсказок', '🧠',
     'tasks_no_hints', 10, 'coins', '200 монет', 200, 140),

    ('no_hints_25', 'perfect', '25 без подсказок', 'Реши 25 задач без подсказок', '🧠',
     'tasks_no_hints', 25, 'coins', '500 монет', 500, 150),

    ('no_hints_50', 'perfect', 'Гений без подсказок', 'Реши 50 задач без подсказок', '🧠',
     'tasks_no_hints', 50, 'coins', '1000 монет', 1000, 160)
ON CONFLICT (id) DO NOTHING;
