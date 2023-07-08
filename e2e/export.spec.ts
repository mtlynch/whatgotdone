import { expect, test } from "@playwright/test";
import { mockLoginAsUser } from "./helpers/test_apis";

test("can export the logged-in user's data", async ({ page }) => {
  await mockLoginAsUser(page, "staging_jimmy");

  await page.getByTestId("account-dropdown").click();
  await page.getByTestId("profile-link").click();

  await expect(page).toHaveURL("/staging_jimmy");

  await page.getByTestId("export-data-btn").click();
});

test("can't see an export data button for other users", async ({ page }) => {
  await mockLoginAsUser(page, "staging_jimmy");

  await page.goto("/leader_lenny");

  await expect(page.getByTestId("export-data-btn")).toHaveCount(0);
});
