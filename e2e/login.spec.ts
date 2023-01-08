import { expect, test } from "@playwright/test";

test('clicking "Post Update" before authenticating prompts login', async ({
  page,
}) => {
  await page.goto("/");

  await expect(page.locator("nav .account-dropdown")).toHaveCount(0);

  await page.locator("nav .post-update").click();

  await expect(page).toHaveURL("/login");
});

test("back button should work if the user decides not to login/sign up", async ({
  page,
}) => {
  await page.goto("/");

  await page.locator("nav .post-update").click();

  await expect(page).toHaveURL("/login");

  await page.goBack();

  await expect(page).toHaveURL("/");
});

/*


it("logs in and signs out", () => {
  cy.visit("/");
  cy.get("nav .post-update").click();
  cy.completeLoginForm("staging_jimmy");

  cy.location("pathname").should("include", "/entry/edit");

  cy.get(".account-dropdown").click();
  cy.get(".sign-out-link a").click();
  cy.location("pathname").should("eq", "/");

  cy.get("nav .account-dropdown").should("not.exist");

  // Try signing in again.
  cy.get("nav .post-update").click();
  cy.completeLoginForm("staging_jimmy");

  cy.location("pathname").should("include", "/entry/edit");

  cy.get(".account-dropdown").click();
  cy.get(".sign-out-link a").click();
  cy.location("pathname").should("eq", "/");
});

it("bare route should redirect authenticated user to their edit entry page", () => {
  cy.intercept("/api/user/me").as("getUsername");

  cy.visit("/");
  cy.location("pathname").should("eq", "/");

  // Clicking the navbar brand should point to /about page.
  cy.get(".navbar .navbar-brand").click();
  cy.location("pathname").should("eq", "/");

  cy.login("staging_jimmy");
  cy.location("pathname").should("include", "/entry/edit");
  cy.wait("@getUsername");

  // Navigating back to bare route should redirect to edit entry page.
  cy.visit("/");
  cy.reload();
  cy.location("pathname").should("contain", "/entry/edit/");

  // Clicking navbar brand should point to edit entry page.
  cy.visit("/");
  cy.get(".navbar .navbar-brand").click();
  cy.location("pathname").should("contain", "/entry/edit/");

  // Log out
  cy.get(".account-dropdown").click();
  cy.get(".sign-out-link a").click();

  // Post-logout, user should be on bare route.
  cy.location("pathname").should("eq", "/");

  // Clicking navbar brand should point to bare route.
  cy.get(".navbar .navbar-brand").click();
  cy.location("pathname").should("eq", "/");
});

it("visiting authenticated page after UserKit token expires should redirect to login", () => {
  cy.visit("/");
  cy.get(".post-update").click();
  cy.completeLoginForm("joe123");

  cy.location("pathname").should("contain", "/entry/edit");
  cy.get(".account-dropdown").click();
  cy.get(".preferences-link a").click();

  cy.location("pathname").should("eq", "/preferences");

  // Simulate a UserKit cookie going stale.
  cy.setCookie("userkit_auth_token", "");

  cy.reload();

  cy.location("pathname").should("eq", "/login");
  cy.completeLoginForm("joe123");

  // Redirect to where the user was before the redirect.
  cy.location("pathname").should("eq", "/preferences");
});*/
