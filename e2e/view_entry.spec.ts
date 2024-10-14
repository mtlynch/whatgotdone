import { expect, test } from "@playwright/test";
import { wipeDB, populateDummyData } from "./helpers/test_apis.js";

test.beforeEach(async ({ page }) => {
  await wipeDB(page);
  await populateDummyData(page);
});

test("renders the date correctly", async ({ page }) => {
  await page.goto("/staging_jimmy/2019-06-28");
  await expect(page).toHaveTitle(
    "staging_jimmy's What Got Done for the week of 2019-06-28"
  );

  await expect(page.locator(".journal").first()).toContainText(
    "staging_jimmy's update for the week ending on Friday, Jun 28, 2019"
  );
});
