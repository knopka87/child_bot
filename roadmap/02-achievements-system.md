# 02: Система достижений (Achievements)

> Фаза 0 | Приоритет: P0 | Сложность: Низкая | Срок: 3-4 дня

## Цель

Внедрить систему бейджей/достижений для мотивации и долгосрочного engagement. Достижения визуализируют прогресс и дают цели.

## Каталог достижений (MVP)

| ID | Название | Условие | Редкость |
|----|----------|---------|----------|
| `first_task` | Первые шаги | 1 задача | common |
| `tasks_5` | Меткий стрелок | 5 задач | common |
| `tasks_10` | Умник | 10 задач | common |
| `tasks_25` | Ракета | 25 задач | rare |
| `tasks_50` | Чемпион | 50 задач | rare |
| `tasks_100` | Легенда | 100 задач | legendary |
| `streak_7` | Огонь! | 7 дней streak | rare |
| `streak_14` | Двухнедельник | 14 дней streak | rare |
| `streak_30` | Месячный марафон | 30 дней streak | legendary |
| `no_hints` | Профессор | Решил без подсказок | rare |
| `night_owl` | Сова | Решил после 21:00 | common |
| `early_bird` | Ранняя пташка | Решил до 8:00 | common |
| `speed_demon` | Молния | 3 задачи за день | rare |
| `perfect_week` | Идеальная неделя | 7 дней, все верно | legendary |
| `comeback` | Возвращение | Вернулся после 7 дней | common |

## Миграция базы данных

```sql
-- migrations/028_achievements.up.sql

-- Каталог достижений
CREATE TABLE achievement (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    icon_key TEXT NOT NULL,           -- ключ для стикера/иконки
    rarity TEXT NOT NULL DEFAULT 'common', -- common, rare, legendary
    condition_type TEXT NOT NULL,      -- tasks_count, streak, time_based, etc.
    condition_value JSONB NOT NULL,    -- {"count": 5} или {"hour_min": 21}
    xp_reward INT NOT NULL DEFAULT 0,
    sort_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Разблокированные достижения пользователей
CREATE TABLE user_achievement (
    user_id BIGINT NOT NULL REFERENCES "user"(chat_id) ON DELETE CASCADE,
    achievement_id TEXT NOT NULL REFERENCES achievement(id),
    unlocked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    notified BOOLEAN NOT NULL DEFAULT false,
    PRIMARY KEY (user_id, achievement_id)
);

CREATE INDEX idx_user_achievement_user ON user_achievement(user_id);
CREATE INDEX idx_user_achievement_unlocked ON user_achievement(unlocked_at);

-- Счётчики для условий достижений
CREATE TABLE user_achievement_progress (
    user_id BIGINT NOT NULL REFERENCES "user"(chat_id) ON DELETE CASCADE,
    counter_type TEXT NOT NULL,        -- tasks_total, tasks_today, correct_streak
    counter_value INT NOT NULL DEFAULT 0,
    last_updated TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, counter_type)
);

-- Начальные достижения
INSERT INTO achievement (id, name, description, icon_key, rarity, condition_type, condition_value, xp_reward, sort_order) VALUES
('first_task', 'Первые шаги', 'Реши свою первую задачу', 'star', 'common', 'tasks_count', '{"count": 1}', 10, 1),
('tasks_5', 'Меткий стрелок', 'Реши 5 задач', 'target', 'common', 'tasks_count', '{"count": 5}', 25, 2),
('tasks_10', 'Умник', 'Реши 10 задач', 'brain', 'common', 'tasks_count', '{"count": 10}', 50, 3),
('tasks_25', 'Ракета', 'Реши 25 задач', 'rocket', 'rare', 'tasks_count', '{"count": 25}', 100, 4),
('tasks_50', 'Чемпион', 'Реши 50 задач', 'trophy', 'rare', 'tasks_count', '{"count": 50}', 200, 5),
('tasks_100', 'Легенда', 'Реши 100 задач', 'crown', 'legendary', 'tasks_count', '{"count": 100}', 500, 6),
('streak_7', 'Огонь!', '7 дней подряд', 'fire', 'rare', 'streak', '{"days": 7}', 100, 10),
('streak_14', 'Двухнедельник', '14 дней подряд', 'fire_double', 'rare', 'streak', '{"days": 14}', 200, 11),
('streak_30', 'Месячный марафон', '30 дней подряд', 'fire_max', 'legendary', 'streak', '{"days": 30}', 500, 12),
('no_hints', 'Профессор', 'Реши задачу без подсказок', 'graduation', 'rare', 'no_hints', '{}', 50, 20),
('night_owl', 'Сова', 'Реши задачу после 21:00', 'owl', 'common', 'time_based', '{"hour_min": 21, "hour_max": 24}', 15, 30),
('early_bird', 'Ранняя пташка', 'Реши задачу до 8:00', 'sun', 'common', 'time_based', '{"hour_min": 5, "hour_max": 8}', 15, 31),
('speed_demon', 'Молния', '3 задачи за один день', 'lightning', 'rare', 'daily_count', '{"count": 3}', 75, 40),
('comeback', 'Возвращение', 'Вернись после 7 дней отсутствия', 'comeback', 'common', 'comeback', '{"days": 7}', 25, 50);
```

```sql
-- migrations/028_achievements.down.sql
DROP TABLE IF EXISTS user_achievement_progress;
DROP TABLE IF EXISTS user_achievement;
DROP TABLE IF EXISTS achievement;
```

## Структуры данных

```go
// internal/store/achievement.go
package store

import (
    "context"
    "encoding/json"
    "time"
)

type Achievement struct {
    ID             string          `json:"id"`
    Name           string          `json:"name"`
    Description    string          `json:"description"`
    IconKey        string          `json:"icon_key"`
    Rarity         string          `json:"rarity"` // common, rare, legendary
    ConditionType  string          `json:"condition_type"`
    ConditionValue json.RawMessage `json:"condition_value"`
    XPReward       int             `json:"xp_reward"`
    SortOrder      int             `json:"sort_order"`
    IsActive       bool            `json:"is_active"`
}

type UserAchievement struct {
    UserID        int64     `json:"user_id"`
    AchievementID string    `json:"achievement_id"`
    UnlockedAt    time.Time `json:"unlocked_at"`
    Notified      bool      `json:"notified"`
}

type AchievementWithStatus struct {
    Achievement
    Unlocked   bool       `json:"unlocked"`
    UnlockedAt *time.Time `json:"unlocked_at,omitempty"`
}

type AchievementUnlock struct {
    Achievement Achievement
    IsNew       bool
}
```

## Store методы

```go
// internal/store/achievement.go

func (s *Store) GetAllAchievements(ctx context.Context) ([]Achievement, error) {
    rows, err := s.DB.QueryContext(ctx, `
        SELECT id, name, description, icon_key, rarity,
               condition_type, condition_value, xp_reward, sort_order
        FROM achievement
        WHERE is_active = true
        ORDER BY sort_order
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var achievements []Achievement
    for rows.Next() {
        var a Achievement
        if err := rows.Scan(
            &a.ID, &a.Name, &a.Description, &a.IconKey, &a.Rarity,
            &a.ConditionType, &a.ConditionValue, &a.XPReward, &a.SortOrder,
        ); err != nil {
            return nil, err
        }
        achievements = append(achievements, a)
    }
    return achievements, rows.Err()
}

func (s *Store) GetUserAchievements(ctx context.Context, userID int64) ([]AchievementWithStatus, error) {
    rows, err := s.DB.QueryContext(ctx, `
        SELECT a.id, a.name, a.description, a.icon_key, a.rarity,
               a.condition_type, a.condition_value, a.xp_reward, a.sort_order,
               ua.unlocked_at IS NOT NULL as unlocked,
               ua.unlocked_at
        FROM achievement a
        LEFT JOIN user_achievement ua ON a.id = ua.achievement_id AND ua.user_id = $1
        WHERE a.is_active = true
        ORDER BY a.sort_order
    `, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var achievements []AchievementWithStatus
    for rows.Next() {
        var a AchievementWithStatus
        var unlockedAt *time.Time
        if err := rows.Scan(
            &a.ID, &a.Name, &a.Description, &a.IconKey, &a.Rarity,
            &a.ConditionType, &a.ConditionValue, &a.XPReward, &a.SortOrder,
            &a.Unlocked, &unlockedAt,
        ); err != nil {
            return nil, err
        }
        a.UnlockedAt = unlockedAt
        achievements = append(achievements, a)
    }
    return achievements, rows.Err()
}

func (s *Store) UnlockAchievement(ctx context.Context, userID int64, achievementID string) (*AchievementUnlock, error) {
    // Проверяем, не разблокировано ли уже
    var exists bool
    err := s.DB.QueryRowContext(ctx, `
        SELECT EXISTS(
            SELECT 1 FROM user_achievement
            WHERE user_id = $1 AND achievement_id = $2
        )
    `, userID, achievementID).Scan(&exists)
    if err != nil {
        return nil, err
    }

    if exists {
        return &AchievementUnlock{IsNew: false}, nil
    }

    // Разблокируем
    _, err = s.DB.ExecContext(ctx, `
        INSERT INTO user_achievement (user_id, achievement_id)
        VALUES ($1, $2)
        ON CONFLICT DO NOTHING
    `, userID, achievementID)
    if err != nil {
        return nil, err
    }

    // Получаем данные достижения
    var a Achievement
    err = s.DB.QueryRowContext(ctx, `
        SELECT id, name, description, icon_key, rarity, xp_reward
        FROM achievement WHERE id = $1
    `, achievementID).Scan(&a.ID, &a.Name, &a.Description, &a.IconKey, &a.Rarity, &a.XPReward)
    if err != nil {
        return nil, err
    }

    return &AchievementUnlock{Achievement: a, IsNew: true}, nil
}

func (s *Store) IncrementCounter(ctx context.Context, userID int64, counterType string, delta int) (int, error) {
    var newValue int
    err := s.DB.QueryRowContext(ctx, `
        INSERT INTO user_achievement_progress (user_id, counter_type, counter_value, last_updated)
        VALUES ($1, $2, $3, NOW())
        ON CONFLICT (user_id, counter_type) DO UPDATE
        SET counter_value = user_achievement_progress.counter_value + $3,
            last_updated = NOW()
        RETURNING counter_value
    `, userID, counterType, delta).Scan(&newValue)
    return newValue, err
}

func (s *Store) GetCounter(ctx context.Context, userID int64, counterType string) (int, error) {
    var value int
    err := s.DB.QueryRowContext(ctx, `
        SELECT COALESCE(counter_value, 0)
        FROM user_achievement_progress
        WHERE user_id = $1 AND counter_type = $2
    `, userID, counterType).Scan(&value)
    if err == sql.ErrNoRows {
        return 0, nil
    }
    return value, err
}

func (s *Store) ResetDailyCounter(ctx context.Context, userID int64, counterType string) error {
    _, err := s.DB.ExecContext(ctx, `
        UPDATE user_achievement_progress
        SET counter_value = 0, last_updated = NOW()
        WHERE user_id = $1 AND counter_type = $2
    `, userID, counterType)
    return err
}
```

## Сервис проверки достижений

```go
// internal/service/achievement/checker.go
package achievement

import (
    "context"
    "encoding/json"
    "time"

    "child_bot/api/internal/store"
)

type Checker struct {
    store *store.Store
}

func NewChecker(s *store.Store) *Checker {
    return &Checker{store: s}
}

type TaskEvent struct {
    UserID       int64
    UsedHints    int
    IsCorrect    bool
    SolvedAt     time.Time
}

func (c *Checker) OnTaskSolved(ctx context.Context, event TaskEvent) ([]store.AchievementUnlock, error) {
    var unlocks []store.AchievementUnlock

    // 1. Увеличиваем счётчик задач
    totalTasks, err := c.store.IncrementCounter(ctx, event.UserID, "tasks_total", 1)
    if err != nil {
        return nil, err
    }

    // 2. Увеличиваем дневной счётчик
    dailyTasks, err := c.store.IncrementCounter(ctx, event.UserID, "tasks_today", 1)
    if err != nil {
        return nil, err
    }

    // 3. Проверяем достижения по количеству задач
    taskMilestones := map[int]string{
        1:   "first_task",
        5:   "tasks_5",
        10:  "tasks_10",
        25:  "tasks_25",
        50:  "tasks_50",
        100: "tasks_100",
    }

    if achievementID, ok := taskMilestones[totalTasks]; ok {
        unlock, err := c.store.UnlockAchievement(ctx, event.UserID, achievementID)
        if err != nil {
            return nil, err
        }
        if unlock.IsNew {
            unlocks = append(unlocks, *unlock)
        }
    }

    // 4. Проверяем "Молния" (3 задачи за день)
    if dailyTasks == 3 {
        unlock, err := c.store.UnlockAchievement(ctx, event.UserID, "speed_demon")
        if err != nil {
            return nil, err
        }
        if unlock.IsNew {
            unlocks = append(unlocks, *unlock)
        }
    }

    // 5. Проверяем "Профессор" (без подсказок)
    if event.UsedHints == 0 && event.IsCorrect {
        unlock, err := c.store.UnlockAchievement(ctx, event.UserID, "no_hints")
        if err != nil {
            return nil, err
        }
        if unlock.IsNew {
            unlocks = append(unlocks, *unlock)
        }
    }

    // 6. Проверяем время суток
    hour := event.SolvedAt.Hour()
    if hour >= 21 || hour < 5 {
        unlock, _ := c.store.UnlockAchievement(ctx, event.UserID, "night_owl")
        if unlock != nil && unlock.IsNew {
            unlocks = append(unlocks, *unlock)
        }
    }
    if hour >= 5 && hour < 8 {
        unlock, _ := c.store.UnlockAchievement(ctx, event.UserID, "early_bird")
        if unlock != nil && unlock.IsNew {
            unlocks = append(unlocks, *unlock)
        }
    }

    return unlocks, nil
}

func (c *Checker) OnStreakUpdate(ctx context.Context, userID int64, streak int) (*store.AchievementUnlock, error) {
    streakMilestones := map[int]string{
        7:  "streak_7",
        14: "streak_14",
        30: "streak_30",
    }

    if achievementID, ok := streakMilestones[streak]; ok {
        return c.store.UnlockAchievement(ctx, userID, achievementID)
    }

    return nil, nil
}

func (c *Checker) OnUserReturn(ctx context.Context, userID int64, daysSinceLastActivity int) (*store.AchievementUnlock, error) {
    if daysSinceLastActivity >= 7 {
        return c.store.UnlockAchievement(ctx, userID, "comeback")
    }
    return nil, nil
}
```

## Интеграция в Telegram Bot

```go
// internal/v2/telegram/achievements.go
package telegram

import (
    "context"
    "fmt"
    "strings"
    "time"

    "child_bot/api/internal/service/achievement"
    "child_bot/api/internal/store"
)

func (r *Router) checkAndNotifyAchievements(ctx context.Context, chatID int64, event achievement.TaskEvent) {
    unlocks, err := r.achievementChecker.OnTaskSolved(ctx, event)
    if err != nil {
        log.Printf("achievement check failed: %v", err)
        return
    }

    for _, unlock := range unlocks {
        r.notifyAchievement(chatID, unlock.Achievement)
    }
}

func (r *Router) notifyAchievement(chatID int64, a store.Achievement) {
    rarityEmoji := map[string]string{
        "common":    "⭐",
        "rare":      "🌟",
        "legendary": "💫",
    }

    emoji := rarityEmoji[a.Rarity]
    if emoji == "" {
        emoji = "🏅"
    }

    msg := fmt.Sprintf(
        "%s *Достижение разблокировано!*\n\n"+
        "*%s*\n"+
        "_%s_\n\n"+
        "+%d XP",
        emoji, a.Name, a.Description, a.XPReward,
    )

    // Отправляем с соответствующим стикером
    r.sendAchievementNotification(chatID, msg, a.IconKey)
}

func (r *Router) sendAchievementNotification(chatID int64, text string, iconKey string) {
    // Сначала стикер (если есть)
    if stickerID := r.getAchievementSticker(iconKey); stickerID != "" {
        r.bot.Send(tgbotapi.NewSticker(chatID, tgbotapi.FileID(stickerID)))
    }

    // Затем текст
    msg := tgbotapi.NewMessage(chatID, text)
    msg.ParseMode = "Markdown"
    r.bot.Send(msg)
}

func (r *Router) getAchievementSticker(iconKey string) string {
    // Маппинг iconKey -> Telegram sticker file_id
    // Заполняется после загрузки стикеров
    stickers := map[string]string{
        "star":       "", // TODO: добавить file_id
        "target":     "",
        "brain":      "",
        "rocket":     "",
        "trophy":     "",
        "crown":      "",
        "fire":       "",
        "graduation": "",
        "owl":        "",
        "sun":        "",
        "lightning":  "",
    }
    return stickers[iconKey]
}
```

## Команда просмотра достижений

```go
// internal/v2/telegram/commands.go
func (r *Router) handleAchievementsCommand(chatID int64) {
    achievements, err := r.store.GetUserAchievements(r.ctx, chatID)
    if err != nil {
        r.sendMessage(chatID, "Ошибка загрузки достижений")
        return
    }

    var unlocked, locked []string

    for _, a := range achievements {
        icon := r.rarityIcon(a.Rarity)
        if a.Unlocked {
            unlocked = append(unlocked, fmt.Sprintf("%s %s", icon, a.Name))
        } else {
            locked = append(locked, fmt.Sprintf("🔒 %s", a.Name))
        }
    }

    var sb strings.Builder
    sb.WriteString("🏆 *Твои достижения*\n\n")

    if len(unlocked) > 0 {
        sb.WriteString("*Разблокировано:*\n")
        for _, u := range unlocked {
            sb.WriteString(u + "\n")
        }
    }

    if len(locked) > 0 {
        sb.WriteString("\n*Ещё предстоит:*\n")
        for _, l := range locked {
            sb.WriteString(l + "\n")
        }
    }

    sb.WriteString(fmt.Sprintf("\n_%d из %d_", len(unlocked), len(achievements)))

    msg := tgbotapi.NewMessage(chatID, sb.String())
    msg.ParseMode = "Markdown"
    r.bot.Send(msg)
}

func (r *Router) rarityIcon(rarity string) string {
    switch rarity {
    case "legendary":
        return "💫"
    case "rare":
        return "🌟"
    default:
        return "⭐"
    }
}
```

## Сброс дневных счётчиков (Cron Job)

```go
// cmd/bot/main.go или отдельный worker

func startDailyCounterReset(store *store.Store) {
    ticker := time.NewTicker(1 * time.Hour)
    go func() {
        for range ticker.C {
            now := time.Now()
            // Сбрасываем в полночь по Москве
            if now.Hour() == 0 {
                ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
                store.ResetAllDailyCounters(ctx)
                cancel()
            }
        }
    }()
}

// internal/store/achievement.go
func (s *Store) ResetAllDailyCounters(ctx context.Context) error {
    _, err := s.DB.ExecContext(ctx, `
        UPDATE user_achievement_progress
        SET counter_value = 0, last_updated = NOW()
        WHERE counter_type LIKE '%_today'
    `)
    return err
}
```

## Тестирование

```go
// internal/service/achievement/checker_test.go
func TestOnTaskSolved(t *testing.T) {
    tests := []struct {
        name            string
        totalTasksBefore int
        event           TaskEvent
        wantAchievements []string
    }{
        {
            name:             "first task",
            totalTasksBefore: 0,
            event:            TaskEvent{UsedHints: 0, IsCorrect: true},
            wantAchievements: []string{"first_task", "no_hints"},
        },
        {
            name:             "fifth task with hints",
            totalTasksBefore: 4,
            event:            TaskEvent{UsedHints: 2, IsCorrect: true},
            wantAchievements: []string{"tasks_5"},
        },
        {
            name:             "night owl",
            totalTasksBefore: 10,
            event:            TaskEvent{UsedHints: 1, IsCorrect: true, SolvedAt: time.Date(2024, 1, 1, 22, 0, 0, 0, time.UTC)},
            wantAchievements: []string{"night_owl"},
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

- [ ] Создать миграцию `028_achievements.up.sql`
- [ ] Реализовать `store/achievement.go`
- [ ] Создать `service/achievement/checker.go`
- [ ] Интегрировать в `check.go` после верного ответа
- [ ] Связать со streak-системой (01-streak-system)
- [ ] Добавить команду `/achievements` или `/badges`
- [ ] Реализовать cron для сброса дневных счётчиков
- [ ] Написать unit-тесты
- [ ] Подготовить ТЗ на стикеры для дизайнера
- [ ] Тестирование на staging

## Связанные шаги

- [01-streak-system.md](./01-streak-system.md) — streak milestones триггерят achievements
- [03-daily-reports.md](./03-daily-reports.md) — новые достижения в отчёте

---

[← Streak-система](./01-streak-system.md) | [Назад к Roadmap](./roadmap.md) | [Далее: Daily Reports →](./03-daily-reports.md)
