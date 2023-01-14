import { expect, test } from "@playwright/test";
import { wipeDB, populateDummyData } from "./helpers/db.js";
import { loginAsUser } from "./helpers/login.js";

test.beforeEach(async ({ page }) => {
  await wipeDB(page);
  await populateDummyData(page);
});

test("shows recent posts", async ({ page }) => {
  await page.goto("/");
  await page.locator("data-test-id=recent-link").click();
  await expect(page).toHaveURL("/recent");

  await expect(page.locator(".journal").first()).toHaveText(
    "staging_jimmy's update"
  );
});
