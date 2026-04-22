# Backend Migration — Миграция с Telegram Bot на REST API

**Дата:** 2026-03-29
**Версия:** 1.0
**Цель:** Извлечь бизнес-логику из Telegram бота и создать REST API для миниапп (VK, Max, Telegram)

---

## 🎯 Обзор

Текущий backend построен как **Telegram Bot**, который:
- Получает updates от Telegram API
- Обрабатывает фото и текст
- Общается с LLM-сервером
- Хранит сессии в PostgreSQL

**Задача:** Извлечь всю бизнес-логику и создать REST API, убрав зависимость от Telegram Bot API.

---

## 📊 Текущая архитектура

### Структура проекта

```
api/
├── cmd/bot/main.go                    # Точка входа Telegram бота
├── internal/
│   ├── v2/
│   │   ├── llmclient/client.go       # HTTP клиент для LLM-сервера ✅ ПЕРЕИСПОЛЬЗУЕМ
│   │   ├── types/                     # TypeScript-подобные типы ✅ ПЕРЕИСПОЛЬЗУЕМ
│   │   │   ├── hint.go               # HintRequest/Response
│   │   │   ├── check.go              # CheckRequest/Response
│   │   │   ├── analogue.go           # AnalogueRequest/Response
│   │   │   ├── detect.go             # DetectRequest/Response
│   │   │   └── parse.go              # ParseRequest/Response
│   │   └── telegram/                  # ❌ УДАЛИТЬ (Telegram-специфичный код)
│   │       ├── session.go            # → Извлечь Session Management
│   │       ├── photo.go              # → Извлечь Photo Processing
│   │       ├── callbacks.go          # → Убрать (UI-специфично)
│   │       ├── ui.go                 # → Убрать (UI-специфично)
│   │       ├── check.go              # → Извлечь Check flow
│   │       ├── analogue.go           # → Извлечь Analogue flow
│   │       └── hints.go              # → Извлечь Hints flow
│   ├── store/                         # PostgreSQL store ✅ ПЕРЕИСПОЛЬЗУЕМ
│   │   ├── session.go                # Сессии
│   │   ├── parse.go                  # История парсинга
│   │   ├── user.go                   # Пользователи
│   │   ├── hint.go                   # Подсказки
│   │   └── textbook_search.go        # Поиск по учебникам
│   └── config/                        # Конфиг ✅ ПЕРЕИСПОЛЬЗУЕМ
```

---

## 🔄 LLM-сервер (существующий)

LLM-сервер уже работает как **независимый HTTP API**. Frontend миниаппа будет общаться с ним через backend proxy.

### Существующие endpoints LLM-сервера

| Endpoint | Метод | Описание |
|----------|-------|----------|
| `/v2/detect` | POST | Определение задания на фото (subject, grade) |
| `/v2/parse` | POST | Парсинг текста задания из фото |
| `/v2/hint` | POST | Генерация подсказок (L1, L2, L3) |
| `/v2/check_solution` | POST | Проверка решения ученика |
| `/v2/analogue_solution` | POST | Генерация аналогичного задания |

**Формат запросов:**
```json
{
  "llm_name": "gpt-4o-mini",
  "image": "base64...",
  "task": {...},
  "options": {...}
}
```

**HTTP Client уже реализован:** `api/internal/v2/llmclient/client.go`

---

## 🗄️ База данных (PostgreSQL)

### Существующие таблицы

| Таблица | Назначение | Использование |
|---------|------------|---------------|
| `sessions` | Хранение сессий пользователей | ✅ Переиспользуем |
| `parse_history` | История парсинга заданий | ✅ Переиспользуем |
| `users` | Профили пользователей | ✅ Переиспользуем |
| `textbook_search` | Кэш поиска по учебникам | ✅ Переиспользуем |
| `hint_context` | Контекст подсказок | ✅ Переиспользуем |

**Store interface:** `api/internal/store/store.go`

```go
type Store interface {
    // Session management
    FindSession(ctx context.Context, chatID int64) (*Session, error)
    UpsertSessionID(ctx context.Context, chatID int64, sessionID string) error

    // Parse history
    SaveParseHistory(ctx context.Context, h *ParseHistory) error
    GetParseHistory(ctx context.Context, sessionID string) (*ParseHistory, error)

    // User profile
    GetUser(ctx context.Context, userID int64) (*User, error)
    UpsertUser(ctx context.Context, u *User) error
}
```

---

## ❌ Что нужно удалить

### 1. Telegram Bot Dependencies

```go
// cmd/bot/main.go
import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken) // ❌ Удалить
```

### 2. Telegram-специфичный код

**Удалить целиком:**
- `api/internal/v2/telegram/ui.go` - формирование Telegram кнопок
- `api/internal/v2/telegram/callbacks.go` - обработка callback_query
- Все references на `tgbotapi.*`

---

## ✅ Что нужно извлечь и переписать

### 1. **Session Management** (из `telegram/session.go`)

**Текущая реализация:**
```go
// Telegram-специфично: chatID как ключ
func (r *Router) setSession(cid int64, sid string) {
    sessionByChat.Store(cid, sid)
    _ = r.Store.UpsertSessionID(context.Background(), cid, sid)
}
```

**Новая реализация для REST API:**
```go
// Platform-agnostic: userID + platform
type SessionManager struct {
    cache *TTLCache
    store Store
}

func (sm *SessionManager) GetSession(ctx context.Context, userID string) (*Session, error) {
    // 1. Try cache
    if cached := sm.cache.Get(userID); cached != nil {
        return cached.(*Session), nil
    }

    // 2. Try database
    session, err := sm.store.FindSessionByUserID(ctx, userID)
    if err == nil {
        sm.cache.Set(userID, session)
        return session, nil
    }

    // 3. Create new
    session = &Session{
        ID:        uuid.NewString(),
        UserID:    userID,
        CreatedAt: time.Now(),
    }
    sm.store.UpsertSession(ctx, session)
    sm.cache.Set(userID, session)
    return session, nil
}
```

### 2. **Photo Processing** (из `telegram/photo.go`)

**Текущая реализация:**
- Download фото через Telegram API
- Debounce (сборка батча фото)
- Склейка множества фото в одно

**Новая реализация:**
```go
type PhotoProcessor struct {
    maxSizeBytes int
    maxPixels    int
}

func (pp *PhotoProcessor) Process(images [][]byte) ([]byte, error) {
    if len(images) == 0 {
        return nil, errors.New("no images")
    }

    if len(images) == 1 {
        return images[0], nil
    }

    // Combine multiple images vertically
    return pp.combineVertically(images)
}

func (pp *PhotoProcessor) combineVertically(images [][]byte) ([]byte, error) {
    // Decode all images
    // Find max width, sum heights
    // Create canvas
    // Draw images vertically centered
    // Resize if needed
    // Encode as JPEG
}
```

### 3. **Detect + Parse Flow** (из `telegram/photo.go`)

**Текущая реализация:**
```go
func (r *Router) runDetectThenParse(ctx context.Context, chatID int64, userID *int64, photo []byte, mediaGroupID string) {
    // 1. Detect subject + grade
    detectResp, err := r.LLMClient.Detect(ctx, llmName, detectReq)

    // 2. Parse task structure
    parseResp, err := r.LLMClient.Parse(ctx, llmName, parseReq)

    // 3. Save to database
    r.Store.SaveParseHistory(ctx, history)

    // 4. Show UI (Telegram-specific)
    r.showTaskUI(chatID, parseResp)
}
```

**Новая реализация (REST API):**
```go
// POST /api/v1/tasks/upload
func (h *Handler) UploadTask(w http.ResponseWriter, r *http.Request) {
    // 1. Parse multipart/form-data (images)
    files := r.MultipartForm.File["images"]
    images := [][]byte{}
    for _, fh := range files {
        data, _ := io.ReadAll(fh.Open())
        images = append(images, data)
    }

    // 2. Process images (combine if multiple)
    photo, err := h.photoProcessor.Process(images)

    // 3. Get user from JWT
    userID := getUserIDFromContext(r.Context())

    // 4. Detect + Parse
    detectResp, _ := h.llmClient.Detect(ctx, llmName, detectReq)
    parseResp, _ := h.llmClient.Parse(ctx, llmName, parseReq)

    // 5. Create attempt
    attempt := &Attempt{
        ID:        uuid.NewString(),
        UserID:    userID,
        Photo:     photo,
        TaskText:  parseResp.Task.Text,
        Status:    "parsed",
    }
    h.store.SaveAttempt(ctx, attempt)

    // 6. Return JSON
    json.NewEncoder(w).Encode(UploadTaskResponse{
        AttemptID: attempt.ID,
        Task:      parseResp,
    })
}
```

### 4. **Hint Flow** (из `telegram/hints.go`)

**Текущая state machine:**
```
Hints → ShowHint(L1) → ShowHint(L2) → ShowHint(L3) → EnterAnswer
```

**Новая реализация (REST API):**

```go
// POST /api/v1/attempts/:id/hints
func (h *Handler) GetHints(w http.ResponseWriter, r *http.Request) {
    attemptID := chi.URLParam(r, "id")
    userID := getUserIDFromContext(r.Context())

    // 1. Load attempt
    attempt, _ := h.store.GetAttempt(ctx, attemptID, userID)

    // 2. Get hints from LLM
    hintResp, _ := h.llmClient.Hint(ctx, llmName, hintReq)

    // 3. Update attempt
    attempt.Hints = hintResp.Items
    attempt.HintsUnlocked = 0 // Initially 0 visible
    h.store.UpdateAttempt(ctx, attempt)

    // 4. Return JSON
    json.NewEncoder(w).Encode(GetHintsResponse{
        Hints:         hintResp.Items,
        MaxHints:      len(hintResp.Items[0].Hints),
        UnlockedCount: 0,
    })
}

// POST /api/v1/attempts/:id/hints/unlock
func (h *Handler) UnlockHint(w http.ResponseWriter, r *http.Request) {
    attemptID := chi.URLParam(r, "id")
    level := r.FormValue("level") // "L1" | "L2" | "L3"

    // 1. Load attempt
    attempt, _ := h.store.GetAttempt(ctx, attemptID, userID)

    // 2. Check if user has coins/subscription
    if !h.canUnlockHint(userID, level) {
        http.Error(w, "insufficient coins", http.StatusPaymentRequired)
        return
    }

    // 3. Unlock hint
    attempt.HintsUnlocked++
    h.store.UpdateAttempt(ctx, attempt)

    // 4. Deduct coins or count towards subscription limit
    h.deductCoins(ctx, userID, hintCost[level])

    // 5. Return updated hints
    json.NewEncoder(w).Encode(UnlockHintResponse{
        UnlockedCount: attempt.HintsUnlocked,
        Hints:         attempt.Hints,
    })
}
```

### 5. **Check Flow** (из `telegram/check.go`)

**Новая реализация:**

```go
// POST /api/v1/attempts/:id/check
func (h *Handler) CheckAnswer(w http.ResponseWriter, r *http.Request) {
    attemptID := chi.URLParam(r, "id")
    userID := getUserIDFromContext(r.Context())

    var req CheckAnswerRequest
    json.NewDecoder(r.Body).Decode(&req)

    // 1. Load attempt
    attempt, _ := h.store.GetAttempt(ctx, attemptID, userID)

    // 2. Check solution via LLM
    checkResp, _ := h.llmClient.CheckSolution(ctx, llmName, checkReq)

    // 3. Update attempt
    attempt.Status = "checked"
    attempt.CheckResult = checkResp
    attempt.IsCorrect = checkResp.Decision == "correct"
    h.store.UpdateAttempt(ctx, attempt)

    // 4. Update user stats
    if attempt.IsCorrect {
        h.incrementCorrectCount(ctx, userID)
        h.awardCoins(ctx, userID, 10)
        h.updateVillainHealth(ctx, userID, -10) // Damage villain
    }

    // 5. Return result
    json.NewEncoder(w).Encode(CheckAnswerResponse{
        Status:     checkResp.Status,
        Decision:   checkResp.Decision,
        Feedback:   checkResp.Feedback,
        IsCorrect:  attempt.IsCorrect,
        CoinsEarned: 10,
    })
}
```

### 6. **Analogue Flow** (из `telegram/analogue.go`)

```go
// POST /api/v1/attempts/:id/analogue
func (h *Handler) GetAnalogue(w http.ResponseWriter, r *http.Request) {
    attemptID := chi.URLParam(r, "id")
    userID := getUserIDFromContext(r.Context())

    // 1. Load original attempt
    attempt, _ := h.store.GetAttempt(ctx, attemptID, userID)

    // 2. Generate analogue via LLM
    analogueResp, _ := h.llmClient.AnalogueSolution(ctx, llmName, analogueReq)

    // 3. Create new attempt
    newAttempt := &Attempt{
        ID:           uuid.NewString(),
        UserID:       userID,
        TaskText:     analogueResp.TaskText,
        Status:       "analogue_generated",
        OriginalID:   attemptID,
    }
    h.store.SaveAttempt(ctx, newAttempt)

    // 4. Return JSON
    json.NewEncoder(w).Encode(GetAnalogueResponse{
        AttemptID: newAttempt.ID,
        TaskText:  analogueResp.TaskText,
        Items:     analogueResp.Items,
    })
}
```

---

## 🏗️ Новая архитектура REST API

### Слои архитектуры

```
┌─────────────────────────────────────┐
│   Frontend MiniApp (VK/Max/TG)      │
│   React + TypeScript + VKUI         │
└─────────────────┬───────────────────┘
                  │ HTTPS REST API (JWT)
┌─────────────────▼───────────────────┐
│   API Gateway / REST API Server     │
│   - JWT Authentication              │
│   - Rate Limiting                   │
│   - Request Validation              │
│   - Response Formatting             │
└─────────────────┬───────────────────┘
                  │
┌─────────────────▼───────────────────┐
│   Business Logic Layer              │
│   - Session Manager                 │
│   - Photo Processor                 │
│   - Attempt Manager                 │
│   - Coin/Subscription Manager       │
│   - Achievement Manager             │
└───┬─────────────────────────────┬───┘
    │                             │
┌───▼─────────────────┐  ┌────────▼─────────────┐
│   PostgreSQL        │  │   LLM Server         │
│   - sessions        │  │   - /v2/detect       │
│   - attempts        │  │   - /v2/parse        │
│   - users           │  │   - /v2/hint         │
│   - achievements    │  │   - /v2/check        │
│   - subscriptions   │  │   - /v2/analogue     │
└─────────────────────┘  └──────────────────────┘
```

---

## 📝 Итоговый список изменений

### ✅ Переиспользуем без изменений

- `api/internal/v2/llmclient/` - HTTP клиент для LLM-сервера
- `api/internal/v2/types/` - типы данных (Hint, Check, Analogue, etc.)
- `api/internal/store/` - PostgreSQL store interface
- `api/internal/config/` - конфигурация

### 🔄 Извлекаем и переписываем

- **Session Management** → platform-agnostic версия
- **Photo Processing** → библиотека без Telegram API
- **Detect + Parse Flow** → REST API endpoint
- **Hint Flow** → REST API endpoints
- **Check Flow** → REST API endpoint
- **Analogue Flow** → REST API endpoint

### ❌ Удаляем полностью

- `cmd/bot/main.go` - точка входа Telegram бота
- `api/internal/v2/telegram/` - вся директория
- Dependencies: `github.com/go-telegram-bot-api/telegram-bot-api`

### ➕ Создаём новое

- `cmd/api/main.go` - точка входа REST API сервера
- `api/internal/handler/` - HTTP handlers
- `api/internal/middleware/` - JWT auth, rate limiting
- `api/internal/service/` - бизнес-логика (Session, Photo, Attempt)

---

## 📊 Оценка трудозатрат

| Задача | Сложность | Время |
|--------|-----------|-------|
| Создать REST API skeleton | Средняя | 1 день |
| JWT Authentication middleware | Низкая | 0.5 дня |
| Извлечь Session Management | Низкая | 0.5 дня |
| Извлечь Photo Processing | Средняя | 1 день |
| Реализовать Upload Task endpoint | Средняя | 1 день |
| Реализовать Hints endpoints | Средняя | 1 день |
| Реализовать Check endpoint | Средняя | 1 день |
| Реализовать Analogue endpoint | Низкая | 0.5 дня |
| Добавить Profile endpoints | Низкая | 0.5 дня |
| Добавить Achievements endpoints | Средняя | 1 день |
| Добавить Friends/Referral endpoints | Средняя | 1 день |
| Тестирование и отладка | Высокая | 2 дня |
| **ИТОГО** | | **11-12 дней** |

---

## 🎯 Следующие шаги

1. **Прочитай [16_BACKEND_API.md](./16_BACKEND_API.md)** - детальный roadmap по backend
2. **Прочитай обновлённый [API_DATA_REQUIREMENTS.md](./API_DATA_REQUIREMENTS.md)** - все новые endpoints
3. **Согласуй архитектуру** с frontend командой
4. **Создай ветку** `feature/rest-api-migration`
5. **Начни разработку** последовательно по endpoints

---

**Готово к миграции!** ✅

**Next:** [16_BACKEND_API.md](./16_BACKEND_API.md)
