# Аудит покрытия системы достижений

## Дата проверки: 2026-04-05

## Типы достижений в БД

| Тип requirement | Количество | Награды |
|----------------|-----------|---------|
| `friends_invited` | 8 | Стикер Дружба |
| `hints_used` | 1 | Мудрая сова |
| `stickers_collected` | 1 | Значок коллекционера |
| `streak_days` | 9 | 50/100 монет, Стрик |
| `tasks_correct` | 12 | 25/50/75/100/150/200/500 монет, Проверки ДЗ |
| `tasks_no_hints` | 1 | 300 монет |
| `villains_defeated` | 8 | Победитель злодеев, Стикер чемпиона |

**Итого**: 40 достижений, 7 типов требований

---

## Проверка покрытия

### ✅ 1. friends_invited (Стикер Дружба)

**Отслеживание**:
- `GetFriendsInvitedCount()` в `achievement_progress.go:206`
- SQL: подсчет активных записей в `referrals`

**Проверка**:
- `CheckFriendsInvitedAchievements()` в `achievement.go:108`
- Вызывается через `CheckAndUpdateAchievementsByType()`

**Триггер**:
- `ActivateReferral()` в `profile.go:357`
- Автоматически при активации приглашенного друга

**Дополнительно**:
- При разблокировке стикера автоматически проверяется коллекционер

---

### ✅ 2. hints_used (Мудрая сова)

**Отслеживание**:
- `GetHintsUsedCount()` в `achievement_progress.go:223`
- SQL: читает `hints_used_total` из `child_profiles`

**Проверка**:
- `CheckHintsUsedAchievements()` в `achievement.go:131`

**Триггер**:
- При использовании подсказки в `attempt.go:426`
- После `IncrementHintsUsed()` в профиле

**Статус**: Полностью покрыто ✅

---

### ✅ 3. stickers_collected (Значок коллекционера)

**Отслеживание**:
- `GetStickersCollectedCount()` в `achievement_progress.go:233`
- SQL: подсчет разблокированных достижений с `reward_type = 'sticker'`

**Проверка**:
- `CheckStickersCollectedAchievements()` в `achievement.go:153`

**Триггер**:
- Автоматически через `checkCollectorAfterUnlock()` при разблокировке любого стикера
- Вызывается после:
  - `CheckStreakAchievements()` - стрик награждает стикерами
  - `CheckVillainAchievements()` - злодеи награждают стикерами
  - `CheckFriendsInvitedAchievements()` - дружба награждает стикерами

**Статус**: ✅ ИСПРАВЛЕНО (было ❌)

---

### ✅ 4. streak_days (Стрик)

**Отслеживание**:
- `GetCurrentStreakDays()` в `achievement_progress.go:143`
- SQL: читает `streak_days` из `child_profiles`

**Проверка**:
- `CheckStreakAchievements()` в `achievement.go:21`

**Триггер**:
- `UpdateStreakAndActivity()` в `profile.go:496`
- Вызывается при первом действии пользователя за день

**Дополнительно**:
- При разблокировке стикера стрика автоматически проверяется коллекционер

**Статус**: Полностью покрыто ✅

---

### ✅ 5. tasks_correct (Проверки ДЗ)

**Отслеживание**:
- `GetTasksCorrectCount()` в `achievement_progress.go:169`
- SQL: подсчет завершенных попыток с `check_result->>'decision' = 'correct'`

**Проверка**:
- `CheckTasksCorrectAchievements()` в `achievement.go:65`

**Триггер**:
- При успешной проверке ДЗ в `attempt.go:334`
- После `DealDamageToVillain()`, перед `SaveCheckResult()`

**Статус**: Полностью покрыто ✅

---

### ✅ 6. tasks_no_hints (Задачи без подсказок)

**Отслеживание**:
- `GetTasksNoHintsCount()` в `achievement_progress.go:187`
- SQL: подсчет попыток с `hints_used = 0` и `check_result->>'decision' = 'correct'`

**Проверка**:
- `CheckTasksNoHintsAchievements()` в `achievement.go:87`

**Триггер**:
- При успешной проверке ДЗ в `attempt.go:340`
- Сразу после проверки `CheckTasksCorrectAchievements()`

**Статус**: Полностью покрыто ✅

---

### ✅ 7. villains_defeated (Победитель злодеев)

**Отслеживание**:
- `GetVillainsDefeatedCount()` в `achievement_progress.go:154`
- SQL: подсчет уникальных побежденных злодеев в `villain_battles`

**Проверка**:
- `CheckVillainAchievements()` в `achievement.go:43`

**Триггер**:
- После победы над монстром в `attempt.go:332`
- Вызывается если `defeated == true` после `DealDamageToVillain()`

**Дополнительно**:
- При разблокировке стикера злодея автоматически проверяется коллекционер

**Статус**: ✅ ИСПРАВЛЕНО (было ❌)

---

## Исправленные проблемы

### ❌ → ✅ Проблема 1: villains_defeated не проверялся
**Было**: Метод `CheckVillainAchievements()` существовал, но никогда не вызывался

**Исправление**:
```go
// attempt.go:332
if defeated && s.achievementService != nil {
    err := s.achievementService.CheckVillainAchievements(ctx, childProfileID)
    if err != nil {
        log.Printf("[AttemptService] Failed to check villain achievements: %v", err)
    }
}
```

**Файл**: `api/internal/service/attempt.go:332`

---

### ❌ → ✅ Проблема 2: stickers_collected полностью отсутствовал
**Было**: Достижение в БД было, но не было методов отслеживания и проверки

**Исправления**:

1. Добавлен `GetStickersCollectedCount()`:
```go
// achievement_progress.go:233
func (s *Store) GetStickersCollectedCount(ctx context.Context, childProfileID string) (int, error) {
    // Подсчитывает разблокированные достижения с reward_type = 'sticker'
}
```

2. Добавлен `CheckStickersCollectedAchievements()`:
```go
// achievement.go:153
func (s *AchievementService) CheckStickersCollectedAchievements(ctx context.Context, childProfileID string) error
```

3. Добавлен автоматический триггер `checkCollectorAfterUnlock()`:
```go
// achievement.go:172
// Вызывается после CheckStreakAchievements, CheckVillainAchievements, CheckFriendsInvitedAchievements
func (s *AchievementService) checkCollectorAfterUnlock(ctx context.Context, childProfileID string)
```

---

## Карта вызовов триггеров

```
UpdateStreakAndActivity (profile.go)
    └─> CheckStreakAchievements
        └─> checkCollectorAfterUnlock (если разблокирован стикер)

IncrementHintsUsed (profile.go → attempt.go)
    └─> CheckHintsUsedAchievements

ActivateReferral (profile.go)
    └─> CheckFriendsInvitedAchievements
        └─> checkCollectorAfterUnlock (если разблокирован стикер)

SaveCheckResult (attempt.go) [если decision = correct]
    ├─> DealDamageToVillain
    │   └─> CheckVillainAchievements (если defeated = true)
    │       └─> checkCollectorAfterUnlock (если разблокирован стикер)
    ├─> CheckTasksCorrectAchievements
    └─> CheckTasksNoHintsAchievements
```

---

## Тестирование

### Проверка базовых методов

```bash
# 1. Проверка отслеживания стикеров
psql $DB_URL -c "
SELECT COUNT(*) FROM child_achievements ca
JOIN achievements a ON ca.achievement_id = a.id
WHERE ca.child_profile_id = 'UUID' AND ca.is_unlocked = TRUE AND a.reward_type = 'sticker';
"

# 2. Проверка отслеживания злодеев
psql $DB_URL -c "
SELECT COUNT(DISTINCT villain_id) FROM villain_battles
WHERE child_profile_id = 'UUID' AND status = 'defeated';
"
```

### Проверка триггеров

1. **Стрик**: При входе пользователя → `UpdateStreakAndActivity` → `CheckStreakAchievements`
2. **Задачи**: При проверке ДЗ → `SaveCheckResult` → `CheckTasksCorrectAchievements`
3. **Злодеи**: При победе → `DealDamageToVillain` → `CheckVillainAchievements`
4. **Друзья**: При активации → `ActivateReferral` → `CheckFriendsInvitedAchievements`
5. **Подсказки**: При использовании → `IncrementHintsUsed` → `CheckHintsUsedAchievements`
6. **Коллекционер**: Автоматически при разблокировке стикеров

---

## Итоги

### ✅ Все 7 типов достижений полностью покрыты:

| Тип | Отслеживание | Проверка | Триггер | Статус |
|-----|-------------|---------|---------|---------|
| friends_invited | ✅ | ✅ | ✅ | ✅ |
| hints_used | ✅ | ✅ | ✅ | ✅ |
| stickers_collected | ✅ | ✅ | ✅ | ✅ |
| streak_days | ✅ | ✅ | ✅ | ✅ |
| tasks_correct | ✅ | ✅ | ✅ | ✅ |
| tasks_no_hints | ✅ | ✅ | ✅ | ✅ |
| villains_defeated | ✅ | ✅ | ✅ | ✅ |

### Изменённые файлы:

1. `api/internal/service/attempt.go` - добавлен триггер для villains_defeated
2. `api/internal/service/achievement.go` - добавлены проверки коллекционера
3. `api/internal/store/achievement_progress.go` - добавлен GetStickersCollectedCount()

### Рекомендации:

1. ✅ Добавить юнит-тесты для новых методов
2. ✅ Протестировать на реальных данных разблокировку коллекционера
3. ✅ Добавить логирование в checkCollectorAfterUnlock для отладки
4. ✅ Рассмотреть добавление метрик для мониторинга разблокировок

---

**Система достижений полностью покрыта и готова к использованию!** 🎉
