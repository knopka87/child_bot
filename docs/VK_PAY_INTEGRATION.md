# VK Pay Integration

## Обзор

Полная интеграция платежей через VK Pay для подписок в приложении. Система поддерживает создание платежей, обработку webhooks от VK и автоматическую активацию подписок.

## Архитектура

```
┌─────────────┐
│   Frontend  │
│ (VK Bridge) │
└──────┬──────┘
       │ 1. POST /subscription/subscribe
       ▼
┌──────────────────┐
│ Subscription     │
│ Handler          │
└──────┬───────────┘
       │ 2. CreatePayment()
       ▼
┌──────────────────┐
│ VK Pay Service   │
└──────┬───────────┘
       │ 3. INSERT payment
       ▼
┌──────────────────┐
│ PostgreSQL       │
│ - payments       │
│ - subscriptions  │
└──────────────────┘

      VK Pay ──────────────┐
          │ 4. Webhook     │
          ▼                │
┌──────────────────┐       │
│ VKPayWebhook     │       │
│ Handler          │       │
└──────┬───────────┘       │
       │ 5. ProcessWebhook()│
       ▼                   │
┌──────────────────┐       │
│ VK Pay Service   │       │
└──────┬───────────┘       │
       │ 6. Update payment │
       │ 7. Activate subscription
       ▼
┌──────────────────┐
│ PostgreSQL       │
└──────────────────┘
```

## База данных

### Таблица payments

Хранит все транзакции платежей:

```sql
CREATE TABLE payments (
    id UUID PRIMARY KEY,
    subscription_id BIGINT REFERENCES subscriptions(id),
    child_profile_id UUID NOT NULL,
    plan_id VARCHAR(100) NOT NULL,

    -- Сумма
    amount_cents INTEGER NOT NULL,
    currency VARCHAR(10) DEFAULT 'RUB',

    -- VK Pay данные
    vk_order_id VARCHAR(255) UNIQUE,
    vk_transaction_id VARCHAR(255),
    vk_user_id VARCHAR(255),

    -- Статус
    status VARCHAR(20) NOT NULL,
    -- pending, processing, completed, failed, refunded, cancelled

    -- Метаданные
    ip_address VARCHAR(45),
    user_agent TEXT,
    metadata JSONB,

    -- Даты
    paid_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

### Таблица payment_events

Аудит всех изменений платежа:

```sql
CREATE TABLE payment_events (
    id BIGSERIAL PRIMARY KEY,
    payment_id UUID NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    old_status VARCHAR(20),
    new_status VARCHAR(20),
    vk_event_type VARCHAR(100),
    vk_event_data JSONB,
    error_code VARCHAR(100),
    error_message TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

## API Endpoints

### 1. Получить статус подписки

**GET /subscription/status**

Headers:
```
X-Platform-ID: vk
X-Child-Profile-ID: {uuid}
```

Response:
```json
{
  "status": "active",
  "plan_id": "monthly",
  "plan_name": "Месячная подписка",
  "features": ["unlimited_tasks", "hints", "achievements"],
  "expires_at": "2026-05-19T12:00:00Z",
  "renews_at": "2026-05-19T12:00:00Z",
  "can_cancel": true,
  "can_resume": false
}
```

**Статусы:**
- `trial` - триальный период
- `active` - активная подписка
- `expired` - истекшая подписка
- `cancelled` - отмененная (но действует до конца периода)

### 2. Получить планы подписки

**GET /subscription/plans**

Response:
```json
[
  {
    "id": "monthly",
    "name": "Месячная подписка",
    "description": "Полный доступ ко всем функциям на 1 месяц",
    "price": 49900,
    "currency": "RUB",
    "duration": "month",
    "features": [
      "Неограниченное количество задач",
      "Умные подсказки",
      "Проверка решений",
      "Достижения и награды"
    ],
    "is_popular": true,
    "trial_days": 7,
    "discount_percent": 0
  },
  {
    "id": "yearly",
    "name": "Годовая подписка",
    "description": "Выгодная подписка на целый год - экономия 33%",
    "price": 399900,
    "currency": "RUB",
    "duration": "year",
    "features": [...],
    "is_popular": false,
    "trial_days": 14,
    "discount_percent": 33
  }
]
```

### 3. Оформить подписку

**POST /subscription/subscribe**

Headers:
```
X-Platform-ID: vk
X-Child-Profile-ID: {uuid}
```

Query params (от VK):
```
?vk_user_id=12345&sign={signature}&...
```

Request body:
```json
{
  "plan_id": "monthly",
  "payment_method": "vk_pay"
}
```

Response:
```json
{
  "payment_id": "uuid",
  "order_id": "order_abc123_1234567890",
  "vk_pay_url": "https://vk.com/app12345#order_id=order_abc123_1234567890",
  "amount": 49900,
  "currency": "RUB",
  "status": "pending",
  "expires_at": "2026-04-19T13:30:00Z",
  "metadata": {
    "plan_name": "Месячная подписка",
    "duration": 30
  }
}
```

**Frontend integration:**

```typescript
// 1. Получаем данные для платежа от API
const response = await fetch('/subscription/subscribe', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'X-Platform-ID': 'vk',
    'X-Child-Profile-ID': childProfileId,
  },
  body: JSON.stringify({
    plan_id: 'monthly',
    payment_method: 'vk_pay'
  })
});

const paymentData = await response.json();

// 2. Открываем VK Pay форму через VK Bridge
import bridge from '@vkontakte/vk-bridge';

const result = await bridge.send('VKWebAppOpenPayForm', {
  app_id: VK_APP_ID,
  action: 'pay-to-group',
  params: {
    amount: paymentData.amount / 100, // В рублях
    description: 'Подписка: Месячная',
    order_id: paymentData.order_id,
  }
});

// 3. Обрабатываем результат
if (result.status === 'success') {
  // Платеж успешен - webhook обработает его
  console.log('Payment successful:', result);
} else {
  // Платеж отменен или ошибка
  console.error('Payment failed:', result);
}
```

### 4. Отменить подписку

**DELETE /subscription/cancel**

Headers:
```
X-Platform-ID: vk
X-Child-Profile-ID: {uuid}
```

Response:
```json
{
  "status": "active",
  "cancelled_at": "2026-04-19T12:00:00Z",
  "expires_at": "2026-05-19T23:59:59Z",
  "message": "Подписка отменена. Доступ сохраняется до конца оплаченного периода."
}
```

### 5. Возобновить подписку

**POST /subscription/resume**

Headers:
```
X-Platform-ID: vk
X-Child-Profile-ID: {uuid}
```

Response:
```json
{
  "status": "active",
  "renews_at": "2026-05-19T23:59:59Z",
  "message": "Подписка возобновлена. Автопродление включено."
}
```

## Webhooks

### Endpoint

**POST /webhooks/vk-pay**

VK будет отправлять уведомления о платежах на этот endpoint.

### Типы уведомлений

#### 1. Confirmation

Первое событие при настройке webhook в VK:

Request:
```json
{
  "type": "confirmation",
  "group_id": 12345
}
```

Response (plain text):
```
your_confirmation_code_here
```

**Настройка:**
В VK App Settings → Payments → Notification URL нужно указать:
```
https://your-domain.com/webhooks/vk-pay
```

VK отправит confirmation событие и ожидает код подтверждения.

#### 2. Order Status Change

Изменение статуса заказа:

Request:
```json
{
  "type": "order_status_change",
  "object": {
    "order_id": "order_abc123_1234567890",
    "status": "charged",
    "amount": 499,
    "user_id": 12345
  },
  "group_id": 12345,
  "event_id": "event_123"
}
```

Response (plain text):
```
ok
```

**VK Статусы → Наши статусы:**
- `chargeable` → `processing`
- `charged` → `completed` (активирует подписку)
- `refunded` → `refunded`
- `declined` → `failed`
- `cancelled` → `cancelled`

## VK Pay Service

### Создание платежа

```go
service := NewVKPayService(store, VKPayConfig{
    AppID:       os.Getenv("VK_APP_ID"),
    AppSecret:   os.Getenv("VK_APP_SECRET"),
    CallbackURL: "https://your-domain.com/webhooks/vk-pay",
})

req := CreatePaymentRequest{
    ChildProfileID: childProfileID,
    PlanID:         "monthly",
    IPAddress:      "192.168.1.1",
    UserAgent:      "Mozilla/5.0...",
    VKUserID:       "12345",
}

paymentOrder, err := service.CreatePayment(ctx, req)
```

**Что происходит:**
1. Загружает план подписки из БД
2. Генерирует уникальный order_id
3. Создает запись в таблице `payments` со статусом `pending`
4. Логирует событие "created" в `payment_events`
5. Формирует VK Pay URL для frontend
6. Платеж действителен 30 минут

### Обработка webhook

```go
err := service.ProcessWebhook(ctx, vkUserID, orderID, notificationType, payload)
```

**Что происходит:**
1. Загружает платеж по `order_id`
2. Логирует событие "webhook_received"
3. Обрабатывает тип уведомления
4. Для "order_status_change":
   - Маппит VK статус на наш статус
   - Обновляет платеж в БД
   - Логирует "status_changed"
   - Если статус `completed` → активирует подписку

### Активация подписки

```go
err := service.activateSubscription(ctx, payment)
```

**Что происходит:**
1. Загружает план подписки
2. Проверяет наличие активной подписки
3. **Если подписка существует:**
   - Продлевает на duration_days от текущего expires_at
   - Устанавливает статус `active`
   - Включает auto_renew
4. **Если подписки нет:**
   - Создает новую подписку
   - started_at = NOW()
   - expires_at = started_at + duration_days
5. Связывает платеж с подпиской

## Переменные окружения

```bash
# VK App credentials
VK_APP_ID=12345678
VK_APP_SECRET=your_app_secret_here

# VK Confirmation code для webhook
VK_CONFIRMATION_CODE=your_confirmation_code
```

## Настройка в VK

1. **Создать VK Mini App**
   - Перейти в https://vk.com/apps?act=manage
   - Создать новое приложение
   - Скопировать App ID и App Secret

2. **Настроить Payments**
   - Settings → Payments
   - Включить VK Pay
   - Notification URL: `https://your-domain.com/webhooks/vk-pay`
   - VK отправит confirmation событие
   - Вернуть confirmation code
   - Сохранить настройки

3. **Тестирование**
   - Использовать тестовый режим VK Pay
   - Webhooks будут с типом `order_status_change_test`

## Безопасность

### 1. Валидация VK Sign

Frontend передает VK параметры с подписью:
```
?vk_user_id=12345&vk_app_id=12345678&sign={hmac_signature}
```

VKAuthMiddleware проверяет подпись:
```go
// Собираем все vk_* параметры (кроме sign)
params := collectVKParams(query)

// Сортируем по ключу
sort.Strings(params)

// Создаем строку для подписи
queryString := strings.Join(params, "&")

// Вычисляем HMAC-SHA256
mac := hmac.New(sha256.New, []byte(VK_APP_SECRET))
mac.Write([]byte(queryString))
expectedSign := base64URLEncode(mac.Sum(nil))

// Сравниваем constant-time
hmac.Equal(expectedSign, receivedSign)
```

### 2. IP Whitelist

Опционально можно добавить whitelist VK IP адресов для webhooks.

### 3. Логирование

Все платежные события логируются в `payment_events`:
- Создание платежа
- Получение webhook
- Изменение статуса
- Ошибки

### 4. Идемпотентность

Webhook может быть отправлен несколько раз. Обработка идемпотентна:
- Используем `vk_order_id` как уникальный ключ
- Проверяем текущий статус перед обновлением
- Transaction для атомарности операций

## Мониторинг

### Метрики для отслеживания:

1. **Payment metrics:**
   - Количество платежей по статусам
   - Conversion rate (created → completed)
   - Average payment time
   - Failed payment rate

2. **Webhook metrics:**
   - Webhook response time
   - Webhook error rate
   - Duplicate webhook count

3. **Subscription metrics:**
   - Active subscriptions
   - Churn rate
   - Revenue по планам

### SQL запросы:

```sql
-- Платежи по статусам за последние 24 часа
SELECT status, COUNT(*) as count, SUM(amount_cents) as total
FROM payments
WHERE created_at > NOW() - INTERVAL '24 hours'
GROUP BY status;

-- Conversion rate
SELECT
  COUNT(*) FILTER (WHERE status = 'completed') * 100.0 / COUNT(*) as conversion_rate
FROM payments
WHERE created_at > NOW() - INTERVAL '7 days';

-- Активные подписки
SELECT plan_id, COUNT(*) as active_count
FROM subscriptions
WHERE status IN ('trial', 'active')
GROUP BY plan_id;

-- Проблемные платежи (pending > 1 час)
SELECT id, vk_order_id, created_at
FROM payments
WHERE status = 'pending'
  AND created_at < NOW() - INTERVAL '1 hour';
```

## Troubleshooting

### Платеж не активировал подписку

**Проверить:**
1. Логи webhook handler: `docker logs api | grep VKPayWebhook`
2. События платежа: `SELECT * FROM payment_events WHERE payment_id = '{uuid}'`
3. Статус платежа: `SELECT * FROM payments WHERE id = '{uuid}'`

**Возможные причины:**
- Webhook не дошел до сервера
- Ошибка при обработке webhook
- Неправильный маппинг VK статуса

### Webhook возвращает ошибку

**Проверить:**
1. VK_APP_SECRET правильно настроен
2. Order ID существует в БД
3. Логи: `docker logs api | grep VKPayService`

### Платеж висит в pending

**Действия:**
1. Проверить в VK dashboard статус платежа
2. Если в VK платеж завершен, но у нас pending - вручную вызвать ProcessWebhook
3. Настроить автоматическую очистку expired pending платежей:

```sql
UPDATE payments
SET status = 'cancelled'
WHERE status = 'pending'
  AND expires_at < NOW();
```

## Cron Jobs

Добавить периодические задачи:

```go
// Каждый час: истечение pending платежей
func ExpirePendingPayments(ctx context.Context, store *store.Store) error {
    query := `
        UPDATE payments
        SET status = 'cancelled', updated_at = NOW()
        WHERE status = 'pending' AND expires_at < NOW()
    `
    result, err := store.DB.ExecContext(ctx, query)
    // ...
}

// Каждый день: истечение подписок
func ExpireSubscriptions(ctx context.Context, store *store.Store) error {
    rows, err := store.ExpireSubscriptions(ctx)
    log.Printf("Expired %d subscriptions", rows)
    return err
}
```

## Тестирование

### 1. Unit тесты

```go
func TestCreatePayment(t *testing.T) {
    // Создать payment
    // Проверить что создан с правильными параметрами
    // Проверить событие "created"
}

func TestProcessWebhook_Completed(t *testing.T) {
    // Создать pending payment
    // Отправить webhook с status = charged
    // Проверить payment стал completed
    // Проверить подписка активирована
}
```

### 2. Integration тесты

```bash
# 1. Создать платеж
curl -X POST http://localhost:8080/subscription/subscribe \
  -H "Content-Type: application/json" \
  -H "X-Platform-ID: vk" \
  -H "X-Child-Profile-ID: {uuid}" \
  -d '{"plan_id":"monthly","payment_method":"vk_pay"}'

# 2. Имитировать webhook
curl -X POST http://localhost:8080/webhooks/vk-pay \
  -H "Content-Type: application/json" \
  -d '{
    "type": "order_status_change",
    "object": {
      "order_id": "order_abc123_1234567890",
      "status": "charged",
      "user_id": 12345
    }
  }'

# 3. Проверить статус подписки
curl http://localhost:8080/subscription/status \
  -H "X-Platform-ID: vk" \
  -H "X-Child-Profile-ID: {uuid}"
```

## Миграция данных

Для применения изменений в БД:

```bash
# Применить миграцию
make migrate-up

# Или вручную
psql -d child_bot_db -f api/migrations/056_payments.up.sql
```

## Sources

- [VK Mini Apps API](https://github.com/VKCOM/vk-mini-apps-api)
- [VK Bridge](https://github.com/VKCOM/vk-bridge)
- [VK Pay Coverage](https://payatlas.com/payment-method/vk-pay-5226)
