-- Rollback checks series

DELETE FROM achievements WHERE id IN (
    'checks_series_1',
    'checks_series_10',
    'checks_series_100',
    'checks_series_500',
    'checks_series_1000'
);
