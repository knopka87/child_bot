/**
 * Утилиты для работы с аутентификацией и профилем пользователя
 *
 * TODO: Заменить на реальную систему аутентификации:
 * - Telegram WebApp initData
 * - VK Mini Apps bridge
 * - Или обычная JWT аутентификация
 */

/**
 * Получает ID профиля текущего пользователя
 *
 * @returns {string | null} UUID профиля или null если не авторизован
 *
 * TODO: Реализовать реальное получение профиля:
 * 1. Для Telegram: из window.Telegram.WebApp.initDataUnsafe
 * 2. Для VK: из vkBridge.send('VKWebAppGetUserInfo')
 * 3. Для Web: из localStorage после логина или из JWT токена
 */
export function getCurrentChildProfileId(): string | null {
  // ВРЕМЕННОЕ РЕШЕНИЕ: возвращаем тестовый профиль из БД
  // В продакшене это должно получаться из:
  // - localStorage.getItem('childProfileId')
  // - или из контекста приложения
  // - или из Telegram/VK WebApp данных

  return "cb569db4-cee2-438e-9ab1-90101901abb0"; // Друг 1 из БД
}

/**
 * Сохраняет ID профиля при входе в приложение
 *
 * @param {string} profileId - UUID профиля
 */
export function setCurrentChildProfileId(profileId: string): void {
  // TODO: Сохранять в localStorage или контекст
  // localStorage.setItem('childProfileId', profileId);
  console.warn('setCurrentChildProfileId not implemented, profileId:', profileId);
}

/**
 * Очищает данные профиля при выходе
 */
export function clearCurrentChildProfile(): void {
  // TODO: Очистить localStorage и контекст
  // localStorage.removeItem('childProfileId');
  console.warn('clearCurrentChildProfile not implemented');
}

/**
 * Проверяет, авторизован ли пользователь
 *
 * @returns {boolean}
 */
export function isAuthenticated(): boolean {
  // TODO: Реальная проверка авторизации
  return getCurrentChildProfileId() !== null;
}
