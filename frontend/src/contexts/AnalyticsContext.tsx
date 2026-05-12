// src/contexts/AnalyticsContext.tsx
import { createContext, useEffect, useState, ReactNode } from 'react';
import { AnalyticsService } from '@/services/analytics/AnalyticsService';
import { PlatformBridge } from '@/services/platform/PlatformBridge';
import config from '@/config';
import type { AnalyticsConfig } from '@/types/analytics';
import type { PlatformType } from '@/types/platform';

const defaultConfig: AnalyticsConfig = {
  enabled: config.analytics.enabled,
  debug: config.analytics.debug,
  batchSize: config.analytics.batchSize,
  batchInterval: config.analytics.batchInterval,
  retryAttempts: config.analytics.retryAttempts,
  retryDelay: config.analytics.retryDelay,
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
  const [analytics, setAnalytics] = useState<AnalyticsService | null>(null);

  useEffect(() => {
    // Получаем тип платформы
    let platformType: PlatformType;
    try {
      const platformBridge = new PlatformBridge();
      platformType = platformBridge.getPlatformType();
    } catch (error) {
      // Если PlatformBridge выбросил ошибку (например, VK_ONLY_ACCESS)
      // используем fallback платформу и не инициализируем аналитику
      if (error instanceof Error && error.message.includes('VK_ONLY_ACCESS')) {
        console.log('[AnalyticsContext] VK_ONLY_ACCESS detected, skipping analytics initialization');
        // Не инициализируем аналитику для экрана редиректа
        return;
      }

      // Для других ошибок используем fallback
      console.error('[AnalyticsContext] PlatformBridge failed, using fallback:', error);
      platformType = 'vk';
    }

    const mergedConfig = { ...defaultConfig, ...config };
    const analyticsService = new AnalyticsService(mergedConfig, platformType);
    void analyticsService.init();
    setAnalytics(analyticsService);

    const handleBeforeUnload = () => {
      void analyticsService.flush();
    };

    window.addEventListener('beforeunload', handleBeforeUnload);

    return () => {
      window.removeEventListener('beforeunload', handleBeforeUnload);
      void analyticsService.flush();
      analyticsService.destroy();
    };
  }, []);

  return (
    <AnalyticsContext.Provider value={analytics}>
      {children}
    </AnalyticsContext.Provider>
  );
}
