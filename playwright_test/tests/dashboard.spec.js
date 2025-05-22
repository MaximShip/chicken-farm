import { test, expect } from '@playwright/test';
import { TestHelpers } from '../utils/test-helpers.js';

test.describe('Dashboard Page', () => {
  let helpers;

  test.beforeEach(async ({ page }) => {
    helpers = new TestHelpers(page);
    await page.goto('/');
    await helpers.waitForPageLoad();
  });

  test('отображает корректные статистические данные', async ({ page }) => {
    await helpers.waitForLoadingToFinish();

    const statValues = await page.locator('.stat-value').allTextContents();
    
    for (const value of statValues) {
      const numericValue = parseFloat(value.replace(/[^\d.-]/g, ''));
      if (!isNaN(numericValue)) {
        expect(numericValue).toBeGreaterThanOrEqual(0);
      }
    }
  });

  test('обрабатывает ошибки загрузки данных', async ({ page }) => {
    await page.route('**/api/chickens', async route => {
      await route.fulfill({
        status: 500,
        contentType: 'application/json',
        body: JSON.stringify({ error: 'Server Error' })
      });
    });

    await page.goto('/');
    await helpers.waitForPageLoad();

    await expect(page.locator('text=Ошибка при загрузке')).toBeVisible();
  });
});
