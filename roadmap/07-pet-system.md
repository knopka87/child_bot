# 07: Питомец-компаньон (Pet System)

> Фаза 1 | Приоритет: P1 | Сложность: Низкая | Срок: 2-3 дня

## Цель

Создать систему виртуального питомца для эмоциональной привязки пользователя. Питомец "кормится" при решении задач.

## MVP Scope

- **1 вид питомца** (дракончик)
- **3 состояния**: hungry → fed → happy
- **Кормление** после каждой решённой задачи
- **Деградация** состояния со временем

## Миграция базы данных

```sql
-- migrations/031_pet_system.up.sql

CREATE TABLE user_pet (
    user_id BIGINT PRIMARY KEY REFERENCES "user"(chat_id) ON DELETE CASCADE,
    pet_type TEXT NOT NULL DEFAULT 'dragon',
    state TEXT NOT NULL DEFAULT 'hungry', -- hungry, fed, happy
    happiness INT NOT NULL DEFAULT 0,     -- 0-100
    last_fed_at TIMESTAMPTZ,
    evolution_stage INT NOT NULL DEFAULT 1, -- Фаза 2: 1-4
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Инициализация для существующих пользователей
INSERT INTO user_pet (user_id)
SELECT chat_id FROM "user"
ON CONFLICT DO NOTHING;
```

## Структуры данных

```go
// internal/store/pet.go
package store

import (
    "context"
    "database/sql"
    "time"
)

type PetState string

const (
    PetStateHungry PetState = "hungry"
    PetStateFed    PetState = "fed"
    PetStateHappy  PetState = "happy"
)

type UserPet struct {
    UserID         int64
    PetType        string
    State          PetState
    Happiness      int
    LastFedAt      sql.NullTime
    EvolutionStage int
    CreatedAt      time.Time
    UpdatedAt      time.Time
}

type FeedResult struct {
    PreviousState PetState
    NewState      PetState
    HappinessGain int
    NewHappiness  int
}
```

## Store методы

```go
// internal/store/pet.go

func (s *Store) GetUserPet(ctx context.Context, userID int64) (*UserPet, error) {
    var pet UserPet
    err := s.DB.QueryRowContext(ctx, `
        SELECT user_id, pet_type, state, happiness, last_fed_at,
               evolution_stage, created_at, updated_at
        FROM user_pet
        WHERE user_id = $1
    `, userID).Scan(
        &pet.UserID, &pet.PetType, &pet.State, &pet.Happiness,
        &pet.LastFedAt, &pet.EvolutionStage, &pet.CreatedAt, &pet.UpdatedAt,
    )
    if err == sql.ErrNoRows {
        // Create pet for new user
        return s.createPet(ctx, userID)
    }
    return &pet, err
}

func (s *Store) createPet(ctx context.Context, userID int64) (*UserPet, error) {
    pet := &UserPet{
        UserID:         userID,
        PetType:        "dragon",
        State:          PetStateHungry,
        Happiness:      0,
        EvolutionStage: 1,
        CreatedAt:      time.Now(),
        UpdatedAt:      time.Now(),
    }

    _, err := s.DB.ExecContext(ctx, `
        INSERT INTO user_pet (user_id, pet_type, state, happiness, evolution_stage)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (user_id) DO NOTHING
    `, userID, pet.PetType, pet.State, pet.Happiness, pet.EvolutionStage)

    return pet, err
}

func (s *Store) FeedPet(ctx context.Context, userID int64, happinessGain int) (*FeedResult, error) {
    tx, err := s.DB.BeginTx(ctx, nil)
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()

    var pet UserPet
    err = tx.QueryRowContext(ctx, `
        SELECT state, happiness FROM user_pet WHERE user_id = $1 FOR UPDATE
    `, userID).Scan(&pet.State, &pet.Happiness)
    if err == sql.ErrNoRows {
        // Create pet
        _, err = tx.ExecContext(ctx, `
            INSERT INTO user_pet (user_id, pet_type, state, happiness, last_fed_at)
            VALUES ($1, 'dragon', 'fed', $2, NOW())
        `, userID, happinessGain)
        if err != nil {
            return nil, err
        }
        tx.Commit()
        return &FeedResult{
            PreviousState: PetStateHungry,
            NewState:      PetStateFed,
            HappinessGain: happinessGain,
            NewHappiness:  happinessGain,
        }, nil
    }
    if err != nil {
        return nil, err
    }

    result := &FeedResult{
        PreviousState: pet.State,
        HappinessGain: happinessGain,
    }

    // Calculate new happiness
    newHappiness := pet.Happiness + happinessGain
    if newHappiness > 100 {
        newHappiness = 100
    }
    result.NewHappiness = newHappiness

    // Determine new state
    var newState PetState
    switch {
    case newHappiness >= 80:
        newState = PetStateHappy
    case newHappiness >= 30:
        newState = PetStateFed
    default:
        newState = PetStateHungry
    }
    result.NewState = newState

    _, err = tx.ExecContext(ctx, `
        UPDATE user_pet
        SET state = $2, happiness = $3, last_fed_at = NOW(), updated_at = NOW()
        WHERE user_id = $1
    `, userID, newState, newHappiness)
    if err != nil {
        return nil, err
    }

    return result, tx.Commit()
}

func (s *Store) DecayPetHappiness(ctx context.Context) (int, error) {
    // Decrease happiness by 5 for pets not fed in last 6 hours
    result, err := s.DB.ExecContext(ctx, `
        UPDATE user_pet
        SET happiness = GREATEST(0, happiness - 5),
            state = CASE
                WHEN happiness - 5 < 30 THEN 'hungry'
                WHEN happiness - 5 < 80 THEN 'fed'
                ELSE state
            END,
            updated_at = NOW()
        WHERE last_fed_at < NOW() - INTERVAL '6 hours'
          AND happiness > 0
    `)
    if err != nil {
        return 0, err
    }
    affected, _ := result.RowsAffected()
    return int(affected), nil
}
```

## Pet Service

```go
// internal/service/gamification/pet.go
package gamification

import (
    "context"

    "child_bot/api/internal/store"
)

type PetService struct {
    store *store.Store
}

func NewPetService(s *store.Store) *PetService {
    return &PetService{store: s}
}

// GetPet returns current pet state (with decay applied)
func (p *PetService) GetPet(ctx context.Context, userID int64) (*store.UserPet, error) {
    pet, err := p.store.GetUserPet(ctx, userID)
    if err != nil {
        return nil, err
    }

    // Apply decay if needed
    if pet.LastFedAt.Valid {
        hoursSinceFed := time.Since(pet.LastFedAt.Time).Hours()
        if hoursSinceFed > 6 {
            decayAmount := int(hoursSinceFed / 6) * 5
            if decayAmount > 0 && pet.Happiness > 0 {
                pet.Happiness = max(0, pet.Happiness-decayAmount)
                pet.State = p.calculateState(pet.Happiness)
            }
        }
    }

    return pet, nil
}

// Feed adds happiness to pet
func (p *PetService) Feed(ctx context.Context, userID int64) (*store.FeedResult, error) {
    happinessGain := 15 // базовое значение
    return p.store.FeedPet(ctx, userID, happinessGain)
}

// FeedBonus adds extra happiness (for streak, achievements)
func (p *PetService) FeedBonus(ctx context.Context, userID int64, bonus int) (*store.FeedResult, error) {
    return p.store.FeedPet(ctx, userID, bonus)
}

func (p *PetService) calculateState(happiness int) store.PetState {
    switch {
    case happiness >= 80:
        return store.PetStateHappy
    case happiness >= 30:
        return store.PetStateFed
    default:
        return store.PetStateHungry
    }
}
```

## REST API Handler

```go
// internal/api/handlers/pet.go
package handlers

import (
    "encoding/json"
    "net/http"

    "child_bot/api/internal/api/middleware"
    "child_bot/api/internal/service/gamification"
)

type PetHandler struct {
    petService *gamification.PetService
}

func NewPetHandler(svc *gamification.PetService) *PetHandler {
    return &PetHandler{petService: svc}
}

type PetResponse struct {
    Type           string `json:"type"`
    State          string `json:"state"`
    Happiness      int    `json:"happiness"`
    EvolutionStage int    `json:"evolution_stage"`
    LastFedAt      string `json:"last_fed_at,omitempty"`
}

func (h *PetHandler) GetPet(w http.ResponseWriter, r *http.Request) {
    user := middleware.GetUser(r.Context())

    pet, err := h.petService.GetPet(r.Context(), user.UserID)
    if err != nil {
        http.Error(w, `{"error": "failed to get pet"}`, http.StatusInternalServerError)
        return
    }

    resp := PetResponse{
        Type:           pet.PetType,
        State:          string(pet.State),
        Happiness:      pet.Happiness,
        EvolutionStage: pet.EvolutionStage,
    }
    if pet.LastFedAt.Valid {
        resp.LastFedAt = pet.LastFedAt.Time.Format(time.RFC3339)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

type FeedResponse struct {
    PreviousState string `json:"previous_state"`
    NewState      string `json:"new_state"`
    HappinessGain int    `json:"happiness_gain"`
    NewHappiness  int    `json:"new_happiness"`
    StateChanged  bool   `json:"state_changed"`
}

func (h *PetHandler) Feed(w http.ResponseWriter, r *http.Request) {
    user := middleware.GetUser(r.Context())

    result, err := h.petService.Feed(r.Context(), user.UserID)
    if err != nil {
        http.Error(w, `{"error": "failed to feed pet"}`, http.StatusInternalServerError)
        return
    }

    resp := FeedResponse{
        PreviousState: string(result.PreviousState),
        NewState:      string(result.NewState),
        HappinessGain: result.HappinessGain,
        NewHappiness:  result.NewHappiness,
        StateChanged:  result.PreviousState != result.NewState,
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
        // ... streak, XP ...

        // Feed pet
        result, err := s.petService.Feed(ctx, userID)
        if err == nil && result.NewState == store.PetStateHappy && result.PreviousState != store.PetStateHappy {
            // Pet became happy! Could trigger notification
            s.eventBus.Publish(events.Event{
                Type:   events.EventPetHappy,
                UserID: userID,
            })
        }
    }
}
```

## Decay Scheduler

```go
// cmd/bot/main.go

func startPetDecayScheduler(store *store.Store) {
    ticker := time.NewTicker(1 * time.Hour)
    go func() {
        for range ticker.C {
            ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
            affected, err := store.DecayPetHappiness(ctx)
            if err != nil {
                log.Printf("pet decay error: %v", err)
            } else if affected > 0 {
                log.Printf("decayed happiness for %d pets", affected)
            }
            cancel()
        }
    }()
}
```

## Telegram Bot Integration

```go
// internal/v2/telegram/pet.go
package telegram

func (r *Router) sendPetStatus(chatID int64) {
    pet, err := r.petService.GetPet(r.ctx, chatID)
    if err != nil {
        return
    }

    emoji := r.petStateEmoji(pet.State)
    bar := r.happinessBar(pet.Happiness)

    msg := fmt.Sprintf(
        "%s *Твой питомец*\n\n"+
        "Состояние: %s\n"+
        "Счастье: %s %d%%\n",
        emoji, r.petStateName(pet.State), bar, pet.Happiness,
    )

    if pet.State == store.PetStateHungry {
        msg += "\n_Покорми меня! Реши задачку_ 🥺"
    }

    r.sendMessage(chatID, msg)
}

func (r *Router) petStateEmoji(state store.PetState) string {
    switch state {
    case store.PetStateHappy:
        return "🐉✨"
    case store.PetStateFed:
        return "🐉"
    default:
        return "🐉😢"
    }
}

func (r *Router) petStateName(state store.PetState) string {
    switch state {
    case store.PetStateHappy:
        return "Счастлив!"
    case store.PetStateFed:
        return "Сытый"
    default:
        return "Голодный"
    }
}

func (r *Router) happinessBar(happiness int) string {
    filled := happiness / 10
    empty := 10 - filled
    return strings.Repeat("█", filled) + strings.Repeat("░", empty)
}

func (r *Router) notifyPetFed(chatID int64, result *store.FeedResult) {
    if result.PreviousState != result.NewState {
        switch result.NewState {
        case store.PetStateHappy:
            r.sendMessage(chatID, "🐉✨ Твой питомец счастлив!")
        case store.PetStateFed:
            r.sendMessage(chatID, "🐉 Питомец наелся!")
        }
    }
}
```

## Тестирование

```go
// internal/store/pet_test.go
func TestFeedPet(t *testing.T) {
    store := setupTestStore(t)
    ctx := context.Background()
    userID := int64(123)

    // Create user
    store.UpsertUser(ctx, store.UserUpsert{ChatID: userID})

    tests := []struct {
        name             string
        initialHappiness int
        gain             int
        wantState        store.PetState
    }{
        {
            name:             "hungry to fed",
            initialHappiness: 20,
            gain:             15,
            wantState:        store.PetStateFed,
        },
        {
            name:             "fed to happy",
            initialHappiness: 70,
            gain:             15,
            wantState:        store.PetStateHappy,
        },
        {
            name:             "cap at 100",
            initialHappiness: 95,
            gain:             15,
            wantState:        store.PetStateHappy,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Set initial state
            store.DB.ExecContext(ctx, `
                UPDATE user_pet SET happiness = $2 WHERE user_id = $1
            `, userID, tt.initialHappiness)

            result, err := store.FeedPet(ctx, userID, tt.gain)
            if err != nil {
                t.Fatalf("FeedPet error: %v", err)
            }

            if result.NewState != tt.wantState {
                t.Errorf("want state %s, got %s", tt.wantState, result.NewState)
            }
        })
    }
}
```

## Чек-лист

- [ ] Создать миграцию `031_pet_system.up.sql`
- [ ] Реализовать `store/pet.go`
- [ ] Создать `service/gamification/pet.go`
- [ ] Интегрировать с event bus (кормление при решении задачи)
- [ ] Добавить decay scheduler
- [ ] REST API endpoints (`GET /pet`, `POST /pet/feed`)
- [ ] Telegram команда `/pet`
- [ ] Написать unit-тесты
- [ ] Подготовить ТЗ на анимации питомца для дизайнера

## Связанные шаги

- [06-service-refactoring.md](./06-service-refactoring.md) — интеграция в gamification service
- [12-pet-evolution.md](./12-pet-evolution.md) — эволюция питомца (Фаза 2)

---

[← Service Refactoring](./06-service-refactoring.md) | [Назад к Roadmap](./roadmap.md) | [Далее: Boss System →](./08-boss-system.md)
