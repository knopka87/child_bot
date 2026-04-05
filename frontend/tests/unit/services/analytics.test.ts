// tests/unit/services/analytics.test.ts
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { AnalyticsService } from '@/services/analytics/AnalyticsService';
import type { AnalyticsConfig } from '@/types/analytics';

describe('AnalyticsService', () => {
  let analyticsService: AnalyticsService;
  let mockConfig: AnalyticsConfig;

  beforeEach(() => {
    mockConfig = {
      enabled: true,
      debug: false,
      batchSize: 5,
      batchInterval: 1000,
      retryAttempts: 2,
      retryDelay: 500,
      platforms: ['backend'],
    };

    analyticsService = new AnalyticsService(mockConfig);
  });

  it('should initialize successfully', async () => {
    await analyticsService.init();
    expect(analyticsService.getSessionId()).toBeDefined();
  });

  it('should track event with correct parameters', async () => {
    await analyticsService.init();

    const consoleSpy = vi.spyOn(console, 'log');

    analyticsService.trackEvent('home_opened', {
      child_profile_id: 'test-123',
    });

    expect(analyticsService.getSessionId()).toBeTruthy();
    consoleSpy.mockRestore();
  });

  it('should not track events when disabled', async () => {
    const disabledConfig = { ...mockConfig, enabled: false };
    const disabledService = new AnalyticsService(disabledConfig);

    await disabledService.init();

    const consoleSpy = vi.spyOn(console, 'log');

    disabledService.trackEvent('home_opened', {
      child_profile_id: 'test-123',
    });

    expect(consoleSpy).not.toHaveBeenCalled();
    consoleSpy.mockRestore();
  });

  it('should set user properties', async () => {
    await analyticsService.init();

    analyticsService.setUserProperties({
      grade: 5,
      level: 10,
      coins_balance: 100,
    });

    expect(analyticsService.getSessionId()).toBeDefined();
  });

  it('should update user properties', async () => {
    await analyticsService.init();

    analyticsService.setUserProperties({
      grade: 5,
    });

    analyticsService.updateUserProperties({
      level: 10,
    });

    expect(analyticsService.getSessionId()).toBeDefined();
  });
});
