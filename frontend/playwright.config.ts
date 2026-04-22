import { defineConfig, devices } from '@playwright/test';

/**
 * E2E тесты для критических flow
 * Запуск: npm run test:e2e
 */
export default defineConfig({
  testDir: './e2e',

  /* Максимальное время выполнения теста */
  timeout: 30 * 1000,

  /* Параллельное выполнение тестов */
  fullyParallel: true,

  /* Не падать на первой ошибке */
  forbidOnly: !!process.env.CI,

  /* Retry on CI только */
  retries: process.env.CI ? 2 : 0,

  /* Параллельность в CI */
  workers: process.env.CI ? 1 : undefined,

  /* Reporter */
  reporter: [
    ['html'],
    ['list'],
    ['json', { outputFile: 'test-results/results.json' }]
  ],

  /* Shared settings */
  use: {
    /* Base URL */
    baseURL: process.env.PLAYWRIGHT_BASE_URL || 'http://localhost:5173',

    /* Collect trace on failure */
    trace: 'on-first-retry',

    /* Screenshot on failure */
    screenshot: 'only-on-failure',

    /* Video on failure */
    video: 'retain-on-failure',

    /* Timeout для actions */
    actionTimeout: 10 * 1000,
  },

  /* Configure projects for major browsers */
  projects: [
    {
      name: 'chromium',
      use: {
        ...devices['Desktop Chrome'],
        viewport: { width: 393, height: 852 }, // iPhone 14 Pro size
      },
    },

    {
      name: 'firefox',
      use: {
        ...devices['Desktop Firefox'],
        viewport: { width: 393, height: 852 },
      },
    },

    {
      name: 'webkit',
      use: {
        ...devices['Desktop Safari'],
        viewport: { width: 393, height: 852 },
      },
    },

    /* Mobile viewports */
    {
      name: 'Mobile Chrome',
      use: { ...devices['Pixel 5'] },
    },
    {
      name: 'Mobile Safari',
      use: { ...devices['iPhone 13 Pro'] },
    },
  ],

  /* Run dev server before tests */
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:5173',
    reuseExistingServer: !process.env.CI,
    timeout: 120 * 1000,
  },
});
