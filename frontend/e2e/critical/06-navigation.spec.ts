import { test, expect } from '../fixtures/auth';

/**
 * CRITICAL FLOW #6: Navigation
 *
 * Проверяет навигацию между основными экранами
 */
test.describe('Navigation Flow', () => {
  test('should navigate between main tabs', async ({ page }) => {
    await page.goto('/');

    // 1. Проверяем что на Home
    await expect(page).toHaveURL(/\/home/);

    // 2. Переходим в Achievements
    await page.click('[data-testid="nav-achievements"]');
    await expect(page).toHaveURL(/\/achievements/);
    await expect(page.locator('text=Достижения')).toBeVisible();

    // 3. Переходим в Profile
    await page.click('[data-testid="nav-profile"]');
    await expect(page).toHaveURL(/\/profile/);
    await expect(page.locator('text=Профиль')).toBeVisible();

    // 4. Переходим в Friends
    await page.click('[data-testid="nav-friends"]');
    await expect(page).toHaveURL(/\/friends/);
    await expect(page.locator('text=Друзья')).toBeVisible();

    // 5. Возвращаемся на Home
    await page.click('[data-testid="nav-home"]');
    await expect(page).toHaveURL(/\/home/);
  });

  test('should navigate to history from profile', async ({ page }) => {
    await page.goto('/profile');

    // Кликаем на "История"
    await page.click('text=История');

    // Проверяем переход
    await expect(page).toHaveURL(/\/profile\/history/);
    await expect(page.locator('text=История')).toBeVisible();

    // Проверяем back button
    await page.click('[data-testid="back-button"]');
    await expect(page).toHaveURL(/\/profile/);
  });

  test('should navigate back from deep pages', async ({ page }) => {
    // Переход: Home → Help → Processing → Result
    await page.goto('/help/result/attempt-123');

    // Нажимаем "На главную"
    await page.click('button:has-text("На главную")');
    await expect(page).toHaveURL(/\/home/);
  });

  test('should preserve state on tab switch', async ({ page }) => {
    // 1. Открываем достижения, скроллим
    await page.goto('/achievements');
    await page.evaluate(() => window.scrollTo(0, 500));
    const scrollPos = await page.evaluate(() => window.scrollY);

    // 2. Переходим на другую вкладку
    await page.click('[data-testid="nav-home"]');

    // 3. Возвращаемся
    await page.click('[data-testid="nav-achievements"]');

    // 4. Проверяем что позиция скролла сохранилась
    const newScrollPos = await page.evaluate(() => window.scrollY);
    expect(newScrollPos).toBe(scrollPos);
  });

  test('should handle browser back/forward', async ({ page }) => {
    await page.goto('/');

    // Переходим на profile
    await page.click('[data-testid="nav-profile"]');
    await expect(page).toHaveURL(/\/profile/);

    // Используем browser back
    await page.goBack();
    await expect(page).toHaveURL(/\/home/);

    // Используем browser forward
    await page.goForward();
    await expect(page).toHaveURL(/\/profile/);
  });
});
