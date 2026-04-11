-- Rollback errors fixed series

DELETE FROM achievements WHERE id IN (
    'errors_fixed_1',
    'errors_fixed_5',
    'errors_fixed_10',
    'errors_fixed_50',
    'errors_fixed_100',
    'errors_fixed_500'
);

-- Восстанавливаем старое достижение "5 ошибок исправлено"
INSERT INTO achievements (id, type, title, description, icon, requirement_type, requirement_value, reward_type, reward_name, reward_amount, priority)
VALUES
    ('achievement_fixes_5', 'tasks', '5 ошибок исправлено', 'Исправил 5 ошибок после проверки', '🛠️',
     'tasks_correct', 5, 'coins', '50 монет', 50, 30)
ON CONFLICT (id) DO NOTHING;
