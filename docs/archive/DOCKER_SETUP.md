# Docker Setup Guide

Инструкции по запуску проекта через Docker Compose.

---

## 📁 Структура Docker конфигураций

```
child_bot/
├── docker-compose.yml              # Production конфигурация
├── docker/
│   └── docker-compose.dev.yml     # Development конфигурация
├── .env.example                    # Шаблон для development
├── .env.production.example         # Шаблон для production
└── .env                            # Ваш локальный .env (не коммитится)
```

---

## 🚀 Quick Start

### Development (рекомендуется для локальной разработки)

```bash
# 1. Скопировать .env
cp .env.example .env

# 2. Отредактировать .env (опционально)
# Для development значения по умолчанию подходят

# 3. Запустить development окружение
docker compose -f docker/docker-compose.dev.yml up -d

# 4. Посмотреть логи
docker compose -f docker/docker-compose.dev.yml logs -f

# 5. Остановить
docker compose -f docker/docker-compose.dev.yml down
```

**Доступ:**
- Frontend: http://localhost:5173 (Vite dev server с hot reload)
- Backend: http://localhost:8080
- PostgreSQL: localhost:5432
- Redis: localhost:6379

### Production (для production деплоя)

```bash
# 1. Скопировать production шаблон
cp .env.production.example .env

# 2. ОБЯЗАТЕЛЬНО изменить пароли и секреты в .env!
nano .env

# 3. Запустить production окружение
docker compose up -d

# 4. Посмотреть логи
docker compose logs -f

# 5. Остановить
docker compose down
```

**Доступ:**
- Frontend: http://localhost:80 (Nginx с production build)
- Backend: http://localhost:8080
- PostgreSQL: localhost:5432
- Redis: localhost:6379

---

## 🔧 Конфигурации

### Development (`docker/docker-compose.dev.yml`)

**Особенности:**
- ✅ Frontend: Vite dev server с hot reload
- ✅ Backend: автоматическое применение миграций через entrypoint.sh
- ✅ Volumes для hot reload кода
- ✅ Быстрые healthchecks (5s interval)
- ✅ Debug логирование
- ✅ CORS разрешён для localhost

**Сервисы:**
- `postgres` → `child_bot_postgres_dev`
- `redis` → `child_bot_redis_dev`
- `backend` → `child_bot_backend_dev`
- `frontend` → `child_bot_frontend_dev`

**Порты:**
- Frontend: 5173
- Backend: 8080
- PostgreSQL: 5432
- Redis: 6379

### Production (`docker-compose.yml`)

**Особенности:**
- ✅ Frontend: Nginx с optimized production build
- ✅ Backend: компилированный Go binary
- ✅ Автоматические миграции при старте
- ✅ Healthchecks для reliability
- ✅ Restart policies
- ✅ Требует настройки .env с реальными credentials

**Сервисы:**
- `postgres` → `child_bot_postgres`
- `redis` → `child_bot_redis`
- `backend` → `child_bot_backend`
- `frontend` → `child_bot_frontend`

**Порты:**
- Frontend: 80
- Backend: 8080
- PostgreSQL: 5432
- Redis: 6379

---

## 📋 Переменные окружения

### Критичные для production:

```bash
# ОБЯЗАТЕЛЬНО изменить!
POSTGRES_PASSWORD=your_strong_password_here
REDIS_PASSWORD=your_strong_redis_password_here

# Настроить URLs
LLM_SERVER_URL=http://your-llm-server:8000
APP_URL=https://your-domain.com
ALLOWED_ORIGINS=https://your-domain.com,https://vk.com

# VK App ID
VITE_VK_APP_ID=your_vk_app_id

# Legal URLs
VITE_PRIVACY_POLICY_URL=https://your-domain.com/legal/privacy
VITE_TERMS_OF_SERVICE_URL=https://your-domain.com/legal/terms
```

### Опциональные:

```bash
# Внешние сервисы
SENTRY_DSN=https://your-sentry-dsn
AMPLITUDE_API_KEY=your_amplitude_key
VK_SECRET_KEY=your_vk_secret_key
```

---

## 🔨 Полезные команды

### Development

```bash
# Запуск только БД и Redis
docker compose -f docker/docker-compose.dev.yml up -d postgres redis

# Пересоздать контейнеры (после изменения Dockerfile)
docker compose -f docker/docker-compose.dev.yml up -d --build

# Посмотреть логи конкретного сервиса
docker compose -f docker/docker-compose.dev.yml logs -f backend
docker compose -f docker/docker-compose.dev.yml logs -f frontend

# Зайти внутрь контейнера
docker compose -f docker/docker-compose.dev.yml exec backend sh
docker compose -f docker/docker-compose.dev.yml exec postgres psql -U child_bot -d child_bot

# Остановить и удалить volumes (ОСТОРОЖНО: удаляет данные!)
docker compose -f docker/docker-compose.dev.yml down -v
```

### Production

```bash
# Запуск с пересборкой
docker compose up -d --build

# Перезапуск конкретного сервиса
docker compose restart backend

# Просмотр статуса
docker compose ps

# Остановка без удаления volumes
docker compose down

# Полная очистка (удалить volumes)
docker compose down -v
```

### Миграции

```bash
# Development (автоматически применяются при старте backend)
# Если нужно применить вручную:
docker compose -f docker/docker-compose.dev.yml exec backend /app/entrypoint.sh

# Production (автоматически применяются при старте backend)
# Или вручную:
docker compose exec backend migrate -source "file:///app/migrations" \
  -database "$DATABASE_URL" up
```

### Логи

```bash
# Все сервисы
docker compose -f docker/docker-compose.dev.yml logs -f

# Только backend
docker compose -f docker/docker-compose.dev.yml logs -f backend

# Только errors
docker compose -f docker/docker-compose.dev.yml logs -f | grep -i error

# Последние 100 строк
docker compose -f docker/docker-compose.dev.yml logs --tail=100
```

---

## 🏗️ Архитектура

### Development Flow
```
┌─────────────┐     ┌──────────────┐     ┌────────────┐
│   Browser   │────▶│ Vite Dev     │────▶│  Backend   │
│             │     │ (port 5173)  │     │ (port 8080)│
└─────────────┘     └──────────────┘     └─────┬──────┘
                                                │
                                    ┌───────────┴────────┐
                                    │                    │
                              ┌─────▼─────┐      ┌──────▼──────┐
                              │ PostgreSQL│      │    Redis    │
                              │ (port 5432)      │ (port 6379) │
                              └───────────┘      └─────────────┘
```

### Production Flow
```
┌─────────────┐     ┌──────────────┐     ┌────────────┐
│   Browser   │────▶│    Nginx     │────▶│  Backend   │
│             │     │  (port 80)   │     │ (port 8080)│
└─────────────┘     └──────────────┘     └─────┬──────┘
                                                │
                                    ┌───────────┴────────┐
                                    │                    │
                              ┌─────▼─────┐      ┌──────▼──────┐
                              │ PostgreSQL│      │    Redis    │
                              │ (port 5432)      │ (port 6379) │
                              └───────────┘      └─────────────┘
```

---

## 🐛 Troubleshooting

### Порты заняты

```bash
# Проверить что занимает порт
lsof -i :5173
lsof -i :8080
lsof -i :5432

# Изменить порты в .env
FRONTEND_PORT=5174
BACKEND_PORT=8081
POSTGRES_PORT=5433
```

### База данных не запускается

```bash
# Проверить логи
docker compose -f docker/docker-compose.dev.yml logs postgres

# Удалить volume и пересоздать
docker compose -f docker/docker-compose.dev.yml down -v
docker compose -f docker/docker-compose.dev.yml up -d
```

### Backend не подключается к БД

```bash
# Проверить что БД готова
docker compose -f docker/docker-compose.dev.yml exec postgres pg_isready

# Проверить DATABASE_URL
docker compose -f docker/docker-compose.dev.yml exec backend env | grep DATABASE_URL

# Проверить сеть
docker network ls
docker network inspect child_bot_child_bot_network
```

### Миграции не применяются

```bash
# Проверить логи backend при старте
docker compose -f docker/docker-compose.dev.yml logs backend | grep -i migration

# Применить вручную
docker compose -f docker/docker-compose.dev.yml exec backend \
  migrate -source "file:///app/migrations" -database "$DATABASE_URL" up

# Проверить версию миграций
docker compose -f docker/docker-compose.dev.yml exec backend \
  migrate -source "file:///app/migrations" -database "$DATABASE_URL" version
```

### Frontend не собирается

```bash
# Проверить node_modules
docker compose -f docker/docker-compose.dev.yml exec frontend ls -la /app/node_modules

# Пересобрать
docker compose -f docker/docker-compose.dev.yml down
docker compose -f docker/docker-compose.dev.yml up -d --build frontend
```

---

## 📊 Мониторинг

### Healthchecks

```bash
# Проверить статус всех сервисов
docker compose -f docker/docker-compose.dev.yml ps

# Backend health
curl http://localhost:8080/health

# Frontend (production)
curl http://localhost/

# PostgreSQL
docker compose -f docker/docker-compose.dev.yml exec postgres pg_isready -U child_bot
```

### Ресурсы

```bash
# Использование ресурсов контейнерами
docker stats

# Размер volumes
docker system df -v
```

---

## 🔒 Безопасность

### Production Checklist

- [ ] Изменены все пароли в `.env`
- [ ] `POSTGRES_PASSWORD` != `dev_secret`
- [ ] `REDIS_PASSWORD` != `dev_redis_secret`
- [ ] `ALLOWED_ORIGINS` содержит только реальные домены
- [ ] `VITE_API_BASE_URL` указывает на production API
- [ ] SSL сертификаты настроены (если нужен HTTPS)
- [ ] Firewall правила настроены
- [ ] `.env` файл не коммитится в git
- [ ] Backup стратегия для PostgreSQL настроена

---

## 📚 Дополнительные ресурсы

- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [PostgreSQL Docker Image](https://hub.docker.com/_/postgres)
- [Redis Docker Image](https://hub.docker.com/_/redis)
- [Nginx Docker Image](https://hub.docker.com/_/nginx)

---

**Последнее обновление:** 2026-04-04
**Версия:** 1.0
