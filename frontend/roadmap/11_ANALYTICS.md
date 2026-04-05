# Phase 10: Интеграция аналитики (Analytics Integration)

**Длительность:** 2-3 дня
**Приоритет:** Критический
**Зависимости:** Все предыдущие фазы

---

## Цель

Создать полную систему аналитики: архитектура analytics service, отправка всех событий из реестра, user properties, валидация событий, дебаг режим, batch отправка, retry логика.

---

## Архитектура Analytics

### Структура компонентов

```
Analytics Service
├── Event Queue (batch)
├── Event Validator
├── Platform Adapters
│   ├── VK Analytics
│   ├── Max Analytics
│   └── Custom Backend
├── User Properties Manager
├── Session Manager
└── Debug Console
```

---

## Часть 1: Типы и Интерфейсы

### 1.1. TypeScript Types

**Файл:** `src/types/analytics.ts`

```typescript
export type AnalyticsEventName =
  // Onboarding
  | 'onboarding_opened'
  | 'registration_opened'
  | 'consent_screen_opened'
  | 'grade_selected'
  | 'avatar_selected'
  | 'display_name_entered'
  | 'adult_consent_checked'
  | 'privacy_policy_opened'
  | 'privacy_policy_accepted'
  | 'terms_opened'
  | 'terms_accepted'
  | 'email_entered'
  | 'email_verification_sent'
  | 'email_verification_success'
  | 'onboarding_completed'
  // Home
  | 'home_opened'
  | 'level_bar_viewed'
  | 'coins_balance_viewed'
  | 'tasks_correct_count_viewed'
  | 'home_help_clicked'
  | 'home_check_clicked'
  | 'unfinished_attempt_modal_shown'
  | 'unfinished_attempt_continue_clicked'
  | 'unfinished_attempt_new_task_clicked'
  | 'mascot_clicked'
  | 'villain_clicked'
  | 'recent_attempt_clicked'
  | 'recent_attempts_view_all_clicked'
  // ... (все остальные события из реестра)
  // Achievements
  | 'achievements_opened'
  | 'achievement_clicked'
  | 'achievement_unlocked'
  // Friends
  | 'friends_opened'
  | 'referral_link_copied'
  | 'referral_share_sent'
  // Profile
  | 'profile_opened'
  | 'history_opened'
  // Villain
  | 'villain_screen_opened'
  | 'victory_screen_opened'
  // Support
  | 'support_opened'
  | 'support_message_sent'
  // Paywall
  | 'paywall_opened'
  | 'payment_started'
  | 'payment_success'
  | 'payment_failed';

export interface AnalyticsEvent {
  name: AnalyticsEventName;
  timestamp: number;
  sessionId: string;
  params: Record<string, any>;
}

export interface AnalyticsConfig {
  enabled: boolean;
  debug: boolean;
  batchSize: number;
  batchInterval: number; // ms
  retryAttempts: number;
  retryDelay: number; // ms
  platforms: AnalyticsPlatform[];
}

export type AnalyticsPlatform = 'vk' | 'max' | 'backend';

export interface UserProperties {
  // Parent user properties
  platform_type?: 'vk' | 'max' | 'web';
  subscription_status?: 'trial' | 'active' | 'expired' | 'cancelled';
  trial_status?: string;
  email_verified?: boolean;
  weekly_report_enabled?: boolean;
  report_archive_enabled?: boolean;

  // Child profile properties
  grade?: number;
  level?: number;
  coins_balance?: number;
  tasks_solved_correct_count?: number;
  wins_count?: number;
  checks_correct_count?: number;
  current_streak_days?: number;
  has_unfinished_attempt?: boolean;
  active_villain_id?: string;
  active_villain_health_percent?: number;
  invited_count_total?: number;
  achievements_unlocked_count?: number;
}

export interface AnalyticsSession {
  id: string;
  startedAt: number;
  platform: string;
  appVersion: string;
  userId?: string;
  childProfileId?: string;
}
```

---

## Часть 2: Analytics Service

### 2.1. Core Analytics Service

**Файл:** `src/services/analytics/AnalyticsService.ts`

```typescript
import { EventQueue } from './EventQueue';
import { EventValidator } from './EventValidator';
import { PlatformAdapters } from './PlatformAdapters';
import { UserPropertiesManager } from './UserPropertiesManager';
import { SessionManager } from './SessionManager';
import type {
  AnalyticsEvent,
  AnalyticsEventName,
  AnalyticsConfig,
  UserProperties,
} from '@/types/analytics';

export class AnalyticsService {
  private config: AnalyticsConfig;
  private eventQueue: EventQueue;
  private validator: EventValidator;
  private adapters: PlatformAdapters;
  private userProperties: UserPropertiesManager;
  private sessionManager: SessionManager;
  private isInitialized = false;

  constructor(config: AnalyticsConfig) {
    this.config = config;
    this.eventQueue = new EventQueue(config);
    this.validator = new EventValidator();
    this.adapters = new PlatformAdapters(config);
    this.userProperties = new UserPropertiesManager();
    this.sessionManager = new SessionManager();
  }

  /**
   * Инициализация сервиса
   */
  async init(): Promise<void> {
    if (this.isInitialized) return;

    try {
      // Start session
      await this.sessionManager.startSession();

      // Initialize platform adapters
      await this.adapters.init();

      // Setup event queue flush
      this.eventQueue.onFlush(async (events) => {
        await this.sendBatch(events);
      });

      this.isInitialized = true;

      if (this.config.debug) {
        console.log('[Analytics] Service initialized');
      }
    } catch (error) {
      console.error('[Analytics] Failed to initialize:', error);
    }
  }

  /**
   * Отправить событие
   */
  trackEvent(
    name: AnalyticsEventName,
    params: Record<string, any> = {}
  ): void {
    if (!this.config.enabled) return;

    try {
      // Validate event
      const validationResult = this.validator.validate(name, params);
      if (!validationResult.isValid) {
        console.warn(
          `[Analytics] Invalid event "${name}":`,
          validationResult.errors
        );
        if (!this.config.debug) return;
      }

      // Create event
      const event: AnalyticsEvent = {
        name,
        timestamp: Date.now(),
        sessionId: this.sessionManager.getSessionId(),
        params: {
          ...params,
          app_version: import.meta.env.VITE_APP_VERSION,
          platform_type: this.sessionManager.getPlatform(),
        },
      };

      // Add to queue
      this.eventQueue.push(event);

      // Debug log
      if (this.config.debug) {
        console.log('[Analytics] Event tracked:', event);
      }
    } catch (error) {
      console.error('[Analytics] Failed to track event:', error);
    }
  }

  /**
   * Установить user properties
   */
  setUserProperties(properties: UserProperties): void {
    if (!this.config.enabled) return;

    try {
      this.userProperties.set(properties);

      // Send to platforms
      this.adapters.setUserProperties(properties);

      if (this.config.debug) {
        console.log('[Analytics] User properties set:', properties);
      }
    } catch (error) {
      console.error('[Analytics] Failed to set user properties:', error);
    }
  }

  /**
   * Обновить user properties
   */
  updateUserProperties(properties: Partial<UserProperties>): void {
    if (!this.config.enabled) return;

    try {
      this.userProperties.update(properties);

      // Send to platforms
      this.adapters.setUserProperties(
        this.userProperties.getAll()
      );

      if (this.config.debug) {
        console.log('[Analytics] User properties updated:', properties);
      }
    } catch (error) {
      console.error('[Analytics] Failed to update user properties:', error);
    }
  }

  /**
   * Получить ID сессии
   */
  getSessionId(): string {
    return this.sessionManager.getSessionId();
  }

  /**
   * Отправить batch событий
   */
  private async sendBatch(events: AnalyticsEvent[]): Promise<void> {
    if (events.length === 0) return;

    try {
      // Send to all platforms
      await this.adapters.sendBatch(events);

      if (this.config.debug) {
        console.log(`[Analytics] Batch sent: ${events.length} events`);
      }
    } catch (error) {
      console.error('[Analytics] Failed to send batch:', error);

      // Retry logic
      if (this.config.retryAttempts > 0) {
        await this.retryBatch(events, this.config.retryAttempts);
      }
    }
  }

  /**
   * Retry отправки batch
   */
  private async retryBatch(
    events: AnalyticsEvent[],
    attemptsLeft: number
  ): Promise<void> {
    if (attemptsLeft === 0) {
      console.error('[Analytics] Batch send failed after all retries');
      return;
    }

    await new Promise((resolve) =>
      setTimeout(resolve, this.config.retryDelay)
    );

    try {
      await this.adapters.sendBatch(events);
      if (this.config.debug) {
        console.log('[Analytics] Batch sent after retry');
      }
    } catch (error) {
      await this.retryBatch(events, attemptsLeft - 1);
    }
  }

  /**
   * Flush всех событий
   */
  async flush(): Promise<void> {
    await this.eventQueue.flush();
  }

  /**
   * Очистить все данные
   */
  clear(): void {
    this.eventQueue.clear();
    this.userProperties.clear();
    this.sessionManager.endSession();
  }
}
```

---

### 2.2. Event Queue

**Файл:** `src/services/analytics/EventQueue.ts`

```typescript
import type { AnalyticsEvent, AnalyticsConfig } from '@/types/analytics';

export class EventQueue {
  private queue: AnalyticsEvent[] = [];
  private flushTimer: NodeJS.Timeout | null = null;
  private flushCallback: ((events: AnalyticsEvent[]) => Promise<void>) | null =
    null;

  constructor(private config: AnalyticsConfig) {
    this.startFlushTimer();
  }

  /**
   * Добавить событие в очередь
   */
  push(event: AnalyticsEvent): void {
    this.queue.push(event);

    // Flush if batch size reached
    if (this.queue.length >= this.config.batchSize) {
      this.flush();
    }
  }

  /**
   * Установить callback для flush
   */
  onFlush(callback: (events: AnalyticsEvent[]) => Promise<void>): void {
    this.flushCallback = callback;
  }

  /**
   * Flush всех событий
   */
  async flush(): Promise<void> {
    if (this.queue.length === 0 || !this.flushCallback) return;

    const events = [...this.queue];
    this.queue = [];

    await this.flushCallback(events);
  }

  /**
   * Очистить очередь
   */
  clear(): void {
    this.queue = [];
  }

  /**
   * Запустить таймер flush
   */
  private startFlushTimer(): void {
    this.flushTimer = setInterval(() => {
      this.flush();
    }, this.config.batchInterval);
  }

  /**
   * Остановить таймер
   */
  destroy(): void {
    if (this.flushTimer) {
      clearInterval(this.flushTimer);
      this.flushTimer = null;
    }
  }
}
```

---

### 2.3. Event Validator

**Файл:** `src/services/analytics/EventValidator.ts`

```typescript
import type { AnalyticsEventName } from '@/types/analytics';
import { ANALYTICS_SCHEMA } from './schema';

interface ValidationResult {
  isValid: boolean;
  errors: string[];
}

export class EventValidator {
  /**
   * Валидация события
   */
  validate(
    name: AnalyticsEventName,
    params: Record<string, any>
  ): ValidationResult {
    const schema = ANALYTICS_SCHEMA[name];

    if (!schema) {
      return {
        isValid: false,
        errors: [`Unknown event: ${name}`],
      };
    }

    const errors: string[] = [];

    // Check required params
    schema.required?.forEach((param) => {
      if (!(param in params)) {
        errors.push(`Missing required parameter: ${param}`);
      }
    });

    // Check param types
    Object.entries(params).forEach(([key, value]) => {
      const expectedType = schema.params[key];
      if (!expectedType) return;

      const actualType = typeof value;
      if (actualType !== expectedType) {
        errors.push(
          `Invalid type for ${key}: expected ${expectedType}, got ${actualType}`
        );
      }
    });

    return {
      isValid: errors.length === 0,
      errors,
    };
  }
}
```

---

### 2.4. Analytics Schema

**Файл:** `src/services/analytics/schema.ts`

```typescript
import type { AnalyticsEventName } from '@/types/analytics';

interface EventSchema {
  required?: string[];
  optional?: string[];
  params: Record<string, 'string' | 'number' | 'boolean'>;
}

export const ANALYTICS_SCHEMA: Record<AnalyticsEventName, EventSchema> = {
  onboarding_opened: {
    required: ['platform_type', 'session_id'],
    optional: ['entry_point'],
    params: {
      platform_type: 'string',
      session_id: 'string',
      entry_point: 'string',
    },
  },
  grade_selected: {
    required: ['grade'],
    optional: ['child_profile_id'],
    params: {
      grade: 'number',
      child_profile_id: 'string',
    },
  },
  home_opened: {
    required: ['child_profile_id'],
    optional: ['entry_point'],
    params: {
      child_profile_id: 'string',
      entry_point: 'string',
    },
  },
  // ... (добавить все остальные события из реестра)
};
```

---

### 2.5. Platform Adapters

**Файл:** `src/services/analytics/PlatformAdapters.ts`

```typescript
import type {
  AnalyticsEvent,
  AnalyticsConfig,
  UserProperties,
} from '@/types/analytics';
import { VKAdapter } from './adapters/VKAdapter';
import { MaxAdapter } from './adapters/MaxAdapter';
import { BackendAdapter } from './adapters/BackendAdapter';

export class PlatformAdapters {
  private adapters: Map<string, any> = new Map();

  constructor(private config: AnalyticsConfig) {
    // Initialize adapters based on platform
    if (config.platforms.includes('vk')) {
      this.adapters.set('vk', new VKAdapter());
    }
    if (config.platforms.includes('max')) {
      this.adapters.set('max', new MaxAdapter());
    }
    if (config.platforms.includes('backend')) {
      this.adapters.set('backend', new BackendAdapter());
    }
  }

  /**
   * Инициализация адаптеров
   */
  async init(): Promise<void> {
    const promises = Array.from(this.adapters.values()).map((adapter) =>
      adapter.init()
    );
    await Promise.all(promises);
  }

  /**
   * Отправить batch событий
   */
  async sendBatch(events: AnalyticsEvent[]): Promise<void> {
    const promises = Array.from(this.adapters.values()).map((adapter) =>
      adapter.sendBatch(events)
    );
    await Promise.all(promises);
  }

  /**
   * Установить user properties
   */
  setUserProperties(properties: UserProperties): void {
    this.adapters.forEach((adapter) => {
      adapter.setUserProperties(properties);
    });
  }
}
```

---

### 2.6. Backend Adapter

**Файл:** `src/services/analytics/adapters/BackendAdapter.ts`

```typescript
import { apiClient } from '@/api/client';
import type { AnalyticsEvent, UserProperties } from '@/types/analytics';

export class BackendAdapter {
  async init(): Promise<void> {
    // No initialization needed
  }

  /**
   * Отправить batch событий
   */
  async sendBatch(events: AnalyticsEvent[]): Promise<void> {
    try {
      await apiClient.post('/analytics/events', { events });
    } catch (error) {
      console.error('[BackendAdapter] Failed to send batch:', error);
      throw error;
    }
  }

  /**
   * Установить user properties
   */
  setUserProperties(properties: UserProperties): void {
    try {
      apiClient.post('/analytics/properties', { properties });
    } catch (error) {
      console.error('[BackendAdapter] Failed to set properties:', error);
    }
  }
}
```

---

## Часть 3: React Hook

### 3.1. useAnalytics Hook

**Файл:** `src/hooks/useAnalytics.ts`

```typescript
import { useContext } from 'react';
import { AnalyticsContext } from '@/contexts/AnalyticsContext';
import type { AnalyticsEventName } from '@/types/analytics';

export function useAnalytics() {
  const analytics = useContext(AnalyticsContext);

  if (!analytics) {
    throw new Error('useAnalytics must be used within AnalyticsProvider');
  }

  return {
    trackEvent: (name: AnalyticsEventName, params?: Record<string, any>) => {
      analytics.trackEvent(name, params);
    },
    setUserProperties: analytics.setUserProperties.bind(analytics),
    updateUserProperties: analytics.updateUserProperties.bind(analytics),
    sessionId: analytics.getSessionId(),
  };
}
```

---

### 3.2. AnalyticsProvider

**Файл:** `src/contexts/AnalyticsContext.tsx`

```typescript
import { createContext, useEffect, useRef, ReactNode } from 'react';
import { AnalyticsService } from '@/services/analytics/AnalyticsService';
import type { AnalyticsConfig } from '@/types/analytics';

const defaultConfig: AnalyticsConfig = {
  enabled: import.meta.env.PROD,
  debug: import.meta.env.DEV,
  batchSize: 10,
  batchInterval: 10000, // 10 seconds
  retryAttempts: 3,
  retryDelay: 2000, // 2 seconds
  platforms: ['backend'],
};

export const AnalyticsContext = createContext<AnalyticsService | null>(null);

interface AnalyticsProviderProps {
  children: ReactNode;
  config?: Partial<AnalyticsConfig>;
}

export function AnalyticsProvider({
  children,
  config = {},
}: AnalyticsProviderProps) {
  const analyticsRef = useRef<AnalyticsService | null>(null);

  useEffect(() => {
    const mergedConfig = { ...defaultConfig, ...config };
    analyticsRef.current = new AnalyticsService(mergedConfig);
    analyticsRef.current.init();

    // Flush on page unload
    const handleBeforeUnload = () => {
      analyticsRef.current?.flush();
    };

    window.addEventListener('beforeunload', handleBeforeUnload);

    return () => {
      window.removeEventListener('beforeunload', handleBeforeUnload);
      analyticsRef.current?.flush();
    };
  }, [config]);

  return (
    <AnalyticsContext.Provider value={analyticsRef.current}>
      {children}
    </AnalyticsContext.Provider>
  );
}
```

---

## Часть 4: Примеры интеграции

### 4.1. Интеграция в компонент

```typescript
import { useEffect } from 'react';
import { useAnalytics } from '@/hooks/useAnalytics';

export function HomePage() {
  const analytics = useAnalytics();

  useEffect(() => {
    // Track page view
    analytics.trackEvent('home_opened', {
      child_profile_id: 'xxx',
      entry_point: 'direct',
    });
  }, []);

  const handleHelpClick = () => {
    analytics.trackEvent('home_help_clicked', {
      child_profile_id: 'xxx',
    });
  };

  return (
    <div>
      <button onClick={handleHelpClick}>Помоги разобраться</button>
    </div>
  );
}
```

---

### 4.2. Обновление user properties

```typescript
import { useEffect } from 'react';
import { useAnalytics } from '@/hooks/useAnalytics';
import { useProfileStore } from '@/stores/profileStore';

export function App() {
  const analytics = useAnalytics();
  const profile = useProfileStore((state) => state.profile);

  useEffect(() => {
    if (profile) {
      analytics.updateUserProperties({
        grade: profile.grade,
        level: profile.level,
        coins_balance: profile.coinsBalance,
      });
    }
  }, [profile, analytics]);

  return <div>...</div>;
}
```

---

## Часть 5: Debug Console

### 5.1. Debug Panel Component

**Файл:** `src/components/debug/AnalyticsDebugPanel.tsx`

```typescript
import { useState, useEffect } from 'react';
import { useAnalytics } from '@/hooks/useAnalytics';
import styles from './AnalyticsDebugPanel.module.css';

export function AnalyticsDebugPanel() {
  const analytics = useAnalytics();
  const [events, setEvents] = useState<any[]>([]);
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    // Listen to console.log for analytics events
    const originalLog = console.log;
    console.log = (...args) => {
      if (args[0]?.includes('[Analytics]')) {
        setEvents((prev) => [...prev, { timestamp: Date.now(), args }]);
      }
      originalLog(...args);
    };

    return () => {
      console.log = originalLog;
    };
  }, []);

  if (!import.meta.env.DEV) return null;

  return (
    <>
      <button
        className={styles.toggleButton}
        onClick={() => setIsVisible(!isVisible)}
      >
        📊
      </button>

      {isVisible && (
        <div className={styles.panel}>
          <div className={styles.header}>
            <h3>Analytics Debug</h3>
            <button onClick={() => setEvents([])}>Clear</button>
          </div>
          <div className={styles.events}>
            {events.map((event, i) => (
              <div key={i} className={styles.event}>
                <pre>{JSON.stringify(event.args, null, 2)}</pre>
              </div>
            ))}
          </div>
        </div>
      )}
    </>
  );
}
```

---

## Чеклист задач

### Core Service
- [ ] Создать AnalyticsService
- [ ] Создать EventQueue с batch отправкой
- [ ] Создать EventValidator
- [ ] Создать SessionManager
- [ ] Создать UserPropertiesManager

### Platform Adapters
- [ ] Реализовать BackendAdapter
- [ ] Реализовать VKAdapter
- [ ] Реализовать MaxAdapter
- [ ] Добавить retry логику

### React Integration
- [ ] Создать AnalyticsProvider
- [ ] Создать useAnalytics hook
- [ ] Интегрировать в App

### Schema & Validation
- [ ] Создать полную схему всех событий
- [ ] Добавить валидацию параметров
- [ ] Добавить типы для всех событий

### Debug Tools
- [ ] Создать AnalyticsDebugPanel
- [ ] Добавить debug logging
- [ ] Добавить event inspector

### Testing
- [ ] Протестировать отправку событий
- [ ] Протестировать batch отправку
- [ ] Протестировать retry логику
- [ ] Протестировать валидацию

---

## Следующий этап

После завершения Analytics переходи к **12_TESTING.md** для настройки тестирования.
