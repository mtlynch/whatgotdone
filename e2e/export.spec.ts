import { expect, test } from "@playwright/test";
import { loginAsUser } from "./helpers/login.js";

test("can export the logged-in user's data", async ({ page }) => {
  await loginAsUser(page, "staging_jimmy");

  await page.locator("data-test-id=account-dropdown").click();
  await page.locator("data-test-id=profile-link").click();

  await expect(page).toHaveURL("/staging_jimmy");

  await page.locator("data-test-id=export-data-btn").click();
});

test("can't see an export data button for other users", async ({ page }) => {
  await loginAsUser(page, "staging_jimmy");

  await page.goto("/leader_lenny");

  await expect(page.locator("data-test-id=export-data-btn")).toHaveCount(0);
});

/*it("can export the logged-in user's data", () => {
  cy.login("staging_jimmy");

  cy.get(".account-dropdown").click();
  cy.get(".profile-link a").click();

  cy.location("pathname").should("eq", "/staging_jimmy");

  cy.get('[data-test-id="export-data-btn"]').click();
});

it("can't see an export data button for other users", () => {
  cy.login("staging_jimmy");

  cy.visit("/leader_lenny");

  cy.get('[data-test-id="export-data-btn"]').should("not.exist");
});*/
