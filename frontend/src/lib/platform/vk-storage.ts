// src/lib/platform/vk-storage.ts
import bridge from './bridge';

// Проверяем режим разработки
const isDev = import.meta.env.DEV;
const isInsideVK = typeof window !== 'undefined' && window.location.search.includes('vk_');
const useLocalStorageOnly = isDev && !isInsideVK;

export const vkStorage = {
  /**
   * Получить значение по ключу
   */
  async getItem(key: string): Promise<string | null> {
    // В dev режиме вне VK сразу используем localStorage
    if (useLocalStorageOnly) {
      return localStorage.getItem(key);
    }

    try {
      // Добавляем timeout на VK Bridge запросы
      const storagePromise = bridge.send('VKWebAppStorageGet', { keys: [key] });
      const timeoutPromise = new Promise<never>((_, reject) =>
        setTimeout(() => reject(new Error('VKWebAppStorageGet timeout')), 1000)
      );

      const data = await Promise.race([storagePromise, timeoutPromise]);
      return data.keys[0]?.value || null;
    } catch (error) {
      console.warn('[VKStorage] Fallback to localStorage:', error);
      return localStorage.getItem(key);
    }
  },

  /**
   * Сохранить значение
   */
  async setItem(key: string, value: string): Promise<void> {
    // В dev режиме вне VK сразу используем localStorage
    if (useLocalStorageOnly) {
      localStorage.setItem(key, value);
      console.log('[VKStorage] Saved to localStorage (dev mode):', key, value.substring(0, 50));
      return;
    }

    try {
      const setPromise = bridge.send('VKWebAppStorageSet', { key, value });
      const timeoutPromise = new Promise<never>((_, reject) =>
        setTimeout(() => reject(new Error('VKWebAppStorageSet timeout')), 1000)
      );
      await Promise.race([setPromise, timeoutPromise]);
      // Дублируем в localStorage для надежности
      localStorage.setItem(key, value);
    } catch (error) {
      console.warn('[VKStorage] Fallback to localStorage:', error);
      localStorage.setItem(key, value);
    }
  },

  /**
   * Удалить значение
   */
  async removeItem(key: string): Promise<void> {
    // В dev режиме вне VK сразу используем localStorage
    if (useLocalStorageOnly) {
      localStorage.removeItem(key);
      return;
    }

    try {
      const removePromise = bridge.send('VKWebAppStorageSet', { key, value: '' });
      const timeoutPromise = new Promise<never>((_, reject) =>
        setTimeout(() => reject(new Error('VKWebAppStorageSet timeout')), 1000)
      );
      await Promise.race([removePromise, timeoutPromise]);
      // Дублируем в localStorage
      localStorage.removeItem(key);
    } catch (error) {
      console.warn('[VKStorage] Fallback to localStorage:', error);
      localStorage.removeItem(key);
    }
  },

  /**
   * Получить несколько значений сразу
   */
  async getItems(keys: string[]): Promise<Record<string, string>> {
    // В dev режиме вне VK сразу используем localStorage
    if (useLocalStorageOnly) {
      const result: Record<string, string> = {};
      for (const key of keys) {
        const value = localStorage.getItem(key);
        if (value) result[key] = value;
      }
      return result;
    }

    try {
      const getPromise = bridge.send('VKWebAppStorageGet', { keys });
      const timeoutPromise = new Promise<never>((_, reject) =>
        setTimeout(() => reject(new Error('VKWebAppStorageGet timeout')), 1000)
      );

      const data = await Promise.race([getPromise, timeoutPromise]);
      const result: Record<string, string> = {};
      for (const item of data.keys) {
        if (item.value) {
          result[item.key] = item.value;
        }
      }
      return result;
    } catch (error) {
      console.warn('[VKStorage] Fallback to localStorage:', error);
      const result: Record<string, string> = {};
      for (const key of keys) {
        const value = localStorage.getItem(key);
        if (value) result[key] = value;
      }
      return result;
    }
  },
};

// Typed wrappers для типобезопасности
export const storageKeys = {
  TOKEN: 'auth_token',
  USER_ID: 'user_id',
  PROFILE_ID: 'profile_id',
  ONBOARDING_COMPLETED: 'onboarding_completed',
  REFERRAL_CODE: 'referral_code',
  // Onboarding progress
  ONBOARDING_STEP: 'onboarding_step',
  ONBOARDING_GRADE: 'onboarding_grade',
  ONBOARDING_AVATAR: 'onboarding_avatar',
  ONBOARDING_EMAIL: 'onboarding_email',
  ONBOARDING_EMAIL_VERIFIED: 'onboarding_email_verified',
  ONBOARDING_DISPLAY_NAME: 'onboarding_display_name',
  ONBOARDING_CONSENTS: 'onboarding_consents',
} as const;

/**
 * Type-safe wrapper для работы с storage
 */
export async function getStorageValue<T>(
  key: string,
  defaultValue: T
): Promise<T> {
  const value = await vkStorage.getItem(key);
  if (!value) return defaultValue;

  try {
    return JSON.parse(value) as T;
  } catch {
    return value as T;
  }
}

export async function setStorageValue<T>(key: string, value: T): Promise<void> {
  const stringValue = typeof value === 'string' ? value : JSON.stringify(value);
  await vkStorage.setItem(key, stringValue);
}
