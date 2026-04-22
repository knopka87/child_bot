import { test, expect } from '@playwright/test';

/**
 * CRITICAL FLOW #1: Onboarding
 *
 * Проверяет полный процесс онбординга нового пользователя
 */
test.describe('Onboarding Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Очищаем storage перед каждым тестом
    await page.context().clearCookies();
    await page.evaluate(() => localStorage.clear());
  });

  test('should complete full onboarding flow', async ({ page }) => {
    // 1. Открываем приложение
    await page.goto('/');

    // 2. Проверяем что показался Welcome экран
    await expect(page.locator('text=Добро пожаловать')).toBeVisible({ timeout: 10000 });

    // 3. Нажимаем "Начать"
    await page.click('button:has-text("Начать")');

    // 4. Вводим имя
    await expect(page.locator('text=Как тебя зовут')).toBeVisible();
    const nameInput = page.locator('input[type="text"]').first();
    await nameInput.fill('Тестовый Ученик');
    await page.click('button:has-text("Далее")');

    // 5. Выбираем класс
    await expect(page.locator('text=Выбери класс')).toBeVisible();
    await page.click('button:has-text("5 класс")');

    // 6. Выбираем аватар
    await expect(page.locator('text=Выбери аватар')).toBeVisible();
    // Ждем загрузки аватаров
    await page.waitForTimeout(1000);
    const firstAvatar = page.locator('[data-testid="avatar-option"]').first();
    await firstAvatar.click();
    await page.click('button:has-text("Далее")');

    // 7. Вводим email родителя
    await expect(page.locator('text=Email родителя')).toBeVisible();
    const emailInput = page.locator('input[type="email"]');
    await emailInput.fill('parent@example.com');
    await page.click('button:has-text("Далее")');

    // 8. Принимаем согласия
    await expect(page.locator('text=Согласие')).toBeVisible();

    // Отмечаем все чекбоксы
    const checkboxes = page.locator('input[type="checkbox"]');
    const count = await checkboxes.count();
    for (let i = 0; i < count; i++) {
      await checkboxes.nth(i).check();
    }

    await page.click('button:has-text("Принять")');

    // 9. Проверяем что попали на главную страницу
    await expect(page).toHaveURL(/\/home/, { timeout: 15000 });

    // 10. Проверяем основные элементы главной страницы
    await expect(page.locator('text=Тестовый Ученик')).toBeVisible();
    await expect(page.locator('[data-testid="help-button"]')).toBeVisible();
    await expect(page.locator('[data-testid="check-button"]')).toBeVisible();
  });

  test('should validate required fields', async ({ page }) => {
    await page.goto('/');
    await page.click('button:has-text("Начать")');

    // Пытаемся перейти дальше без ввода имени
    const nextButton = page.locator('button:has-text("Далее")');
    await expect(nextButton).toBeDisabled();

    // Вводим имя - кнопка активируется
    await page.locator('input[type="text"]').first().fill('Тест');
    await expect(nextButton).toBeEnabled();
  });

  test('should allow navigation back', async ({ page }) => {
    await page.goto('/');
    await page.click('button:has-text("Начать")');

    // Вводим имя и идем дальше
    await page.locator('input[type="text"]').first().fill('Тест');
    await page.click('button:has-text("Далее")');

    // Проверяем что на экране выбора класса
    await expect(page.locator('text=Выбери класс')).toBeVisible();

    // Нажимаем назад
    await page.click('[data-testid="back-button"]');

    // Проверяем что вернулись к вводу имени
    await expect(page.locator('text=Как тебя зовут')).toBeVisible();
    // И имя сохранилось
    await expect(page.locator('input[type="text"]').first()).toHaveValue('Тест');
  });
});
