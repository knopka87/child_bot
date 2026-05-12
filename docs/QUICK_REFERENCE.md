# Quick Reference Guide

Быстрая справка по проекту для оптимизации работы.

## SSH & Server

```bash
# SSH подключение
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149

# Быстрая команда на сервере
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149 "команда"
```

**Сервер:** `77.222.60.149`
**Домен:** `https://vk.obyasnyatel.ru`

## Database

**Connection:**
```bash
docker exec -i child_bot_postgres psql -U child_bot -d child_bot
```

**Quick Scripts:**
```bash
# Получить информацию о профиле
./scripts/get-profile.sh <profile_id>

# Найти пользователя по VK ID
./scripts/find-vk-user.sh <vk_user_id>

# Выполнить SQL запрос
./scripts/db-query.sh "SELECT * FROM child_profiles LIMIT 5"

# Логи за сегодня
./scripts/logs-today.sh [search_pattern]
```

## Docker Containers

```bash
child_bot_backend   - Backend API (порт 8080)
child_bot_postgres  - PostgreSQL
child_bot_redis     - Redis
```

**Команды:**
```bash
# Логи бекенда
docker logs child_bot_backend --tail 100 -f

# Логи за сегодня с фильтром
docker logs child_bot_backend 2>&1 | grep "2026/05/12"

# Статус контейнеров
docker ps
```

## Database Schema

### Main Table: child_profiles

**Key Fields:**
- `id` - UUID (PK)
- `display_name` - Имя пользователя
- `platform_id` - Платформа (vk, telegram, max, web)
- `platform_user_id` - ID пользователя на платформе
- `level` - Уровень
- `experience_points` - Опыт
- `coins_balance` - Баланс монет
- `streak_days` - Серия дней
- `last_activity_at` - Последняя активность

**UNIQUE constraint:** `(platform_id, platform_user_id)`

### Related Tables
- `attempts` - Попытки решения задач
- `child_achievements` - Достижения пользователя
- `referral_codes` - Реферальные коды
- `referrals` - Реферальные связи
- `villain_battles` - Битвы со злодеями
- `subscriptions` - Подписки

**Note:** Все связи `ON DELETE CASCADE` - при удалении профиля удаляются все связанные данные.

## Common SQL Queries

### Get Profile
```sql
SELECT id, display_name, platform_id, platform_user_id,
       level, experience_points, coins_balance,
       created_at, last_activity_at
FROM child_profiles
WHERE id = 'uuid-here';
```

### Find by VK User
```sql
SELECT id, display_name
FROM child_profiles
WHERE platform_id = 'vk' AND platform_user_id = '123456';
```

### Check Related Data
```sql
SELECT
  'attempts' as table_name, COUNT(*) FROM attempts WHERE child_profile_id = 'uuid'
UNION ALL
SELECT 'achievements', COUNT(*) FROM child_achievements WHERE child_profile_id = 'uuid'
UNION ALL
SELECT 'referrals', COUNT(*) FROM referrals WHERE referrer_id = 'uuid';
```

### Delete Profile (CAREFUL!)
```sql
DELETE FROM child_profiles WHERE id = 'uuid'
RETURNING id, display_name, platform_id;
```

## API Endpoints

### Authentication Headers
```
X-Platform-ID: vk|telegram|max|web
X-Child-Profile-ID: <uuid>
```

### Key Endpoints
- `GET /api/home/{childProfileId}` - Главный экран
- `GET /api/profile` - Профиль пользователя
- `POST /api/profiles/child` - Создание профиля
- `GET /api/profiles/by-platform` - Получение профиля по platform credentials
- `GET /api/achievements` - Достижения

## Платформы

**Поддерживается только VK Mini Apps** (с 12 мая 2026)

Backend валидация:
- `platform_id` должен быть `"vk"`
- Web, Telegram, MAX больше не поддерживаются
- При попытке использовать другую платформу: HTTP 400

Frontend:
- Если открыто не через VK → показывается VKOnlyAccess компонент
- VK параметры: `vk_user_id`, `vk_app_id`, `vk_platform`, `sign`

См. подробности: `docs/VK_ONLY_ACCESS.md`

## Middleware Auth

**Public paths (auth not required):**
- `/health`, `/api/health`
- `/onboarding/*`, `/api/onboarding/*`
- `/avatars`, `/api/avatars`
- `/analytics/events`, `/api/analytics/events`
- `/legal/*`, `/api/legal/*`

**Requires Platform-ID:** Все остальные эндпоинты (только `vk`)

**Requires Child-Profile-ID:** Все эндпоинты кроме:
- `/api/profiles/child` (создание)
- `/api/profiles/by-platform` (получение по credentials)
- `/api/consent`
- `/api/email/*`

## Troubleshooting

### Повторный онбординг
**Симптом:** При заходе в приложение просит пройти онбординг снова, хотя профиль существует.

**Причина:** VK Bridge не успевает инициализироваться, `getVKUserId()` возвращает `null`.

**Решение (внедрено):**
- Кэширование VK user ID в `sessionStorage`
- Fallback на кэшированный профиль при ошибках
- Не очищаем кэш при сетевых ошибках (только при явном 404)

См. подробности: `docs/ONBOARDING_FIX.md`

### Служебные значения в реферальном коде
**Симптом:** При онбординге в поле реферального кода появляется "recs" или другие непонятные значения.

**Причина:** VK передаёт служебные параметры через Launch Params (vk_ref, vk_fragment), которые не являются реферальными кодами.

**Решение (внедрено):**
- Фильтрация служебных значений: `other`, `recs`, `recommendations`
- Автоматическая очистка из storage
- Блокировка ручного ввода служебных значений

См. подробности: `docs/REFERRAL_CODE_FILTER.md`

### Приложение не работает с VPN
**Симптом:** Белый экран, spinner бесконечно крутится, ошибка "Не удалось подключиться к VK".

**Причина:** VK блокирует VPN и прокси для защиты от мошенничества. VK Bridge не может инициализироваться.

**Решение:**
1. Отключить VPN/прокси
2. Перезапустить приложение VK
3. Попробовать снова

**Техническая информация:**
- VK проверяет подпись запросов (sign parameter) на основе IP
- VPN меняет IP → подпись становится невалидной → запросы отклоняются
- Это ограничение платформы VK, обойти невозможно

См. подробности: `docs/VPN_LIMITATION.md`

### 401 Unauthorized
Проверить заголовки:
1. `X-Platform-ID` должен быть установлен
2. `X-Child-Profile-ID` должен быть установлен (если требуется для endpoint)
3. Проверить что `platform_id` в заголовке совпадает с `platform_id` профиля в БД

### Duplicate Profiles
```sql
-- Найти web-профили (потенциальные дубликаты)
SELECT id, display_name, created_at, last_activity_at
FROM child_profiles
WHERE platform_id = 'web'
ORDER BY created_at DESC;

-- Найти профили созданные близко по времени
SELECT cp1.id, cp1.platform_id, cp2.id, cp2.platform_id, cp1.created_at
FROM child_profiles cp1
JOIN child_profiles cp2 ON cp1.id < cp2.id
WHERE ABS(EXTRACT(EPOCH FROM (cp1.created_at - cp2.created_at))) < 3600
  AND cp1.platform_id != cp2.platform_id;
```

### Platform Detection Issue
Если пользователь заходит не через VK Mini App:
1. В URL нет VK параметров (`vk_user_id`, `vk_platform`, `vk_app_id`)
2. Приложение определяет платформу как `web`
3. Создается новый профиль с `platform_id='web'` и сгенерированным `platform_user_id='web-<uuid>'`

**Решение:** Пользователь должен заходить через VK Mini App: https://vk.com/app54517931

## Logs Analysis

### Find Profile Activity
```bash
# За сегодня
docker logs child_bot_backend 2>&1 | grep -a "2026/05/12" | grep -a "profile-id"

# Все логи профиля
docker logs child_bot_backend 2>&1 | grep -a "profile-id" | tail -50
```

### Error Patterns
```bash
# 401 ошибки
docker logs child_bot_backend 2>&1 | grep "401"

# Создание профилей
docker logs child_bot_backend 2>&1 | grep "CreateChildProfile"

# Auth проблемы
docker logs child_bot_backend 2>&1 | grep "\[Auth\]"
```

## Files Location

```
docs/
  ├── DB_SCHEMA.md          - Схема БД
  ├── SQL_TEMPLATES.md      - Шаблоны SQL запросов
  └── QUICK_REFERENCE.md    - Эта справка

scripts/
  ├── db-query.sh           - Выполнить SQL запрос
  ├── get-profile.sh        - Получить информацию о профиле
  ├── find-vk-user.sh       - Найти профиль по VK ID
  └── logs-today.sh         - Логи за сегодня

DEPLOYMENT.md               - Полная инструкция по деплою
```

## Environment Variables

**Production (.env.production):**
```
POSTGRES_DB=child_bot
POSTGRES_USER=child_bot
POSTGRES_PASSWORD=0xVQhgOz8E7EPVQZu0E5x7ZzixdwoX5d

VK_APP_ID=54517931
VK_APP_SECRET=BNE8tONS2h0rRjHyx9wk

APP_URL=https://77.222.60.149
```

## Support

**VK App:** https://vk.com/app54517931
**Domain:** https://vk.obyasnyatel.ru
**GitHub:** https://github.com/knopka87/child_bot
