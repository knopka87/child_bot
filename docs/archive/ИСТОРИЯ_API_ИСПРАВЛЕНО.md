# История - API эндпоинт исправлен

**Дата:** 2026-04-04

## Проблема

Фронтенд запрашивал неправильный эндпоинт:

```
GET /attempts/history?child_profile_id=xxx&limit=50&offset=0&filter=all
```

**Ошибка:** 405 Method Not Allowed

## Причина

Бэкенд реализовал эндпоинт по другому пути:

```
GET /profile/history
```

**Отличия:**
1. Путь: `/profile/history` вместо `/attempts/history`
2. Параметр `child_profile_id` берётся из middleware автоматически (не нужен в query)
3. Формат данных: snake_case вместо camelCase

## ✅ Исправление

### 1. Обновлён метод `getHistory()` в `api/profile.ts`

**Было:**
```typescript
const response = await apiClient.get<BackendHistoryResponse>(
  '/attempts/history',
  { params: { child_profile_id: childProfileId, ... } }
);
```

**Стало:**
```typescript
const response = await apiClient.get<BackendHistoryAttempt[]>(
  '/profile/history',
  { params: { mode: '...', status: '...' } }
);
```

### 2. Добавлена конвертация snake_case → camelCase

Бэкенд возвращает:
```json
{
  "id": "uuid",
  "mode": "help",
  "status": "success",
  "scenario_type": "single_photo",
  "created_at": "2026-04-04T10:00:00Z",
  "completed_at": "2026-04-04T10:05:00Z",
  "images": [...],
  "result": {...},
  "hints_used": 1
}
```

Фронтенд получает:
```typescript
{
  id: "uuid",
  mode: "help",
  status: "success",
  scenarioType: "single_photo",
  createdAt: "2026-04-04T10:00:00Z",
  completedAt: "2026-04-04T10:05:00Z",
  images: [...],
  result: {...},
  hintsUsed: 1
}
```

### 3. Упрощён метод `getHistoryDetail()`

Теперь использует данные из `getHistory()` вместо отдельного эндпоинта.

## Query параметры

| Параметр | Тип | Описание |
|----------|-----|----------|
| `mode` | 'help' \| 'check' | Фильтр по режиму (опционально) |
| `status` | 'success' \| 'error' \| 'in_progress' | Фильтр по статусу (опционально) |
| `date_from` | string | Начальная дата (опционально) |
| `date_to` | string | Конечная дата (опционально) |

**Примечание:** Параметр `child_profile_id` НЕ нужен - бэкенд берёт его из JWT токена через middleware.

## Backend эндпоинт

**Файл:** `api/internal/api/handler/profile.go:228-250`

**Роут:** `GET /profile/history` (зарегистрирован в `router.go:102`)

**Структура ответа:**
```go
type HistoryAttempt struct {
	ID           string         `json:"id"`
	Mode         string         `json:"mode"`
	Status       string         `json:"status"`
	ScenarioType string         `json:"scenario_type,omitempty"`
	CreatedAt    string         `json:"created_at"`
	CompletedAt  string         `json:"completed_at,omitempty"`
	Images       []HistoryImage `json:"images"`
	Result       *HistoryResult `json:"result,omitempty"`
	HintsUsed    int            `json:"hints_used,omitempty"`
}
```

## ✅ Проверка

```bash
npm run typecheck
# ✅ Успешно
```

## Тестирование

1. Открыть страницу истории:
```
http://localhost:5173/profile/history
```

2. Проверить консоль:
```
[profileAPI] History response: [...]
```

3. Должны загрузиться попытки (если они есть в БД)

## Следующие шаги

**Backend TODO (Phase 4):**

В `profile.go:228-250` метод `GetHistory()` сейчас возвращает пустой массив:

```go
// TODO: Phase 4 - получение истории через service layer
// history, err := h.service.GetHistory(r.Context(), childProfileID, filters)

// Placeholder
history := []HistoryAttempt{}

response.OK(w, history)
```

Нужно реализовать:
1. Service method `ProfileService.GetHistory()`
2. Загрузка попыток из БД
3. Применение фильтров (mode, status, даты)
4. Заполнение всех полей (изображения, результаты, подсказки)

## Статус

✅ Фронтенд исправлен
✅ Типы валидны
✅ Эндпоинт правильный
⏳ Backend placeholder (нужна реализация Phase 4)

---

**Готово к тестированию** после реализации бэкенд логики!
