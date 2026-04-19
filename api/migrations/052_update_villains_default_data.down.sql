-- Откатываем к старым значениям
UPDATE villains SET max_hp = 100, damage_per_correct_task = 5, reward_coins = 100, level = 1 WHERE id = 'count_error';
UPDATE villains SET max_hp = 150, damage_per_correct_task = 5, reward_coins = 150, level = 2 WHERE id = 'baron_confusion';
UPDATE villains SET max_hp = 200, damage_per_correct_task = 5, reward_coins = 200, level = 2 WHERE id = 'duchess_distraction';
UPDATE villains SET max_hp = 120, damage_per_correct_task = 5, reward_coins = 120, level = 3 WHERE id = 'sir_procrastination';
UPDATE villains SET max_hp = 140, damage_per_correct_task = 5, reward_coins = 140, level = 3 WHERE id = 'madame_mistake';
UPDATE villains SET max_hp = 160, damage_per_correct_task = 5, reward_coins = 160, level = 4 WHERE id = 'lord_laziness';
UPDATE villains SET max_hp = 300, damage_per_correct_task = 10, reward_coins = 500, level = 5 WHERE id = 'boss_week_chaos';
