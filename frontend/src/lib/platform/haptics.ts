// src/lib/platform/haptics.ts
import bridge from './bridge';

export const haptics = {
  /**
   * Impact vibration - для кликов и жестов
   */
  impact(style: 'light' | 'medium' | 'heavy' = 'medium'): void {
    bridge
      .send('VKWebAppTapticImpactOccurred', { style })
      .catch((error) => {
        console.debug('[Haptics] Impact not available:', error);
      });
  },

  /**
   * Notification vibration - для результатов операций
   */
  notification(type: 'error' | 'success' | 'warning'): void {
    bridge
      .send('VKWebAppTapticNotificationOccurred', { type })
      .catch((error) => {
        console.debug('[Haptics] Notification not available:', error);
      });
  },

  /**
   * Selection vibration - для переключения между элементами
   */
  selection(): void {
    bridge
      .send('VKWebAppTapticSelectionChanged', {})
      .catch((error) => {
        console.debug('[Haptics] Selection not available:', error);
      });
  },
};

// Использование в компонентах
export function useHaptics() {
  return {
    /**
     * Легкая вибрация для кнопок
     */
    onButtonClick: () => haptics.impact('light'),

    /**
     * Средняя вибрация для важных действий
     */
    onImportantAction: () => haptics.impact('medium'),

    /**
     * Сильная вибрация для критичных действий
     */
    onCriticalAction: () => haptics.impact('heavy'),

    /**
     * Вибрация успеха
     */
    onSuccess: () => haptics.notification('success'),

    /**
     * Вибрация ошибки
     */
    onError: () => haptics.notification('error'),

    /**
     * Вибрация предупреждения
     */
    onWarning: () => haptics.notification('warning'),

    /**
     * Вибрация при переключении
     */
    onSelect: () => haptics.selection(),
  };
}
