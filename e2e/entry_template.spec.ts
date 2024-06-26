import { expect, test } from "@playwright/test";
import { mockLoginAsUser, wipeDB } from "./helpers/test_apis.js";

test.beforeEach(async ({ page }) => {
  await wipeDB(page);
});

test("uses the entry template for new drafts", async ({ page }) => {
  let apiDraftGet = page.waitForRequest(
    (request) =>
      request.url().includes("/api/draft") && request.method() === "GET"
  );
  await mockLoginAsUser(page, "staging_jimmy");

  await expect(page).toHaveURL(/\/entry\/edit\/.+/g);

  // Wait for page to pull down any previous entry.
  await apiDraftGet;

  apiDraftGet = page.waitForRequest(
    (request) =>
      request.url().includes("/api/draft") && request.method() === "GET"
  );

  await page.goto("/entry/edit/2020-03-06");

  // Wait for page to pull down any previous entry.
  await apiDraftGet;

  expect(await page.locator(".editor-content .ProseMirror")).toContainText("");

  // Set an entry template on the preferences page.
  await page.getByTestId("account-dropdown").click();
  await page.getByTestId("preferences-link").click();
  await expect(page).toHaveURL("/preferences");

  await page
    .locator("id=entry-template-input")
    .fill("# Example project\n\n* Item A\n* Item B");
  await page.locator(".btn-primary").click();

  await expect(page.locator(".alert-success")).toBeVisible();

  // Verify new entries start with the template.

  apiDraftGet = page.waitForRequest(
    (request) =>
      request.url().includes("/api/draft") && request.method() === "GET"
  );
  await page.goto("/entry/edit/2020-03-06");
  await apiDraftGet;

  await page.locator(".switch-mode .btn").click();
  await expect(page.locator(".markdown-editor .editor-textarea")).toHaveValue(
    "# Example project\n\n* Item A\n* Item B"
  );
});
