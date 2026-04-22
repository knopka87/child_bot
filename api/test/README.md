# Testing Documentation

Полная документация по тестированию REST API backend.

## Содержание

- [Обзор](#обзор)
- [Уровни тестирования](#уровни-тестирования)
- [Настройка тестового окружения](#настройка-тестового-окружения)
- [Запуск тестов](#запуск-тестов)
- [Структура тестов](#структура-тестов)
- [Best Practices](#best-practices)

---

## Обзор

Проект использует три уровня тестирования:

1. **Unit Tests** - быстрые изолированные тесты с моками
2. **Integration Tests** - тесты с реальной БД
3. **E2E Tests** - полные end-to-end тесты через HTTP

**Покрытие:**
- ✅ Handlers (8 endpoints)
- ✅ Middleware (Auth, CORS, Recovery, Logging)
- ✅ Service Layer (AttemptService с LLM интеграцией)
- ✅ REST API E2E (с моками и реальным LLM)

---

## Уровни тестирования

### 1. Unit Tests (Юнит-тесты)

**Цель:** Проверить логику отдельных компонентов изолированно.

**Особенности:**
- Используют моки для зависимостей
- Очень быстрые (< 1 секунда)
- Не требуют внешних зависимостей (БД, LLM)
- Запускаются с флагом `-short`

**Расположение:**
```
api/internal/api/handler/*_test.go
api/internal/api/middleware/middleware_test.go
```

**Примеры:**
```go
// handler_test.go
func TestAttemptHandler_Create(t *testing.T) {
    mockService := &mockAttemptService{...}
    handler := NewAttemptHandler(mockService)
    // ... test logic
}
```

**Запуск:**
```bash
# Все юнит-тесты
make test

# Только handlers
make test-handlers

# Только middleware
make test-middleware

# С покрытием
make test-cover
```

---

### 2. Integration Tests (Интеграционные тесты)

**Цель:** Проверить взаимодействие компонентов с реальной БД.

**Особенности:**
- Используют реальную тестовую БД
- Средняя скорость (2-10 секунд)
- Требуют `TEST_DATABASE_URL`
- Автоматическая очистка данных через `t.Cleanup()`

**Расположение:**
```
api/internal/service/*_test.go
```

**Примеры:**
```go
// attempt_test.go
func TestAttemptService_ProcessHelp_Integration(t *testing.T) {
    db := setupTestDB(t)
    st := setupTestStore(t)
    mockLLM := &mockLLMClient{...}
    service := NewAttemptService(st, mockLLM, "gpt-4")
    // ... test with real DB
}
```

**Запуск:**
```bash
# Настройка test DB
make db-test-create
make migrate-test-up

# Запуск интеграционных тестов
make test-integration
```

---

### 3. E2E Tests (End-to-End тесты)

**Цель:** Проверить полный flow через HTTP API.

**Особенности:**
- Используют `httptest.Server`
- Полный stack: Router → Handler → Service → Store → DB
- Можно использовать mock или real LLM
- Долгие (10-30 минут с real LLM)

**Расположение:**
```
api/test/e2e/rest_api_test.go   # REST API E2E
api/test/e2e/hint_flow_test.go  # Telegram E2E (legacy)
```

**Примеры:**
```go
// rest_api_test.go
func TestE2E_AttemptFlow_Help(t *testing.T) {
    server, db := setupE2EServer(t, cfg)

    // Create attempt via HTTP
    resp := makeE2ERequest(t, server, "POST", "/attempts", createReq, ...)

    // Upload image
    // Process
    // Get hints
    // ...
}
```

**Запуск:**
```bash
# Быстрый E2E (с mock LLM)
make test-e2e

# С реальным LLM (требует LLM proxy)
make test-e2e-rest-real

# Все E2E тесты (включая legacy Telegram)
make test-e2e-all
```

---

## Настройка тестового окружения

### 1. Создать тестовую БД

```bash
make db-test-create
```

Это создаст базу `childbot_test` в PostgreSQL.

### 2. Применить миграции

```bash
make migrate-test-up
```

Все миграции из `api/migrations/` будут применены к test DB.

### 3. Проверить статус

```bash
# Проверить подключение
docker compose exec db psql -U childbot -d childbot_test -c "SELECT 1"

# Проверить таблицы
docker compose exec db psql -U childbot -d childbot_test -c "\dt"
```

### 4. Environment Variables

Создайте `.env.test` (опционально):

```bash
TEST_DATABASE_URL=postgres://childbot:root@localhost:5432/childbot_test?sslmode=disable
LLM_PROXY_URL=http://localhost:8081
USE_REAL_LLM=false
```

---

## Запуск тестов

### Быстрая проверка (Unit tests)

```bash
# Все unit тесты
make test

# С race detector (важно!)
make test-race

# С покрытием
make test-cover
open api/coverage.html
```

### Интеграционные тесты

```bash
# Настроить окружение (один раз)
make test-e2e-setup

# Запустить integration tests
make test-integration
```

### E2E тесты

```bash
# Быстрый E2E (mock LLM, 1-2 минуты)
make test-e2e

# Полный E2E с реальным LLM (30+ минут)
make test-e2e-rest-real
```

### Все тесты

```bash
# Unit + Integration + E2E (fast)
make test-all

# Для CI: unit + race + coverage
make test-ci
```

---

## Структура тестов

```
api/
├── internal/
│   ├── api/
│   │   ├── handler/
│   │   │   ├── attempt.go
│   │   │   ├── attempt_test.go           # Unit tests
│   │   │   ├── handler_test.go           # Test helpers + mocks
│   │   │   └── ...
│   │   └── middleware/
│   │       ├── auth.go
│   │       └── middleware_test.go        # Unit tests
│   └── service/
│       ├── attempt.go
│       ├── attempt_test.go               # Integration tests
│       └── service_test.go               # Test helpers + mocks
└── test/
    └── e2e/
        ├── rest_api_test.go              # REST API E2E tests
        ├── hint_flow_test.go             # Telegram E2E (legacy)
        ├── check_answer_test.go          # Telegram E2E (legacy)
        ├── helpers.go
        ├── config.go
        └── testdata/
            ├── tasks/                     # Task images
            └── answers/                   # Answer images
```

---

## Best Practices

### 1. Используйте Table-Driven Tests

```go
func TestAttemptHandler_Create(t *testing.T) {
    tests := []struct {
        name           string
        requestBody    interface{}
        mockCreate     func(...) (string, error)
        expectedStatus int
    }{
        {name: "success", ...},
        {name: "validation error", ...},
        {name: "service error", ...},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test logic
        })
    }
}
```

### 2. Используйте t.Helper()

```go
func setupTestDB(t *testing.T) *sql.DB {
    t.Helper()  // Правильные номера строк в ошибках
    // ...
}
```

### 3. Всегда делайте Cleanup

```go
func createTestProfile(t *testing.T, db *sql.DB) string {
    // ...
    t.Cleanup(func() {
        db.Exec("DELETE FROM child_profiles WHERE id = $1", profileID)
    })
    return profileID
}
```

### 4. Используйте -short для быстрых тестов

```go
func TestIntegration_Something(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    // ... slow test with DB
}
```

### 5. Mock только на границах

❌ **Плохо:** Мокировать все внутри service
```go
// Don't mock internal implementation details
```

✅ **Хорошо:** Мокировать внешние зависимости (LLM, external APIs)
```go
mockLLM := &mockLLMClient{
    detectFunc: func(...) { return mockResponse, nil }
}
service := NewAttemptService(realStore, mockLLM, "gpt-4")
```

### 6. Используйте Race Detector

```bash
# Всегда запускайте с -race в CI
make test-race
```

### 7. Проверяйте покрытие

```bash
make test-cover
open api/coverage.html
```

**Цель:** 80%+ покрытие для критического кода (handlers, service layer).

---

## Troubleshooting

### Ошибка: "TEST_DATABASE_URL not set"

```bash
# Убедитесь, что переменная установлена
export TEST_DATABASE_URL="postgres://childbot:root@localhost:5432/childbot_test?sslmode=disable"

# Или используйте Makefile (читает .env)
make test-integration
```

### Ошибка: "connection refused" к БД

```bash
# Проверьте, что PostgreSQL запущен
docker compose ps

# Создайте test DB
make db-test-create

# Примените миграции
make migrate-test-up
```

### Ошибка: "LLM proxy not available"

```bash
# Для E2E с mock LLM (не нужен proxy)
make test-e2e

# Для E2E с real LLM (нужен proxy)
export LLM_PROXY_URL="http://localhost:8081"
make test-e2e-rest-real
```

### Тесты зависают

```bash
# Установите timeout
go test -timeout 5m ./...

# Или используйте Makefile с таймаутами
make test-integration  # 5m timeout
make test-e2e          # 10m timeout
```

### Race condition обнаружен

```bash
# Найдите race condition
go test -race ./...

# Типичные проблемы:
# - Goroutines пишут в shared state без sync.Mutex
# - Capture loop variable перед goroutine: v := v
```

---

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_USER: childbot
          POSTGRES_PASSWORD: root
          POSTGRES_DB: childbot_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432

    steps:
      - uses: actions/checkout@v3

      - uses: actions/setup-go@v4
        with:
          go-version: '1.24'

      - name: Install migrate
        run: make migrate-install

      - name: Run migrations
        run: make migrate-test-up
        env:
          TEST_DATABASE_URL: postgres://childbot:root@localhost:5432/childbot_test?sslmode=disable

      - name: Run tests
        run: make test-ci
        env:
          TEST_DATABASE_URL: postgres://childbot:root@localhost:5432/childbot_test?sslmode=disable

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./api/coverage.out
```

---

## Metrics

**Current Test Coverage (Phase 6):**

| Component           | Unit | Integration | E2E | Coverage |
|---------------------|------|-------------|-----|----------|
| Handlers            | ✅   | -           | ✅  | 85%      |
| Middleware          | ✅   | -           | ✅  | 90%      |
| Service Layer       | ✅   | ✅          | ✅  | 80%      |
| Router              | -    | -           | ✅  | 70%      |
| **Total**           | ✅   | ✅          | ✅  | **82%**  |

**Test Execution Time:**

- Unit tests: ~2 seconds
- Integration tests: ~10 seconds
- E2E tests (mock LLM): ~1-2 minutes
- E2E tests (real LLM): ~30 minutes

---

## Next Steps (Phase 7)

После завершения Phase 7 (Cleanup Telegram Code):

- [ ] Удалить `test/e2e/hint_flow_test.go`
- [ ] Удалить `test/e2e/check_answer_test.go`
- [ ] Удалить `internal/v2/telegram/*_test.go`
- [ ] Обновить Makefile (убрать `test-telegram`)
- [ ] Обновить CI/CD конфиг

---

## Полезные команды

```bash
# Проверка тестов
make test                    # Быстрая проверка
make test-race              # С race detector
make test-cover             # С покрытием
make test-all               # Все тесты

# Отдельные компоненты
make test-handlers          # Handler unit tests
make test-middleware        # Middleware unit tests
make test-service           # Service integration tests

# E2E тесты
make test-e2e-setup         # Подготовка (один раз)
make test-e2e               # REST API E2E (fast)
make test-e2e-rest-real     # REST API E2E (с LLM)

# CI
make test-ci                # Race + Coverage для CI
```

---

## Контакты

При проблемах с тестами:
1. Проверьте этот README
2. Проверьте `make help`
3. Проверьте logs: `docker compose logs -f`
4. Создайте issue в репозитории
