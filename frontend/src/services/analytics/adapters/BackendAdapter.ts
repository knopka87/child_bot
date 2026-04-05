// src/services/analytics/adapters/BackendAdapter.ts
import { apiClient } from '@/api/client';
import { API_ROUTES } from '@/api/routes';
import { createLogger } from '@/lib/logger';
import type { StoredAnalyticsEvent, UserProperties } from '@/types/analytics';

const logger = createLogger('BackendAdapter');

export class BackendAdapter {
  async init(): Promise<void> {
    logger.debug('Initialized');
  }

  async sendBatch(events: StoredAnalyticsEvent[]): Promise<void> {
    try {
      await apiClient.post(API_ROUTES.analytics.events, { events });
      logger.debug('Batch sent successfully', { eventsCount: events.length });
    } catch (error) {
      logger.error('Failed to send batch', {
        eventsCount: events.length,
        error,
      });
      throw error;
    }
  }

  setUserProperties(properties: UserProperties): void {
    apiClient
      .post(API_ROUTES.analytics.properties, { properties })
      .catch((error) => {
        logger.error('Failed to set properties', { error });
      });
  }
}
