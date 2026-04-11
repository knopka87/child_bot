-- Rollback friendship stickers

DELETE FROM achievements WHERE id IN (
    'friendship_sticker_5',
    'friendship_sticker_10',
    'friendship_sticker_15',
    'friendship_sticker_20',
    'friendship_sticker_25',
    'friendship_sticker_30',
    'friendship_sticker_40',
    'friendship_sticker_50'
);
