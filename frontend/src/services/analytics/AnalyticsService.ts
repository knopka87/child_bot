// src/services/analytics/AnalyticsService.ts
import { EventQueue } from './EventQueue';
import { EventValidator } from './EventValidator';
import { PlatformAdapters } from './PlatformAdapters';
import { UserPropertiesManager } from './UserPropertiesManager';
import { SessionManager } from './SessionManager';
import { createLogger } from '@/lib/logger';
import type {
  AnalyticsEventName,
  AnalyticsEventParams,
  AnalyticsConfig,
  UserProperties,
  StoredAnalyticsEvent,
} from '@/types/analytics';
import type { PlatformType } from '@/types/platform';

const logger = createLogger('Analytics');
const FAILED_BATCHES_KEY = 'analytics_failed_batches';

export class AnalyticsService {
  private config: AnalyticsConfig;
  private eventQueue: EventQueue;
  private validator: EventValidator;
  private adapters: PlatformAdapters;
  private userProperties: UserPropertiesManager;
  private sessionManager: SessionManager;
  private isInitialized = false;

  constructor(config: AnalyticsConfig, platformType?: PlatformType) {
    this.config = config;
    this.eventQueue = new EventQueue(config);
    this.validator = new EventValidator();
    this.adapters = new PlatformAdapters(config);
    this.userProperties = new UserPropertiesManager();
    this.sessionManager = new SessionManager(platformType);
  }

  async init(): Promise<void> {
    if (this.isInitialized) return;

    try {
      await this.sessionManager.startSession();

      await this.adapters.init();

      this.eventQueue.onFlush(async (events) => {
        await this.sendBatch(events);
      });

      // Пытаемся отправить failed batches из предыдущих сессий
      await this.retryFailedBatches();

      this.isInitialized = true;

      logger.info('Service initialized', {
        sessionId: this.sessionManager.getSessionId(),
        platform: this.sessionManager.getPlatform(),
      });
    } catch (error) {
      logger.error('Failed to initialize', { error });
    }
  }

  trackEvent<T extends AnalyticsEventName>(
    name: T,
    params: Omit<AnalyticsEventParams<T>, 'app_version' | 'platform_type' | 'child_profile_id'> = {} as any
  ): void {
    if (!this.config.enabled) return;

    try {
      const validationResult = this.validator.validate(name, params);
      if (!validationResult.isValid) {
        logger.warn(`Invalid event "${name}"`, {
          errors: validationResult.errors,
        });
        if (!this.config.debug) return;
      }

      const event: StoredAnalyticsEvent = {
        name,
        timestamp: Date.now(),
        sessionId: this.sessionManager.getSessionId(),
        params: {
          ...params,
          app_version: import.meta.env.VITE_APP_VERSION || '0.1.0',
          platform_type: this.sessionManager.getPlatform(),
        },
      };

      this.eventQueue.push(event);

      logger.debug('Event tracked', { name, params });
    } catch (error) {
      logger.error('Failed to track event', { name, error });
    }
  }

  setUserProperties(properties: UserProperties): void {
    if (!this.config.enabled) return;

    try {
      this.userProperties.set(properties);
      this.adapters.setUserProperties(properties);

      logger.debug('User properties set', { properties });
    } catch (error) {
      logger.error('Failed to set user properties', { error });
    }
  }

  updateUserProperties(properties: Partial<UserProperties>): void {
    if (!this.config.enabled) return;

    try {
      this.userProperties.update(properties);
      this.adapters.setUserProperties(this.userProperties.getAll());

      logger.debug('User properties updated', { properties });
    } catch (error) {
      logger.error('Failed to update user properties', { error });
    }
  }

  getSessionId(): string {
    return this.sessionManager.getSessionId();
  }

  private async sendBatch(events: StoredAnalyticsEvent[]): Promise<void> {
    if (events.length === 0) return;

    try {
      await this.adapters.sendBatch(events);

      logger.debug('Batch sent', { eventsCount: events.length });
    } catch (error) {
      logger.error('Failed to send batch', {
        eventsCount: events.length,
        error,
      });

      if (this.config.retryAttempts > 0) {
        await this.retryBatch(events, this.config.retryAttempts);
      } else {
        // Нет retry - сохраняем сразу
        this.saveFailedBatch(events);
      }
    }
  }

  /**
   * Retry batch с exponential backoff и сохранением в localStorage при неудаче
   */
  private async retryBatch(
    events: StoredAnalyticsEvent[],
    maxAttempts: number
  ): Promise<void> {
    let attempt = 0;
    let lastError: Error | null = null;

    while (attempt < maxAttempts) {
      attempt++;

      // Exponential backoff: 2s, 4s, 8s...
      const delay = this.config.retryDelay * Math.pow(2, attempt - 1);
      await new Promise((resolve) => setTimeout(resolve, delay));

      try {
        await this.adapters.sendBatch(events);

        logger.info('Batch sent after retry', {
          attempt,
          eventsCount: events.length,
        });

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

    // Все попытки исчерпаны - сохраняем в localStorage
    this.saveFailedBatch(events);

    logger.error('Batch send failed after all retries', {
      attempts: maxAttempts,
      eventsCount: events.length,
      lastError: lastError?.message,
    });
  }

  /**
   * Сохраняет failed batch в localStorage для последующей отправки
   */
  private saveFailedBatch(events: StoredAnalyticsEvent[]): void {
    try {
      const failedBatches = JSON.parse(
        localStorage.getItem(FAILED_BATCHES_KEY) || '[]'
      ) as Array<{ timestamp: number; events: StoredAnalyticsEvent[] }>;

      failedBatches.push({
        timestamp: Date.now(),
        events,
      });

      // Храним максимум 10 failed batches (чтобы не переполнить localStorage)
      if (failedBatches.length > 10) {
        failedBatches.shift();
      }

      localStorage.setItem(FAILED_BATCHES_KEY, JSON.stringify(failedBatches));

      logger.warn('Saved failed batch to localStorage', {
        eventsCount: events.length,
        totalFailedBatches: failedBatches.length,
      });
    } catch (error) {
      logger.error('Failed to save failed batch', { error });
    }
  }

  /**
   * Пытается отправить failed batches из localStorage
   */
  private async retryFailedBatches(): Promise<void> {
    try {
      const failedBatches = JSON.parse(
        localStorage.getItem(FAILED_BATCHES_KEY) || '[]'
      ) as Array<{ timestamp: number; events: StoredAnalyticsEvent[] }>;

      if (failedBatches.length === 0) return;

      logger.info('Retrying failed batches', {
        batchesCount: failedBatches.length,
      });

      // Пытаемся отправить каждый batch
      const successfullyRetried: number[] = [];

      for (let i = 0; i < failedBatches.length; i++) {
        const batch = failedBatches[i];

        try {
          await this.adapters.sendBatch(batch.events);

          logger.info('Failed batch sent successfully', {
            batchIndex: i,
            eventsCount: batch.events.length,
          });

          successfullyRetried.push(i);
        } catch (error) {
          logger.warn('Failed to retry batch', {
            batchIndex: i,
            error,
          });
        }
      }

      // Удаляем успешно отправленные batches
      if (successfullyRetried.length > 0) {
        const remainingBatches = failedBatches.filter(
          (_, index) => !successfullyRetried.includes(index)
        );

        localStorage.setItem(
          FAILED_BATCHES_KEY,
          JSON.stringify(remainingBatches)
        );

        logger.info('Failed batches cleanup complete', {
          retriedCount: successfullyRetried.length,
          remainingCount: remainingBatches.length,
        });
      }
    } catch (error) {
      logger.error('Failed to retry failed batches', { error });
    }
  }

  async flush(): Promise<void> {
    await this.eventQueue.flush();
  }

  clear(): void {
    this.eventQueue.clear();
    this.userProperties.clear();
    this.sessionManager.endSession();
  }

  destroy(): void {
    this.eventQueue.destroy();
    this.clear();
  }
}
