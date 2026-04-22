import { test as base } from '@playwright/test';

/**
 * Fixture для авторизованного пользователя
 */
export const test = base.extend({
  // Автоматически устанавливаем storage перед каждым тестом
  storageState: async ({ browser }, use) => {
    const context = await browser.newContext();
    const page = await context.newPage();

    // Имитируем VK storage
    await page.addInitScript(() => {
      // Mock VK Bridge
      (window as any).vkBridge = {
        send: (method: string, params?: any) => {
          console.log('[Mock VK Bridge]', method, params);

          // Mock responses
          if (method === 'VKWebAppStorageGet') {
            const keys = params?.keys || [];
            const items = keys.map((key: string) => {
              if (key === 'child_profile_id') {
                return { key, value: 'test-profile-id-123' };
              }
              if (key === 'platform_id') {
                return { key, value: 'vk' };
              }
              return { key, value: '' };
            });
            return Promise.resolve({ keys: items });
          }

          if (method === 'VKWebAppStorageSet') {
            return Promise.resolve({ result: true });
          }

          return Promise.resolve({});
        },
        subscribe: () => {},
        unsubscribe: () => {},
      };
    });

    // Устанавливаем storage state
    await page.evaluate(() => {
      localStorage.setItem('child_profile_id', 'test-profile-id-123');
      localStorage.setItem('platform_id', 'vk');
    });

    const storage = await page.context().storageState();
    await context.close();

    await use(storage);
  },
});

export { expect } from '@playwright/test';
