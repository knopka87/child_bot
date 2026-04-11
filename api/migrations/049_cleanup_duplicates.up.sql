-- Очистка дубликатов и переименование наград

-- 1. Удаляем дубликат Мудрой совы (75 монет за 5 подсказок)
DELETE FROM achievements WHERE id = 'achievement_wise_owl';

-- 2. Удаляем дубликат "5 дней подряд" (есть коллекционный Стрик)
DELETE FROM achievements WHERE id = 'achievement_streak_5';

-- 3. Удаляем дубликат "Победитель злодеев" (Стикер чемпиона за 3 злодеев)
DELETE FROM achievements WHERE id = 'achievement_villain_defeater';

-- 4. Переименовываем "Стикер Дружба" → "Дружба"
UPDATE achievements
SET title = 'Дружба',
    reward_name = 'Дружба'
WHERE reward_name = 'Стикер Дружба';
