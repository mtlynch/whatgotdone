import { expect, test } from "@playwright/test";
import { mockLoginAsUser } from "./helpers/test_apis";

test("can export the logged-in user's data", async ({ page }) => {
  await mockLoginAsUser(page, "staging_jimmy");

  await page.getByRole("button", { name: "Account" }).click();
  await page.getByRole("menuitem", { name: "Export" }).click();

  await expect(page).toHaveURL("/export");

  await page.getByRole("button", { name: "Download (JSON)" }).click();
});

test("can't see an export data button for other users", async ({ page }) => {
  await mockLoginAsUser(page, "staging_jimmy");

  await page.goto("/leader_lenny");

  await expect(
    page.getByRole("button", { name: "Download (JSON)" })
  ).not.toBeVisible();
});
