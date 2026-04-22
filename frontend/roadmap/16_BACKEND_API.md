# Phase 15: Backend REST API

**Длительность:** 11-12 дней
**Приоритет:** Критический (блокирует frontend разработку)
**Зависимости:** Нет

---

## 🎯 Цель Phase

Мигрировать backend с Telegram Bot на REST API для поддержки VK/Max/Telegram миниаппов.

**Что делаем:**
- ❌ Удаляем Telegram Bot зависимости
- ✅ Создаём REST API Server
- ✅ Извлекаем бизнес-логику из `internal/v2/telegram/`
- ✅ Добавляем JWT аутентификацию
- ✅ Реализуем все endpoints для миниаппа

---

## 📋 Содержание

- [Архитектура](#архитектура)
- [Step 1: Setup REST API Server](#step-1-setup-rest-api-server)
- [Step 2: JWT Authentication](#step-2-jwt-authentication)
- [Step 3: Session Management](#step-3-session-management)
- [Step 4: Photo Processing](#step-4-photo-processing)
- [Step 5: Task Upload Endpoint](#step-5-task-upload-endpoint)
- [Step 6: Hints Endpoints](#step-6-hints-endpoints)
- [Step 7: Check Endpoint](#step-7-check-endpoint)
- [Step 8: Analogue Endpoint](#step-8-analogue-endpoint)
- [Step 9: Profile Endpoints](#step-9-profile-endpoints)
- [Step 10: Achievements Endpoints](#step-10-achievements-endpoints)
- [Step 11: Friends/Referral Endpoints](#step-11-friendsreferral-endpoints)
- [Step 12: Testing](#step-12-testing)
- [Чеклист](#чеклист)

---

## 🏗️ Архитектура

### Новая структура проекта

```
api/
├── cmd/
│   └── api/
│       └── main.go                    # ➕ REST API server entry point
├── internal/
│   ├── handler/                        # ➕ HTTP handlers
│   │   ├── auth.go                    # VK sign validation, JWT
│   │   ├── task.go                    # Upload, Get tasks
│   │   ├── hint.go                    # Get/Unlock hints
│   │   ├── check.go                   # Check answer
│   │   ├── analogue.go                # Generate analogue
│   │   ├── profile.go                 # User profile
│   │   ├── achievement.go             # Achievements
│   │   └── friend.go                  # Referrals
│   ├── middleware/                     # ➕ HTTP middleware
│   │   ├── auth.go                    # JWT authentication
│   │   ├── cors.go                    # CORS headers
│   │   ├── ratelimit.go               # Rate limiting
│   │   └── logger.go                  # Request logging
│   ├── service/                        # ➕ Business logic (platform-agnostic)
│   │   ├── session.go                 # Session management
│   │   ├── photo.go                   # Photo processing
│   │   ├── attempt.go                 # Attempt management
│   │   ├── coin.go                    # Coins & subscription
│   │   ├── achievement.go             # Achievement unlocking
│   │   └── villain.go                 # Villain health
│   ├── v2/
│   │   ├── llmclient/                 # ✅ REUSE (no changes)
│   │   ├── types/                     # ✅ REUSE (no changes)
│   │   └── telegram/                  # ❌ DELETE (deprecated)
│   ├── store/                          # ✅ REUSE + EXTEND
│   │   ├── session.go                 # ✅ Reuse
│   │   ├── attempt.go                 # ➕ New: attempts table
│   │   ├── user.go                    # ✅ Reuse + extend
│   │   ├── achievement.go             # ➕ New
│   │   ├── referral.go                # ➕ New
│   │   └── subscription.go            # ➕ New
│   └── util/                           # ✅ REUSE
└── migrations/                         # ➕ Database migrations
    ├── 001_create_attempts.sql
    ├── 002_create_achievements.sql
    ├── 003_create_referrals.sql
    └── 004_create_subscriptions.sql
```

### Слои взаимодействия

```
┌─────────────────────────────────────┐
│         Frontend MiniApp            │
└──────────────┬──────────────────────┘
               │ JWT Token (Authorization header)
               │
┌──────────────▼──────────────────────┐
│         HTTP Router (chi)           │
│    GET  /api/v1/profile/me          │
│    POST /api/v1/tasks/upload        │
│    POST /api/v1/attempts/:id/hints  │
│    POST /api/v1/attempts/:id/check  │
└──────────────┬──────────────────────┘
               │ Middleware chain
┌──────────────▼──────────────────────┐
│         Middleware Stack            │
│    1. CORS                          │
│    2. Logger                        │
│    3. JWT Auth ← extract user_id   │
│    4. Rate Limiter                  │
└──────────────┬──────────────────────┘
               │ Context with user_id
┌──────────────▼──────────────────────┐
│           Handler Layer             │
│    handler.GetProfile(w, r)         │
│    handler.UploadTask(w, r)         │
│    handler.GetHints(w, r)           │
└──────────────┬──────────────────────┘
               │ Service calls
┌──────────────▼──────────────────────┐
│          Service Layer              │
│    attemptService.Create()          │
│    photoService.Process()           │
│    coinService.Deduct()             │
└──┬───────────────────────────────┬──┘
   │                               │
┌──▼─────────────────┐  ┌──────────▼──────────┐
│   Store (DB)       │  │   LLM Client        │
│   PostgreSQL       │  │   HTTP → LLM Server │
└────────────────────┘  └─────────────────────┘
```

---

## Step 1: Setup REST API Server

**Длительность:** 1 день

### 1.1. Создать точку входа

**Файл:** `cmd/api/main.go`

```go
package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"

	"child-bot/api/internal/config"
	"child-bot/api/internal/handler"
	"child-bot/api/internal/llmclient"
	mid "child-bot/api/internal/middleware"
	"child-bot/api/internal/service"
	"child-bot/api/internal/store"
	"child-bot/api/internal/v2/llmclient"
)

func main() {
	cfg := config.Load()

	// Connect to PostgreSQL
	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db connection failed: %v", err)
	}
	defer db.Close()

	// Initialize stores
	st := store.NewStore(db)

	// Initialize LLM client
	llmClient := llmclient.New(&llmclient.Client{
		Base: cfg.LLMServerURL,
		HC:   &http.Client{Timeout: 3 * time.Minute},
	})

	// Initialize services
	sessionSvc := service.NewSessionService(st)
	photoSvc := service.NewPhotoService(10*1024*1024, 4096*4096)
	attemptSvc := service.NewAttemptService(st, llmClient)
	coinSvc := service.NewCoinService(st)
	achievementSvc := service.NewAchievementService(st)
	villainSvc := service.NewVillainService(st)

	// Initialize handlers
	h := handler.NewHandler(handler.HandlerDeps{
		Store:          st,
		LLMClient:      llmClient,
		SessionSvc:     sessionSvc,
		PhotoSvc:       photoSvc,
		AttemptSvc:     attemptSvc,
		CoinSvc:        coinSvc,
		AchievementSvc: achievementSvc,
		VillainSvc:     villainSvc,
		JWTSecret:      cfg.JWTSecret,
		VKSecretKey:    cfg.VKSecretKey,
	})

	// Setup router
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(mid.Logger)
	r.Use(middleware.Recoverer)
	r.Use(mid.CORS())

	// Health check (no auth)
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Public routes
	r.Post("/api/v1/auth/vk", h.AuthVK)
	r.Post("/api/v1/auth/telegram", h.AuthTelegram)

	// Protected routes (require JWT)
	r.Group(func(r chi.Router) {
		r.Use(mid.JWTAuth(cfg.JWTSecret))

		// Profile
		r.Get("/api/v1/profile/me", h.GetProfile)
		r.Patch("/api/v1/profile/me", h.UpdateProfile)

		// Tasks & Attempts
		r.Post("/api/v1/tasks/upload", h.UploadTask)
		r.Get("/api/v1/attempts", h.ListAttempts)
		r.Get("/api/v1/attempts/{id}", h.GetAttempt)

		// Hints
		r.Post("/api/v1/attempts/{id}/hints", h.GetHints)
		r.Post("/api/v1/attempts/{id}/hints/unlock", h.UnlockHint)

		// Check
		r.Post("/api/v1/attempts/{id}/check", h.CheckAnswer)

		// Analogue
		r.Post("/api/v1/attempts/{id}/analogue", h.GenerateAnalogue)

		// Achievements
		r.Get("/api/v1/achievements", h.ListAchievements)

		// Friends
		r.Get("/api/v1/referrals", h.GetReferralInfo)
		r.Post("/api/v1/referrals/apply", h.ApplyReferralCode)

		// Villain
		r.Get("/api/v1/villain", h.GetVillainStatus)
	})

	// HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		log.Printf("API server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
	log.Println("server stopped")
}
```

### 1.2. Обновить `config/config.go`

```go
package config

type Config struct {
	// Server
	Port        string
	Environment string // "development" | "production"

	// Database
	DatabaseURL string

	// LLM Server
	LLMServerURL string
	DefaultLLM   string // "gpt-4o-mini"

	// Auth
	JWTSecret     string // For signing JWT tokens
	VKSecretKey   string // For validating VK sign
	TGBotToken    string // For validating Telegram WebApp initData

	// VK
	VKAppID int64
}

func Load() *Config {
	return &Config{
		Port:          getenv("PORT", "8080"),
		Environment:   getenv("ENVIRONMENT", "development"),
		DatabaseURL:   mustGetenv("DATABASE_URL"),
		LLMServerURL:  mustGetenv("LLM_SERVER_URL"),
		DefaultLLM:    getenv("DEFAULT_LLM", "gpt-4o-mini"),
		JWTSecret:     mustGetenv("JWT_SECRET"),
		VKSecretKey:   mustGetenv("VK_SECRET_KEY"),
		TGBotToken:    getenv("TELEGRAM_BOT_TOKEN", ""),
		VKAppID:       mustGetenvInt64("VK_APP_ID"),
	}
}
```

---

## Step 2: JWT Authentication

**Длительность:** 0.5 дня

### 2.1. Создать JWT middleware

**Файл:** `internal/middleware/auth.go`

```go
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDContextKey contextKey = "user_id"

// JWTAuth middleware validates JWT token and extracts user_id
func JWTAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			// Parse and validate token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			// Extract user_id from claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
				return
			}

			userID, ok := claims["user_id"].(string)
			if !ok || userID == "" {
				http.Error(w, "missing user_id in token", http.StatusUnauthorized)
				return
			}

			// Add user_id to context
			ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID extracts user_id from context
func GetUserID(r *http.Request) string {
	userID, _ := r.Context().Value(UserIDContextKey).(string)
	return userID
}
```

### 2.2. Создать VK Sign Validation

**Файл:** `internal/handler/auth.go`

```go
package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthVKRequest struct {
	// VK Launch Params (from VKWebAppGetLaunchParams)
	VKUserID       int64  `json:"vk_user_id"`
	VKAppID        int64  `json:"vk_app_id"`
	VKAccessToken  string `json:"vk_access_token_settings"`
	VKLanguage     string `json:"vk_language"`
	VKPlatform     string `json:"vk_platform"`
	VKAreNotifications int `json:"vk_are_notifications_enabled"`
	VKRef          string `json:"vk_ref"`
	VKTs           int64  `json:"vk_ts"`
	Sign           string `json:"sign"` // HMAC-SHA256 signature
}

type AuthVKResponse struct {
	Token  string `json:"token"`  // JWT token
	UserID string `json:"user_id"` // child_profile_id
}

func (h *Handler) AuthVK(w http.ResponseWriter, r *http.Request) {
	var req AuthVKRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// 1. Validate VK sign
	if !h.validateVKSign(req) {
		http.Error(w, "invalid vk sign", http.StatusUnauthorized)
		return
	}

	// 2. Get or create user
	userID := fmt.Sprintf("vk_%d", req.VKUserID)
	user, err := h.Store.GetUserByPlatformID(r.Context(), "vk", req.VKUserID)
	if err != nil {
		// Create new user
		user = &store.User{
			ID:         userID,
			Platform:   "vk",
			PlatformID: req.VKUserID,
			Language:   req.VKLanguage,
			CreatedAt:  time.Now(),
		}
		if err := h.Store.CreateUser(r.Context(), user); err != nil {
			http.Error(w, "failed to create user", http.StatusInternalServerError)
			return
		}
	}

	// 3. Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"platform":  "vk",
		"vk_user_id": req.VKUserID,
		"exp":       time.Now().Add(30 * 24 * time.Hour).Unix(), // 30 days
		"iat":       time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.JWTSecret))
	if err != nil {
		http.Error(w, "failed to sign token", http.StatusInternalServerError)
		return
	}

	// 4. Return token
	json.NewEncoder(w).Encode(AuthVKResponse{
		Token:  tokenString,
		UserID: user.ID,
	})
}

// validateVKSign validates HMAC-SHA256 signature from VK
func (h *Handler) validateVKSign(req AuthVKRequest) bool {
	// Build query string from launch params (excluding sign)
	params := url.Values{}
	params.Add("vk_user_id", fmt.Sprint(req.VKUserID))
	params.Add("vk_app_id", fmt.Sprint(req.VKAppID))
	if req.VKAccessToken != "" {
		params.Add("vk_access_token_settings", req.VKAccessToken)
	}
	params.Add("vk_language", req.VKLanguage)
	params.Add("vk_platform", req.VKPlatform)
	params.Add("vk_are_notifications_enabled", fmt.Sprint(req.VKAreNotifications))
	if req.VKRef != "" {
		params.Add("vk_ref", req.VKRef)
	}
	params.Add("vk_ts", fmt.Sprint(req.VKTs))

	// Sort keys alphabetically
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build query string
	var queryParts []string
	for _, k := range keys {
		queryParts = append(queryParts, k+"="+params.Get(k))
	}
	queryString := strings.Join(queryParts, "&")

	// Compute HMAC-SHA256
	mac := hmac.New(sha256.New, []byte(h.VKSecretKey))
	mac.Write([]byte(queryString))
	expectedSign := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	// Remove padding
	expectedSign = strings.TrimRight(expectedSign, "=")
	receivedSign := strings.TrimRight(req.Sign, "=")

	return hmac.Equal([]byte(expectedSign), []byte(receivedSign))
}
```

---

## Step 3: Session Management

**Длительность:** 0.5 дня

### 3.1. Создать Session Service

**Файл:** `internal/service/session.go`

```go
package service

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	"child-bot/api/internal/store"
)

type SessionService struct {
	store store.Store
	cache *sync.Map // userID -> *Session
}

func NewSessionService(store store.Store) *SessionService {
	return &SessionService{
		store: store,
		cache: &sync.Map{},
	}
}

type Session struct {
	ID        string
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GetOrCreate returns existing session or creates a new one
func (s *SessionService) GetOrCreate(ctx context.Context, userID string) (*Session, error) {
	// Try cache first
	if cached, ok := s.cache.Load(userID); ok {
		return cached.(*Session), nil
	}

	// Try database
	dbSession, err := s.store.FindSessionByUserID(ctx, userID)
	if err == nil && dbSession != nil {
		session := &Session{
			ID:        dbSession.ID,
			UserID:    dbSession.UserID,
			CreatedAt: dbSession.CreatedAt,
			UpdatedAt: dbSession.UpdatedAt,
		}
		s.cache.Store(userID, session)
		return session, nil
	}

	// Create new session
	session := &Session{
		ID:        uuid.NewString(),
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to database
	if err := s.store.CreateSession(ctx, &store.Session{
		ID:        session.ID,
		UserID:    session.UserID,
		CreatedAt: session.CreatedAt,
	}); err != nil {
		return nil, err
	}

	// Cache
	s.cache.Store(userID, session)
	return session, nil
}

// Clear removes session from cache (useful for logout)
func (s *SessionService) Clear(userID string) {
	s.cache.Delete(userID)
}
```

---

## Step 4: Photo Processing

**Длительность:** 1 день

### 4.1. Создать Photo Service

**Файл:** `internal/service/photo.go`

```go
package service

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math"
)

type PhotoService struct {
	maxSizeBytes int
	maxPixels    int
}

func NewPhotoService(maxSizeBytes, maxPixels int) *PhotoService {
	return &PhotoService{
		maxSizeBytes: maxSizeBytes,
		maxPixels:    maxPixels,
	}
}

// Process handles single or multiple images
func (ps *PhotoService) Process(images [][]byte) ([]byte, error) {
	if len(images) == 0 {
		return nil, errors.New("no images provided")
	}

	if len(images) == 1 {
		return images[0], nil
	}

	return ps.combineVertically(images)
}

// combineVertically combines multiple images into one (vertical stack)
func (ps *PhotoService) combineVertically(images [][]byte) ([]byte, error) {
	// Decode all images
	decoded := make([]image.Image, 0, len(images))
	widths := make([]int, 0, len(images))
	heights := make([]int, 0, len(images))

	for _, imgBytes := range images {
		img, err := ps.decodeImage(imgBytes)
		if err != nil {
			return nil, fmt.Errorf("decode error: %w", err)
		}
		decoded = append(decoded, img)
		bounds := img.Bounds()
		widths = append(widths, bounds.Dx())
		heights = append(heights, bounds.Dy())
	}

	// Calculate dimensions
	maxWidth := 0
	totalHeight := 0
	for i := range decoded {
		if widths[i] > maxWidth {
			maxWidth = widths[i]
		}
		totalHeight += heights[i]
	}

	if maxWidth == 0 || totalHeight == 0 {
		return nil, errors.New("empty images")
	}

	// Create canvas
	canvas := image.NewRGBA(image.Rect(0, 0, maxWidth, totalHeight))
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	// Draw images vertically centered
	y := 0
	for i, img := range decoded {
		w := widths[i]
		h := heights[i]
		x := (maxWidth - w) / 2
		rect := image.Rect(x, y, x+w, y+h)
		draw.Draw(canvas, rect, img, img.Bounds().Min, draw.Over)
		y += h
	}

	// Resize if too large
	totalPixels := maxWidth * totalHeight
	final := image.Image(canvas)
	if totalPixels > ps.maxPixels {
		scale := math.Sqrt(float64(ps.maxPixels) / float64(totalPixels))
		newW := int(float64(maxWidth)*scale + 0.5)
		newH := int(float64(totalHeight)*scale + 0.5)
		if newW < 1 {
			newW = 1
		}
		if newH < 1 {
			newH = 1
		}
		final = ps.scaleDown(canvas, newW, newH)
	}

	// Encode as JPEG
	var out bytes.Buffer
	if err := jpeg.Encode(&out, final, &jpeg.Options{Quality: 90}); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func (ps *PhotoService) decodeImage(data []byte) (image.Image, error) {
	// Try JPEG
	if len(data) >= 2 && data[0] == 0xFF && data[1] == 0xD8 {
		return jpeg.Decode(bytes.NewReader(data))
	}

	// Try PNG
	if len(data) >= 8 &&
		data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		return png.Decode(bytes.NewReader(data))
	}

	// Fallback
	img, _, err := image.Decode(bytes.NewReader(data))
	return img, err
}

func (ps *PhotoService) scaleDown(src image.Image, newW, newH int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	sb := src.Bounds()
	srcW := sb.Dx()
	srcH := sb.Dy()

	for y := 0; y < newH; y++ {
		sy := sb.Min.Y + (y * srcH) / newH
		for x := 0; x < newW; x++ {
			sx := sb.Min.X + (x * srcW) / newW
			dst.Set(x, y, src.At(sx, sy))
		}
	}

	return dst
}
```

---

## Step 5: Task Upload Endpoint

**Длительность:** 1 день

### 5.1. Создать Attempt Service

**Файл:** `internal/service/attempt.go`

```go
package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"

	"child-bot/api/internal/store"
	"child-bot/api/internal/v2/llmclient"
	"child-bot/api/internal/v2/types"
)

type AttemptService struct {
	store     store.Store
	llmClient *llmclient.Client
}

func NewAttemptService(store store.Store, llmClient *llmclient.Client) *AttemptService {
	return &AttemptService{
		store:     store,
		llmClient: llmClient,
	}
}

type CreateAttemptRequest struct {
	UserID string
	Photo  []byte
}

type Attempt struct {
	ID          string
	UserID      string
	SessionID   string
	Photo       []byte
	Subject     string
	Grade       int
	TaskText    string
	ParsedData  *types.ParseResponse
	Status      string // "uploaded", "parsed", "hints_requested", "checked"
	CreatedAt   time.Time
}

// Create processes photo and creates attempt
func (as *AttemptService) Create(ctx context.Context, req CreateAttemptRequest) (*Attempt, error) {
	// 1. Encode photo as base64
	imageB64 := base64.StdEncoding.EncodeToString(req.Photo)

	// 2. Detect subject and grade
	detectResp, err := as.llmClient.Detect(ctx, "gpt-4o-mini", types.DetectRequest{
		Image: imageB64,
	})
	if err != nil {
		return nil, fmt.Errorf("detect failed: %w", err)
	}

	// 3. Parse task structure
	parseResp, err := as.llmClient.Parse(ctx, "gpt-4o-mini", types.ParseRequest{
		Image:   imageB64,
		Subject: detectResp.Subject,
		Grade:   detectResp.Grade,
		Locale:  "ru",
	})
	if err != nil {
		return nil, fmt.Errorf("parse failed: %w", err)
	}

	// 4. Create attempt
	attempt := &Attempt{
		ID:         uuid.NewString(),
		UserID:     req.UserID,
		Photo:      req.Photo,
		Subject:    string(detectResp.Subject),
		Grade:      detectResp.Grade,
		TaskText:   parseResp.Task.Text,
		ParsedData: &parseResp,
		Status:     "parsed",
		CreatedAt:  time.Now(),
	}

	// 5. Save to database
	if err := as.store.CreateAttempt(ctx, &store.Attempt{
		ID:         attempt.ID,
		UserID:     attempt.UserID,
		Photo:      attempt.Photo,
		Subject:    attempt.Subject,
		Grade:      attempt.Grade,
		TaskText:   attempt.TaskText,
		ParsedData: parseResp,
		Status:     attempt.Status,
		CreatedAt:  attempt.CreatedAt,
	}); err != nil {
		return nil, fmt.Errorf("save attempt failed: %w", err)
	}

	return attempt, nil
}

// Get retrieves attempt by ID (ensures user owns it)
func (as *AttemptService) Get(ctx context.Context, attemptID, userID string) (*Attempt, error) {
	dbAttempt, err := as.store.GetAttempt(ctx, attemptID)
	if err != nil {
		return nil, err
	}

	// Security: verify ownership
	if dbAttempt.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	return &Attempt{
		ID:         dbAttempt.ID,
		UserID:     dbAttempt.UserID,
		Photo:      dbAttempt.Photo,
		Subject:    dbAttempt.Subject,
		Grade:      dbAttempt.Grade,
		TaskText:   dbAttempt.TaskText,
		ParsedData: dbAttempt.ParsedData,
		Status:     dbAttempt.Status,
		CreatedAt:  dbAttempt.CreatedAt,
	}, nil
}
```

### 5.2. Создать Upload Handler

**Файл:** `internal/handler/task.go`

```go
package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"child-bot/api/internal/middleware"
	"child-bot/api/internal/service"
)

type UploadTaskResponse struct {
	AttemptID  string                `json:"attempt_id"`
	Subject    string                `json:"subject"`
	Grade      int                   `json:"grade"`
	TaskText   string                `json:"task_text"`
	ParsedData *types.ParseResponse  `json:"parsed_data"`
}

func (h *Handler) UploadTask(w http.ResponseWriter, r *http.Request) {
	// 1. Get user ID from JWT
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// 2. Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	// 3. Get uploaded files
	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		http.Error(w, "no images uploaded", http.StatusBadRequest)
		return
	}

	// 4. Read images
	images := make([][]byte, 0, len(files))
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "failed to open image", http.StatusBadRequest)
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "failed to read image", http.StatusBadRequest)
			return
		}
		images = append(images, data)
	}

	// 5. Process photos (combine if multiple)
	photo, err := h.PhotoSvc.Process(images)
	if err != nil {
		http.Error(w, "failed to process images", http.StatusInternalServerError)
		return
	}

	// 6. Create attempt
	attempt, err := h.AttemptSvc.Create(r.Context(), service.CreateAttemptRequest{
		UserID: userID,
		Photo:  photo,
	})
	if err != nil {
		http.Error(w, "failed to create attempt", http.StatusInternalServerError)
		return
	}

	// 7. Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UploadTaskResponse{
		AttemptID:  attempt.ID,
		Subject:    attempt.Subject,
		Grade:      attempt.Grade,
		TaskText:   attempt.TaskText,
		ParsedData: attempt.ParsedData,
	})
}
```

---

## Step 6: Hints Endpoints

**Длительность:** 1 день

### 6.1. Extend Attempt Service

**Добавить в** `internal/service/attempt.go`:

```go
// GetHints generates hints for attempt
func (as *AttemptService) GetHints(ctx context.Context, attemptID, userID string) (*types.HintResponse, error) {
	// 1. Get attempt
	attempt, err := as.Get(ctx, attemptID, userID)
	if err != nil {
		return nil, err
	}

	// 2. Check if hints already generated
	if attempt.ParsedData.Hints != nil {
		return attempt.ParsedData.Hints, nil
	}

	// 3. Generate hints via LLM
	hintResp, err := as.llmClient.Hint(ctx, "gpt-4o-mini", types.HintRequest{
		Task:  attempt.ParsedData.Task,
		Items: attempt.ParsedData.Items,
		Mode:  "learn",
	})
	if err != nil {
		return nil, fmt.Errorf("hint generation failed: %w", err)
	}

	// 4. Update attempt in database
	attempt.ParsedData.Hints = &hintResp
	if err := as.store.UpdateAttemptHints(ctx, attemptID, &hintResp); err != nil {
		return nil, err
	}

	return &hintResp, nil
}

// UnlockHint unlocks next hint level
func (as *AttemptService) UnlockHint(ctx context.Context, attemptID, userID, level string) error {
	// 1. Get attempt
	attempt, err := as.Get(ctx, attemptID, userID)
	if err != nil {
		return err
	}

	// 2. Increment unlocked count
	unlockedCount := attempt.HintsUnlocked + 1

	// 3. Update database
	return as.store.UpdateAttemptHintsUnlocked(ctx, attemptID, unlockedCount)
}
```

### 6.2. Создать Hint Handler

**Файл:** `internal/handler/hint.go`

```go
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"child-bot/api/internal/middleware"
)

type GetHintsResponse struct {
	Hints         []types.HintItem `json:"hints"`
	UnlockedCount int              `json:"unlocked_count"`
	MaxHints      int              `json:"max_hints"`
}

func (h *Handler) GetHints(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	attemptID := chi.URLParam(r, "id")

	// Get hints
	hintResp, err := h.AttemptSvc.GetHints(r.Context(), attemptID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get attempt to check unlocked count
	attempt, _ := h.AttemptSvc.Get(r.Context(), attemptID, userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetHintsResponse{
		Hints:         hintResp.Items,
		UnlockedCount: attempt.HintsUnlocked,
		MaxHints:      len(hintResp.Items[0].Hints),
	})
}

type UnlockHintRequest struct {
	Level string `json:"level"` // "L1" | "L2" | "L3"
}

type UnlockHintResponse struct {
	UnlockedCount int `json:"unlocked_count"`
	CoinsSpent    int `json:"coins_spent"`
}

func (h *Handler) UnlockHint(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	attemptID := chi.URLParam(r, "id")

	var req UnlockHintRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Check if user can unlock (has coins or subscription)
	canUnlock, cost := h.CoinSvc.CanUnlockHint(r.Context(), userID, req.Level)
	if !canUnlock {
		http.Error(w, "insufficient coins", http.StatusPaymentRequired)
		return
	}

	// Unlock hint
	if err := h.AttemptSvc.UnlockHint(r.Context(), attemptID, userID, req.Level); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Deduct coins (if not subscription)
	if cost > 0 {
		h.CoinSvc.Deduct(r.Context(), userID, cost, "hint_unlock")
	}

	// Get updated attempt
	attempt, _ := h.AttemptSvc.Get(r.Context(), attemptID, userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UnlockHintResponse{
		UnlockedCount: attempt.HintsUnlocked,
		CoinsSpent:    cost,
	})
}
```

---

## Step 7: Check Endpoint

**Длительность:** 1 день

### 7.1. Extend Attempt Service

**Добавить в** `internal/service/attempt.go`:

```go
type CheckAnswerRequest struct {
	AttemptID  string
	UserID     string
	Photo      []byte // Photo of student's solution
}

type CheckAnswerResult struct {
	Status     string
	Decision   string
	IsCorrect  bool
	Feedback   string
	ErrorSpans []types.ErrorSpan
	Confidence float64
}

func (as *AttemptService) CheckAnswer(ctx context.Context, req CheckAnswerRequest) (*CheckAnswerResult, error) {
	// 1. Get attempt
	attempt, err := as.Get(ctx, req.AttemptID, req.UserID)
	if err != nil {
		return nil, err
	}

	// 2. Encode photo as base64
	answerPhotoB64 := base64.StdEncoding.EncodeToString(req.Photo)

	// 3. Check solution via LLM
	checkResp, err := as.llmClient.CheckSolution(ctx, "gpt-4o-mini", types.CheckRequest{
		Image:      answerPhotoB64,
		TaskStruct: types.TaskStructCheck{
			TaskTextClean: attempt.TaskText,
			Items:         attempt.ParsedData.Items,
		},
		RawTaskText: attempt.TaskText,
		Student: types.StudentCheck{
			Grade:   int64(attempt.Grade),
			Subject: attempt.Subject,
			Locale:  "ru",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("check failed: %w", err)
	}

	// 4. Update attempt
	isCorrect := checkResp.Decision == types.CheckDecisionCorrect
	if err := as.store.UpdateAttemptCheckResult(ctx, req.AttemptID, checkResp, isCorrect); err != nil {
		return nil, err
	}

	confidence := 0.0
	if checkResp.Confidence != nil {
		confidence = *checkResp.Confidence
	}

	return &CheckAnswerResult{
		Status:     string(checkResp.Status),
		Decision:   string(checkResp.Decision),
		IsCorrect:  isCorrect,
		Feedback:   checkResp.Feedback,
		ErrorSpans: checkResp.ErrorSpans,
		Confidence: confidence,
	}, nil
}
```

### 7.2. Создать Check Handler

**Файл:** `internal/handler/check.go`

```go
package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"child-bot/api/internal/middleware"
	"child-bot/api/internal/service"
)

type CheckAnswerResponse struct {
	Status     string              `json:"status"`
	Decision   string              `json:"decision"`
	IsCorrect  bool                `json:"is_correct"`
	Feedback   string              `json:"feedback"`
	ErrorSpans []types.ErrorSpan   `json:"error_spans"`
	Confidence float64             `json:"confidence"`
	CoinsEarned int                `json:"coins_earned"`
	VillainDamage int              `json:"villain_damage"`
}

func (h *Handler) CheckAnswer(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	attemptID := chi.URLParam(r, "id")

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	// Get answer photo
	file, _, err := r.FormFile("answer_photo")
	if err != nil {
		http.Error(w, "missing answer_photo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	photo, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "failed to read photo", http.StatusBadRequest)
		return
	}

	// Check answer
	result, err := h.AttemptSvc.CheckAnswer(r.Context(), service.CheckAnswerRequest{
		AttemptID: attemptID,
		UserID:    userID,
		Photo:     photo,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Reward if correct
	coinsEarned := 0
	villainDamage := 0
	if result.IsCorrect {
		coinsEarned = 10
		villainDamage = 10

		// Award coins
		h.CoinSvc.Award(r.Context(), userID, coinsEarned, "correct_answer")

		// Damage villain
		h.VillainSvc.DealDamage(r.Context(), userID, villainDamage)

		// Check achievements
		h.AchievementSvc.CheckAndUnlock(r.Context(), userID)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CheckAnswerResponse{
		Status:        result.Status,
		Decision:      result.Decision,
		IsCorrect:     result.IsCorrect,
		Feedback:      result.Feedback,
		ErrorSpans:    result.ErrorSpans,
		Confidence:    result.Confidence,
		CoinsEarned:   coinsEarned,
		VillainDamage: villainDamage,
	})
}
```

---

## Step 8: Analogue Endpoint

**Длительность:** 0.5 дня

### 8.1. Extend Attempt Service

```go
func (as *AttemptService) GenerateAnalogue(ctx context.Context, attemptID, userID string) (*Attempt, error) {
	// 1. Get original attempt
	original, err := as.Get(ctx, attemptID, userID)
	if err != nil {
		return nil, err
	}

	// 2. Generate analogue via LLM
	analogueResp, err := as.llmClient.AnalogueSolution(ctx, "gpt-4o-mini", types.AnalogueRequest{
		Subject:    types.Subject(original.Subject),
		Grade:      original.Grade,
		TaskStruct: original.ParsedData.Task,
	})
	if err != nil {
		return nil, fmt.Errorf("analogue generation failed: %w", err)
	}

	// 3. Create new attempt
	newAttempt := &Attempt{
		ID:        uuid.NewString(),
		UserID:    userID,
		Subject:   original.Subject,
		Grade:     original.Grade,
		TaskText:  analogueResp.Task.Text,
		ParsedData: &types.ParseResponse{
			Task:  analogueResp.Task,
			Items: analogueResp.Items,
		},
		Status:     "analogue_generated",
		CreatedAt:  time.Now(),
	}

	// 4. Save to database
	if err := as.store.CreateAttempt(ctx, &store.Attempt{
		ID:         newAttempt.ID,
		UserID:     newAttempt.UserID,
		Subject:    newAttempt.Subject,
		Grade:      newAttempt.Grade,
		TaskText:   newAttempt.TaskText,
		ParsedData: newAttempt.ParsedData,
		Status:     newAttempt.Status,
		CreatedAt:  newAttempt.CreatedAt,
	}); err != nil {
		return nil, err
	}

	return newAttempt, nil
}
```

### 8.2. Создать Analogue Handler

**Файл:** `internal/handler/analogue.go`

```go
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"child-bot/api/internal/middleware"
)

type GenerateAnalogueResponse struct {
	AttemptID  string               `json:"attempt_id"`
	TaskText   string               `json:"task_text"`
	ParsedData *types.ParseResponse `json:"parsed_data"`
}

func (h *Handler) GenerateAnalogue(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	attemptID := chi.URLParam(r, "id")

	// Generate analogue
	newAttempt, err := h.AttemptSvc.GenerateAnalogue(r.Context(), attemptID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GenerateAnalogueResponse{
		AttemptID:  newAttempt.ID,
		TaskText:   newAttempt.TaskText,
		ParsedData: newAttempt.ParsedData,
	})
}
```

---

## Step 9: Profile Endpoints

**Длительность:** 0.5 дня

### 9.1. Создать Profile Handler

**Файл:** `internal/handler/profile.go`

```go
package handler

import (
	"encoding/json"
	"net/http"

	"child-bot/api/internal/middleware"
)

type GetProfileResponse struct {
	UserID                  string `json:"user_id"`
	Name                    string `json:"name"`
	Avatar                  string `json:"avatar"`
	Grade                   int    `json:"grade"`
	Level                   int    `json:"level"`
	LevelProgress           int    `json:"level_progress_percent"`
	CoinsBalance            int    `json:"coins_balance"`
	TasksSolvedCorrectCount int    `json:"tasks_solved_correct_count"`
	VillainHealthPercent    int    `json:"villain_health_percent"`
	HasSubscription         bool   `json:"has_subscription"`
	SubscriptionExpiresAt   *time.Time `json:"subscription_expires_at"`
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	// Get user from database
	user, err := h.Store.GetUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// Get stats
	stats, _ := h.Store.GetUserStats(r.Context(), userID)

	// Get villain status
	villain, _ := h.VillainSvc.GetStatus(r.Context(), userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetProfileResponse{
		UserID:                  user.ID,
		Name:                    user.Name,
		Avatar:                  user.Avatar,
		Grade:                   user.Grade,
		Level:                   stats.Level,
		LevelProgress:           stats.LevelProgress,
		CoinsBalance:            stats.CoinsBalance,
		TasksSolvedCorrectCount: stats.CorrectCount,
		VillainHealthPercent:    villain.HealthPercent,
		HasSubscription:         user.HasSubscription,
		SubscriptionExpiresAt:   user.SubscriptionExpiresAt,
	})
}

type UpdateProfileRequest struct {
	Name   *string `json:"name"`
	Avatar *string `json:"avatar"`
	Grade  *int    `json:"grade"`
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Update user
	if err := h.Store.UpdateUser(r.Context(), userID, store.UpdateUserParams{
		Name:   req.Name,
		Avatar: req.Avatar,
		Grade:  req.Grade,
	}); err != nil {
		http.Error(w, "failed to update profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
```

---

## Step 10: Achievements Endpoints

**Длительность:** 1 день

### 10.1. Создать Achievement Service

**Файл:** `internal/service/achievement.go`

```go
package service

import (
	"context"

	"child-bot/api/internal/store"
)

type AchievementService struct {
	store store.Store
}

func NewAchievementService(store store.Store) *AchievementService {
	return &AchievementService{store: store}
}

// CheckAndUnlock проверяет и разблокирует достижения
func (as *AchievementService) CheckAndUnlock(ctx context.Context, userID string) error {
	// Получить статистику пользователя
	stats, err := as.store.GetUserStats(ctx, userID)
	if err != nil {
		return err
	}

	// Получить все достижения
	allAchievements, err := as.store.ListAchievements(ctx)
	if err != nil {
		return err
	}

	// Получить уже разблокированные
	unlocked, err := as.store.GetUnlockedAchievements(ctx, userID)
	if err != nil {
		return err
	}

	unlockedMap := make(map[string]bool)
	for _, a := range unlocked {
		unlockedMap[a.ID] = true
	}

	// Проверить каждое достижение
	for _, achievement := range allAchievements {
		if unlockedMap[achievement.ID] {
			continue // Already unlocked
		}

		// Check condition
		if as.checkCondition(achievement, stats) {
			// Unlock achievement
			if err := as.store.UnlockAchievement(ctx, userID, achievement.ID); err != nil {
				continue // Skip on error
			}

			// Award coins
			if achievement.RewardCoins > 0 {
				as.store.AddCoins(ctx, userID, achievement.RewardCoins, "achievement_"+achievement.ID)
			}
		}
	}

	return nil
}

func (as *AchievementService) checkCondition(achievement *store.Achievement, stats *store.UserStats) bool {
	switch achievement.Condition {
	case "correct_answers_5":
		return stats.CorrectCount >= 5
	case "correct_answers_10":
		return stats.CorrectCount >= 10
	case "correct_answers_50":
		return stats.CorrectCount >= 50
	case "streak_3":
		return stats.CurrentStreak >= 3
	case "streak_7":
		return stats.CurrentStreak >= 7
	case "level_5":
		return stats.Level >= 5
	case "level_10":
		return stats.Level >= 10
	case "invited_friends_3":
		return stats.InvitedFriends >= 3
	case "defeat_villain_1":
		return stats.VillainsDefeated >= 1
	default:
		return false
	}
}
```

### 10.2. Создать Achievement Handler

**Файл:** `internal/handler/achievement.go`

```go
package handler

import (
	"encoding/json"
	"net/http"

	"child-bot/api/internal/middleware"
)

type AchievementDTO struct {
	ID          string `json:"id"`
	Icon        string `json:"icon"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	IsUnlocked  bool   `json:"is_unlocked"`
	UnlockedAt  *time.Time `json:"unlocked_at"`
	RewardCoins int    `json:"reward_coins"`
}

type ListAchievementsResponse struct {
	Achievements      []AchievementDTO `json:"achievements"`
	TotalCount        int              `json:"total_count"`
	UnlockedCount     int              `json:"unlocked_count"`
}

func (h *Handler) ListAchievements(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	// Get all achievements
	all, err := h.Store.ListAchievements(r.Context())
	if err != nil {
		http.Error(w, "failed to list achievements", http.StatusInternalServerError)
		return
	}

	// Get unlocked
	unlocked, err := h.Store.GetUnlockedAchievements(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to get unlocked achievements", http.StatusInternalServerError)
		return
	}

	unlockedMap := make(map[string]*store.UnlockedAchievement)
	for i := range unlocked {
		unlockedMap[unlocked[i].AchievementID] = &unlocked[i]
	}

	// Build response
	achievements := make([]AchievementDTO, 0, len(all))
	unlockedCount := 0

	for _, a := range all {
		dto := AchievementDTO{
			ID:          a.ID,
			Icon:        a.Icon,
			Name:        a.Name,
			Description: a.Description,
			Category:    a.Category,
			IsUnlocked:  false,
			RewardCoins: a.RewardCoins,
		}

		if ua, ok := unlockedMap[a.ID]; ok {
			dto.IsUnlocked = true
			dto.UnlockedAt = &ua.UnlockedAt
			unlockedCount++
		}

		achievements = append(achievements, dto)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ListAchievementsResponse{
		Achievements:  achievements,
		TotalCount:    len(all),
		UnlockedCount: unlockedCount,
	})
}
```

---

## Step 11: Friends/Referral Endpoints

**Длительность:** 1 день

### 11.1. Создать Referral Handler

**Файл:** `internal/handler/friend.go`

```go
package handler

import (
	"encoding/json"
	"net/http"

	"child-bot/api/internal/middleware"
)

type GetReferralInfoResponse struct {
	ReferralCode  string `json:"referral_code"`
	ReferralLink  string `json:"referral_link"`
	InvitedCount  int    `json:"invited_count"`
	InvitedTarget int    `json:"invited_target"`
	RewardCoins   int    `json:"reward_coins_per_friend"`
}

func (h *Handler) GetReferralInfo(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	// Get user
	user, err := h.Store.GetUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// Get referral stats
	invitedCount, err := h.Store.GetReferralCount(r.Context(), userID)
	if err != nil {
		invitedCount = 0
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetReferralInfoResponse{
		ReferralCode:  user.ReferralCode,
		ReferralLink:  fmt.Sprintf("https://vk.com/app%d#ref=%s", h.VKAppID, user.ReferralCode),
		InvitedCount:  invitedCount,
		InvitedTarget: 5,
		RewardCoins:   50,
	})
}

type ApplyReferralCodeRequest struct {
	ReferralCode string `json:"referral_code"`
}

func (h *Handler) ApplyReferralCode(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req ApplyReferralCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Find inviter by referral code
	inviter, err := h.Store.GetUserByReferralCode(r.Context(), req.ReferralCode)
	if err != nil {
		http.Error(w, "invalid referral code", http.StatusNotFound)
		return
	}

	// Apply referral
	if err := h.Store.ApplyReferral(r.Context(), userID, inviter.ID); err != nil {
		http.Error(w, "failed to apply referral", http.StatusInternalServerError)
		return
	}

	// Reward inviter
	h.Store.AddCoins(r.Context(), inviter.ID, 50, "referral_invite")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
```

---

## Step 12: Testing

**Длительность:** 2 дня

### 12.1. Unit Tests

**Файл:** `internal/service/attempt_test.go`

```go
package service_test

import (
	"context"
	"testing"

	"child-bot/api/internal/service"
	"child-bot/api/internal/store/mock"
)

func TestAttemptService_Create(t *testing.T) {
	mockStore := mock.NewStore()
	mockLLMClient := mock.NewLLMClient()

	svc := service.NewAttemptService(mockStore, mockLLMClient)

	// Mock LLM responses
	mockLLMClient.DetectResponse = types.DetectResponse{
		Subject: "math",
		Grade:   5,
	}
	mockLLMClient.ParseResponse = types.ParseResponse{
		Task: types.ParseTask{
			Text: "Решите уравнение: 2x + 5 = 15",
		},
	}

	// Test
	attempt, err := svc.Create(context.Background(), service.CreateAttemptRequest{
		UserID: "test_user",
		Photo:  []byte("fake_photo"),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if attempt.Subject != "math" {
		t.Errorf("expected subject=math, got %s", attempt.Subject)
	}
}
```

### 12.2. Integration Tests

**Файл:** `test/integration/upload_task_test.go`

```go
package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"child-bot/api/cmd/api"
)

func TestUploadTask_Integration(t *testing.T) {
	// Setup test server
	srv := httptest.NewServer(api.SetupRouter())
	defer srv.Close()

	// Login to get JWT token
	token := loginTestUser(t, srv.URL)

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("images", "test.jpg")
	part.Write(loadTestImage(t, "testdata/math_problem.jpg"))
	writer.Close()

	// Send request
	req, _ := http.NewRequest("POST", srv.URL+"/api/v1/tasks/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	// Parse response
	var result handler.UploadTaskResponse
	json.NewDecoder(resp.Body).Decode(&result)

	if result.AttemptID == "" {
		t.Error("expected attempt_id, got empty")
	}

	if result.Subject == "" {
		t.Error("expected subject, got empty")
	}
}
```

---

## 📝 Чеклист

### Phase 15: Backend REST API

#### Step 1: Setup REST API Server
- [ ] Создан `cmd/api/main.go`
- [ ] Настроен chi router
- [ ] Добавлены middleware (CORS, Logger, Recoverer)
- [ ] Создан healthcheck endpoint
- [ ] Настроен graceful shutdown

#### Step 2: JWT Authentication
- [ ] Создан JWT middleware
- [ ] Реализована VK sign validation
- [ ] Создан endpoint `/api/v1/auth/vk`
- [ ] Создан endpoint `/api/v1/auth/telegram` (опционально)
- [ ] Добавлена защита всех endpoints с помощью JWT

#### Step 3: Session Management
- [ ] Создан Session Service
- [ ] Реализован TTL cache для сессий
- [ ] Добавлена таблица `sessions` в БД
- [ ] Написаны unit tests

#### Step 4: Photo Processing
- [ ] Создан Photo Service
- [ ] Реализована склейка множественных фото
- [ ] Реализовано ресайзинг изображений
- [ ] Написаны unit tests

#### Step 5: Task Upload Endpoint
- [ ] Создан Attempt Service
- [ ] Реализован endpoint `POST /api/v1/tasks/upload`
- [ ] Интегрирован LLM Client (Detect + Parse)
- [ ] Добавлена таблица `attempts` в БД
- [ ] Написаны integration tests

#### Step 6: Hints Endpoints
- [ ] Реализован метод `GetHints` в Attempt Service
- [ ] Реализован метод `UnlockHint` в Attempt Service
- [ ] Создан endpoint `POST /api/v1/attempts/:id/hints`
- [ ] Создан endpoint `POST /api/v1/attempts/:id/hints/unlock`
- [ ] Добавлена проверка coins/subscription
- [ ] Написаны tests

#### Step 7: Check Endpoint
- [ ] Реализован метод `CheckAnswer` в Attempt Service
- [ ] Создан endpoint `POST /api/v1/attempts/:id/check`
- [ ] Интегрирован LLM Client (CheckSolution)
- [ ] Добавлено начисление coins за правильные ответы
- [ ] Добавлено нанесение урона злодею
- [ ] Написаны tests

#### Step 8: Analogue Endpoint
- [ ] Реализован метод `GenerateAnalogue` в Attempt Service
- [ ] Создан endpoint `POST /api/v1/attempts/:id/analogue`
- [ ] Интегрирован LLM Client (AnalogueSolution)
- [ ] Написаны tests

#### Step 9: Profile Endpoints
- [ ] Создан endpoint `GET /api/v1/profile/me`
- [ ] Создан endpoint `PATCH /api/v1/profile/me`
- [ ] Добавлена таблица `user_stats` в БД
- [ ] Написаны tests

#### Step 10: Achievements Endpoints
- [ ] Создан Achievement Service
- [ ] Создан endpoint `GET /api/v1/achievements`
- [ ] Добавлена таблица `achievements` в БД
- [ ] Добавлена таблица `unlocked_achievements` в БД
- [ ] Реализована автоматическая разблокировка при условиях
- [ ] Написаны tests

#### Step 11: Friends/Referral Endpoints
- [ ] Создан endpoint `GET /api/v1/referrals`
- [ ] Создан endpoint `POST /api/v1/referrals/apply`
- [ ] Добавлена таблица `referrals` в БД
- [ ] Реализовано начисление coins за приглашения
- [ ] Написаны tests

#### Step 12: Testing
- [ ] Написаны unit tests для всех services
- [ ] Написаны integration tests для всех endpoints
- [ ] Проверена security (JWT validation, user ownership)
- [ ] Проверен rate limiting
- [ ] Проведено load testing (100+ rps)

#### Удаление Telegram Bot
- [ ] Удалена директория `internal/v2/telegram/`
- [ ] Удалён `cmd/bot/main.go`
- [ ] Удалены dependencies на `telegram-bot-api`

#### Документация
- [ ] Обновлён API_DATA_REQUIREMENTS.md
- [ ] Создана Swagger/OpenAPI спецификация
- [ ] Обновлён README.md

---

## ✅ Definition of Done

**Phase 15 считается завершённым, когда:**

1. ✅ REST API сервер запускается и отвечает на healthcheck
2. ✅ Все 15+ endpoints реализованы и работают
3. ✅ JWT аутентификация работает
4. ✅ VK sign validation работает
5. ✅ LLM Client интегрирован (Detect, Parse, Hint, Check, Analogue)
6. ✅ PostgreSQL store расширен новыми таблицами
7. ✅ Unit tests покрывают >80% кода services
8. ✅ Integration tests покрывают все endpoints
9. ✅ Telegram Bot полностью удалён
10. ✅ API документация обновлена

---

**Готово к разработке!** 🚀

**Next Phase:** [00_OVERVIEW.md](./00_OVERVIEW.md) (обновлённый с Phase 15)
