// src/services/analytics/PlatformAdapters.ts
import type {
  StoredAnalyticsEvent,
  AnalyticsConfig,
  UserProperties,
} from '@/types/analytics';
import { BackendAdapter } from './adapters/BackendAdapter';

export class PlatformAdapters {
  private adapters: Map<string, BackendAdapter> = new Map();

  constructor(config: AnalyticsConfig) {
    if (config.platforms.includes('backend')) {
      this.adapters.set('backend', new BackendAdapter());
    }
  }

  async init(): Promise<void> {
    const promises = Array.from(this.adapters.values()).map((adapter) =>
      adapter.init()
    );
    await Promise.all(promises);
  }

  async sendBatch(events: StoredAnalyticsEvent[]): Promise<void> {
    const promises = Array.from(this.adapters.values()).map((adapter) =>
      adapter.sendBatch(events)
    );
    await Promise.all(promises);
  }

  setUserProperties(properties: UserProperties): void {
    this.adapters.forEach((adapter) => {
      adapter.setUserProperties(properties);
    });
  }
}
