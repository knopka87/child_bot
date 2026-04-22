-- Rollback villains series

DELETE FROM achievements WHERE id IN (
    'villains_series_1',
    'villains_series_5',
    'villains_series_10',
    'villains_series_50',
    'villains_series_100',
    'villains_series_500',
    'villains_series_1000'
);
