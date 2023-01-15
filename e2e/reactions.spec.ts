import { expect, test } from "@playwright/test";
import {
  mockLoginAsUser,
  wipeDB,
  populateDummyData,
} from "./helpers/test_apis.js";

test.beforeEach(async ({ page }) => {
  await wipeDB(page);
  await populateDummyData(page);
});

test("logs in and reacts to an entry", async ({ page, baseURL }) => {
  // Try reacting to an entry before logging in to ensure What Got Done prompts
  // a login.
  await page.goto("/staging_jimmy/2019-06-28");
  await page.locator(".reaction-buttons .btn").first().click();
  await expect(page).toHaveURL("/login");

  // Log in
  await mockLoginAsUser(page, "reacting_tommy");
  await expect(page).toHaveURL(/\/entry\/edit\/.+/g);

  // Go back to the entry and react.
  await page.goto("/staging_jimmy/2019-06-28");
  await page.locator(".reaction-buttons .btn").first().click();

  await expect(page.locator(".reactions .reaction")).toHaveText(
    "reacting_tommy reacted with a 👍"
  );

  // Verify reaction persists after reload.
  await page.reload();
  await expect(page.locator(".reactions .reaction")).toHaveText(
    "reacting_tommy reacted with a 👍"
  );
});
