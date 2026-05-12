# Руководство по локальному запуску миниаппа

## Текущее состояние проекта

### ✅ Что готово

**Frontend (React + TypeScript + Vite)**
- ✅ Полная реализация UI компонентов
- ✅ Все страницы (Home, Tasks, Villain, Victory, Profile)
- ✅ Сервисы (Analytics, Platform adapters)
- ✅ Конфигурация (config, routes, validation)
- ✅ TypeScript без ошибок
- ✅ Production build: 508 KB gzipped
- ✅ Тесты: 9/9 passing
- ✅ Структурированное логирование
- ✅ Runtime валидация с zod
- ✅ Request deduplication и retry logic

**Backend (Go Telegram Bot)**
- ✅ Telegram bot с версиями v1 и v2
- ✅ База данных PostgreSQL с миграциями
- ✅ Redis для кэширования
- ✅ Интеграция с LLM (OpenAI)
- ✅ Health check endpoint `/healthz`
- ✅ Docker конфигурация

### ❌ Что отсутствует

**REST API для миниаппа**
- ❌ Backend реализован только как Telegram bot (webhook/polling)
- ❌ Нет HTTP API endpoints для frontend:
  - `/api/v1/tasks` - список заданий
  - `/api/v1/tasks/:id` - детали задания
  - `/api/v1/tasks/:id/submit` - отправка ответа
  - `/api/v1/villains/:id` - информация о злодее
  - `/api/v1/profile` - профиль пользователя
  - `/api/v1/analytics/events` - отправка аналитики
  - И другие ~30 endpoints из `frontend/src/api/routes.ts`

**Проблема**: Frontend расчитан на REST API, но backend - это только Telegram bot без HTTP API layer.

---

## Варианты запуска

### Вариант 1: Frontend только для демонстрации UI (без данных)

Запуск frontend без backend для просмотра интерфейса.

#### Требования
- Node.js 20+ (LTS)
- npm 10+

#### Шаги

1. **Установка зависимостей**
```bash
cd /Users/a.yanover/Xsolla/child_bot/frontend
npm install
```

2. **Создание `.env` файла**
```bash
cp .env.example .env
```

3. **Редактирование `.env`**
```bash
# Укажите базовый URL для API (будет недоступен без backend)
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_API_TIMEOUT=30000

# Включите debug для разработки
VITE_ANALYTICS_DEBUG=true

# Platform IDs (для тестирования в VK/Telegram/Max)
VITE_VK_APP_ID=your_vk_app_id
VITE_MAX_APP_ID=your_max_app_id
VITE_TELEGRAM_BOT_USERNAME=your_bot_username
```

4. **Запуск dev сервера**
```bash
npm run dev
```

5. **Открыть в браузере**
```
http://localhost:5173
```

**Ограничения**:
- ❌ API запросы будут падать с ошибками (backend недоступен)
- ❌ Нельзя загрузить задания, профиль, данные о злодеях
- ✅ Можно посмотреть UI компоненты и навигацию
- ✅ Работают platform adapters (VK Bridge, Telegram WebApp, MAX Bridge)

---

### Вариант 2: Telegram bot + База данных (без миниаппа)

Запуск только Telegram бота для работы через Telegram.

#### Требования
- Docker 24+ и Docker Compose
- Telegram Bot Token (от @BotFather)
- OpenAI API Key (для LLM)

#### Шаги

1. **Создание root `.env` файла**
```bash
cd /Users/a.yanover/Xsolla/child_bot
cat > .env << 'EOF'
# Database
POSTGRES_DB=homework
POSTGRES_USER=homework
POSTGRES_PASSWORD=homework_secret
POSTGRES_PORT=5432

# Redis
REDIS_PASSWORD=redis_secret
REDIS_PORT=6379

# Backend
BACKEND_PORT=8080
ENV=development
LOG_LEVEL=debug

# Telegram Bot
TELEGRAM_BOT_TOKEN=your_bot_token_here
TELEGRAM_BOT_VERSION=v2
# TELEGRAM_WEBHOOK_URL=  # Оставьте пустым для polling mode

# OpenAI
OPENAI_API_KEY=your_openai_api_key_here
OPENAI_MODEL=gpt-4

# JWT
JWT_SECRET=super_secret_jwt_key_change_in_production
JWT_EXPIRES_IN=15m

# S3 (опционально)
# S3_ENDPOINT=
# S3_ACCESS_KEY=
# S3_SECRET_KEY=
# S3_BUCKET=homework-images
EOF
```

2. **Запуск через Docker Compose**
```bash
docker compose up -d postgres redis backend
```

3. **Проверка логов**
```bash
docker compose logs -f backend
```

4. **Проверка health**
```bash
curl http://localhost:8080/healthz
# Ответ: ok
```

5. **Тестирование через Telegram**
Откройте Telegram и напишите вашему боту `/start`

**Что работает**:
- ✅ Telegram bot принимает сообщения
- ✅ База данных с миграциями
- ✅ Кэширование в Redis
- ✅ LLM интеграция для проверки заданий
- ❌ Веб-миниапп не работает (нет REST API)

---

### Вариант 3: Полноценный запуск миниаппа (требует REST API)

⚠️ **Недоступен**: Требуется реализация REST API layer для frontend.

#### Что нужно сделать

1. **Создать HTTP API handlers в Go backend**
   - Реализовать REST endpoints из `frontend/src/api/routes.ts`
   - Добавить middleware для CORS, authentication
   - Интегрировать с существующей логикой Telegram бота

2. **Архитектура**
```
┌─────────────────┐
│  Frontend       │
│  (React SPA)    │
│  Port: 5173     │
└────────┬────────┘
         │ HTTP REST API
         ↓
┌─────────────────┐
│  Backend API    │  ← Нужно создать
│  (Go HTTP)      │
│  Port: 8080     │
└────────┬────────┘
         │
    ┌────┴─────┐
    ↓          ↓
┌────────┐ ┌────────┐
│ Postgres│ │ Redis  │
└────────┘ └────────┘

┌─────────────────┐
│  Telegram Bot   │  ← Уже существует
│  (Go Webhook)   │
└────────┬────────┘
         │
    ┌────┴─────┐
    ↓          ↓
┌────────┐ ┌────────┐
│ Postgres│ │ Redis  │
└────────┘ └────────┘
```

3. **Требуемые endpoints (примеры)**

```go
// api/cmd/webapp/main.go

// Tasks
GET    /api/v1/tasks
GET    /api/v1/tasks/:id
POST   /api/v1/tasks/:id/submit

// Villains
GET    /api/v1/villains
GET    /api/v1/villains/:id
GET    /api/v1/villains/:id/attempts

// Profile
GET    /api/v1/profile
PUT    /api/v1/profile
GET    /api/v1/profile/stats
GET    /api/v1/profile/achievements

// Analytics
POST   /api/v1/analytics/events
POST   /api/v1/analytics/properties

// Hints
GET    /api/v1/hints
POST   /api/v1/hints
GET    /api/v1/hints/:id

// Auth
POST   /api/v1/auth/login
POST   /api/v1/auth/refresh
POST   /api/v1/auth/logout
```

4. **Примерная оценка работы**: 40-60 часов
   - Создание HTTP handlers: ~15h
   - Middleware (auth, CORS, logging): ~8h
   - Интеграция с существующей БД: ~10h
   - Тестирование: ~10h
   - Документация API: ~5h
   - Frontend интеграция: ~5h

---

## Текущие скрипты

### Frontend (`package.json`)

```bash
npm run dev          # Dev сервер (Vite) на http://localhost:5173
npm run build        # Production build → dist/
npm run preview      # Preview production build
npm run lint         # ESLint проверка
npm run format       # Prettier форматирование
npm run typecheck    # TypeScript type checking
npm run test         # Vitest тесты
```

### Docker Compose

```bash
# Запуск всех сервисов
docker compose up -d

# Только база данных
docker compose up -d postgres redis

# Только backend
docker compose up -d postgres redis backend

# Логи
docker compose logs -f backend

# Остановка
docker compose down

# Полная очистка (включая volumes)
docker compose down -v
```

---

## Проверка готовности

### Frontend Build

```bash
cd frontend
npm install
npm run typecheck  # ✅ No errors
npm run build      # ✅ Build successful: 508 KB gzipped
npm run test       # ✅ 9/9 tests passing
```

### Backend Build

```bash
cd api
go build -o bin/bot ./cmd/bot
# ✅ Компилируется без ошибок
```

### База данных

```bash
docker compose up -d postgres
docker compose exec postgres psql -U homework -d homework -c '\dt'
# ✅ 10+ таблиц из миграций
```

---

## Рекомендации

### Для демонстрации UI
→ Используйте **Вариант 1** (Frontend only)

### Для тестирования Telegram бота
→ Используйте **Вариант 2** (Bot + DB)

### Для полноценного миниаппа
→ Требуется реализация REST API (**Вариант 3**)

---

## Дополнительная информация

### Документация проекта
- `DESIGN_ANALYSIS.md` - архитектура и дизайн решений
- `MINIAPPS_PLATFORMS.md` - интеграция с VK, Telegram, MAX
- `CODE_REVIEW_FIXES.md` - исправленные проблемы кода
- `REMAINING_ISSUES_DETAILED.md` - оставшиеся задачи (7 LOW priority)

### Структура frontend
```
frontend/
├── src/
│   ├── components/     # UI компоненты
│   ├── pages/          # Страницы (Home, Tasks, Villain, etc.)
│   ├── services/       # Сервисы (API, Analytics, Platform)
│   ├── stores/         # Zustand stores
│   ├── hooks/          # Custom React hooks
│   ├── config/         # Конфигурация (config, routes, assets)
│   ├── lib/            # Утилиты (logger, validation)
│   └── types/          # TypeScript типы
├── tests/              # Vitest тесты
├── index.html          # HTML entry point
├── vite.config.ts      # Vite конфигурация
└── package.json        # NPM зависимости
```

### Платформы миниаппов

Frontend поддерживает запуск на:
- **VK Mini Apps** (VKUI + VK Bridge)
- **Telegram Mini Apps** (Telegram WebApp)
- **MAX Mini Apps** (MAX Bridge)
- **Веб-браузер** (fallback)

Platform detection автоматический через `PlatformBridge`.

---

## Вопросы и поддержка

### Ошибка: "Network Error" при API запросах
→ Backend не запущен или недоступен. Проверьте `VITE_API_BASE_URL` в `.env`

### Ошибка: "Cannot find module"
→ Запустите `npm install` в директории frontend

### Ошибка: Database connection failed
→ Проверьте `POSTGRES_*` переменные в root `.env` и статус контейнера

### Telegram bot не отвечает
→ Проверьте `TELEGRAM_BOT_TOKEN` и статус webhook: `docker compose logs backend`

---

## Следующие шаги

**Краткосрочные (1-2 дня)**:
1. Исправить 7 LOW priority issues из code review
2. Добавить E2E тесты для frontend

**Среднесрочные (1-2 недели)**:
1. Реализовать REST API layer для миниаппа
2. Интегрировать frontend с backend
3. Добавить authentication/authorization
4. Настроить production deployment

**Долгосрочные (1+ месяц)**:
1. Добавить real-time обновления (WebSockets)
2. Реализовать offline support с Service Workers
3. Оптимизация производительности
4. Расширенная аналитика и мониторинг
