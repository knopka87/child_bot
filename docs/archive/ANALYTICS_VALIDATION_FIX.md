# Исправление ошибки валидации аналитики

## Дата: 2026-04-04

## Проблема

При открытии страницы онбординга в консоли появлялась ошибка:

```
[Analytics] Invalid event "onboarding_opened"
{
  errors: [
    "Missing required parameter: platform_type",
    "Missing required parameter: session_id"
  ]
}
```

## Причина

### Как работает аналитика

1. **Пользователь вызывает событие:**
   ```typescript
   analytics.trackEvent('onboarding_opened', {});
   ```

2. **AnalyticsService валидирует параметры:**
   ```typescript
   const validationResult = this.validator.validate(name, params);
   // params = {} - пустой объект
   // Валидатор проверяет: есть ли platform_type? НЕТ ❌
   // Валидатор проверяет: есть ли session_id? НЕТ ❌
   ```

3. **ТОЛЬКО ПОТОМ добавляются автоматические параметры:**
   ```typescript
   const event = {
     name,
     sessionId: this.sessionManager.getSessionId(), // Добавляется здесь
     params: {
       ...params,
       platform_type: this.sessionManager.getPlatform(), // Добавляется здесь
     },
   };
   ```

### Проблема в порядке операций

```
┌─────────────────────────────────────────┐
│ 1. Пользователь: trackEvent({})        │
│                                         │
│ 2. ❌ Валидация: platform_type? НЕТ    │
│    ❌ Валидация: session_id? НЕТ       │
│                                         │
│ 3. ⚠️ Показать warning в консоли       │
│                                         │
│ 4. ✅ Добавить platform_type            │
│    ✅ Добавить session_id               │
│                                         │
│ 5. Отправить событие                   │
└─────────────────────────────────────────┘
```

**Вывод:** Валидация происходит ДО добавления автоматических параметров.

## Решение

Параметры, которые **добавляются автоматически системой**, не должны требоваться в схеме валидации.

### Изменения в schema.ts

**Было:**
```typescript
onboarding_opened: {
  required: ['platform_type', 'session_id'], // ❌ Требуются от пользователя
  optional: ['entry_point'],
  params: {
    platform_type: 'string',
    session_id: 'string',
    entry_point: 'string',
  },
}
```

**Стало:**
```typescript
onboarding_opened: {
  required: [], // ✅ platform_type и session_id добавляются автоматически
  optional: ['entry_point'],
  params: {
    platform_type: 'string', // Описаны в params для документации
    session_id: 'string',    // но не required
    entry_point: 'string',
  },
}
```

## Логика системы

### Автоматически добавляемые параметры

AnalyticsService **автоматически добавляет** следующие параметры ко ВСЕМ событиям:

| Параметр | Источник | Когда добавляется |
|----------|----------|-------------------|
| `session_id` | `sessionManager.getSessionId()` | В поле event.sessionId |
| `platform_type` | `sessionManager.getPlatform()` | В event.params |
| `app_version` | `import.meta.env.VITE_APP_VERSION` | В event.params |

**Код из AnalyticsService.ts:**
```typescript
const event: StoredAnalyticsEvent = {
  name,
  timestamp: Date.now(),
  sessionId: this.sessionManager.getSessionId(), // Автоматически
  params: {
    ...params, // Пользовательские параметры
    app_version: import.meta.env.VITE_APP_VERSION || '0.1.0', // Автоматически
    platform_type: this.sessionManager.getPlatform(), // Автоматически
  },
};
```

### Какие параметры требовать в схеме?

**Правило:**
- ✅ **required** - только то, что **ДОЛЖЕН передать пользователь**
- ❌ **НЕ required** - то, что **добавляется автоматически**

**Примеры:**

```typescript
// ✅ ПРАВИЛЬНО: пользователь должен передать grade
grade_selected: {
  required: ['grade'], // Пользователь передаёт
  params: {
    grade: 'number',
    child_profile_id: 'string', // Может добавиться автоматически
  },
}

// ✅ ПРАВИЛЬНО: пользователь должен передать child_profile_id
home_opened: {
  required: ['child_profile_id'], // Пользователь передаёт
  params: {
    child_profile_id: 'string',
    entry_point: 'string', // optional
  },
}

// ✅ ПРАВИЛЬНО: ничего не требуется от пользователя
onboarding_opened: {
  required: [], // Всё добавляется автоматически
  optional: ['entry_point'],
  params: {
    platform_type: 'string', // Автоматически
    session_id: 'string',    // Автоматически
    entry_point: 'string',   // Опционально от пользователя
  },
}
```

## Проверка исправления

### Тест 1: Консоль браузера

1. Откройте http://localhost:5173
2. Откройте DevTools → Console
3. ✅ **Ожидается:** НЕТ ошибок `[Analytics] Invalid event "onboarding_opened"`

### Тест 2: Отправка события

Проверьте в Network tab (DevTools):

```json
POST /analytics/events
{
  "events": [
    {
      "name": "onboarding_opened",
      "timestamp": 1234567890,
      "sessionId": "abc-def-123", // ✅ Добавлено автоматически
      "params": {
        "platform_type": "web",   // ✅ Добавлено автоматически
        "app_version": "0.1.0"    // ✅ Добавлено автоматически
      }
    }
  ]
}
```

### Тест 3: Backend логи

```bash
docker compose -f docker/docker-compose.dev.yml logs backend | grep analytics
```

Должны видеть события с корректными параметрами.

## Изменённые файлы

1. **frontend/src/services/analytics/schema.ts**
   - Убрал `platform_type` и `session_id` из `required` для `onboarding_opened`
   - Добавил комментарий объясняющий почему

## Best Practices

### При создании новых событий

**Спрашивайте себя:**

1. **Этот параметр добавляется автоматически?**
   - Если ДА → НЕ включать в `required`
   - Если НЕТ → включить в `required` (если обязателен)

2. **Пользователь должен передавать этот параметр?**
   - Если ДА → включить в `required`
   - Если НЕТ → `optional` или вообще не требовать

**Автоматические параметры (НЕ требовать):**
- `platform_type` - добавляется из SessionManager
- `session_id` - добавляется из SessionManager
- `app_version` - добавляется из env

**Пользовательские параметры (требовать если обязательны):**
- `child_profile_id` - если событие для профиля
- `grade` - если событие о выборе класса
- `villain_id` - если событие о злодее
- и т.д.

## Альтернативные решения (не реализованы)

### Вариант 1: Валидация после добавления параметров

Изменить порядок в AnalyticsService:

```typescript
// Сначала создать event
const event = {
  name,
  sessionId: this.sessionManager.getSessionId(),
  params: {
    ...params,
    platform_type: this.sessionManager.getPlatform(),
  },
};

// Потом валидировать
const validationResult = this.validator.validate(name, event.params);
```

**Минусы:**
- Сложнее логика
- Нужно передавать sessionId отдельно

### Вариант 2: Умный валидатор

Передавать context в валидатор:

```typescript
this.validator.validate(name, params, {
  sessionId: this.sessionManager.getSessionId(),
  platformType: this.sessionManager.getPlatform(),
});
```

**Минусы:**
- Усложняет API валидатора
- Требует изменений в нескольких местах

### Выбранное решение: Правильная схема

✅ **Самое простое и правильное:**
- Схема отражает реальность
- Не требуем то, что добавляется автоматически
- Минимальные изменения

## Заключение

✅ **Исправлено:** Убрана ошибка валидации для `onboarding_opened`

✅ **Принцип:** Автоматические параметры не должны быть в `required`

✅ **Результат:** События отправляются без warnings в консоли

---

**Статус:** ✅ Готово
**Протестировано:** Да
**Breaking changes:** Нет
