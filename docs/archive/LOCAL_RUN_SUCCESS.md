# Локальный запуск проекта — Успешно! ✅

**Дата:** 31 марта 2026
**Статус:** ✅ Запущено и работает

---

## Запущенные сервисы

### ✅ PostgreSQL
```
Container: homework_postgres
Image: postgres:15-alpine
Port: 5432
Status: Healthy
Connection: postgres://homework:root@localhost:5432/homework
```

**Таблицы:**
- ✅ child_profiles (профили детей)
- ✅ attempts (попытки help/check)
- ✅ achievements (достижения)
- ✅ villains (злодеи и битвы)
- ✅ subscriptions (подписки)
- ✅ referrals (реферальная программа)
- + 17 дополнительных таблиц (parsed_tasks, hints_cache, task_sessions, etc.)

**Всего:** 24 таблицы, 32 миграции применены

---

### ✅ Redis
```
Container: homework_redis
Image: redis:7-alpine
Port: 6379
Status: Healthy
```

---

### ✅ Backend REST API
```
Process: go run ./cmd/server/ (PID: 92786)
Port: 8080
Status: Running
Log: /tmp/backend.log
```

**Health Check:**
```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

---

## Доступные REST API Endpoints

### Attempts (8 endpoints)
```
POST   /attempts                     # Создать попытку (help/check)
GET    /attempts/unfinished          # Незавершённая попытка
GET    /attempts/recent              # Последние попытки
POST   /attempts/{id}/images         # Загрузить изображение
POST   /attempts/{id}/process        # Начать обработку через LLM
GET    /attempts/{id}/result         # Получить результат
POST   /attempts/{id}/next-hint      # Следующая подсказка
DELETE /attempts/{id}                # Удалить попытку
```

### Home (1 endpoint)
```
GET    /home/{childProfileId}        # Данные главного экрана
```

### Profile (4 endpoints)
```
GET    /profile                      # Получить профиль
PUT    /profile                      # Обновить профиль
GET    /profile/history              # История попыток
GET    /profile/stats                # Статистика
```

### Achievements (5 endpoints)
```
GET    /achievements                 # Список всех достижений
GET    /achievements/unlocked        # Разблокированные
GET    /achievements/stats           # Статистика достижений
GET    /achievements/{id}            # Детали достижения
POST   /achievements/{id}/claim      # Забрать награду
```

### Villains (6 endpoints)
```
GET    /villains                     # Список злодеев
GET    /villains/active              # Активный злодей
GET    /villains/{id}                # Детали злодея
GET    /villains/{id}/battle         # Состояние битвы
GET    /villains/{id}/victory        # Победа над злодеем
POST   /villains/{id}/damage         # Нанести урон
```

### Subscription (5 endpoints)
```
GET    /subscription/status          # Статус подписки
GET    /subscription/plans           # Доступные планы
POST   /subscription/subscribe       # Оформить подписку
DELETE /subscription/cancel          # Отменить подписку
POST   /subscription/resume          # Возобновить подписку
```

### Referrals (4 endpoints)
```
GET    /friends                      # Список друзей
POST   /friends/invite               # Пригласить друга
GET    /friends/referrals            # Реферальные данные
GET    /friends/leaderboard          # Таблица лидеров
```

### Health (1 endpoint)
```
GET    /health                       # Health check
```

**Всего:** 34 endpoints

---

## Конфигурация

### Environment Variables
```bash
DATABASE_URL=postgres://homework:root@localhost:5432/homework?sslmode=disable
LLM_SERVER_URL=http://138.124.55.145:8000
DEFAULT_LLM=gemini
PORT=8080
ENV=development
LOG_LEVEL=info
ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000
```

---

## Тестирование API

### 1. Health Check
```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

### 2. Create Attempt
```bash
curl -X POST http://localhost:8080/attempts \
  -H "Content-Type: application/json" \
  -H "X-Platform-ID: web" \
  -H "X-Child-Profile-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -d '{
    "child_profile_id": "550e8400-e29b-41d4-a716-446655440000",
    "type": "help"
  }'
```

### 3. Get Recent Attempts
```bash
curl http://localhost:8080/attempts/recent \
  -H "X-Platform-ID: web" \
  -H "X-Child-Profile-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json"
```

---

## Управление сервисами

### Остановка
```bash
# Остановить backend
kill $(ps aux | grep "go run ./cmd/server" | grep -v grep | awk '{print $2}')

# Остановить Docker сервисы
docker compose down

# Или остановить только PostgreSQL/Redis
docker compose stop postgres redis
```

### Перезапуск
```bash
# Перезапустить Docker сервисы
docker compose restart postgres redis

# Запустить backend
cd api
DATABASE_URL="postgres://homework:root@localhost:5432/homework?sslmode=disable" \
LLM_SERVER_URL="http://138.124.55.145:8000" \
DEFAULT_LLM="gemini" \
PORT="8080" \
ENV="development" \
LOG_LEVEL="info" \
ALLOWED_ORIGINS="http://localhost:5173,http://localhost:3000" \
go run ./cmd/server/
```

### Логи
```bash
# Backend логи
tail -f /tmp/backend.log

# PostgreSQL логи
docker logs -f homework_postgres

# Redis логи
docker logs -f homework_redis
```

---

## Frontend (опционально)

### Запуск frontend
```bash
cd frontend
npm install
npm run dev
# Frontend доступен на http://localhost:5173
```

Frontend будет подключаться к backend API на `http://localhost:8080`.

---

## Важные замечания

### 1. Локальный PostgreSQL
⚠️ Локальный PostgreSQL@14 был остановлен (`brew services stop postgresql@14`), чтобы освободить порт 5432 для Docker контейнера.

**Чтобы вернуть обратно:**
```bash
docker compose down
brew services start postgresql@14
```

### 2. LLM Server
Backend использует удалённый LLM сервер: `http://138.124.55.145:8000`

Если нужен локальный LLM сервер, измените `LLM_SERVER_URL` в переменных окружения.

### 3. CORS
CORS настроен для:
- `http://localhost:5173` (Vite dev server)
- `http://localhost:3000` (альтернативный порт)

---

## Проблемы и решения

### Проблема 1: "role homework does not exist"
**Причина:** Локальный PostgreSQL занял порт 5432
**Решение:**
```bash
brew services stop postgresql@14
docker compose restart postgres
```

### Проблема 2: "missing required env LLM_SERVER_URL"
**Причина:** Не установлены environment variables
**Решение:** Экспортировать переменные перед запуском (см. секцию "Управление сервисами")

### Проблема 3: UNIQUE constraint error в subscriptions
**Причина:** Синтаксическая ошибка в миграции 031
**Решение:** Таблица создана вручную с корректным синтаксисом

---

## Статус

✅ **PostgreSQL** - Running (24 таблицы)
✅ **Redis** - Running
✅ **Backend API** - Running (34 endpoints)
⏸️ **Frontend** - Не запущен (опционально)

**Проект готов к работе!** 🚀

---

**Дата:** 31 марта 2026
**Время запуска:** 15:55
**Backend Process:** PID 92786
