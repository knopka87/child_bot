# Чеклист деплоя на swtest.ru

## ✅ Быстрая проверка перед деплоем

### 1. Git репозиторий готов
- [ ] Все изменения закоммичены
- [ ] Push в ветку `main` выполнен
- [ ] Файл `.env` НЕ в репозитории (проверено)
- [ ] Файл `docker-compose.production.yml` существует

### 2. Переменные окружения подготовлены
- [ ] `POSTGRES_PASSWORD` - сгенерирован (`openssl rand -base64 32`)
- [ ] `REDIS_PASSWORD` - сгенерирован (`openssl rand -base64 32`)
- [ ] `VK_APP_SECRET` - получен из VK Admin
- [ ] `EMAIL_API_KEY` - получен от Mailtrap
- [ ] `LLM_SERVER_URL` - актуальный адрес

### 3. VK Mini App настроен
- [ ] Приложение создано в https://vk.com/apps?act=manage
- [ ] Callback API URL: `https://swtest.ru/api/v1/vk/callback`
- [ ] Домен `swtest.ru` в Trusted domains
- [ ] `VK_APP_ID` и `VK_APP_SECRET` скопированы

### 4. Email провайдер настроен
- [ ] Зарегистрированы на https://mailtrap.io
- [ ] API ключ получен из Email Sending → API Keys
- [ ] Домен добавлен (опционально)

---

## 📝 Список переменных для Portainer

**Минимальный набор (обязательно):**

```
ENV=production
LOG_LEVEL=info
POSTGRES_DB=child_bot
POSTGRES_USER=child_bot
POSTGRES_PASSWORD=<сгенерированный>
REDIS_PASSWORD=<сгенерированный>
PORT=8080
DEFAULT_LLM=gpt
LLM_SERVER_URL=http://138.124.55.145:8000
ALLOWED_ORIGINS=https://swtest.ru,https://vk.com
APP_URL=https://swtest.ru
VK_APP_ID=54517931
VK_APP_SECRET=<из VK Admin>
EMAIL_PROVIDER=mailtrap
EMAIL_API_KEY=<из Mailtrap>
EMAIL_FROM=noreply@swtest.ru
VITE_API_BASE_URL=/api/v1
VITE_APP_VERSION=1.0.0
VITE_ANALYTICS_ENABLED=true
VITE_VK_APP_ID=54517931
```

**Опционально:**
```
VITE_MAX_APP_ID=<ваш MAX App ID>
VITE_TELEGRAM_BOT_USERNAME=<ваш Telegram bot>
```

---

## 🚀 Шаги деплоя

1. ✅ Войдите в https://swtest.ru
2. ✅ Создайте новый Stack с именем `child_bot`
3. ✅ Метод: Git repository
4. ✅ URL: `https://github.com/knopka87/child_bot.git`
5. ✅ Ветка: `main`
6. ✅ Compose file: `docker-compose.production.yml`
7. ✅ Включите автообновление (15 минут)
8. ✅ Введите все переменные окружения вручную
9. ✅ Deploy!

---

## 🔍 Проверка после деплоя

### Проверка контейнеров
```bash
# Все сервисы должны быть Running
✅ postgres - Running
✅ redis - Running
✅ backend - Running
✅ frontend - Running
```

### Проверка endpoints
```bash
# Healthcheck backend
curl https://swtest.ru/api/v1/health
# Ожидается: {"status":"ok"}

# Frontend
curl https://swtest.ru/
# Ожидается: HTML страница
```

### Проверка логов
1. Откройте логи backend контейнера
2. Найдите строки:
   - `[server] starting at :8080` ✅
   - `[migrate] applied` ✅ (миграции применились)
   - Нет ошибок подключения к БД ✅

---

## ⚠️ Частые проблемы

| Проблема | Решение |
|----------|---------|
| `env file not found` | ✅ Убедитесь что используете `docker-compose.production.yml`, а не `docker-compose.yml` |
| Миграции не применились | Перезапустите backend контейнер |
| CORS ошибка | Проверьте `ALLOWED_ORIGINS=https://swtest.ru,https://vk.com` |
| Email не отправляется | Проверьте `EMAIL_API_KEY` и `EMAIL_PROVIDER=mailtrap` |
| VK авторизация не работает | Проверьте `VK_APP_SECRET` и Callback API URL |

---

## 📚 Полная документация

См. `docs/deployment-swtest.ru.md` для подробной инструкции со скриншотами.
