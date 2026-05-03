# Deployment Changes Log

## 2026-04-27 - Analytics & Auth Fix

### Backend Changes
- **auth.go**: Добавлен `/api/analytics/events` в whitelist путей, не требующих `X-Child-Profile-ID`
  - Путь: `api/internal/api/middleware/auth.go`
  - Изменение: добавлена строка `/api/analytics/events` в функцию `requiresChildProfile()`

### Environment Variables
- **POSTGRES_PASSWORD**: `0xVQhgOz8E7EPVQZu0E5x7ZzixdwoX5d` (обновлено - убраны спецсимволы)
- **REDIS_PASSWORD**: `G8uESoxbJxZdiNnTJbwlSYtAu1xPGn37` (обновлено - убраны спецсимволы)
- **VITE_API_BASE_URL**: `/` (для production, т.к. пути в routes.ts уже содержат /api/)

### Frontend Changes
- **ReportPage.tsx**: Изменен `API_BASE_URL` на использование `config.isProduction`
  - Development: `http://localhost:8080`
  - Production: `` (пустая строка для относительных путей)

- **.env.production** (frontend):
  ```env
  VITE_API_BASE_URL=/
  ```

### Docker Compose
- Обновлен до `docker-compose.production.yml`
- Backend использует образ из GHCR: `ghcr.io/knopka87/child_bot-backend:latest`
- Добавлен проброс портов `8080:8080` для backend
- Healthcheck обновлен на `/api/health` с исправленными параметрами

### Deployment Process
1. Обновить код на GitHub (push в ветку `prod-v1`)
2. GitHub Actions автоматически соберет и опубликует образы в GHCR
3. На сервере:
   ```bash
   cd /root/child_bot
   docker compose pull backend
   docker compose up -d backend
   ```
4. Для frontend (ручная сборка):
   ```bash
   # Локально
   cd frontend
   VITE_API_BASE_URL=/ npm run build

   # Загрузка на сервер
   scp -r dist root@77.222.60.149:/root/child_bot/frontend/dist.new
   ssh root@77.222.60.149 "rm -rf /root/child_bot/frontend/dist && mv /root/child_bot/frontend/dist.new /root/child_bot/frontend/dist && nginx -s reload"
   ```

### Server Status
- Backend: ✅ Running (port 8080)
- PostgreSQL: ✅ Healthy
- Redis: ✅ Healthy
- Nginx: ✅ Proxying /api/ → backend:8080
- Frontend: ✅ Serving from /root/child_bot/frontend/dist

### URLs
- Production: https://vk.obyasnyatel.ru
- VK Mini App: https://vk.com/app54517931
