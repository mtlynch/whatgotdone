import { expect, test } from "@playwright/test";
import { wipeDB, populateDummyData } from "./helpers/test_apis.js";
import { mockLoginAsUser } from "./helpers/test_apis";

test.beforeEach(async ({ page }) => {
  await wipeDB(page);
  await populateDummyData(page);
});

test("views an existing user's profile", async ({ page }) => {
  await page.goto("/staging_jimmy");

  await expect(page.locator("h1")).toHaveText("staging_jimmy");
});

test("gets 404 for non-existent user's profile page", async ({ page }) => {
  const response = await page.goto("/nonExistentUser");

  expect(response?.status()).toBe(404);
});

test("logs in and updates profile", async ({ page }) => {
  await mockLoginAsUser(page, "staging_jimmy");

  await page.locator("data-test-id=account-dropdown").click();
  await page.locator("data-test-id=profile-link").click();

  await expect(page).toHaveURL("/staging_jimmy");

  let apiUserGet = page.waitForRequest(
    (request) =>
      request.url().endsWith("/api/user/staging_jimmy") &&
      request.method() === "GET"
  );

  await page.locator("data-test-id=edit-btn").click();
  await expect(page).toHaveURL("/profile/edit");

  // Wait for page to pull down existing profile.
  await apiUserGet;

  await page.locator("#user-bio").fill("Hello, my name is staging_jimmy!");
  await page.locator("#email-address").fill("jimmy@example.com");
  await page.locator("#twitter-handle").fill("@jimmy");
  await page.locator("#mastodon-address").fill("jimmy@masto.example.com");
  await page.locator("#save-profile").click();

  await expect(page).toHaveURL("/staging_jimmy");
  await expect(page.locator("data-test-id=user-bio")).toHaveText(
    "Hello, my name is staging_jimmy!"
  );
  await expect(page.locator("data-test-id=email-address")).toHaveText(
    "jimmy@example.com"
  );
  await expect(page.locator("data-test-id=twitter-handle")).toHaveText(
    "@jimmy"
  );
  await expect(page.locator("data-test-id=mastodon-address")).toHaveText(
    "jimmy@masto.example.com"
  );
});

test("logs in and sets profile photo", async ({ page }) => {
  // TODO
});
