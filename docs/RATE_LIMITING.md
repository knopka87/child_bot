## Rate Limiting

## Обзор

Rate limiting ограничивает количество запросов от одного клиента в единицу времени. Это защищает API от:
- DDoS атак
- Abuse (злоупотребления)
- Случайных петель в клиентском коде
- Перегрузки сервера

## Реализация

### Алгоритм: Sliding Window

Используется **Sliding Window** (скользящее окно) для точного контроля rate:

```
Window = 1 минута
Limit = 100 запросов

Timeline:
|-----|-----|-----|-----|
0     15s   30s   45s   60s

Запросы в 0-15s: 30
Запросы в 15-30s: 40
Запросы в 30-45s: 20
Запросы в 45-60s: 10

В момент 30s:
- Окно: [0s - 30s]
- Запросов: 30 + 40 = 70 ✅ Allowed
```

**Преимущества:**
- Точный контроль rate
- Нет burst проблем
- Справедливое распределение

**Альтернативы:**
- Fixed Window - проще, но допускает bursts
- Token Bucket - более гибкий, но сложнее

## Preset конфигурации

### RateLimitDefault

```go
RequestsPerWindow: 100
WindowDuration: 1 minute
```

**Использование:** Стандартный лимит для большинства endpoints

**Кого затронет:**
- 100 запросов/минута = ~1.67 запросов/секунду
- Достаточно для обычного использования
- Блокирует только aggressive abuse

### RateLimitStrict

```go
RequestsPerWindow: 10
WindowDuration: 1 minute
```

**Использование:** Критические операции (email, платежи, пароли)

**Кого затронет:**
- 10 запросов/минута = ~1 запрос/6 секунд
- Защищает от brute-force атак
- Лимитирует дорогие операции

### RateLimitRelaxed

```go
RequestsPerWindow: 300
WindowDuration: 1 minute
```

**Использование:** Read-only операции, публичные endpoints

**Кого затронет:**
- 300 запросов/минута = 5 запросов/секунду
- Подходит для polling
- Не мешает легитимному использованию

## Конфигурация по endpoint

### Глобальный лимит

В router применяется глобальный default лимит:

```go
middleware.RateLimit(middleware.RateLimitDefault)
```

Применяется ко всем endpoints.

### Кастомные лимиты

Для конкретных endpoint групп можно создать отдельные middleware:

```go
// В router.go
strictLimiter := middleware.RateLimit(middleware.RateLimitStrict)

// Применяем к конкретным handlers
mux.Handle("/email/verify/send", strictLimiter(http.HandlerFunc(emailHandler.SendVerification)))
mux.Handle("/subscription/subscribe", strictLimiter(http.HandlerFunc(subscriptionHandler.Subscribe)))
```

### Пример: разные лимиты для разных операций

```go
func New(deps *Dependencies) http.Handler {
    mux := http.NewServeMux()

    // Создаем лимитеры
    defaultLimiter := middleware.RateLimit(middleware.RateLimitDefault)
    strictLimiter := middleware.RateLimit(middleware.RateLimitStrict)
    relaxedLimiter := middleware.RateLimit(middleware.RateLimitRelaxed)

    // Read-only endpoints - relaxed
    mux.Handle("GET /profile", relaxedLimiter(http.HandlerFunc(profileHandler.Get)))
    mux.Handle("GET /achievements", relaxedLimiter(http.HandlerFunc(achievementHandler.List)))

    // Write operations - default
    mux.Handle("POST /attempts", defaultLimiter(http.HandlerFunc(attemptHandler.Create)))
    mux.Handle("PUT /profile", defaultLimiter(http.HandlerFunc(profileHandler.Update)))

    // Critical operations - strict
    mux.Handle("POST /email/verify/send", strictLimiter(http.HandlerFunc(emailHandler.SendVerification)))
    mux.Handle("POST /subscription/subscribe", strictLimiter(http.HandlerFunc(subscriptionHandler.Subscribe)))

    // Глобальный chain без rate limit (применяется через handlers выше)
    return middleware.Chain(
        middleware.Recovery,
        middleware.Logging,
        // ... другие middleware
    )(mux)
}
```

## Client identification

### Per-IP limiting

Rate limit применяется **per-IP address**:

```go
clientID := getClientIP(r)
```

Один IP = один лимит.

### IP extraction logic

1. **X-Forwarded-For** (приоритет 1)
   - Используется если присутствует
   - Первый IP в списке = реальный client IP
   - Полезно за proxy/load balancer

2. **X-Real-IP** (приоритет 2)
   - Если X-Forwarded-For отсутствует
   - Используется некоторыми proxy

3. **RemoteAddr** (fallback)
   - Прямое подключение
   - Без proxy

### Пример headers

```http
# За Nginx reverse proxy:
GET /api/v1/profile HTTP/1.1
Host: api.example.com
X-Forwarded-For: 203.0.113.1
X-Real-IP: 203.0.113.1

# Прямое подключение:
GET /api/v1/profile HTTP/1.1
Host: api.example.com
RemoteAddr: 203.0.113.1:54321
```

### Shared IP issues

**Проблема:** Пользователи за NAT/корпоративным proxy имеют один IP

**Решение:**
- Увеличить лимиты для production
- Рассмотреть per-user limiting (требует auth)
- Использовать комбинацию IP + User-Agent

## Response headers

При каждом запросе API возвращает rate limit headers:

```http
HTTP/1.1 200 OK
X-RateLimit-Limit: 100
```

При превышении лимита:

```http
HTTP/1.1 429 Too Many Requests
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
Retry-After: 60

{"error": "Rate limit exceeded. Please try again later."}
```

### Header spec

| Header | Описание | Пример |
|--------|----------|--------|
| `X-RateLimit-Limit` | Максимум запросов в окне | `100` |
| `X-RateLimit-Remaining` | Оставшиеся запросы (только при 429) | `0` |
| `Retry-After` | Секунд до сброса лимита | `60` |

## Frontend integration

### Обработка 429 ответов

```typescript
async function apiRequest(url: string, options?: RequestInit) {
  const response = await fetch(url, options)

  if (response.status === 429) {
    // Получаем Retry-After header
    const retryAfter = response.headers.get('Retry-After')
    const seconds = retryAfter ? parseInt(retryAfter) : 60

    throw new RateLimitError(`Too many requests. Retry after ${seconds}s`, seconds)
  }

  return response
}

class RateLimitError extends Error {
  constructor(message: string, public retryAfter: number) {
    super(message)
    this.name = 'RateLimitError'
  }
}
```

### Автоматический retry с backoff

```typescript
async function fetchWithRetry(url: string, options?: RequestInit, maxRetries = 3) {
  for (let i = 0; i < maxRetries; i++) {
    try {
      return await apiRequest(url, options)
    } catch (error) {
      if (error instanceof RateLimitError) {
        if (i === maxRetries - 1) {
          throw error // Последняя попытка
        }

        // Ждем указанное время + небольшой jitter
        const delay = error.retryAfter * 1000 + Math.random() * 1000
        await new Promise(resolve => setTimeout(resolve, delay))
        continue
      }

      throw error
    }
  }
}
```

### UI индикация

```typescript
function RateLimitWarning({ error }: { error: RateLimitError }) {
  const [countdown, setCountdown] = useState(error.retryAfter)

  useEffect(() => {
    const timer = setInterval(() => {
      setCountdown(prev => {
        if (prev <= 1) {
          clearInterval(timer)
          return 0
        }
        return prev - 1
      })
    }, 1000)

    return () => clearInterval(timer)
  }, [])

  return (
    <Alert variant="warning">
      Слишком много запросов. Повторите через {countdown} секунд.
    </Alert>
  )
}
```

## Development vs Production

### Development (`ENV=development`)

- Rate limiting **отключен**
- Все запросы проходят без ограничений
- Удобно для локальной разработки

```bash
ENV=development go run cmd/api/main.go
```

### Production (`ENV=production`)

- Rate limiting **активен**
- Лимиты применяются согласно конфигурации
- Логируются все превышения лимита

```bash
ENV=production go run cmd/api/main.go
```

## Monitoring

### Логи

```
# Превышение лимита:
[RateLimit] Rate limit exceeded for 203.0.113.1 on POST /api/v1/attempts

# Development mode:
(логов нет - middleware пропускает все)
```

### Метрики для отслеживания

1. **Rate limit hits**
   - Количество 429 ответов
   - По endpoint
   - По IP

2. **Top offenders**
   - IP с наибольшим количеством блокировок
   - Возможные боты/abuse

3. **Legitimate blocks**
   - Реальные пользователи которых затронуло
   - Признак что лимиты слишком строгие

### Пример мониторинга

```go
// TODO: Интеграция с Prometheus/StatsD
var (
    rateLimitHitsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "rate_limit_hits_total",
            Help: "Total number of rate limit hits",
        },
        []string{"endpoint", "client_ip"},
    )
)
```

## Performance

### Memory usage

Каждый уникальный IP хранит:
- Slice запросов в текущем окне
- Last access time
- Mutex для concurrency

**Оценка:**
- 1000 активных IP
- Средний 50 запросов/окно
- ~200KB памяти

### Cleanup

Неактивные entries автоматически удаляются:

```go
// Каждые 5 минут
cleanupInterval := 5 * time.Minute

// Удаляются entries старше 10 минут
cleanupThreshold := 10 * time.Minute
```

### Concurrency

- `sync.RWMutex` для entries map
- Per-entry `sync.Mutex` для requests slice
- Минимальная contention
- Safe для concurrent requests

## Testing

### Unit tests

```bash
cd api
go test -v ./internal/api/middleware/ratelimit_test.go ./internal/api/middleware/ratelimit.go
```

Покрытие:
- ✅ Requests within limit pass
- ✅ Requests exceeding limit blocked
- ✅ Sliding window correctly enforced
- ✅ Different clients have separate limits
- ✅ Development mode bypasses limiting
- ✅ X-Forwarded-For header respected
- ✅ Remaining requests calculation

### Manual testing

```bash
# 1. Запустить в production режиме
ENV=production go run cmd/api/main.go

# 2. Отправить запросы в цикле (достичь лимита)
for i in {1..101}; do
  curl -v http://localhost:8080/api/v1/profile \
    -H "X-Child-Profile-ID: test" \
    -H "X-Platform-ID: vk"
  echo "Request $i"
done

# 3. Проверить что 101-й запрос получил 429
# HTTP/1.1 429 Too Many Requests
# Retry-After: 60
```

### Load testing

```bash
# С помощью Apache Bench
ab -n 1000 -c 10 http://localhost:8080/api/v1/profile

# С помощью wrk
wrk -t12 -c400 -d30s http://localhost:8080/api/v1/profile
```

## Troubleshooting

### 429 Too Many Requests для легитимных пользователей

**Проблема:** Реальные пользователи получают rate limit

**Решения:**
1. Увеличить лимиты:
   ```go
   RateLimitDefault = RateLimitConfig{
       RequestsPerWindow: 200, // было 100
       WindowDuration:    time.Minute,
   }
   ```

2. Использовать per-user limiting (требует auth)

3. Whitelist для известных IPs:
   ```go
   func isWhitelisted(ip string) bool {
       whitelist := []string{"10.0.0.1", "192.168.1.1"}
       // ...
   }
   ```

### Shared IP проблема (NAT, корпоративный proxy)

**Проблема:** Множество пользователей за одним IP

**Решения:**
1. Комбинированный ключ: IP + User ID
2. Увеличить лимиты для production
3. Использовать fingerprinting (User-Agent, Accept headers)

### Memory leak

**Проблема:** Entries не удаляются, память растет

**Диагностика:**
```go
// Добавить логирование в cleanup
log.Printf("[RateLimit] Cleanup: %d entries", len(rl.entries))
```

**Решение:**
- Уменьшить cleanup threshold
- Увеличить cleanup frequency
- Проверить что cleanup горутина работает

### Rate limit не применяется в dev

**Проблема:** Даже в production режиме лимиты не работают

**Проверка:**
```bash
# Убедитесь что ENV=production
echo $ENV
grep ENV .env

# Перезапустите backend
docker compose restart backend
```

## Best practices

### 1. Разные лимиты для разных операций

```go
// ❌ Плохо: один лимит для всех
middleware.RateLimit(middleware.RateLimitDefault)

// ✅ Хорошо: дифференцированные лимиты
// Read operations - relaxed
// Write operations - default
// Critical operations - strict
```

### 2. Graceful degradation

```go
// В dev режиме не блокировать
if isDev {
    next.ServeHTTP(w, r)
    return
}
```

### 3. Информативные error messages

```go
// ✅ Хорошо
http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)

// ❌ Плохо
http.Error(w, "Error", http.StatusTooManyRequests)
```

### 4. Retry-After header

```go
// Всегда указывайте когда можно повторить
w.Header().Set("Retry-After", fmt.Sprintf("%d", seconds))
```

### 5. Monitoring и алерты

- Отслеживайте rate limit hits
- Алерты при резком росте блокировок
- Проверяйте легитимность blocked IPs

## Security considerations

### IP spoofing

**Угроза:** Злоумышленник может подделать X-Forwarded-For header

**Защита:**
1. Доверяйте X-Forwarded-For только от известных proxy
2. Validate IP format
3. Используйте middleware для sanitization

```go
// В Nginx
proxy_set_header X-Forwarded-For $remote_addr;
```

### Distributed attacks

**Угроза:** Атака с множества IP адресов (distributed DDoS)

**Защита:**
- Per-IP limiting не поможет
- Используйте cloud-based DDoS protection
- Rate limiting на уровне endpoint (global)

### Bypass через reset

**Угроза:** Злоумышленник меняет IP чтобы сбросить лимит

**Защита:**
- Используйте per-user limiting для авторизованных запросов
- Fingerprinting (User-Agent, TLS fingerprint)
- Требуйте auth для критических операций

## Future improvements

### Redis-based distributed limiting

Для multi-instance deployment:

```go
import "github.com/go-redis/redis"

type RedisRateLimiter struct {
    client *redis.Client
}

func (rl *RedisRateLimiter) Allow(clientID string, config RateLimitConfig) bool {
    key := fmt.Sprintf("ratelimit:%s", clientID)
    // INCR key
    // EXPIRE key config.WindowDuration
    // ...
}
```

### Per-user limiting

Для авторизованных пользователей:

```go
// Вместо IP используем user ID
clientID := getUserID(r)
if clientID == "" {
    clientID = getClientIP(r) // fallback to IP
}
```

### Dynamic limits

Адаптивные лимиты based on load:

```go
func getDynamicLimit() int {
    load := getSystemLoad()
    if load > 0.8 {
        return 50 // Снижаем при высокой нагрузке
    }
    return 100
}
```

### Whitelist/Blacklist

```go
var (
    whitelist = []string{"10.0.0.0/8", "192.168.0.0/16"}
    blacklist = []string{"1.2.3.4", "5.6.7.8"}
)
```

## Links

- [OWASP API Security](https://owasp.org/www-project-api-security/)
- [Rate Limiting Strategies](https://cloud.google.com/architecture/rate-limiting-strategies-techniques)
- [Sliding Window Algorithm](https://en.wikipedia.org/wiki/Sliding_window_protocol)

## Checklist

- [ ] Rate limiting включен в production
- [ ] Лимиты настроены адекватно
- [ ] Frontend обрабатывает 429 ответы
- [ ] Retry-After header устанавливается
- [ ] Логи мониторятся на rate limit hits
- [ ] Dev режим отключает limiting
- [ ] X-Forwarded-For обрабатывается корректно
- [ ] Cleanup горутина работает
