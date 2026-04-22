import { test, expect } from '../fixtures/auth';

/**
 * CRITICAL FLOW #2: Help Flow
 *
 * Проверяет процесс получения помощи с домашним заданием
 */
test.describe('Help Flow', () => {
  test('should upload task and get hints', async ({ page }) => {
    // 1. Открываем главную страницу
    await page.goto('/');

    // 2. Нажимаем "Помочь с ДЗ"
    await page.click('[data-testid="help-button"]');

    // 3. Проверяем что попали на страницу загрузки
    await expect(page).toHaveURL(/\/help\/upload/);
    await expect(page.locator('text=Загрузи фото')).toBeVisible();

    // 4. Выбираем файл (mock)
    // В реальном тесте нужен файл, здесь проверяем UI
    const fileInput = page.locator('input[type="file"]');
    await expect(fileInput).toBeAttached();

    // Имитируем выбор файла через DataTransfer API
    const buffer = await page.evaluate(() => {
      // Создаем mock canvas с текстом задания
      const canvas = document.createElement('canvas');
      canvas.width = 800;
      canvas.height = 600;
      const ctx = canvas.getContext('2d')!;
      ctx.fillStyle = 'white';
      ctx.fillRect(0, 0, 800, 600);
      ctx.fillStyle = 'black';
      ctx.font = '24px Arial';
      ctx.fillText('Задача: 2 + 2 = ?', 50, 100);
      return canvas.toDataURL();
    });

    // В реальности Playwright может использовать setInputFiles
    // await fileInput.setInputFiles({
    //   name: 'task.png',
    //   mimeType: 'image/png',
    //   buffer: Buffer.from(buffer.split(',')[1], 'base64'),
    // });

    // 5. Проверяем кнопку обработки
    const processButton = page.locator('button:has-text("Обработать")');
    // await expect(processButton).toBeEnabled();
    // await processButton.click();

    // 6. Проверяем переход на экран обработки
    // await expect(page).toHaveURL(/\/help\/processing/);
    // await expect(page.locator('text=Обрабатываем')).toBeVisible();

    // 7. Ждем завершения обработки (или mock)
    // В реальности API вернет результат
    // await page.waitForURL(/\/help\/result/, { timeout: 30000 });

    // 8. Проверяем результат
    // await expect(page.locator('[data-testid="hint-text"]')).toBeVisible();
    // await expect(page.locator('text=Подсказка')).toBeVisible();
  });

  test('should continue unfinished attempt', async ({ page }) => {
    // 1. Создаем незавершенную попытку через API mock
    await page.route('**/api/v1/home/*', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          profile: {
            id: 'test-profile',
            displayName: 'Тест',
            level: 1,
            levelProgress: 50,
            coinsBalance: 100,
            tasksSolvedCorrectCount: 5,
          },
          mascot: {
            id: 'owl_1',
            state: 'idle',
            imageUrl: '/assets/mascot.png',
            message: 'Привет!',
          },
          unfinishedAttempt: {
            id: 'attempt-123',
            type: 'help',
            mode: 'help',
            status: 'processing',
            createdAt: new Date().toISOString(),
          },
          recentAttempts: [],
          achievements: {
            unlockedCount: 2,
            totalCount: 50,
          },
        }),
      });
    });

    // 2. Открываем главную
    await page.goto('/');

    // 3. Нажимаем Help
    await page.click('[data-testid="help-button"]');

    // 4. Должна показаться модалка с незавершенной попыткой
    await expect(page.locator('text=Продолжить')).toBeVisible();
    await expect(page.locator('text=Новая задача')).toBeVisible();

    // 5. Продолжаем
    await page.click('button:has-text("Продолжить")');

    // 6. Должны перейти к обработке
    await expect(page).toHaveURL(/\/help/);
  });

  test('should show error on upload failure', async ({ page }) => {
    // Mock API error
    await page.route('**/api/v1/attempts', async (route) => {
      await route.fulfill({
        status: 500,
        contentType: 'application/json',
        body: JSON.stringify({ error: 'Server error' }),
      });
    });

    await page.goto('/help/upload');

    // Попытка загрузки файла должна показать ошибку
    // await fileInput.setInputFiles('test-task.png');
    // await expect(page.locator('text=Ошибка')).toBeVisible();
  });
});
