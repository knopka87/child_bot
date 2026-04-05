// src/services/analytics/SessionManager.ts
import type { AnalyticsSession } from '@/types/analytics';
import type { PlatformType } from '@/types/platform';

export class SessionManager {
  private session: AnalyticsSession | null = null;
  private platformType: PlatformType = 'web';

  constructor(platformType?: PlatformType) {
    this.platformType = platformType || 'web';
  }

  async startSession(): Promise<void> {
    this.session = {
      id: this.generateSessionId(),
      startedAt: Date.now(),
      platform: this.platformType,
      appVersion: import.meta.env.VITE_APP_VERSION || '0.1.0',
    };
  }

  getSessionId(): string {
    return this.session?.id || '';
  }

  getPlatform(): string {
    return this.session?.platform || 'web';
  }

  setUserId(userId: string): void {
    if (this.session) {
      this.session.userId = userId;
    }
  }

  setChildProfileId(childProfileId: string): void {
    if (this.session) {
      this.session.childProfileId = childProfileId;
    }
  }

  endSession(): void {
    this.session = null;
  }

  private generateSessionId(): string {
    return `${Date.now()}-${Math.random().toString(36).substring(2, 11)}`;
  }
}
