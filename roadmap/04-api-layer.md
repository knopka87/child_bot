# 04: REST API Layer

> Фаза 1 | Приоритет: P0 | Сложность: Высокая | Срок: 5-7 дней

## Цель

Создать REST API слой для Mini App. Это фундамент для всех последующих шагов Фазы 1+.

## Архитектура

```
api/internal/
├── api/                      # NEW: HTTP API
│   ├── router.go             # Chi/Gin router setup
│   ├── middleware/
│   │   ├── auth.go           # JWT validation
│   │   ├── ratelimit.go      # Rate limiting
│   │   ├── cors.go           # CORS
│   │   ├── logging.go        # Request logging
│   │   └── recovery.go       # Panic recovery
│   ├── handlers/
│   │   ├── auth.go           # POST /auth/telegram
│   │   ├── user.go           # GET/PATCH /me
│   │   ├── task.go           # POST /tasks/detect, /tasks/parse, etc.
│   │   ├── gamification.go   # GET /streak, /achievements
│   │   ├── pet.go            # GET/POST /pet
│   │   ├── boss.go           # GET /boss
│   │   └── parent.go         # Parent portal endpoints
│   ├── dto/                  # Request/Response DTOs
│   │   ├── auth.go
│   │   ├── user.go
│   │   ├── task.go
│   │   └── gamification.go
│   └── errors/
│       └── errors.go         # API error types
```

## Зависимости

```go
// go.mod additions
require (
    github.com/go-chi/chi/v5 v5.0.12
    github.com/go-chi/cors v1.2.1
    github.com/golang-jwt/jwt/v5 v5.2.0
    golang.org/x/time v0.5.0  // rate limiting
)
```

## Router Setup

```go
// internal/api/router.go
package api

import (
    "net/http"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/cors"

    "child_bot/api/internal/api/handlers"
    mw "child_bot/api/internal/api/middleware"
    "child_bot/api/internal/store"
)

type Config struct {
    JWTSecret      string
    BotToken       string
    RateLimitRPS   int
    AllowedOrigins []string
}

func NewRouter(cfg Config, store *store.Store, llmClient *llmclient.Client) http.Handler {
    r := chi.NewRouter()

    // Global middleware
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(mw.Logger())
    r.Use(mw.Recoverer())
    r.Use(middleware.Timeout(60 * time.Second))

    // CORS
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   cfg.AllowedOrigins,
        AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
        ExposedHeaders:   []string{"X-Request-ID"},
        AllowCredentials: true,
        MaxAge:           300,
    }))

    // Rate limiting
    r.Use(mw.RateLimiter(cfg.RateLimitRPS))

    // Handlers
    authHandler := handlers.NewAuthHandler(store, cfg.BotToken, cfg.JWTSecret)
    userHandler := handlers.NewUserHandler(store)
    taskHandler := handlers.NewTaskHandler(store, llmClient)
    gamificationHandler := handlers.NewGamificationHandler(store)
    petHandler := handlers.NewPetHandler(store)
    bossHandler := handlers.NewBossHandler(store)

    // Public routes
    r.Post("/api/v1/auth/telegram", authHandler.AuthTelegram)

    // Health check
    r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("ok"))
    })

    // Protected routes
    r.Group(func(r chi.Router) {
        r.Use(mw.JWTAuth(cfg.JWTSecret))

        // User
        r.Get("/api/v1/me", userHandler.GetMe)
        r.Patch("/api/v1/me", userHandler.UpdateMe)
        r.Get("/api/v1/me/stats", userHandler.GetStats)

        // Tasks
        r.Post("/api/v1/tasks/detect", taskHandler.Detect)
        r.Post("/api/v1/tasks/parse", taskHandler.Parse)
        r.Post("/api/v1/tasks/{sessionID}/hint", taskHandler.GetHint)
        r.Post("/api/v1/tasks/{sessionID}/check", taskHandler.CheckSolution)
        r.Get("/api/v1/tasks/history", taskHandler.GetHistory)

        // Gamification
        r.Get("/api/v1/streak", gamificationHandler.GetStreak)
        r.Get("/api/v1/achievements", gamificationHandler.GetAchievements)

        // Pet
        r.Get("/api/v1/pet", petHandler.GetPet)
        r.Post("/api/v1/pet/feed", petHandler.Feed)

        // Boss
        r.Get("/api/v1/boss/current", bossHandler.GetCurrent)
        r.Post("/api/v1/boss/attack", bossHandler.Attack)
    })

    // Parent portal (separate auth)
    r.Group(func(r chi.Router) {
        r.Use(mw.JWTAuth(cfg.JWTSecret))
        r.Use(mw.RequireRole("parent"))

        r.Get("/api/v1/parent/children", parentHandler.GetChildren)
        r.Get("/api/v1/parent/child/{childID}/stats", parentHandler.GetChildStats)
        r.Post("/api/v1/parent/child/{childID}/praise", parentHandler.SendPraise)
    })

    return r
}
```

## Middleware

### JWT Authentication

```go
// internal/api/middleware/auth.go
package middleware

import (
    "context"
    "net/http"
    "strings"

    "github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserContextKey contextKey = "user"

type UserClaims struct {
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"` // "student" | "parent"
    jwt.RegisteredClaims
}

func JWTAuth(secret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, `{"error": "missing authorization header"}`, http.StatusUnauthorized)
                return
            }

            parts := strings.Split(authHeader, " ")
            if len(parts) != 2 || parts[0] != "Bearer" {
                http.Error(w, `{"error": "invalid authorization header"}`, http.StatusUnauthorized)
                return
            }

            tokenString := parts[1]

            token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
                return []byte(secret), nil
            })
            if err != nil || !token.Valid {
                http.Error(w, `{"error": "invalid token"}`, http.StatusUnauthorized)
                return
            }

            claims, ok := token.Claims.(*UserClaims)
            if !ok {
                http.Error(w, `{"error": "invalid token claims"}`, http.StatusUnauthorized)
                return
            }

            ctx := context.WithValue(r.Context(), UserContextKey, claims)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func GetUser(ctx context.Context) *UserClaims {
    user, _ := ctx.Value(UserContextKey).(*UserClaims)
    return user
}

func RequireRole(role string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user := GetUser(r.Context())
            if user == nil || user.Role != role {
                http.Error(w, `{"error": "forbidden"}`, http.StatusForbidden)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

### Rate Limiter

```go
// internal/api/middleware/ratelimit.go
package middleware

import (
    "net/http"
    "sync"

    "golang.org/x/time/rate"
)

type ipRateLimiter struct {
    ips map[string]*rate.Limiter
    mu  *sync.RWMutex
    r   rate.Limit
    b   int
}

func newIPRateLimiter(rps int) *ipRateLimiter {
    return &ipRateLimiter{
        ips: make(map[string]*rate.Limiter),
        mu:  &sync.RWMutex{},
        r:   rate.Limit(rps),
        b:   rps * 2,
    }
}

func (i *ipRateLimiter) getLimiter(ip string) *rate.Limiter {
    i.mu.Lock()
    defer i.mu.Unlock()

    limiter, exists := i.ips[ip]
    if !exists {
        limiter = rate.NewLimiter(i.r, i.b)
        i.ips[ip] = limiter
    }

    return limiter
}

func RateLimiter(rps int) func(http.Handler) http.Handler {
    limiter := newIPRateLimiter(rps)

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ip := r.RemoteAddr
            if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
                ip = forwarded
            }

            if !limiter.getLimiter(ip).Allow() {
                http.Error(w, `{"error": "too many requests"}`, http.StatusTooManyRequests)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

### Logging

```go
// internal/api/middleware/logging.go
package middleware

import (
    "log"
    "net/http"
    "time"

    "github.com/go-chi/chi/v5/middleware"
)

func Logger() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

            defer func() {
                log.Printf(
                    "%s %s %d %s %s",
                    r.Method,
                    r.URL.Path,
                    ww.Status(),
                    time.Since(start),
                    r.Header.Get("X-Request-ID"),
                )
            }()

            next.ServeHTTP(ww, r)
        })
    }
}
```

### Recovery

```go
// internal/api/middleware/recovery.go
package middleware

import (
    "log"
    "net/http"
    "runtime/debug"
)

func Recoverer() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            defer func() {
                if rvr := recover(); rvr != nil {
                    log.Printf("panic: %v\n%s", rvr, debug.Stack())
                    http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
                }
            }()

            next.ServeHTTP(w, r)
        })
    }
}
```

## DTOs

```go
// internal/api/dto/common.go
package dto

type ErrorResponse struct {
    Error   string `json:"error"`
    Code    string `json:"code,omitempty"`
    Details any    `json:"details,omitempty"`
}

type PaginationRequest struct {
    Limit  int `json:"limit"`
    Offset int `json:"offset"`
}

type PaginatedResponse[T any] struct {
    Data       []T  `json:"data"`
    Total      int  `json:"total"`
    HasMore    bool `json:"has_more"`
    NextOffset int  `json:"next_offset,omitempty"`
}
```

```go
// internal/api/dto/task.go
package dto

type DetectRequest struct {
    Image  string `json:"image" validate:"required"` // base64
    Locale string `json:"locale,omitempty"`
}

type DetectResponse struct {
    SessionID      string  `json:"session_id"`
    Subject        string  `json:"subject"`
    SubjectConfidence string `json:"subject_confidence"`
    Quality        Quality `json:"quality"`
}

type Quality struct {
    Status string   `json:"status"` // "ok" | "recommend_retake"
    Issues []string `json:"issues,omitempty"`
}

type ParseRequest struct {
    SessionID string `json:"session_id" validate:"required"`
}

type ParseResponse struct {
    SessionID  string     `json:"session_id"`
    Task       ParsedTask `json:"task"`
    Items      []TaskItem `json:"items"`
    HintPolicy HintPolicy `json:"hint_policy"`
}

type HintRequest struct {
    Level int `json:"level" validate:"min=1,max=3"` // 1, 2, or 3
}

type HintResponse struct {
    Level    int    `json:"level"`
    HintText string `json:"hint_text"`
    HasMore  bool   `json:"has_more"`
}

type CheckRequest struct {
    Image string `json:"image" validate:"required"` // base64 solution photo
}

type CheckResponse struct {
    Status     string      `json:"status"` // "correct" | "incorrect" | "cannot_evaluate"
    Feedback   string      `json:"feedback"`
    ErrorSpans []ErrorSpan `json:"error_spans,omitempty"`
}

type ErrorSpan struct {
    Start int    `json:"start"`
    End   int    `json:"end"`
    Text  string `json:"text"`
}
```

## Handlers

```go
// internal/api/handlers/task.go
package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/go-chi/chi/v5"

    "child_bot/api/internal/api/dto"
    mw "child_bot/api/internal/api/middleware"
    "child_bot/api/internal/llmclient"
    "child_bot/api/internal/store"
)

type TaskHandler struct {
    store     *store.Store
    llmClient *llmclient.Client
}

func NewTaskHandler(s *store.Store, llm *llmclient.Client) *TaskHandler {
    return &TaskHandler{store: s, llmClient: llm}
}

func (h *TaskHandler) Detect(w http.ResponseWriter, r *http.Request) {
    user := mw.GetUser(r.Context())
    if user == nil {
        http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
        return
    }

    var req dto.DetectRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
        return
    }

    // Вызываем LLM server
    result, err := h.llmClient.Detect(r.Context(), req.Image, req.Locale)
    if err != nil {
        http.Error(w, `{"error": "detection failed"}`, http.StatusInternalServerError)
        return
    }

    // Создаём сессию
    sessionID, err := h.store.CreateTaskSession(r.Context(), user.UserID)
    if err != nil {
        http.Error(w, `{"error": "session creation failed"}`, http.StatusInternalServerError)
        return
    }

    resp := dto.DetectResponse{
        SessionID:         sessionID,
        Subject:           result.Classification.Subject,
        SubjectConfidence: result.Classification.Confidence,
        Quality: dto.Quality{
            Status: result.Quality.Status,
            Issues: result.Quality.Issues,
        },
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func (h *TaskHandler) Parse(w http.ResponseWriter, r *http.Request) {
    user := mw.GetUser(r.Context())

    var req dto.ParseRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
        return
    }

    // Получаем сессию и проверяем ownership
    session, err := h.store.GetTaskSession(r.Context(), req.SessionID)
    if err != nil || session.ChatID != user.UserID {
        http.Error(w, `{"error": "session not found"}`, http.StatusNotFound)
        return
    }

    // Вызываем LLM server для парсинга
    // ... реализация аналогична v2/telegram/parse.go

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func (h *TaskHandler) GetHint(w http.ResponseWriter, r *http.Request) {
    user := mw.GetUser(r.Context())
    sessionID := chi.URLParam(r, "sessionID")

    var req dto.HintRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        req.Level = 1 // default
    }

    // Проверяем сессию
    session, err := h.store.GetTaskSession(r.Context(), sessionID)
    if err != nil || session.ChatID != user.UserID {
        http.Error(w, `{"error": "session not found"}`, http.StatusNotFound)
        return
    }

    // Получаем подсказку (из кэша или генерируем)
    // ... реализация аналогична v2/telegram/hints.go

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func (h *TaskHandler) CheckSolution(w http.ResponseWriter, r *http.Request) {
    user := mw.GetUser(r.Context())
    sessionID := chi.URLParam(r, "sessionID")

    var req dto.CheckRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
        return
    }

    // Проверяем сессию
    session, err := h.store.GetTaskSession(r.Context(), sessionID)
    if err != nil || session.ChatID != user.UserID {
        http.Error(w, `{"error": "session not found"}`, http.StatusNotFound)
        return
    }

    // Проверяем решение
    // ... реализация аналогична v2/telegram/check.go

    // Триггерим gamification events (async)
    go h.onTaskCompleted(r.Context(), user.UserID, sessionID, result.IsCorrect)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func (h *TaskHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
    user := mw.GetUser(r.Context())

    limit := 20
    offset := 0
    // Parse query params...

    sessions, total, err := h.store.GetUserTaskHistory(r.Context(), user.UserID, limit, offset)
    if err != nil {
        http.Error(w, `{"error": "failed to get history"}`, http.StatusInternalServerError)
        return
    }

    resp := dto.PaginatedResponse[dto.TaskHistoryItem]{
        Data:    sessions,
        Total:   total,
        HasMore: offset+limit < total,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func (h *TaskHandler) onTaskCompleted(ctx context.Context, userID int64, sessionID string, isCorrect bool) {
    // Streak update
    // Achievement check
    // Pet feed
    // Boss attack
    // etc.
}
```

## Интеграция в main.go

```go
// cmd/bot/main.go

func main() {
    cfg := config.Load()

    // ... existing code ...

    // API server (параллельно с Telegram bot)
    if cfg.APIEnabled {
        apiCfg := api.Config{
            JWTSecret:      cfg.JWTSecret,
            BotToken:       cfg.TelegramBotToken,
            RateLimitRPS:   cfg.RateLimitRPS,
            AllowedOrigins: cfg.CORSOrigins,
        }

        apiRouter := api.NewRouter(apiCfg, store, llmClient)

        apiServer := &http.Server{
            Addr:         ":" + cfg.APIPort,
            Handler:      apiRouter,
            ReadTimeout:  15 * time.Second,
            WriteTimeout: 60 * time.Second,
            IdleTimeout:  60 * time.Second,
        }

        go func() {
            log.Printf("API server starting on port %s", cfg.APIPort)
            if err := apiServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
                log.Fatalf("API server error: %v", err)
            }
        }()

        // Graceful shutdown for API server
        defer func() {
            ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
            defer cancel()
            apiServer.Shutdown(ctx)
        }()
    }

    // ... telegram bot code ...
}
```

## Config изменения

```go
// internal/config/config.go

type Config struct {
    // Existing...
    Port               string
    TelegramBotToken   string
    // ...

    // New API config
    APIEnabled     bool
    APIPort        string
    JWTSecret      string
    RateLimitRPS   int
    CORSOrigins    []string
}

func Load() *Config {
    return &Config{
        // Existing...

        APIEnabled:   os.Getenv("API_ENABLED") == "true",
        APIPort:      getEnvOrDefault("API_PORT", "8081"),
        JWTSecret:    os.Getenv("JWT_SECRET"),
        RateLimitRPS: getEnvInt("RATE_LIMIT_RPS", 10),
        CORSOrigins:  strings.Split(os.Getenv("CORS_ORIGINS"), ","),
    }
}
```

## Environment Variables

```env
# .env additions
API_ENABLED=true
API_PORT=8081
JWT_SECRET=your-secret-key-min-32-chars-long
RATE_LIMIT_RPS=10
CORS_ORIGINS=https://your-mini-app.com,http://localhost:3000
```

## Тестирование

```go
// internal/api/handlers/task_test.go
func TestDetect(t *testing.T) {
    // Setup test server
    store := setupTestStore(t)
    llm := setupMockLLM(t)
    handler := NewTaskHandler(store, llm)

    router := chi.NewRouter()
    router.Use(mw.JWTAuth("test-secret"))
    router.Post("/api/v1/tasks/detect", handler.Detect)

    // Test cases
    tests := []struct {
        name       string
        token      string
        body       string
        wantStatus int
    }{
        {
            name:       "valid request",
            token:      validToken,
            body:       `{"image": "base64..."}`,
            wantStatus: http.StatusOK,
        },
        {
            name:       "missing auth",
            token:      "",
            body:       `{"image": "base64..."}`,
            wantStatus: http.StatusUnauthorized,
        },
        {
            name:       "invalid body",
            token:      validToken,
            body:       `{invalid}`,
            wantStatus: http.StatusBadRequest,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest("POST", "/api/v1/tasks/detect", strings.NewReader(tt.body))
            if tt.token != "" {
                req.Header.Set("Authorization", "Bearer "+tt.token)
            }
            req.Header.Set("Content-Type", "application/json")

            rec := httptest.NewRecorder()
            router.ServeHTTP(rec, req)

            if rec.Code != tt.wantStatus {
                t.Errorf("want status %d, got %d", tt.wantStatus, rec.Code)
            }
        })
    }
}
```

## Чек-лист

- [ ] Добавить зависимости в go.mod
- [ ] Создать структуру папок `internal/api/`
- [ ] Реализовать middleware (auth, ratelimit, logging, recovery)
- [ ] Реализовать DTOs
- [ ] Реализовать handlers (auth, user, task)
- [ ] Интегрировать в main.go
- [ ] Добавить config для API
- [ ] Написать unit-тесты для handlers
- [ ] Написать integration tests
- [ ] Документировать API (OpenAPI/Swagger)
- [ ] Тестирование на staging

## Связанные шаги

- [05-auth-system.md](./05-auth-system.md) — детали авторизации
- [06-service-refactoring.md](./06-service-refactoring.md) — выделение бизнес-логики

---

[← Daily Reports](./03-daily-reports.md) | [Назад к Roadmap](./roadmap.md) | [Далее: Auth System →](./05-auth-system.md)
