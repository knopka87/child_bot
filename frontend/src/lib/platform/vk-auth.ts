// src/lib/platform/vk-auth.ts
/**
 * Утилиты для работы с VK авторизацией
 */
import bridge from './bridge';
import { isDevMode, isVKBridgeReady, initVKBridge } from './bridge';

export interface VKUserInfo {
  id: number;
  first_name: string;
  last_name: string;
  photo_100?: string;
  photo_200?: string;
}

export interface VKLaunchParams {
  vk_user_id?: string;
  vk_app_id?: string;
  vk_is_app_user?: string;
  vk_are_notifications_enabled?: string;
  vk_language?: string;
  vk_ref?: string;
  vk_access_token_settings?: string;
  vk_group_id?: string;
  vk_viewer_group_role?: string;
  sign?: string;
}

/**
 * Получить информацию о текущем пользователе VK
 */
export async function getVKUserInfo(): Promise<VKUserInfo | null> {
  try {
    // Ждём инициализации bridge
    if (!isVKBridgeReady()) {
      await initVKBridge();
    }

    // В dev режиме возвращаем mock данные
    if (isDevMode()) {
      console.log('[VK Auth] Dev mode: returning mock user info');
      return {
        id: 123456789,
        first_name: 'Тест',
        last_name: 'Пользователь',
        photo_100: 'https://via.placeholder.com/100',
        photo_200: 'https://via.placeholder.com/200',
      };
    }

    // Запрашиваем данные пользователя
    const userInfo = await bridge.send('VKWebAppGetUserInfo');
    console.log('[VK Auth] User info received:', userInfo);

    return userInfo as VKUserInfo;
  } catch (error) {
    console.error('[VK Auth] Failed to get user info:', error);
    return null;
  }
}

/**
 * Получить VK user ID текущего пользователя
 */
export async function getVKUserId(): Promise<number | null> {
  const userInfo = await getVKUserInfo();
  return userInfo?.id ?? null;
}

/**
 * Получить launch параметры из URL
 * Эти параметры используются для валидации sign на backend
 */
export function getVKLaunchParams(): VKLaunchParams {
  if (typeof window === 'undefined') {
    return {};
  }

  const params = new URLSearchParams(window.location.search);
  const vkParams: VKLaunchParams = {};

  // Собираем все vk_* параметры
  params.forEach((value, key) => {
    if (key.startsWith('vk_') || key === 'sign') {
      vkParams[key as keyof VKLaunchParams] = value;
    }
  });

  return vkParams;
}

/**
 * Получить sign параметр для валидации на backend
 */
export function getVKSign(): string | null {
  const params = getVKLaunchParams();
  return params.sign || null;
}

/**
 * Получить все VK параметры как query string для отправки на backend
 */
export function getVKParamsQueryString(): string {
  const params = getVKLaunchParams();
  const searchParams = new URLSearchParams();

  Object.entries(params).forEach(([key, value]) => {
    if (value) {
      searchParams.append(key, value);
    }
  });

  return searchParams.toString();
}

/**
 * Проверить что приложение запущено из VK
 */
export function isLaunchedFromVK(): boolean {
  const params = getVKLaunchParams();
  return !!(params.vk_user_id && params.vk_app_id);
}
