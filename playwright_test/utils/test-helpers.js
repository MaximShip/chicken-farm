export class TestHelpers {
  constructor(page) {
    this.page = page;
  }

  async waitForPageLoad() {
    await this.page.waitForLoadState('networkidle');
    await this.page.waitForTimeout(500);
  }


  async fillChickenForm(chickenData) {
    await this.page.fill('[name="cage_id"]', chickenData.cage_id.toString());
    await this.page.fill('[name="weight"]', chickenData.weight.toString());
    await this.page.fill('[name="age"]', chickenData.age.toString());
    await this.page.fill('[name="egg_per_month"]', chickenData.egg_per_month.toString());
    await this.page.fill('[name="breed"]', chickenData.breed);
  }

  async fillEmployeeForm(employeeData) {
    await this.page.fill('[name="full_name"]', employeeData.full_name);
    await this.page.fill('[name="passport_data"]', employeeData.passport_data);
    await this.page.fill('[name="salary"]', employeeData.salary.toString());
    if (employeeData.cages) {
      await this.page.fill('[name="cages"]', employeeData.cages);
    }
  }

  async waitForSuccessNotification(text) {
    const selector = `.success:has-text("${text}"), .alert-success:has-text("${text}")`;
    await this.page.waitForSelector(selector, { timeout: 15000 });
  }

  async waitForErrorNotification(text) {
    const selector = `.error:has-text("${text}"), .alert-error:has-text("${text}")`;
    await this.page.waitForSelector(selector, { timeout: 15000 });
  }

  async waitForModal() {
    await this.page.waitForSelector('.modal', { state: 'visible' });
  }

  async waitForModalClose() {
    await this.page.waitForSelector('.modal', { state: 'hidden' });
  }

  async takeScreenshot(name) {
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
    const path = `screenshots/${name}-${timestamp}.png`;
    await this.page.screenshot({ 
      path,
      fullPage: true 
    });
    console.log(`📸 Скриншот сохранен: ${path}`);
  }

  async waitForLoadingToFinish() {
    const loadingSelectors = [
      'text=Загрузка...',
      '.loading',
      '.spinner',
      '[data-testid="loading"]'
    ];

    for (const selector of loadingSelectors) {
      try {
        await this.page.waitForSelector(selector, { state: 'hidden', timeout: 5000 });
      } catch (e) {
    }
    }
  }

  async waitForTableData() {
    await this.page.waitForSelector('tbody tr', { timeout: 10000 });
  }

  async getTableRowCount() {
    return await this.page.locator('tbody tr').count();
  }

  async clickAndWait(selector, waitForSelector = null) {
    await this.page.click(selector);
    if (waitForSelector) {
      await this.page.waitForSelector(waitForSelector);
    }
    await this.page.waitForTimeout(300);
  }

  async fillField(selector, value) {
    await this.page.fill(selector, '');
    await this.page.fill(selector, value);
  }
}