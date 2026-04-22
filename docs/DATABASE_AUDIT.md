# Детальный аудит базы данных - Homework Helper Bot

**Дата:** 2026-04-05
**Всего таблиц:** 28
**Версия:** 2.0 (полное описание всех полей)

---

## 🟢 АКТИВНЫЕ ТАБЛИЦЫ - REST API (16 таблиц)

### 1. child_profiles (96 kB)
**Миграция:** `027_child_profiles.up.sql`
**Назначение:** Профили детей-пользователей для REST API
**Комментарий:** Профили детей для REST API (независимо от Telegram)

#### Поля:

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | UUID PRIMARY KEY | Уникальный идентификатор профиля |
| `display_name` | VARCHAR(100) NOT NULL | Отображаемое имя ребёнка |
| `avatar_id` | VARCHAR(50) | ID выбранного аватара |
| `grade` | INTEGER NOT NULL | Класс обучения (1-11), CHECK: >= 1 AND <= 11 |
| `email` | VARCHAR(255) | Email родителя (опционально) |
| `email_verified` | BOOLEAN DEFAULT FALSE | Подтверждён ли email |
| `platform_id` | VARCHAR(20) NOT NULL | Платформа регистрации: 'vk', 'telegram', 'max', 'web' |
| `platform_user_id` | VARCHAR(255) | ID пользователя на платформе (опционально) |
| `level` | INTEGER DEFAULT 1 | Уровень пользователя в gamification, CHECK: >= 1 |
| `experience_points` | INTEGER DEFAULT 0 | Очки опыта для прогресса уровня, CHECK: >= 0 |
| `coins_balance` | INTEGER DEFAULT 0 | Баланс монет (геймификация), CHECK: >= 0 |
| `tasks_solved_total` | INTEGER DEFAULT 0 | Всего решённых задач, CHECK: >= 0 |
| `tasks_solved_correct` | INTEGER DEFAULT 0 | Правильно решённых задач, CHECK: >= 0 |
| `hints_used_total` | INTEGER DEFAULT 0 | Всего использовано подсказок, CHECK: >= 0 |
| `streak_days` | INTEGER DEFAULT 0 | Серия дней подряд, CHECK: >= 0 |
| `last_activity_at` | TIMESTAMPTZ | Дата последней активности |
| `created_at` | TIMESTAMPTZ NOT NULL DEFAULT NOW() | Дата создания профиля |
| `updated_at` | TIMESTAMPTZ NOT NULL DEFAULT NOW() | Дата последнего обновления |

**Индексы:**
- `idx_child_profiles_platform` - (platform_id, platform_user_id)
- `idx_child_profiles_activity` - (last_activity_at DESC)
- `idx_child_profiles_created` - (created_at DESC)

**Constraints:**
- UNIQUE (platform_id, platform_user_id) - один профиль на платформу+пользователя

**Триггер:** `child_profiles_updated_at` - автообновление updated_at

---

### 2. attempts (448 kB)
**Миграция:** `028_attempts.up.sql`
**Назначение:** Попытки решения заданий (help и check)
**Комментарий:** Унифицированная таблица для help (подсказки) и check (проверка решения)

#### Поля:

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | UUID PRIMARY KEY | Уникальный идентификатор попытки |
| `child_profile_id` | UUID NOT NULL | ID профиля ребёнка, FK → child_profiles(id) ON DELETE CASCADE |
| `attempt_type` | VARCHAR(20) NOT NULL | Тип попытки: 'help' или 'check', CHECK: IN ('help', 'check') |
| `status` | VARCHAR(20) NOT NULL DEFAULT 'created' | Статус: created, processing, completed, failed |
| `task_image_url` | TEXT | URL или путь к изображению задания |
| `answer_image_url` | TEXT | URL или путь к изображению ответа (для check) |
| `detect_result` | JSONB | JSON результат Detect API (определение типа задания) |
| `parse_result` | JSONB | JSON результат Parse API (извлечение условия задачи) |
| `hints_result` | JSONB | JSON результат Hint API (подсказки для help) |
| `check_result` | JSONB | JSON результат CheckSolution API (проверка для check) |
| `current_hint_index` | INTEGER DEFAULT 0 | Текущий индекс подсказки (для help) |
| `hints_used` | INTEGER DEFAULT 0 | Количество использованных подсказок |
| `time_spent_seconds` | INTEGER | Время решения в секундах |
| `is_correct` | BOOLEAN | Правильность решения (для check) |
| `has_errors` | BOOLEAN | Наличие ошибок (для check) |
| `created_at` | TIMESTAMPTZ NOT NULL DEFAULT NOW() | Дата создания попытки |
| `updated_at` | TIMESTAMPTZ NOT NULL DEFAULT NOW() | Дата последнего обновления |
| `completed_at` | TIMESTAMPTZ | Дата завершения попытки |

**Индексы:**
- `idx_attempts_child_profile` - (child_profile_id, created_at DESC)
- `idx_attempts_status` - (status, created_at DESC)
- `idx_attempts_type` - (attempt_type, created_at DESC)
- `idx_attempts_unfinished` - (child_profile_id, status) WHERE status IN ('created', 'processing')

**Триггер:** `attempts_updated_at` - автообновление updated_at

---

### 3. achievements (48 kB)
**Миграция:** `029_achievements.up.sql`
**Назначение:** Справочник всех достижений

#### Поля:

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | VARCHAR(100) PRIMARY KEY | ID достижения, например: 'streak_7', 'tasks_10', 'villain_1_defeated' |
| `type` | VARCHAR(50) NOT NULL | Тип: 'streak', 'tasks', 'fixes', 'villain_defeater' |
| `title` | VARCHAR(200) NOT NULL | Название достижения |
| `description` | TEXT NOT NULL | Описание достижения |
| `icon` | VARCHAR(200) NOT NULL | Emoji или URL изображения иконки |
| `requirement_type` | VARCHAR(50) NOT NULL | Тип условия: 'streak_days', 'tasks_count', 'villain_defeated' |
| `requirement_value` | INTEGER NOT NULL | Целевое значение для разблокировки |
| `reward_type` | VARCHAR(50) NOT NULL | Тип награды: 'coins', 'sticker', 'avatar', 'badge' |
| `reward_id` | VARCHAR(100) | ID награды |
| `reward_name` | VARCHAR(200) | Название награды |
| `reward_amount` | INTEGER | Количество (для coins) |
| `shelf_order` | INTEGER DEFAULT 1 | Номер полки для отображения (1-3) |
| `position_in_shelf` | INTEGER DEFAULT 0 | Позиция на полке (0-3) |
| `created_at` | TIMESTAMPTZ NOT NULL DEFAULT NOW() | Дата создания |

**Индексы:**
- `idx_achievements_type` - (type)

**Примеры данных:**
- `streak_3` - Три дня подряд (🔥, 30 монет)
- `streak_7` - Неделя успеха (🔥, 50 монет)
- `tasks_10` - Начинающий (🎯, 20 монет)
- `villain_1_defeated` - Победитель Графа (⚔️, стикер)

---

### 4. child_achievements (32 kB)
**Миграция:** `029_achievements.up.sql`
**Назначение:** Прогресс достижений для каждого ребенка

#### Поля:

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | BIGSERIAL PRIMARY KEY | Уникальный идентификатор записи |
| `child_profile_id` | UUID NOT NULL | ID профиля ребёнка, FK → child_profiles(id) ON DELETE CASCADE |
| `achievement_id` | VARCHAR(100) NOT NULL | ID достижения, FK → achievements(id) ON DELETE CASCADE |
| `current_progress` | INTEGER DEFAULT 0 | Текущий прогресс, CHECK: >= 0 |
| `is_unlocked` | BOOLEAN DEFAULT FALSE | Разблокировано ли достижение |
| `is_claimed` | BOOLEAN DEFAULT FALSE | Забрана ли награда пользователем |
| `unlocked_at` | TIMESTAMPTZ | Дата разблокировки |
| `claimed_at` | TIMESTAMPTZ | Дата получения награды |
| `created_at` | TIMESTAMPTZ NOT NULL DEFAULT NOW() | Дата создания записи |
| `updated_at` | TIMESTAMPTZ NOT NULL DEFAULT NOW() | Дата обновления |

**Индексы:**
- `idx_child_achievements_profile` - (child_profile_id, is_unlocked)
- `idx_child_achievements_unlocked` - (child_profile_id, unlocked_at DESC) WHERE is_unlocked = TRUE

**Constraints:**
- UNIQUE (child_profile_id, achievement_id) - одно достижение на профиль

**Триггер:** `child_achievements_updated_at` - автообновление updated_at

---

### 5. villains (48 kB)
**Миграция:** `030_villains.up.sql`
**Назначение:** Справочник злодеев (геймификация)

#### Поля:

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | VARCHAR(100) PRIMARY KEY | ID злодея, например: 'count_error', 'baron_confusion' |
| `name` | VARCHAR(200) NOT NULL | Имя злодея |
| `description` | TEXT NOT NULL | Описание злодея |
| `image_url` | TEXT NOT NULL | URL изображения злодея |
| `max_hp` | INTEGER NOT NULL | Максимальное здоровье, CHECK: > 0 |
| `level` | INTEGER NOT NULL DEFAULT 1 | Уровень сложности, CHECK: >= 1 |
| `damage_per_correct_task` | INTEGER NOT NULL DEFAULT 5 | Урон за одну правильно решенную задачу |
| `unlock_order` | INTEGER NOT NULL DEFAULT 1 | Порядок появления злодеев |
| `reward_coins` | INTEGER DEFAULT 100 | Награда монетами за победу |
| `reward_achievement_id` | VARCHAR(100) | ID достижения за победу |
| `created_at` | TIMESTAMPTZ NOT NULL DEFAULT NOW() | Дата создания |

**Индексы:**
- `idx_villains_unlock_order` - (unlock_order)

**Примеры данных:**
- `count_error` - Граф Ошибок (HP: 100, урон: 5, награда: 100 монет)
- `baron_confusion` - Барон Путаница (HP: 150, урон: 5, награда: 150 монет)
- `duchess_distraction` - Герцогиня Отвлечения (HP: 200, урон: 5, награда: 200 монет)

---

### 6. villain_battles (72 kB)
**Миграция:** `030_villains.up.sql`
**Назначение:** Битвы пользователей со злодеями

#### Поля:

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | BIGSERIAL PRIMARY KEY | Уникальный идентификатор битвы |
| `child_profile_id` | UUID NOT NULL | ID профиля ребёнка, FK → child_profiles(id) ON DELETE CASCADE |
| `villain_id` | VARCHAR(100) NOT NULL | ID злодея, FK → villains(id) ON DELETE CASCADE |
| `status` | VARCHAR(20) NOT NULL DEFAULT 'active' | Статус: active, defeated, abandoned |
| `current_hp` | INTEGER NOT NULL | Текущее здоровье злодея |
| `total_damage_dealt` | INTEGER DEFAULT 0 | Всего нанесено урона, CHECK: >= 0 |
| `correct_tasks_count` | INTEGER DEFAULT 0 | Количество правильно решённых задач, CHECK: >= 0 |
| `rewards_claimed` | BOOLEAN DEFAULT FALSE | Получены ли награды за победу |
| `started_at` | TIMESTAMPTZ NOT NULL DEFAULT NOW() | Дата начала битвы |
| `defeated_at` | TIMESTAMPTZ | Дата победы над злодеем |
| `updated_at` | TIMESTAMPTZ NOT NULL DEFAULT NOW() | Дата обновления |

**Индексы:**
- `idx_villain_battles_profile_status` - (child_profile_id, status)
- `idx_villain_battles_active` - (child_profile_id) WHERE status = 'active'

**Constraints:**
- UNIQUE (child_profile_id, villain_id, status)

**Триггер:** `villain_battles_updated_at` - автообновление updated_at

---

### 7. damage_events (16 kB)
**Миграция:** `030_villains.up.sql`
**Назначение:** История нанесения урона злодеям

#### Поля:

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | BIGSERIAL PRIMARY KEY | Уникальный идентификатор события |
| `battle_id` | BIGINT NOT NULL | ID битвы, FK → villain_battles(id) ON DELETE CASCADE |
| `attempt_id` | UUID | ID попытки, FK → attempts(id) ON DELETE SET NULL |
| `damage` | INTEGER NOT NULL | Нанесённый урон, CHECK: >= 0 |
| `task_type` | VARCHAR(20) NOT NULL | Тип задачи: 'help' или 'check', CHECK: IN ('help', 'check') |
| `created_at` | TIMESTAMPTZ NOT NULL DEFAULT NOW() | Дата события |

**Индексы:**
- `idx_damage_events_battle` - (battle_id, created_at DESC)

---

### 8. subscription_plans (48 kB)
**Миграция:** `031_subscriptions.up.sql`
**Назначение:** Планы подписок (тарифы)

#### Поля:

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | VARCHAR(100) PRIMARY KEY | ID плана, например: 'monthly', 'yearly' |
| `name` | VARCHAR(200) NOT NULL | Название плана |
| `description` | TEXT NOT NULL | Описание плана |
| `price_cents` | INTEGER NOT NULL | Цена в копейках (499 руб = 49900), CHECK: >= 0 |
| `currency` | VARCHAR(10) NOT NULL DEFAULT 'RUB' | Валюта |
| `duration_days` | INTEGER NOT NULL | Длительность подписки в днях, CHECK: > 0 |
| `trial_days` | INTEGER DEFAULT 0 | Количество дней пробного периода, CHECK: >= 0 |
| `discount_percent` | INTEGER DEFAULT 0 | Процент скидки, CHECK: >= 0 AND <= 100 |
| `is_popular` | BOOLEAN DEFAULT FALSE | Популярный ли план (для UI) |
| `display_order` | INTEGER DEFAULT 0 | Порядок отображения |
| `is_active` | BOOLEAN DEFAULT TRUE | Активен ли план |
| `created_at` | TIMESTAMPTZ NOT NULL DEFAULT NOW() | Дата создания |

**Индексы:**
- `idx_subscription_plans_active` - (is_active, display_order)

**Примеры данных:**
- `monthly` - Месячная подписка (499 руб, 30 дней, 7 дней trial)
- `yearly` - Годовая подписка (3999 руб, 365 дней, 14 дней trial, скидка 33%)

---

### 9. subscriptions (32 kB)
**Миграция:** `031_subscriptions.up.sql`
**Назначение:** Активные подписки пользователей

#### Поля:

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | BIGSERIAL PRIMARY KEY | Уникальный идентификатор подписки |
| `child_profile_id` | UUID NOT NULL | ID профиля ребёнка, FK → child_profiles(id) ON DELETE CASCADE |
| `plan_id` | VARCHAR(100) NOT NULL | ID плана, FK → subscription_plans(id) |
| `status` | VARCHAR(20) NOT NULL DEFAULT 'trial' | Статус: trial, active, expired, cancelled |
| `started_at` | TIMESTAMPTZ NOT NULL DEFAULT NOW() | Дата начала подписки |
| `trial_ends_at` | TIMESTAMPTZ | Дата окончания trial периода |
| `expires_at` | TIMESTAMPTZ NOT NULL | Дата окончания подписки |
| `cancelled_at` | TIMESTAMPTZ | Дата отмены подписки |
| `auto_renew` | BOOLEAN DEFAULT TRUE | Автоматическое продление |
| `payment_provider` | VARCHAR(50) | Платёжная система: 'yookassa', 'stripe' |
| `payment_external_id` | VARCHAR(255) | ID в платежной системе |
| `created_at` | TIMESTAMPTZ NOT NULL DEFAULT NOW() | Дата создания |
| `updated_at` | TIMESTAMPTZ NOT NULL DEFAULT NOW() | Дата обновления |

**Индексы:**
- `idx_subscriptions_active_unique` - UNIQUE (child_profile_id) WHERE status IN ('trial', 'active')
- `idx_subscriptions_profile` - (child_profile_id, status)
- `idx_subscriptions_expires` - (expires_at) WHERE status IN ('trial', 'active')

**Триггер:** `subscriptions_updated_at` - автообновление updated_at

---

### 10. referral_codes (72 kB)
**Миграция:** `032_referrals.up.sql`
**Назначение:** Реферальные коды пользователей

#### Поля:

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | BIGSERIAL PRIMARY KEY | Уникальный идентификатор |
| `child_profile_id` | UUID NOT NULL | ID профиля владельца кода, FK → child_profiles(id) ON DELETE CASCADE |
| `code` | VARCHAR(20) NOT NULL UNIQUE | Уникальный код приглашения (например: 'ABCD1234') |
| `uses_count` | INTEGER DEFAULT 0 | Количество использований, CHECK: >= 0 |
| `max_uses` | INTEGER | Максимум использований (NULL = unlimited) |
| `created_at` | TIMESTAMPTZ NOT NULL DEFAULT NOW() | Дата создания |
| `expires_at` | TIMESTAMPTZ | Дата истечения кода |

**Индексы:**
- `idx_referral_codes_code` - (code)

**Constraints:**
- UNIQUE (child_profile_id) - один код на профиль

**Триггер:** Автоматическое создание кода при создании профиля (функция `generate_referral_code()`)

---

### 11. referrals (32 kB)
**Миграция:** `032_referrals.up.sql`
**Назначение:** Реферальные связи между пользователями

#### Поля:

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | BIGSERIAL PRIMARY KEY | Уникальный идентификатор |
| `referrer_id` | UUID NOT NULL | Кто пригласил, FK → child_profiles(id) ON DELETE CASCADE |
| `referred_id` | UUID NOT NULL | Кого пригласили, FK → child_profiles(id) ON DELETE CASCADE |
| `is_active` | BOOLEAN DEFAULT FALSE | Активировался ли приглашенный пользователь |
| `reward_coins` | INTEGER DEFAULT 50 | Награда за приглашение |
| `reward_claimed` | BOOLEAN DEFAULT FALSE | Забрана ли награда |
| `invited_at` | TIMESTAMPTZ NOT NULL DEFAULT NOW() | Дата приглашения |
| `activated_at` | TIMESTAMPTZ | Когда приглашенный стал активным |
| `reward_claimed_at` | TIMESTAMPTZ | Дата получения награды |

**Индексы:**
- `idx_referrals_referrer` - (referrer_id, is_active)
- `idx_referrals_referred` - (referred_id)

**Constraints:**
- UNIQUE (referrer_id, referred_id)

---

### 12. referral_milestones (48 kB)
**Миграция:** `032_referrals.up.sql`
**Назначение:** Награды за количество приглашенных друзей

#### Поля:

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | BIGSERIAL PRIMARY KEY | Уникальный идентификатор |
| `friends_count` | INTEGER NOT NULL UNIQUE | Количество друзей для достижения, CHECK: > 0 |
| `reward_coins` | INTEGER NOT NULL | Награда монетами, CHECK: > 0 |
| `description` | TEXT NOT NULL | Описание milestone |
| `created_at` | TIMESTAMPTZ NOT NULL DEFAULT NOW() | Дата создания |

**Примеры данных:**
- 1 друг → 50 монет
- 3 друга → 100 монет
- 5 друзей → 200 монет
- 10 друзей → 500 монет

---

### 13. child_referral_milestones (24 kB)
**Миграция:** `032_referrals.up.sql`
**Назначение:** Прогресс milestone для каждого пользователя

#### Поля:

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | BIGSERIAL PRIMARY KEY | Уникальный идентификатор |
| `child_profile_id` | UUID NOT NULL | ID профиля, FK → child_profiles(id) ON DELETE CASCADE |
| `milestone_id` | BIGINT NOT NULL | ID milestone, FK → referral_milestones(id) ON DELETE CASCADE |
| `is_claimed` | BOOLEAN DEFAULT FALSE | Забрана ли награда |
| `claimed_at` | TIMESTAMPTZ | Дата получения награды |

**Индексы:**
- `idx_child_referral_milestones_profile` - (child_profile_id, is_claimed)

**Constraints:**
- UNIQUE (child_profile_id, milestone_id)

---

### 14. parent_consents (80 kB)
**Миграция:** `033_parent_consents.up.sql`
**Назначение:** Согласия родителей на обработку данных (COPPA/GDPR compliance)

#### Поля:

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | UUID PRIMARY KEY | Уникальный идентификатор |
| `parent_user_id` | VARCHAR(255) NOT NULL | ID родителя на платформе (platform_user_id) |
| `platform_id` | VARCHAR(20) NOT NULL | Платформа: 'vk', 'telegram', 'max', 'web' |
| `privacy_policy_version` | VARCHAR(20) NOT NULL | Версия политики конфиденциальности (например: 1.0) |
| `privacy_policy_accepted` | BOOLEAN NOT NULL DEFAULT TRUE | Принята ли политика |
| `privacy_policy_accepted_at` | TIMESTAMPTZ NOT NULL DEFAULT NOW() | Дата принятия |
| `terms_version` | VARCHAR(20) NOT NULL | Версия условий использования (например: 1.0) |
| `terms_accepted` | BOOLEAN NOT NULL DEFAULT TRUE | Приняты ли условия |
| `terms_accepted_at` | TIMESTAMPTZ NOT NULL DEFAULT NOW() | Дата принятия |
| `adult_consent` | BOOLEAN NOT NULL DEFAULT FALSE | Подтверждение совершеннолетия родителя |
| `adult_consent_at` | TIMESTAMPTZ | Дата подтверждения |
| `ip_address` | VARCHAR(45) | IP-адрес (IPv4 или IPv6) для аудита |
| `user_agent` | TEXT | User-Agent для аудита |
| `created_at` | TIMESTAMPTZ NOT NULL DEFAULT NOW() | Дата создания |
| `updated_at` | TIMESTAMPTZ NOT NULL DEFAULT NOW() | Дата обновления |

**Индексы:**
- `idx_parent_consents_parent` - (platform_id, parent_user_id)
- `idx_parent_consents_created` - (created_at DESC)

**Constraints:**
- UNIQUE (platform_id, parent_user_id)

**Триггер:** `parent_consents_updated_at` - автообновление updated_at

---

### 15. email_verifications (80 kB)
**Миграция:** `034_legal_documents.up.sql`
**Назначение:** Верификация email адресов родителей

#### Поля:

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | UUID PRIMARY KEY | Уникальный идентификатор |
| `email` | VARCHAR(255) NOT NULL | Email для верификации |
| `verification_code` | VARCHAR(6) NOT NULL | 6-значный код для подтверждения |
| `is_verified` | BOOLEAN NOT NULL DEFAULT FALSE | Подтверждён ли email |
| `verified_at` | TIMESTAMPTZ | Дата подтверждения |
| `expires_at` | TIMESTAMPTZ NOT NULL | Код действителен 15 минут |
| `send_attempts` | INTEGER NOT NULL DEFAULT 1 | Количество отправок кода |
| `verify_attempts` | INTEGER NOT NULL DEFAULT 0 | Количество попыток ввода кода |
| `parent_user_id` | VARCHAR(255) | ID родителя для связи |
| `platform_id` | VARCHAR(20) | Платформа: 'vk', 'telegram', etc. |
| `ip_address` | VARCHAR(45) | IP-адрес для аудита |
| `created_at` | TIMESTAMPTZ NOT NULL DEFAULT NOW() | Дата создания |
| `updated_at` | TIMESTAMPTZ NOT NULL DEFAULT NOW() | Дата обновления |

**Индексы:**
- `idx_email_verifications_email` - (email)
- `idx_email_verifications_code` - (verification_code)
- `idx_email_verifications_expires` - (expires_at) WHERE is_verified = FALSE

**Триггер:** `email_verifications_updated_at` - автообновление updated_at

---

### 16. legal_documents (96 kB)
**Миграция:** `034_legal_documents.up.sql`
**Назначение:** Юридические документы (политика конфиденциальности, условия использования)

#### Поля:

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | UUID PRIMARY KEY | Уникальный идентификатор |
| `document_type` | VARCHAR(50) NOT NULL | Тип: 'privacy_policy' или 'terms_of_service' |
| `version` | VARCHAR(20) NOT NULL | Версия документа: '1.0', '1.1', '2.0' |
| `title` | VARCHAR(500) NOT NULL | Заголовок документа |
| `content` | TEXT NOT NULL | Полный текст документа в markdown |
| `language` | VARCHAR(10) NOT NULL DEFAULT 'ru' | Язык: 'ru', 'en' |
| `is_active` | BOOLEAN NOT NULL DEFAULT TRUE | Активная версия для отображения пользователям |
| `effective_date` | DATE NOT NULL | Дата вступления в силу |
| `created_at` | TIMESTAMPTZ NOT NULL DEFAULT NOW() | Дата создания |
| `updated_at` | TIMESTAMPTZ NOT NULL DEFAULT NOW() | Дата обновления |

**Индексы:**
- `idx_legal_documents_type_active` - (document_type, is_active) WHERE is_active = TRUE
- `idx_legal_documents_effective_date` - (effective_date DESC)

**Constraints:**
- UNIQUE (document_type, version, language)

**Триггер:** `legal_documents_updated_at` - автообновление updated_at

---

## 🟢 УЧЕБНИКИ - Индексированная база заданий (3 таблицы)

### 17. textbooks (48 kB)
**Миграция:** `020_textbooks.up.sql`
**Назначение:** Справочник учебников

#### Поля:

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | BIGSERIAL PRIMARY KEY | Уникальный идентификатор |
| `subject` | VARCHAR(50) NOT NULL | Предмет: 'math', 'russian', etc. |
| `grade` | INTEGER NOT NULL | Класс обучения |
| `authors` | VARCHAR(255) NOT NULL | Авторы учебника |
| `title` | VARCHAR(255) NOT NULL | Название учебника |
| `part` | INTEGER DEFAULT NULL | Часть учебника (1, 2, 3, NULL) |
| `year` | INTEGER DEFAULT NULL | Год издания |
| `publisher` | VARCHAR(255) DEFAULT NULL | Издательство |
| `source_url` | VARCHAR(500) DEFAULT NULL | URL источника |
| `created_at` | TIMESTAMPTZ NOT NULL DEFAULT NOW() | Дата создания |

**Constraints:**
- UNIQUE (subject, grade, authors, part)

---

### 18. textbook_tasks (2.4 MB) 🔥 САМАЯ БОЛЬШАЯ
**Миграция:** `020_textbooks.up.sql` + 022-024 (данные Peterson)
**Назначение:** База заданий из учебников с метаданными

#### Поля:

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | BIGSERIAL PRIMARY KEY | Уникальный идентификатор |
| `textbook_id` | BIGINT NOT NULL | ID учебника, FK → textbooks(id) ON DELETE CASCADE |
| `page_number` | INTEGER NOT NULL | Номер страницы |
| `task_number` | VARCHAR(20) NOT NULL | Номер задания |
| `task_order` | INTEGER NOT NULL DEFAULT 0 | Порядок задания на странице |
| `condition_text` | TEXT | Условие задачи (plain text) |
| `condition_html` | TEXT | Условие задачи (HTML) |
| `solution_text` | TEXT | Решение (plain text) |
| `solution_html` | TEXT | Решение (HTML) |
| `hints_text` | TEXT | Подсказки (plain text) |
| `hints_html` | TEXT | Подсказки (HTML) |
| `has_sub_items` | BOOLEAN NOT NULL DEFAULT FALSE | Есть ли подпункты (a, б, в) |
| `sub_items_json` | JSONB DEFAULT NULL | JSON с подпунктами |
| `source_url` | VARCHAR(500) DEFAULT NULL | URL источника |
| `created_at` | TIMESTAMPTZ NOT NULL DEFAULT NOW() | Дата создания |

**Индексы:**
- `idx_textbook_tasks_textbook_page` - (textbook_id, page_number)
- `idx_textbook_tasks_task_number` - (textbook_id, task_number)

**Данные:** ~5000+ заданий из Peterson 3 класс

---

### 19. textbook_task_images (336 kB)
**Миграция:** `020_textbooks.up.sql`
**Назначение:** Изображения для заданий

#### Поля:

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | BIGSERIAL PRIMARY KEY | Уникальный идентификатор |
| `task_id` | BIGINT NOT NULL | ID задания, FK → textbook_tasks(id) ON DELETE CASCADE |
| `image_type` | VARCHAR(20) NOT NULL | Тип: 'condition', 'solution', 'hint' |
| `image_order` | INTEGER NOT NULL DEFAULT 0 | Порядок изображения |
| `sub_item_letter` | VARCHAR(5) DEFAULT NULL | Буква подпункта (a, б, в) |
| `original_url` | VARCHAR(500) NOT NULL | Оригинальный URL |
| `local_path` | VARCHAR(500) DEFAULT NULL | Локальный путь |
| `alt_text` | VARCHAR(500) DEFAULT NULL | Альтернативный текст |
| `width` | INTEGER DEFAULT NULL | Ширина изображения |
| `height` | INTEGER DEFAULT NULL | Высота изображения |
| `file_size` | INTEGER DEFAULT NULL | Размер файла в байтах |
| `created_at` | TIMESTAMPTZ NOT NULL DEFAULT NOW() | Дата создания |

**Индексы:**
- `idx_textbook_task_images_task` - (task_id, image_type, image_order)

---

## 🔴 DEPRECATED - Таблицы Telegram бота (7 таблиц)

### 20. parsed_tasks (40 kB)
**Миграция:** `001_parsed_tasks.up.sql`
**Назначение:** Распознанные задания из Telegram бота (v1)
**Статус:** 🔴 DEPRECATED - заменена на `attempts`

### 21. hints_cache (16 kB)
**Миграция:** `002_hints_cache.up.sql`
**Назначение:** Кеш подсказок для заданий (Telegram v1)
**Статус:** 🔴 DEPRECATED - данные теперь в `attempts.hints_result`

### 22. metrics_events (56 kB)
**Миграция:** `003_metrics.up.sql`
**Назначение:** События метрик (Telegram v1)
**Статус:** 🔴 DEPRECATED - заменена на analytics в REST API

### 23. timeline_events (40 kB)
**Миграция:** `006_timeline_events.up.sql`
**Назначение:** События в timeline пользователя (Telegram v1)
**Статус:** 🔴 DEPRECATED - timeline для Telegram бота

### 24. task_sessions (24 kB)
**Миграция:** `007_task_sessions.up.sql`
**Назначение:** Сессии решения заданий (Telegram v2)
**Статус:** 🔴 DEPRECATED - заменены на `attempts`

### 25. chat (16 kB)
**Миграция:** `015_chat.up.sql`
**Назначение:** Сообщения чата (Telegram)
**Статус:** 🔴 DEPRECATED - история чата Telegram бота

### 26. user (8 kB)
**Миграция:** `016_user.up.sql`
**Назначение:** Пользователи Telegram бота
**Статус:** 🔴 DEPRECATED - заменена на `child_profiles`

---

## 📊 СЛУЖЕБНЫЕ

### 27. schema_migrations (24 kB)
**Назначение:** Отслеживание применённых миграций
**Статус:** ✅ СЛУЖЕБНАЯ - необходима для миграций

---

## 📊 Итоговая статистика

### По статусу:
- ✅ **Активные (REST API):** 16 таблиц
- ✅ **Активные (Учебники):** 3 таблицы
- 🔴 **Deprecated (Telegram):** 7 таблиц (можно удалить)
- ✅ **Служебные:** 1 таблица

### По размеру (Top-5):
1. `textbook_task_index` - 2.6 MB (индексы)
2. `textbook_tasks` - 2.4 MB (5000+ заданий)
3. `attempts` - 448 kB (попытки пользователей)
4. `textbook_task_images` - 336 kB (изображения)
5. `legal_documents` - 96 kB (юридические документы)

### Рекомендации по оптимизации:
1. **attempts** - добавить партиционирование по дате для старых записей
2. **textbook_task_index** - рассмотреть сжатие embeddings
3. **child_achievements** - добавить индекс для незавершённых достижений

---

## 🔥 Следующие шаги

Согласно `CLEANUP_RECOMMENDATIONS.md`:

**Этап 2: Удаление deprecated таблиц БД** (~200 kB)
- `parsed_tasks`, `hints_cache`, `metrics_events`
- `timeline_events`, `task_sessions`, `chat`, `user`

**Требования:**
1. Создать backup базы данных
2. Создать миграцию `036_cleanup_telegram_tables.up.sql`
3. Применить миграцию
4. Проверить работу REST API
