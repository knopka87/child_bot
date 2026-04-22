# Инструкция по деплою приложения на swtest.ru

---

## 📋 Общая информация

swtest.ru использует веб-интерфейс управления контейнерами на базе **Portainer / Portainer Business**. Деплой осуществляется через веб-интерфейс без доступа к SSH консоли сервера.

✅ Приложение полностью адаптировано для деплоя через веб-интерфейс swtest.ru

⚠️ **КРИТИЧНО ПО БЕЗОПАСНОСТИ:**
- Проект содержит production конфигурацию в `docker-compose.production.yml`
- Файл `.env` НЕ должен коммититься в Git (добавлен в `.gitignore`)
- Все секретные ключи должны быть уникальными для production
- **НИКОГДА не коммитьте реальные `VK_APP_SECRET`, `EMAIL_API_KEY`, `POSTGRES_PASSWORD` в Git**
- Используйте `.env.example` как шаблон с placeholder значениями

---

## ✅ Чеклист готовности к деплою

Перед началом деплоя убедитесь что:

**Конфигурация:**
- [x] Файл `docker-compose.production.yml` создан в корне проекта
- [x] Dockerfile для backend готов (multi-stage build с Alpine)
- [x] Dockerfile для frontend готов (build + nginx)
- [x] nginx.conf настроен с security headers и API proxy
- [x] Миграции базы данных в `api/migrations/`

**Безопасность:**
- [x] `.env` добавлен в `.gitignore`
- [ ] Создан файл `.env.production` с уникальными секретами
- [ ] `POSTGRES_PASSWORD` - случайный пароль минимум 16 символов
- [ ] `REDIS_PASSWORD` - случайный пароль минимум 16 символов
- [ ] `VK_APP_SECRET` получен из VK Admin панели
- [ ] `EMAIL_API_KEY` получен от email провайдера

**VK Mini App:**
- [ ] Приложение создано в VK Admin
- [ ] Callback API настроен на `https://swtest.ru/api/v1/vk/callback`
- [ ] VK Pay настроен (если используется платная подписка)
- [ ] Домен `swtest.ru` добавлен в Trusted domains

**Email:**
- [ ] Email провайдер настроен: ✅ **Mailtrap** (рекомендуется)
- [ ] Зарегистрируйтесь на https://mailtrap.io
- [ ] Получите API ключ в разделе **Email Sending → API Keys**
- [ ] `EMAIL_PROVIDER=mailtrap` (или не указывать - mailtrap по умолчанию)
- [ ] `EMAIL_API_KEY` указан в переменных окружения
- [ ] `EMAIL_FROM=noreply@swtest.ru`
- [ ] ✅ Бесплатный тариф Mailtrap: 1000 писем/месяц
- [ ] Альтернативы: sendgrid, ses, mailgun - все поддерживаются

---

## 🔑 Предварительная подготовка

✅ **swtest.ru автоматически собирает образы напрямую из GitHub!** Локальная сборка не требуется.

### 1. Подготовка репозитория
```bash
# Коммитьте все исправления
git add .
git commit -m "production deploy preparation"
git push origin main
```

### 2. Подготовка переменных окружения
Скопируйте и заполните файл `.env.production`:
```bash
cp .env.production.example .env.production
```

⚠️ **Обязательно измените:**
- `POSTGRES_PASSWORD` - случайный пароль 32+ символов
- `REDIS_PASSWORD` - случайный пароль 32+ символов
- `APP_URL=https://swtest.ru`
- `ALLOWED_ORIGINS=https://swtest.ru,https://vk.com`
- `LLM_SERVER_URL` - актуальный адрес LLM сервера
- `VK_APP_ID` - ID вашего VK Mini App
- `VK_APP_SECRET` - Secret Key из настроек VK Mini App
- `VITE_VK_APP_ID` - ID приложения для фронтенда
- `EMAIL_API_KEY` - API ключ от Mailtrap (https://mailtrap.io)
- `EMAIL_FROM` - email адрес отправителя (noreply@swtest.ru)

📝 **Как сгенерировать безопасные пароли:**
```bash
# PostgreSQL пароль
openssl rand -base64 32

# Redis пароль
openssl rand -base64 32
```

### 3. Проверка переменных окружения
Перед деплоем проверьте что все критичные переменные заполнены:

```bash
# Проверка что все секреты установлены
grep -E "CHANGE_ME|your_|YOUR_" .env.production

# Если команда вернула совпадения - значит есть незаполненные переменные
# Заполните их перед деплоем!
```

✅ **Контрольный список переменных:**
- [ ] `POSTGRES_PASSWORD` - установлен уникальный пароль
- [ ] `REDIS_PASSWORD` - установлен уникальный пароль
- [ ] `VK_APP_SECRET` - указан Secret Key из VK Admin
- [ ] `EMAIL_API_KEY` - указан ключ от email провайдера
- [ ] `LLM_SERVER_URL` - указан актуальный адрес
- [ ] `APP_URL` - установлен на `https://swtest.ru`
- [ ] `ALLOWED_ORIGINS` - содержит `https://swtest.ru,https://vk.com`

### 4. Как получить API ключ от Mailtrap (рекомендуется)

📧 **Mailtrap** - удобный email провайдер с бесплатным тарифом 1000 писем/месяц.

**Шаги получения API ключа:**

1. Зарегистрируйтесь на https://mailtrap.io
2. Подтвердите email адрес
3. Войдите в панель управления
4. Перейдите в раздел **Email Sending** (не Testing!)
5. Нажмите **Domains** → **Add Domain**
6. Добавьте домен `swtest.ru` (потребуется подтверждение через DNS записи)
   - Или используйте тестовый домен Mailtrap (для начала)
7. Перейдите в **Settings** → **API Tokens**
8. Нажмите **Create Token**
9. Скопируйте токен - это ваш `EMAIL_API_KEY`

**Настройка в .env.production:**
```bash
EMAIL_PROVIDER=mailtrap  # Использует нативный Mailtrap API
EMAIL_API_KEY=your_mailtrap_token_here
EMAIL_FROM=noreply@swtest.ru
```

💡 **Почему Mailtrap:**
- ✅ Нативная поддержка в коде (отдельная реализация)
- ✅ 1000 писем/месяц бесплатно
- ✅ Простая настройка без верификации домена
- ✅ Удобная панель управления

💡 **Альтернативные провайдеры:**
| Провайдер | Бесплатный лимит | EMAIL_PROVIDER |
|-----------|------------------|----------------|
| **Mailtrap** (рекомендуется) | 1000 писем/месяц | `mailtrap` (default) |
| SendGrid | 100 писем/день | `sendgrid` |
| AWS SES | 62000 писем/месяц | `ses` |
| Mailgun | 5000 писем/месяц | `mailgun` |

---

## 🚀 Пошаговый деплой через веб-интерфейс swtest.ru

---

### 📌 Шаг 1. Авторизация
1. Откройте в браузере: `https://swtest.ru`
2. Введите логин и пароль
3. Перейдите в раздел **Stacks / Стек**

---

### 📌 Шаг 2. Создание нового стека
1. Нажмите кнопку **Add stack / Добавить стек**
2. Введите имя стека: `child_bot`
3. Выберите метод: **Git repository**
4. Введите адрес репозитория: `https://github.com/knopka87/child_bot.git`
5. Укажите ветку: `main`
6. Укажите путь к compose файлу: `docker-compose.production.yml`
7. Включите опцию **Automatic updates / Автоматическое обновление**
8. Укажите интервал обновления: `15 минут`

---

### 📌 Шаг 3. Подготовка файла `.env.production` (на локальной машине)

⚠️ **Этот шаг нужно выполнить на вашем компьютере ПЕРЕД загрузкой в swtest.ru**

```bash
# 1. Перейдите в директорию проекта
cd /path/to/child_bot

# 2. Скопируйте шаблон
cp .env.production.example .env.production

# 3. Сгенерируйте безопасные пароли
echo "POSTGRES_PASSWORD=$(openssl rand -base64 32)"
echo "REDIS_PASSWORD=$(openssl rand -base64 32)"

# 4. Откройте файл в редакторе
nano .env.production  # или code .env.production
```

**Заполните следующие переменные:**

```bash
# Environment
ENV=production
LOG_LEVEL=info

# Database
POSTGRES_PASSWORD=сгенерированный_пароль_postgresql
REDIS_PASSWORD=сгенерированный_пароль_redis

# Backend
PORT=8080
DEFAULT_LLM=gpt
LLM_SERVER_URL=http://138.124.55.145:8000  # или ваш адрес
ALLOWED_ORIGINS=https://swtest.ru,https://vk.com
APP_URL=https://swtest.ru

# VK Mini App (из VK Admin)
VK_APP_ID=54517931
VK_APP_SECRET=ваш_vk_app_secret_из_vk_admin

# Email (Mailtrap рекомендуется - 1000 писем/месяц бесплатно)
EMAIL_PROVIDER=mailtrap  # Или: sendgrid, ses, mailgun
EMAIL_API_KEY=ваш_mailtrap_api_key_из_mailtrap_io
EMAIL_FROM=noreply@swtest.ru

# Frontend (Vite build args)
VITE_API_BASE_URL=/api/v1
VITE_APP_VERSION=1.0.0
VITE_ANALYTICS_ENABLED=true
VITE_VK_APP_ID=54517931
VITE_MAX_APP_ID=ваш_max_app_id
VITE_TELEGRAM_BOT_USERNAME=ваш_bot_username
```

**Проверка перед использованием:**
```bash
# Убедитесь что нет незаполненных значений
grep -E "CHANGE_ME|your_vk|your_max|your_bot|your-llm|ваш_|сгенерированный_" .env.production

# Результат должен быть пустым!
```

💾 **Сохраните файл `.env.production` - он понадобится в следующем шаге**

---

### 📌 Шаг 4. Настройка переменных окружения (в swtest.ru)

1. Перейдите на вкладку **Environment variables / Переменные окружения**
2. Нажмите **Load variables from .env file**
3. Загрузите подготовленный файл `.env.production` (созданный в Шаге 3)
4. Проверьте что все переменные загрузились корректно

⚠️ **Если загрузка файла не работает:**
- Введите переменные вручную по одной через кнопку **Add variable**
- Скопируйте имя и значение из вашего `.env.production` файла

---

### 📌 Шаг 5. Настройка сети
1. Перейдите на вкладку **Network / Сеть**
2. Выберите сеть: `public-network`
3. Включите опцию **Enable public access / Включить публичный доступ**
4. Укажите домен: `swtest.ru`
5. Включите **Automatic SSL / Автоматический SSL**

---

### 📌 Шаг 6. Запуск деплоя
1. Проверьте все настройки
2. Нажмите кнопку **Deploy the stack / Развернуть стек**
3. Дождитесь завершения деплоя (занимает ~2 минуты)

---

### 📌 Шаг 7. Проверка статуса сервисов
После завершения деплоя перейдите на страницу стека `child_bot`:

✅ Все контейнеры должны быть в статусе **Running / Запущен**:
| Сервис | Статус |
|--------|--------|
| `postgres` | ✅ Running |
| `redis` | ✅ Running |
| `backend` | ✅ Running |
| `frontend` | ✅ Running |

✅ Проверка healthcheck:
1. Откройте логи каждого контейнера
2. Убедитесь что нет ошибок
3. Проверьте что миграции базы данных применились успешно

---

## ✅ Пост-деплойные проверки

### 1. Проверка доступности сайта
Откройте в браузере: `https://swtest.ru`

✅ Ожидаемый результат:
- Сайт открывается без ошибок
- Нет ошибок в консоли разработчика
- Все статические файлы загружаются

### 2. Проверка работоспособности API
```bash
# Проверка healthcheck бэкенда
curl https://swtest.ru/api/health
```

✅ Ожидаемый ответ: `{"status": "ok"}`

### 3. Проверка функциональности
1. ✅ Авторизация через VK
2. ✅ Создание профиля ребенка
3. ✅ Отправка задания на проверку
4. ✅ Работа чата
5. ✅ Страница отчетов
6. ✅ Страница достижений

---

## ⚠️ Известные проблемы веб-интерфейса swtest.ru

| Проблема | Решение |
|----------|---------|
| ❌ Ошибка при загрузке .env файла | Вводите переменные вручную по одной |
| ❌ Контейнер не запускается | Увеличьте лимит памяти на сервисе |
| ❌ Не применяются миграции | Перезапустите контейнер backend вручную |
| ❌ Ошибка CORS | Добавьте домен в `ALLOWED_ORIGINS` и перезапустите стек |
| ❌ Вебсокет не работает | Включите опцию `WebSocket support` в настройках сети |

---

## 🔄 Автоматическое обновление приложения

✅ swtest.ru умеет автоматически обновлять приложение при коммитах в GitHub:

1. Приложение автоматически проверяет обновления в репозитории каждые 15 минут
2. При обнаружении новых коммитов автоматически пересобираются образы
3. Происходит бесшовное обновление контейнеров без простоя сервиса
4. Healthcheck автоматически откатывает обновление при ошибках

✅ Для ручного обновления:
1. Перейдите на страницу стека `child_bot`
2. Нажмите кнопку **Force update / Принудительное обновление**
3. Подтвердите обновление

---

## 📊 Мониторинг и логи

Через веб-интерфейс доступно:
1. ✅ Просмотр логов каждого контейнера в реальном времени
2. ✅ Графики использования CPU, RAM, диска
3. ✅ Статус healthcheck сервисов
4. ✅ История событий и перезапусков

---

## 📌 Полезные ссылки

- Документация swtest.ru: https://swtest.ru/docs
- Панель управления: https://swtest.ru/panel
- Поддержка: support@swtest.ru

---

### 📌 Дополнительная проверка
✅ В корне репозитория находится `docker-compose.production.yml` для production деплоя
✅ Backend использует Alpine Linux образ с минимальным размером
✅ Frontend собирается в production режиме с nginx
✅ Все сервисы имеют healthcheck для автоматического мониторинга
✅ Миграции базы данных применяются автоматически при старте backend

---

## 🔐 Проверка безопасности перед деплоем

Убедитесь что выполнены следующие требования безопасности:

| Проверка | Статус |
|----------|--------|
| Файл `.env` добавлен в `.gitignore` | ✅ |
| `VK_APP_SECRET` установлен из VK Admin | ⚠️ Проверьте |
| `POSTGRES_PASSWORD` - уникальный случайный пароль | ⚠️ Проверьте |
| `REDIS_PASSWORD` - уникальный случайный пароль | ⚠️ Проверьте |
| `EMAIL_API_KEY` установлен | ⚠️ Проверьте |
| HTTPS настроен через swtest.ru | ✅ |
| CORS настроен только для `swtest.ru` и `vk.com` | ✅ |
| Security headers включены в nginx | ✅ |

---

✅ Деплой завершён. Приложение готово к работе на продакшене.
