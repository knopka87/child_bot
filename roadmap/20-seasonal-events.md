# 20: Сезонные события

> Фаза 3 | Приоритет: P3 | Сложность: Средняя | Срок: 3-4 дня (на событие)

## Цель

Ограниченные по времени события с уникальными наградами для повышения engagement.

## Календарь событий

| Событие | Период | Тема |
|---------|--------|------|
| Новый год | 20 дек - 10 янв | Зима, подарки |
| 1 сентября | 25 авг - 5 сен | Школа, знания |
| День знаний | 1-7 сен | Учёба |
| Хэллоуин | 25-31 окт | Тыквы, призраки |
| 8 марта | 1-8 мар | Весна, цветы |

## Компоненты события

1. **Тематические предметы** — ограниченные cosmetics
2. **Особые достижения** — только во время события
3. **Тематический босс** — уникальный еженедельный босс
4. **Миссии события** — специальные задания

## Миграция

```sql
-- migrations/041_seasonal_events.up.sql

CREATE TABLE seasonal_event (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    theme_color TEXT,
    boss_type_id TEXT REFERENCES boss_type(id),
    is_active BOOLEAN DEFAULT false
);

CREATE TABLE seasonal_mission (
    id TEXT PRIMARY KEY,
    event_id TEXT REFERENCES seasonal_event(id),
    name TEXT NOT NULL,
    description TEXT,
    target_type TEXT NOT NULL,   -- tasks, streak, correct_answers
    target_value INT NOT NULL,
    reward_item_id TEXT REFERENCES cosmetic_item(id),
    reward_xp INT DEFAULT 0
);

CREATE TABLE user_seasonal_progress (
    user_id BIGINT REFERENCES "user"(chat_id),
    mission_id TEXT REFERENCES seasonal_mission(id),
    current_value INT DEFAULT 0,
    completed_at TIMESTAMPTZ,
    reward_claimed BOOLEAN DEFAULT false,
    PRIMARY KEY (user_id, mission_id)
);

-- Пример события
INSERT INTO seasonal_event VALUES
('newyear_2025', 'Новогодний марафон', 'Помоги Деду Морозу собрать знания!',
 '2024-12-20', '2025-01-10', '#1E40AF', 'frost_dragon', false);

INSERT INTO seasonal_mission VALUES
('ny25_tasks_10', 'newyear_2025', 'Снежные задачки', 'Реши 10 задач', 'tasks', 10, 'hat_santa', 50),
('ny25_streak_5', 'newyear_2025', 'Новогодняя серия', '5 дней подряд', 'streak', 5, 'scarf_red', 100),
('ny25_perfect', 'newyear_2025', 'Без ошибок', '5 задач подряд правильно', 'correct_streak', 5, 'star_gold', 150);
```

## Store методы

```go
func (s *Store) GetActiveEvent(ctx context.Context) (*SeasonalEvent, error) {
    var event SeasonalEvent
    err := s.DB.QueryRowContext(ctx, `
        SELECT id, name, description, start_date, end_date, theme_color, boss_type_id
        FROM seasonal_event
        WHERE is_active = true
          AND CURRENT_DATE BETWEEN start_date AND end_date
    `).Scan(/* ... */)
    return &event, err
}

func (s *Store) GetEventMissions(ctx context.Context, eventID string, userID int64) ([]MissionWithProgress, error) {
    rows, err := s.DB.QueryContext(ctx, `
        SELECT sm.id, sm.name, sm.description, sm.target_type, sm.target_value,
               COALESCE(usp.current_value, 0), usp.completed_at IS NOT NULL
        FROM seasonal_mission sm
        LEFT JOIN user_seasonal_progress usp ON sm.id = usp.mission_id AND usp.user_id = $2
        WHERE sm.event_id = $1
    `, eventID, userID)
    // ...
}

func (s *Store) UpdateMissionProgress(ctx context.Context, userID int64, missionID string, delta int) (*MissionProgress, error) {
    var progress MissionProgress
    err := s.DB.QueryRowContext(ctx, `
        INSERT INTO user_seasonal_progress (user_id, mission_id, current_value)
        VALUES ($1, $2, $3)
        ON CONFLICT (user_id, mission_id) DO UPDATE
        SET current_value = user_seasonal_progress.current_value + $3
        RETURNING current_value, completed_at IS NOT NULL
    `, userID, missionID, delta).Scan(&progress.CurrentValue, &progress.Completed)

    // Проверяем достижение цели
    if !progress.Completed {
        var targetValue int
        s.DB.QueryRowContext(ctx, `SELECT target_value FROM seasonal_mission WHERE id = $1`, missionID).Scan(&targetValue)

        if progress.CurrentValue >= targetValue {
            s.DB.ExecContext(ctx, `
                UPDATE user_seasonal_progress
                SET completed_at = NOW()
                WHERE user_id = $1 AND mission_id = $2
            `, userID, missionID)
            progress.Completed = true
            progress.JustCompleted = true
        }
    }

    return &progress, err
}
```

## API Endpoints

```
GET  /api/v1/events/current            # Текущее событие
GET  /api/v1/events/{id}/missions      # Миссии события
POST /api/v1/events/missions/{id}/claim # Забрать награду
```

## Интеграция

```go
// При решении задачи обновляем прогресс всех активных миссий
func (s *GamificationService) onTaskCorrect(event events.Event) {
    ctx := context.Background()
    userID := event.UserID

    // Получаем активные миссии пользователя
    activeEvent, _ := s.store.GetActiveEvent(ctx)
    if activeEvent == nil {
        return
    }

    missions, _ := s.store.GetEventMissions(ctx, activeEvent.ID, userID)

    for _, mission := range missions {
        if mission.Completed {
            continue
        }

        var delta int
        switch mission.TargetType {
        case "tasks":
            delta = 1
        case "correct_streak":
            // Проверяем, не было ли ошибок
            delta = 1 // или 0 если была ошибка
        }

        if delta > 0 {
            progress, _ := s.store.UpdateMissionProgress(ctx, userID, mission.ID, delta)
            if progress.JustCompleted {
                // Уведомление через WebSocket
                s.hub.SendToUser(userID, "mission_completed", map[string]any{
                    "mission": mission,
                })
            }
        }
    }
}
```

## Admin API (для управления событиями)

```
POST /api/admin/events              # Создать событие
PUT  /api/admin/events/{id}/activate # Активировать
PUT  /api/admin/events/{id}/deactivate # Деактивировать
POST /api/admin/events/{id}/missions # Добавить миссию
```

## Автоматизация

```go
// Cron job для активации/деактивации событий
func (s *EventScheduler) Run() {
    ticker := time.NewTicker(1 * time.Hour)
    for range ticker.C {
        s.store.DB.ExecContext(ctx, `
            UPDATE seasonal_event
            SET is_active = (CURRENT_DATE BETWEEN start_date AND end_date)
        `)
    }
}
```

## Чек-лист

- [ ] Миграция `041_seasonal_events.up.sql`
- [ ] Store методы
- [ ] REST API endpoints
- [ ] Интеграция с gamification service
- [ ] WebSocket уведомления о миссиях
- [ ] Admin API для управления
- [ ] Scheduler для автоактивации
- [ ] Тематические ассеты (ТЗ дизайнеру)
- [ ] Создать первое событие (Новый год)

---

[← Mini Games](./19-mini-games.md) | [Назад к Roadmap](./roadmap.md)

---

# Завершение Roadmap

Поздравляем! Вы дошли до конца roadmap. После реализации всех фаз у вас будет:

- **Полноценный REST API** для Mini App
- **Богатая геймификация**: streak, achievements, pet, boss
- **Социальные функции**: parent portal, leaderboard, family quests
- **Продвинутые фичи**: voice, mini-games, seasonal events

## Дальнейшие шаги

1. Начните с [Фазы 0](./01-streak-system.md) — Quick Wins
2. Параллельно готовьте дизайн для Фазы 1
3. Регулярно тестируйте на staging
4. Собирайте обратную связь от пользователей
5. Итерируйте!
