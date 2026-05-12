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
 * С кэшированием для избежания повторных запросов
 */
export async function getVKUserId(): Promise<number | null> {
  const VK_USER_ID_CACHE_KEY = 'vk_user_id_cache';

  // ВАЖНО: Сначала проверяем кэш
  const cachedUserId = sessionStorage.getItem(VK_USER_ID_CACHE_KEY);
  if (cachedUserId) {
    const userId = parseInt(cachedUserId, 10);
    if (!isNaN(userId)) {
      console.log('[VK Auth] Using cached VK user ID:', userId);
      return userId;
    }
  }

  // Пробуем получить через VK Bridge
  const userInfo = await getVKUserInfo();
  if (userInfo?.id) {
    // Кэшируем для будущих вызовов
    sessionStorage.setItem(VK_USER_ID_CACHE_KEY, userInfo.id.toString());
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
      // Кэшируем для будущих вызовов
      sessionStorage.setItem(VK_USER_ID_CACHE_KEY, userId.toString());
      return userId;
    }
  }

  console.warn('[VK Auth] Failed to get VK user ID - neither from cache, Bridge nor URL params');
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
 * VK приглашения работают через официальный механизм VKWebAppShowInviteBox:
 * 1. Отправитель вызывает: VKWebAppShowInviteBox({ requestKey: 'CODE' })
 * 2. Получатель запускает приложение по приглашению
 * 3. VK передаёт параметр vk_request_key = "CODE" в Launch Params
 *
 * Документация: https://dev.vk.com/ru/games/promotion/game-mechanics/invites
 */
/**
 * Служебные значения VK, которые не являются реферальными кодами
 * Эти значения могут приходить в Launch Params, но их нужно игнорировать
 */
const VK_SERVICE_VALUES = ['other', 'recs', 'recommendations', 'null', 'undefined', ''];

/**
 * Проверяет, является ли значение валидным реферальным кодом
 */
function isValidReferralCode(code: string | null | undefined): boolean {
  if (!code) return false;
  const normalized = code.toLowerCase().trim();
  return !VK_SERVICE_VALUES.includes(normalized);
}

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

    // DEBUG: Выводим ВСЕ Launch Params для диагностики
    console.log('[VK Auth] ===== LAUNCH PARAMS DEBUG =====');
    console.log('[VK Auth] Launch Params keys:', Object.keys(launchParams));
    console.log('[VK Auth] Full Launch Params:', JSON.stringify(launchParams, null, 2));
    console.log('[VK Auth] vk_request_key:', (launchParams as any).vk_request_key);
    console.log('[VK Auth] vk_ref:', (launchParams as any).vk_ref);
    console.log('[VK Auth] vk_start:', (launchParams as any).vk_start);
    console.log('[VK Auth] vk_fragment:', (launchParams as any).vk_fragment);
    console.log('[VK Auth] =====================================');

    // ПРИОРИТЕТ 1: vk_request_key - от VKWebAppShowInviteBox
    // Это ЕДИНСТВЕННЫЙ способ передать данные в VK Mini Apps (iframe)!
    // Документация: https://dev.vk.com/ru/bridge/VKWebAppShowInviteBox
    const requestKey = (launchParams as any).vk_request_key;
    if (isValidReferralCode(requestKey)) {
      console.log('[VK Auth] ✅ Referral code found in vk_request_key:', requestKey);
      return requestKey;
    }

    // ПРИОРИТЕТ 2: vk_ref - НЕ работает в iframe, но может работать в других случаях
    const refParam = (launchParams as any).vk_ref;
    if (isValidReferralCode(refParam)) {
      console.log('[VK Auth] ✅ Referral code found in vk_ref:', refParam);
      return refParam;
    }

    // ПРИОРИТЕТ 3: vk_fragment (часть URL после #)
    const fragmentParam = (launchParams as any).vk_fragment;
    if (fragmentParam && fragmentParam !== 'null' && fragmentParam !== 'undefined') {
      console.log('[VK Auth] Found vk_fragment:', fragmentParam);

      // Парсим fragment: может быть "start=CODE" или просто "CODE"
      let extractedCode: string | null = null;
      if (fragmentParam.startsWith('start=')) {
        extractedCode = fragmentParam.substring(6); // Убираем "start="
      } else if (fragmentParam.startsWith('ref=')) {
        extractedCode = fragmentParam.substring(4); // Убираем "ref="
      } else {
        // Если fragment не содержит =, считаем что это сам код
        extractedCode = fragmentParam;
      }

      if (isValidReferralCode(extractedCode)) {
        console.log('[VK Auth] ✅ Referral code extracted from vk_fragment:', extractedCode);
        return extractedCode;
      } else {
        console.log('[VK Auth] ⚠️ Service value ignored from vk_fragment:', extractedCode);
      }
    }

    // ПРИОРИТЕТ 4: vk_start (альтернативный способ)
    const startParam = (launchParams as any).vk_start;
    if (isValidReferralCode(startParam)) {
      console.log('[VK Auth] ✅ Referral code found in vk_start:', startParam);
      return startParam;
    }

    console.log('[VK Auth] ⚠️ Referral code not found in Launch Params');
    return null;
  } catch (error) {
    console.error('[VK Auth] Failed to get referral code:', error);
    return null;
  }
}
