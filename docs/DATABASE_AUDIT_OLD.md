# Аудит базы данных - Homework Helper Bot

**Дата:** 2026-04-04
**Всего таблиц:** 28

## Классификация таблиц

### 🟢 АКТИВНЫЕ - Используются в текущей версии (REST API)

#### 1. **child_profiles** (96 kB)
**Миграция:** 027_child_profiles.up.sql
**Назначение:** Профили детей-пользователей
**Основные поля:**
- `id` UUID - уникальный ID профиля
- `display_name` VARCHAR(100) - имя ребёнка
- `avatar_url` TEXT - URL аватара
- `grade` INTEGER (1-4) - класс обучения
- `platform_id` VARCHAR(50) - платформа (web, vk, max, telegram)
- `energy` INTEGER - энергия для геймификации (0-100)
- Счетчики: `coins`, `streak_days`, `total_tasks_solved`, и т.д.

**Статус:** ✅ АКТИВНА - центральная таблица пользователей

---

#### 2. **attempts** (448 kB)
**Миграция:** 028_attempts.up.sql
**Назначение:** Попытки решения заданий (help/check)
**Основные поля:**
- `id` UUID - ID попытки
- `child_profile_id` UUID - связь с профилем
- `type` VARCHAR(20) - тип (help, check)
- `status` VARCHAR(20) - статус (created, processing, completed, failed)
- `task_image_data` TEXT - base64 изображение задания
- `answer_image_data` TEXT - base64 изображение ответа
- `hints_result` JSONB - результаты LLM (подсказки)
- `check_result` JSONB - результаты проверки

**Статус:** ✅ АКТИВНА - основная таблица для help/check flow

---

#### 3. **achievements** (48 kB)
**Миграция:** 029_achievements.up.sql
**Назначение:** Справочник достижений
**Основные поля:**
- `id` VARCHAR(100) - ID достижения
- `name`, `description`, `icon_url`
- `type` VARCHAR(50) - категория
- `condition_type` VARCHAR(50) - условие получения
- `threshold` INTEGER - порог для получения

**Статус:** ✅ АКТИВНА - справочник достижений

---

#### 4. **child_achievements** (32 kB)
**Миграция:** 029_achievements.up.sql
**Назначение:** Связь детей и их достижений
**Основные поля:**
- `child_profile_id` UUID
- `achievement_id` VARCHAR(100)
- `earned_at` TIMESTAMPTZ
- `progress` INTEGER - прогресс к получению

**Статус:** ✅ АКТИВНА - прогресс достижений

---

#### 5. **villains** (48 kB)
**Миграция:** 030_villains.up.sql
**Назначение:** Справочник злодеев (геймификация)
**Основные поля:**
- `id` VARCHAR(100) - ID злодея
- `name`, `description`, `image_url`
- `max_hp` INTEGER - максимальное здоровье
- `level` INTEGER - уровень сложности
- `damage_per_correct_task` INTEGER - урон за правильный ответ
- `unlock_order` INTEGER - порядок появления
- `reward_coins`, `reward_achievement_id`

**Статус:** ✅ АКТИВНА - геймификация

---

#### 6. **villain_battles** (72 kB)
**Миграция:** 030_villains.up.sql
**Назначение:** Битвы пользователей со злодеями
**Основные поля:**
- `child_profile_id` UUID
- `villain_id` VARCHAR(100)
- `status` VARCHAR(20) - active, defeated, abandoned
- `current_hp` INTEGER - текущее здоровье злодея
- `total_damage_dealt` INTEGER
- `correct_tasks_count` INTEGER

**Статус:** ✅ АКТИВНА - прогресс битв

---

#### 7. **damage_events** (16 kB)
**Миграция:** 030_villains.up.sql
**Назначение:** История нанесения урона злодеям
**Основные поля:**
- `battle_id` BIGINT - ссылка на битву
- `attempt_id` UUID - ссылка на попытку
- `damage` INTEGER
- `task_type` VARCHAR(20) - help/check

**Статус:** ✅ АКТИВНА - история урона

---

#### 8. **subscription_plans** (48 kB)
**Миграция:** 031_subscriptions.up.sql
**Назначение:** Планы подписок (тарифы)
**Основные поля:**
- `id` VARCHAR(100)
- `name`, `description`
- `price_monthly` INTEGER - цена в копейках
- `features` JSONB - список возможностей

**Статус:** ✅ АКТИВНА - монетизация

---

#### 9. **subscriptions** (32 kB)
**Миграция:** 031_subscriptions.up.sql
**Назначение:** Активные подписки пользователей
**Основные поля:**
- `child_profile_id` UUID
- `plan_id` VARCHAR(100)
- `status` VARCHAR(20) - active, cancelled, expired
- `current_period_start`, `current_period_end`
- `auto_renew` BOOLEAN

**Статус:** ✅ АКТИВНА - подписки пользователей

---

#### 10. **referral_codes** (72 kB)
**Миграция:** 032_referrals.up.sql
**Назначение:** Реферальные коды
**Основные поля:**
- `code` VARCHAR(20) UNIQUE - уникальный код
- `owner_id` UUID - владелец кода
- `uses_count` INTEGER - сколько раз использован
- `max_uses` INTEGER - лимит использований

**Статус:** ✅ АКТИВНА - реферальная система

---

#### 11. **referrals** (32 kB)
**Миграция:** 032_referrals.up.sql
**Назначение:** Связи рефералов
**Основные поля:**
- `referrer_id` UUID - кто пригласил
- `referee_id` UUID - кого пригласили
- `referral_code` VARCHAR(20)
- `activated` BOOLEAN - активирован ли реферал
- `activation_threshold` INTEGER - порог активации

**Статус:** ✅ АКТИВНА - отслеживание рефералов

---

#### 12. **referral_milestones** (48 kB)
**Миграция:** 032_referrals.up.sql
**Назначение:** Справочник milestone'ов реферальной программы
**Основные поля:**
- `id` BIGSERIAL
- `milestone_count` INTEGER - количество рефералов
- `reward_type` VARCHAR(50) - тип награды
- `reward_value` INTEGER

**Статус:** ✅ АКТИВНА - вознаграждения за рефералов

---

#### 13. **child_referral_milestones** (24 kB)
**Миграция:** 032_referrals.up.sql
**Назначение:** Достигнутые milestone'ы детей
**Основные поля:**
- `child_profile_id` UUID
- `milestone_id` BIGINT
- `achieved_at` TIMESTAMPTZ

**Статус:** ✅ АКТИВНА - прогресс milestone'ов

---

#### 14. **parent_consents** (80 kB)
**Миграция:** 033_parent_consents.up.sql
**Назначение:** Согласия родителей
**Основные поля:**
- `child_profile_id` UUID
- `parent_email` VARCHAR(255)
- `email_verified` BOOLEAN
- `consent_given` BOOLEAN
- `consent_text` TEXT - текст согласия
- `ip_address` VARCHAR(45)

**Статус:** ✅ АКТИВНА - COPPA/GDPR compliance

---

#### 15. **email_verifications** (80 kB)
**Миграция:** 033_parent_consents.up.sql (создаётся внутри)
**Назначение:** Токены для верификации email
**Основные поля:**
- `child_profile_id` UUID
- `email` VARCHAR(255)
- `verification_token` VARCHAR(100)
- `expires_at` TIMESTAMPTZ
- `verified_at` TIMESTAMPTZ

**Статус:** ✅ АКТИВНА - верификация email родителей

---

#### 16. **legal_documents** (96 kB)
**Миграция:** 034_legal_documents.up.sql
**Назначение:** Версии юридических документов
**Основные поля:**
- `id` BIGSERIAL
- `document_type` VARCHAR(50) - privacy_policy, terms_of_service
- `version` VARCHAR(20)
- `content` TEXT - содержимое документа
- `effective_date` TIMESTAMPTZ

**Статус:** ✅ АКТИВНА - юридические документы

---

### 🟡 TELEGRAM BOT - Старые таблицы из Telegram-версии

#### 17. **parsed_tasks** (40 kB)
**Миграция:** 001_parsed_tasks.up.sql
**Назначение:** Распознанные задания из Telegram бота (v1)
**Основные поля:**
- `id` BIGSERIAL
- `user_id` BIGINT - Telegram user_id
- `raw_text` TEXT - сырой текст задания
- `parsed_data` JSONB - результат парсинга
- `state` VARCHAR(50)

**Статус:** 🔴 DEPRECATED - использовалась в Telegram боте v1, заменена на `attempts`

---

#### 18. **hints_cache** (16 kB)
**Миграция:** 002_hints_cache.up.sql
**Назначение:** Кеш подсказок для заданий (Telegram v1)
**Основные поля:**
- `task_id` BIGINT - ссылка на parsed_tasks
- `hint_level` INTEGER (1-3)
- `hint_text` TEXT

**Статус:** 🔴 DEPRECATED - кеш для старого бота, данные теперь в `attempts.hints_result`

---

#### 19. **metrics_events** (56 kB)
**Миграция:** 003_metrics.up.sql
**Назначение:** События метрик (Telegram v1)
**Основные поля:**
- `event_type` VARCHAR(100)
- `user_id` BIGINT - Telegram user_id
- `event_data` JSONB

**Статус:** 🔴 DEPRECATED - старая система метрик, заменена на analytics в REST API

---

#### 20. **timeline_events** (40 kB)
**Миграция:** 006_timeline_events.up.sql
**Назначение:** События в timeline пользователя (Telegram v1)
**Основные поля:**
- `user_id` BIGINT
- `event_type` VARCHAR(100)
- `event_data` JSONB

**Статус:** 🔴 DEPRECATED - timeline для Telegram бота

---

#### 21. **task_sessions** (24 kB)
**Миграция:** 007_task_sessions.up.sql
**Назначение:** Сессии решения заданий (Telegram v2)
**Основные поля:**
- `session_id` UUID
- `user_id` BIGINT
- `state` VARCHAR(50)
- `session_data` JSONB

**Статус:** 🔴 DEPRECATED - сессии для Telegram бота v2, заменены на `attempts`

---

#### 22. **chat** (16 kB)
**Миграция:** 015_chat.up.sql
**Назначение:** Сообщения чата (Telegram)
**Основные поля:**
- `user_id` BIGINT
- `message_text` TEXT
- `role` VARCHAR(20) - user/assistant

**Статус:** 🔴 DEPRECATED - история чата Telegram бота

---

#### 23. **user** (8 kB)
**Миграция:** 016_user.up.sql
**Назначение:** Пользователи Telegram бота
**Основные поля:**
- `user_id` BIGINT PRIMARY KEY - Telegram user_id
- `username` VARCHAR(255)
- `first_name`, `last_name`

**Статус:** 🔴 DEPRECATED - заменена на `child_profiles`

---

### 🟢 УЧЕБНИКИ - Индексированная база заданий Peterson

#### 24. **textbooks** (48 kB)
**Миграция:** 020_textbooks.up.sql
**Назначение:** Справочник учебников
**Основные поля:**
- `id` VARCHAR(100)
- `title` VARCHAR(255) - "Петерсон 3 класс"
- `author` VARCHAR(255)
- `grade` INTEGER

**Статус:** ✅ АКТИВНА - справочник учебников для поиска

---

#### 25. **textbook_tasks** (2.4 MB) 🔥 САМАЯ БОЛЬШАЯ
**Миграция:** 020_textbooks.up.sql + 022-024 (данные Peterson)
**Назначение:** База заданий из учебников с метаданными
**Основные поля:**
- `id` BIGSERIAL
- `textbook_id` VARCHAR(100)
- `page_number` INTEGER
- `task_number` VARCHAR(20)
- `content` TEXT - описание задания
- `topics` VARCHAR(255)[] - массив тем
- `difficulty` VARCHAR(20)

**Статус:** ✅ АКТИВНА - база для поиска похожих заданий
**Данные:** ~5000+ заданий из Peterson 3 класс

---

#### 26. **textbook_task_images** (336 kB)
**Миграция:** 020_textbooks.up.sql
**Назначение:** Изображения заданий
**Основные поля:**
- `task_id` BIGINT
- `image_url` TEXT
- `image_type` VARCHAR(50)

**Статус:** ✅ АКТИВНА - изображения для заданий

---

#### 27. **textbook_task_index** (2.6 MB) 🔥 САМАЯ БОЛЬШАЯ
**Миграция:** 021_textbook_task_index.up.sql
**Назначение:** Полнотекстовый индекс для поиска заданий
**Основные поля:**
- `task_id` BIGINT
- `search_vector` tsvector - индекс для поиска
- `embedding` vector(1536) - векторные эмбеддинги

**Статус:** ✅ АКТИВНА - индекс для быстрого поиска похожих заданий

---

### 📊 СЛУЖЕБНЫЕ

#### 28. **schema_migrations** (24 kB)
**Назначение:** Отслеживание применённых миграций
**Статус:** ✅ СЛУЖЕБНАЯ - необходима для миграций

---

## 📊 Итоговая статистика

### По статусу:
- ✅ **Активные (REST API):** 16 таблиц
- ✅ **Активные (Учебники):** 3 таблицы
- 🔴 **Deprecated (Telegram):** 7 таблиц
- ✅ **Служебные:** 1 таблица

### По размеру (Top-5):
1. `textbook_task_index` - 2.6 MB (индексы)
2. `textbook_tasks` - 2.4 MB (данные заданий)
3. `attempts` - 448 kB (попытки пользователей)
4. `textbook_task_images` - 336 kB (изображения)
5. `legal_documents` - 96 kB (юридические документы)

---

## 🔥 Рекомендации по очистке

### Можно удалить (Telegram v1/v2):
```sql
-- Старые таблицы Telegram бота (сохранить данные в архив перед удалением!)
DROP TABLE IF EXISTS parsed_tasks CASCADE;
DROP TABLE IF EXISTS hints_cache CASCADE;
DROP TABLE IF EXISTS metrics_events CASCADE;
DROP TABLE IF EXISTS timeline_events CASCADE;
DROP TABLE IF EXISTS task_sessions CASCADE;
DROP TABLE IF EXISTS chat CASCADE;
DROP TABLE IF EXISTS user CASCADE;
```

**Экономия места:** ~200 kB (минимум, так как таблицы почти пустые)

**⚠️ ВАЖНО:** Перед удалением создать backup:
```bash
pg_dump -t parsed_tasks -t hints_cache -t metrics_events \
  -t timeline_events -t task_sessions -t chat -t user \
  "$DATABASE_URL" > telegram_bot_archive.sql
```

---

## 📝 Комментарии

### Дубликаты отсутствуют
Все таблицы имеют уникальное назначение. Нет дублирования функциональности между `user` (Telegram) и `child_profiles` (REST API) - это разные системы аутентификации.

### Потенциальные улучшения
1. **attempts** (448 kB) - можно добавить партиционирование по дате для старых записей
2. **textbook_task_index** - рассмотреть сжатие или удаление старых embeddings
3. Добавить индексы для `child_achievements.progress` для быстрого поиска незавершённых достижений

---

**Итого:** БД в хорошем состоянии. 7 устаревших таблиц от Telegram бота можно безопасно удалить после архивирования.
