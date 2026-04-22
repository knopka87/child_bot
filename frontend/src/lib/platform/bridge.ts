// src/lib/platform/bridge.ts
/**
 * Централизованный экспорт VK Bridge с инициализацией
 * Автоматически выбирает mock или реальный bridge в зависимости от окружения
 */
import VKBridgeDefault from '@vkontakte/vk-bridge';
import VKBridgeMockDefault from '@vkontakte/vk-bridge-mock';

// Определяем окружение
const isDev = import.meta.env.DEV;
const isInsideVK = typeof window !== 'undefined' && window.location.search.includes('vk_');

// Выбираем bridge
let bridge: typeof VKBridgeDefault;

if (isDev && !isInsideVK) {
  // Используем mock для локальной разработки
  bridge = VKBridgeMockDefault;
  console.log('[VK Bridge] Using mock for local development');
} else {
  // Используем реальный VK Bridge
  bridge = VKBridgeDefault;
  console.log('[VK Bridge] Using real VK Bridge');
}

// Флаг инициализации
let isInitialized = false;
let initPromise: Promise<void> | null = null;

/**
 * Инициализация VK Bridge
 * Вызывается автоматически при импорте, но можно вызвать повторно
 */
export async function initVKBridge(): Promise<void> {
  if (isInitialized) {
    console.log('[VK Bridge] Already initialized');
    return;
  }

  if (initPromise) {
    return initPromise;
  }

  initPromise = (async () => {
    try {
      console.log('[VK Bridge] Initializing...');
      await bridge.send('VKWebAppInit');
      isInitialized = true;
      console.log('[VK Bridge] Initialized successfully');
    } catch (error) {
      console.error('[VK Bridge] Initialization failed:', error);
      throw error;
    }
  })();

  return initPromise;
}

/**
 * Проверка что VK Bridge готов к работе
 */
export function isVKBridgeReady(): boolean {
  return isInitialized;
}

/**
 * Проверка что приложение запущено внутри VK
 */
export function isRunningInsideVK(): boolean {
  return isInsideVK;
}

/**
 * Проверка что используется dev режим
 */
export function isDevMode(): boolean {
  return isDev && !isInsideVK;
}

// Автоматическая инициализация при загрузке модуля
initVKBridge().catch(err => {
  console.error('[VK Bridge] Auto-initialization failed:', err);
});

export default bridge;
