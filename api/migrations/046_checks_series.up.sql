-- Добавляем серию достижений за проверки ДЗ
-- Показывается только максимальное разблокированное

INSERT INTO achievements (id, type, title, description, icon, requirement_type, requirement_value, reward_type, reward_name, reward_amount, priority)
VALUES
    -- 1 проверка ДЗ (priority 400)
    ('checks_series_1', 'tasks', 'Проверки ДЗ', 'За 1 проверку домашнего задания', '✅',
     'tasks_correct', 1, 'sticker', 'Проверки ДЗ', 0, 400),

    -- 10 проверок ДЗ
    ('checks_series_10', 'tasks', 'Проверки ДЗ', 'За 10 проверок домашних заданий', '✅',
     'tasks_correct', 10, 'sticker', 'Проверки ДЗ', 0, 401),

    -- 100 проверок ДЗ
    ('checks_series_100', 'tasks', 'Проверки ДЗ', 'За 100 проверок домашних заданий', '✅',
     'tasks_correct', 100, 'sticker', 'Проверки ДЗ', 0, 402),

    -- 500 проверок ДЗ
    ('checks_series_500', 'tasks', 'Проверки ДЗ', 'За 500 проверок домашних заданий', '✅',
     'tasks_correct', 500, 'sticker', 'Проверки ДЗ', 0, 403),

    -- 1000 проверок ДЗ
    ('checks_series_1000', 'tasks', 'Проверки ДЗ', 'За 1000 проверок домашних заданий', '✅',
     'tasks_correct', 1000, 'sticker', 'Проверки ДЗ', 0, 404)
ON CONFLICT (id) DO NOTHING;
