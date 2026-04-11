-- Добавляем достижения "Стикер Дружба" за каждых 5 приглашённых друзей
-- Они появляются друг за другом и отображаются как звёздочка с цифрой

INSERT INTO achievements (id, type, title, description, icon, requirement_type, requirement_value, reward_type, reward_name, reward_amount, priority)
VALUES
    -- Дружба за 5 друзей (priority 200 - после основных 12 достижений)
    ('friendship_sticker_5', 'social', 'Стикер Дружба', 'За 5 приглашённых друзей', '⭐',
     'friends_invited', 5, 'sticker', 'Стикер Дружба', 0, 200),

    -- Дружба за 10 друзей
    ('friendship_sticker_10', 'social', 'Стикер Дружба', 'За 10 приглашённых друзей', '⭐',
     'friends_invited', 10, 'sticker', 'Стикер Дружба', 0, 201),

    -- Дружба за 15 друзей
    ('friendship_sticker_15', 'social', 'Стикер Дружба', 'За 15 приглашённых друзей', '⭐',
     'friends_invited', 15, 'sticker', 'Стикер Дружба', 0, 202),

    -- Дружба за 20 друзей
    ('friendship_sticker_20', 'social', 'Стикер Дружба', 'За 20 приглашённых друзей', '⭐',
     'friends_invited', 20, 'sticker', 'Стикер Дружба', 0, 203),

    -- Дружба за 25 друзей
    ('friendship_sticker_25', 'social', 'Стикер Дружба', 'За 25 приглашённых друзей', '⭐',
     'friends_invited', 25, 'sticker', 'Стикер Дружба', 0, 204),

    -- Дружба за 30 друзей
    ('friendship_sticker_30', 'social', 'Стикер Дружба', 'За 30 приглашённых друзей', '⭐',
     'friends_invited', 30, 'sticker', 'Стикер Дружба', 0, 205),

    -- Дружба за 40 друзей
    ('friendship_sticker_40', 'social', 'Стикер Дружба', 'За 40 приглашённых друзей', '⭐',
     'friends_invited', 40, 'sticker', 'Стикер Дружба', 0, 206),

    -- Дружба за 50 друзей
    ('friendship_sticker_50', 'social', 'Стикер Дружба', 'За 50 приглашённых друзей', '⭐',
     'friends_invited', 50, 'sticker', 'Стикер Дружба', 0, 207)
ON CONFLICT (id) DO NOTHING;

COMMENT ON COLUMN achievements.requirement_type IS
'Типы требований: streak_days, villains_defeated, tasks_correct, tasks_no_hints, stickers_collected, friends_invited';
