# REST API Endpoints

## Статус: Phase 3 завершена ✅

Реализовано **34 endpoints** с полной валидацией и placeholder responses.

---

## Authentication

Большинство endpoints требуют следующие заголовки:
```
X-Platform-ID: vk|telegram|max|web
X-Child-Profile-ID: <uuid профиля ребенка>
```

**Public endpoints** (без auth):
- `GET /health`
- `POST /onboarding/start` (будет добавлен позже)
- `POST /onboarding/complete` (будет добавлен позже)

---

## Endpoints

### Health Check

#### `GET /health`
Проверка состояния сервера

**Response:**
```json
{
  "status": "ok"
}
```

---

### Attempts (Попытки решения)

#### `POST /attempts`
Создать новую попытку

**Request:**
```json
{
  "type": "help|check",
  "child_profile_id": "uuid"
}
```

**Response:**
```json
{
  "attempt_id": "uuid",
  "status": "created"
}
```

#### `GET /attempts/unfinished`
Получить незавершенную попытку

**Query params:**
- `childProfileId` - UUID профиля

**Response:** `Attempt` или `null`

#### `GET /attempts/recent`
Получить последние попытки

**Query params:**
- `childProfileId` - UUID профиля
- `limit` - количество (default: 3)

**Response:** Array of `RecentAttempt`

#### `POST /attempts/{id}/images`
Загрузить изображение

**Request:**
```json
{
  "image_type": "task|answer",
  "image_data": "data:image/png;base64,..."
}
```

#### `POST /attempts/{id}/process`
Начать обработку через LLM

**Response:**
```json
{
  "status": "processing",
  "message": "Attempt is being processed"
}
```

#### `GET /attempts/{id}/result`
Получить результат

**Response:**
```json
{
  "attempt_id": "uuid",
  "type": "help|check",
  "status": "completed",
  "result": {...}
}
```

#### `POST /attempts/{id}/next-hint`
Получить следующую подсказку

**Response:**
```json
{
  "hint": "text",
  "hint_index": 1,
  "total_hints": 3,
  "has_more_hints": true
}
```

#### `DELETE /attempts/{id}`
Удалить попытку

**Response:** `204 No Content`

---

### Home

#### `GET /home/{childProfileId}`
Получить данные для главного экрана

**Response:**
```json
{
  "profile": {
    "id": "uuid",
    "display_name": "string",
    "level": 5,
    "level_progress": 60,
    "coins_balance": 150,
    "tasks_solved_correct_count": 42
  },
  "mascot": {
    "id": "string",
    "state": "idle|happy|thinking|celebrating",
    "image_url": "string",
    "message": "string"
  },
  "villain": {...},
  "unfinished_attempt": {...},
  "recent_attempts": [...]
}
```

---

### Profile

#### `GET /profile`
Получить профиль

**Response:**
```json
{
  "id": "uuid",
  "display_name": "string",
  "avatar_id": "string",
  "avatar_url": "string",
  "grade": 5,
  "subscription": {
    "status": "trial|active|expired|cancelled",
    "trial_days_remaining": 7
  }
}
```

#### `PUT /profile`
Обновить профиль

**Request:**
```json
{
  "display_name": "string",
  "avatar_id": "string",
  "grade": 5
}
```

#### `GET /profile/history`
Получить историю попыток

**Query params:**
- `mode` - help|check|all
- `status` - success|error|in_progress|all
- `date_from` - ISO date
- `date_to` - ISO date

**Response:** Array of `HistoryAttempt`

#### `GET /profile/stats`
Получить статистику

**Response:**
```json
{
  "total_attempts": 50,
  "successful_attempts": 42,
  "errors_fixed": 35,
  "streak_days": 7,
  "average_accuracy": 84.0,
  "total_hints_used": 15
}
```

---

### Achievements

#### `GET /achievements`
Список всех достижений

**Query params:**
- `category` - тип достижения

**Response:** Array of `Achievement`

#### `GET /achievements/unlocked`
Только разблокированные

**Response:** Array of `Achievement`

#### `GET /achievements/stats`
Статистика достижений

**Response:**
```json
{
  "unlocked_count": 5,
  "total_count": 25,
  "progress_percent": 20.0
}
```

#### `GET /achievements/{id}`
Информация о достижении

**Response:** `Achievement`

#### `POST /achievements/{id}/claim`
Забрать награду

**Response:**
```json
{
  "claimed": true,
  "reward_type": "coins",
  "amount": 50,
  "message": "Получено 50 монет!"
}
```

---

### Villains

#### `GET /villains`
Список злодеев

**Response:** Array of `Villain`

#### `GET /villains/active`
Активный злодей

**Response:** `Villain`

#### `GET /villains/{id}`
Информация о злодее

**Response:** `Villain`

#### `GET /villains/{id}/battle`
Информация о битве

**Response:**
```json
{
  "villain_id": "string",
  "battle_stats": {
    "total_damage_dealt": 25,
    "correct_tasks_count": 5,
    "damage_per_task": 5,
    "progress_percent": 25.0
  },
  "recent_damage": [...],
  "can_damage_now": true
}
```

#### `GET /villains/{id}/victory`
Информация о победе

**Response:**
```json
{
  "villain_id": "string",
  "villain_name": "string",
  "defeated_at": "ISO date",
  "total_damage": 100,
  "tasks_completed": 20,
  "rewards": [...]
}
```

#### `POST /villains/{id}/damage`
Нанести урон

**Request:**
```json
{
  "attempt_id": "uuid",
  "damage": 5
}
```

**Response:**
```json
{
  "damage_dealt": 5,
  "villain_hp": 70,
  "is_defeated": false,
  "message": "Нанесено 5 урона!"
}
```

---

### Subscription

#### `GET /subscription/status`
Статус подписки

**Response:**
```json
{
  "status": "trial|active|expired|cancelled",
  "features": ["unlimited_tasks", "hints"],
  "trial_days_remaining": 7,
  "can_cancel": false,
  "can_resume": false
}
```

#### `GET /subscription/plans`
Доступные планы

**Response:** Array of `SubscriptionPlan`

#### `POST /subscription/subscribe`
Оформить подписку

**Request:**
```json
{
  "plan_id": "string",
  "payment_method": "card|yookassa"
}
```

**Response:**
```json
{
  "payment_url": "string",
  "payment_id": "string",
  "status": "pending",
  "expires_at": "ISO date"
}
```

#### `DELETE /subscription/cancel`
Отменить подписку

**Response:**
```json
{
  "status": "cancelled",
  "cancelled_at": "ISO date",
  "expires_at": "ISO date",
  "message": "string"
}
```

#### `POST /subscription/resume`
Возобновить подписку

**Response:**
```json
{
  "status": "active",
  "renews_at": "ISO date",
  "message": "Подписка возобновлена"
}
```

---

### Friends / Referral

#### `GET /friends`
Список друзей

**Response:** Array of `InvitedFriend`

#### `POST /friends/invite`
Создать приглашение

**Response:**
```json
{
  "referral_link": "string",
  "share_text": "string",
  "platform": "link"
}
```

#### `GET /friends/referrals`
Данные реферальной программы

**Response:**
```json
{
  "referral_code": "string",
  "referral_link": "string",
  "total_invited": 3,
  "active_invited": 2,
  "total_rewards": 150,
  "invited_friends": [...],
  "reward_milestones": [...]
}
```

#### `GET /friends/leaderboard`
Leaderboard друзей

**Response:** Array of `LeaderboardEntry`

---

## Error Responses

Все ошибки возвращаются в формате:
```json
{
  "error": "Error message",
  "code": "ERROR_CODE"
}
```

**HTTP статусы:**
- `400` - Bad Request (валидация не прошла)
- `401` - Unauthorized (отсутствуют auth заголовки)
- `403` - Forbidden (нет доступа)
- `404` - Not Found (ресурс не найден)
- `409` - Conflict (конфликт данных)
- `500` - Internal Server Error

---

## TODO Phase 4: Service Layer

Все endpoints возвращают placeholder данные. В Phase 4 будет реализована бизнес-логика:
- Интеграция с Store (PostgreSQL)
- Интеграция с LLMClient (AI обработка)
- Реальная логика для всех handlers
