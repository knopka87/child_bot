-- Rollback achievement examples

DELETE FROM achievements WHERE id IN (
    -- Streak achievements
    'streak_3', 'streak_7', 'streak_14', 'streak_30',
    -- Villain achievements
    'villains_1', 'villains_3', 'villains_5', 'villains_10',
    -- Tasks correct achievements
    'tasks_10', 'tasks_25', 'tasks_50', 'tasks_100',
    -- Tasks no hints achievements
    'no_hints_5', 'no_hints_10', 'no_hints_25', 'no_hints_50'
);
