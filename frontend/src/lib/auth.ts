/**
 * Утилиты для работы с аутентификацией и профилем пользователя через VK Mini Apps
 */

import { getVKUserId, getVKParamsQueryString } from './platform/vk-auth';
import { isDevMode } from './platform/bridge';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';
const STORAGE_KEY = 'child_profile_id';
const VK_USER_ID_KEY = 'vk_user_id';

/**
 * Получает ID профиля текущего пользователя
 *
 * В production: получает через VK Bridge и запрашивает профиль с backend
 * В development: использует захардкоженный тестовый профиль
 *
 * @returns {Promise<string | null>} UUID профиля или null если не авторизован
 */
export async function getCurrentChildProfileId(): Promise<string | null> {
  // В dev режиме возвращаем тестовый профиль
  if (isDevMode()) {
    const devProfileId = "1c84c913-19b3-40f7-b3ab-a94f90ce374f"; // Web профиль из БД
    console.log('[Auth] Dev mode: using test profile:', devProfileId);
    return devProfileId;
  }

  try {
    // Проверяем кэш в sessionStorage (безопаснее чем localStorage)
    const cachedProfileId = sessionStorage.getItem(STORAGE_KEY);
    if (cachedProfileId) {
      console.log('[Auth] Using cached profile ID');
      return cachedProfileId;
    }

    // Получаем VK user ID через vk-auth
    const vkUserId = await getVKUserId();
    if (!vkUserId) {
      console.error('[Auth] Failed to get VK user ID');
      return null;
    }

    console.log('[Auth] VK user ID:', vkUserId);

    // Запрашиваем child_profile_id у backend
    const response = await fetch(
      `${API_BASE_URL}/profiles/by-platform?platform_id=vk&platform_user_id=${vkUserId}`,
      {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'X-Platform-ID': 'vk',
        },
      }
    );

    if (response.ok) {
      const data = await response.json();
      const profileId = data.child_profile_id;

      // Сохраняем в sessionStorage
      sessionStorage.setItem(STORAGE_KEY, profileId);
      sessionStorage.setItem(VK_USER_ID_KEY, vkUserId.toString());

      console.log('[Auth] Profile found:', profileId);
      return profileId;
    } else if (response.status === 404) {
      // Профиль не найден - пользователь должен пройти онбординг
      console.log('[Auth] Profile not found, user needs onboarding');
      return null;
    } else {
      console.error('[Auth] Failed to get profile:', response.status);
      return null;
    }
  } catch (error) {
    console.error('[Auth] Failed to get profile:', error);
    return null;
  }
}

/**
 * Синхронная версия для компонентов, которые не могут использовать async
 * Читает из sessionStorage
 *
 * @returns {string | null} UUID профиля из кэша или null
 */
export function getCurrentChildProfileIdSync(): string | null {
  if (isDevMode()) {
    return "1c84c913-19b3-40f7-b3ab-a94f90ce374f";
  }
  return sessionStorage.getItem(STORAGE_KEY);
}

/**
 * Сохраняет ID профиля после создания/онбординга
 *
 * @param {string} profileId - UUID профиля
 */
export function setCurrentChildProfileId(profileId: string): void {
  sessionStorage.setItem(STORAGE_KEY, profileId);
  console.log('[Auth] Profile ID saved:', profileId);
}

/**
 * Получить VK user ID из кэша (без нового запроса)
 */
export function getCachedVKUserId(): string | null {
  if (isDevMode()) {
    return '123456789';
  }
  return sessionStorage.getItem(VK_USER_ID_KEY);
}

/**
 * Очищает данные профиля при выходе
 */
export function clearCurrentChildProfile(): void {
  sessionStorage.removeItem(STORAGE_KEY);
  sessionStorage.removeItem(VK_USER_ID_KEY);
  console.log('[Auth] Profile data cleared');
}

/**
 * Проверяет, авторизован ли пользователь (асинхронная версия)
 *
 * @returns {Promise<boolean>}
 */
export async function isAuthenticated(): Promise<boolean> {
  const profileId = await getCurrentChildProfileId();
  return profileId !== null;
}

/**
 * Проверяет, авторизован ли пользователь (синхронная версия из кэша)
 *
 * @returns {boolean}
 */
export function isAuthenticatedSync(): boolean {
  return getCurrentChildProfileIdSync() !== null;
}

/**
 * Получить VK параметры для отправки на backend (для валидации sign)
 */
export function getVKAuthParams(): string {
  return getVKParamsQueryString();
}
