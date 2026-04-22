# Full Stack Running — Frontend + Backend ✅

**Дата:** 31 марта 2026
**Статус:** ✅ Оба сервиса запущены и работают

---

## Запущенные сервисы

### ✅ Backend REST API
```
Process: go run ./cmd/server/ (PID: 92786)
Port: 8080
URL: http://localhost:8080
Health: http://localhost:8080/health → {"status":"ok"}
Log: /tmp/backend.log
```

**Endpoints:** 34 REST API endpoints доступны

---

### ✅ Frontend (Vite Dev Server)
```
Process: node vite (PID: 99381)
Port: 3000
URL: http://localhost:3000/
Protocol: HTTP (для локальной разработки)
Log: /tmp/frontend.log
```

**Network addresses:**
- Local: http://localhost:3000/
- Network: http://192.168.1.65:3000/

---

### ✅ PostgreSQL
```
Container: homework_postgres
Image: postgres:15-alpine
Port: 5432
Status: Healthy
Tables: 24 таблицы
```

---

### ✅ Redis
```
Container: homework_redis
Image: redis:7-alpine
Port: 6379
Status: Healthy
```

---

## Быстрый доступ

| Сервис | URL | Описание |
|--------|-----|----------|
| **Frontend** | http://localhost:3000/ | React приложение |
| **Backend API** | http://localhost:8080 | REST API |
| **Health Check** | http://localhost:8080/health | Мониторинг |

---

## Конфигурация Frontend

### Environment Variables (.env)
```bash
VITE_API_BASE_URL=http://localhost:8080
VITE_ANALYTICS_ENABLED=true
VITE_ANALYTICS_DEBUG=true
VITE_FEATURE_VILLAIN=true
VITE_FEATURE_ACHIEVEMENTS=true
VITE_FEATURE_REFERRALS=true
VITE_APP_VERSION=1.0.0
VITE_APP_NAME=Homework Helper
```

### API Integration
Frontend настроен на обращение к backend по адресу: `http://localhost:8080`

---

## Тестирование Full Stack

### 1. Проверка Backend Health
```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

### 2. Открыть Frontend
```bash
# В браузере открыть:
http://localhost:3000/

# Должна открыться React приложение без проблем с SSL
```

### 3. Проверка API из Frontend
Откройте Developer Tools (F12) в браузере и проверьте Network вкладку.
Все запросы к API должны идти на `http://localhost:8080`.

---

## Управление

### Остановить всё
```bash
# Остановить frontend
kill $(ps aux | grep "node.*vite" | grep -v grep | awk '{print $2}')

# Остановить backend
kill $(ps aux | grep "go run ./cmd/server" | grep -v grep | awk '{print $2}')

# Остановить PostgreSQL и Redis
docker compose down
```

### Перезапустить Frontend
```bash
cd /Users/a.yanover/Xsolla/child_bot/frontend
npm run dev > /tmp/frontend.log 2>&1 &
```

### Перезапустить Backend
```bash
cd /Users/a.yanover/Xsolla/child_bot/api
DATABASE_URL="postgres://homework:root@localhost:5432/homework?sslmode=disable" \
LLM_SERVER_URL="http://138.124.55.145:8000" \
DEFAULT_LLM="gemini" \
PORT="8080" \
ENV="development" \
LOG_LEVEL="info" \
ALLOWED_ORIGINS="http://localhost:5173,http://localhost:3000" \
go run ./cmd/server/ > /tmp/backend.log 2>&1 &
```

### Просмотр логов
```bash
# Frontend логи
tail -f /tmp/frontend.log

# Backend логи
tail -f /tmp/backend.log

# PostgreSQL логи
docker logs -f homework_postgres
```

---

## Проверка работы

### Запущенные процессы
```bash
ps aux | grep -E "(go run|vite)" | grep -v grep

# Должно показать:
# - go run ./cmd/server/  (Backend)
# - node .../vite         (Frontend)
```

### Порты
```bash
lsof -i :8080 -i :3000 | grep LISTEN

# Должно показать:
# - server (port 8080) - Backend
# - node (port 3000)   - Frontend
```

### Docker контейнеры
```bash
docker compose ps

# Должно показать:
# - homework_postgres (Up, healthy)
# - homework_redis    (Up, healthy)
```

---

## Разработка

### Hot Reload
Оба сервиса поддерживают hot reload:
- **Frontend**: Vite автоматически перезагружает при изменении файлов в `frontend/src/`
- **Backend**: При изменении `.go` файлов нужен перезапуск (можно использовать `air` для auto-reload)

### Отладка

#### Frontend
1. Откройте https://localhost:3000/
2. Нажмите F12 (Developer Tools)
3. Проверьте Console и Network вкладки

#### Backend
1. Проверьте логи: `tail -f /tmp/backend.log`
2. Используйте `curl` для тестирования endpoints
3. Добавьте `LOG_LEVEL=debug` для детальных логов

---

## Архитектура

```
┌─────────────────────┐
│   Browser           │
│   https://localhost │
│   :3000             │
└──────────┬──────────┘
           │
           │ HTTP Requests
           ▼
┌─────────────────────┐
│   Frontend (Vite)   │
│   React + TypeScript│
│   Port: 3000        │
└──────────┬──────────┘
           │
           │ API Calls
           │ http://localhost:8080
           ▼
┌─────────────────────┐
│   Backend (Go)      │
│   REST API          │
│   Port: 8080        │
└──────┬──────────────┘
       │
       ├──────────────────┐
       │                  │
       ▼                  ▼
┌─────────────┐    ┌─────────────┐
│  PostgreSQL │    │    Redis    │
│  Port: 5432 │    │  Port: 6379 │
└─────────────┘    └─────────────┘
       │
       │
       ▼
┌─────────────────────┐
│   LLM Server        │
│   (Remote)          │
│   138.124.55.145    │
└─────────────────────┘
```

---

## HTTP vs HTTPS в Development

### Локальная разработка (текущая настройка)
Frontend использует **HTTP** без SSL для простоты разработки.

### Production
В production HTTPS будет настроен через:
- Nginx reverse proxy с Let's Encrypt сертификатом
- Или CloudFlare для автоматического HTTPS

### Если нужен HTTPS локально
Измените в `vite.config.ts`:
```typescript
server: {
  https: true, // Включить HTTPS
}
```
Но для большинства задач разработки HTTP достаточно.

---

## Next Steps

### 1. Тестирование интеграции
- Проверить все 34 API endpoints через frontend
- Протестировать authentication flow
- Проверить analytics events

### 2. Добавить фичи
- Onboarding flow
- Help и Check функциональность
- Achievements и Villains
- Referral система

### 3. Оптимизация
- Настроить CORS для production
- Добавить rate limiting
- Настроить caching
- Добавить monitoring

---

## Troubleshooting

### Frontend не открывается
```bash
# Проверить что процесс запущен
ps aux | grep vite

# Проверить логи
tail -f /tmp/frontend.log

# Перезапустить
cd frontend && npm run dev
```

### Backend возвращает 500
```bash
# Проверить логи
tail -f /tmp/backend.log

# Проверить подключение к БД
psql "postgres://homework:root@localhost:5432/homework" -c "SELECT 1"
```

### CORS ошибки
Убедитесь что в backend настроен правильный `ALLOWED_ORIGINS`:
```bash
export ALLOWED_ORIGINS="http://localhost:5173,http://localhost:3000"
```

---

## Статус сервисов

| Сервис | PID | Port | Status |
|--------|-----|------|--------|
| Frontend (Vite) | 95093 | 3000 | ✅ Running |
| Backend (Go) | 92786 | 8080 | ✅ Running |
| PostgreSQL | Docker | 5432 | ✅ Healthy |
| Redis | Docker | 6379 | ✅ Healthy |

---

**Полный стек запущен и готов к разработке!** 🚀

**Дата:** 31 марта 2026
**Время:** 15:59
