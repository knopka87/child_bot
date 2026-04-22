# 01: Streak-система (Серия дней)

> Фаза 0 | Приоритет: P0 | Сложность: Низкая | Срок: 2-3 дня

## Цель

Внедрить механику "серии дней" для повышения ежедневного возврата пользователей. Streak — один из самых эффективных retention-инструментов (см. Duolingo).

## Бизнес-логика

### Правила streak

1. **Увеличение**: Streak +1 при решении хотя бы 1 задачи за день
2. **Сброс**: Streak = 0 если пользователь пропустил день
3. **Milestone**: Особые награды на 7, 14, 30, 100 дней
4. **Freeze** (Фаза 2): Возможность "заморозить" streak на 1 день

### Определение "дня"

- Используем таймзону пользователя (или UTC+3 Moscow по умолчанию)
- День считается с 00:00 до 23:59 локального времени
- Граница дня: `last_activity_date < today`

## Миграция базы данных

```sql
-- migrations/027_user_streak.up.sql
CREATE TABLE user_streak (
    user_id BIGINT PRIMARY KEY REFERENCES "user"(chat_id) ON DELETE CASCADE,
    current_streak INT NOT NULL DEFAULT 0,
    max_streak INT NOT NULL DEFAULT 0,
    last_activity_date DATE,
    streak_started_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_user_streak_last_activity ON user_streak(last_activity_date);

-- Инициализация для существующих пользователей
INSERT INTO user_streak (user_id, current_streak, max_streak)
SELECT chat_id, 0, 0 FROM "user"
ON CONFLICT DO NOTHING;
```

```sql
-- migrations/027_user_streak.down.sql
DROP TABLE IF EXISTS user_streak;
```

## Структуры данных

```go
// internal/store/streak.go
package store

import (
    "context"
    "database/sql"
    "time"
)

type UserStreak struct {
    UserID          int64
    CurrentStreak   int
    MaxStreak       int
    LastActivityDate sql.NullTime
    StreakStartedAt sql.NullTime
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

type StreakUpdate struct {
    PreviousStreak int
    CurrentStreak  int
    IsNewRecord    bool
    MilestoneHit   int // 0 если не milestone
}
```

## Store методы

```go
// internal/store/streak.go

func (s *Store) GetUserStreak(ctx context.Context, userID int64) (*UserStreak, error) {
    var streak UserStreak
    err := s.DB.QueryRowContext(ctx, `
        SELECT user_id, current_streak, max_streak,
               last_activity_date, streak_started_at, created_at, updated_at
        FROM user_streak
        WHERE user_id = $1
    `, userID).Scan(
        &streak.UserID, &streak.CurrentStreak, &streak.MaxStreak,
        &streak.LastActivityDate, &streak.StreakStartedAt,
        &streak.CreatedAt, &streak.UpdatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, nil
    }
    return &streak, err
}

func (s *Store) UpdateStreak(ctx context.Context, userID int64, today time.Time) (*StreakUpdate, error) {
    todayDate := today.Truncate(24 * time.Hour)
    yesterdayDate := todayDate.AddDate(0, 0, -1)

    tx, err := s.DB.BeginTx(ctx, nil)
    if err != nil {
        return nil, fmt.Errorf("begin tx: %w", err)
    }
    defer tx.Rollback()

    var streak UserStreak
    err = tx.QueryRowContext(ctx, `
        SELECT current_streak, max_streak, last_activity_date
        FROM user_streak
        WHERE user_id = $1
        FOR UPDATE
    `, userID).Scan(&streak.CurrentStreak, &streak.MaxStreak, &streak.LastActivityDate)

    if err == sql.ErrNoRows {
        // Создаём новую запись
        _, err = tx.ExecContext(ctx, `
            INSERT INTO user_streak (user_id, current_streak, max_streak, last_activity_date, streak_started_at)
            VALUES ($1, 1, 1, $2, $2)
        `, userID, todayDate)
        if err != nil {
            return nil, fmt.Errorf("insert streak: %w", err)
        }
        tx.Commit()
        return &StreakUpdate{
            PreviousStreak: 0,
            CurrentStreak:  1,
            IsNewRecord:    true,
            MilestoneHit:   0,
        }, nil
    }

    update := &StreakUpdate{PreviousStreak: streak.CurrentStreak}

    // Уже активен сегодня — ничего не меняем
    if streak.LastActivityDate.Valid &&
       streak.LastActivityDate.Time.Truncate(24*time.Hour).Equal(todayDate) {
        update.CurrentStreak = streak.CurrentStreak
        return update, nil
    }

    var newStreak int
    var streakStartedAt time.Time

    if streak.LastActivityDate.Valid &&
       streak.LastActivityDate.Time.Truncate(24*time.Hour).Equal(yesterdayDate) {
        // Продолжаем streak
        newStreak = streak.CurrentStreak + 1
        // streak_started_at не меняем
    } else {
        // Streak прерван, начинаем заново
        newStreak = 1
        streakStartedAt = todayDate
    }

    newMaxStreak := streak.MaxStreak
    if newStreak > streak.MaxStreak {
        newMaxStreak = newStreak
        update.IsNewRecord = true
    }

    // Проверяем milestone
    milestones := []int{7, 14, 30, 50, 100}
    for _, m := range milestones {
        if newStreak == m {
            update.MilestoneHit = m
            break
        }
    }

    if streakStartedAt.IsZero() {
        _, err = tx.ExecContext(ctx, `
            UPDATE user_streak
            SET current_streak = $2, max_streak = $3,
                last_activity_date = $4, updated_at = NOW()
            WHERE user_id = $1
        `, userID, newStreak, newMaxStreak, todayDate)
    } else {
        _, err = tx.ExecContext(ctx, `
            UPDATE user_streak
            SET current_streak = $2, max_streak = $3,
                last_activity_date = $4, streak_started_at = $5, updated_at = NOW()
            WHERE user_id = $1
        `, userID, newStreak, newMaxStreak, todayDate, streakStartedAt)
    }

    if err != nil {
        return nil, fmt.Errorf("update streak: %w", err)
    }

    update.CurrentStreak = newStreak

    return update, tx.Commit()
}

func (s *Store) CheckStreakBroken(ctx context.Context, userID int64, today time.Time) (bool, int, error) {
    todayDate := today.Truncate(24 * time.Hour)
    yesterdayDate := todayDate.AddDate(0, 0, -1)

    var lastActivity sql.NullTime
    var currentStreak int

    err := s.DB.QueryRowContext(ctx, `
        SELECT last_activity_date, current_streak
        FROM user_streak
        WHERE user_id = $1
    `, userID).Scan(&lastActivity, &currentStreak)

    if err == sql.ErrNoRows {
        return false, 0, nil
    }
    if err != nil {
        return false, 0, err
    }

    if !lastActivity.Valid {
        return false, 0, nil
    }

    lastDate := lastActivity.Time.Truncate(24 * time.Hour)

    // Streak сломан если последняя активность раньше вчера
    if lastDate.Before(yesterdayDate) && currentStreak > 0 {
        return true, currentStreak, nil
    }

    return false, currentStreak, nil
}
```

## Интеграция в Telegram Bot (v2)

```go
// internal/v2/telegram/streak.go
package telegram

import (
    "context"
    "fmt"
    "time"
)

func (r *Router) updateStreakOnTaskComplete(ctx context.Context, chatID int64) error {
    update, err := r.store.UpdateStreak(ctx, chatID, time.Now())
    if err != nil {
        return fmt.Errorf("update streak: %w", err)
    }

    // Формируем сообщение о streak
    var streakMsg string

    if update.CurrentStreak == 1 && update.PreviousStreak == 0 {
        streakMsg = "🔥 Начинаем серию! День 1"
    } else if update.CurrentStreak > update.PreviousStreak {
        streakMsg = fmt.Sprintf("🔥 Серия: %d %s подряд!",
            update.CurrentStreak, r.pluralizeDays(update.CurrentStreak))

        if update.IsNewRecord {
            streakMsg += " 🏆 Новый рекорд!"
        }

        if update.MilestoneHit > 0 {
            streakMsg += fmt.Sprintf("\n🎉 Достижение: %d дней!", update.MilestoneHit)
            // Триггерим achievement
            go r.unlockAchievement(ctx, chatID, fmt.Sprintf("streak_%d", update.MilestoneHit))
        }
    }

    if streakMsg != "" {
        r.sendStreakNotification(chatID, streakMsg)
    }

    return nil
}

func (r *Router) pluralizeDays(n int) string {
    if n%10 == 1 && n%100 != 11 {
        return "день"
    }
    if n%10 >= 2 && n%10 <= 4 && (n%100 < 10 || n%100 >= 20) {
        return "дня"
    }
    return "дней"
}

func (r *Router) sendStreakNotification(chatID int64, msg string) {
    // Отправляем как часть ответа или отдельным сообщением
    // В зависимости от контекста
}
```

## Интеграция в существующий flow

Добавить вызов в `hints.go` после успешного показа подсказки или в `check.go` после верного ответа:

```go
// internal/v2/telegram/check.go
func (r *Router) onCheckCorrect(ctx context.Context, chatID int64, sessionID string) {
    // ... существующая логика ...

    // Обновляем streak
    if err := r.updateStreakOnTaskComplete(ctx, chatID); err != nil {
        log.Printf("streak update failed: %v", err)
        // Не прерываем flow, streak — не критичен
    }
}
```

## Проверка потерянного streak при старте

```go
// internal/v2/telegram/router.go
func (r *Router) HandleUpdate(upd tgbotapi.Update) {
    chatID := extractChatID(upd)

    // Проверяем потерянный streak при первом сообщении за день
    if r.shouldCheckStreakLoss(chatID) {
        broken, oldStreak, _ := r.store.CheckStreakBroken(r.ctx, chatID, time.Now())
        if broken && oldStreak >= 3 {
            r.sendMessage(chatID, fmt.Sprintf(
                "💔 Ой, серия из %d дней прервалась...\nНо ничего! Начнём новую прямо сейчас! 🚀",
                oldStreak,
            ))
        }
    }

    // ... остальная логика ...
}

func (r *Router) shouldCheckStreakLoss(chatID int64) bool {
    // Используем TTL-кэш чтобы проверять только раз в день
    key := fmt.Sprintf("streak_check_%d_%s", chatID, time.Now().Format("2006-01-02"))
    if _, exists := r.streakCheckCache.Get(key); exists {
        return false
    }
    r.streakCheckCache.Set(key, true, 24*time.Hour)
    return true
}
```

## UI сообщения

### При увеличении streak
```
🔥 Серия: 5 дней подряд!
Ещё 2 дня до нового достижения 💪
```

### При достижении milestone
```
🔥🔥🔥 ОГОНЬ! 7 дней подряд! 🔥🔥🔥
🏆 Достижение разблокировано: "Недельный марафон"
```

### При потере streak
```
💔 Ой, серия из 5 дней прервалась...
Но ничего! Начнём новую прямо сейчас! 🚀
```

### При новом рекорде
```
🔥 Серия: 15 дней подряд!
🏆 Новый личный рекорд! Так держать!
```

## Тестирование

```go
// internal/store/streak_test.go
func TestUpdateStreak(t *testing.T) {
    tests := []struct {
        name           string
        lastActivity   time.Time
        currentStreak  int
        today          time.Time
        wantStreak     int
        wantNewRecord  bool
        wantMilestone  int
    }{
        {
            name:          "first activity ever",
            lastActivity:  time.Time{},
            currentStreak: 0,
            today:         time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
            wantStreak:    1,
            wantNewRecord: true,
        },
        {
            name:          "continue streak",
            lastActivity:  time.Date(2024, 1, 14, 18, 0, 0, 0, time.UTC),
            currentStreak: 5,
            today:         time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
            wantStreak:    6,
        },
        {
            name:          "streak broken",
            lastActivity:  time.Date(2024, 1, 12, 18, 0, 0, 0, time.UTC),
            currentStreak: 5,
            today:         time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
            wantStreak:    1,
            wantNewRecord: false,
        },
        {
            name:          "same day no change",
            lastActivity:  time.Date(2024, 1, 15, 8, 0, 0, 0, time.UTC),
            currentStreak: 5,
            today:         time.Date(2024, 1, 15, 18, 0, 0, 0, time.UTC),
            wantStreak:    5,
        },
        {
            name:          "milestone hit",
            lastActivity:  time.Date(2024, 1, 14, 18, 0, 0, 0, time.UTC),
            currentStreak: 6,
            today:         time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
            wantStreak:    7,
            wantMilestone: 7,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // ... test implementation
        })
    }
}
```

## Чек-лист

- [ ] Создать миграцию `027_user_streak.up.sql`
- [ ] Реализовать `store/streak.go`
- [ ] Добавить методы в `Store` interface
- [ ] Интегрировать в `check.go` (после правильного ответа)
- [ ] Добавить проверку потери streak при старте
- [ ] Написать unit-тесты
- [ ] Подготовить стикеры для дизайнера (передать ТЗ)
- [ ] Тестирование на staging

## Связанные шаги

- [02-achievements-system.md](./02-achievements-system.md) — streak milestones дают achievements
- [03-daily-reports.md](./03-daily-reports.md) — streak показывается в отчёте

---

[← Назад к Roadmap](./roadmap.md) | [Далее: Achievements →](./02-achievements-system.md)
