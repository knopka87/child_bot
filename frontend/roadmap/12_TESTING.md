# Phase 11: Тестирование (Testing Strategy)

**Длительность:** 4-5 дней
**Приоритет:** Критический
**Зависимости:** Все предыдущие фазы

---

## Цель

Создать полную стратегию тестирования: unit тесты (Vitest + RTL), integration тесты, E2E тесты (Playwright) для критических потоков, тестирование на реальных устройствах, performance и accessibility testing.

---

## Архитектура тестирования

### Пирамида тестирования

```
       E2E Tests (Playwright)
      /                      \
   Integration Tests
  /                          \
Unit Tests (Vitest + RTL)
```

### Структура тестов

```
tests/
├── unit/
│   ├── components/
│   ├── hooks/
│   ├── utils/
│   └── services/
├── integration/
│   ├── flows/
│   └── api/
├── e2e/
│   ├── onboarding.spec.ts
│   ├── help-flow.spec.ts
│   ├── check-flow.spec.ts
│   └── paywall.spec.ts
└── performance/
    └── lighthouse.config.js
```

---

## Часть 1: Настройка инструментов

### 1.1. Vitest Configuration

**Файл:** `vitest.config.ts`

```typescript
import { defineConfig } from 'vitest/config';
import react from '@vitejs/plugin-react';
import path from 'path';

export default defineConfig({
  plugins: [react()],
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./tests/setup.ts'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      include: ['src/**/*.{ts,tsx}'],
      exclude: [
        'src/**/*.test.{ts,tsx}',
        'src/**/*.spec.{ts,tsx}',
        'src/types/**',
        'src/**/*.d.ts',
      ],
      thresholds: {
        lines: 80,
        functions: 80,
        branches: 75,
        statements: 80,
      },
    },
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
});
```

---

### 1.2. Test Setup

**Файл:** `tests/setup.ts`

```typescript
import '@testing-library/jest-dom';
import { cleanup } from '@testing-library/react';
import { afterEach, vi } from 'vitest';

// Cleanup after each test
afterEach(() => {
  cleanup();
});

// Mock window.matchMedia
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation((query) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
});

// Mock IntersectionObserver
global.IntersectionObserver = class IntersectionObserver {
  constructor() {}
  disconnect() {}
  observe() {}
  takeRecords() {
    return [];
  }
  unobserve() {}
} as any;

// Mock navigator.clipboard
Object.assign(navigator, {
  clipboard: {
    writeText: vi.fn(() => Promise.resolve()),
    readText: vi.fn(() => Promise.resolve('')),
  },
});
```

---

### 1.3. Playwright Configuration

**Файл:** `playwright.config.ts`

```typescript
import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './tests/e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: [
    ['html'],
    ['json', { outputFile: 'playwright-report/results.json' }],
  ],
  use: {
    baseURL: 'http://localhost:5173',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
  },

  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },
    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
    },
    {
      name: 'Mobile Chrome',
      use: { ...devices['Pixel 5'] },
    },
    {
      name: 'Mobile Safari',
      use: { ...devices['iPhone 12'] },
    },
  ],

  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:5173',
    reuseExistingServer: !process.env.CI,
  },
});
```

---

## Часть 2: Unit Tests

### 2.1. Component Tests

**Файл:** `src/components/ui/Button/Button.test.tsx`

```typescript
import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { Button } from './Button';

describe('Button', () => {
  it('renders with text', () => {
    render(<Button>Click me</Button>);
    expect(screen.getByText('Click me')).toBeInTheDocument();
  });

  it('calls onClick when clicked', () => {
    const onClick = vi.fn();
    render(<Button onClick={onClick}>Click me</Button>);

    fireEvent.click(screen.getByText('Click me'));
    expect(onClick).toHaveBeenCalledTimes(1);
  });

  it('is disabled when disabled prop is true', () => {
    render(<Button disabled>Click me</Button>);
    const button = screen.getByText('Click me');

    expect(button).toBeDisabled();
  });

  it('shows loading spinner when isLoading is true', () => {
    render(<Button isLoading>Click me</Button>);
    expect(screen.getByRole('status')).toBeInTheDocument();
  });

  it('applies correct variant class', () => {
    const { rerender } = render(<Button variant="primary">Primary</Button>);
    expect(screen.getByText('Primary')).toHaveClass('primary');

    rerender(<Button variant="secondary">Secondary</Button>);
    expect(screen.getByText('Secondary')).toHaveClass('secondary');
  });

  it('applies fullWidth class when isFullWidth is true', () => {
    render(<Button isFullWidth>Full Width</Button>);
    expect(screen.getByText('Full Width')).toHaveClass('fullWidth');
  });
});
```

---

### 2.2. Hook Tests

**Файл:** `src/hooks/useAnalytics.test.ts`

```typescript
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useAnalytics } from './useAnalytics';
import { AnalyticsProvider } from '@/contexts/AnalyticsContext';

describe('useAnalytics', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('throws error when used outside AnalyticsProvider', () => {
    expect(() => renderHook(() => useAnalytics())).toThrow(
      'useAnalytics must be used within AnalyticsProvider'
    );
  });

  it('tracks event with correct parameters', () => {
    const { result } = renderHook(() => useAnalytics(), {
      wrapper: AnalyticsProvider,
    });

    act(() => {
      result.current.trackEvent('home_opened', {
        child_profile_id: 'test-123',
      });
    });

    // Verify event was tracked (would need to mock analytics service)
    expect(result.current.sessionId).toBeDefined();
  });

  it('sets user properties', () => {
    const { result } = renderHook(() => useAnalytics(), {
      wrapper: AnalyticsProvider,
    });

    act(() => {
      result.current.setUserProperties({
        grade: 5,
        level: 10,
      });
    });

    // Verify properties were set
  });
});
```

---

### 2.3. Utility Tests

**Файл:** `src/utils/validation.test.ts`

```typescript
import { describe, it, expect } from 'vitest';
import { validateEmail, validateDisplayName } from './validation';

describe('validateEmail', () => {
  it('validates correct email', () => {
    expect(validateEmail('test@example.com')).toBe(true);
    expect(validateEmail('user.name+tag@example.co.uk')).toBe(true);
  });

  it('rejects invalid email', () => {
    expect(validateEmail('invalid')).toBe(false);
    expect(validateEmail('test@')).toBe(false);
    expect(validateEmail('@example.com')).toBe(false);
    expect(validateEmail('test @example.com')).toBe(false);
  });

  it('rejects empty email', () => {
    expect(validateEmail('')).toBe(false);
  });
});

describe('validateDisplayName', () => {
  it('validates correct name', () => {
    expect(validateDisplayName('Артём')).toBe(true);
    expect(validateDisplayName('Anna-Maria')).toBe(true);
  });

  it('rejects too short name', () => {
    expect(validateDisplayName('A')).toBe(false);
  });

  it('rejects too long name', () => {
    expect(validateDisplayName('A'.repeat(21))).toBe(false);
  });

  it('rejects empty name', () => {
    expect(validateDisplayName('')).toBe(false);
    expect(validateDisplayName('   ')).toBe(false);
  });
});
```

---

## Часть 3: Integration Tests

### 3.1. Flow Integration Test

**Файл:** `tests/integration/onboarding-flow.test.tsx`

```typescript
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { OnboardingFlow } from '@/pages/Onboarding/OnboardingFlow';
import { onboardingAPI } from '@/api/onboarding';

vi.mock('@/api/onboarding');

describe('Onboarding Flow Integration', () => {
  beforeEach(() => {
    vi.mocked(onboardingAPI.getAvatars).mockResolvedValue([
      { id: '1', imageUrl: '/avatar1.png', name: 'Fox', isPremium: false },
      { id: '2', imageUrl: '/avatar2.png', name: 'Cat', isPremium: false },
    ]);
  });

  it('completes full onboarding flow', async () => {
    render(
      <BrowserRouter>
        <OnboardingFlow />
      </BrowserRouter>
    );

    // Step 1: Welcome screen
    expect(screen.getByText(/добро пожаловать/i)).toBeInTheDocument();
    fireEvent.click(screen.getByText(/продолжить/i));

    // Step 2: Grade selection
    await waitFor(() => {
      expect(screen.getByText(/в каком ты классе/i)).toBeInTheDocument();
    });
    fireEvent.click(screen.getByText('5'));

    // Step 3: Avatar selection
    await waitFor(() => {
      expect(screen.getByText(/выбери своего героя/i)).toBeInTheDocument();
    });
    const avatar = await screen.findByAltText('Fox');
    fireEvent.click(avatar);

    // Step 4: Display name
    await waitFor(() => {
      expect(screen.getByText(/как тебя зовут/i)).toBeInTheDocument();
    });
    const nameInput = screen.getByPlaceholderText(/твоё имя/i);
    fireEvent.change(nameInput, { target: { value: 'Артём' } });
    fireEvent.click(screen.getByText(/продолжить/i));

    // Step 5: Consent screen
    await waitFor(() => {
      expect(screen.getByText(/согласие/i)).toBeInTheDocument();
    });
    fireEvent.click(screen.getByLabelText(/являюсь родителем/i));
    fireEvent.click(screen.getByLabelText(/политику конфиденциальности/i));
    fireEvent.click(screen.getByLabelText(/пользовательское соглашение/i));
    fireEvent.click(screen.getByText(/продолжить/i));

    // Verify completion
    await waitFor(() => {
      expect(vi.mocked(onboardingAPI.completeOnboarding)).toHaveBeenCalled();
    });
  });

  it('validates display name input', async () => {
    render(
      <BrowserRouter>
        <OnboardingFlow />
      </BrowserRouter>
    );

    // Navigate to display name screen
    // ... (skip intermediate steps)

    const nameInput = screen.getByPlaceholderText(/твоё имя/i);
    const submitButton = screen.getByText(/продолжить/i);

    // Test too short
    fireEvent.change(nameInput, { target: { value: 'A' } });
    fireEvent.click(submitButton);
    expect(await screen.findByText(/слишком короткое/i)).toBeInTheDocument();

    // Test too long
    fireEvent.change(nameInput, { target: { value: 'A'.repeat(21) } });
    fireEvent.click(submitButton);
    expect(await screen.findByText(/слишком длинное/i)).toBeInTheDocument();

    // Test valid
    fireEvent.change(nameInput, { target: { value: 'Артём' } });
    fireEvent.click(submitButton);
    await waitFor(() => {
      expect(screen.queryByText(/ошибка/i)).not.toBeInTheDocument();
    });
  });
});
```

---

## Часть 4: E2E Tests

### 4.1. Onboarding E2E Test

**Файл:** `tests/e2e/onboarding.spec.ts`

```typescript
import { test, expect } from '@playwright/test';

test.describe('Onboarding Flow', () => {
  test('completes full onboarding', async ({ page }) => {
    await page.goto('/');

    // Welcome screen
    await expect(page.getByText(/добро пожаловать/i)).toBeVisible();
    await page.getByRole('button', { name: /продолжить/i }).click();

    // Grade selection
    await expect(page.getByText(/в каком ты классе/i)).toBeVisible();
    await page.getByText('5', { exact: true }).click();

    // Avatar selection
    await expect(page.getByText(/выбери своего героя/i)).toBeVisible();
    await page.locator('[alt="Fox"]').first().click();

    // Display name
    await expect(page.getByText(/как тебя зовут/i)).toBeVisible();
    await page.getByPlaceholder(/твоё имя/i).fill('Артём');
    await page.getByRole('button', { name: /продолжить/i }).click();

    // Consent screen
    await expect(page.getByText(/согласие/i)).toBeVisible();
    await page.getByLabel(/являюсь родителем/i).check();
    await page.getByLabel(/политику конфиденциальности/i).check();
    await page.getByLabel(/пользовательское соглашение/i).check();
    await page.getByRole('button', { name: /продолжить/i }).click();

    // Email input
    await expect(page.getByText(/email родителя/i)).toBeVisible();
    await page.getByPlaceholder(/example.com/i).fill('parent@example.com');
    await page.getByRole('button', { name: /отправить письмо/i }).click();

    // Should redirect to home after completion
    await expect(page).toHaveURL(/\/home/);
  });

  test('shows validation errors', async ({ page }) => {
    await page.goto('/onboarding');

    // Navigate to display name
    // ... (skip intermediate steps)

    // Test empty name
    await page.getByRole('button', { name: /продолжить/i }).click();
    await expect(page.getByText(/введи своё имя/i)).toBeVisible();

    // Test too short name
    await page.getByPlaceholder(/твоё имя/i).fill('A');
    await page.getByRole('button', { name: /продолжить/i }).click();
    await expect(page.getByText(/слишком короткое/i)).toBeVisible();
  });
});
```

---

### 4.2. Help Flow E2E Test

**Файл:** `tests/e2e/help-flow.spec.ts`

```typescript
import { test, expect } from '@playwright/test';
import path from 'path';

test.describe('Help Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Login and navigate to home
    await page.goto('/');
    // ... (authentication steps)
  });

  test('completes help flow with image upload', async ({ page }) => {
    // Click help button
    await page.getByRole('button', { name: /помоги разобраться/i }).click();

    // Upload screen
    await expect(page.getByText(/загрузи фото задания/i)).toBeVisible();

    // Upload image
    const fileInput = page.locator('input[type="file"]');
    const testImagePath = path.join(__dirname, '../fixtures/task-image.jpg');
    await fileInput.setInputFiles(testImagePath);

    // Quality check screen
    await expect(page.getByText(/всё ли видно/i)).toBeVisible();
    await page.getByRole('button', { name: /всё видно/i }).click();

    // Processing screen
    await expect(page.getByText(/обрабатываем/i)).toBeVisible();

    // Result screen (with timeout for processing)
    await expect(page.getByText(/результат/i)).toBeVisible({
      timeout: 30000,
    });

    // Should show hints
    await expect(page.getByText(/подсказка/i)).toBeVisible();
  });

  test('allows retaking photo in quality check', async ({ page }) => {
    await page.getByRole('button', { name: /помоги разобраться/i }).click();

    // Upload image
    const fileInput = page.locator('input[type="file"]');
    await fileInput.setInputFiles('../fixtures/task-image.jpg');

    // Click retake
    await page.getByRole('button', { name: /переснять/i }).click();

    // Should return to upload screen
    await expect(page.getByText(/выбрать изображение/i)).toBeVisible();
  });
});
```

---

## Часть 5: Performance Tests

### 5.1. Lighthouse Configuration

**Файл:** `tests/performance/lighthouse.config.js`

```javascript
module.exports = {
  ci: {
    collect: {
      url: [
        'http://localhost:5173/',
        'http://localhost:5173/home',
        'http://localhost:5173/achievements',
      ],
      numberOfRuns: 3,
      settings: {
        preset: 'desktop',
      },
    },
    assert: {
      assertions: {
        'categories:performance': ['error', { minScore: 0.9 }],
        'categories:accessibility': ['error', { minScore: 0.9 }],
        'categories:best-practices': ['error', { minScore: 0.9 }],
        'categories:seo': ['error', { minScore: 0.9 }],
        'first-contentful-paint': ['error', { maxNumericValue: 2000 }],
        'largest-contentful-paint': ['error', { maxNumericValue: 2500 }],
        'cumulative-layout-shift': ['error', { maxNumericValue: 0.1 }],
        'total-blocking-time': ['error', { maxNumericValue: 300 }],
      },
    },
    upload: {
      target: 'temporary-public-storage',
    },
  },
};
```

---

## Часть 6: Accessibility Tests

### 6.1. Accessibility Test

**Файл:** `tests/e2e/accessibility.spec.ts`

```typescript
import { test, expect } from '@playwright/test';
import AxeBuilder from '@axe-core/playwright';

test.describe('Accessibility', () => {
  test('home page has no accessibility violations', async ({ page }) => {
    await page.goto('/home');

    const accessibilityScanResults = await new AxeBuilder({ page }).analyze();

    expect(accessibilityScanResults.violations).toEqual([]);
  });

  test('onboarding has no accessibility violations', async ({ page }) => {
    await page.goto('/onboarding');

    const accessibilityScanResults = await new AxeBuilder({ page }).analyze();

    expect(accessibilityScanResults.violations).toEqual([]);
  });

  test('can navigate with keyboard only', async ({ page }) => {
    await page.goto('/home');

    // Tab through interactive elements
    await page.keyboard.press('Tab');
    await expect(page.locator(':focus')).toBeVisible();

    // Should be able to activate focused element
    await page.keyboard.press('Enter');
  });
});
```

---

## Часть 7: Test Utilities

### 7.1. Test Helpers

**Файл:** `tests/utils/test-utils.tsx`

```typescript
import { ReactElement } from 'react';
import { render, RenderOptions } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { AnalyticsProvider } from '@/contexts/AnalyticsContext';

interface AllTheProvidersProps {
  children: React.ReactNode;
}

function AllTheProviders({ children }: AllTheProvidersProps) {
  return (
    <BrowserRouter>
      <AnalyticsProvider config={{ enabled: false, debug: false }}>
        {children}
      </AnalyticsProvider>
    </BrowserRouter>
  );
}

function customRender(
  ui: ReactElement,
  options?: Omit<RenderOptions, 'wrapper'>
) {
  return render(ui, { wrapper: AllTheProviders, ...options });
}

export * from '@testing-library/react';
export { customRender as render };
```

---

### 7.2. Mock Data

**Файл:** `tests/mocks/data.ts`

```typescript
import type { HomeData } from '@/types/home';
import type { Achievement } from '@/types/achievements';

export const mockHomeData: HomeData = {
  profile: {
    id: 'test-child-1',
    displayName: 'Артём',
    level: 5,
    levelProgress: 60,
    coinsBalance: 1234,
    tasksSolvedCorrectCount: 12,
  },
  mascot: {
    id: 'fox',
    state: 'idle',
    imageUrl: '/mascot-fox.png',
    message: 'Мы справимся! 💪',
  },
  villain: {
    id: 'villain-1',
    name: 'Кракозябра',
    imageUrl: '/villain.png',
    healthPercent: 60,
    currentHealth: 60,
    maxHealth: 100,
    taunt: 'Ха-ха! Попробуй-ка реши задачки!',
  },
  unfinishedAttempt: null,
  recentAttempts: [],
};

export const mockAchievements: Achievement[] = [
  {
    id: 'ach-1',
    type: 'first_task',
    title: 'Первое задание',
    description: 'Выполни своё первое задание',
    icon: '🎯',
    isUnlocked: true,
    unlockedAt: '2026-03-20T10:00:00Z',
    progress: { current: 1, total: 1, percent: 100 },
    reward: {
      type: 'coins',
      id: 'reward-1',
      name: '50 монет',
      amount: 50,
    },
    shelfOrder: 1,
    positionInShelf: 0,
  },
  // ... more achievements
];
```

---

## Чеклист тестирования

### Unit Tests
- [ ] Протестировать все UI компоненты
- [ ] Протестировать все hooks
- [ ] Протестировать все utility функции
- [ ] Протестировать сервисы (analytics, api)
- [ ] Достичь coverage > 80%

### Integration Tests
- [ ] Протестировать onboarding flow
- [ ] Протестировать help flow
- [ ] Протестировать check flow
- [ ] Протестировать paywall flow
- [ ] Протестировать navigation между экранами

### E2E Tests
- [ ] Протестировать критические пользовательские потоки
- [ ] Протестировать на разных браузерах
- [ ] Протестировать на мобильных устройствах
- [ ] Протестировать error states
- [ ] Протестировать loading states

### Performance Tests
- [ ] Lighthouse CI для всех страниц
- [ ] Bundle size анализ
- [ ] Memory leak detection
- [ ] Network waterfall analysis

### Accessibility Tests
- [ ] Axe accessibility scan
- [ ] Keyboard navigation
- [ ] Screen reader compatibility
- [ ] Color contrast checks

### Manual Testing
- [ ] Тестирование на реальных устройствах VK
- [ ] Тестирование на реальных устройствах Max
- [ ] Cross-browser testing
- [ ] Edge case scenarios

---

## Следующий этап

После завершения Testing переходи к **13_PLATFORMS.md** для адаптации под платформы.
