-- Обновляем злодеев актуальными данными
-- unlock_order: 1=понедельник, 7=воскресенье

UPDATE villains SET 
    max_hp = 100, 
    damage_per_correct_task = 20, 
    reward_coins = 100, 
    level = 1,
    image_url = '/assets/villains/count_error.png'
WHERE id = 'count_error';

UPDATE villains SET 
    max_hp = 120, 
    damage_per_correct_task = 20, 
    reward_coins = 150, 
    level = 2,
    image_url = '/assets/villains/baron_confusion.png'
WHERE id = 'baron_confusion';

UPDATE villains SET 
    max_hp = 120, 
    damage_per_correct_task = 20, 
    reward_coins = 200, 
    level = 2,
    image_url = '/assets/villains/duchess_distraction.png'
WHERE id = 'duchess_distraction';

UPDATE villains SET 
    max_hp = 140, 
    damage_per_correct_task = 20, 
    reward_coins = 120, 
    level = 3,
    image_url = '/assets/villains/sir_procrastination.png'
WHERE id = 'sir_procrastination';

UPDATE villains SET 
    max_hp = 140, 
    damage_per_correct_task = 20, 
    reward_coins = 140, 
    level = 3,
    image_url = '/assets/villains/madame_mistake.png'
WHERE id = 'madame_mistake';

UPDATE villains SET 
    max_hp = 160, 
    damage_per_correct_task = 20, 
    reward_coins = 160, 
    level = 4,
    image_url = '/assets/villains/lord_laziness.png'
WHERE id = 'lord_laziness';

UPDATE villains SET 
    max_hp = 200, 
    damage_per_correct_task = 20, 
    reward_coins = 500, 
    level = 5,
    image_url = '/assets/villains/boss_week_chaos.png'
WHERE id = 'boss_week_chaos';
