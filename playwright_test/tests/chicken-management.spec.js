import { test, expect } from '@playwright/test';
import { TestHelpers } from '../utils/test-helpers.js';
import { testData } from '../utils/test-data.js';

test.describe('Управление курами', () => {
  let helpers;

  test.beforeEach(async ({ page }) => {
    helpers = new TestHelpers(page);
    await page.goto('/chickens');
    await helpers.waitForPageLoad();
  });


  test('корректно закрывает модальное окно по кнопке "Отмена"', async ({ page }) => {
    await page.click('button:has-text("Добавить курицу")');
    await helpers.waitForModal();

    await page.fill('[name="breed"]', 'Тестовая порода');

    await page.click('button:has-text("Отмена")');

    await helpers.waitForModalClose();
  });

  test('отображает пустое состояние когда нет кур', async ({ page }) => {
    await page.route('**/api/chickens', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([])
      });
    });

    await page.goto('/chickens');
    await helpers.waitForPageLoad();

    await expect(page.locator('text=Пока нет ни одной курицы')).toBeVisible();
  });
});
