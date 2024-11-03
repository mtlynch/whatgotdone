import { expect, test } from "@playwright/test";
import { wipeDB, populateDummyData } from "./helpers/test_apis.js";

test.beforeEach(async ({ page }) => {
  await wipeDB(page);
  await populateDummyData(page);
});

test("shows recent posts", async ({ page }) => {
  await page.goto("/");
  await page.getByRole("link", { name: "Recent" }).click();
  await expect(page).toHaveURL("/recent");

  await expect(page.locator(".journal").first()).toContainText(
    "staging_jimmy's update for the week ending on Friday, Jun 28, 2019"
  );
  await expect(page.locator(".journal").first()).toContainText(
    "Today was a productive day. I created a test data manager for What Got Done!"
  );
});
