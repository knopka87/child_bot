-- Добавляем серию достижений за streak (дни подряд)
-- Показывается только максимальное разблокированное

INSERT INTO achievements (id, type, title, description, icon, requirement_type, requirement_value, reward_type, reward_name, reward_amount, priority)
VALUES
    -- 1 день подряд (priority 300 - после дружбы)
    ('streak_series_1', 'streak', 'Стрик', 'За 1 день занятий', '🔥',
     'streak_days', 1, 'sticker', 'Стрик', 0, 300),

    -- 3 дня подряд
    ('streak_series_3', 'streak', 'Стрик', 'За 3 дня подряд', '🔥',
     'streak_days', 3, 'sticker', 'Стрик', 0, 301),

    -- 7 дней подряд
    ('streak_series_7', 'streak', 'Стрик', 'За неделю занятий подряд', '🔥',
     'streak_days', 7, 'sticker', 'Стрик', 0, 302),

    -- 30 дней подряд
    ('streak_series_30', 'streak', 'Стрик', 'За месяц занятий подряд', '🔥',
     'streak_days', 30, 'sticker', 'Стрик', 0, 303),

    -- 90 дней подряд
    ('streak_series_90', 'streak', 'Стрик', 'За 3 месяца занятий подряд', '🔥',
     'streak_days', 90, 'sticker', 'Стрик', 0, 304),

    -- 180 дней подряд
    ('streak_series_180', 'streak', 'Стрик', 'За 6 месяцев занятий подряд', '🔥',
     'streak_days', 180, 'sticker', 'Стрик', 0, 305),

    -- 365 дней подряд
    ('streak_series_365', 'streak', 'Стрик', 'За год занятий подряд', '🔥',
     'streak_days', 365, 'sticker', 'Стрик', 0, 306)
ON CONFLICT (id) DO NOTHING;
