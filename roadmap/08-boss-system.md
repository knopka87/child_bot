# 08: Босс недели (Weekly Boss)

> Фаза 1 | Приоритет: P1 | Сложность: Средняя | Срок: 2-3 дня

## Цель

Создать систему еженедельных боссов для коллективной мотивации. Пользователи "атакуют" босса решая задачи.

## MVP Scope

- **3 типа боссов** в ротации
- **HP босса** уменьшается при решении задач
- **Еженедельная ротация** (понедельник)
- **Награды** за победу над боссом

## Типы боссов

| Тип | Название | Цвет | HP |
|-----|----------|------|-----|
| `dragon` | Дракон Забывания | Синий | 1000 |
| `error_king` | Король Ошибок | Красный | 1000 |
| `laziness` | Туман Лени | Серый | 1000 |

## Миграция базы данных

```sql
-- migrations/032_boss_system.up.sql

-- Типы боссов
CREATE TABLE boss_type (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    color TEXT,
    base_hp INT NOT NULL DEFAULT 1000,
    defeat_xp_reward INT NOT NULL DEFAULT 100,
    defeat_achievement_id TEXT REFERENCES achievement(id)
);

INSERT INTO boss_type (id, name, description, color, base_hp, defeat_xp_reward) VALUES
('dragon', 'Дракон Забывания', 'Заставляет забыть выученное', '#3B82F6', 1000, 100),
('error_king', 'Король Ошибок', 'Подсовывает ошибки в решения', '#EF4444', 1000, 100),
('laziness', 'Туман Лени', 'Мешает сосредоточиться', '#6B7280', 1000, 100);

-- Еженедельные боссы
CREATE TABLE weekly_boss (
    id SERIAL PRIMARY KEY,
    boss_type_id TEXT NOT NULL REFERENCES boss_type(id),
    week_start DATE NOT NULL,           -- Понедельник недели
    hp_total INT NOT NULL,              -- Начальное HP
    hp_remaining INT NOT NULL,          -- Текущее HP
    defeated_at TIMESTAMPTZ,            -- Когда побеждён
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(week_start)
);

CREATE INDEX idx_weekly_boss_week ON weekly_boss(week_start);

-- Прогресс пользователей
CREATE TABLE user_boss_progress (
    user_id BIGINT NOT NULL REFERENCES "user"(chat_id) ON DELETE CASCADE,
    boss_id INT NOT NULL REFERENCES weekly_boss(id) ON DELETE CASCADE,
    damage_dealt INT NOT NULL DEFAULT 0,
    attacks_count INT NOT NULL DEFAULT 0,
    reward_claimed BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, boss_id)
);

CREATE INDEX idx_user_boss_progress_boss ON user_boss_progress(boss_id);
```

## Структуры данных

```go
// internal/store/boss.go
package store

import (
    "context"
    "database/sql"
    "time"
)

type BossType struct {
    ID                   string
    Name                 string
    Description          string
    Color                string
    BaseHP               int
    DefeatXPReward       int
    DefeatAchievementID  *string
}

type WeeklyBoss struct {
    ID          int
    BossTypeID  string
    WeekStart   time.Time
    HPTotal     int
    HPRemaining int
    DefeatedAt  sql.NullTime
    CreatedAt   time.Time

    // Joined
    BossType *BossType
}

type UserBossProgress struct {
    UserID        int64
    BossID        int
    DamageDealt   int
    AttacksCount  int
    RewardClaimed bool
}

type AttackResult struct {
    DamageDealt    int
    BossHPBefore   int
    BossHPAfter    int
    BossDefeated   bool
    TotalUserDamage int
    UserRank       int // Позиция в топе атакующих
}
```

## Store методы

```go
// internal/store/boss.go

func (s *Store) GetCurrentBoss(ctx context.Context) (*WeeklyBoss, error) {
    // Получаем понедельник текущей недели
    now := time.Now()
    weekStart := now.AddDate(0, 0, -int(now.Weekday())+1) // Monday
    weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, time.UTC)

    var boss WeeklyBoss
    var bossType BossType

    err := s.DB.QueryRowContext(ctx, `
        SELECT wb.id, wb.boss_type_id, wb.week_start, wb.hp_total, wb.hp_remaining,
               wb.defeated_at, wb.created_at,
               bt.id, bt.name, bt.description, bt.color, bt.base_hp, bt.defeat_xp_reward
        FROM weekly_boss wb
        JOIN boss_type bt ON wb.boss_type_id = bt.id
        WHERE wb.week_start = $1
    `, weekStart).Scan(
        &boss.ID, &boss.BossTypeID, &boss.WeekStart, &boss.HPTotal, &boss.HPRemaining,
        &boss.DefeatedAt, &boss.CreatedAt,
        &bossType.ID, &bossType.Name, &bossType.Description, &bossType.Color,
        &bossType.BaseHP, &bossType.DefeatXPReward,
    )

    if err == sql.ErrNoRows {
        // Создаём нового босса
        return s.createWeeklyBoss(ctx, weekStart)
    }
    if err != nil {
        return nil, err
    }

    boss.BossType = &bossType
    return &boss, nil
}

func (s *Store) createWeeklyBoss(ctx context.Context, weekStart time.Time) (*WeeklyBoss, error) {
    // Выбираем случайного босса
    var bossType BossType
    err := s.DB.QueryRowContext(ctx, `
        SELECT id, name, description, color, base_hp, defeat_xp_reward
        FROM boss_type
        ORDER BY RANDOM()
        LIMIT 1
    `).Scan(&bossType.ID, &bossType.Name, &bossType.Description, &bossType.Color,
        &bossType.BaseHP, &bossType.DefeatXPReward)
    if err != nil {
        return nil, err
    }

    var bossID int
    err = s.DB.QueryRowContext(ctx, `
        INSERT INTO weekly_boss (boss_type_id, week_start, hp_total, hp_remaining)
        VALUES ($1, $2, $3, $3)
        ON CONFLICT (week_start) DO UPDATE SET boss_type_id = weekly_boss.boss_type_id
        RETURNING id
    `, bossType.ID, weekStart, bossType.BaseHP).Scan(&bossID)
    if err != nil {
        return nil, err
    }

    return &WeeklyBoss{
        ID:          bossID,
        BossTypeID:  bossType.ID,
        WeekStart:   weekStart,
        HPTotal:     bossType.BaseHP,
        HPRemaining: bossType.BaseHP,
        BossType:    &bossType,
    }, nil
}

func (s *Store) AttackBoss(ctx context.Context, userID int64, bossID int, damage int) (*AttackResult, error) {
    tx, err := s.DB.BeginTx(ctx, nil)
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()

    // Lock boss row
    var hpRemaining int
    var defeatedAt sql.NullTime
    err = tx.QueryRowContext(ctx, `
        SELECT hp_remaining, defeated_at FROM weekly_boss WHERE id = $1 FOR UPDATE
    `, bossID).Scan(&hpRemaining, &defeatedAt)
    if err != nil {
        return nil, err
    }

    result := &AttackResult{
        BossHPBefore: hpRemaining,
    }

    // Если босс уже побеждён, урон не наносим
    if defeatedAt.Valid {
        result.BossHPAfter = 0
        result.BossDefeated = true
        result.DamageDealt = 0
        return result, nil
    }

    // Наносим урон
    actualDamage := damage
    if actualDamage > hpRemaining {
        actualDamage = hpRemaining
    }
    newHP := hpRemaining - actualDamage

    result.DamageDealt = actualDamage
    result.BossHPAfter = newHP

    // Обновляем HP босса
    if newHP <= 0 {
        result.BossDefeated = true
        _, err = tx.ExecContext(ctx, `
            UPDATE weekly_boss SET hp_remaining = 0, defeated_at = NOW() WHERE id = $1
        `, bossID)
    } else {
        _, err = tx.ExecContext(ctx, `
            UPDATE weekly_boss SET hp_remaining = $2 WHERE id = $1
        `, bossID, newHP)
    }
    if err != nil {
        return nil, err
    }

    // Обновляем прогресс пользователя
    err = tx.QueryRowContext(ctx, `
        INSERT INTO user_boss_progress (user_id, boss_id, damage_dealt, attacks_count)
        VALUES ($1, $2, $3, 1)
        ON CONFLICT (user_id, boss_id) DO UPDATE
        SET damage_dealt = user_boss_progress.damage_dealt + $3,
            attacks_count = user_boss_progress.attacks_count + 1,
            updated_at = NOW()
        RETURNING damage_dealt
    `, userID, bossID, actualDamage).Scan(&result.TotalUserDamage)
    if err != nil {
        return nil, err
    }

    // Получаем ранг пользователя
    err = tx.QueryRowContext(ctx, `
        SELECT COUNT(*) + 1 FROM user_boss_progress
        WHERE boss_id = $1 AND damage_dealt > $2
    `, bossID, result.TotalUserDamage).Scan(&result.UserRank)
    if err != nil {
        result.UserRank = 0
    }

    return result, tx.Commit()
}

func (s *Store) GetBossLeaderboard(ctx context.Context, bossID int, limit int) ([]UserBossProgress, error) {
    rows, err := s.DB.QueryContext(ctx, `
        SELECT ubp.user_id, ubp.damage_dealt, ubp.attacks_count,
               c.first_name, c.username
        FROM user_boss_progress ubp
        JOIN chat c ON ubp.user_id = c.id
        WHERE ubp.boss_id = $1
        ORDER BY ubp.damage_dealt DESC
        LIMIT $2
    `, bossID, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var leaderboard []UserBossProgress
    for rows.Next() {
        var p UserBossProgress
        var firstName, username string
        if err := rows.Scan(&p.UserID, &p.DamageDealt, &p.AttacksCount, &firstName, &username); err != nil {
            return nil, err
        }
        leaderboard = append(leaderboard, p)
    }
    return leaderboard, rows.Err()
}

func (s *Store) GetUserBossProgress(ctx context.Context, userID int64, bossID int) (*UserBossProgress, error) {
    var p UserBossProgress
    err := s.DB.QueryRowContext(ctx, `
        SELECT user_id, boss_id, damage_dealt, attacks_count, reward_claimed
        FROM user_boss_progress
        WHERE user_id = $1 AND boss_id = $2
    `, userID, bossID).Scan(&p.UserID, &p.BossID, &p.DamageDealt, &p.AttacksCount, &p.RewardClaimed)
    if err == sql.ErrNoRows {
        return &UserBossProgress{UserID: userID, BossID: bossID}, nil
    }
    return &p, err
}
```

## Boss Service

```go
// internal/service/gamification/boss.go
package gamification

import (
    "context"

    "child_bot/api/internal/store"
)

type BossService struct {
    store *store.Store
}

func NewBossService(s *store.Store) *BossService {
    return &BossService{store: s}
}

func (b *BossService) GetCurrentBoss(ctx context.Context) (*store.WeeklyBoss, error) {
    return b.store.GetCurrentBoss(ctx)
}

func (b *BossService) Attack(ctx context.Context, userID int64, damage int) (*store.AttackResult, error) {
    boss, err := b.store.GetCurrentBoss(ctx)
    if err != nil {
        return nil, err
    }

    return b.store.AttackBoss(ctx, userID, boss.ID, damage)
}

func (b *BossService) GetUserProgress(ctx context.Context, userID int64) (*store.UserBossProgress, error) {
    boss, err := b.store.GetCurrentBoss(ctx)
    if err != nil {
        return nil, err
    }

    return b.store.GetUserBossProgress(ctx, userID, boss.ID)
}

func (b *BossService) GetLeaderboard(ctx context.Context, limit int) ([]store.UserBossProgress, error) {
    boss, err := b.store.GetCurrentBoss(ctx)
    if err != nil {
        return nil, err
    }

    return b.store.GetBossLeaderboard(ctx, boss.ID, limit)
}
```

## REST API Handler

```go
// internal/api/handlers/boss.go
package handlers

type BossHandler struct {
    bossService *gamification.BossService
}

type BossResponse struct {
    ID          int    `json:"id"`
    Type        string `json:"type"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Color       string `json:"color"`
    HPTotal     int    `json:"hp_total"`
    HPRemaining int    `json:"hp_remaining"`
    HPPercent   int    `json:"hp_percent"`
    IsDefeated  bool   `json:"is_defeated"`
    WeekStart   string `json:"week_start"`
}

func (h *BossHandler) GetCurrent(w http.ResponseWriter, r *http.Request) {
    boss, err := h.bossService.GetCurrentBoss(r.Context())
    if err != nil {
        http.Error(w, `{"error": "failed to get boss"}`, http.StatusInternalServerError)
        return
    }

    resp := BossResponse{
        ID:          boss.ID,
        Type:        boss.BossTypeID,
        Name:        boss.BossType.Name,
        Description: boss.BossType.Description,
        Color:       boss.BossType.Color,
        HPTotal:     boss.HPTotal,
        HPRemaining: boss.HPRemaining,
        HPPercent:   boss.HPRemaining * 100 / boss.HPTotal,
        IsDefeated:  boss.DefeatedAt.Valid,
        WeekStart:   boss.WeekStart.Format("2006-01-02"),
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

type AttackResponse struct {
    DamageDealt     int  `json:"damage_dealt"`
    BossHPRemaining int  `json:"boss_hp_remaining"`
    BossDefeated    bool `json:"boss_defeated"`
    TotalUserDamage int  `json:"total_user_damage"`
    UserRank        int  `json:"user_rank"`
}

func (h *BossHandler) Attack(w http.ResponseWriter, r *http.Request) {
    user := middleware.GetUser(r.Context())

    // Базовый урон за задачу
    damage := 10

    result, err := h.bossService.Attack(r.Context(), user.UserID, damage)
    if err != nil {
        http.Error(w, `{"error": "failed to attack boss"}`, http.StatusInternalServerError)
        return
    }

    resp := AttackResponse{
        DamageDealt:     result.DamageDealt,
        BossHPRemaining: result.BossHPAfter,
        BossDefeated:    result.BossDefeated,
        TotalUserDamage: result.TotalUserDamage,
        UserRank:        result.UserRank,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}
```

## Интеграция с Event Bus

```go
// internal/service/gamification/service.go

func (s *GamificationService) onTaskChecked(event events.Event) {
    ctx := context.Background()
    userID := event.UserID
    isCorrect := event.Payload["is_correct"].(bool)

    if isCorrect {
        // ... streak, XP, pet ...

        // Attack boss
        damage := 10
        result, err := s.bossService.Attack(ctx, userID, damage)
        if err == nil && result.BossDefeated {
            // Boss defeated! Award all participants
            s.awardBossDefeatRewards(ctx, result.BossID)
        }
    }
}
```

## Telegram Bot Integration

```go
// internal/v2/telegram/boss.go
package telegram

func (r *Router) sendBossStatus(chatID int64) {
    boss, err := r.bossService.GetCurrentBoss(r.ctx)
    if err != nil {
        return
    }

    hpPercent := boss.HPRemaining * 100 / boss.HPTotal
    bar := r.hpBar(hpPercent)

    userProgress, _ := r.bossService.GetUserProgress(r.ctx, chatID)

    msg := fmt.Sprintf(
        "👹 *Босс недели: %s*\n\n"+
        "HP: %s %d%%\n"+
        "Твой урон: ⚔️ %d\n",
        boss.BossType.Name, bar, hpPercent, userProgress.DamageDealt,
    )

    if boss.DefeatedAt.Valid {
        msg += "\n🎉 *ПОБЕЖДЁН!*"
    } else {
        msg += "\n_Решай задачки, чтобы атаковать!_"
    }

    r.sendMessage(chatID, msg)
}

func (r *Router) hpBar(percent int) string {
    filled := percent / 10
    empty := 10 - filled
    return "❤️" + strings.Repeat("█", filled) + strings.Repeat("░", empty)
}

func (r *Router) notifyBossAttack(chatID int64, result *store.AttackResult) {
    msg := fmt.Sprintf("⚔️ -%d урона боссу!", result.DamageDealt)

    if result.BossDefeated {
        msg += "\n\n🎉 *БОСС ПОБЕЖДЁН!* 🎉\nВсе участники получают награду!"
    } else if result.BossHPAfter < result.BossHPBefore/2 && result.BossHPBefore >= result.BossHPBefore/2 {
        msg += "\n\n⚡ Босс ослаблен! Осталось меньше 50% HP!"
    }

    r.sendMessage(chatID, msg)
}
```

## Тестирование

```go
// internal/store/boss_test.go
func TestAttackBoss(t *testing.T) {
    store := setupTestStore(t)
    ctx := context.Background()
    userID := int64(123)

    // Setup
    store.UpsertUser(ctx, store.UserUpsert{ChatID: userID})
    boss, _ := store.GetCurrentBoss(ctx)

    tests := []struct {
        name         string
        damage       int
        wantDefeated bool
    }{
        {
            name:         "normal attack",
            damage:       10,
            wantDefeated: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := store.AttackBoss(ctx, userID, boss.ID, tt.damage)
            if err != nil {
                t.Fatalf("AttackBoss error: %v", err)
            }

            if result.BossDefeated != tt.wantDefeated {
                t.Errorf("want defeated=%v, got %v", tt.wantDefeated, result.BossDefeated)
            }
        })
    }
}
```

## Чек-лист

- [ ] Создать миграцию `032_boss_system.up.sql`
- [ ] Реализовать `store/boss.go`
- [ ] Создать `service/gamification/boss.go`
- [ ] REST API endpoints (`GET /boss/current`, `POST /boss/attack`)
- [ ] Интегрировать с event bus (атака при решении задачи)
- [ ] Telegram команда `/boss`
- [ ] Еженедельная ротация (scheduler или on-demand)
- [ ] Награды за победу
- [ ] Unit-тесты
- [ ] Подготовить ТЗ на визуал боссов для дизайнера

## Связанные шаги

- [06-service-refactoring.md](./06-service-refactoring.md) — интеграция в gamification service
- [16-leaderboard.md](./16-leaderboard.md) — общий leaderboard

---

[← Pet System](./07-pet-system.md) | [Назад к Roadmap](./roadmap.md) | [Далее: WebSocket →](./09-websocket.md)
