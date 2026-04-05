# 15: Spaced Repetition (Интервальное повторение)

> Фаза 2 | Приоритет: P2 | Сложность: Средняя | Срок: 3-4 дня

## Цель

Система "Машина времени" — возврат к старым задачам для закрепления материала.

## Алгоритм

Основан на SM-2 (SuperMemo):

```
interval(n) = interval(n-1) * EF
EF = max(1.3, EF + (0.1 - (5-q) * (0.08 + (5-q) * 0.02)))

где q = качество ответа (0-5)
```

## Упрощённая версия для MVP

| Результат | Следующее повторение |
|-----------|---------------------|
| Верно (без подсказок) | +7 дней |
| Верно (с подсказками) | +3 дня |
| Неверно | +1 день |

## Миграция

```sql
-- migrations/037_spaced_repetition.up.sql

CREATE TABLE user_review_queue (
    id SERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES "user"(chat_id) ON DELETE CASCADE,
    session_id TEXT NOT NULL,
    topic_id TEXT REFERENCES topic(id),
    next_review DATE NOT NULL,
    interval_days INT DEFAULT 1,
    ease_factor FLOAT DEFAULT 2.5,
    review_count INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_review_queue_next ON user_review_queue(user_id, next_review);
```

## Store методы

```go
func (s *Store) GetDueReviews(ctx context.Context, userID int64, limit int) ([]ReviewItem, error) {
    rows, err := s.DB.QueryContext(ctx, `
        SELECT urq.id, urq.session_id, pt.task_text_clean, urq.topic_id
        FROM user_review_queue urq
        JOIN parsed_tasks pt ON urq.session_id = pt.session_id
        WHERE urq.user_id = $1 AND urq.next_review <= CURRENT_DATE
        ORDER BY urq.next_review
        LIMIT $2
    `, userID, limit)
    // ...
}

func (s *Store) ScheduleReview(ctx context.Context, userID int64, sessionID string, topicID string) error {
    _, err := s.DB.ExecContext(ctx, `
        INSERT INTO user_review_queue (user_id, session_id, topic_id, next_review)
        VALUES ($1, $2, $3, CURRENT_DATE + INTERVAL '1 day')
        ON CONFLICT DO NOTHING
    `, userID, sessionID, topicID)
    return err
}

func (s *Store) UpdateReview(ctx context.Context, reviewID int64, quality int) error {
    // quality: 0-2 = неверно, 3-4 = с трудом, 5 = легко

    _, err := s.DB.ExecContext(ctx, `
        UPDATE user_review_queue
        SET interval_days = CASE
                WHEN $2 < 3 THEN 1
                WHEN $2 < 5 THEN interval_days + 3
                ELSE interval_days * 2
            END,
            next_review = CURRENT_DATE + interval_days,
            review_count = review_count + 1
        WHERE id = $1
    `, reviewID, quality)
    return err
}
```

## API Endpoints

```
GET  /api/v1/reviews/due           # Задачи для повторения
POST /api/v1/reviews/{id}/complete # Отметить повторение
GET  /api/v1/reviews/calendar      # Календарь повторений
```

## Интеграция

После решения задачи добавляем в очередь повторения:

```go
func (s *TaskService) onTaskCompleted(ctx context.Context, userID int64, sessionID string) {
    topicID := s.getTopicFromSession(ctx, sessionID)
    s.store.ScheduleReview(ctx, userID, sessionID, topicID)
}
```

## Уведомления

Ежедневное напоминание если есть задачи для повторения:

```go
func (s *NotificationService) sendDailyReviewReminder() {
    users, _ := s.store.GetUsersWithDueReviews(ctx)
    for _, userID := range users {
        count, _ := s.store.CountDueReviews(ctx, userID)
        if count > 0 {
            s.hub.SendToUser(userID, "review_reminder", map[string]int{"count": count})
        }
    }
}
```

## Чек-лист

- [ ] Миграция `037_spaced_repetition.up.sql`
- [ ] Store методы
- [ ] REST API
- [ ] Добавление задач в очередь после решения
- [ ] Ежедневные уведомления
- [ ] UI календаря повторений

---

[← Knowledge Map](./14-knowledge-map.md) | [Назад к Roadmap](./roadmap.md) | [Далее: Leaderboard →](./16-leaderboard.md)
