-- Rollback streak series

DELETE FROM achievements WHERE id IN (
    'streak_series_1',
    'streak_series_3',
    'streak_series_7',
    'streak_series_30',
    'streak_series_90',
    'streak_series_180',
    'streak_series_365'
);
