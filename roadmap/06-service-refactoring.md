# 06: Рефакторинг сервисного слоя

> Фаза 1 | Приоритет: P0 | Сложность: Высокая | Срок: 4-5 дней

## Цель

Выделить бизнес-логику из Telegram-специфичного кода в переиспользуемый сервисный слой. Это позволит использовать одну логику для Telegram Bot и Mini App API.

## Текущая проблема

```
# Сейчас: логика вшита в Telegram handlers
internal/v2/telegram/
├── parse.go      # PARSE API + Telegram UI + state management
├── hints.go      # HINT API + Telegram UI + caching
├── check.go      # CHECK API + Telegram UI + gamification
└── ...

# Проблема: дублирование при создании REST API
internal/api/handlers/
├── task.go       # Копировать логику из telegram/parse.go?
```

## Целевая архитектура

```
# После: чистое разделение
internal/
├── service/                    # Бизнес-логика (переиспользуемая)
│   ├── task/
│   │   ├── service.go          # TaskService interface + impl
│   │   ├── detect.go           # Detect logic
│   │   ├── parse.go            # Parse logic
│   │   ├── hint.go             # Hint logic (with caching)
│   │   └── check.go            # Check logic
│   ├── gamification/
│   │   ├── service.go          # GamificationService
│   │   ├── streak.go
│   │   ├── achievement.go
│   │   ├── pet.go
│   │   └── boss.go
│   ├── user/
│   │   └── service.go          # UserService
│   └── events/
│       └── bus.go              # Event bus for async processing
│
├── v2/telegram/                # Telegram-specific (UI only)
│   ├── router.go
│   ├── ui.go                   # Message formatting
│   └── handlers.go             # Thin handlers calling services
│
└── api/handlers/               # REST API (thin handlers)
    └── task.go                 # Calling same services
```

## Service Interfaces

```go
// internal/service/task/service.go
package task

import (
    "context"

    "child_bot/api/internal/store"
)

// TaskService handles task-related business logic
type TaskService interface {
    // Detect classifies the image and creates a session
    Detect(ctx context.Context, req DetectRequest) (*DetectResult, error)

    // Parse parses the task from the image
    Parse(ctx context.Context, req ParseRequest) (*ParseResult, error)

    // GetHint returns a hint for the task (from cache or generates new)
    GetHint(ctx context.Context, req HintRequest) (*HintResult, error)

    // CheckSolution verifies the student's solution
    CheckSolution(ctx context.Context, req CheckRequest) (*CheckResult, error)

    // GetSession returns session data
    GetSession(ctx context.Context, sessionID string) (*store.TaskSession, error)

    // GetHistory returns user's task history
    GetHistory(ctx context.Context, userID int64, limit, offset int) ([]TaskHistoryItem, int, error)
}

// DetectRequest for Detect method
type DetectRequest struct {
    UserID int64
    Image  []byte // raw image bytes
    Locale string
}

// DetectResult from Detect method
type DetectResult struct {
    SessionID         string
    Subject           string
    SubjectConfidence string
    Quality           Quality
    NeedsRetake       bool
}

type Quality struct {
    Status string
    Issues []string
}

// ParseRequest for Parse method
type ParseRequest struct {
    UserID    int64
    SessionID string
    // Image already stored in session
}

// ParseResult from Parse method
type ParseResult struct {
    Task       ParsedTask
    Items      []TaskItem
    HintPolicy HintPolicy
}

// HintRequest for GetHint method
type HintRequest struct {
    UserID    int64
    SessionID string
    Level     int // 1, 2, or 3
    Engine    string // "gemini" or "gpt"
}

// HintResult from GetHint method
type HintResult struct {
    Level    int
    HintText string
    HasMore  bool
    FromCache bool
}

// CheckRequest for CheckSolution method
type CheckRequest struct {
    UserID    int64
    SessionID string
    Image     []byte // solution image
}

// CheckResult from CheckSolution method
type CheckResult struct {
    Status     string // "correct", "incorrect", "cannot_evaluate"
    Feedback   string
    ErrorSpans []ErrorSpan
    IsCorrect  bool
}
```

## Service Implementation

```go
// internal/service/task/impl.go
package task

import (
    "context"
    "fmt"
    "time"

    "child_bot/api/internal/llmclient"
    "child_bot/api/internal/store"
)

type serviceImpl struct {
    store     *store.Store
    llmClient *llmclient.Client
    hintCache *HintCache
    eventBus  EventPublisher
}

func NewService(store *store.Store, llm *llmclient.Client, eventBus EventPublisher) TaskService {
    return &serviceImpl{
        store:     store,
        llmClient: llm,
        hintCache: NewHintCache(2 * time.Hour),
        eventBus:  eventBus,
    }
}

func (s *serviceImpl) Detect(ctx context.Context, req DetectRequest) (*DetectResult, error) {
    // 1. Call LLM server
    llmResp, err := s.llmClient.Detect(ctx, llmclient.DetectRequest{
        Image:  req.Image,
        Locale: req.Locale,
    })
    if err != nil {
        return nil, fmt.Errorf("llm detect: %w", err)
    }

    // 2. Create session
    sessionID, err := s.store.CreateTaskSession(ctx, req.UserID)
    if err != nil {
        return nil, fmt.Errorf("create session: %w", err)
    }

    // 3. Store image reference
    if err := s.store.StoreSessionImage(ctx, sessionID, req.Image); err != nil {
        return nil, fmt.Errorf("store image: %w", err)
    }

    // 4. Log event
    s.store.LogTimelineEvent(ctx, store.TimelineEvent{
        ChatID:        req.UserID,
        TaskSessionID: sessionID,
        EventType:     "api_detect",
        Direction:     "api",
        OK:            true,
    })

    return &DetectResult{
        SessionID:         sessionID,
        Subject:           llmResp.Classification.Subject,
        SubjectConfidence: llmResp.Classification.Confidence,
        Quality: Quality{
            Status: llmResp.Quality.Status,
            Issues: llmResp.Quality.Issues,
        },
        NeedsRetake: llmResp.Quality.Status == "recommend_retake",
    }, nil
}

func (s *serviceImpl) Parse(ctx context.Context, req ParseRequest) (*ParseResult, error) {
    // 1. Get session
    session, err := s.store.GetTaskSession(ctx, req.SessionID)
    if err != nil {
        return nil, fmt.Errorf("get session: %w", err)
    }

    // 2. Verify ownership
    if session.ChatID != req.UserID {
        return nil, fmt.Errorf("session not found")
    }

    // 3. Get user grade
    user, err := s.store.GetUser(ctx, req.UserID)
    if err != nil {
        return nil, fmt.Errorf("get user: %w", err)
    }

    // 4. Get stored image
    image, err := s.store.GetSessionImage(ctx, req.SessionID)
    if err != nil {
        return nil, fmt.Errorf("get image: %w", err)
    }

    // 5. Call LLM server
    llmResp, err := s.llmClient.Parse(ctx, llmclient.ParseRequest{
        Image: image,
        Grade: user.Grade,
        // ... other fields from session
    })
    if err != nil {
        return nil, fmt.Errorf("llm parse: %w", err)
    }

    // 6. Cache parsed task
    if err := s.store.CacheParsedTask(ctx, req.SessionID, llmResp); err != nil {
        // Non-critical, just log
    }

    // 7. Log event
    s.store.LogTimelineEvent(ctx, store.TimelineEvent{
        ChatID:        req.UserID,
        TaskSessionID: req.SessionID,
        EventType:     "api_parse",
        Direction:     "api",
        OK:            true,
    })

    return &ParseResult{
        Task:       mapToParseTask(llmResp.Task),
        Items:      mapToTaskItems(llmResp.Items),
        HintPolicy: mapToHintPolicy(llmResp.Items),
    }, nil
}

func (s *serviceImpl) GetHint(ctx context.Context, req HintRequest) (*HintResult, error) {
    // 1. Check cache first
    cacheKey := fmt.Sprintf("%s:%s:L%d", req.SessionID, req.Engine, req.Level)
    if cached, ok := s.hintCache.Get(cacheKey); ok {
        return &HintResult{
            Level:     req.Level,
            HintText:  cached.HintText,
            HasMore:   req.Level < 3,
            FromCache: true,
        }, nil
    }

    // 2. Check DB cache
    dbHint, err := s.store.GetCachedHint(ctx, req.SessionID, req.Engine, req.Level)
    if err == nil && dbHint != nil {
        return &HintResult{
            Level:     req.Level,
            HintText:  dbHint.HintText,
            HasMore:   req.Level < 3,
            FromCache: true,
        }, nil
    }

    // 3. Get parsed task from session
    parsedTask, err := s.store.GetParsedTask(ctx, req.SessionID)
    if err != nil {
        return nil, fmt.Errorf("get parsed task: %w", err)
    }

    // 4. Call LLM server
    llmResp, err := s.llmClient.Hint(ctx, llmclient.HintRequest{
        Task:   parsedTask,
        Level:  req.Level,
        Engine: req.Engine,
    })
    if err != nil {
        return nil, fmt.Errorf("llm hint: %w", err)
    }

    hint := extractHintByLevel(llmResp, req.Level)

    // 5. Cache result
    s.hintCache.Set(cacheKey, hint, 2*time.Hour)
    s.store.CacheHint(ctx, req.SessionID, req.Engine, req.Level, hint)

    // 6. Log event
    s.store.LogTimelineEvent(ctx, store.TimelineEvent{
        ChatID:        req.UserID,
        TaskSessionID: req.SessionID,
        EventType:     "api_hint",
        Direction:     "api",
        OK:            true,
    })

    return &HintResult{
        Level:     req.Level,
        HintText:  hint.HintText,
        HasMore:   req.Level < 3,
        FromCache: false,
    }, nil
}

func (s *serviceImpl) CheckSolution(ctx context.Context, req CheckRequest) (*CheckResult, error) {
    // 1. Get session and parsed task
    session, err := s.store.GetTaskSession(ctx, req.SessionID)
    if err != nil || session.ChatID != req.UserID {
        return nil, fmt.Errorf("session not found")
    }

    parsedTask, err := s.store.GetParsedTask(ctx, req.SessionID)
    if err != nil {
        return nil, fmt.Errorf("get parsed task: %w", err)
    }

    // 2. Call LLM server
    llmResp, err := s.llmClient.Check(ctx, llmclient.CheckRequest{
        Image:      req.Image,
        ParsedTask: parsedTask,
    })
    if err != nil {
        return nil, fmt.Errorf("llm check: %w", err)
    }

    // 3. Determine correctness
    isCorrect := llmResp.Decision == "correct"

    // 4. Publish event for gamification (async)
    s.eventBus.Publish(Event{
        Type:   EventTaskChecked,
        UserID: req.UserID,
        Payload: map[string]any{
            "session_id": req.SessionID,
            "is_correct": isCorrect,
            "hints_used": session.HintsUsed,
        },
    })

    // 5. Log event
    s.store.LogTimelineEvent(ctx, store.TimelineEvent{
        ChatID:        req.UserID,
        TaskSessionID: req.SessionID,
        EventType:     "api_check",
        Direction:     "api",
        OK:            isCorrect,
    })

    return &CheckResult{
        Status:     llmResp.Decision,
        Feedback:   llmResp.Feedback,
        ErrorSpans: mapErrorSpans(llmResp.ErrorSpans),
        IsCorrect:  isCorrect,
    }, nil
}
```

## Event Bus

```go
// internal/service/events/bus.go
package events

import (
    "sync"
)

type EventType string

const (
    EventTaskDetected EventType = "task.detected"
    EventTaskParsed   EventType = "task.parsed"
    EventHintUsed     EventType = "hint.used"
    EventTaskChecked  EventType = "task.checked"
    EventTaskCorrect  EventType = "task.correct"
    EventTaskIncorrect EventType = "task.incorrect"
)

type Event struct {
    Type    EventType
    UserID  int64
    Payload map[string]any
}

type Handler func(event Event)

type EventBus struct {
    handlers map[EventType][]Handler
    mu       sync.RWMutex
}

func NewEventBus() *EventBus {
    return &EventBus{
        handlers: make(map[EventType][]Handler),
    }
}

func (b *EventBus) Subscribe(eventType EventType, handler Handler) {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.handlers[eventType] = append(b.handlers[eventType], handler)
}

func (b *EventBus) Publish(event Event) {
    b.mu.RLock()
    handlers := b.handlers[event.Type]
    b.mu.RUnlock()

    // Fire handlers asynchronously
    for _, h := range handlers {
        go func(handler Handler) {
            defer func() {
                if r := recover(); r != nil {
                    log.Printf("event handler panic: %v", r)
                }
            }()
            handler(event)
        }(h)
    }
}
```

## Gamification Service

```go
// internal/service/gamification/service.go
package gamification

import (
    "context"

    "child_bot/api/internal/service/events"
    "child_bot/api/internal/store"
)

type GamificationService struct {
    store *store.Store
}

func NewService(s *store.Store, eventBus *events.EventBus) *GamificationService {
    svc := &GamificationService{store: s}

    // Subscribe to task events
    eventBus.Subscribe(events.EventTaskChecked, svc.onTaskChecked)

    return svc
}

func (s *GamificationService) onTaskChecked(event events.Event) {
    ctx := context.Background()
    userID := event.UserID
    isCorrect := event.Payload["is_correct"].(bool)
    hintsUsed := event.Payload["hints_used"].(int)

    if isCorrect {
        // Update streak
        s.updateStreak(ctx, userID)

        // Add XP
        xp := 25
        if hintsUsed == 0 {
            xp = 50 // Bonus for no hints
        }
        s.addXP(ctx, userID, xp)

        // Feed pet
        s.feedPet(ctx, userID)

        // Attack boss
        s.attackBoss(ctx, userID, 10)

        // Check achievements
        s.checkAchievements(ctx, userID, hintsUsed)
    }
}

func (s *GamificationService) updateStreak(ctx context.Context, userID int64) error {
    // Implementation from 01-streak-system.md
    return nil
}

func (s *GamificationService) addXP(ctx context.Context, userID int64, xp int) error {
    _, err := s.store.DB.ExecContext(ctx, `
        UPDATE "user" SET xp = xp + $2 WHERE chat_id = $1
    `, userID, xp)
    return err
}

// ... other methods
```

## Обновление Telegram Router

```go
// internal/v2/telegram/router.go
package telegram

import (
    "child_bot/api/internal/service/task"
    "child_bot/api/internal/service/gamification"
)

type Router struct {
    bot       *tgbotapi.BotAPI
    store     *store.Store

    // Services (NEW)
    taskService task.TaskService
    gamificationService *gamification.GamificationService
}

func NewRouter(bot *tgbotapi.BotAPI, store *store.Store, taskSvc task.TaskService, gamSvc *gamification.GamificationService) *Router {
    return &Router{
        bot:                 bot,
        store:               store,
        taskService:         taskSvc,
        gamificationService: gamSvc,
    }
}

// Thin handler - delegates to service
func (r *Router) handleParse(chatID int64, sessionID string) {
    result, err := r.taskService.Parse(r.ctx, task.ParseRequest{
        UserID:    chatID,
        SessionID: sessionID,
    })
    if err != nil {
        r.sendError(chatID, err)
        return
    }

    // Only UI logic here
    r.sendParsedTask(chatID, result)
}

func (r *Router) sendParsedTask(chatID int64, result *task.ParseResult) {
    msg := r.formatParsedTask(result)
    keyboard := r.buildHintKeyboard(result.HintPolicy)

    r.bot.Send(tgbotapi.NewMessage(chatID, msg))
    // ...
}
```

## Обновление REST API Handler

```go
// internal/api/handlers/task.go
package handlers

type TaskHandler struct {
    taskService task.TaskService
}

func NewTaskHandler(taskSvc task.TaskService) *TaskHandler {
    return &TaskHandler{taskService: taskSvc}
}

func (h *TaskHandler) Parse(w http.ResponseWriter, r *http.Request) {
    user := mw.GetUser(r.Context())

    var req dto.ParseRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        // ...
    }

    // Same service as Telegram bot!
    result, err := h.taskService.Parse(r.Context(), task.ParseRequest{
        UserID:    user.UserID,
        SessionID: req.SessionID,
    })
    if err != nil {
        // ...
    }

    // Only DTO mapping here
    resp := dto.ParseResponse{
        SessionID: req.SessionID,
        Task:      mapToDTO(result.Task),
        Items:     mapItemsToDTO(result.Items),
    }

    json.NewEncoder(w).Encode(resp)
}
```

## Инициализация в main.go

```go
// cmd/bot/main.go
func main() {
    cfg := config.Load()
    store := setupStore(cfg)
    llmClient := llmclient.New(cfg.LLMServerURL)

    // Event bus
    eventBus := events.NewEventBus()

    // Services
    taskService := task.NewService(store, llmClient, eventBus)
    gamificationService := gamification.NewService(store, eventBus)

    // Telegram bot
    bot, _ := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
    telegramRouter := telegram.NewRouter(bot, store, taskService, gamificationService)

    // REST API
    if cfg.APIEnabled {
        taskHandler := handlers.NewTaskHandler(taskService)
        apiRouter := api.NewRouter(cfg, taskHandler, ...)
        // ...
    }

    // ...
}
```

## Чек-лист

- [ ] Создать `internal/service/task/` с интерфейсами
- [ ] Вынести логику из `v2/telegram/parse.go` в `service/task/parse.go`
- [ ] Вынести логику из `v2/telegram/hints.go` в `service/task/hint.go`
- [ ] Вынести логику из `v2/telegram/check.go` в `service/task/check.go`
- [ ] Создать event bus `internal/service/events/`
- [ ] Создать `internal/service/gamification/`
- [ ] Обновить Telegram Router для использования сервисов
- [ ] Обновить REST handlers для использования сервисов
- [ ] Написать unit-тесты для сервисов
- [ ] Интеграционные тесты
- [ ] Убедиться, что Telegram bot работает как раньше

## Связанные шаги

- [04-api-layer.md](./04-api-layer.md) — использует сервисы
- [07-pet-system.md](./07-pet-system.md) — добавляется в gamification service
- [08-boss-system.md](./08-boss-system.md) — добавляется в gamification service

---

[← Auth System](./05-auth-system.md) | [Назад к Roadmap](./roadmap.md) | [Далее: Pet System →](./07-pet-system.md)
