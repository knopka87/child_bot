# E2E Testing Guide

## Обзор

End-to-End тесты для проверки критических пользовательских flow перед релизом и модерацией VK.

## Технологии

- **Playwright** - framework для E2E тестирования
- **TypeScript** - типизированные тесты
- **Multi-browser** - Chrome, Firefox, Safari, Mobile

## Установка

```bash
cd frontend

# Установить зависимости
npm install

# Установить браузеры Playwright
npm run playwright:install
```

## Запуск тестов

### Все тесты

```bash
# Headless mode (CI)
npm run test:e2e

# Headed mode (видимый браузер)
npm run test:e2e:headed

# UI mode (интерактивный)
npm run test:e2e:ui

# Debug mode
npm run test:e2e:debug
```

### Отдельные тесты

```bash
# Только onboarding
npx playwright test 01-onboarding

# Только help flow
npx playwright test 02-help-flow

# Specific browser
npx playwright test --project=chromium
```

### Отчеты

```bash
# Открыть HTML отчет
npm run test:e2e:report
```

## Структура тестов

```
frontend/
├── e2e/
│   ├── fixtures/
│   │   └── auth.ts           # Fixture для авторизованного пользователя
│   └── critical/
│       ├── 01-onboarding.spec.ts      # Онбординг
│       ├── 02-help-flow.spec.ts       # Помощь с ДЗ
│       ├── 03-check-flow.spec.ts      # Проверка решения
│       ├── 04-achievements.spec.ts    # Достижения
│       ├── 05-villain-battle.spec.ts  # Битвы со злодеями
│       └── 06-navigation.spec.ts      # Навигация
├── playwright.config.ts      # Конфигурация Playwright
└── package.json
```

## Критические Flow

### 1. Onboarding Flow

**Файл:** `e2e/critical/01-onboarding.spec.ts`

**Что тестируется:**
- ✅ Полный процесс онбординга (Welcome → Name → Grade → Avatar → Email → Consent → Home)
- ✅ Валидация обязательных полей
- ✅ Навигация назад с сохранением состояния
- ✅ Переход на главную после завершения

**Тесты:**
```typescript
test('should complete full onboarding flow')
test('should validate required fields')
test('should allow navigation back')
```

### 2. Help Flow

**Файл:** `e2e/critical/02-help-flow.spec.ts`

**Что тестируется:**
- ✅ Загрузка фото задания
- ✅ Обработка через AI
- ✅ Получение подсказок
- ✅ Продолжение незавершенной попытки
- ✅ Обработка ошибок

**Тесты:**
```typescript
test('should upload task and get hints')
test('should continue unfinished attempt')
test('should show error on upload failure')
```

### 3. Check Flow

**Файл:** `e2e/critical/03-check-flow.spec.ts`

**Что тестируется:**
- ✅ Выбор сценария (одно/два фото)
- ✅ Загрузка задания и решения
- ✅ Проверка решения
- ✅ Валидация файлов

**Тесты:**
```typescript
test('should select scenario and upload solution')
test('should upload two photos for comparison')
test('should show validation errors')
```

### 4. Achievements Flow

**Файл:** `e2e/critical/04-achievements.spec.ts`

**Что тестируется:**
- ✅ Отображение достижений
- ✅ Разблокированные vs заблокированные
- ✅ Прогресс серийных достижений
- ✅ Модальное окно с деталями

**Тесты:**
```typescript
test('should display achievements page')
test('should show locked achievements')
test('should display progress for serial achievements')
```

### 5. Villain Battle Flow

**Файл:** `e2e/critical/05-villain-battle.spec.ts`

**Что тестируется:**
- ✅ Отображение активного злодея
- ✅ Начало битвы
- ✅ Нанесение урона
- ✅ Экран победы

**Тесты:**
```typescript
test('should display active villain')
test('should start battle and damage villain')
test('should show victory screen when villain defeated')
test('should apply damage after correct solution')
```

### 6. Navigation Flow

**Файл:** `e2e/critical/06-navigation.spec.ts`

**Что тестируется:**
- ✅ Переключение между табами
- ✅ Навигация внутри разделов
- ✅ Back button
- ✅ Browser back/forward
- ✅ Сохранение состояния

**Тесты:**
```typescript
test('should navigate between main tabs')
test('should navigate to history from profile')
test('should navigate back from deep pages')
test('should preserve state on tab switch')
test('should handle browser back/forward')
```

## Fixtures

### Auth Fixture

**Файл:** `e2e/fixtures/auth.ts`

Автоматически устанавливает storage state для авторизованного пользователя:

```typescript
import { test, expect } from '../fixtures/auth';

test('my authenticated test', async ({ page }) => {
  // Пользователь уже авторизован
  // child_profile_id = 'test-profile-id-123'
  // platform_id = 'vk'
});
```

**Mock VK Bridge:**
- `VKWebAppStorageGet` - возвращает test profile
- `VKWebAppStorageSet` - сохраняет в localStorage
- `VKWebAppInit` - инициализация

## Конфигурация

### playwright.config.ts

```typescript
export default defineConfig({
  testDir: './e2e',
  timeout: 30 * 1000,

  use: {
    baseURL: 'http://localhost:5173',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
  },

  projects: [
    { name: 'chromium' },
    { name: 'firefox' },
    { name: 'webkit' },
    { name: 'Mobile Chrome' },
    { name: 'Mobile Safari' },
  ],

  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:5173',
    reuseExistingServer: !process.env.CI,
  },
});
```

**Ключевые настройки:**
- **timeout**: 30 секунд на тест
- **fullyParallel**: параллельное выполнение
- **retries**: 2 попытки в CI
- **webServer**: автозапуск dev сервера

## Best Practices

### 1. Используйте data-testid

```tsx
// В компоненте
<button data-testid="help-button">Помочь</button>

// В тесте
await page.click('[data-testid="help-button"]');
```

### 2. Ждите элементы явно

```typescript
// ❌ Плохо
await page.click('button');

// ✅ Хорошо
await expect(page.locator('button:has-text("Далее")')).toBeVisible();
await page.click('button:has-text("Далее")');
```

### 3. Mock API для стабильности

```typescript
await page.route('**/api/v1/home/*', async (route) => {
  await route.fulfill({
    status: 200,
    body: JSON.stringify({ ... }),
  });
});
```

### 4. Проверяйте URL transitions

```typescript
await page.click('button:has-text("Далее")');
await expect(page).toHaveURL(/\/next-step/);
```

### 5. Используйте custom fixtures

```typescript
// Создайте fixture для часто используемого состояния
export const test = base.extend({
  authenticatedPage: async ({ page }, use) => {
    // Setup
    await page.goto('/');
    await performAuth(page);

    await use(page);

    // Cleanup
  },
});
```

## CI Integration

### GitHub Actions

```yaml
name: E2E Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Node
        uses: actions/setup-node@v3
        with:
          node-version: '18'

      - name: Install dependencies
        run: |
          cd frontend
          npm ci

      - name: Install Playwright
        run: |
          cd frontend
          npx playwright install --with-deps chromium

      - name: Run E2E tests
        run: |
          cd frontend
          npm run test:e2e

      - name: Upload test results
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: playwright-report
          path: frontend/playwright-report/
```

## Debugging

### UI Mode

Лучший способ для разработки тестов:

```bash
npm run test:e2e:ui
```

**Возможности:**
- 👁️ Визуальный просмотр шагов
- ⏸️ Pause на любом шаге
- 🔍 Inspect элементов
- 📸 Скриншоты каждого шага
- 📹 Video replay

### Debug Mode

Пошаговое выполнение:

```bash
npm run test:e2e:debug
```

**Или в коде:**

```typescript
test('debug example', async ({ page }) => {
  await page.pause(); // Breakpoint
  await page.click('button');
});
```

### Traces

Записывайте trace для воспроизведения:

```typescript
// playwright.config.ts
use: {
  trace: 'on', // Always record
}
```

Просмотр:
```bash
npx playwright show-trace trace.zip
```

## Troubleshooting

### Тесты падают с timeout

**Проблема:** `Test timeout of 30000ms exceeded`

**Решения:**
1. Увеличьте timeout:
```typescript
test('slow test', async ({ page }) => {
  test.setTimeout(60000);
  // ...
});
```

2. Проверьте что dev server запущен:
```bash
npm run dev
```

3. Используйте waitForURL вместо sleep:
```typescript
// ❌ Плохо
await page.waitForTimeout(5000);

// ✅ Хорошо
await page.waitForURL(/\/result/);
```

### Элемент не найден

**Проблема:** `locator.click: Target closed`

**Решения:**
1. Ждите загрузку:
```typescript
await page.waitForLoadState('networkidle');
```

2. Проверьте selector:
```typescript
// Debug selector
const element = page.locator('button');
console.log(await element.count()); // 0 = not found
```

3. Используйте более специфичный selector:
```typescript
// ❌ Слишком общий
await page.click('button');

// ✅ Специфичный
await page.click('[data-testid="submit-button"]');
```

### Mock не работает

**Проблема:** API запрос проходит мимо mock

**Решение:**
```typescript
// Установите route ДО навигации
await page.route('**/api/**', route => route.fulfill(...));
await page.goto('/');
```

### Screenshots не сохраняются

**Проблема:** Нет скриншотов в test-results

**Решение:**
```typescript
// Принудительный screenshot
await page.screenshot({ path: 'debug.png' });

// В config
use: {
  screenshot: 'on', // Всегда
}
```

## Performance

### Оптимизация тестов

**1. Parallel execution**
```typescript
test.describe.configure({ mode: 'parallel' });
```

**2. Reuse context**
```typescript
// Для быстрых тестов
test.use({ storageState: 'auth.json' });
```

**3. Skip unnecessary waits**
```typescript
// ❌ Медленно
await page.waitForTimeout(3000);

// ✅ Быстро
await page.waitForSelector('[data-testid="element"]');
```

**4. Disable animations**
```typescript
await page.addStyleTag({
  content: '* { animation: none !important; transition: none !important; }'
});
```

### Метрики

Типичное время выполнения:

| Test Suite | Tests | Duration |
|------------|-------|----------|
| Onboarding | 3 | ~45s |
| Help Flow | 3 | ~30s |
| Check Flow | 3 | ~30s |
| Achievements | 3 | ~20s |
| Villain Battle | 4 | ~35s |
| Navigation | 5 | ~25s |
| **Total** | **21** | **~3min** |

## Coverage

### Покрытие критических path

- ✅ **Onboarding**: 100% (все шаги)
- ✅ **Help**: 90% (file upload mock)
- ✅ **Check**: 90% (file upload mock)
- ✅ **Achievements**: 100%
- ✅ **Villain**: 95% (damage calculation)
- ✅ **Navigation**: 100%

### Что НЕ тестируется E2E

- ❌ Unit логика компонентов (используйте Vitest)
- ❌ API backend (используйте Go тесты)
- ❌ Стилизация (используйте visual regression)
- ❌ Accessibility (используйте axe-core)

## Maintenance

### Обновление тестов

После изменения UI:

1. Запустите тесты:
```bash
npm run test:e2e:ui
```

2. Обновите selectors где нужно

3. Обновите ожидаемые тексты

4. Проверьте на всех браузерах

### Добавление новых тестов

```typescript
// e2e/critical/07-new-feature.spec.ts
import { test, expect } from '../fixtures/auth';

test.describe('New Feature Flow', () => {
  test('should work correctly', async ({ page }) => {
    await page.goto('/new-feature');
    // ...
  });
});
```

## Resources

- [Playwright Documentation](https://playwright.dev)
- [Best Practices](https://playwright.dev/docs/best-practices)
- [Debugging Guide](https://playwright.dev/docs/debug)
- [CI/CD Integration](https://playwright.dev/docs/ci)
