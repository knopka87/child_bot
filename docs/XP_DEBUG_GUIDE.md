# Руководство по отладке системы XP

## Обзор системы

Система XP начисляет опыт за различные действия пользователя:

| Событие | XP | Где вызывается | Файл:строка |
|---------|-----|----------------|-------------|
| Правильный ответ | 50 | При проверке задания (is_correct=true) | `api/internal/service/attempt.go:431` |
| Исправление ошибок | 20 | При проверке задания (is_correct=false) | `api/internal/service/attempt.go:457` |
| Запрос подсказки | 10 | При каждой подсказке | `api/internal/service/attempt.go:569` |
| Ежедневный вход | 30 | При первом входе в день | `api/internal/service/profile.go:640` |
| Победа над злодеем | 100 | При нанесении финального урона | `api/internal/service/attempt.go:406` |
| Разблокировка достижения | 50 | При получении любого достижения | `api/internal/service/achievement.go:49` |

**Бонус:** При повышении уровня дополнительно начисляется **100 монет**.

## Логирование

### Добавленные логи для отладки

Теперь каждый Award* метод выводит подробные логи:

1. **🎯 Starting** - Начало начисления XP
2. **❌ Failed** - Ошибка при начислении
3. **✅ Successfully** - Успешное начисление (без повышения уровня)
4. **🎉 Level up** - Повышение уровня
5. **⚠️ NIL** - profileService не установлен

### Пример успешного начисления XP

```
[AttemptService] 🎯 Answer is correct, calling AwardCorrectAnswer for child abc123
[ProfileService.AwardCorrectAnswer] 🎯 Starting to award 50 XP for correct answer to child abc123
[Store.AddXP] 🎯 Starting to add 50 XP to child abc123
[Store.AddXP] 📊 Current state: child=abc123, XP=0, level=1
[Store] ✅ XP updated: child=abc123, XP: 0 -> 50, level: 1
[ProfileService.AwardCorrectAnswer] ✅ Successfully awarded 50 XP for correct answer to child abc123 (no level up)
[AttemptService] ✅ Successfully called AwardCorrectAnswer for child abc123
```

### Пример повышения уровня

```
[Store.AddXP] 🎯 Starting to add 50 XP to child abc123
[Store.AddXP] 📊 Current state: child=abc123, XP=80, level=1
[Store] 🎉 Level up! child=abc123, level 1 -> 2 (XP: 130)
[Store] 🪙 Level up reward: child=abc123, level=2, coins=100
[Store] ✅ XP updated: child=abc123, XP: 80 -> 130, level: 1 -> 2, coins_reward=100
[ProfileService] 🎉 Level up from correct answer! child=abc123, level=2
```

### Пример ошибки

```
[AttemptService] 🎯 Answer is correct, calling AwardCorrectAnswer for child abc123
[ProfileService.AwardCorrectAnswer] 🎯 Starting to award 50 XP for correct answer to child abc123
[Store.AddXP] 🎯 Starting to add 50 XP to child abc123
[Store.AddXP] ❌ Failed to get current XP and level for abc123: sql: no rows in result set
[ProfileService] ❌ Failed to award correct answer XP for abc123: get current xp and level: sql: no rows in result set
[AttemptService] ❌ Failed to award correct answer XP for abc123: get current xp and level: sql: no rows in result set
```

### Пример NIL profileService

```
[AttemptService] ⚠️ profileService is NIL - cannot award correct answer XP
```

## Как проверить начисление XP

### 1. Проверка логов сервера

```bash
# Для Docker
docker logs child-bot-api-1 | grep -E "🎯|❌|✅|🎉|⚠️" | tail -n 50

# Для локального запуска
tail -f logs/app.log | grep -E "🎯|❌|✅|🎉|⚠️"
```

### 2. Проверка в базе данных

```sql
-- Проверить текущий XP и уровень
SELECT id, name, xp_total, level, coins_balance
FROM child_profiles
WHERE id = 'your-profile-id';

-- Посмотреть историю изменений (если есть таблица логов)
SELECT * FROM child_profiles
WHERE id = 'your-profile-id'
ORDER BY updated_at DESC
LIMIT 10;
```

### 3. Проверка через API

```bash
# Получить профиль
curl -X GET "http://localhost:8080/api/profile" \
  -H "X-Child-Profile-ID: your-profile-id"
```

## Типичные проблемы и решения

### Проблема 1: XP не начисляется вообще

**Симптомы:**
- В логах нет сообщений 🎯 Starting
- В логах есть ⚠️ profileService is NIL

**Решение:**
1. Проверьте `api/internal/api/router/router.go:48` - должно быть:
   ```go
   achievementService.SetProfileService(profileService)
   ```
2. Проверьте что backend перезапущен после изменений
3. Убедитесь что сервисы инициализируются в правильном порядке

### Проблема 2: XP начисляется, но не сохраняется

**Симптомы:**
- В логах есть ✅ Successfully awarded
- В базе данных XP не изменился

**Решение:**
1. Проверьте что транзакция коммитится в `store/xp.go:85`
2. Проверьте права доступа к базе данных
3. Проверьте что `child_profile_id` корректный

### Проблема 3: Ошибка при получении XP из БД

**Симптомы:**
- В логах есть ❌ Failed to get current XP and level
- Ошибка: `sql: no rows in result set`

**Решение:**
1. Проверьте что профиль существует в таблице `child_profiles`
2. Проверьте что `child_profile_id` корректный
3. Создайте профиль если его нет:
   ```sql
   INSERT INTO child_profiles (id, name, xp_total, level)
   VALUES ('your-id', 'Test User', 0, 1);
   ```

### Проблема 4: Level up не происходит

**Симптомы:**
- XP начисляется
- XP достаточно для level up
- Но уровень не повышается

**Решение:**
1. Проверьте формулу расчета XP в `store/xp.go:113`:
   ```go
   XPForLevel(level) = 50 × level² + 50 × level
   ```
2. Для уровня 1→2 нужно: 50×1²+50×1 = 100 XP
3. Проверьте что цикл level up работает корректно (строки 41-52)

### Проблема 5: Ошибка при начислении монет за level up

**Симптомы:**
- Level up произошёл
- Но ошибка: `add level up coins: ...`

**Решение:**
1. Проверьте что поле `coins_balance` существует в таблице
2. Проверьте что поле имеет тип INTEGER или BIGINT
3. Проверьте права на UPDATE для таблицы

## Тестирование системы XP

### Тест 1: Правильный ответ (+50 XP)

1. Создать попытку
2. Загрузить изображение задания и решения
3. Проверить задание (правильное решение)
4. Проверить логи - должно быть:
   ```
   [AttemptService] 🎯 Answer is correct, calling AwardCorrectAnswer
   [ProfileService.AwardCorrectAnswer] ✅ Successfully awarded 50 XP
   ```

### Тест 2: Неправильный ответ (+20 XP)

1. Создать попытку
2. Загрузить изображение с ошибками
3. Проверить задание (неправильное решение)
4. Проверить логи - должно быть:
   ```
   [AttemptService] 🎯 Answer is incorrect (errors found), calling AwardFixErrors
   [ProfileService.AwardFixErrors] ✅ Successfully awarded 20 XP
   ```

### Тест 3: Запрос подсказки (+10 XP)

1. Создать попытку с заданием
2. Запросить подсказку через `/api/attempts/{id}/next-hint`
3. Проверить логи - должно быть:
   ```
   [AttemptService.GetNextHint] IncrementHintsUsed succeeded, now calling AwardHintRequest
   [ProfileService.AwardHintRequest] ✅ Successfully awarded 10 XP
   ```

### Тест 4: Ежедневный вход (+30 XP)

1. Войти в приложение первый раз за день
2. Вызвать `/api/profile` или `/api/home/{childProfileId}`
3. Проверить логи - должно быть:
   ```
   [UpdateStreakAndActivity] 🎯 New day detected, calling AwardDailyLogin
   [ProfileService.AwardDailyLogin] ✅ Successfully awarded 30 XP
   ```

### Тест 5: Победа над злодеем (+100 XP)

1. Получить активного злодея
2. Решить достаточно заданий чтобы нанести финальный урон
3. Проверить логи - должно быть:
   ```
   [AttemptService] 🎯 Villain defeated! Calling AwardVillainDefeat
   [ProfileService.AwardVillainDefeat] ✅ Successfully awarded 100 XP
   ```

### Тест 6: Разблокировка достижения (+50 XP)

1. Выполнить условие для достижения (например, 5 дней подряд)
2. Проверить логи - должно быть:
   ```
   [AchievementService] 🎯 profileService is available, awarding XP for 1 achievements
   [ProfileService.AwardAchievementUnlock] ✅ Successfully awarded 50 XP
   ```

## Формула расчета уровней

```
XP для уровня N = 50 × N² + 50 × N

Уровень 1→2: 100 XP
Уровень 2→3: 250 XP
Уровень 3→4: 500 XP
Уровень 4→5: 850 XP
Уровень 5→6: 1300 XP
```

## Контакты для поддержки

Если проблема не решается:
1. Соберите логи за последние 30 минут
2. Укажите `child_profile_id` у которого проблема
3. Укажите какое действие выполнялось
4. Приложите скриншоты логов с эмодзи маркерами
