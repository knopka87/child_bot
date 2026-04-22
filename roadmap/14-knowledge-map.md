# 14: Карта знаний

> Фаза 2 | Приоритет: P1 | Сложность: Высокая | Срок: 5-7 дней

## Цель

Визуализация прогресса ученика по темам. Дерево знаний, которое "расцветает".

## Структура тем

```
Математика 3 класс
├── Числа и арифметика
│   ├── Сложение и вычитание
│   ├── Умножение
│   ├── Деление
│   └── Порядок действий
├── Геометрия
│   ├── Фигуры
│   ├── Периметр
│   └── Площадь
└── Задачи
    ├── Задачи на движение
    ├── Задачи на сравнение
    └── Составные задачи
```

## Уровни освоения

| Уровень | Цвет | Требование |
|---------|------|------------|
| 0 | Серый | Не начато |
| 1 | Жёлтый | 1-2 задачи |
| 2 | Зелёный | 3-5 задач |
| 3 | Золотой | 6+ задач без ошибок |

## Миграция

```sql
-- migrations/036_knowledge_map.up.sql

CREATE TABLE topic (
    id TEXT PRIMARY KEY,
    subject TEXT NOT NULL,
    grade INT NOT NULL,
    name TEXT NOT NULL,
    parent_id TEXT REFERENCES topic(id),
    sort_order INT DEFAULT 0
);

CREATE TABLE user_topic_progress (
    user_id BIGINT REFERENCES "user"(chat_id) ON DELETE CASCADE,
    topic_id TEXT REFERENCES topic(id),
    tasks_solved INT DEFAULT 0,
    tasks_correct INT DEFAULT 0,
    mastery_level INT DEFAULT 0,  -- 0-3
    last_practiced TIMESTAMPTZ,
    PRIMARY KEY (user_id, topic_id)
);

CREATE INDEX idx_user_topic_progress_user ON user_topic_progress(user_id);
```

## Store методы

```go
func (s *Store) GetKnowledgeMap(ctx context.Context, userID int64, subject string, grade int) ([]TopicWithProgress, error) {
    rows, err := s.DB.QueryContext(ctx, `
        SELECT t.id, t.name, t.parent_id, t.sort_order,
               COALESCE(utp.tasks_solved, 0),
               COALESCE(utp.tasks_correct, 0),
               COALESCE(utp.mastery_level, 0)
        FROM topic t
        LEFT JOIN user_topic_progress utp ON t.id = utp.topic_id AND utp.user_id = $1
        WHERE t.subject = $2 AND t.grade = $3
        ORDER BY t.sort_order
    `, userID, subject, grade)
    // ...
}

func (s *Store) UpdateTopicProgress(ctx context.Context, userID int64, topicID string, correct bool) error {
    _, err := s.DB.ExecContext(ctx, `
        INSERT INTO user_topic_progress (user_id, topic_id, tasks_solved, tasks_correct, last_practiced)
        VALUES ($1, $2, 1, $3::int, NOW())
        ON CONFLICT (user_id, topic_id) DO UPDATE
        SET tasks_solved = user_topic_progress.tasks_solved + 1,
            tasks_correct = user_topic_progress.tasks_correct + $3::int,
            mastery_level = CASE
                WHEN user_topic_progress.tasks_correct + $3::int >= 6 THEN 3
                WHEN user_topic_progress.tasks_solved + 1 >= 3 THEN 2
                ELSE 1
            END,
            last_practiced = NOW()
    `, userID, topicID, correct)
    return err
}
```

## API Endpoints

```
GET /api/v1/knowledge-map
    ?subject=math
    &grade=3
```

## Интеграция

При решении задачи определяем топик из PARSE response и обновляем прогресс:

```go
func (s *GamificationService) onTaskChecked(event events.Event) {
    topicID := event.Payload["topic_id"].(string)
    isCorrect := event.Payload["is_correct"].(bool)

    s.store.UpdateTopicProgress(ctx, event.UserID, topicID, isCorrect)
}
```

## Чек-лист

- [ ] Миграция `036_knowledge_map.up.sql`
- [ ] Заполнить таблицу topics для 1-4 классов
- [ ] Store методы
- [ ] REST API
- [ ] Интеграция с task checking
- [ ] Маппинг template_id → topic_id
- [ ] UI дерева знаний (ТЗ дизайнеру)

---

[← Customization](./13-customization.md) | [Назад к Roadmap](./roadmap.md) | [Далее: Spaced Repetition →](./15-spaced-repetition.md)
