# Система Достижений

## Обзор

Автоматическая система проверки и разблокировки достижений. Проверки происходят в реальном времени при выполнении пользователем определённых действий.

## Типы Достижений

### 1. `streak_days` - Дни подряд
**Триггер:** Обновление streak при заходе пользователя на главную страницу
**Метод проверки:** `AchievementService.CheckStreakAchievements()`
**Когда вызывается:** После успешного обновления `streak_days` в `ProfileService.UpdateStreakAndActivity()`

**Примеры:**
- `streak_3` - 3 дня подряд → 30 монет
- `streak_7` - 7 дней подряд → 70 монет
- `streak_14` - 14 дней подряд → 140 монет
- `streak_30` - 30 дней подряд → 300 монет

### 2. `villains_defeated` - Побеждённые монстры
**Триггер:** Победа над монстром (HP <= 0)
**Метод проверки:** `AchievementService.CheckVillainAchievements()`
**Когда вызывается:** После победы над монстром в `VillainService.DealDamageToVillain()`

**Примеры:**
- `villains_1` - Первый монстр → 50 монет
- `villains_3` - 3 монстра → 150 монет
- `villains_5` - 5 монстров → 250 монет
- `villains_10` - 10 монстров → 500 монет

### 3. `tasks_correct` - Правильно решённые задачи
**Триггер:** Правильное решение задачи (check_result.decision = 'correct')
**Метод проверки:** `AchievementService.CheckTasksCorrectAchievements()`
**Когда вызывается:** После правильного ответа в `AttemptService.ProcessCheck()`

**Примеры:**
- `tasks_10` - 10 задач → 100 монет
- `tasks_25` - 25 задач → 250 монет
- `tasks_50` - 50 задач → 500 монет
- `tasks_100` - 100 задач → 1000 монет

### 4. `tasks_no_hints` - Задачи без подсказок
**Триггер:** Правильное решение задачи без использования подсказок (hints_used = 0)
**Метод проверки:** `AchievementService.CheckTasksNoHintsAchievements()`
**Когда вызывается:** После правильного ответа в `AttemptService.ProcessCheck()`

**Примеры:**
- `no_hints_5` - 5 задач без подсказок → 100 монет
- `no_hints_10` - 10 задач без подсказок → 200 монет
- `no_hints_25` - 25 задач без подсказок → 500 монет
- `no_hints_50` - 50 задач без подсказок → 1000 монет

### 5. `stickers_collected` - Собранные стикеры
**Триггер:** Получение нового стикера
**Метод проверки:** Будет реализовано в Phase 6
**Статус:** TODO

## Архитектура

### Store Layer (`store/achievement_progress.go`)

**`UpdateAchievementProgress(childProfileID, achievementID, newProgress)`**
- Обновляет прогресс конкретного достижения
- Автоматически разблокирует если `newProgress >= requirement_value`
- Возвращает `wasUnlocked bool` для отслеживания новых разблокировок

**`CheckAndUpdateAchievementsByType(childProfileID, requirementType, currentValue)`**
- Проверяет ВСЕ достижения заданного типа
- Обновляет прогресс для каждого
- Возвращает список ID разблокированных достижений

**Вспомогательные методы:**
- `GetCurrentStreakDays()` - текущий streak пользователя
- `GetVillainsDefeatedCount()` - количество побеждённых монстров
- `GetTasksCorrectCount()` - количество правильных задач
- `GetTasksNoHintsCount()` - количество задач без подсказок

### Service Layer (`service/achievement.go`)

**`AchievementService`** - централизованный сервис для проверки достижений

Методы:
- `CheckStreakAchievements(childProfileID)` - проверка достижений за streak
- `CheckVillainAchievements(childProfileID)` - проверка достижений за монстров
- `CheckTasksCorrectAchievements(childProfileID)` - проверка достижений за задачи
- `CheckTasksNoHintsAchievements(childProfileID)` - проверка достижений за задачи без подсказок

Каждый метод:
1. Получает текущее значение метрики (streak, villains count, tasks count)
2. Вызывает `CheckAndUpdateAchievementsByType()`
3. Логирует разблокированные достижения

### Интеграция

**ProfileService:**
```go
func (s *ProfileService) UpdateStreakAndActivity(ctx, childProfileID) error {
    // ... обновление streak_days ...

    if s.achievementService != nil && newStreak != currentStreak {
        s.achievementService.CheckStreakAchievements(ctx, childProfileID)
    }
}
```

**VillainService:**
```go
func (s *VillainService) DealDamageToVillain(ctx, childProfileID, attemptID, taskType) (defeated, coins, error) {
    // ... нанесение урона, проверка победы ...

    if defeated && s.achievementService != nil {
        s.achievementService.CheckVillainAchievements(ctx, childProfileID)
    }
}
```

**AttemptService:**
```go
func (s *AttemptService) ProcessCheck(ctx, attemptID, childProfileID, ...) error {
    // ... проверка решения ...

    if checkResp.Decision == types.CheckDecisionCorrect && s.achievementService != nil {
        s.achievementService.CheckTasksCorrectAchievements(ctx, childProfileID)
        s.achievementService.CheckTasksNoHintsAchievements(ctx, childProfileID)
    }
}
```

## Database Schema

### `achievements` table
```sql
CREATE TABLE achievements (
    id VARCHAR(100) PRIMARY KEY,
    type VARCHAR(50),           -- 'streak', 'villain', 'tasks', 'perfect'
    title VARCHAR(200),
    description TEXT,
    icon VARCHAR(50),           -- emoji
    requirement_type VARCHAR(50), -- 'streak_days', 'villains_defeated', 'tasks_correct', 'tasks_no_hints'
    requirement_value INTEGER,   -- число для достижения
    reward_type VARCHAR(50),     -- 'coins', 'sticker', 'avatar', 'badge'
    reward_name VARCHAR(200),
    reward_amount INTEGER,
    priority INTEGER DEFAULT 100 -- для сортировки (меньше = выше)
);
```

### `child_achievements` table
```sql
CREATE TABLE child_achievements (
    id BIGSERIAL PRIMARY KEY,
    child_profile_id UUID REFERENCES child_profiles(id),
    achievement_id VARCHAR(100) REFERENCES achievements(id),
    current_progress INTEGER DEFAULT 0,
    is_unlocked BOOLEAN DEFAULT FALSE,
    is_claimed BOOLEAN DEFAULT FALSE,
    unlocked_at TIMESTAMPTZ,
    claimed_at TIMESTAMPTZ
);
```

## Миграции

- **029** - Создание таблиц achievements, child_achievements
- **037** - Удаление shelves, добавление priority
- **038** - Автоматическое начисление наград (триггеры)
- **042** - Примеры достижений всех типов

## Логирование

Все разблокировки логируются с emoji 🎉:
```
[Store] 🎉 Achievement unlocked! child=uuid, achievement=streak_7, type=streak_days, value=7
[AchievementService] 🎉 Unlocked 1 streak achievements for child uuid: [streak_7]
```

## Тестирование

1. **Streak achievements:** Зайти на главную несколько дней подряд
2. **Villain achievements:** Победить монстров (решать задачи правильно)
3. **Tasks correct:** Решить задачи правильно через ProcessCheck
4. **Tasks no hints:** Решить задачи НЕ используя подсказки (hints_used=0)

## Roadmap

- [ ] Phase 6: Добавить достижения за стикеры (`stickers_collected`)
- [ ] Phase 7: Уведомления о разблокировке достижений в UI
- [ ] Phase 8: Анимации разблокировки достижений
- [ ] Phase 9: Социальные достижения (реферальная система)
