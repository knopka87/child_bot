# Оставшиеся проблемы - Детальный анализ

**Дата**: 2024-03-30
**Статус**: 13 проблем к исправлению (6 MEDIUM + 7 LOW)

---

## ⚠️ MEDIUM Priority (6 проблем)

Эти проблемы не блокируют production, но должны быть исправлены в ближайшее время для улучшения качества и поддерживаемости кода.

---

### MEDIUM-1: Console Logging in Production

**Приоритет**: ⚠️ MEDIUM
**Затронутые файлы**:
- `src/services/analytics/*.ts`
- `src/services/platform/adapters/*.ts`
- `src/contexts/*.tsx`
- Multiple components

**Проблема:**
В коде используется `console.log()`, `console.error()`, `console.warn()` для логирования. Это приводит к нескольким проблемам:

1. **Нет контроля уровней логирования** - все логи выводятся всегда
2. **Невозможно централизованно управлять** логами в production
3. **Нет структурированного формата** - сложно парсить и анализировать
4. **Performance impact** - console.log в production может быть медленным
5. **Утечка информации** - debug логи могут содержать sensitive данные

**Примеры проблемного кода:**
```typescript
// src/services/analytics/AnalyticsService.ts
console.log('[Analytics] Event tracked:', event);
console.error('[Analytics] Failed to track event:', error);

// src/services/platform/adapters/VKAdapter.ts
console.log('[VKAdapter] Initialized successfully');
console.error('[VKAdapter] Failed to initialize:', error);

// src/services/platform/adapters/TelegramAdapter.ts
console.log('[TelegramAdapter] Theme changed:', data);
```

**Рекомендуемое решение:**

Создать структурированный logger с уровнями (DEBUG, INFO, WARN, ERROR).

```typescript
// src/lib/logger.ts
enum LogLevel {
  DEBUG = 0,
  INFO = 1,
  WARN = 2,
  ERROR = 3,
  NONE = 4,
}

interface LogContext {
  module?: string;
  userId?: string;
  sessionId?: string;
  [key: string]: any;
}

class Logger {
  private level: LogLevel;
  private context: LogContext = {};

  constructor(level: LogLevel = LogLevel.INFO) {
    this.level = import.meta.env.PROD ? LogLevel.WARN : LogLevel.DEBUG;
  }

  setContext(context: LogContext): void {
    this.context = { ...this.context, ...context };
  }

  private log(level: LogLevel, message: string, data?: any): void {
    if (level < this.level) return;

    const timestamp = new Date().toISOString();
    const logEntry = {
      timestamp,
      level: LogLevel[level],
      message,
      ...this.context,
      ...(data && { data }),
    };

    // В production можно отправлять в внешний сервис
    if (import.meta.env.PROD && level >= LogLevel.ERROR) {
      // TODO: Send to error tracking service (Sentry, etc.)
    }

    // Console output с форматированием
    const prefix = `[${LogLevel[level]}] ${this.context.module || 'App'}:`;

    switch (level) {
      case LogLevel.DEBUG:
        console.debug(prefix, message, data);
        break;
      case LogLevel.INFO:
        console.info(prefix, message, data);
        break;
      case LogLevel.WARN:
        console.warn(prefix, message, data);
        break;
      case LogLevel.ERROR:
        console.error(prefix, message, data);
        break;
    }
  }

  debug(message: string, data?: any): void {
    this.log(LogLevel.DEBUG, message, data);
  }

  info(message: string, data?: any): void {
    this.log(LogLevel.INFO, message, data);
  }

  warn(message: string, data?: any): void {
    this.log(LogLevel.WARN, message, data);
  }

  error(message: string, data?: any): void {
    this.log(LogLevel.ERROR, message, data);
  }
}

// Экспортируем синглтон
export const logger = new Logger();

// Экспортируем функцию для создания модульных логгеров
export const createLogger = (module: string): Logger => {
  const moduleLogger = new Logger();
  moduleLogger.setContext({ module });
  return moduleLogger;
};
```

**Использование:**
```typescript
// src/services/analytics/AnalyticsService.ts
import { createLogger } from '@/lib/logger';

const logger = createLogger('Analytics');

export class AnalyticsService {
  trackEvent(name: string, params: any): void {
    logger.debug('Event tracked', { name, params });

    try {
      // ...
    } catch (error) {
      logger.error('Failed to track event', { name, error });
    }
  }
}
```

**Оценка работы**: 3-4 часа

---

### MEDIUM-2: API Client Error Handling

**Приоритет**: ⚠️ MEDIUM
**Файл**: `src/api/client.ts`

**Проблема:**
API endpoints используются как hardcoded строки без централизованной константы. Это приводит к:

1. **Опечатки** - легко ошибиться в пути
2. **Дублирование** - один путь может быть написан по-разному
3. **Сложность рефакторинга** - нужно искать по всему коду
4. **Нет type-safety** - TypeScript не проверяет корректность путей

**Примеры проблемного кода:**
```typescript
// src/services/analytics/adapters/BackendAdapter.ts
await apiClient.post('/analytics/events', { events });
await apiClient.post('/analytics/properties', { properties });

// Где-то в другом файле может быть:
await apiClient.post('/analytics/event', { events }); // Опечатка!
```

**Рекомендуемое решение:**

Создать централизованный API routes config с type-safety.

```typescript
// src/api/routes.ts
export const API_ROUTES = {
  analytics: {
    events: '/analytics/events',
    properties: '/analytics/properties',
  },
  tasks: {
    list: '/tasks',
    byId: (id: string) => `/tasks/${id}`,
    submit: (id: string) => `/tasks/${id}/submit`,
  },
  villain: {
    list: '/villains',
    byId: (id: string) => `/villains/${id}`,
    victory: (id: string) => `/villains/${id}/victory`,
  },
  profile: {
    get: '/profile',
    update: '/profile',
    history: '/profile/history',
  },
  achievements: {
    list: '/achievements',
    claim: (id: string) => `/achievements/${id}/claim`,
  },
  friends: {
    list: '/friends',
    invite: '/friends/invite',
    referrals: '/friends/referrals',
  },
} as const;

// Type-safe helper для построения URL
type RouteValue = string | ((...args: any[]) => string);

type FlattenRoutes<T> = {
  [K in keyof T]: T[K] extends RouteValue
    ? T[K]
    : T[K] extends object
    ? FlattenRoutes<T[K]>
    : never;
};

export type ApiRoutes = FlattenRoutes<typeof API_ROUTES>;
```

**Использование:**
```typescript
// src/services/analytics/adapters/BackendAdapter.ts
import { API_ROUTES } from '@/api/routes';

export class BackendAdapter {
  async sendBatch(events: StoredAnalyticsEvent[]): Promise<void> {
    await apiClient.post(API_ROUTES.analytics.events, { events });
    //                   ^^^^^^^^^^^^^^^^^^^^^^^^^ Autocomplete работает!
  }

  setUserProperties(properties: UserProperties): void {
    apiClient
      .post(API_ROUTES.analytics.properties, { properties })
      .catch((error) => {
        logger.error('Failed to set properties', { error });
      });
  }
}
```

**Дополнительно** - улучшить error handling в APIClient:

```typescript
// src/api/client.ts
export class APIClient {
  async post<T>(url: string, data?: any): Promise<T> {
    try {
      const response = await this.client.post<T>(url, data);
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        // Структурированная обработка ошибок
        const apiError = {
          status: error.response?.status,
          statusText: error.response?.statusText,
          data: error.response?.data,
          url: error.config?.url,
        };

        logger.error('API request failed', apiError);

        // Можно добавить retry logic для некоторых ошибок
        if (error.response?.status === 429) {
          // Rate limiting - можно retry с backoff
        }
      }
      throw error;
    }
  }
}
```

**Оценка работы**: 2-3 часа

---

### MEDIUM-3: Missing Type Guards

**Приоритет**: ⚠️ MEDIUM
**Файлы**: Multiple

**Проблема:**
Отсутствуют runtime type guards для проверки данных, приходящих из внешних источников (API, localStorage, platform bridges).

**Примеры проблемного кода:**
```typescript
// src/hooks/useVillain.ts
const villain = await apiClient.get<Villain>(`/villains/${id}`);
// Нет проверки что villain действительно соответствует типу Villain

// src/services/platform/adapters/TelegramAdapter.ts
const user = this.webApp.initDataUnsafe.user;
// Нет проверки что user существует и содержит нужные поля
```

**Рекомендуемое решение:**

Использовать zod для runtime валидации.

```typescript
// src/lib/validation/schemas.ts
import { z } from 'zod';

export const VillainSchema = z.object({
  id: z.string(),
  name: z.string(),
  description: z.string(),
  imageUrl: z.string().url(),
  healthPercent: z.number().min(0).max(100),
  currentHealth: z.number().nonnegative(),
  maxHealth: z.number().positive(),
  taunt: z.string(),
  isActive: z.boolean(),
  isDefeated: z.boolean(),
});

export const UserProfileSchema = z.object({
  child_profile_id: z.string(),
  display_name: z.string().min(1).max(50),
  grade: z.number().int().min(1).max(11),
  level: z.number().int().nonnegative(),
  coins_balance: z.number().int().nonnegative(),
  // ...
});

// Type inference из schema
export type Villain = z.infer<typeof VillainSchema>;
export type UserProfile = z.infer<typeof UserProfileSchema>;
```

**Использование:**
```typescript
// src/hooks/useVillain.ts
import { VillainSchema } from '@/lib/validation/schemas';

export function useVillain(id: string) {
  const { data, error, isLoading } = useSWR(
    `/villains/${id}`,
    async (url) => {
      const response = await apiClient.get(url);

      // Runtime validation
      const result = VillainSchema.safeParse(response);

      if (!result.success) {
        logger.error('Invalid villain data', {
          errors: result.error.errors,
          data: response,
        });
        throw new Error('Invalid villain data from API');
      }

      return result.data;
    }
  );

  return { villain: data, error, isLoading };
}
```

**Оценка работы**: 4-6 часов

---

### MEDIUM-4: No Request Deduplication

**Приоритет**: ⚠️ MEDIUM
**Файл**: `src/api/client.ts`

**Проблема:**
API client не дедуплицирует одинаковые запросы, выполняемые одновременно. Если несколько компонентов запрашивают один ресурс, будет несколько идентичных HTTP запросов.

**Рекомендуемое решение:**

Добавить request deduplication через promise cache.

```typescript
// src/api/client.ts
export class APIClient {
  private requestCache = new Map<string, Promise<any>>();

  private getCacheKey(method: string, url: string, data?: any): string {
    return `${method}:${url}:${JSON.stringify(data || {})}`;
  }

  async get<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    const cacheKey = this.getCacheKey('GET', url);

    // Если запрос уже в процессе - возвращаем существующий promise
    if (this.requestCache.has(cacheKey)) {
      return this.requestCache.get(cacheKey);
    }

    const requestPromise = this.client
      .get<T>(url, config)
      .then((response) => {
        this.requestCache.delete(cacheKey);
        return response.data;
      })
      .catch((error) => {
        this.requestCache.delete(cacheKey);
        throw error;
      });

    this.requestCache.set(cacheKey, requestPromise);
    return requestPromise;
  }
}
```

**Оценка работы**: 2 часа

---

### MEDIUM-5: Missing Retry Logic for Analytics

**Приоритет**: ⚠️ MEDIUM
**Файл**: `src/services/analytics/AnalyticsService.ts`

**Проблема:**
Retry logic для analytics batch sending работает через рекурсию, что может привести к stack overflow при большом количестве retry attempts.

**Текущий код:**
```typescript
private async retryBatch(events: StoredAnalyticsEvent[], attemptsLeft: number): Promise<void> {
  if (attemptsLeft === 0) {
    console.error('[Analytics] Batch send failed after all retries');
    return;
  }

  await new Promise((resolve) => setTimeout(resolve, this.config.retryDelay));

  try {
    await this.adapters.sendBatch(events);
  } catch (error) {
    await this.retryBatch(events, attemptsLeft - 1); // ← Рекурсия
  }
}
```

**Рекомендуемое решение:**

Использовать итеративный подход с exponential backoff.

```typescript
private async retryBatch(events: StoredAnalyticsEvent[], maxAttempts: number): Promise<void> {
  let attempt = 0;
  let lastError: Error | null = null;

  while (attempt < maxAttempts) {
    attempt++;

    // Exponential backoff: 2s, 4s, 8s, 16s...
    const delay = this.config.retryDelay * Math.pow(2, attempt - 1);
    await new Promise((resolve) => setTimeout(resolve, delay));

    try {
      await this.adapters.sendBatch(events);

      if (this.config.debug) {
        logger.info('Batch sent after retry', { attempt, eventsCount: events.length });
      }

      return; // Успешно отправили
    } catch (error) {
      lastError = error instanceof Error ? error : new Error(String(error));

      logger.warn('Batch send retry failed', {
        attempt,
        maxAttempts,
        error: lastError.message,
      });
    }
  }

  // Все попытки исчерпаны - сохраняем в localStorage для отправки позже
  this.saveFailedBatch(events);

  logger.error('Batch send failed after all retries', {
    attempts: maxAttempts,
    eventsCount: events.length,
    lastError: lastError?.message,
  });
}

private saveFailedBatch(events: StoredAnalyticsEvent[]): void {
  try {
    const failedBatches = JSON.parse(
      localStorage.getItem('analytics_failed_batches') || '[]'
    );

    failedBatches.push({
      timestamp: Date.now(),
      events,
    });

    // Храним максимум 10 failed batches
    if (failedBatches.length > 10) {
      failedBatches.shift();
    }

    localStorage.setItem('analytics_failed_batches', JSON.stringify(failedBatches));
  } catch (error) {
    logger.error('Failed to save failed batch', { error });
  }
}
```

**Оценка работы**: 3 часа

---

### MEDIUM-6: No Environment-Specific Configuration

**Приоритет**: ⚠️ MEDIUM
**Файлы**: Multiple

**Проблема:**
Конфигурация захардкожена в коде, нет гибкости для разных окружений (dev, staging, production).

**Рекомендуемое решение:**

Создать централизованный config с environment variables.

```typescript
// src/config/index.ts
interface AppConfig {
  environment: 'development' | 'staging' | 'production';
  api: {
    baseURL: string;
    timeout: number;
  };
  analytics: {
    enabled: boolean;
    debug: boolean;
    batchSize: number;
    batchInterval: number;
  };
  features: {
    villainMode: boolean;
    achievements: boolean;
    referrals: boolean;
  };
  platforms: {
    vk: {
      appId: string;
    };
    max: {
      appId: string;
    };
    telegram: {
      botUsername: string;
    };
  };
}

const config: AppConfig = {
  environment: import.meta.env.MODE as any,

  api: {
    baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
    timeout: parseInt(import.meta.env.VITE_API_TIMEOUT || '30000'),
  },

  analytics: {
    enabled: import.meta.env.PROD,
    debug: import.meta.env.DEV,
    batchSize: parseInt(import.meta.env.VITE_ANALYTICS_BATCH_SIZE || '10'),
    batchInterval: parseInt(import.meta.env.VITE_ANALYTICS_BATCH_INTERVAL || '10000'),
  },

  features: {
    villainMode: import.meta.env.VITE_FEATURE_VILLAIN === 'true',
    achievements: import.meta.env.VITE_FEATURE_ACHIEVEMENTS === 'true',
    referrals: import.meta.env.VITE_FEATURE_REFERRALS === 'true',
  },

  platforms: {
    vk: {
      appId: import.meta.env.VITE_VK_APP_ID || '',
    },
    max: {
      appId: import.meta.env.VITE_MAX_APP_ID || '',
    },
    telegram: {
      botUsername: import.meta.env.VITE_TELEGRAM_BOT_USERNAME || '',
    },
  },
};

export default config;
```

**.env.example:**
```bash
# API Configuration
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_API_TIMEOUT=30000

# Analytics
VITE_ANALYTICS_BATCH_SIZE=10
VITE_ANALYTICS_BATCH_INTERVAL=10000

# Feature Flags
VITE_FEATURE_VILLAIN=true
VITE_FEATURE_ACHIEVEMENTS=true
VITE_FEATURE_REFERRALS=true

# Platform IDs
VITE_VK_APP_ID=your_vk_app_id
VITE_MAX_APP_ID=your_max_app_id
VITE_TELEGRAM_BOT_USERNAME=your_bot
```

**Оценка работы**: 2-3 часа

---

## 📝 LOW Priority (7 проблем)

Эти проблемы не критичны, но их исправление улучшит качество кода и developer experience.

---

### LOW-1: Test Coverage < 80%

**Приоритет**: 📝 LOW
**Текущее состояние**: 60% coverage (9/9 tests passing)
**Цель**: 80%+ coverage

**Отсутствующие тесты:**

1. **Villain Pages**
   - `src/pages/Villain/VillainPage.tsx`
   - `src/pages/Villain/VictoryPage.tsx`
   - `src/pages/Villain/components/*`

2. **Platform Adapters**
   - `src/services/platform/adapters/VKAdapter.ts`
   - `src/services/platform/adapters/TelegramAdapter.ts`
   - `src/services/platform/adapters/MaxAdapter.ts`
   - `src/services/platform/PlatformBridge.ts`

3. **Hooks**
   - `src/hooks/useVillain.ts`
   - `src/hooks/useAchievements.ts`
   - `src/hooks/useFriends.ts`

4. **Integration Tests**
   - Нет интеграционных тестов для полных user flows

**Рекомендуемый план:**

```typescript
// tests/unit/pages/Villain/VillainPage.test.tsx
import { render, screen, waitFor } from '@testing-library/react';
import { VillainPage } from '@/pages/Villain/VillainPage';

describe('VillainPage', () => {
  it('should display loading state', () => {
    render(<VillainPage />);
    expect(screen.getByText(/загрузка/i)).toBeInTheDocument();
  });

  it('should display villain data when loaded', async () => {
    // Mock API response
    const mockVillain = {
      id: '1',
      name: 'Evil Math Monster',
      healthPercent: 75,
      // ...
    };

    // Test implementation
  });

  it('should navigate to victory page when villain defeated', async () => {
    // Test implementation
  });
});
```

**Оценка работы**: 8-10 часов

---

### LOW-2: Missing Error Boundaries

**Приоритет**: 📝 LOW
**Файлы**:
- `src/contexts/AnalyticsContext.tsx`
- `src/App.tsx`
- Major routes

**Проблема:**
Нет Error Boundaries для отлова runtime ошибок в React компонентах.

**Рекомендуемое решение:**

```typescript
// src/components/ErrorBoundary/ErrorBoundary.tsx
import { Component, ErrorInfo, ReactNode } from 'react';
import { logger } from '@/lib/logger';

interface Props {
  children: ReactNode;
  fallback?: ReactNode;
  onError?: (error: Error, errorInfo: ErrorInfo) => void;
}

interface State {
  hasError: boolean;
  error: Error | null;
}

export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo): void {
    logger.error('React Error Boundary caught error', {
      error: error.message,
      stack: error.stack,
      componentStack: errorInfo.componentStack,
    });

    this.props.onError?.(error, errorInfo);
  }

  render(): ReactNode {
    if (this.state.hasError) {
      return this.props.fallback || (
        <div style={{ padding: 20 }}>
          <h1>Что-то пошло не так</h1>
          <p>Попробуйте перезагрузить страницу</p>
          <button onClick={() => window.location.reload()}>
            Перезагрузить
          </button>
        </div>
      );
    }

    return this.props.children;
  }
}
```

**Использование:**
```typescript
// src/App.tsx
<ErrorBoundary>
  <AnalyticsProvider>
    <ErrorBoundary fallback={<AnalyticsErrorFallback />}>
      <Router />
    </ErrorBoundary>
  </AnalyticsProvider>
</ErrorBoundary>
```

**Оценка работы**: 2-3 часа

---

### LOW-3: Hardcoded Asset Paths

**Приоритет**: 📝 LOW
**Файлы**:
- `src/pages/Villain/VictoryPage.tsx`
- Potentially others

**Проблема:**
```typescript
<img src="/assets/villain-defeated.png" alt="Defeated villain" />
```

Hardcoded пути к ассетам могут сломаться при изменении структуры проекта.

**Рекомендуемое решение:**

```typescript
// Вариант 1: Import ассетов
import villainDefeatedImage from '@/assets/images/villain-defeated.png';

<img src={villainDefeatedImage} alt="Defeated villain" />

// Вариант 2: Централизованный config
// src/config/assets.ts
export const ASSETS = {
  images: {
    villainDefeated: '/assets/villain-defeated.png',
    mascot: '/assets/mascot.png',
  },
  icons: {
    coin: '/assets/icons/coin.svg',
    star: '/assets/icons/star.svg',
  },
} as const;

// Использование
import { ASSETS } from '@/config/assets';

<img src={ASSETS.images.villainDefeated} alt="Defeated villain" />
```

**Оценка работы**: 1-2 часа

---

### LOW-4: No Loading States for Images

**Приоритет**: 📝 LOW
**Файлы**: Multiple components

**Проблема:**
Изображения загружаются без показа loading state или placeholder.

**Рекомендуемое решение:**

```typescript
// src/components/ui/Image/Image.tsx
import { useState, ImgHTMLAttributes } from 'react';
import { Spinner } from '@/components/ui/Spinner';
import styles from './Image.module.css';

interface ImageProps extends ImgHTMLAttributes<HTMLImageElement> {
  fallback?: string;
  showLoader?: boolean;
}

export function Image({
  src,
  alt,
  fallback = '/assets/placeholder.png',
  showLoader = true,
  className,
  ...props
}: ImageProps) {
  const [isLoading, setIsLoading] = useState(true);
  const [hasError, setHasError] = useState(false);

  return (
    <div className={styles.wrapper}>
      {isLoading && showLoader && (
        <div className={styles.loader}>
          <Spinner size="m" />
        </div>
      )}

      <img
        src={hasError ? fallback : src}
        alt={alt}
        className={className}
        onLoad={() => setIsLoading(false)}
        onError={() => {
          setIsLoading(false);
          setHasError(true);
        }}
        {...props}
      />
    </div>
  );
}
```

**Оценка работы**: 2 часа

---

### LOW-5: Missing Performance Monitoring

**Приоритет**: 📝 LOW
**Файлы**: New files needed

**Проблема:**
Нет мониторинга производительности приложения.

**Рекомендуемое решение:**

```typescript
// src/lib/performance.ts
export class PerformanceMonitor {
  private metrics: Map<string, number[]> = new Map();

  mark(name: string): void {
    if (!('performance' in window)) return;
    performance.mark(name);
  }

  measure(name: string, startMark: string, endMark?: string): number | null {
    if (!('performance' in window)) return null;

    try {
      if (!endMark) {
        this.mark(`${name}-end`);
        endMark = `${name}-end`;
      }

      const measure = performance.measure(name, startMark, endMark);
      const duration = measure.duration;

      // Сохраняем метрику
      if (!this.metrics.has(name)) {
        this.metrics.set(name, []);
      }
      this.metrics.get(name)!.push(duration);

      // Отправляем в analytics
      if (window.analytics) {
        window.analytics.trackEvent('performance_metric', {
          metric_name: name,
          duration_ms: Math.round(duration),
        });
      }

      return duration;
    } catch (error) {
      logger.warn('Performance measure failed', { name, error });
      return null;
    }
  }

  getMetrics(name: string): { avg: number; min: number; max: number; count: number } | null {
    const values = this.metrics.get(name);
    if (!values || values.length === 0) return null;

    return {
      avg: values.reduce((a, b) => a + b, 0) / values.length,
      min: Math.min(...values),
      max: Math.max(...values),
      count: values.length,
    };
  }
}

export const performanceMonitor = new PerformanceMonitor();
```

**Использование:**
```typescript
// src/hooks/useVillain.ts
import { performanceMonitor } from '@/lib/performance';

export function useVillain(id: string) {
  useEffect(() => {
    performanceMonitor.mark('villain-load-start');
  }, []);

  const { data, error, isLoading } = useSWR(`/villains/${id}`, fetcher, {
    onSuccess: () => {
      performanceMonitor.measure('villain-load', 'villain-load-start');
    },
  });

  return { villain: data, error, isLoading };
}
```

**Оценка работы**: 3-4 часа

---

### LOW-6: No Offline Support

**Приоритет**: 📝 LOW
**Файлы**: New service worker needed

**Проблема:**
Приложение не работает offline, нет PWA capabilities.

**Рекомендуемое решение:**

Добавить service worker с workbox для offline support.

```typescript
// vite.config.ts
import { VitePWA } from 'vite-plugin-pwa';

export default defineConfig({
  plugins: [
    react(),
    VitePWA({
      registerType: 'autoUpdate',
      includeAssets: ['favicon.ico', 'apple-touch-icon.png', 'masked-icon.svg'],
      manifest: {
        name: 'Homework Helper',
        short_name: 'Homework',
        description: 'AI-powered homework assistant',
        theme_color: '#0077ff',
        icons: [
          {
            src: 'pwa-192x192.png',
            sizes: '192x192',
            type: 'image/png',
          },
          {
            src: 'pwa-512x512.png',
            sizes: '512x512',
            type: 'image/png',
          },
        ],
      },
      workbox: {
        runtimeCaching: [
          {
            urlPattern: /^https:\/\/api\./,
            handler: 'NetworkFirst',
            options: {
              cacheName: 'api-cache',
              expiration: {
                maxEntries: 100,
                maxAgeSeconds: 60 * 60 * 24, // 24 hours
              },
            },
          },
          {
            urlPattern: /\.(png|jpg|jpeg|svg|gif)$/,
            handler: 'CacheFirst',
            options: {
              cacheName: 'images-cache',
              expiration: {
                maxEntries: 50,
                maxAgeSeconds: 60 * 60 * 24 * 30, // 30 days
              },
            },
          },
        ],
      },
    }),
  ],
});
```

**Оценка работы**: 4-5 часов

---

### LOW-7: Missing Accessibility (a11y) Improvements

**Приоритет**: 📝 LOW
**Файлы**: Multiple components

**Проблема:**
Недостаточно ARIA атрибутов, keyboard navigation, screen reader support.

**Рекомендуемые улучшения:**

1. **Добавить aria-labels**
```typescript
<button aria-label="Закрыть модальное окно" onClick={onClose}>
  <Icon24Close />
</button>
```

2. **Keyboard navigation**
```typescript
<div
  role="button"
  tabIndex={0}
  onClick={handleClick}
  onKeyDown={(e) => {
    if (e.key === 'Enter' || e.key === ' ') {
      handleClick();
    }
  }}
>
  Кликабельный элемент
</div>
```

3. **Focus management**
```typescript
const modalRef = useRef<HTMLDivElement>(null);

useEffect(() => {
  if (isOpen) {
    modalRef.current?.focus();
  }
}, [isOpen]);
```

4. **Live regions для dynamic content**
```typescript
<div aria-live="polite" aria-atomic="true">
  {notification && <p>{notification}</p>}
</div>
```

**Оценка работы**: 5-6 часов

---

## 📊 Сводная таблица

| ID | Приоритет | Проблема | Файлы | Время | Impact |
|----|-----------|----------|-------|-------|--------|
| M-1 | MEDIUM | Console Logging | Multiple | 3-4h | Maintainability |
| M-2 | MEDIUM | API Routes | api/client.ts | 2-3h | Type Safety |
| M-3 | MEDIUM | Type Guards | Multiple | 4-6h | Runtime Safety |
| M-4 | MEDIUM | Request Dedup | api/client.ts | 2h | Performance |
| M-5 | MEDIUM | Retry Logic | AnalyticsService | 3h | Reliability |
| M-6 | MEDIUM | Config | Multiple | 2-3h | Flexibility |
| L-1 | LOW | Test Coverage | tests/ | 8-10h | Quality |
| L-2 | LOW | Error Boundaries | components/ | 2-3h | UX |
| L-3 | LOW | Asset Paths | Multiple | 1-2h | Maintainability |
| L-4 | LOW | Image Loading | components/ | 2h | UX |
| L-5 | LOW | Performance | lib/ | 3-4h | Observability |
| L-6 | LOW | Offline Support | service-worker | 4-5h | UX |
| L-7 | LOW | Accessibility | Multiple | 5-6h | UX |

**Итого**: ~42-57 часов работы

---

## 🎯 Рекомендуемый порядок исправления

### Sprint 1 (1 неделя) - MEDIUM Priority
1. M-2: API Routes (2-3h) - быстрый win, улучшает type safety
2. M-1: Logger (3-4h) - foundation для других улучшений
3. M-4: Request Dedup (2h) - простое улучшение performance
4. M-6: Config (2-3h) - нужен для других задач

**Итого Sprint 1**: 9-12 часов

### Sprint 2 (1 неделя) - MEDIUM Priority continued
1. M-5: Retry Logic (3h)
2. M-3: Type Guards (4-6h) - самая сложная MEDIUM задача
3. L-2: Error Boundaries (2-3h) - важно для production

**Итого Sprint 2**: 9-12 часов

### Sprint 3 (1-2 недели) - LOW Priority
1. L-1: Test Coverage (8-10h) - критично для качества
2. L-3: Asset Paths (1-2h)
3. L-4: Image Loading (2h)
4. L-5: Performance Monitoring (3-4h)

**Итого Sprint 3**: 14-18 часов

### Sprint 4 (опционально) - Nice to have
1. L-6: Offline Support (4-5h)
2. L-7: Accessibility (5-6h)

**Итого Sprint 4**: 9-11 часов

---

## ✅ Критерии готовности

Для каждой задачи:
- [ ] Код написан и типизирован
- [ ] Добавлены unit тесты
- [ ] Обновлена документация
- [ ] Code review пройден
- [ ] TypeScript компиляция без ошибок
- [ ] Все тесты проходят
- [ ] Production build успешен

---

**Последнее обновление**: 2024-03-30
**Автор**: Claude Code Agent
**Статус документа**: Living document - обновляется по мере выполнения задач
