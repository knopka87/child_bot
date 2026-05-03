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
 * С fallback на URL параметр vk_user_id если VK Bridge не работает
 */
export async function getVKUserId(): Promise<number | null> {
  // Пробуем получить через VK Bridge
  const userInfo = await getVKUserInfo();
  if (userInfo?.id) {
    return userInfo.id;
  }

  // FALLBACK: Читаем vk_user_id из URL параметров
  // Это работает когда приложение открывается по ссылке с VK параметрами,
  // но VK Bridge ещё не инициализирован или не работает (не в VK iframe)
  const params = getVKLaunchParams();
  if (params.vk_user_id) {
    const userId = parseInt(params.vk_user_id, 10);
    if (!isNaN(userId)) {
      console.log('[VK Auth] Using vk_user_id from URL params:', userId);
      return userId;
    }
  }

  console.warn('[VK Auth] Failed to get VK user ID - neither from Bridge nor URL params');
  return null;
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

/**
 * Получить реферальный код из VK Launch Params
 *
 * VK Mini Apps механизм передачи данных:
 * - URL: https://vk.com/app123#start=CODE (формат без _param!)
 * - Launch Params содержат: vk_start_param = "CODE"
 *
 * Это ЕДИНСТВЕННЫЙ способ передать данные в VK Mini App!
 * Документация: https://dev.vk.com/ru/mini-apps/development/launch-params
 */
export async function getVKRefCode(): Promise<string | null> {
  try {
    // Главный источник: VK Bridge Launch Params
    if (!isVKBridgeReady()) {
      await initVKBridge();
    }

    // В dev режиме возвращаем null
    if (isDevMode()) {
      console.log('[VK Auth] Dev mode: referral code not available');
      return null;
    }

    const launchParams = await bridge.send('VKWebAppGetLaunchParams');

    // DEBUG: Выводим ВСЕ ключи Launch Params
    console.log('[VK Auth] Launch Params keys:', Object.keys(launchParams));
    console.log('[VK Auth] Full Launch Params:', launchParams);

    // ГЛАВНОЕ: vk_start_param содержит наш реферальный код
    // Формат URL: https://vk.com/app54517931#start=CODE
    const startParam = (launchParams as any).vk_start_param;

    if (startParam) {
      console.log('[VK Auth] ✅ Referral code found in vk_start_param:', startParam);
      return startParam;
    }

    // Fallback: проверяем window.location.hash (для прямого доступа не через VK)
    if (typeof window !== 'undefined' && window.location.hash) {
      const hash = window.location.hash.substring(1); // убираем #

      // Проверяем формат #start=CODE (официальный VK формат)
      if (hash.startsWith('start=')) {
        const code = hash.substring('start='.length);
        console.log('[VK Auth] Referral code found in window.location.hash:', code);
        return code;
      }

      // Или query-string формат в hash
      const hashParams = new URLSearchParams(hash);
      const hashRef = hashParams.get('start') || hashParams.get('ref');

      if (hashRef) {
        console.log('[VK Auth] Referral code found in hash params:', hashRef);
        return hashRef;
      }
    }

    console.log('[VK Auth] ⚠️ Referral code not found');
    return null;
  } catch (error) {
    console.error('[VK Auth] Failed to get referral code:', error);
    return null;
  }
}
