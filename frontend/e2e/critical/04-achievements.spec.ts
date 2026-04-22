import { test, expect } from '../fixtures/auth';

/**
 * CRITICAL FLOW #4: Achievements
 *
 * Проверяет систему достижений и наград
 */
test.describe('Achievements Flow', () => {
  test('should display achievements page', async ({ page }) => {
    // Mock achievements data
    await page.route('**/api/v1/achievements*', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([
          {
            id: 'first_task',
            title: 'Первая задача',
            description: 'Решил первую задачу',
            icon_url: '/assets/achievement-1.png',
            is_unlocked: true,
            unlocked_at: new Date().toISOString(),
            xp_reward: 10,
            coins_reward: 5,
          },
          {
            id: 'streak_3',
            title: 'Серия из 3',
            description: 'Решил 3 задачи подряд',
            icon_url: '/assets/achievement-2.png',
            is_unlocked: false,
            progress: 2,
            progress_max: 3,
          },
        ]),
      });
    });

    // 1. Открываем страницу достижений
    await page.goto('/achievements');

    // 2. Проверяем заголовок и статистику
    await expect(page.locator('text=Достижения')).toBeVisible();
    await expect(page.locator('text=Разблокировано')).toBeVisible();

    // 3. Проверяем отображение достижений
    await expect(page.locator('text=Первая задача')).toBeVisible();
    await expect(page.locator('text=Серия из 3')).toBeVisible();

    // 4. Кликаем на разблокированное достижение
    await page.click('text=Первая задача');

    // 5. Проверяем модальное окно с деталями
    await expect(page.locator('[data-testid="achievement-modal"]')).toBeVisible();
    await expect(page.locator('text=+10 XP')).toBeVisible();
    await expect(page.locator('text=+5')).toBeVisible(); // coins
  });

  test('should show locked achievements', async ({ page }) => {
    await page.goto('/achievements');

    // Заблокированные достижения должны быть серыми
    const lockedAchievement = page.locator('[data-achievement-locked="true"]').first();
    await expect(lockedAchievement).toHaveCSS('opacity', '0.6');
  });

  test('should display progress for serial achievements', async ({ page }) => {
    await page.goto('/achievements');

    // Проверяем прогресс-бар
    const progressBar = page.locator('[data-testid="achievement-progress"]').first();
    await expect(progressBar).toBeVisible();

    // Проверяем текст прогресса (например "2/3")
    await expect(page.locator('text=/\\d+\\/\\d+/')).toBeVisible();
  });
});
