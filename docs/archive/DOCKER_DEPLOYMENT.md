# Docker Deployment Guide

Полное руководство по развёртыванию приложения с использованием Docker.

---

## Архитектура

```
┌─────────────────────────────────────────────────────────┐
│                    Docker Network                        │
│                                                          │
│  ┌──────────┐      ┌──────────┐      ┌──────────┐     │
│  │          │      │          │      │          │     │
│  │ Frontend │─────▶│ Backend  │─────▶│ Postgres │     │
│  │  (Nginx) │      │   (Go)   │      │          │     │
│  │          │      │          │      │          │     │
│  └────┬─────┘      └────┬─────┘      └──────────┘     │
│       │                 │                              │
│       │                 │            ┌──────────┐     │
│       │                 └───────────▶│  Redis   │     │
│       │                              │          │     │
│  Port 80                        Port 8080       │     │
│                                                  │     │
└──────────────────────────────────────────────────────┘
```

### Сервисы

1. **Frontend (React + Vite + Nginx)**
   - Port: 80
   - Статические файлы
   - Proxy для API запросов
   - SPA routing

2. **Backend (Go API)**
   - Port: 8080
   - REST API
   - WebSocket support
   - Database migrations

3. **PostgreSQL**
   - Port: 5432
   - Persistent storage
   - Health checks

4. **Redis**
   - Port: 6379
   - Кэширование
   - Session storage

---

## Quick Start

### 1. Подготовка

```bash
# Клонируем репозиторий
git clone <repo-url>
cd child_bot

# Копируем .env файл
cp .env.example .env

# Редактируем .env и заполняем необходимые значения
nano .env
```

### 2. Сборка и запуск

```bash
# Сборка всех контейнеров
docker-compose build

# Запуск в фоновом режиме
docker-compose up -d

# Проверка статуса
docker-compose ps

# Просмотр логов
docker-compose logs -f
```

### 3. Проверка работоспособности

```bash
# Frontend health check
curl http://localhost/health

# Backend health check
curl http://localhost:8080/health

# Полная проверка
docker-compose ps
```

---

## Команды Docker Compose

### Основные команды

```bash
# Запуск всех сервисов
docker-compose up -d

# Остановка всех сервисов
docker-compose down

# Остановка с удалением volumes
docker-compose down -v

# Перезапуск сервиса
docker-compose restart <service-name>

# Пересборка и запуск
docker-compose up -d --build
```

### Работа с отдельными сервисами

```bash
# Запуск только backend
docker-compose up -d backend

# Логи frontend
docker-compose logs -f frontend

# Перезапуск postgres
docker-compose restart postgres

# Выполнить команду в контейнере
docker-compose exec backend sh
docker-compose exec postgres psql -U homework -d homework
```

### Логи и мониторинг

```bash
# Все логи
docker-compose logs -f

# Логи конкретного сервиса
docker-compose logs -f backend

# Последние 100 строк
docker-compose logs --tail=100 frontend

# Статистика контейнеров
docker stats

# Использование дискового пространства
docker system df
```

---

## Environment Variables

### Обязательные переменные

```bash
# Database
POSTGRES_PASSWORD=strong_password_here

# Redis
REDIS_PASSWORD=strong_redis_password

# JWT
JWT_SECRET=very_long_random_secret_key

# Telegram
TELEGRAM_BOT_TOKEN=your_telegram_bot_token

# OpenAI
OPENAI_API_KEY=your_openai_api_key
```

### Опциональные переменные

```bash
# Порты (по умолчанию)
FRONTEND_PORT=80
BACKEND_PORT=8080
POSTGRES_PORT=5432
REDIS_PORT=6379

# Окружение
ENV=production  # development, staging, production
LOG_LEVEL=info  # debug, info, warn, error

# Features
VITE_ENABLE_ANALYTICS=true
```

---

## Database Migrations

### Автоматические миграции

При запуске backend контейнер автоматически применяет миграции через `entrypoint.sh`:

```bash
# entrypoint.sh делает:
migrate -path /app/migrations -database $DATABASE_URL up
```

### Ручные миграции

```bash
# Применить миграции
docker-compose exec backend migrate -path /app/migrations -database $DATABASE_URL up

# Откатить последнюю миграцию
docker-compose exec backend migrate -path /app/migrations -database $DATABASE_URL down 1

# Проверить версию
docker-compose exec backend migrate -path /app/migrations -database $DATABASE_URL version

# Создать новую миграцию (на хосте)
migrate create -ext sql -dir api/migrations -seq add_new_table
```

---

## Backup & Restore

### Backup PostgreSQL

```bash
# Создать backup
docker-compose exec postgres pg_dump -U homework homework > backup_$(date +%Y%m%d_%H%M%S).sql

# Или с docker-compose
docker-compose exec -T postgres pg_dump -U homework homework | gzip > backup_$(date +%Y%m%d_%H%M%S).sql.gz
```

### Restore PostgreSQL

```bash
# Восстановить из backup
cat backup_20240330_120000.sql | docker-compose exec -T postgres psql -U homework -d homework

# Или из gzip
gunzip < backup_20240330_120000.sql.gz | docker-compose exec -T postgres psql -U homework -d homework
```

### Backup Redis

```bash
# Создать snapshot
docker-compose exec redis redis-cli -a $REDIS_PASSWORD SAVE

# Скопировать dump.rdb
docker cp homework_redis:/data/dump.rdb ./redis_backup_$(date +%Y%m%d).rdb
```

---

## Production Deployment

### 1. Подготовка сервера

```bash
# Установка Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# Установка Docker Compose
apt-get update
apt-get install docker-compose-plugin

# Создание пользователя
useradd -m -s /bin/bash homework
usermod -aG docker homework
```

### 2. Настройка firewall

```bash
# UFW
ufw allow 80/tcp    # HTTP
ufw allow 443/tcp   # HTTPS
ufw allow 22/tcp    # SSH
ufw enable

# Проверка
ufw status
```

### 3. SSL сертификат (Let's Encrypt)

```bash
# Установка certbot
apt-get install certbot

# Получение сертификата
certbot certonly --standalone -d your-domain.com

# Автоматическое обновление
crontab -e
# Добавить: 0 0 * * * certbot renew --quiet
```

### 4. Nginx с SSL (опционально, вместо Docker nginx)

```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;

    location / {
        proxy_pass http://localhost:80;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 5. Systemd service

```bash
# Создать файл: /etc/systemd/system/homework.service
[Unit]
Description=Homework App
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/home/homework/child_bot
ExecStart=/usr/bin/docker-compose up -d
ExecStop=/usr/bin/docker-compose down
TimeoutStartSec=0

[Install]
WantedBy=multi-user.target

# Активация
systemctl enable homework
systemctl start homework
```

---

## Monitoring & Logging

### Health Checks

```bash
# Frontend
curl http://localhost/health

# Backend
curl http://localhost:8080/health

# Postgres
docker-compose exec postgres pg_isready -U homework

# Redis
docker-compose exec redis redis-cli -a $REDIS_PASSWORD ping
```

### Логирование

```bash
# Все логи в файл
docker-compose logs -f > logs_$(date +%Y%m%d).log

# Логи с timestamp
docker-compose logs -f --timestamps

# Фильтрация по сервису
docker-compose logs -f backend | grep ERROR
```

### Prometheus + Grafana (опционально)

Добавить в `docker-compose.yml`:

```yaml
  prometheus:
    image: prom/prometheus
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    ports:
      - "9090:9090"

  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
```

---

## Troubleshooting

### Проблема: Контейнер не запускается

```bash
# Проверить логи
docker-compose logs <service-name>

# Проверить статус
docker-compose ps

# Пересоздать контейнер
docker-compose up -d --force-recreate <service-name>
```

### Проблема: Database connection failed

```bash
# Проверить доступность Postgres
docker-compose exec postgres pg_isready

# Проверить переменные окружения
docker-compose exec backend env | grep POSTGRES

# Проверить сеть
docker network inspect child_bot_homework_network
```

### Проблема: Out of memory

```bash
# Проверить использование памяти
docker stats

# Добавить memory limits в docker-compose.yml:
services:
  backend:
    mem_limit: 512m
    mem_reservation: 256m
```

### Проблема: Disk space full

```bash
# Очистить неиспользуемые образы
docker image prune -a

# Очистить volumes
docker volume prune

# Полная очистка
docker system prune -a --volumes
```

---

## Security Best Practices

1. **Secrets Management**
   - Использовать `.env` файл (не коммитить в Git)
   - Docker secrets для production
   - Rotate passwords регулярно

2. **Network Security**
   - Использовать internal networks
   - Ограничить expose портов
   - Настроить firewall

3. **Container Security**
   - Регулярно обновлять образы
   - Использовать non-root пользователей
   - Scan образов на уязвимости

4. **SSL/TLS**
   - Всегда использовать HTTPS
   - Настроить SSL certificates
   - Force HTTPS redirect

---

## Performance Optimization

### 1. Nginx caching

```nginx
# Добавить в nginx.conf
proxy_cache_path /var/cache/nginx levels=1:2 keys_zone=api_cache:10m max_size=100m;

location /api/ {
    proxy_cache api_cache;
    proxy_cache_valid 200 10m;
    proxy_cache_use_stale error timeout updating;
}
```

### 2. Database connection pooling

```yaml
backend:
  environment:
    DB_MAX_CONNECTIONS: 20
    DB_IDLE_CONNECTIONS: 5
```

### 3. Redis optimization

```yaml
redis:
  command: redis-server --maxmemory 256mb --maxmemory-policy allkeys-lru
```

---

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Copy files to server
        uses: appleboy/scp-action@master
        with:
          host: ${{ secrets.SERVER_HOST }}
          username: ${{ secrets.SERVER_USER }}
          key: ${{ secrets.SSH_KEY }}
          source: "."
          target: "/home/homework/child_bot"

      - name: Deploy with Docker Compose
        uses: appleboy/ssh-action@master
        with:
          host: ${{ secrets.SERVER_HOST }}
          username: ${{ secrets.SERVER_USER }}
          key: ${{ secrets.SSH_KEY }}
          script: |
            cd /home/homework/child_bot
            docker-compose pull
            docker-compose up -d --build
            docker-compose exec backend migrate -path /app/migrations -database $DATABASE_URL up
```

---

## Support

- Docker Hub: https://hub.docker.com/
- Docker Docs: https://docs.docker.com/
- Docker Compose: https://docs.docker.com/compose/

---

## Changelog

### v1.0.0 (2024-03-30)
- Initial Docker setup
- Frontend + Backend + Postgres + Redis
- Health checks
- Auto migrations
- Production-ready configuration
