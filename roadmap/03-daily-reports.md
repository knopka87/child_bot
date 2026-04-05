# 03: Ежедневные отчёты для родителей

> Фаза 0 | Приоритет: P1 | Сложность: Низкая | Срок: 2 дня

## Цель

Автоматическая отправка ежедневной статистики родителю. Повышает вовлечённость родителей и демонстрирует ценность продукта.

## Формат отчёта

```
📊 Отчёт за сегодня

✅ Решено задач: 3
📚 Темы: умножение, задачи на движение
⏱ Время занятий: ~15 мин
🔥 Серия: 5 дней подряд!
🏅 Новое достижение: "Меткий стрелок"

Так держать! 💪
```

## Миграция базы данных

```sql
-- migrations/029_daily_reports.up.sql

-- Настройки уведомлений пользователя
CREATE TABLE user_notification_settings (
    user_id BIGINT PRIMARY KEY REFERENCES "user"(chat_id) ON DELETE CASCADE,
    daily_report_enabled BOOLEAN NOT NULL DEFAULT true,
    daily_report_time TIME NOT NULL DEFAULT '20:00',
    timezone TEXT NOT NULL DEFAULT 'Europe/Moscow',
    parent_chat_id BIGINT,  -- куда отправлять отчёт (если отличается)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Лог отправленных отчётов
CREATE TABLE daily_report_log (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES "user"(chat_id),
    report_date DATE NOT NULL,
    sent_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    tasks_count INT,
    topics JSONB,
    streak INT,
    new_achievements JSONB,
    UNIQUE(user_id, report_date)
);

CREATE INDEX idx_daily_report_log_date ON daily_report_log(report_date);
```

## Структуры данных

```go
// internal/store/report.go
package store

import (
    "context"
    "database/sql"
    "encoding/json"
    "time"
)

type NotificationSettings struct {
    UserID             int64
    DailyReportEnabled bool
    DailyReportTime    string // "20:00"
    Timezone           string
    ParentChatID       *int64
}

type DailyReportData struct {
    UserID          int64
    ReportDate      time.Time
    TasksCount      int
    Topics          []string
    TotalTimeMinutes int
    Streak          int
    NewAchievements []string
}
```

## Store методы

```go
// internal/store/report.go

func (s *Store) GetUsersForDailyReport(ctx context.Context, currentHour int) ([]NotificationSettings, error) {
    rows, err := s.DB.QueryContext(ctx, `
        SELECT user_id, daily_report_enabled, daily_report_time, timezone, parent_chat_id
        FROM user_notification_settings
        WHERE daily_report_enabled = true
          AND EXTRACT(HOUR FROM daily_report_time) = $1
    `, currentHour)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var settings []NotificationSettings
    for rows.Next() {
        var s NotificationSettings
        if err := rows.Scan(&s.UserID, &s.DailyReportEnabled, &s.DailyReportTime, &s.Timezone, &s.ParentChatID); err != nil {
            return nil, err
        }
        settings = append(settings, s)
    }
    return settings, rows.Err()
}

func (s *Store) GetDailyStats(ctx context.Context, userID int64, date time.Time) (*DailyReportData, error) {
    startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
    endOfDay := startOfDay.Add(24 * time.Hour)

    report := &DailyReportData{
        UserID:     userID,
        ReportDate: date,
    }

    // Количество задач за день
    err := s.DB.QueryRowContext(ctx, `
        SELECT COUNT(DISTINCT task_session_id)
        FROM timeline_events
        WHERE chat_id = $1
          AND event_type = 'api_check'
          AND ok = true
          AND created_at >= $2
          AND created_at < $3
    `, userID, startOfDay, endOfDay).Scan(&report.TasksCount)
    if err != nil && err != sql.ErrNoRows {
        return nil, err
    }

    // Темы (subjects из parsed_tasks)
    rows, err := s.DB.QueryContext(ctx, `
        SELECT DISTINCT pt.subject
        FROM parsed_tasks pt
        JOIN task_sessions ts ON pt.session_id = ts.session_id
        WHERE ts.chat_id = $1
          AND ts.updated_at >= $2
          AND ts.updated_at < $3
          AND pt.subject IS NOT NULL
    `, userID, startOfDay, endOfDay)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var topic string
        if err := rows.Scan(&topic); err != nil {
            return nil, err
        }
        report.Topics = append(report.Topics, topic)
    }

    // Streak
    streak, err := s.GetUserStreak(ctx, userID)
    if err == nil && streak != nil {
        report.Streak = streak.CurrentStreak
    }

    // Новые достижения за день
    achievementRows, err := s.DB.QueryContext(ctx, `
        SELECT a.name
        FROM user_achievement ua
        JOIN achievement a ON ua.achievement_id = a.id
        WHERE ua.user_id = $1
          AND ua.unlocked_at >= $2
          AND ua.unlocked_at < $3
    `, userID, startOfDay, endOfDay)
    if err != nil {
        return nil, err
    }
    defer achievementRows.Close()

    for achievementRows.Next() {
        var name string
        if err := achievementRows.Scan(&name); err != nil {
            return nil, err
        }
        report.NewAchievements = append(report.NewAchievements, name)
    }

    return report, nil
}

func (s *Store) LogDailyReport(ctx context.Context, report *DailyReportData) error {
    topicsJSON, _ := json.Marshal(report.Topics)
    achievementsJSON, _ := json.Marshal(report.NewAchievements)

    _, err := s.DB.ExecContext(ctx, `
        INSERT INTO daily_report_log (user_id, report_date, tasks_count, topics, streak, new_achievements)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (user_id, report_date) DO NOTHING
    `, report.UserID, report.ReportDate, report.TasksCount, topicsJSON, report.Streak, achievementsJSON)
    return err
}

func (s *Store) WasReportSentToday(ctx context.Context, userID int64, date time.Time) (bool, error) {
    var exists bool
    err := s.DB.QueryRowContext(ctx, `
        SELECT EXISTS(
            SELECT 1 FROM daily_report_log
            WHERE user_id = $1 AND report_date = $2
        )
    `, userID, date.Format("2006-01-02")).Scan(&exists)
    return exists, err
}
```

## Сервис отправки отчётов

```go
// internal/service/report/sender.go
package report

import (
    "context"
    "fmt"
    "strings"
    "time"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "child_bot/api/internal/store"
)

type Sender struct {
    store *store.Store
    bot   *tgbotapi.BotAPI
}

func NewSender(s *store.Store, bot *tgbotapi.BotAPI) *Sender {
    return &Sender{store: s, bot: bot}
}

func (s *Sender) SendDailyReports(ctx context.Context) error {
    currentHour := time.Now().Hour()

    users, err := s.store.GetUsersForDailyReport(ctx, currentHour)
    if err != nil {
        return fmt.Errorf("get users: %w", err)
    }

    today := time.Now()

    for _, user := range users {
        // Проверяем, не отправляли ли уже
        sent, _ := s.store.WasReportSentToday(ctx, user.UserID, today)
        if sent {
            continue
        }

        stats, err := s.store.GetDailyStats(ctx, user.UserID, today)
        if err != nil {
            log.Printf("get stats for %d: %v", user.UserID, err)
            continue
        }

        // Не отправляем пустой отчёт
        if stats.TasksCount == 0 {
            continue
        }

        message := s.formatReport(stats)

        // Отправляем родителю или самому пользователю
        targetChatID := user.UserID
        if user.ParentChatID != nil {
            targetChatID = *user.ParentChatID
        }

        msg := tgbotapi.NewMessage(targetChatID, message)
        msg.ParseMode = "Markdown"

        if _, err := s.bot.Send(msg); err != nil {
            log.Printf("send report to %d: %v", targetChatID, err)
            continue
        }

        // Логируем отправку
        s.store.LogDailyReport(ctx, stats)
    }

    return nil
}

func (s *Sender) formatReport(stats *DailyReportData) string {
    var sb strings.Builder

    sb.WriteString("📊 *Отчёт за сегодня*\n\n")

    // Задачи
    sb.WriteString(fmt.Sprintf("✅ Решено задач: %d\n", stats.TasksCount))

    // Темы
    if len(stats.Topics) > 0 {
        sb.WriteString(fmt.Sprintf("📚 Темы: %s\n", strings.Join(stats.Topics, ", ")))
    }

    // Streak
    if stats.Streak > 0 {
        sb.WriteString(fmt.Sprintf("🔥 Серия: %d %s подряд!\n",
            stats.Streak, pluralizeDays(stats.Streak)))
    }

    // Достижения
    if len(stats.NewAchievements) > 0 {
        for _, a := range stats.NewAchievements {
            sb.WriteString(fmt.Sprintf("🏅 Новое достижение: \"%s\"\n", a))
        }
    }

    // Мотивационная концовка
    sb.WriteString("\n")
    if stats.TasksCount >= 5 {
        sb.WriteString("Отличная работа! 🚀")
    } else if stats.TasksCount >= 3 {
        sb.WriteString("Так держать! 💪")
    } else {
        sb.WriteString("Хорошее начало! ⭐")
    }

    return sb.String()
}

func pluralizeDays(n int) string {
    if n%10 == 1 && n%100 != 11 {
        return "день"
    }
    if n%10 >= 2 && n%10 <= 4 && (n%100 < 10 || n%100 >= 20) {
        return "дня"
    }
    return "дней"
}
```

## Scheduler (Cron Job)

```go
// cmd/bot/main.go

func startReportScheduler(store *store.Store, bot *tgbotapi.BotAPI) {
    sender := report.NewSender(store, bot)

    // Проверяем каждый час
    ticker := time.NewTicker(1 * time.Hour)
    go func() {
        for range ticker.C {
            ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
            if err := sender.SendDailyReports(ctx); err != nil {
                log.Printf("send daily reports: %v", err)
            }
            cancel()
        }
    }()
}

// В main():
func main() {
    // ... инициализация ...

    startReportScheduler(store, bot)

    // ... остальной код ...
}
```

## Команда настройки отчётов

```go
// internal/v2/telegram/commands.go

func (r *Router) handleReportSettingsCommand(chatID int64) {
    keyboard := tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("✅ Включить отчёты", "report_enable"),
            tgbotapi.NewInlineKeyboardButtonData("❌ Выключить", "report_disable"),
        ),
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("🕐 Время: 20:00", "report_time_20"),
            tgbotapi.NewInlineKeyboardButtonData("🕘 Время: 21:00", "report_time_21"),
        ),
    )

    msg := tgbotapi.NewMessage(chatID, "⚙️ *Настройки ежедневного отчёта*\n\nОтчёт будет приходить каждый день в выбранное время.")
    msg.ParseMode = "Markdown"
    msg.ReplyMarkup = keyboard
    r.bot.Send(msg)
}

func (r *Router) handleReportCallback(chatID int64, data string) {
    switch data {
    case "report_enable":
        r.store.UpdateNotificationSettings(r.ctx, chatID, true, "")
        r.sendMessage(chatID, "✅ Ежедневные отчёты включены")
    case "report_disable":
        r.store.UpdateNotificationSettings(r.ctx, chatID, false, "")
        r.sendMessage(chatID, "❌ Ежедневные отчёты выключены")
    case "report_time_20":
        r.store.UpdateNotificationSettings(r.ctx, chatID, true, "20:00")
        r.sendMessage(chatID, "🕐 Время отчёта: 20:00")
    case "report_time_21":
        r.store.UpdateNotificationSettings(r.ctx, chatID, true, "21:00")
        r.sendMessage(chatID, "🕘 Время отчёта: 21:00")
    }
}
```

## Тестирование

```go
// internal/service/report/sender_test.go
func TestFormatReport(t *testing.T) {
    sender := &Sender{}

    tests := []struct {
        name   string
        stats  *store.DailyReportData
        wantContains []string
    }{
        {
            name: "full report",
            stats: &store.DailyReportData{
                TasksCount:      5,
                Topics:          []string{"умножение", "деление"},
                Streak:          7,
                NewAchievements: []string{"Огонь!"},
            },
            wantContains: []string{
                "Решено задач: 5",
                "умножение, деление",
                "7 дней подряд",
                "Огонь!",
                "Отличная работа",
            },
        },
        {
            name: "minimal report",
            stats: &store.DailyReportData{
                TasksCount: 1,
                Topics:     []string{},
                Streak:     1,
            },
            wantContains: []string{
                "Решено задач: 1",
                "1 день подряд",
                "Хорошее начало",
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := sender.formatReport(tt.stats)
            for _, want := range tt.wantContains {
                if !strings.Contains(result, want) {
                    t.Errorf("expected %q in report, got:\n%s", want, result)
                }
            }
        })
    }
}
```

## Чек-лист

- [ ] Создать миграцию `029_daily_reports.up.sql`
- [ ] Реализовать `store/report.go`
- [ ] Создать `service/report/sender.go`
- [ ] Добавить scheduler в `main.go`
- [ ] Добавить команду `/reports` для настройки
- [ ] Обработать callbacks настроек
- [ ] Написать unit-тесты для форматирования
- [ ] Тестирование на staging
- [ ] Документировать для пользователей

## Связанные шаги

- [01-streak-system.md](./01-streak-system.md) — streak показывается в отчёте
- [02-achievements-system.md](./02-achievements-system.md) — новые достижения в отчёте
- [11-parent-child.md](./11-parent-child.md) — отправка отчёта родителю

---

[← Achievements](./02-achievements-system.md) | [Назад к Roadmap](./roadmap.md) | [Далее: API Layer →](./04-api-layer.md)
