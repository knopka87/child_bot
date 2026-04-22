# 05: Авторизация (Telegram Init Data + JWT)

> Фаза 1 | Приоритет: P0 | Сложность: Средняя | Срок: 2-3 дня

## Цель

Реализовать безопасную авторизацию пользователей Mini App через Telegram Init Data с выдачей JWT токенов.

## Как работает Telegram Mini App Auth

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   Mini App      │     │   Backend API   │     │   Telegram      │
│   (Frontend)    │     │   (Go)          │     │   Servers       │
└────────┬────────┘     └────────┬────────┘     └────────┬────────┘
         │                       │                       │
         │ 1. Open Mini App      │                       │
         │ ◄─────────────────────┼───────────────────────│
         │    (initData in URL)  │                       │
         │                       │                       │
         │ 2. POST /auth/telegram│                       │
         │    {initData: "..."}  │                       │
         ├──────────────────────►│                       │
         │                       │                       │
         │                       │ 3. Validate hash      │
         │                       │    using bot_token    │
         │                       ├───────────────────────►
         │                       │                       │
         │                       │ 4. Hash valid         │
         │                       │◄───────────────────────
         │                       │                       │
         │ 5. JWT token          │                       │
         │◄───────────────────────                       │
         │                       │                       │
         │ 6. API requests       │                       │
         │    with JWT           │                       │
         ├──────────────────────►│                       │
         │                       │                       │
```

## Telegram Init Data Format

```
query_id=AAHdF6IQAAAAAN0XohDhrOrc
&user=%7B%22id%22%3A279058397%2C%22first_name%22%3A%22Vladislav%22%2C%22last_name%22%3A%22Kibenko%22%2C%22username%22%3A%22vdkfrost%22%2C%22language_code%22%3A%22ru%22%7D
&auth_date=1662771648
&hash=c501b71e775f74ce10e377dea85a7ea24ecd640b223ea86dfe453e0eaed2e2b2
```

## Реализация

### Auth Handler

```go
// internal/api/handlers/auth.go
package handlers

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "sort"
    "strconv"
    "strings"
    "time"

    "github.com/golang-jwt/jwt/v5"

    "child_bot/api/internal/store"
)

type AuthHandler struct {
    store     *store.Store
    botToken  string
    jwtSecret string
    jwtTTL    time.Duration
}

func NewAuthHandler(s *store.Store, botToken, jwtSecret string) *AuthHandler {
    return &AuthHandler{
        store:     s,
        botToken:  botToken,
        jwtSecret: jwtSecret,
        jwtTTL:    24 * time.Hour,
    }
}

type TelegramAuthRequest struct {
    InitData string `json:"init_data"`
}

type TelegramAuthResponse struct {
    Token     string    `json:"token"`
    ExpiresAt time.Time `json:"expires_at"`
    User      UserInfo  `json:"user"`
}

type UserInfo struct {
    ID        int64  `json:"id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name,omitempty"`
    Username  string `json:"username,omitempty"`
    PhotoURL  string `json:"photo_url,omitempty"`
    Role      string `json:"role"`
}

type TelegramUser struct {
    ID           int64  `json:"id"`
    FirstName    string `json:"first_name"`
    LastName     string `json:"last_name"`
    Username     string `json:"username"`
    LanguageCode string `json:"language_code"`
    IsPremium    bool   `json:"is_premium"`
    PhotoURL     string `json:"photo_url"`
}

func (h *AuthHandler) AuthTelegram(w http.ResponseWriter, r *http.Request) {
    var req TelegramAuthRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.jsonError(w, "invalid request body", http.StatusBadRequest)
        return
    }

    if req.InitData == "" {
        h.jsonError(w, "init_data is required", http.StatusBadRequest)
        return
    }

    // 1. Parse and validate init data
    user, err := h.validateInitData(req.InitData)
    if err != nil {
        h.jsonError(w, fmt.Sprintf("invalid init_data: %v", err), http.StatusUnauthorized)
        return
    }

    // 2. Upsert user in database
    dbUser, err := h.store.UpsertUser(r.Context(), store.UserUpsert{
        ChatID:    user.ID,
        Username:  user.Username,
        FirstName: user.FirstName,
        LastName:  user.LastName,
    })
    if err != nil {
        h.jsonError(w, "failed to create user", http.StatusInternalServerError)
        return
    }

    // 3. Generate JWT
    expiresAt := time.Now().Add(h.jwtTTL)

    claims := &UserClaims{
        UserID:   user.ID,
        Username: user.Username,
        Role:     dbUser.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expiresAt),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "child-bot",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(h.jwtSecret))
    if err != nil {
        h.jsonError(w, "failed to generate token", http.StatusInternalServerError)
        return
    }

    // 4. Return response
    resp := TelegramAuthResponse{
        Token:     tokenString,
        ExpiresAt: expiresAt,
        User: UserInfo{
            ID:        user.ID,
            FirstName: user.FirstName,
            LastName:  user.LastName,
            Username:  user.Username,
            PhotoURL:  user.PhotoURL,
            Role:      dbUser.Role,
        },
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) validateInitData(initData string) (*TelegramUser, error) {
    // Parse query string
    values, err := url.ParseQuery(initData)
    if err != nil {
        return nil, fmt.Errorf("invalid query string: %w", err)
    }

    // Extract hash
    hash := values.Get("hash")
    if hash == "" {
        return nil, fmt.Errorf("hash is missing")
    }
    values.Del("hash")

    // Check auth_date (prevent replay attacks)
    authDateStr := values.Get("auth_date")
    if authDateStr == "" {
        return nil, fmt.Errorf("auth_date is missing")
    }

    authDate, err := strconv.ParseInt(authDateStr, 10, 64)
    if err != nil {
        return nil, fmt.Errorf("invalid auth_date: %w", err)
    }

    // Reject if older than 1 hour
    if time.Now().Unix()-authDate > 3600 {
        return nil, fmt.Errorf("auth_date is too old")
    }

    // Build data-check-string
    var keys []string
    for k := range values {
        keys = append(keys, k)
    }
    sort.Strings(keys)

    var pairs []string
    for _, k := range keys {
        pairs = append(pairs, k+"="+values.Get(k))
    }
    dataCheckString := strings.Join(pairs, "\n")

    // Calculate HMAC-SHA256
    secretKey := hmac.New(sha256.New, []byte("WebAppData"))
    secretKey.Write([]byte(h.botToken))
    secret := secretKey.Sum(nil)

    h2 := hmac.New(sha256.New, secret)
    h2.Write([]byte(dataCheckString))
    calculatedHash := hex.EncodeToString(h2.Sum(nil))

    // Compare hashes
    if calculatedHash != hash {
        return nil, fmt.Errorf("hash mismatch")
    }

    // Parse user data
    userJSON := values.Get("user")
    if userJSON == "" {
        return nil, fmt.Errorf("user data is missing")
    }

    var user TelegramUser
    if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
        return nil, fmt.Errorf("invalid user JSON: %w", err)
    }

    return &user, nil
}

func (h *AuthHandler) jsonError(w http.ResponseWriter, message string, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// UserClaims for JWT
type UserClaims struct {
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}
```

### Store: Upsert User

```go
// internal/store/user.go

type UserUpsert struct {
    ChatID    int64
    Username  string
    FirstName string
    LastName  string
}

type User struct {
    ChatID    int64
    Username  string
    FirstName string
    LastName  string
    Grade     int
    Role      string // "student" | "parent"
    XP        int
    Level     int
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (s *Store) UpsertUser(ctx context.Context, u UserUpsert) (*User, error) {
    var user User

    err := s.DB.QueryRowContext(ctx, `
        INSERT INTO "user" (chat_id, grade, role, xp, level)
        VALUES ($1, 0, 'student', 0, 1)
        ON CONFLICT (chat_id) DO UPDATE
        SET updated_at = NOW()
        RETURNING chat_id, grade, role, xp, level, created_at, updated_at
    `, u.ChatID).Scan(
        &user.ChatID, &user.Grade, &user.Role, &user.XP, &user.Level,
        &user.CreatedAt, &user.UpdatedAt,
    )
    if err != nil {
        return nil, fmt.Errorf("upsert user: %w", err)
    }

    // Update chat info
    _, err = s.DB.ExecContext(ctx, `
        INSERT INTO chat (id, username, first_name, last_name, type)
        VALUES ($1, $2, $3, $4, 'private')
        ON CONFLICT (id) DO UPDATE
        SET username = $2, first_name = $3, last_name = $4
    `, u.ChatID, u.Username, u.FirstName, u.LastName)
    if err != nil {
        return nil, fmt.Errorf("upsert chat: %w", err)
    }

    user.Username = u.Username
    user.FirstName = u.FirstName
    user.LastName = u.LastName

    return &user, nil
}
```

### Миграция для новых полей

```sql
-- migrations/030_user_roles.up.sql

ALTER TABLE "user" ADD COLUMN IF NOT EXISTS role TEXT NOT NULL DEFAULT 'student';
ALTER TABLE "user" ADD COLUMN IF NOT EXISTS xp INT NOT NULL DEFAULT 0;
ALTER TABLE "user" ADD COLUMN IF NOT EXISTS level INT NOT NULL DEFAULT 1;
ALTER TABLE "user" ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE "user" ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_user_role ON "user"(role);
```

### JWT Refresh (опционально)

```go
// internal/api/handlers/auth.go

type RefreshRequest struct {
    Token string `json:"token"`
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
    var req RefreshRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.jsonError(w, "invalid request body", http.StatusBadRequest)
        return
    }

    // Parse token without validation (we allow expired for refresh)
    token, err := jwt.ParseWithClaims(req.Token, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
        return []byte(h.jwtSecret), nil
    })

    claims, ok := token.Claims.(*UserClaims)
    if !ok {
        h.jsonError(w, "invalid token", http.StatusUnauthorized)
        return
    }

    // Check if token is not too old (max 7 days for refresh)
    if claims.ExpiresAt != nil {
        expiry := claims.ExpiresAt.Time
        if time.Since(expiry) > 7*24*time.Hour {
            h.jsonError(w, "token too old for refresh", http.StatusUnauthorized)
            return
        }
    }

    // Issue new token
    expiresAt := time.Now().Add(h.jwtTTL)
    newClaims := &UserClaims{
        UserID:   claims.UserID,
        Username: claims.Username,
        Role:     claims.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expiresAt),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "child-bot",
        },
    }

    newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
    tokenString, err := newToken.SignedString([]byte(h.jwtSecret))
    if err != nil {
        h.jsonError(w, "failed to generate token", http.StatusInternalServerError)
        return
    }

    resp := map[string]interface{}{
        "token":      tokenString,
        "expires_at": expiresAt,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}
```

## Frontend Integration

```typescript
// Mini App frontend example
import { initData } from '@telegram-apps/sdk';

async function authenticate() {
  const response = await fetch('https://api.example.com/api/v1/auth/telegram', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      init_data: initData.raw(),
    }),
  });

  if (!response.ok) {
    throw new Error('Authentication failed');
  }

  const { token, user } = await response.json();

  // Store token
  localStorage.setItem('jwt_token', token);

  return user;
}

// Use token in subsequent requests
async function apiRequest(path: string, options: RequestInit = {}) {
  const token = localStorage.getItem('jwt_token');

  return fetch(`https://api.example.com${path}`, {
    ...options,
    headers: {
      ...options.headers,
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
  });
}
```

## Security Considerations

1. **HTTPS Only**: API должен работать только через HTTPS
2. **Short TTL**: JWT токены действуют 24 часа
3. **Replay Attack Prevention**: auth_date проверяется (max 1 hour old)
4. **Secure Secret**: JWT secret минимум 32 символа
5. **Rate Limiting**: Ограничение запросов авторизации

## Тестирование

```go
// internal/api/handlers/auth_test.go
func TestValidateInitData(t *testing.T) {
    handler := NewAuthHandler(nil, "test-bot-token", "test-secret")

    tests := []struct {
        name      string
        initData  string
        wantError bool
    }{
        {
            name:      "valid init data",
            initData:  generateValidInitData(t, "test-bot-token"),
            wantError: false,
        },
        {
            name:      "invalid hash",
            initData:  "user=%7B%22id%22%3A123%7D&auth_date=1234567890&hash=invalid",
            wantError: true,
        },
        {
            name:      "expired auth_date",
            initData:  generateExpiredInitData(t, "test-bot-token"),
            wantError: true,
        },
        {
            name:      "missing user",
            initData:  "auth_date=1234567890&hash=abc",
            wantError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := handler.validateInitData(tt.initData)
            if (err != nil) != tt.wantError {
                t.Errorf("validateInitData() error = %v, wantError %v", err, tt.wantError)
            }
        })
    }
}

func generateValidInitData(t *testing.T, botToken string) string {
    t.Helper()

    user := `{"id":123,"first_name":"Test","username":"test_user"}`
    authDate := strconv.FormatInt(time.Now().Unix(), 10)

    dataCheckString := fmt.Sprintf("auth_date=%s\nuser=%s", authDate, user)

    secretKey := hmac.New(sha256.New, []byte("WebAppData"))
    secretKey.Write([]byte(botToken))
    secret := secretKey.Sum(nil)

    h := hmac.New(sha256.New, secret)
    h.Write([]byte(dataCheckString))
    hash := hex.EncodeToString(h.Sum(nil))

    return fmt.Sprintf("user=%s&auth_date=%s&hash=%s",
        url.QueryEscape(user), authDate, hash)
}
```

## Чек-лист

- [ ] Реализовать `handlers/auth.go`
- [ ] Добавить `validateInitData` с HMAC-SHA256
- [ ] Реализовать JWT генерацию
- [ ] Добавить refresh token endpoint (опционально)
- [ ] Создать миграцию для user roles
- [ ] Написать unit-тесты
- [ ] Документировать для frontend
- [ ] Добавить rate limiting на auth endpoint
- [ ] Тестирование с реальным Telegram Mini App

## Связанные шаги

- [04-api-layer.md](./04-api-layer.md) — использует auth middleware
- [11-parent-child.md](./11-parent-child.md) — роли parent/student

---

[← API Layer](./04-api-layer.md) | [Назад к Roadmap](./roadmap.md) | [Далее: Service Refactoring →](./06-service-refactoring.md)
