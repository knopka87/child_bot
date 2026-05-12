# CSRF Protection

## Обзор

Приложение использует CSRF (Cross-Site Request Forgery) защиту для предотвращения несанкционированных state-changing операций.

## Что такое CSRF?

CSRF атака заставляет авторизованного пользователя выполнить нежелательное действие в веб-приложении. Например:

1. Пользователь авторизован в вашем приложении
2. Пользователь посещает вредоносный сайт evil.com
3. evil.com отправляет POST запрос к вашему API от имени пользователя
4. Без CSRF защиты этот запрос будет выполнен

## Реализация

### Double Submit Cookie Pattern

Используется **Double Submit Cookie** pattern:

1. Server генерирует random CSRF token
2. Token устанавливается как cookie (НЕ HttpOnly, чтобы JS мог прочитать)
3. Client читает token из cookie
4. Client отправляет token в header `X-CSRF-Token` при каждом state-changing запросе
5. Server сравнивает cookie token и header token
6. Если они совпадают - запрос валиден

### Почему это безопасно?

- evil.com может отправить cookie (браузер делает это автоматически)
- Но evil.com **НЕ может прочитать** cookie из-за Same-Origin Policy
- Без возможности прочитать token, evil.com не может установить header

## Защищенные методы

| Метод | Требует CSRF token | Причина |
|-------|-------------------|---------|
| GET | ❌ Нет | Safe method - не изменяет state |
| HEAD | ❌ Нет | Safe method - не изменяет state |
| OPTIONS | ❌ Нет | Safe method - не изменяет state |
| POST | ✅ Да | Изменяет state |
| PUT | ✅ Да | Изменяет state |
| DELETE | ✅ Да | Изменяет state |
| PATCH | ✅ Да | Изменяет state |

## Исключения

Некоторые endpoints освобождены от CSRF проверки:

```go
exemptPaths := []string{
    "/health",                      // Health check - public
    "/api/v1/analytics/events",     // Analytics - может быть public
}
```

**Когда добавлять исключения:**
- Public endpoints без авторизации
- Read-only endpoints
- Webhooks от внешних сервисов

**Когда НЕ добавлять:**
- Любые state-changing операции
- Endpoints требующие авторизацию
- Операции с пользовательскими данными

## Frontend интеграция

### 1. Получение CSRF token

При первом GET запросе backend автоматически устанавливает CSRF cookie. Также можно явно получить token:

```typescript
// GET /api/v1/csrf-token
const response = await fetch('/api/v1/csrf-token')
const data = await response.json()

console.log(data)
// {
//   "csrfToken": "abc123...",
//   "headerName": "X-CSRF-Token",
//   "cookieName": "csrf_token"
// }
```

### 2. Чтение token из cookie

```typescript
function getCSRFToken(): string | null {
  const cookies = document.cookie.split(';')
  for (const cookie of cookies) {
    const [name, value] = cookie.trim().split('=')
    if (name === 'csrf_token') {
      return value
    }
  }
  return null
}
```

### 3. Отправка token в запросах

**Axios:**

```typescript
import axios from 'axios'

// Interceptor для автоматического добавления CSRF token
axios.interceptors.request.use((config) => {
  // Только для unsafe methods
  if (['post', 'put', 'delete', 'patch'].includes(config.method?.toLowerCase() || '')) {
    const csrfToken = getCSRFToken()
    if (csrfToken) {
      config.headers['X-CSRF-Token'] = csrfToken
    }
  }
  return config
})

// Использование
await axios.post('/api/v1/profile', { name: 'New Name' })
// Header X-CSRF-Token добавляется автоматически
```

**Fetch:**

```typescript
async function secureFetch(url: string, options: RequestInit = {}) {
  const csrfToken = getCSRFToken()

  // Для unsafe methods добавляем CSRF token
  if (options.method && !['GET', 'HEAD', 'OPTIONS'].includes(options.method.toUpperCase())) {
    options.headers = {
      ...options.headers,
      'X-CSRF-Token': csrfToken || '',
    }
  }

  return fetch(url, options)
}

// Использование
await secureFetch('/api/v1/profile', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ name: 'New Name' }),
})
```

### 4. React Hook пример

```typescript
import { useEffect, useState } from 'react'

export function useCSRFToken() {
  const [token, setToken] = useState<string | null>(null)

  useEffect(() => {
    // Получаем token при монтировании
    const csrfToken = getCSRFToken()

    if (!csrfToken) {
      // Если нет cookie, запрашиваем у сервера
      fetch('/api/v1/csrf-token')
        .then(res => res.json())
        .then(data => setToken(data.csrfToken))
    } else {
      setToken(csrfToken)
    }
  }, [])

  return token
}

// Использование в компоненте
function ProfileForm() {
  const csrfToken = useCSRFToken()

  const handleSubmit = async () => {
    await fetch('/api/v1/profile', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-CSRF-Token': csrfToken || '',
      },
      body: JSON.stringify({ /* data */ }),
    })
  }

  return <form onSubmit={handleSubmit}>...</form>
}
```

## Backend

### Middleware

**Файл:** `api/internal/api/middleware/csrf.go`

Middleware автоматически:
- Генерирует и устанавливает CSRF token для GET запросов
- Проверяет token для POST/PUT/DELETE/PATCH запросов
- Освобождает от проверки safe methods
- Освобождает от проверки exempt paths

### Development vs Production

**Development (`ENV=development`):**
- Более гибкая валидация
- Логирует warnings вместо blocking
- Автоматически создает token если нет
- Полезно для разработки без необходимости настройки frontend

**Production (`ENV=production`):**
- Строгая валидация
- Блокирует запросы без валидного token
- Возвращает 403 Forbidden

### Регенерация token

После критических операций (login, logout) рекомендуется регенерировать token:

```go
import "child-bot/api/internal/api/middleware"

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    // ... login logic ...

    // Регенерируем CSRF token для дополнительной безопасности
    newToken, err := middleware.RegenerateCSRFToken(w)
    if err != nil {
        log.Printf("Failed to regenerate CSRF token: %v", err)
    }

    // ... return response ...
}
```

## Тестирование

### Unit Tests

```bash
cd api
go test -v ./internal/api/middleware/csrf_test.go ./internal/api/middleware/csrf.go
```

Тесты покрывают:
- ✅ Safe methods проходят без token
- ✅ Unsafe methods блокируются без token
- ✅ Валидный token пропускается
- ✅ Невалидный token блокируется
- ✅ Exempt paths проходят без token
- ✅ Dev режим более гибкий
- ✅ Token generation cryptographically secure
- ✅ Constant-time token comparison

### Manual Testing

```bash
# 1. Получить CSRF token
curl -c cookies.txt http://localhost:8080/api/v1/csrf-token

# 2. Извлечь token из cookie
TOKEN=$(grep csrf_token cookies.txt | awk '{print $7}')

# 3. Попытка POST без token (должен вернуть 403)
curl -X POST http://localhost:8080/api/v1/profile \
  -H "Content-Type: application/json" \
  -d '{"name": "Test"}'

# 4. POST с валидным token (должен пройти)
curl -X POST http://localhost:8080/api/v1/profile \
  -b cookies.txt \
  -H "Content-Type: application/json" \
  -H "X-CSRF-Token: $TOKEN" \
  -d '{"name": "Test"}'
```

## Security Best Practices

### 1. Cookie параметры

```go
cookie := &http.Cookie{
    Name:     "csrf_token",
    Value:    token,
    Path:     "/",
    MaxAge:   86400,           // 24 hours
    HttpOnly: false,           // ⚠️ НЕ HttpOnly (JS должен читать)
    Secure:   isProduction,    // ✅ HTTPS only в production
    SameSite: http.SameSiteStrictMode, // ✅ CSRF защита
}
```

**Почему HttpOnly = false?**
- JavaScript должен прочитать token из cookie
- Это не представляет угрозу, так как token сам по себе не критичен
- Реальная защита - в необходимости отправить token в header

**SameSite = Strict:**
- Cookie отправляется только для same-site запросов
- Дополнительная защита от CSRF
- Блокирует cross-site requests

### 2. Constant-time comparison

```go
func tokensEqual(a, b string) bool {
    if len(a) != len(b) {
        return false
    }

    result := 0
    for i := 0; i < len(a); i++ {
        result |= int(a[i]) ^ int(b[i])
    }

    return result == 0
}
```

Защищает от **timing attacks** - злоумышленник не может узнать token по времени сравнения.

### 3. Cryptographically secure random

```go
bytes := make([]byte, 32)
_, err := rand.Read(bytes)  // crypto/rand, НЕ math/rand
token := base64.URLEncoding.EncodeToString(bytes)
```

- Используется `crypto/rand` для криптографически стойких случайных чисел
- 32 байта = 256 бит энтропии
- Base64 URL-safe encoding

### 4. Token lifetime

```go
const CSRFCookieMaxAge = 86400  // 24 часа
```

- Token действителен 24 часа
- После истечения автоматически генерируется новый
- Не храните token дольше чем необходимо

## Troubleshooting

### 403 Forbidden на POST запросах

**Проблема:** Все POST/PUT/DELETE запросы возвращают 403

**Решение:**
1. Проверьте что frontend читает token из cookie
2. Проверьте что token отправляется в header `X-CSRF-Token`
3. Проверьте что cookie `csrf_token` установлена
4. В dev режиме проверьте логи backend

```bash
docker logs child_bot_backend_dev --tail 50 | grep CSRF
```

### Cookie не устанавливается

**Проблема:** Cookie `csrf_token` не появляется

**Решение:**
1. Убедитесь что сделали хотя бы один GET запрос
2. Проверьте что CORS настроен правильно
3. В dev режиме `Secure` должен быть false
4. Проверьте Domain и Path cookie

### Token mismatch

**Проблема:** Логи показывают "Token mismatch"

**Решение:**
1. Убедитесь что читаете правильную cookie (`csrf_token`)
2. Убедитесь что не модифицируете token (trim, encoding)
3. Проверьте что используете один и тот же token

### Dev режим не работает

**Проблема:** Даже в dev режиме запросы блокируются

**Решение:**
```bash
# Проверьте переменную ENV
grep ENV .env

# Должно быть:
ENV=development

# Перезапустите backend
docker compose -f docker/docker-compose.dev.yml restart backend
```

## VK Mini App специфика

### VK Bridge запросы

VK Bridge может отправлять запросы от имени приложения. Убедитесь что:

1. VK параметры валидируются через VKAuthMiddleware (до CSRF)
2. CSRF token добавляется ко всем VK Bridge API calls
3. Cookie разрешены в VK WebView

### Cross-origin issues

VK Mini Apps работают в iframe. Проверьте:

```go
// CORS должен разрешать credentials
w.Header().Set("Access-Control-Allow-Credentials", "true")

// SameSite может быть проблемой в iframe
// Рассмотрите использование SameSite=None в production для VK
cookie.SameSite = http.SameSiteNoneMode  // только если необходимо
```

## Monitoring

### Логи

Middleware логирует все CSRF события:

```
# Валидный token:
(нет лога - проходит тихо)

# Отсутствует cookie:
[CSRF] Missing CSRF cookie

# Отсутствует header:
[CSRF] Missing CSRF header for POST /api/v1/profile

# Token mismatch:
[CSRF] Token mismatch for POST /api/v1/profile

# Dev режим warnings:
[CSRF] Warning: No CSRF cookie found, generating new one
[CSRF] Warning: No CSRF header found for POST /api/v1/profile
```

### Метрики

Рекомендуется отслеживать:
- Количество CSRF failures
- Endpoints с наибольшим количеством failures
- Время генерации токенов

## Альтернативные подходы

### Synchronizer Token Pattern

Вместо Double Submit Cookie можно использовать:

**Преимущества:**
- Более безопасен (token хранится server-side)
- Не зависит от cookies

**Недостатки:**
- Требует session storage (Redis, DB)
- Сложнее для stateless API

**Когда использовать:**
Если у вас уже есть session management.

### SameSite Cookie только

**SameSite=Strict** cookie сама по себе защищает от CSRF.

**Недостатки:**
- Не поддерживается старыми браузерами
- Может ломать легитимные cross-site flows

**Рекомендация:**
Используйте SameSite + CSRF token для defense in depth.

## Ссылки

- [OWASP CSRF Prevention Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html)
- [Double Submit Cookie Pattern](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html#double-submit-cookie)
- [SameSite Cookie Explained](https://web.dev/samesite-cookies-explained/)
- [MDN: CSRF](https://developer.mozilla.org/en-US/docs/Glossary/CSRF)

## Checklist перед production

- [ ] CSRF middleware включен в production
- [ ] Frontend отправляет X-CSRF-Token header для всех unsafe methods
- [ ] CSRF cookie имеет SameSite=Strict
- [ ] Secure=true в production (HTTPS only)
- [ ] Exempt paths минимальны и обоснованы
- [ ] Тесты покрывают все сценарии
- [ ] Логи мониторятся на CSRF failures
- [ ] VK Bridge интеграция работает с CSRF
