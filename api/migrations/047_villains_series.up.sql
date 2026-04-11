-- Добавляем серию достижений за побеждённых злодеев
-- Показывается только максимальное разблокированное

INSERT INTO achievements (id, type, title, description, icon, requirement_type, requirement_value, reward_type, reward_name, reward_amount, priority)
VALUES
    -- 1 злодей (priority 500)
    ('villains_series_1', 'villain', 'Победитель злодеев', 'За победу над 1 злодеем', '🦹',
     'villains_defeated', 1, 'sticker', 'Победитель злодеев', 0, 500),

    -- 5 злодеев
    ('villains_series_5', 'villain', 'Победитель злодеев', 'За победу над 5 злодеями', '🦹',
     'villains_defeated', 5, 'sticker', 'Победитель злодеев', 0, 501),

    -- 10 злодеев
    ('villains_series_10', 'villain', 'Победитель злодеев', 'За победу над 10 злодеями', '🦹',
     'villains_defeated', 10, 'sticker', 'Победитель злодеев', 0, 502),

    -- 50 злодеев
    ('villains_series_50', 'villain', 'Победитель злодеев', 'За победу над 50 злодеями', '🦹',
     'villains_defeated', 50, 'sticker', 'Победитель злодеев', 0, 503),

    -- 100 злодеев
    ('villains_series_100', 'villain', 'Победитель злодеев', 'За победу над 100 злодеями', '🦹',
     'villains_defeated', 100, 'sticker', 'Победитель злодеев', 0, 504),

    -- 500 злодеев
    ('villains_series_500', 'villain', 'Победитель злодеев', 'За победу над 500 злодеями', '🦹',
     'villains_defeated', 500, 'sticker', 'Победитель злодеев', 0, 505),

    -- 1000 злодеев
    ('villains_series_1000', 'villain', 'Победитель злодеев', 'За победу над 1000 злодеями', '🦹',
     'villains_defeated', 1000, 'sticker', 'Победитель злодеев', 0, 506)
ON CONFLICT (id) DO NOTHING;
