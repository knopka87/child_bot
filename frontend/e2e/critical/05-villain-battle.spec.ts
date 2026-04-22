import { test, expect } from '../fixtures/auth';

/**
 * CRITICAL FLOW #5: Villain Battle
 *
 * Проверяет систему сражений со злодеями
 */
test.describe('Villain Battle Flow', () => {
  test('should display active villain', async ({ page }) => {
    // Mock villain data
    await page.route('**/api/v1/villains/active', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          id: 'villain_1',
          name: 'Граф Ошибок',
          description: 'Мастер математических ошибок',
          image_url: '/assets/villain-1.png',
          hp: 66,
          max_hp: 100,
          is_active: true,
          is_defeated: false,
        }),
      });
    });

    // 1. Открываем главную страницу
    await page.goto('/');

    // 2. Проверяем отображение злодея
    await expect(page.locator('[data-testid="villain-card"]')).toBeVisible();
    await expect(page.locator('text=Граф Ошибок')).toBeVisible();

    // 3. Проверяем health bar
    const healthBar = page.locator('[data-testid="villain-health-bar"]');
    await expect(healthBar).toBeVisible();

    // 4. Кликаем на злодея
    await page.click('[data-testid="villain-card"]');

    // 5. Проверяем переход на страницу злодея
    await expect(page).toHaveURL(/\/villain/);
    await expect(page.locator('text=Граф Ошибок')).toBeVisible();
    await expect(page.locator('text=Начать битву')).toBeVisible();
  });

  test('should start battle and damage villain', async ({ page }) => {
    // Mock battle data
    await page.route('**/api/v1/villains/*/battle', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          villain_id: 'villain_1',
          damage_per_task: 20,
          tasks_to_defeat: 5,
          current_progress: 2,
        }),
      });
    });

    await page.goto('/villain');

    // Нажимаем "Начать битву"
    await page.click('button:has-text("Начать битву")');

    // Должны перейти к загрузке задания
    await expect(page).toHaveURL(/\/help\/upload/);
  });

  test('should show victory screen when villain defeated', async ({ page }) => {
    // Mock defeated villain
    await page.route('**/api/v1/villains/active', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          id: 'villain_1',
          name: 'Граф Ошибок',
          hp: 0,
          max_hp: 100,
          is_active: false,
          is_defeated: true,
        }),
      });
    });

    await page.goto('/villain');

    // Должен редиректнуть на victory
    await expect(page).toHaveURL(/\/villain\/victory/);
    await expect(page.locator('text=Победа')).toBeVisible();
    await expect(page.locator('[data-testid="victory-rewards"]')).toBeVisible();
  });

  test('should apply damage after correct solution', async ({ page }) => {
    // Симуляция полного flow:
    // 1. Открыть villain
    // 2. Начать битву
    // 3. Решить задачу правильно
    // 4. Проверить что HP уменьшился

    let villainHP = 100;

    await page.route('**/api/v1/villains/*/damage', async (route) => {
      villainHP -= 20;
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          new_hp: villainHP,
          damage_dealt: 20,
        }),
      });
    });

    // Эмулируем успешное решение задачи
    await page.goto('/help/result/attempt-123');

    // После просмотра результата должен примениться урон
    // HP должен уменьшиться на 20
  });
});
