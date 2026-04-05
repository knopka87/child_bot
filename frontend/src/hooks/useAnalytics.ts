// src/hooks/useAnalytics.ts
import { useContext } from 'react';
import { AnalyticsContext } from '@/contexts/AnalyticsContext';
import type { AnalyticsEventName } from '@/types/analytics';

export function useAnalytics() {
  const analytics = useContext(AnalyticsContext);

  // Если analytics еще не инициализирован, возвращаем безопасные заглушки
  if (!analytics) {
    return {
      trackEvent: (_name: AnalyticsEventName, _params?: Record<string, any>) => {
        // No-op: аналитика еще не инициализирована
      },
      setUserProperties: (_properties: Record<string, any>) => {
        // No-op
      },
      updateUserProperties: (_properties: Record<string, any>) => {
        // No-op
      },
      sessionId: 'pending',
    };
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
