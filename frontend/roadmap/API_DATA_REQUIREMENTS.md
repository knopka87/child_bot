# API Data Requirements — Требования к данным с Backend

**Проект:** Объяснятель ДЗ MiniApp
**Дата:** 2026-03-29

---

## Обзор

Этот документ описывает все данные, которые фронтенд ожидает получить с бекенда, их типы, структуры и API endpoints.

---

## 🔐 Security First!

**КРИТИЧНО:** Все endpoints требуют аутентификации и НИКОГДА не возвращают данные других пользователей!

См. детали в **[SECURITY.md](./SECURITY.md)**

### Ключевые правила безопасности:

1. **JWT токен обязателен** для каждого запроса
2. **Backend сам определяет user_id** из JWT токена (НИКОГДА из query параметров!)
3. **Пользователь не может указать чужой ID** в параметрах
4. **Backend проверяет владельца** перед возвратом данных
5. **Email маскируются** на frontend (u***@example.com)

---

## 1. Child Profile (Профиль ребенка)

### Endpoint
```
GET /api/v1/profile/me
```

**⚠️ Security Note:**
- ❌ НЕ используй `/api/v1/child-profile/:id` - уязвимость!
- ✅ Используй `/api/v1/profile/me` - backend сам определит ID из JWT

### Response
```typescript
interface ChildProfile {
  child_profile_id: string
  parent_user_id: string
  display_name: string
  avatar_id: string
  grade: number // 1-11
  level: number // текущий уровень
  level_progress_percent: number // 0-100, прогресс до следующего уровня
  coins_balance: number // баланс монет
  tasks_solved_correct_count: number // количество правильно решенных заданий
  wins_count: number // количество побед над злодеями
  checks_correct_count: number // количество успешных проверок
  current_streak_days: number // текущая серия дней
  has_unfinished_attempt: boolean
  mascot_id: string
  mascot_state: 'idle' | 'happy' | 'thinking' | 'celebrating'
  active_villain_id: string | null
  invited_count_total: number // всего приглашено друзей
  achievements_unlocked_count: number
  created_at: string // ISO 8601
  updated_at: string // ISO 8601
}
```

### Используется в экранах
- Home (header, mascot, villain)
- Profile
- Achievements
- Friends

---

## 2. Parent User (Родительский аккаунт)

### Endpoint
```
GET /api/v1/parent-user/:parentUserId
```

### Response
```typescript
interface ParentUser {
  parent_user_id: string
  email: string
  email_verified: boolean
  platform_type: 'vk' | 'max' | 'telegram' | 'web'
  subscription_status: 'free' | 'trial' | 'active' | 'expired' | 'cancelled'
  trial_status: 'not_started' | 'active' | 'expired'
  trial_days_left: number | null
  weekly_report_enabled: boolean
  report_archive_enabled: boolean
  created_at: string
  updated_at: string
}
```

### Используется в экранах
- Profile
- Paywall
- Report Settings

---

## 3. Attempts (Попытки)

### Создание попытки
```
POST /api/v1/attempts
```

**Request:**
```typescript
interface CreateAttemptRequest {
  child_profile_id: string
  mode: 'help' | 'check'
  scenario_type?: 'single_photo' | 'two_photo' // только для check
}
```

**Response:**
```typescript
interface Attempt {
  attempt_id: string
  child_profile_id: string
  mode: 'help' | 'check'
  scenario_type: 'single_photo' | 'two_photo' | null
  attempt_status: 'created' | 'uploading' | 'processing' | 'completed' | 'failed' | 'cancelled'
  created_at: string
  updated_at: string
}
```

### Загрузка изображения
```
POST /api/v1/attempts/:attemptId/images
```

**Request:**
```typescript
// multipart/form-data
{
  image: File
  image_role: 'task' | 'answer' | 'single' // task/answer для check, single для help
}
```

**Response:**
```typescript
interface AttemptImage {
  attempt_image_id: string
  attempt_id: string
  image_role: 'task' | 'answer' | 'single'
  image_url: string
  thumbnail_url: string
  file_size_bytes: number
  mime_type: string
  width: number
  height: number
  uploaded_at: string
}
```

### Получение незаконченной попытки
```
GET /api/v1/attempts/unfinished?child_profile_id=xxx
```

**Response:**
```typescript
interface UnfinishedAttempt {
  attempt_id: string
  mode: 'help' | 'check'
  scenario_type: 'single_photo' | 'two_photo' | null
  attempt_status: 'uploading' | 'processing'
  images: AttemptImage[]
  created_at: string
}
```

### Запуск обработки
```
POST /api/v1/attempts/:attemptId/process
```

**Response:**
```typescript
interface ProcessResponse {
  attempt_id: string
  attempt_status: 'processing'
  estimated_time_seconds: number // примерное время обработки
}
```

### Получение результата (Help)
```
GET /api/v1/attempts/:attemptId/result
```

**Response:**
```typescript
interface HelpResult {
  attempt_id: string
  attempt_status: 'completed' | 'processing' | 'failed'
  hints: Hint[]
  used_hints_count: number
  answer_submitted: boolean
  answer_text?: string
  reward?: {
    coins_earned: number
    xp_earned: number
    villain_damage: number
  }
}

interface Hint {
  hint_level: 1 | 2 | 3
  hint_text: string
  hint_images?: string[]
  unlocked: boolean
}
```

### Получение результата (Check)
```
GET /api/v1/attempts/:attemptId/result
```

**Response:**
```typescript
interface CheckResult {
  attempt_id: string
  attempt_status: 'completed' | 'processing' | 'failed'
  result_status: 'correct' | 'has_errors' | 'wrong'
  error_count: number
  errors: CheckError[]
  reward?: {
    coins_earned: number
    xp_earned: number
    villain_damage: number
  }
}

interface CheckError {
  error_block_id: string
  step_number?: number
  line_reference?: string
  error_type: 'calculation' | 'logic' | 'formatting' | 'missing_step'
  error_message: string // мягкая формулировка
  error_hint: string // подсказка как исправить
  location_type: 'step' | 'line' | 'general'
}
```

### Отправка ответа (Help)
```
POST /api/v1/attempts/:attemptId/submit-answer
```

**Request:**
```typescript
interface SubmitAnswerRequest {
  answer_text: string
  used_hints_count: number
}
```

**Response:**
```typescript
interface SubmitAnswerResponse {
  is_correct: boolean
  reward: {
    coins_earned: number
    xp_earned: number
    villain_damage: number
  }
  profile_updates: {
    level?: number
    level_progress_percent: number
    coins_balance: number
    tasks_solved_correct_count: number
  }
}
```

### Исправление и повтор (Check)
```
POST /api/v1/attempts/:attemptId/resubmit
```

**Request:**
```typescript
interface ResubmitRequest {
  fixed: boolean // пользователь исправил ошибки
}
```

**Response:**
```typescript
interface ResubmitResponse {
  new_attempt_id: string
  mode: 'check'
}
```

### История попыток
```
GET /api/v1/attempts/history?child_profile_id=xxx&limit=20&offset=0&filter=all
```

**Query Parameters:**
- `filter`: 'all' | 'help' | 'check' | 'correct' | 'errors'

**Response:**
```typescript
interface AttemptsHistoryResponse {
  attempts: HistoryAttempt[]
  total: number
  has_more: boolean
}

interface HistoryAttempt {
  attempt_id: string
  mode: 'help' | 'check'
  attempt_status: 'completed' | 'failed' | 'cancelled'
  result_status?: 'correct' | 'has_errors' | 'wrong'
  error_count?: number
  thumbnail_url: string
  created_at: string
  completed_at: string
}
```

### Последние попытки (для Home)
```
GET /api/v1/attempts/recent?child_profile_id=xxx&limit=3
```

**Response:**
```typescript
interface RecentAttempt {
  attempt_id: string
  mode: 'help' | 'check'
  thumbnail_url: string
  history_status: 'completed' | 'processing'
  result_preview: string // краткое описание результата
  created_at: string
}
```

---

## 4. Achievements (Достижения)

### Список достижений
```
GET /api/v1/achievements
```

**⚠️ Security Note:**
- Backend автоматически определяет child_profile_id из JWT токена
- Пользователь получает ТОЛЬКО свои достижения

**⚠️ ВАЖНО: Динамические данные!**

Список достижений полностью определяется бекендом:
- ✅ Backend может добавлять новые достижения без изменения frontend
- ✅ Иконки, названия, условия - всё с бекенда
- ✅ Frontend не хардкодит список достижений
- ✅ Порядок на полках (shelf_order) определяет backend

**Response:**
```typescript
interface AchievementsResponse {
  achievements: Achievement[]
  unlocked_count: number
  total_count: number
}

interface Achievement {
  achievement_id: string // UUID или slug
  name: string // "5 дней подряд", "10 проверок ДЗ"
  description: string
  icon: string // любая эмодзи или URL изображения
  category: string // 'streak' | 'tasks' | ... | любая новая
  requirement: {
    type: string // любой тип требования
    target: number // целевое значение
    current: number // текущий прогресс
    description: string // описание условия для UI
  }
  is_unlocked: boolean
  unlocked_at?: string // ISO 8601
  reward: {
    coins: number
    sticker_id?: string
    sticker_name?: string
  }
  shelf_order: number // позиция для отображения (0-N)
  sort_priority: number // приоритет сортировки
}
```

### Используется в экранах
- Achievements
- Home (модал уведомления о новом достижении)

---

## 5. Villain (Злодей)

### Активный злодей
```
GET /api/v1/villain/active?child_profile_id=xxx
```

**Response:**
```typescript
interface ActiveVillain {
  villain_id: string
  name: string // "Кракозябра"
  description: string
  image_url: string // белый персонаж с короной
  health_current: number
  health_max: number
  health_percent: number // 0-100
  taunts: string[] // ["Ха-ха! Попробуй-ка реши задачки!", ...]
  defeat_requirement: {
    type: 'correct_answers'
    target: 3 // победа после 3 правильных ответов
    current: number
  }
  reward: {
    coins: number
    sticker_id: string
    achievement_id?: string
  }
}
```

### Обновление здоровья злодея (автоматически с Backend)
После каждого правильного ответа бекенд обновляет здоровье и возвращает это в reward.

### Используется в экранах
- Home (отображение злодея)
- Villain Screen (детальная информация)
- Victory Screen (награды)

---

## 6. Referrals (Реферальная система)

### Информация о рефералах
```
GET /api/v1/referrals?child_profile_id=xxx
```

**Response:**
```typescript
interface ReferralInfo {
  referral_code: string // abc123
  referral_link: string // https://homework.app/invite/abc123
  invited_count_total: number // всего приглашенных
  targets: ReferralTarget[]
}

interface ReferralTarget {
  target_count: number // 5, 10, 20
  current_count: number // 2
  is_reached: boolean
  reward: {
    sticker_id: string
    sticker_name: string // "Редкий стикер «Дружба»"
    sticker_icon: string // "⭐"
  }
  is_claimed: boolean
}
```

### Получение награды
```
POST /api/v1/referrals/:targetId/claim
```

**Response:**
```typescript
interface ClaimRewardResponse {
  reward_claimed: boolean
  reward: {
    sticker_id: string
    coins: number
  }
}
```

### Используется в экранах
- Friends

---

## 7. Reports (Отчеты родителю)

### Настройки отчетов
```
GET /api/v1/reports/settings?parent_user_id=xxx
```

**Response:**
```typescript
interface ReportSettings {
  email: string
  email_verified: boolean
  weekly_report_enabled: boolean
  report_archive_enabled: boolean
}
```

### Обновление настроек
```
PATCH /api/v1/reports/settings
```

**Request:**
```typescript
interface UpdateReportSettingsRequest {
  parent_user_id: string
  email?: string
  weekly_report_enabled?: boolean
  report_archive_enabled?: boolean
}
```

### Архив отчетов
```
GET /api/v1/reports/archive?parent_user_id=xxx
```

**Response:**
```typescript
interface ReportsArchive {
  reports: Report[]
  total: number
}

interface Report {
  report_id: string
  period_start: string
  period_end: string
  generated_at: string
  sent_at?: string
  status: 'generated' | 'sent' | 'failed'
  download_url: string // PDF
  summary: {
    tasks_completed: number
    time_spent_minutes: number
    topics_covered: string[]
  }
}
```

### Используется в экранах
- Profile > Отчет родителю

---

## 8. Subscription (Подписка)

### Информация о подписке
```
GET /api/v1/subscription?parent_user_id=xxx
```

**Response:**
```typescript
interface SubscriptionInfo {
  subscription_status: 'free' | 'trial' | 'active' | 'expired' | 'cancelled'
  trial_status: 'not_started' | 'active' | 'expired'
  trial_days_left: number | null
  current_plan?: {
    billing_plan_id: string
    plan_name: string
    billing_period: 'monthly' | 'quarterly' | 'annual'
    price_amount: number
    currency: string
    expires_at: string
    auto_renew: boolean
  }
  available_plans: BillingPlan[]
}

interface BillingPlan {
  billing_plan_id: string
  plan_name: string
  billing_period: 'monthly' | 'quarterly' | 'annual'
  price_amount: number
  currency: string
  discount_percent?: number
  features: string[]
  is_popular: boolean
}
```

### Создание платежа
```
POST /api/v1/subscription/payment
```

**Request:**
```typescript
interface CreatePaymentRequest {
  parent_user_id: string
  billing_plan_id: string
  return_url: string
}
```

**Response:**
```typescript
interface CreatePaymentResponse {
  payment_url: string // URL для редиректа на оплату
  payment_id: string
}
```

### Используется в экранах
- Profile > Подписка
- Paywall

---

## 9. Support (Поддержка)

### Отправка сообщения в поддержку
```
POST /api/v1/support/message
```

**Request:**
```typescript
interface SupportMessageRequest {
  parent_user_id: string
  message: string
  screen_name?: string // контекст, с какого экрана написали
  child_profile_id?: string
}
```

**Response:**
```typescript
interface SupportMessageResponse {
  message_sent: boolean
  ticket_id: string
}
```

### Используется в экранах
- Profile > Помощь

---

## 10. Analytics (Аналитика)

### Отправка событий
```
POST /api/v1/analytics/events
```

**Request:**
```typescript
interface AnalyticsEventsRequest {
  events: AnalyticsEvent[]
}

interface AnalyticsEvent {
  event_name: string
  event_time: string // ISO 8601
  parent_user_id?: string
  child_profile_id?: string
  platform_type: 'vk' | 'max' | 'telegram' | 'web'
  session_id: string
  app_version: string
  screen_name?: string
  entry_point?: string
  // ... additional parameters based on event type
  [key: string]: any
}
```

**Response:**
```typescript
interface AnalyticsResponse {
  events_received: number
  events_processed: number
}
```

### Batch отправка
Фронтенд собирает события в очередь и отправляет batch (до 10 событий) каждые 10 секунд или при достижении лимита.

### Используется везде
- Все экраны отправляют аналитические события

---

## 11. Onboarding (Онбординг)

### Создание профиля ребенка
```
POST /api/v1/child-profile
```

**Request:**
```typescript
interface CreateChildProfileRequest {
  parent_user_id: string
  display_name: string
  avatar_id: string
  grade: number // 1-11
}
```

**Response:**
```typescript
interface CreateChildProfileResponse {
  child_profile_id: string
  child_profile: ChildProfile
}
```

### Email верификация
```
POST /api/v1/parent-user/verify-email
```

**Request:**
```typescript
interface VerifyEmailRequest {
  parent_user_id: string
  email: string
}
```

**Response:**
```typescript
interface VerifyEmailResponse {
  verification_sent: boolean
  verification_code_length: number // для UI подсказки
}
```

### Подтверждение кода
```
POST /api/v1/parent-user/confirm-email
```

**Request:**
```typescript
interface ConfirmEmailRequest {
  parent_user_id: string
  verification_code: string
}
```

**Response:**
```typescript
interface ConfirmEmailResponse {
  email_verified: boolean
}
```

### Используется в экранах
- Onboarding flow

---

## Резюме API Endpoints

| Endpoint | Method | Описание |
|----------|--------|----------|
| `/api/v1/child-profile/:id` | GET | Профиль ребенка |
| `/api/v1/child-profile` | POST | Создание профиля |
| `/api/v1/child-profile/:id` | PATCH | Обновление профиля |
| `/api/v1/parent-user/:id` | GET | Родительский аккаунт |
| `/api/v1/attempts` | POST | Создание попытки |
| `/api/v1/attempts/:id/images` | POST | Загрузка изображения |
| `/api/v1/attempts/unfinished` | GET | Незаконченная попытка |
| `/api/v1/attempts/:id/process` | POST | Запуск обработки |
| `/api/v1/attempts/:id/result` | GET | Результат попытки |
| `/api/v1/attempts/:id/submit-answer` | POST | Отправка ответа |
| `/api/v1/attempts/:id/resubmit` | POST | Повторная попытка |
| `/api/v1/attempts/history` | GET | История попыток |
| `/api/v1/attempts/recent` | GET | Последние попытки |
| `/api/v1/achievements` | GET | Список достижений |
| `/api/v1/villain/active` | GET | Активный злодей |
| `/api/v1/referrals` | GET | Информация о рефералах |
| `/api/v1/referrals/:id/claim` | POST | Получение награды |
| `/api/v1/reports/settings` | GET | Настройки отчетов |
| `/api/v1/reports/settings` | PATCH | Обновление настроек |
| `/api/v1/reports/archive` | GET | Архив отчетов |
| `/api/v1/subscription` | GET | Информация о подписке |
| `/api/v1/subscription/payment` | POST | Создание платежа |
| `/api/v1/support/message` | POST | Сообщение в поддержку |
| `/api/v1/analytics/events` | POST | Отправка аналитики |
| `/api/v1/parent-user/verify-email` | POST | Верификация email |
| `/api/v1/parent-user/confirm-email` | POST | Подтверждение кода |

---

## Общие принципы

### Аутентификация
Все запросы должны содержать токен авторизации:
```
Authorization: Bearer <token>
```

### Обработка ошибок
Все API ошибки возвращаются в стандартном формате:
```typescript
interface APIError {
  error_code: string // например, "VALIDATION_ERROR", "UNAUTHORIZED"
  error_message: string // user-friendly сообщение
  details?: Record<string, any>
}
```

### Pagination
Для списков используется cursor-based pagination:
```typescript
interface PaginatedResponse<T> {
  data: T[]
  pagination: {
    total: number
    limit: number
    offset: number
    has_more: boolean
  }
}
```

### Кеширование
Фронтенд кеширует данные с помощью React Query:
- Profile данные: stale time 5 минут
- Attempts: refetch при фокусе
- Achievements: stale time 10 минут
- Villain: stale time 1 минута

### WebSocket (optional, для real-time обновлений)
Для long-polling обработки попыток можно использовать WebSocket:
```
ws://api.example.com/ws/attempts/:attemptId
```

**Events:**
- `processing_progress` - прогресс обработки (%)
- `processing_complete` - обработка завершена
- `processing_error` - ошибка обработки

---

## 12. REST API Endpoints для Miniapp

### 12.1. Authentication

#### VK Sign Validation
```
POST /api/v1/auth/vk
```

**Request:**
```typescript
interface VKAuthRequest {
  sign: string
  vk_user_id: string
  vk_app_id: string
  vk_are_notifications_enabled?: string
  vk_is_app_user?: string
  vk_is_favorite?: string
  vk_language?: string
  vk_platform?: string
  vk_ref?: string
  vk_ts?: string
}
```

**Response:**
```typescript
interface AuthResponse {
  access_token: string // JWT token
  token_type: 'Bearer'
  expires_in: number // секунды
  user: {
    parent_user_id: string
    platform_type: 'vk'
    has_child_profile: boolean
    child_profile_id?: string
  }
}
```

**Используется в:**
- VK Mini Apps инициализация

---

#### Telegram WebApp Validation
```
POST /api/v1/auth/telegram
```

**Request:**
```typescript
interface TelegramAuthRequest {
  init_data: string // raw Telegram WebApp initData string
  // Альтернативно, распарсенные поля:
  query_id?: string
  user?: {
    id: number
    first_name: string
    last_name?: string
    username?: string
    language_code?: string
    is_premium?: boolean
    photo_url?: string
  }
  auth_date: number
  hash: string
  start_param?: string
}
```

**Response:**
```typescript
interface AuthResponse {
  access_token: string // JWT token
  token_type: 'Bearer'
  expires_in: number // секунды
  user: {
    parent_user_id: string
    platform_type: 'telegram'
    has_child_profile: boolean
    child_profile_id?: string
  }
}
```

**Используется в:**
- Telegram Mini App инициализация

---

### 12.2. Tasks & Attempts

#### Upload Task Photo
```
POST /api/v1/tasks/upload
```

**⚠️ Security Note:**
- Backend определяет child_profile_id из JWT токена
- Пользователь НЕ может загрузить задание для другого профиля

**Request:**
```typescript
// multipart/form-data
interface TaskUploadRequest {
  images: File[] // 1-3 фото задания
  mode?: 'help' | 'check' // по умолчанию 'help'
}
```

**Response:**
```typescript
interface TaskUploadResponse {
  attempt_id: string
  child_profile_id: string
  mode: 'help' | 'check'
  attempt_status: 'processing'
  parsed_task?: {
    task_text: string // распознанный текст задачи
    task_images: string[] // URL загруженных изображений
    subject?: string // математика, физика, русский и т.д.
    grade?: number // определенный класс (1-11)
  }
  estimated_time_seconds: number // примерное время обработки
  created_at: string
}
```

**Используется в:**
- Home > Upload Photo
- Camera flow

---

#### Get User Attempts
```
GET /api/v1/attempts
```

**⚠️ Security Note:**
- Backend автоматически фильтрует по child_profile_id из JWT
- Возвращает ТОЛЬКО попытки текущего пользователя

**Query Parameters:**
```typescript
interface GetAttemptsParams {
  limit?: number // default 20
  offset?: number // default 0
  status?: 'all' | 'completed' | 'processing' | 'failed'
  mode?: 'all' | 'help' | 'check'
  sort?: 'recent' | 'oldest'
}
```

**Response:**
```typescript
interface AttemptsListResponse {
  attempts: AttemptSummary[]
  pagination: {
    total: number
    limit: number
    offset: number
    has_more: boolean
  }
}

interface AttemptSummary {
  attempt_id: string
  mode: 'help' | 'check'
  attempt_status: 'created' | 'uploading' | 'processing' | 'completed' | 'failed'
  task_preview: string // краткое описание задачи
  thumbnail_url: string
  created_at: string
  completed_at?: string
  result_summary?: {
    is_correct?: boolean
    error_count?: number
    hints_used?: number
  }
}
```

**Используется в:**
- History Screen
- Home (recent attempts)

---

#### Get Attempt Details
```
GET /api/v1/attempts/:id
```

**⚠️ Security Note:**
- Backend проверяет, что attempt принадлежит пользователю из JWT
- Возвращает 404 если attempt чужой или не существует

**Response:**
```typescript
interface AttemptDetails {
  attempt_id: string
  child_profile_id: string
  mode: 'help' | 'check'
  attempt_status: 'created' | 'uploading' | 'processing' | 'completed' | 'failed'
  task: {
    task_text: string
    task_images: string[]
    subject?: string
    grade?: number
  }
  hints?: HintDetails[]
  check_result?: CheckResultDetails
  answer_submitted?: {
    answer_text: string
    submitted_at: string
    is_correct: boolean
  }
  reward?: RewardDetails
  created_at: string
  updated_at: string
  completed_at?: string
}

interface HintDetails {
  hint_id: string
  hint_level: 1 | 2 | 3
  hint_text: string
  hint_images?: string[]
  unlocked: boolean
  unlocked_at?: string
  cost_coins: number
}

interface CheckResultDetails {
  result_status: 'correct' | 'has_errors' | 'wrong'
  error_count: number
  errors: CheckError[]
  checked_at: string
}

interface RewardDetails {
  coins_earned: number
  xp_earned: number
  villain_damage: number
  level_up?: boolean
  new_level?: number
  achievement_unlocked?: {
    achievement_id: string
    achievement_name: string
  }
}
```

**Используется в:**
- Help Flow (детали попытки)
- Check Flow (результаты проверки)
- History (детальный просмотр)

---

### 12.3. Hints

#### Generate Hints
```
POST /api/v1/attempts/:id/hints
```

**⚠️ Security Note:**
- Backend проверяет владельца attempt
- Генерация подсказок только для своих попыток

**Request:**
```typescript
interface GenerateHintsRequest {
  // Пустое тело или дополнительные параметры
  hint_style?: 'gentle' | 'direct' | 'socratic' // стиль подсказок
}
```

**Response:**
```typescript
interface GenerateHintsResponse {
  attempt_id: string
  hints_generated: boolean
  hints_count: number // всегда 3 (L1, L2, L3)
  hints: HintPreview[]
}

interface HintPreview {
  hint_id: string
  hint_level: 1 | 2 | 3
  unlocked: boolean
  cost_coins: number
  preview?: string // первые N символов для L1, если разблокирована
}
```

**Используется в:**
- Help Flow > Hints generation

---

#### Unlock Next Hint
```
POST /api/v1/attempts/:id/hints/unlock
```

**⚠️ Security Note:**
- Backend проверяет владельца attempt
- Проверяет баланс монет перед разблокировкой
- Списывает монеты транзакционно

**Request:**
```typescript
interface UnlockHintRequest {
  hint_level: 1 | 2 | 3
}
```

**Response:**
```typescript
interface UnlockHintResponse {
  hint_unlocked: boolean
  hint: {
    hint_id: string
    hint_level: 1 | 2 | 3
    hint_text: string
    hint_images?: string[]
    unlocked_at: string
  }
  coins_spent: number
  new_balance: number
  next_hint_available: boolean
  next_hint_cost?: number
}
```

**Используется в:**
- Help Flow > Unlock hint button

---

### 12.4. Check Answer

#### Check Student Answer
```
POST /api/v1/attempts/:id/check
```

**⚠️ Security Note:**
- Backend проверяет владельца attempt
- Проверка только своих попыток

**Request:**
```typescript
// multipart/form-data
interface CheckAnswerRequest {
  answer_photo: File // фото с решением ученика
  answer_text?: string // опционально, если есть текстовый ответ
}
```

**Response:**
```typescript
interface CheckAnswerResponse {
  attempt_id: string
  check_status: 'processing' | 'completed'
  result?: {
    result_status: 'correct' | 'has_errors' | 'wrong'
    error_count: number
    errors: CheckError[]
    feedback: string // общий фидбек по решению
  }
  reward?: RewardDetails
  estimated_time_seconds?: number // если processing
}
```

**Используется в:**
- Check Flow > Upload answer photo
- Help Flow > Submit answer

---

### 12.5. Generate Analogue

#### Generate Analogue Task
```
POST /api/v1/attempts/:id/analogue
```

**⚠️ Security Note:**
- Backend проверяет владельца attempt
- Генерация аналога только для своих попыток

**Request:**
```typescript
interface GenerateAnalogueRequest {
  difficulty?: 'same' | 'easier' | 'harder' // по умолчанию 'same'
}
```

**Response:**
```typescript
interface GenerateAnalogueResponse {
  new_attempt_id: string
  analogue_task: {
    task_text: string
    task_images?: string[]
    difficulty_level: string
    subject: string
  }
  original_attempt_id: string
  created_at: string
}
```

**Используется в:**
- Help Flow > Generate analogue button
- Victory Screen > Try similar task

---

### 12.6. Profile

#### Get Current User Profile
```
GET /api/v1/profile/me
```

**⚠️ Security Note:**
- Backend определяет child_profile_id из JWT токена
- Пользователь получает ТОЛЬКО свой профиль
- Endpoint уже описан в секции 1, но продублирован для полноты

**Response:**
```typescript
// См. секцию 1. Child Profile
interface ChildProfile {
  child_profile_id: string
  parent_user_id: string
  display_name: string
  avatar_id: string
  grade: number
  level: number
  level_progress_percent: number
  coins_balance: number
  tasks_solved_correct_count: number
  wins_count: number
  checks_correct_count: number
  current_streak_days: number
  has_unfinished_attempt: boolean
  mascot_id: string
  mascot_state: 'idle' | 'happy' | 'thinking' | 'celebrating'
  active_villain_id: string | null
  invited_count_total: number
  achievements_unlocked_count: number
  created_at: string
  updated_at: string
}
```

---

#### Update User Profile
```
PATCH /api/v1/profile/me
```

**⚠️ Security Note:**
- Backend определяет child_profile_id из JWT токена
- Пользователь может обновлять ТОЛЬКО свой профиль

**Request:**
```typescript
interface UpdateProfileRequest {
  display_name?: string
  avatar_id?: string
  grade?: number // 1-11
  mascot_id?: string
}
```

**Response:**
```typescript
interface UpdateProfileResponse {
  profile_updated: boolean
  profile: ChildProfile
}
```

**Используется в:**
- Profile Screen > Edit profile
- Onboarding > Setup profile

---

### 12.7. Achievements

#### Get All Achievements
```
GET /api/v1/achievements
```

**⚠️ Security Note:**
- Backend определяет child_profile_id из JWT токена
- Возвращает достижения с прогрессом ТОЛЬКО для текущего пользователя
- Endpoint уже описан в секции 4, но продублирован для полноты

**Response:**
```typescript
// См. секцию 4. Achievements
interface AchievementsResponse {
  achievements: Achievement[]
  unlocked_count: number
  total_count: number
}

interface Achievement {
  achievement_id: string
  name: string
  description: string
  icon: string
  category: string
  requirement: {
    type: string
    target: number
    current: number
    description: string
  }
  is_unlocked: boolean
  unlocked_at?: string
  reward: {
    coins: number
    sticker_id?: string
    sticker_name?: string
  }
  shelf_order: number
  sort_priority: number
}
```

**Используется в:**
- Achievements Screen
- Home > Achievement notification

---

### 12.8. Referrals

#### Get Referral Information
```
GET /api/v1/referrals
```

**⚠️ Security Note:**
- Backend определяет child_profile_id из JWT токена
- Возвращает реферальную информацию ТОЛЬКО для текущего пользователя
- Endpoint уже описан в секции 6, но продублирован для полноты

**Response:**
```typescript
// См. секцию 6. Referrals
interface ReferralInfo {
  referral_code: string
  referral_link: string
  invited_count_total: number
  targets: ReferralTarget[]
}

interface ReferralTarget {
  target_count: number
  current_count: number
  is_reached: boolean
  reward: {
    sticker_id: string
    sticker_name: string
    sticker_icon: string
  }
  is_claimed: boolean
}
```

**Используется в:**
- Friends Screen

---

#### Apply Referral Code
```
POST /api/v1/referrals/apply
```

**⚠️ Security Note:**
- Backend определяет child_profile_id из JWT токена
- Применяет код ТОЛЬКО для текущего пользователя
- Проверяет, что пользователь не применял код ранее

**Request:**
```typescript
interface ApplyReferralRequest {
  referral_code: string // код реферала (abc123)
}
```

**Response:**
```typescript
interface ApplyReferralResponse {
  referral_applied: boolean
  referrer_info: {
    display_name: string // замаскированное имя реферера
    avatar_id: string
  }
  reward?: {
    coins_earned: number
    bonus_description: string
  }
  error_code?: 'ALREADY_APPLIED' | 'INVALID_CODE' | 'SELF_REFERRAL'
  error_message?: string
}
```

**Используется в:**
- Friends Screen > Apply referral code
- Onboarding > Enter friend's code

---

### 12.9. Villain

#### Get Current Villain
```
GET /api/v1/villain
```

**⚠️ Security Note:**
- Backend определяет child_profile_id из JWT токена
- Возвращает активного злодея для текущего пользователя
- Endpoint уже описан в секции 5, но продублирован для полноты

**Response:**
```typescript
// См. секцию 5. Villain
interface ActiveVillain {
  villain_id: string
  name: string
  description: string
  image_url: string
  health_current: number
  health_max: number
  health_percent: number
  taunts: string[]
  defeat_requirement: {
    type: 'correct_answers'
    target: number
    current: number
  }
  reward: {
    coins: number
    sticker_id: string
    achievement_id?: string
  }
}
```

**Альтернативный response, если злодей не активен:**
```typescript
interface NoActiveVillain {
  villain_id: null
  message: string // "Нет активного злодея. Реши задачу, чтобы встретить нового!"
}
```

**Используется в:**
- Home Screen > Villain widget
- Villain Screen
- Victory Screen

---

## 13. Обновленная сводная таблица endpoints

| Endpoint | Method | Описание | Секция |
|----------|--------|----------|--------|
| `/api/v1/profile/me` | GET | Профиль текущего пользователя | 1, 12.6 |
| `/api/v1/profile/me` | PATCH | Обновление профиля | 12.6 |
| `/api/v1/child-profile` | POST | Создание профиля | 11 |
| `/api/v1/parent-user/:id` | GET | Родительский аккаунт | 2 |
| `/api/v1/auth/vk` | POST | VK авторизация | 12.1 |
| `/api/v1/auth/telegram` | POST | Telegram авторизация | 12.1 |
| `/api/v1/tasks/upload` | POST | Загрузка фото задания | 12.2 |
| `/api/v1/attempts` | GET | Список попыток пользователя | 12.2 |
| `/api/v1/attempts` | POST | Создание попытки | 3 |
| `/api/v1/attempts/:id` | GET | Детали попытки | 12.2 |
| `/api/v1/attempts/:id/images` | POST | Загрузка изображения | 3 |
| `/api/v1/attempts/unfinished` | GET | Незаконченная попытка | 3 |
| `/api/v1/attempts/:id/process` | POST | Запуск обработки | 3 |
| `/api/v1/attempts/:id/result` | GET | Результат попытки | 3 |
| `/api/v1/attempts/:id/hints` | POST | Генерация подсказок | 12.3 |
| `/api/v1/attempts/:id/hints/unlock` | POST | Разблокировка подсказки | 12.3 |
| `/api/v1/attempts/:id/check` | POST | Проверка ответа | 12.4 |
| `/api/v1/attempts/:id/analogue` | POST | Генерация аналога | 12.5 |
| `/api/v1/attempts/:id/submit-answer` | POST | Отправка ответа | 3 |
| `/api/v1/attempts/:id/resubmit` | POST | Повторная попытка | 3 |
| `/api/v1/attempts/history` | GET | История попыток | 3 |
| `/api/v1/attempts/recent` | GET | Последние попытки | 3 |
| `/api/v1/achievements` | GET | Список достижений | 4, 12.7 |
| `/api/v1/villain` | GET | Текущий злодей | 12.9 |
| `/api/v1/villain/active` | GET | Активный злодей (legacy) | 5 |
| `/api/v1/referrals` | GET | Информация о рефералах | 6, 12.8 |
| `/api/v1/referrals/apply` | POST | Применить реферальный код | 12.8 |
| `/api/v1/referrals/:id/claim` | POST | Получение награды | 6 |
| `/api/v1/reports/settings` | GET | Настройки отчетов | 7 |
| `/api/v1/reports/settings` | PATCH | Обновление настроек | 7 |
| `/api/v1/reports/archive` | GET | Архив отчетов | 7 |
| `/api/v1/subscription` | GET | Информация о подписке | 8 |
| `/api/v1/subscription/payment` | POST | Создание платежа | 8 |
| `/api/v1/support/message` | POST | Сообщение в поддержку | 9 |
| `/api/v1/analytics/events` | POST | Отправка аналитики | 10 |
| `/api/v1/parent-user/verify-email` | POST | Верификация email | 11 |
| `/api/v1/parent-user/confirm-email` | POST | Подтверждение кода | 11 |

---

## 14. Дополнительные требования для Miniapp

### 14.1. Rate Limiting

Все endpoints имеют rate limiting:
- Аутентификация: 10 запросов/минуту
- Загрузка фото: 5 запросов/минуту
- Генерация подсказок/аналогов: 3 запроса/минуту
- Остальные endpoints: 60 запросов/минуту

**Response при превышении лимита:**
```typescript
interface RateLimitError {
  error_code: 'RATE_LIMIT_EXCEEDED'
  error_message: string
  retry_after_seconds: number
}
```

HTTP Status: `429 Too Many Requests`

---

### 14.2. Image Upload Constraints

**Требования к изображениям:**
- Форматы: JPEG, PNG, HEIC
- Максимальный размер: 10 MB на файл
- Максимальное разрешение: 4096x4096 px
- Минимальное разрешение: 200x200 px

**Response при ошибке валидации:**
```typescript
interface ImageValidationError {
  error_code: 'INVALID_IMAGE'
  error_message: string
  details: {
    field: 'images' | 'answer_photo'
    reason: 'TOO_LARGE' | 'UNSUPPORTED_FORMAT' | 'TOO_SMALL'
    max_size_mb?: number
    min_width?: number
    min_height?: number
  }
}
```

HTTP Status: `400 Bad Request`

---

### 14.3. Processing Timeouts

**Максимальное время обработки:**
- Распознавание задачи: 30 секунд
- Генерация подсказок: 60 секунд
- Проверка ответа: 45 секунд
- Генерация аналога: 90 секунд

**Response при таймауте:**
```typescript
interface ProcessingTimeout {
  error_code: 'PROCESSING_TIMEOUT'
  error_message: string
  attempt_id: string
  retry_available: boolean
}
```

HTTP Status: `504 Gateway Timeout`

---

### 14.4. Offline Mode Support

Frontend должен поддерживать offline режим:
- Кеширование профиля пользователя
- Кеширование последних 10 попыток
- Кеширование достижений
- Очередь отложенных запросов (аналитика)

**Sync при восстановлении соединения:**
```
POST /api/v1/sync
```

**Request:**
```typescript
interface SyncRequest {
  last_sync_at: string // ISO 8601
  pending_events: AnalyticsEvent[]
}
```

**Response:**
```typescript
interface SyncResponse {
  profile: ChildProfile
  new_achievements: Achievement[]
  villain_updates?: ActiveVillain
  sync_timestamp: string
}
```
