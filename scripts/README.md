# Скрипты для разработки

## Локальные туннели

Для разработки VK Mini Apps необходим публичный HTTPS URL. Используйте один из вариантов:

### Вариант 1: localtunnel (рекомендуется для dev)

**Преимущества:**
- Полностью бесплатный
- Поддержка фиксированных subdomain (URL не меняется)
- Стабильнее cloudflare tunnels для Vite dev server
- Не требует регистрации

**Недостатки:**
- Медленнее ngrok
- Иногда может быть недоступен

**Установка:**
```bash
npm install -g localtunnel
```

**Использование:**
```bash
# Запуск обоих туннелей (frontend + backend) с фиксированными URL
./scripts/tunnel-lt.sh

# Или вручную:
lt --port 5173 --subdomain childbotfe  # Frontend
lt --port 8080 --subdomain childbotbe  # Backend
```

**URL (фиксированные):**
- Frontend: https://childbotfe.loca.lt
- Backend: https://childbotbe.loca.lt

**Логи:**
- Frontend: `tail -f /tmp/lt-frontend.log`
- Backend: `tail -f /tmp/lt-backend.log`

---

### Вариант 2: ngrok (самый стабильный)

**Преимущества:**
- Стабильный и быстрый
- Красивый веб-интерфейс с логами
- Поддержка custom domains (платно)

**Установка:**
```bash
brew install ngrok  # macOS
# или скачайте с https://ngrok.com/download
```

**Настройка (один раз):**
1. Зарегистрируйтесь: https://dashboard.ngrok.com/signup
2. Получите токен: https://dashboard.ngrok.com/get-started/your-authtoken
3. Выполните: `ngrok config add-authtoken <ваш-токен>`

**Использование:**
```bash
# Frontend (Vite dev server, порт 5173)
./scripts/tunnel.sh frontend

# Backend API (порт 8080)
./scripts/tunnel.sh backend

# Production Frontend (порт 80)
./scripts/tunnel.sh prod
```

**Веб-интерфейс:** http://127.0.0.1:4040

---

### Вариант 3: Cloudflare Tunnel (не рекомендуется)

**⚠️ Внимание:** Бесплатные cloudflare tunnels (без аккаунта) очень нестабильны для Vite dev server. Часто возникают ошибки 524 (timeout) при загрузке JavaScript модулей.

**Преимущества:**
- Полностью бесплатный
- Не требует регистрации
- От Cloudflare (надежно для production)

**Недостатки:**
- ❌ Нестабилен для Vite HMR (горячая перезагрузка)
- ❌ Короткие таймауты → ошибки 524
- Медленнее ngrok
- Случайные URL при каждом запуске
- Нет веб-интерфейса

**Установка:**
```bash
brew install cloudflared  # macOS
# или https://developers.cloudflare.com/cloudflare-one/connections/connect-apps/install-and-setup/installation
```

**Использование:**
```bash
# Frontend
./scripts/tunnel-cloudflare.sh frontend

# Backend
./scripts/tunnel-cloudflare.sh backend

# Production
./scripts/tunnel-cloudflare.sh prod
```

---

## Настройка после запуска туннеля

### Автоматическая настройка (localtunnel)

Если используете `./scripts/tunnel-lt.sh`, настройка уже выполнена автоматически:
- ✅ ALLOWED_ORIGINS обновлен в .env
- ✅ VITE_API_BASE_URL обновлен в docker-compose.dev.yml
- ✅ URL фиксированные (не меняются при перезапуске)

**Осталось только:**
1. Открыть https://dev.vk.com/mini-apps/management/settings?id=54517931
2. Установить **Адрес приложения** для всех трех пунктов: `https://childbotfe.loca.lt`
3. Сохранить и обновить страницу (Ctrl+F5)

### Ручная настройка (ngrok, cloudflare)

### 1. Получите HTTPS URL
После запуска туннеля скопируйте URL (например, `https://abc123.ngrok.io`)

### 2. Настройте VK Mini App
1. Откройте https://dev.vk.com/mini-apps/management/settings?id=54517931
2. Установите **Адрес приложения** для всех трех пунктов
3. Сохраните изменения

### 3. Обновите ALLOWED_ORIGINS
Добавьте frontend URL в `.env`:
```env
ALLOWED_ORIGINS=https://abc123.ngrok.io,http://localhost:5173
```

### 4. Обновите VITE_API_BASE_URL
Если туннелируете backend отдельно, обновите в `docker/docker-compose.dev.yml`:
```yaml
environment:
  VITE_API_BASE_URL: https://your-backend-url.ngrok.io
```

### 5. Перезапустите контейнеры
```bash
docker compose -f docker/docker-compose.dev.yml restart backend frontend
```

---

## Проверка

### Frontend туннель
```bash
curl https://your-url.ngrok.io
# Должен вернуть HTML вашего приложения
```

### Backend туннель
```bash
curl https://your-url.ngrok.io/health
# Должен вернуть {"status":"ok"}
```

---

## Troubleshooting

### localtunnel: "503 Service Unavailable"
Туннель упал или не запустился. Решение:
```bash
# Перезапустите туннели
./scripts/tunnel-lt.sh
```

### localtunnel: "Subdomain is already in use"
Кто-то занял ваш subdomain. Варианты:
1. Подождите 5-10 минут и попробуйте снова
2. Используйте другой subdomain в скрипте

### cloudflare: Ошибка 524 (Timeout)
Cloudflare tunnels без аккаунта нестабильны для Vite. Используйте localtunnel или ngrok.

### ngrok: "ERR_NGROK_108"
Не добавлен authtoken. Выполните:
```bash
ngrok config add-authtoken <ваш-токен>
```

### ngrok: "ERR_NGROK_334" (endpoint already online)
У вас уже запущен туннель с этим URL. Решение:
```bash
pkill ngrok  # Остановите старый туннель
./scripts/tunnel.sh frontend  # Запустите новый
```

### CORS ошибки
Убедитесь что туннельный URL добавлен в `ALLOWED_ORIGINS` (.env) и контейнеры перезапущены:
```bash
docker compose -f docker/docker-compose.dev.yml restart backend
```

### VK: "Приложение не инициализировано"
Возможные причины:
1. Туннель не работает → проверьте `curl https://your-url.loca.lt`
2. URL в VK неправильный → проверьте настройки VK Mini App
3. Кэш VK → обновите страницу через Ctrl+F5 или откройте в режиме инкогнито

### Медленная работа
- localtunnel: это нормально, он медленнее ngrok
- ngrok: проверьте регион в конфиге (~/.ngrok2/ngrok.yml)
- cloudflared: используйте localtunnel или ngrok

### Туннель закрывается при закрытии терминала
Туннели запускаются через `nohup` и работают в фоне. Если они все равно падают:
```bash
# Используйте tmux
tmux new -s tunnel
./scripts/tunnel-lt.sh
# Ctrl+B, затем D для отсоединения
# tmux attach -t tunnel для возврата
```

---

## Сравнение туннелей

| Критерий | localtunnel | ngrok | cloudflared |
|----------|-------------|-------|-------------|
| Цена | Бесплатно | Бесплатно (1 туннель) / $10/мес | Бесплатно |
| Регистрация | Нет | Да | Нет |
| Стабильность для Vite | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐ |
| Фиксированные URL | Да (subdomain) | Да (платно) | Нет |
| Веб-интерфейс | Нет | Да | Нет |
| Скорость | Средняя | Высокая | Низкая |

**Рекомендация:**
- Для разработки: **localtunnel** (бесплатно + фиксированные URL)
- Для production: **ngrok** с платным планом (максимальная стабильность)

## Другие альтернативы

Если вышеперечисленное не подходит:
- **serveo**: `ssh -R 80:localhost:5173 serveo.net` (часто недоступен)
- **telebit**: https://telebit.cloud/ (требует регистрацию)
- **bore**: https://github.com/ekzhang/bore (self-hosted)