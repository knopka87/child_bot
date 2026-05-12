# Database Schema Reference

**Database:** `child_bot`
**User:** `child_bot`
**Connection:** `docker exec -i child_bot_postgres psql -U child_bot -d child_bot`

## Tables Overview

```
achievements           - Справочник достижений
attempts               - Попытки решения задач
child_achievements     - Разблокированные достижения пользователей
child_profiles         - Профили детей (основная таблица)
damage_events          - События урона по злодеям
email_verifications    - Верификация email
legal_documents        - Юридические документы
parent_consent_history - История согласий родителей
parent_consents        - Согласия родителей
payment_events         - События платежей
payments               - Платежи
referral_codes         - Реферальные коды
referrals              - Реферальные связи
report_settings        - Настройки отчетов
subscription_plans     - Тарифные планы
subscriptions          - Подписки
textbook_task_images   - Изображения задач из учебников
textbook_task_index    - Индекс задач учебников
textbook_tasks         - Задачи из учебников
textbooks              - Учебники
villain_battles        - Битвы со злодеями
villains               - Злодеи
weekly_reports         - Еженедельные отчеты
```

## Main Table: child_profiles

### Columns
```
id                          UUID PRIMARY KEY (auto-generated)
display_name                VARCHAR(100) NOT NULL
avatar_id                   VARCHAR(50)
grade                       INTEGER NOT NULL (1-4)
email                       VARCHAR(255)
email_verified              BOOLEAN DEFAULT false
platform_id                 VARCHAR(20) NOT NULL (vk, telegram, max, web)
platform_user_id            VARCHAR(255)
level                       INTEGER DEFAULT 1
experience_points           INTEGER DEFAULT 0
coins_balance               INTEGER DEFAULT 0
tasks_solved_total          INTEGER DEFAULT 0
tasks_solved_correct        INTEGER DEFAULT 0
hints_used_total            INTEGER DEFAULT 0
streak_days                 INTEGER DEFAULT 0
last_activity_at            TIMESTAMPTZ
created_at                  TIMESTAMPTZ NOT NULL DEFAULT now()
updated_at                  TIMESTAMPTZ NOT NULL DEFAULT now()
achievements_last_viewed_at TIMESTAMPTZ DEFAULT now()
```

### Indexes
- PRIMARY KEY: `id`
- UNIQUE: `(platform_id, platform_user_id)`
- INDEX: `last_activity_at DESC`
- INDEX: `created_at DESC`
- INDEX: `level DESC`
- INDEX: `(platform_id, platform_user_id)`

### Foreign Keys (Referenced by)
- attempts.child_profile_id
- child_achievements.child_profile_id
- payments.child_profile_id
- referral_codes.child_profile_id
- referrals.referred_id
- referrals.referrer_id
- report_settings.child_profile_id
- subscriptions.child_profile_id
- villain_battles.child_profile_id
- weekly_reports.user_id

### Constraints
- grade: 1-4
- level: >= 1
- coins_balance: >= 0
- experience_points: >= 0
- streak_days: >= 0
- All counter fields: >= 0

## Common Queries

### Find profile by ID
```sql
SELECT * FROM child_profiles WHERE id = 'uuid-here';
```

### Find profile by platform credentials
```sql
SELECT * FROM child_profiles
WHERE platform_id = 'vk' AND platform_user_id = '123456';
```

### Get profile with stats
```sql
SELECT
  id, display_name, platform_id, platform_user_id,
  level, experience_points, coins_balance,
  tasks_solved_correct, streak_days,
  created_at, last_activity_at
FROM child_profiles
WHERE id = 'uuid-here';
```

### Check for duplicates (same user, different platforms)
```sql
SELECT platform_id, platform_user_id, id, display_name, created_at
FROM child_profiles
WHERE platform_user_id = '123456'
ORDER BY created_at;
```

### Get profile with all related data
```sql
SELECT
  (SELECT COUNT(*) FROM attempts WHERE child_profile_id = cp.id) as attempts_count,
  (SELECT COUNT(*) FROM child_achievements WHERE child_profile_id = cp.id) as achievements_count,
  (SELECT COUNT(*) FROM referrals WHERE referrer_id = cp.id) as referrals_count,
  cp.*
FROM child_profiles cp
WHERE cp.id = 'uuid-here';
```
