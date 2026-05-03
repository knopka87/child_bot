# Инструкция по деплою на сервер

## Архитектура

**Сервер:** `77.222.60.149`
**Домен:** `https://vk.obyasnyatel.ru`

### Компоненты:
- **Nginx** - работает напрямую на сервере (порты 80/443), раздает статику фронтенда
- **Backend API** - Docker контейнер на порту 8080
- **PostgreSQL** - Docker контейнер
- **Redis** - Docker контейнер
- **Frontend** - статические файлы в `/root/child_bot/frontend/dist/`

### Важно:
- Frontend НЕ работает в контейнере, Nginx раздает файлы напрямую из директории
- Backend работает в контейнере и доступен через Nginx reverse proxy на `/api/*`

---

## Предварительные требования

### 1. SSH ключ
```bash
# Проверьте наличие SSH ключа
ls -la /Users/a.yanover/Downloads/id_rsa_1/id_rsa

# Установите права (если нужно)
chmod 600 /Users/a.yanover/Downloads/id_rsa_1/id_rsa
```

### 2. GitHub Container Registry
Убедитесь, что образы собраны в GitHub Actions:
- `ghcr.io/knopka87/child_bot-backend:latest`
- `ghcr.io/knopka87/child_bot-frontend:latest`

Проверить можно здесь: https://github.com/knopka87/child_bot/pkgs/container/child_bot-backend

---

## Деплой Backend

### Шаг 1: Подключиться к серверу
```bash
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149
```

### Шаг 2: Обновить образ и перезапустить контейнер
```bash
cd /root/child_bot
docker compose pull backend
docker compose up -d backend
```

### Шаг 3: Проверить статус
```bash
# Проверить статус контейнера
docker compose ps

# Посмотреть логи
docker compose logs -f backend --tail=50

# Проверить health check
curl http://localhost:8080/api/health
```

### Быстрая команда (одной строкой с локальной машины):
```bash
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149 \
  "cd /root/child_bot && docker compose pull backend && docker compose up -d backend && docker compose logs backend --tail=30"
```

---

## Деплой Frontend

### Шаг 1: Загрузить образ и скопировать файлы
```bash
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149 \
  "cd /root/child_bot && \
   docker compose pull frontend && \
   docker create --name temp_frontend ghcr.io/knopka87/child_bot-frontend:latest && \
   docker cp temp_frontend:/usr/share/nginx/html/. /root/child_bot/frontend/dist/ && \
   docker rm temp_frontend"
```

### Шаг 2: Перезагрузить Nginx
```bash
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149 \
  "nginx -t && nginx -s reload || (killall nginx; sleep 2; systemctl start nginx)"
```

### Шаг 3: Проверить деплой
```bash
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149 \
  "ls -lh /root/child_bot/frontend/dist/index.html && \
   curl -I https://vk.obyasnyatel.ru"
```

### Быстрая команда (всё в одной строке):
```bash
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149 \
  "cd /root/child_bot && \
   docker compose pull frontend && \
   docker create --name temp_frontend ghcr.io/knopka87/child_bot-frontend:latest && \
   docker cp temp_frontend:/usr/share/nginx/html/. /root/child_bot/frontend/dist/ && \
   docker rm temp_frontend && \
   nginx -s reload && \
   echo '✅ Frontend deployed successfully'"
```

---

## Деплой Frontend + Backend одновременно

```bash
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149 << 'EOF'
cd /root/child_bot

# Обновить бекенд
echo "📦 Updating backend..."
docker compose pull backend
docker compose up -d backend

# Обновить фронтенд
echo "📦 Updating frontend..."
docker compose pull frontend
docker create --name temp_frontend ghcr.io/knopka87/child_bot-frontend:latest
docker cp temp_frontend:/usr/share/nginx/html/. /root/child_bot/frontend/dist/
docker rm temp_frontend

# Перезагрузить Nginx
echo "🔄 Reloading Nginx..."
nginx -s reload || (killall nginx && sleep 2 && systemctl start nginx)

# Проверить статусы
echo "✅ Deployment complete!"
docker compose ps
systemctl status nginx --no-pager -l
EOF
```

---

## Проверка работоспособности

### 1. Проверить все сервисы
```bash
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149 \
  "cd /root/child_bot && docker compose ps && systemctl status nginx --no-pager"
```

### 2. Проверить API
```bash
# Health check
curl https://vk.obyasnyatel.ru/api/health

# Avatars endpoint
curl https://vk.obyasnyatel.ru/api/avatars
```

### 3. Проверить фронтенд
```bash
# Главная страница
curl -I https://vk.obyasnyatel.ru

# Проверить загрузку JS файлов
curl -I https://vk.obyasnyatel.ru/assets/index-*.js
```

### 4. Посмотреть логи
```bash
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149

# Логи бекенда
docker compose logs -f backend --tail=100

# Логи Nginx
tail -f /var/log/nginx/vk_access.log
tail -f /var/log/nginx/vk_error.log

# Логи всех контейнеров
docker compose logs -f --tail=50
```

---

## Troubleshooting

### Backend не запускается
```bash
# Проверить логи
docker compose logs backend --tail=100

# Проверить подключение к БД
docker compose exec backend sh
# Внутри контейнера:
wget -O- http://localhost:8080/api/health

# Перезапустить с пересборкой
docker compose down backend
docker compose up -d backend
```

### Frontend показывает старую версию
```bash
# Очистить кеш браузера (Ctrl+Shift+R)

# На сервере проверить дату файлов
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149 \
  "stat /root/child_bot/frontend/dist/index.html && \
   ls -lh /root/child_bot/frontend/dist/assets/*.js | head -5"

# Принудительно перезагрузить Nginx
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149 \
  "killall nginx && sleep 2 && systemctl start nginx"
```

### Nginx не запускается
```bash
# Проверить что занимает порты
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149 \
  "lsof -i :80 && lsof -i :443"

# Проверить конфигурацию
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149 "nginx -t"

# Посмотреть логи systemd
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149 \
  "journalctl -xeu nginx.service -n 50"

# Убить все процессы Nginx и запустить заново
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149 \
  "killall nginx; sleep 2; systemctl start nginx"
```

### 403 Forbidden на API запросах
```bash
# Проверить что CSRF токен устанавливается
# В браузере: DevTools → Application → Cookies → https://vk.obyasnyatel.ru
# Должна быть cookie "csrf_token"

# Проверить логи бекенда на CSRF ошибки
docker compose logs backend | grep CSRF

# Проверить что в коде используется apiClient, а не fetch
# Все API вызовы должны идти через profileAPI, checkAPI, helpAPI и т.д.
```

### База данных недоступна
```bash
# Проверить статус PostgreSQL
docker compose ps postgres

# Проверить логи
docker compose logs postgres --tail=50

# Подключиться к БД
docker compose exec postgres psql -U child_bot -d child_bot

# Перезапустить БД (ОСТОРОЖНО!)
docker compose restart postgres
```

---

## Откат на предыдущую версию

### Backend
```bash
# Посмотреть доступные версии образов
# https://github.com/knopka87/child_bot/pkgs/container/child_bot-backend

# Изменить тег в docker-compose.yml или использовать конкретный SHA
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149 \
  "cd /root/child_bot && \
   docker compose pull backend && \
   docker compose up -d backend"
```

### Frontend
```bash
# Восстановить из бэкапа
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149 \
  "cd /root/child_bot/frontend && \
   rm -rf dist && \
   cp -r dist.backup dist && \
   nginx -s reload"
```

---

## Обновление .env файлов

### Backend .env
```bash
# Редактировать на сервере
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149 \
  "nano /root/child_bot/.env"

# После изменений - перезапустить бекенд
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149 \
  "cd /root/child_bot && docker compose up -d backend"
```

### Frontend .env.production
```bash
# НЕ ЗАБУДЬТЕ обновить локальный файл перед сборкой образа!
# frontend/.env.production должен быть закоммичен в git
# При сборке образа в GitHub Actions он копируется в образ

# Проверить текущие env переменные
cat frontend/.env.production
```

---

## Полезные команды

```bash
# Подключение к серверу
alias ssh-child-bot='ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149'

# Быстрый деплой бекенда
alias deploy-backend='ssh-child-bot "cd /root/child_bot && docker compose pull backend && docker compose up -d backend"'

# Быстрый деплой фронтенда
alias deploy-frontend='ssh-child-bot "cd /root/child_bot && docker compose pull frontend && docker create --name temp_frontend ghcr.io/knopka87/child_bot-frontend:latest && docker cp temp_frontend:/usr/share/nginx/html/. /root/child_bot/frontend/dist/ && docker rm temp_frontend && nginx -s reload"'

# Просмотр логов
alias logs-backend='ssh-child-bot "docker compose logs -f backend"'
alias logs-nginx='ssh-child-bot "tail -f /var/log/nginx/vk_error.log"'
```

Добавьте эти алиасы в `~/.zshrc` или `~/.bashrc` для удобства.
