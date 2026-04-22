# 16: Таблица лидеров

> Фаза 2 | Приоритет: P2 | Сложность: Низкая | Срок: 2 дня

## Цель

Таблица лидеров для здоровой конкуренции. Опциональная, включается родителем.

## Типы лидербордов

| Тип | Период | Метрика |
|-----|--------|---------|
| Daily | День | Задачи сегодня |
| Weekly | Неделя | Задачи за неделю |
| All-time | Всё время | Общий XP |
| Boss | Неделя | Урон боссу |

## Миграция

```sql
-- migrations/038_leaderboard.up.sql

CREATE TABLE user_leaderboard_settings (
    user_id BIGINT PRIMARY KEY REFERENCES "user"(chat_id),
    show_in_leaderboard BOOLEAN DEFAULT false,
    display_name TEXT  -- NULL = использовать first_name
);

-- Материализованное представление для производительности
CREATE MATERIALIZED VIEW leaderboard_weekly AS
SELECT
    u.chat_id as user_id,
    COALESCE(uls.display_name, c.first_name) as display_name,
    COUNT(DISTINCT te.task_session_id) as tasks_count,
    u.xp
FROM "user" u
JOIN chat c ON u.chat_id = c.id
LEFT JOIN user_leaderboard_settings uls ON u.chat_id = uls.user_id
LEFT JOIN timeline_events te ON u.chat_id = te.chat_id
    AND te.event_type = 'api_check'
    AND te.ok = true
    AND te.created_at >= date_trunc('week', CURRENT_DATE)
WHERE COALESCE(uls.show_in_leaderboard, false) = true
GROUP BY u.chat_id, uls.display_name, c.first_name, u.xp
ORDER BY tasks_count DESC;

CREATE UNIQUE INDEX idx_leaderboard_weekly_user ON leaderboard_weekly(user_id);

-- Обновление каждые 5 минут через cron или pg_cron
```

## Store методы

```go
func (s *Store) GetLeaderboard(ctx context.Context, period string, limit int) ([]LeaderboardEntry, error) {
    var query string
    switch period {
    case "weekly":
        query = `SELECT user_id, display_name, tasks_count, xp FROM leaderboard_weekly LIMIT $1`
    case "alltime":
        query = `
            SELECT u.chat_id, COALESCE(uls.display_name, c.first_name), 0, u.xp
            FROM "user" u
            JOIN chat c ON u.chat_id = c.id
            LEFT JOIN user_leaderboard_settings uls ON u.chat_id = uls.user_id
            WHERE COALESCE(uls.show_in_leaderboard, false) = true
            ORDER BY u.xp DESC
            LIMIT $1
        `
    }
    // ...
}

func (s *Store) GetUserRank(ctx context.Context, userID int64, period string) (int, error) {
    // ...
}

func (s *Store) SetLeaderboardVisibility(ctx context.Context, userID int64, visible bool) error {
    _, err := s.DB.ExecContext(ctx, `
        INSERT INTO user_leaderboard_settings (user_id, show_in_leaderboard)
        VALUES ($1, $2)
        ON CONFLICT (user_id) DO UPDATE SET show_in_leaderboard = $2
    `, userID, visible)
    return err
}
```

## API Endpoints

```
GET  /api/v1/leaderboard?period=weekly&limit=10
GET  /api/v1/leaderboard/me?period=weekly
POST /api/v1/leaderboard/settings
```

## Родительский контроль

Родитель включает/выключает участие ребёнка в лидерборде:

```go
func (h *ParentHandler) SetChildLeaderboardVisibility(w http.ResponseWriter, r *http.Request) {
    // Проверяем parent-child связь
    // ...
    h.store.SetLeaderboardVisibility(ctx, childID, visible)
}
```

## Чек-лист

- [ ] Миграция `038_leaderboard.up.sql`
- [ ] Store методы
- [ ] REST API
- [ ] Обновление materialized view (cron)
- [ ] Родительский контроль
- [ ] Privacy: опт-ин по умолчанию

---

[← Spaced Repetition](./15-spaced-repetition.md) | [Назад к Roadmap](./roadmap.md) | [Далее: Family Quests →](./17-family-quests.md)
