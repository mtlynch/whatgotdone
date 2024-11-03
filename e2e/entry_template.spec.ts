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

  await expect(page.getByRole("textbox")).toContainText("");

  // Navigate to preferences page using the navigation menu.
  await page.getByRole("button", { name: "Account" }).click();
  await page.getByRole("menuitem", { name: "Preferences" }).click();
  await expect(page).toHaveURL("/preferences");

  await expect(page.getByRole("textbox")).toContainText("");

  // Set template text.
  await page
    .getByRole("textbox")
    .fill("# Example project\n\n* Item A\n* Item B");
  await page.getByRole("button", { name: /save/i }).click();

  // Verify success message.
  await expect(page.getByRole("alert")).toBeVisible();

  // Verify new entries start with the template.
  apiDraftGet = page.waitForRequest(
    (request) =>
      request.url().includes("/api/draft") && request.method() === "GET"
  );
  await page.goto("/entry/edit/2020-03-06");
  await apiDraftGet;

  await expect(page.getByRole("textbox")).toHaveValue(
    "# Example project\n\n* Item A\n* Item B"
  );
});
