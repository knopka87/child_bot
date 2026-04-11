-- Удаляем старое достижение "5 ошибок исправлено"
DELETE FROM achievements WHERE id = 'achievement_fixes_5';

-- Добавляем серию достижений "Исправленные ошибки"
-- Показывается только максимальное разблокированное

INSERT INTO achievements (id, type, title, description, icon, requirement_type, requirement_value, reward_type, reward_name, reward_amount, priority)
VALUES
    -- 1 ошибка (priority 700)
    ('errors_fixed_1', 'tasks', 'Исправленные ошибки', 'Нашёл и исправил 1 ошибку', '📝',
     'errors_found', 1, 'sticker', 'Исправленные ошибки', 0, 700),

    -- 5 ошибок
    ('errors_fixed_5', 'tasks', 'Исправленные ошибки', 'Нашёл и исправил 5 ошибок', '📝',
     'errors_found', 5, 'sticker', 'Исправленные ошибки', 0, 701),

    -- 10 ошибок
    ('errors_fixed_10', 'tasks', 'Исправленные ошибки', 'Нашёл и исправил 10 ошибок', '📝',
     'errors_found', 10, 'sticker', 'Исправленные ошибки', 0, 702),

    -- 50 ошибок
    ('errors_fixed_50', 'tasks', 'Исправленные ошибки', 'Нашёл и исправил 50 ошибок', '📝',
     'errors_found', 50, 'sticker', 'Исправленные ошибки', 0, 703),

    -- 100 ошибок
    ('errors_fixed_100', 'tasks', 'Исправленные ошибки', 'Нашёл и исправил 100 ошибок', '📝',
     'errors_found', 100, 'sticker', 'Исправленные ошибки', 0, 704),

    -- 500 ошибок
    ('errors_fixed_500', 'tasks', 'Исправленные ошибки', 'Нашёл и исправил 500 ошибок', '📝',
     'errors_found', 500, 'sticker', 'Исправленные ошибки', 0, 705)
ON CONFLICT (id) DO NOTHING;

COMMENT ON COLUMN achievements.requirement_type IS
'Типы требований: streak_days, villains_defeated, tasks_correct, tasks_no_hints, hints_used, stickers_collected, friends_invited, errors_found';
