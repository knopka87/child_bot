# Database Migrations

## Обзор

Проект использует [golang-migrate](https://github.com/golang-migrate/migrate) для управления миграциями БД.

## Существующие миграции

### Telegram Bot (001-026)
Миграции 001-026 относятся к Telegram боту и будут удалены/переработаны после полной миграции на REST API.

**Ключевые таблицы:**
- `parsed_tasks` - распарсенные задачи (используется session_id)
- `hints_cache` - кэш подсказок
- `task_sessions` - сессии Telegram (chat_id)
- `user` - пользователи Telegram (chat_id, grade)
- `metrics_events` - события метрик
- `timeline_events` - события timeline
- `textbooks`, `textbook_tasks`, `textbook_task_images` - учебники
- `textbook_task_index` - индекс задач учебников

### REST API (027-032) ✨ NEW

#### 027_child_profiles
**Профили детей** - основная таблица пользователей для REST API

**Таблица:** `child_profiles`

**Ключевые поля:**
- `id` (UUID) - уникальный ID профиля
- `display_name` - имя для отображения
- `avatar_id` - ID аватара
- `grade` - класс (1-11)
- `email` - email (опционально)
- `platform_id` - платформа: vk, telegram, max, web
- `platform_user_id` - ID пользователя на платформе
- `level`, `experience_points`, `coins_balance` - gamification
- `tasks_solved_total`, `tasks_solved_correct` - статистика
- `hints_used_total`, `streak_days` - статистика
- `last_activity_at` - последняя активность

**Уникальность:** `(platform_id, platform_user_id)` - один профиль на платформу+пользователя

---

#### 028_attempts
**Попытки решения задач** - unified таблица для help и check

**Таблица:** `attempts`

**Ключевые поля:**
- `id` (UUID) - уникальный ID попытки
- `child_profile_id` - ссылка на профиль
- `attempt_type` - 'help' или 'check'
- `status` - 'created', 'processing', 'completed', 'failed'
- `task_image_url`, `answer_image_url` - ссылки на изображения
- `detect_result`, `parse_result`, `hints_result`, `check_result` - JSON результаты LLM
- `current_hint_index` - текущая подсказка (для help)
- `hints_used` - количество использованных подсказок
- `is_correct`, `has_errors` - результат (для check)

**Индексы:**
- По профилю + дате
- По статусу
- По незавершенным попыткам

---

#### 029_achievements
**Достижения и прогресс**

**Таблицы:**
- `achievements` - справочник всех достижений
- `child_achievements` - прогресс для каждого ребенка

**achievements:**
- `id` (VARCHAR) - уникальный ID достижения
- `type` - тип: 'streak', 'tasks', 'fixes', 'villain_defeater'
- `title`, `description`, `icon` - отображение
- `requirement_type`, `requirement_value` - условия разблокировки
- `reward_type`, `reward_id`, `reward_amount` - награда
- `shelf_order`, `position_in_shelf` - UI layout

**child_achievements:**
- `child_profile_id` - ссылка на профиль
- `achievement_id` - ссылка на достижение
- `current_progress` - текущий прогресс
- `is_unlocked`, `is_claimed` - статусы
- `unlocked_at`, `claimed_at` - временные метки

**Предустановленные достижения:**
- `streak_3`, `streak_7` - серии дней
- `tasks_10`, `tasks_50` - количество задач
- `villain_1_defeated` - победа над Графом Ошибок

---

#### 030_villains
**Злодеи и битвы**

**Таблицы:**
- `villains` - справочник злодеев
- `villain_battles` - битвы пользователей
- `damage_events` - история урона

**villains:**
- `id` (VARCHAR) - уникальный ID злодея
- `name`, `description`, `image_url` - информация
- `max_hp`, `level` - характеристики
- `damage_per_correct_task` - урон за задачу
- `unlock_order` - порядок появления
- `reward_coins`, `reward_achievement_id` - награды

**villain_battles:**
- `child_profile_id` - кто сражается
- `villain_id` - с кем сражается
- `status` - 'active', 'defeated', 'abandoned'
- `current_hp` - текущее HP злодея
- `total_damage_dealt` - всего нанесено урона
- `correct_tasks_count` - правильных задач
- `rewards_claimed` - забраны ли награды

**damage_events:**
- `battle_id` - ссылка на битву
- `attempt_id` - ссылка на попытку
- `damage` - нанесенный урон
- `task_type` - 'help' или 'check'

**Предустановленные злодеи:**
- `count_error` - Граф Ошибок (HP: 100)
- `baron_confusion` - Барон Путаница (HP: 150)
- `duchess_distraction` - Герцогиня Отвлечения (HP: 200)

---

#### 031_subscriptions
**Подписки и планы**

**Таблицы:**
- `subscription_plans` - планы подписок
- `subscriptions` - подписки пользователей

**subscription_plans:**
- `id` (VARCHAR) - ID плана
- `name`, `description` - информация
- `price_cents` - цена в копейках
- `currency` - валюта (RUB)
- `duration_days` - длительность
- `trial_days` - пробный период
- `discount_percent` - скидка
- `is_popular` - популярный план

**subscriptions:**
- `child_profile_id` - ссылка на профиль
- `plan_id` - ссылка на план
- `status` - 'trial', 'active', 'expired', 'cancelled'
- `started_at`, `trial_ends_at`, `expires_at`, `cancelled_at` - даты
- `auto_renew` - автопродление
- `payment_provider`, `payment_external_id` - платежная информация

**Constraint:** Только одна активная подписка на профиль

**Предустановленные планы:**
- `monthly` - Месячная (499 руб, 7 дней trial)
- `yearly` - Годовая (3999 руб, 14 дней trial, скидка 33%)

---

#### 032_referrals
**Реферальная программа**

**Таблицы:**
- `referrals` - реферальные связи
- `referral_codes` - уникальные коды
- `referral_milestones` - награды за количество
- `child_referral_milestones` - прогресс milestone

**referrals:**
- `referrer_id` - кто пригласил
- `referred_id` - кого пригласили
- `is_active` - активировался ли приглашенный
- `reward_coins` - награда (50 по умолчанию)
- `reward_claimed` - забрана ли награда

**referral_codes:**
- `child_profile_id` - владелец кода
- `code` - уникальный код (8 символов)
- `uses_count` - количество использований
- `max_uses` - максимум (NULL = unlimited)

**referral_milestones:**
- `friends_count` - количество друзей
- `reward_coins` - награда
- `description` - описание

**Предустановленные milestone:**
- 1 друг - 50 монет
- 3 друга - 100 монет
- 5 друзей - 200 монет
- 10 друзей - 500 монет

**Автоматика:**
- При создании `child_profile` автоматически генерируется уникальный реферальный код

---

## Команды

### Применить миграции
```bash
make migrate-up
# или
migrate -source "file://api/migrations" -database "$DATABASE_URL" up
```

### Откатить последнюю миграцию
```bash
make migrate-down
# или
migrate -source "file://api/migrations" -database "$DATABASE_URL" down 1
```

### Создать новую миграцию
```bash
make migrate-create NAME=description
# или
migrate create -ext sql -dir api/migrations -seq description
```

### Проверить статус
```bash
make migrate-status
# или
migrate -source "file://api/migrations" -database "$DATABASE_URL" version
```

---

## ERD (Entity Relationship Diagram)

```
child_profiles (UUID)
    ├─→ attempts (1:N)
    ├─→ child_achievements (1:N)
    ├─→ villain_battles (1:N)
    ├─→ subscriptions (1:1 active)
    ├─→ referrals (1:N as referrer)
    ├─→ referrals (1:N as referred)
    ├─→ referral_codes (1:1)
    └─→ child_referral_milestones (1:N)

achievements
    └─→ child_achievements (1:N)

villains
    └─→ villain_battles (1:N)

villain_battles
    └─→ damage_events (1:N)

subscription_plans
    └─→ subscriptions (1:N)

referral_milestones
    └─→ child_referral_milestones (1:N)
```

---

## Примечания

### Совместимость с Telegram Bot
- Существующие таблицы `parsed_tasks`, `hints_cache`, `task_sessions` могут использоваться REST API через `session_id`
- После полной миграции на REST API, Telegram-специфичные таблицы будут удалены (Phase 7)

### UUID vs BIGSERIAL
- **UUID** используется для публичных ID (child_profiles, attempts) - защита от перебора
- **BIGSERIAL** используется для внутренних связей (battles, achievements progress)

### JSONB колонки
- `detect_result`, `parse_result`, `hints_result`, `check_result` в `attempts` - результаты LLM API
- Позволяют хранить полные responses без нормализации
- Можно индексировать и делать запросы по JSON (при необходимости)

### Триггеры
- `updated_at` автоматически обновляется при UPDATE
- `referral_code` автоматически генерируется при создании профиля

---

## TODO

- [ ] Phase 6: Написать integration tests для миграций
- [ ] Phase 6: Протестировать rollback (down migrations)
- [ ] Phase 7: Удалить Telegram-специфичные таблицы после миграции
