// src/lib/platform/bridge.ts
/**
 * Централизованный экспорт VK Bridge
 * Автоматически выбирает mock или реальный bridge в зависимости от окружения
 */
import VKBridgeDefault from '@vkontakte/vk-bridge';
import VKBridgeMockDefault from '@vkontakte/vk-bridge-mock';

// Определяем окружение
const isDev = import.meta.env.DEV;
const isInsideVK = typeof window !== 'undefined' && window.location.search.includes('vk_');

// Выбираем bridge синхронно
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

// Инициализируем bridge
bridge.send('VKWebAppInit');

export default bridge;
