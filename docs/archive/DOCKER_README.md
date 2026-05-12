# 🐳 Docker Setup - Quick Start

## Быстрый старт за 3 шага

### 1️⃣ Подготовка

```bash
# Клонировать репозиторий
git clone <your-repo-url>
cd child_bot

# Создать .env файл
cp .env.example .env

# Отредактировать .env (обязательно!)
nano .env
```

**Минимальные переменные для запуска:**
- `POSTGRES_PASSWORD` - пароль для БД
- `REDIS_PASSWORD` - пароль для Redis
- `JWT_SECRET` - секретный ключ для JWT

### 2️⃣ Запуск

```bash
# Вариант A: Использовать Makefile (рекомендуется)
make install    # Первый запуск
make up         # Последующие запуски

# Вариант B: Docker Compose напрямую
docker-compose up -d
```

### 3️⃣ Проверка

```bash
# Проверить статус
make ps

# Или
docker-compose ps

# Открыть в браузере
open http://localhost       # Frontend
open http://localhost:8080  # Backend API
```

---

## 📦 Структура контейнеров

| Сервис | Порт | Описание |
|--------|------|----------|
| **frontend** | 80 | React + Vite + Nginx |
| **backend** | 8080 | Go API |
| **postgres** | 5432 | PostgreSQL 15 |
| **redis** | 6379 | Redis Cache |

---

## 🚀 Полезные команды (Makefile)

```bash
make help           # Показать все команды
make up             # Запустить все сервисы
make down           # Остановить все сервисы
make restart        # Перезапустить
make logs           # Показать логи
make ps             # Статус контейнеров
make health         # Проверить здоровье сервисов
make clean          # Остановить и удалить все данные
make rebuild        # Полная пересборка
```

### Работа с БД

```bash
make migrate-up     # Применить миграции
make migrate-down   # Откатить миграцию
make backup-db      # Создать backup
make restore-db FILE=backups/backup.sql.gz  # Восстановить
make shell-postgres # Открыть psql
```

### Разработка

```bash
make dev            # Режим разработки (hot reload)
make test-frontend  # Тесты frontend
make test-backend   # Тесты backend
```

---

## 🔧 Конфигурация

### Environment Variables

Основные переменные в `.env`:

```bash
# Database
POSTGRES_DB=homework
POSTGRES_USER=homework
POSTGRES_PASSWORD=your_secure_password

# Redis
REDIS_PASSWORD=your_redis_password

# Backend
JWT_SECRET=your_jwt_secret_key
OPENAI_API_KEY=your_openai_key
TELEGRAM_BOT_TOKEN=your_telegram_token

# Frontend
VITE_API_BASE_URL=/api
```

### Порты (настраиваемые)

```bash
FRONTEND_PORT=80
BACKEND_PORT=8080
POSTGRES_PORT=5432
REDIS_PORT=6379
```

---

## 🏗️ Архитектура

```
Frontend (Nginx) :80
    ↓ proxy /api/* → Backend :8080
    ↓                   ↓
    ↓              PostgreSQL :5432
    ↓              Redis :6379
```

**Особенности:**
- Frontend в Nginx контейнере (оптимизирован для production)
- Backend с автоматическими миграциями БД
- Health checks для всех сервисов
- Persistent volumes для данных
- Internal network для безопасности

---

## 📝 Логи и мониторинг

```bash
# Все логи
docker-compose logs -f

# Конкретный сервис
docker-compose logs -f backend

# Последние 100 строк
docker-compose logs --tail=100 frontend

# Статистика контейнеров
docker stats
```

---

## 🔐 Security Checklist

- [ ] Изменить `POSTGRES_PASSWORD` в `.env`
- [ ] Изменить `REDIS_PASSWORD` в `.env`
- [ ] Изменить `JWT_SECRET` в `.env`
- [ ] Не коммитить `.env` в Git
- [ ] Использовать SSL в production
- [ ] Настроить firewall (UFW/iptables)
- [ ] Регулярно обновлять Docker образы

---

## 🐛 Troubleshooting

### Проблема: "port already in use"

```bash
# Проверить какой процесс использует порт
lsof -i :80
lsof -i :8080

# Остановить конфликтующий сервис или изменить порт в .env
FRONTEND_PORT=8081 docker-compose up -d
```

### Проблема: "database connection failed"

```bash
# Проверить статус Postgres
docker-compose exec postgres pg_isready

# Проверить логи
docker-compose logs postgres

# Пересоздать контейнер
docker-compose down
docker-compose up -d
```

### Проблема: "out of disk space"

```bash
# Очистить неиспользуемые ресурсы
docker system prune -a --volumes

# Проверить использование
docker system df
```

---

## 📚 Дополнительная документация

- [`DOCKER_DEPLOYMENT.md`](./DOCKER_DEPLOYMENT.md) - Полное руководство по deployment
- [`.env.example`](./.env.example) - Пример конфигурации
- [`frontend/nginx.conf`](./frontend/nginx.conf) - Конфигурация Nginx
- [`docker-compose.yml`](./docker-compose.yml) - Orchestration

---

## 🎯 Production Deployment

Для production развёртывания:

1. Настроить SSL сертификаты (Let's Encrypt)
2. Использовать внешний managed PostgreSQL и Redis
3. Настроить CI/CD pipeline (GitHub Actions)
4. Добавить мониторинг (Prometheus + Grafana)
5. Настроить backup стратегию

Подробнее в [`DOCKER_DEPLOYMENT.md`](./DOCKER_DEPLOYMENT.md)

---

## 💡 Tips & Tricks

```bash
# Быстрый рестарт только frontend
docker-compose restart frontend

# Выполнить команду в контейнере
docker-compose exec backend sh

# Посмотреть переменные окружения
docker-compose exec backend env

# Проверить сеть
docker network inspect child_bot_homework_network

# Скопировать файлы из контейнера
docker cp homework_backend:/app/logs ./logs
```

---

## 🤝 Support

Проблемы с Docker setup? Открывай issue в GitHub!

- Docker Documentation: https://docs.docker.com/
- Docker Compose: https://docs.docker.com/compose/
