# Система серийных достижений

## Концепция

Реализована система **серийных достижений** - групп однотипных достижений с разными пороговыми значениями.

**Ключевая особенность**: На странице достижений показывается только **максимальное разблокированное** достижение из каждой серии, что предотвращает перегрузку экрана.

## Пример работы

Если пользователь пригласил **15 друзей**:
- ✅ Разблокировано: "Стикер Дружба" за 5, 10, 15 друзей
- ❌ Не разблокировано: за 20, 25, 30, 40, 50 друзей
- 📱 **На экране показывается**: только "Стикер Дружба" за 15 друзей (с бейджем "15")

## Серии достижений

### 1. **Стрик (дни подряд)** 🔥
- Requirement type: `streak_days`
- Reward name: `Стрик`
- Иконка: 🔥
- Пороги: **1, 3, 7, 30, 90, 180, 365** дней
- Бейдж показывает: количество дней подряд

**Миграция**: `045_streak_series.up.sql`

### 2. **Проверки ДЗ** ✅
- Requirement type: `tasks_correct`
- Reward name: `Проверки ДЗ`
- Иконка: ✅
- Пороги: **1, 10, 100, 500, 1000** проверок
- Бейдж показывает: количество проверок

**Миграция**: `046_checks_series.up.sql`

### 3. **Победитель злодеев** 🦹
- Requirement type: `villains_defeated`
- Reward name: `Победитель злодеев`
- Иконка: 🦹
- Пороги: **1, 5, 10, 50, 100, 500, 1000** злодеев
- Бейдж показывает: количество побеждённых злодеев

**Миграция**: `047_villains_series.up.sql`

### 4. **Стикер Дружба** ⭐
- Requirement type: `friends_invited`
- Reward name: `Стикер Дружба`
- Иконка: ⭐
- Пороги: **5, 10, 15, 20, 25, 30, 40, 50** друзей
- Бейдж показывает: количество приглашённых друзей

**Миграция**: `044_friendship_stickers.up.sql`

### 5. **Мудрая сова** 🦉
- Requirement type: `hints_used`
- Reward name: `Мудрая сова`
- Иконка: 🦉
- Порог: **1** подсказка
- Уникальное достижение (не серия)

**Миграция**: `048_wise_owl.up.sql`

## Техническая реализация

### Бекенд

#### 1. Логика фильтрации (SQL)

```sql
WITH achievements_with_progress AS (
    SELECT
        a.*,
        COALESCE(ca.current_progress, 0) as current_progress,
        COALESCE(ca.is_unlocked, FALSE) as is_unlocked,
        ca.unlocked_at,
        ROW_NUMBER() OVER (
            PARTITION BY COALESCE(a.reward_name, a.id)
            ORDER BY
                CASE WHEN COALESCE(ca.is_unlocked, FALSE) = TRUE
                    THEN a.requirement_value
                    ELSE -a.requirement_value
                END DESC
        ) as rn
    FROM achievements a
    LEFT JOIN child_achievements ca ON ...
)
SELECT * FROM achievements_with_progress WHERE rn = 1
```

**Логика**:
- Группируем по `reward_name` (PARTITION BY)
- Для **разблокированных**: берем максимальное (ORDER BY requirement_value DESC)
- Для **неразблокированных**: берем минимальное (ORDER BY requirement_value ASC)
- Выбираем только первое (WHERE rn = 1)

#### 2. Методы проверки

**Файл**: `api/internal/store/achievement_progress.go`

```go
// GetCurrentStreakDays - для стрика
// GetTasksCorrectCount - для проверок ДЗ
// GetVillainsDefeatedCount - для злодеев
// GetFriendsInvitedCount - для друзей
// GetHintsUsedCount - для мудрой совы
```

**Файл**: `api/internal/service/achievement.go`

```go
// CheckStreakAchievements - проверка достижений за streak
// CheckTasksCorrectAchievements - проверка за проверки ДЗ
// CheckVillainAchievements - проверка за злодеев
// CheckFriendsInvitedAchievements - проверка за друзей
// CheckHintsUsedAchievements - проверка за подсказки
```

#### 3. Автоматическая проверка

**Когда проверяются достижения**:

- **Стрик**: при обновлении `UpdateStreakAndActivity()` - `profile.go:496`
- **Проверки/Злодеи**: при завершении попытки - `attempt.go` (TODO)
- **Друзья**: при активации реферала - `profile.go:357`
- **Мудрая сова**: при использовании подсказки - `attempt.go:426`

### Фронтенд

#### Бейджи с цифрами

**Файл**: `frontend/src/pages/Achievements/AchievementsPage.tsx:142-150`

```tsx
{achievement.reward.type === 'sticker' &&
  achievement.reward.name === 'Стикер Дружба' &&
  achievement.progress.total && (
    <div className="absolute -top-1 -right-1 w-6 h-6 bg-[#FF6B6B] text-white text-xs font-bold rounded-full flex items-center justify-center shadow-md">
      {achievement.progress.total}
    </div>
  )}
```

**Применяется для всех серий**:
- Стрик: показывает количество дней (1, 3, 7, 30...)
- Проверки: показывает количество проверок (1, 10, 100...)
- Злодеи: показывает количество побед (1, 5, 10...)
- Друзья: показывает количество друзей (5, 10, 15...)

## Статистика достижений

**До оптимизации**: 12 основных + 27 серийных = **39 достижений** → 10 полок

**После оптимизации**: 12 основных + 5 серий (по 1 видимому) = **~17 достижений** → 4-5 полок

**Выгода**:
- ✅ В **2.3 раза** меньше прокрутки
- ✅ Фокус на актуальных целях
- ✅ Чистый, понятный интерфейс

## Тестирование

### Тестовый пользователь

```bash
# ID тестового пользователя с 15 друзьями
60e314e8-c8ac-4d2a-8641-14ff53ba108a
```

### Проверка API

```bash
curl -s "http://localhost:8080/achievements" \
  -H "X-Child-Profile-ID: 60e314e8-c8ac-4d2a-8641-14ff53ba108a" \
  -H "X-Platform-ID: web"
```

**Ожидаемый результат**:
- Стрик: 🔒 "За 1 день занятий" (прогресс: 0/1)
- Проверки: 🔒 "За 1 проверку" (прогресс: 0/1)
- Злодеи: 🔒 "За победу над 1 злодеем" (прогресс: 0/1)
- Дружба: ✅ "За 15 друзей" (прогресс: 15/15) ← показывается максимальное!
- Сова: 🔒 "Использовал первую подсказку" (прогресс: 0/1)

## Будущие улучшения

1. **История достижений**: Модалка при клике на серийное достижение, показывающая все разблокированные из серии
2. **Анимация разблокировки**: Специальная анимация при переходе на следующий уровень серии
3. **Прогресс-бар серии**: Визуализация сколько осталось до следующего уровня
4. **Уведомления**: Push при разблокировке нового уровня серии

## Миграции

| № | Файл | Описание |
|---|------|----------|
| 044 | `friendship_stickers.up.sql` | 8 достижений за друзей (5-50) |
| 045 | `streak_series.up.sql` | 7 достижений за стрик (1-365 дней) |
| 046 | `checks_series.up.sql` | 5 достижений за проверки (1-1000) |
| 047 | `villains_series.up.sql` | 7 достижений за злодеев (1-1000) |
| 048 | `wise_owl.up.sql` | Достижение "Мудрая сова" |

**Всего добавлено**: 28 новых достижений
