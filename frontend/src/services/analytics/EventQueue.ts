// src/services/analytics/EventQueue.ts
import type { StoredAnalyticsEvent, AnalyticsConfig } from '@/types/analytics';

export class EventQueue {
  private queue: StoredAnalyticsEvent[] = [];
  private flushTimer: ReturnType<typeof setInterval> | null = null;
  private flushCallback: ((events: StoredAnalyticsEvent[]) => Promise<void>) | null =
    null;

  constructor(private config: AnalyticsConfig) {
    this.startFlushTimer();
  }

  push(event: StoredAnalyticsEvent): void {
    this.queue.push(event);

    if (this.queue.length >= this.config.batchSize) {
      void this.flush();
    }
  }

  onFlush(callback: (events: StoredAnalyticsEvent[]) => Promise<void>): void {
    this.flushCallback = callback;
  }

  async flush(): Promise<void> {
    if (this.queue.length === 0 || !this.flushCallback) return;

    const events = [...this.queue];
    this.queue = [];

    await this.flushCallback(events);
  }

  clear(): void {
    this.queue = [];
  }

  private startFlushTimer(): void {
    this.flushTimer = setInterval(() => {
      void this.flush();
    }, this.config.batchInterval);
  }

  destroy(): void {
    if (this.flushTimer) {
      clearInterval(this.flushTimer);
      this.flushTimer = null;
    }
  }
}
