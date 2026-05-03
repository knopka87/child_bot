# Consent System (Система согласий)

## Обзор

Система согласий обеспечивает соответствие требованиям законодательства (GDPR, COPPA, российские законы о персональных данных) и требованиям VK модерации.

## Требования законодательства

### GDPR (General Data Protection Regulation)

**Основные требования:**
- Явное согласие на обработку данных
- Возможность отозвать согласие
- Полная история согласий (audit trail)
- Информирование о цели сбора данных

### COPPA (Children's Online Privacy Protection Act)

**Основные требования:**
- Согласие родителей для детей до 13 лет
- Верификация возраста родителя
- Минимизация собираемых данных

### Российское законодательство

**152-ФЗ "О персональных данных":**
- Согласие на обработку ПД
- Информирование об операторе ПД
- Цели обработки данных
- Срок хранения данных

## Архитектура

### Таблицы БД

#### parent_consents

Текущее состояние согласий:

```sql
CREATE TABLE parent_consents (
    id UUID PRIMARY KEY,
    parent_user_id VARCHAR(255) NOT NULL,
    platform_id VARCHAR(20) NOT NULL,

    -- Согласия
    privacy_policy_version VARCHAR(20) NOT NULL,
    privacy_policy_accepted BOOLEAN NOT NULL,
    privacy_policy_accepted_at TIMESTAMPTZ NOT NULL,

    terms_version VARCHAR(20) NOT NULL,
    terms_accepted BOOLEAN NOT NULL,
    terms_accepted_at TIMESTAMPTZ NOT NULL,

    adult_consent BOOLEAN NOT NULL,
    adult_consent_at TIMESTAMPTZ,

    -- Audit
    ip_address VARCHAR(45),
    user_agent TEXT,

    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    UNIQUE (platform_id, parent_user_id)
);
```

#### parent_consent_history

Полная история изменений (audit trail):

```sql
CREATE TABLE parent_consent_history (
    id UUID PRIMARY KEY,
    consent_id UUID NOT NULL,
    parent_user_id VARCHAR(255) NOT NULL,
    platform_id VARCHAR(20) NOT NULL,

    -- Тип изменения
    action VARCHAR(20) NOT NULL, -- 'created', 'updated', 'revoked'

    -- Snapshot на момент изменения
    privacy_policy_version VARCHAR(20) NOT NULL,
    privacy_policy_accepted BOOLEAN NOT NULL,
    terms_version VARCHAR(20) NOT NULL,
    terms_accepted BOOLEAN NOT NULL,
    adult_consent BOOLEAN NOT NULL,

    -- Что изменилось (для action='updated')
    changed_fields TEXT[],
    previous_values JSONB,

    -- Audit metadata
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL
);
```

### API Endpoints

#### POST /api/v1/consent

Сохранение согласия родителя.

**Request:**
```json
{
  "parentUserId": "vk_user_12345",
  "privacyPolicyVersion": "2.0",
  "termsVersion": "1.5",
  "adultConsent": true
}
```

**Response:**
```json
{
  "message": "Consent saved successfully"
}
```

**Логика:**
1. Проверяет платформу через middleware
2. Валидирует все поля
3. Получает IP и User-Agent
4. **Начинает транзакцию**
5. Проверяет существующее согласие
6. Определяет изменения (created/updated)
7. Сохраняет/обновляет согласие (UPSERT)
8. Записывает в историю
9. **Коммитит транзакцию**

#### GET /api/v1/consent?parentUserId=xxx

Получение текущего согласия.

**Response:**
```json
{
  "id": "uuid",
  "parentUserId": "vk_user_12345",
  "privacyPolicyVersion": "2.0",
  "privacyPolicyAccepted": true,
  "termsVersion": "1.5",
  "termsAccepted": true,
  "adultConsent": true,
  "createdAt": "2024-04-19T12:00:00Z",
  "updatedAt": "2024-04-19T12:00:00Z"
}
```

#### GET /api/v1/consent/check?parentUserId=xxx

Быстрая проверка наличия валидного согласия.

**Response:**
```json
{
  "hasConsent": true
}
```

**Используется для:**
- Проверка перед показом контента
- Быстрая валидация без полных данных
- Условная логика в приложении

#### GET /api/v1/consent/history?parentUserId=xxx&limit=50

Получение истории изменений согласий (audit trail).

**Response:**
```json
{
  "history": [
    {
      "id": "uuid1",
      "action": "updated",
      "privacyPolicyVersion": "2.0",
      "privacyPolicyAccepted": true,
      "termsVersion": "1.5",
      "termsAccepted": true,
      "adultConsent": true,
      "changedFields": ["privacy_policy_version"],
      "previousValues": {
        "privacy_policy_version": "1.0"
      },
      "createdAt": "2024-04-19T14:00:00Z"
    },
    {
      "id": "uuid2",
      "action": "created",
      "privacyPolicyVersion": "1.0",
      "privacyPolicyAccepted": true,
      "termsVersion": "1.0",
      "termsAccepted": true,
      "adultConsent": true,
      "createdAt": "2024-04-19T12:00:00Z"
    }
  ],
  "count": 2
}
```

## Версионирование документов

### Формат версий

Используется **Semantic Versioning** для документов:

```
MAJOR.MINOR
```

**Примеры:**
- `1.0` - первая версия
- `1.1` - minor изменения (исправления, уточнения)
- `2.0` - major изменения (существенные изменения условий)

### Когда увеличивать версию

**MAJOR (2.0, 3.0):**
- Существенные изменения условий
- Новые обязательства пользователя
- Изменение политики обработки данных
- Изменение прав и обязанностей сторон

**MINOR (1.1, 1.2):**
- Исправление опечаток
- Уточнения формулировок
- Добавление примеров
- Косметические изменения

### Обновление версий

При обновлении документов:

1. Увеличиваем версию в коде
2. Пользователи видят новую версию при следующем входе
3. Если major изменение - требуем повторное согласие
4. История сохраняется автоматически

**Пример:**
```typescript
// В коде приложения
const CURRENT_PRIVACY_POLICY_VERSION = "2.0"
const CURRENT_TERMS_VERSION = "1.5"

// При изменении документа - просто меняем константу
const CURRENT_PRIVACY_POLICY_VERSION = "3.0" // major изменение
```

## Frontend Integration

### Первичное согласие (Onboarding)

```typescript
interface ConsentData {
  parentUserId: string
  privacyPolicyVersion: string
  termsVersion: string
  adultConsent: boolean
}

async function saveConsent(data: ConsentData) {
  const response = await fetch('/api/v1/consent', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-Platform-ID': 'vk',
    },
    body: JSON.stringify(data),
  })

  if (!response.ok) {
    throw new Error('Failed to save consent')
  }

  return response.json()
}

// Использование
await saveConsent({
  parentUserId: 'vk_12345',
  privacyPolicyVersion: '2.0',
  termsVersion: '1.5',
  adultConsent: true,
})
```

### Проверка согласия

```typescript
async function checkConsent(parentUserId: string): Promise<boolean> {
  const response = await fetch(
    `/api/v1/consent/check?parentUserId=${parentUserId}`,
    {
      headers: {
        'X-Platform-ID': 'vk',
      },
    }
  )

  const data = await response.json()
  return data.hasConsent
}

// Использование в защищенном route
if (!(await checkConsent(parentUserId))) {
  // Редирект на страницу согласий
  router.push('/consent')
}
```

### Просмотр истории

```typescript
interface ConsentHistoryEntry {
  id: string
  action: 'created' | 'updated' | 'revoked'
  privacyPolicyVersion: string
  termsVersion: string
  adultConsent: boolean
  changedFields?: string[]
  previousValues?: Record<string, any>
  createdAt: string
}

async function getConsentHistory(
  parentUserId: string,
  limit: number = 50
): Promise<ConsentHistoryEntry[]> {
  const response = await fetch(
    `/api/v1/consent/history?parentUserId=${parentUserId}&limit=${limit}`,
    {
      headers: {
        'X-Platform-ID': 'vk',
      },
    }
  )

  const data = await response.json()
  return data.history
}
```

### React компонент для согласия

```typescript
function ConsentForm({ onComplete }: { onComplete: () => void }) {
  const [accepted, setAccepted] = useState({
    privacy: false,
    terms: false,
    adult: false,
  })

  const [loading, setLoading] = useState(false)

  const canSubmit = accepted.privacy && accepted.terms && accepted.adult

  const handleSubmit = async () => {
    setLoading(true)
    try {
      await saveConsent({
        parentUserId: getCurrentUserId(),
        privacyPolicyVersion: '2.0',
        termsVersion: '1.5',
        adultConsent: true,
      })
      onComplete()
    } catch (error) {
      console.error('Failed to save consent:', error)
      alert('Ошибка при сохранении согласия')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="consent-form">
      <h2>Согласие на обработку данных</h2>

      <label>
        <input
          type="checkbox"
          checked={accepted.privacy}
          onChange={(e) => setAccepted({ ...accepted, privacy: e.target.checked })}
        />
        Я согласен с{' '}
        <a href="/legal/privacy" target="_blank">
          Политикой конфиденциальности v2.0
        </a>
      </label>

      <label>
        <input
          type="checkbox"
          checked={accepted.terms}
          onChange={(e) => setAccepted({ ...accepted, terms: e.target.checked })}
        />
        Я согласен с{' '}
        <a href="/legal/terms" target="_blank">
          Условиями использования v1.5
        </a>
      </label>

      <label>
        <input
          type="checkbox"
          checked={accepted.adult}
          onChange={(e) => setAccepted({ ...accepted, adult: e.target.checked })}
        />
        Я подтверждаю, что мне исполнилось 18 лет
      </label>

      <button onClick={handleSubmit} disabled={!canSubmit || loading}>
        {loading ? 'Сохранение...' : 'Принять и продолжить'}
      </button>
    </div>
  )
}
```

## Audit Trail

### Что записывается

Для каждого изменения согласия сохраняется:

1. **Тип действия** (created/updated/revoked)
2. **Полный snapshot** состояния
3. **Изменённые поля** (для updated)
4. **Предыдущие значения** (для updated)
5. **IP адрес** пользователя
6. **User-Agent** браузера
7. **Timestamp** с точностью до миллисекунды

### Пример записи в истории

**Создание:**
```json
{
  "action": "created",
  "privacyPolicyVersion": "1.0",
  "termsVersion": "1.0",
  "adultConsent": true,
  "ipAddress": "203.0.113.1",
  "userAgent": "Mozilla/5.0...",
  "createdAt": "2024-04-19T12:00:00.123Z"
}
```

**Обновление:**
```json
{
  "action": "updated",
  "privacyPolicyVersion": "2.0",
  "termsVersion": "1.5",
  "adultConsent": true,
  "changedFields": ["privacy_policy_version", "terms_version"],
  "previousValues": {
    "privacy_policy_version": "1.0",
    "terms_version": "1.0"
  },
  "ipAddress": "203.0.113.1",
  "userAgent": "Mozilla/5.0...",
  "createdAt": "2024-04-19T14:30:00.456Z"
}
```

### Хранение данных

**Retention policy:**
- История согласий хранится **бессрочно**
- Требуется для соответствия GDPR и законодательству РФ
- При удалении пользователя - история **НЕ** удаляется (legal requirement)

**Backup:**
- Ежедневные бэкапы базы данных
- Отдельные бэкапы таблицы `parent_consent_history`
- Хранение в защищённом S3 bucket

## Security

### IP Address sanitization

IP адреса используются только для audit:

```go
ipAddress := r.Header.Get("X-Forwarded-For")
if ipAddress == "" {
    ipAddress = r.RemoteAddr
}
// Сохраняем как есть - не логируем в application logs
```

**Best practices:**
- IP сохраняется только в БД
- Не логируется в application logs
- Доступ к истории только через API (авторизация)

### User-Agent

User-Agent помогает идентифицировать:
- Устройство (Desktop/Mobile)
- Браузер
- ОС

**Privacy:**
- Не содержит PII
- Используется только для audit
- Не передаётся третьим лицам

### Access control

Доступ к согласиям:

```go
// Только владелец может просматривать свои согласия
platformID := middleware.GetPlatformID(r.Context())
// platformID проверяется middleware, подделка невозможна
```

## Compliance Reports

### GDPR Right to Access

Предоставить пользователю все его данные:

```sql
-- Текущее согласие
SELECT * FROM parent_consents
WHERE platform_id = 'vk' AND parent_user_id = 'user_123';

-- История
SELECT * FROM parent_consent_history
WHERE platform_id = 'vk' AND parent_user_id = 'user_123'
ORDER BY created_at DESC;
```

### GDPR Right to be Forgotten

При удалении пользователя:

```sql
-- Удаляем текущее согласие
DELETE FROM parent_consents
WHERE platform_id = 'vk' AND parent_user_id = 'user_123';

-- История НЕ удаляется (legal requirement)
-- Но можно анонимизировать:
UPDATE parent_consent_history
SET parent_user_id = 'DELETED_USER',
    ip_address = NULL,
    user_agent = NULL
WHERE platform_id = 'vk' AND parent_user_id = 'user_123';
```

### Proof of Consent

Для доказательства наличия согласия (regulatory audit):

```sql
SELECT
    pc.parent_user_id,
    pc.privacy_policy_version,
    pc.privacy_policy_accepted_at,
    pc.terms_version,
    pc.terms_accepted_at,
    pc.adult_consent,
    pc.adult_consent_at,
    pc.ip_address,
    pch.created_at as first_consent_at
FROM parent_consents pc
LEFT JOIN (
    SELECT parent_user_id, platform_id, MIN(created_at) as created_at
    FROM parent_consent_history
    WHERE action = 'created'
    GROUP BY parent_user_id, platform_id
) pch ON pc.parent_user_id = pch.parent_user_id
      AND pc.platform_id = pch.platform_id
WHERE pc.platform_id = 'vk';
```

## Monitoring

### Метрики

**Отслеживать:**
1. Количество новых согласий в день
2. Количество обновлений согласий
3. Процент пользователей с валидными согласиями
4. Время на прохождение consent flow

**Alerts:**
- Резкое падение процента согласий
- Много отказов от согласия
- Ошибки сохранения в БД

### Dashboard queries

```sql
-- Согласия за сегодня
SELECT COUNT(*)
FROM parent_consent_history
WHERE action = 'created'
  AND created_at >= CURRENT_DATE;

-- Процент валидных согласий
SELECT
    COUNT(CASE WHEN privacy_policy_accepted AND terms_accepted AND adult_consent THEN 1 END)::FLOAT /
    NULLIF(COUNT(*), 0) * 100 as consent_percentage
FROM parent_consents;
```

## Best Practices

### 1. Версионируйте документы

```typescript
// ❌ Плохо: хардкод версии
await saveConsent({
  privacyPolicyVersion: "1.0",
  termsVersion: "1.0",
})

// ✅ Хорошо: константы
const CURRENT_PRIVACY_VERSION = "2.0"
const CURRENT_TERMS_VERSION = "1.5"

await saveConsent({
  privacyPolicyVersion: CURRENT_PRIVACY_VERSION,
  termsVersion: CURRENT_TERMS_VERSION,
})
```

### 2. Показывайте ссылки на документы

```tsx
// ✅ Хорошо: пользователь может прочитать перед согласием
<label>
  <input type="checkbox" />
  Я согласен с{' '}
  <a href="/legal/privacy" target="_blank">
    Политикой конфиденциальности v{CURRENT_PRIVACY_VERSION}
  </a>
</label>
```

### 3. Требуйте явного согласия

```typescript
// ❌ Плохо: согласие по умолчанию
const [accepted, setAccepted] = useState(true)

// ✅ Хорошо: пользователь должен явно согласиться
const [accepted, setAccepted] = useState(false)
```

### 4. Логируйте ошибки, но не PII

```go
// ❌ Плохо: логируем IP
log.Printf("Consent saved for IP: %s", ipAddress)

// ✅ Хорошо: логируем факт без PII
log.Printf("Consent saved successfully")
```

### 5. Используйте транзакции

```go
// ✅ Обязательно: согласие и история в одной транзакции
tx, _ := db.Begin()
saveConsent(tx, ...)
saveHistory(tx, ...)
tx.Commit()
```

## Troubleshooting

### Consent not saving

**Проблема:** POST /consent возвращает 500

**Проверка:**
```bash
# Проверьте логи backend
docker logs child_bot_backend --tail 100 | grep -i consent

# Проверьте что миграции применены
docker exec child_bot_postgres psql -U child_bot -d child_bot -c "\dt parent_consent*"
```

**Решение:**
```bash
# Применить миграции
make migrate-up
```

### History not recorded

**Проблема:** История не сохраняется

**Проверка:**
```sql
-- Проверьте что таблица существует
SELECT COUNT(*) FROM parent_consent_history;

-- Проверьте транзакции
SELECT * FROM pg_stat_activity
WHERE datname = 'child_bot' AND state = 'idle in transaction';
```

### Invalid version format

**Проблема:** Версия документа невалидна

**Решение:**
Используйте формат `MAJOR.MINOR`:
- ✅ `1.0`, `2.0`, `1.5`
- ❌ `v1.0`, `1`, `1.0.0`

## Links

- [GDPR Official Text](https://gdpr-info.eu/)
- [COPPA Compliance Guide](https://www.ftc.gov/business-guidance/resources/complying-coppa-frequently-asked-questions)
- [152-ФЗ О персональных данных](http://www.consultant.ru/document/cons_doc_LAW_61801/)
- [VK Moderation Requirements](https://dev.vk.com/ru/mini-apps/overview)
