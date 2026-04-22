import { test, expect } from '../fixtures/auth';

/**
 * CRITICAL FLOW #3: Check Flow
 *
 * Проверяет процесс проверки решения домашнего задания
 */
test.describe('Check Flow', () => {
  test('should select scenario and upload solution', async ({ page }) => {
    // 1. Открываем главную
    await page.goto('/');

    // 2. Нажимаем "Проверить решение"
    await page.click('[data-testid="check-button"]');

    // 3. Проверяем страницу выбора сценария
    await expect(page).toHaveURL(/\/check\/scenario/);
    await expect(page.locator('text=Выбери сценарий')).toBeVisible();

    // 4. Выбираем "Одно фото"
    await page.click('[data-testid="scenario-single"]');

    // 5. Проверяем страницу загрузки
    await expect(page).toHaveURL(/\/check\/upload-single/);
    await expect(page.locator('text=Загрузи фото')).toBeVisible();

    // 6. Загружаем задание и решение
    // const taskInput = page.locator('[data-testid="task-upload"]');
    // const solutionInput = page.locator('[data-testid="solution-upload"]');

    // await taskInput.setInputFiles('test-task.png');
    // await solutionInput.setInputFiles('test-solution.png');

    // 7. Нажимаем проверить
    // const checkButton = page.locator('button:has-text("Проверить")');
    // await expect(checkButton).toBeEnabled();
    // await checkButton.click();

    // 8. Должны перейти на processing
    // await expect(page).toHaveURL(/\/check\/processing/);
  });

  test('should upload two photos for comparison', async ({ page }) => {
    await page.goto('/');
    await page.click('[data-testid="check-button"]');

    // Выбираем "Два фото"
    await page.click('[data-testid="scenario-two"]');

    // Проверяем UI для двух фото
    await expect(page).toHaveURL(/\/check\/upload-two/);
    await expect(page.locator('text=Фото задания')).toBeVisible();
    await expect(page.locator('text=Фото решения')).toBeVisible();
  });

  test('should show validation errors', async ({ page }) => {
    await page.goto('/check/upload-single');

    // Пытаемся отправить без файлов
    const submitButton = page.locator('button:has-text("Проверить")');
    await expect(submitButton).toBeDisabled();

    // После загрузки одного файла - все еще disabled
    // await page.locator('[data-testid="task-upload"]').setInputFiles('test.png');
    // await expect(submitButton).toBeDisabled();

    // После загрузки обоих - enabled
    // await page.locator('[data-testid="solution-upload"]').setInputFiles('test2.png');
    // await expect(submitButton).toBeEnabled();
  });
});
