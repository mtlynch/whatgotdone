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

  await page.getByTestId("account-dropdown").click();
  await page.getByTestId("profile-link").click();

  await expect(page).toHaveURL("/staging_jimmy");

  let apiUserGet = page.waitForRequest(
    (request) =>
      request.url().endsWith("/api/user/staging_jimmy") &&
      request.method() === "GET"
  );

  await page.getByRole("link", { name: "Edit" }).click();
  await expect(page).toHaveURL("/profile/edit");

  // Wait for page to pull down existing profile.
  await apiUserGet;

  await page.locator("#user-bio").fill("Hello, my name is staging_jimmy!");
  await page.locator("#email-address").fill("jimmy@example.com");
  await page.locator("#twitter-handle").fill("@jimmy");
  await page.locator("#mastodon-address").fill("jimmy@masto.example.com");
  await page.locator("#save-profile").click();

  await expect(page).toHaveURL("/staging_jimmy");
  await expect(page.getByTestId("user-bio")).toHaveText(
    "Hello, my name is staging_jimmy!"
  );
  await expect(page.getByTestId("email-address")).toHaveText(
    "jimmy@example.com"
  );
  await expect(page.getByTestId("twitter-handle")).toHaveText("@jimmy");
  await expect(page.getByTestId("mastodon-address")).toHaveText(
    "jimmy@masto.example.com"
  );
});

test("logs in and sets profile photo", async ({ page }) => {
  await mockLoginAsUser(page, "staging_jimmy");

  await page.getByTestId("account-dropdown").click();
  await page.getByTestId("profile-link").click();

  await expect(page).toHaveURL("/staging_jimmy");

  await expect(page.locator(".profile img")).toHaveScreenshot(
    "unknown-person.png"
  );

  const apiUserGet = page.waitForRequest(
    (request) =>
      request.url().endsWith("/api/user/staging_jimmy") &&
      request.method() === "GET"
  );

  await page.getByRole("link", { name: "Edit" }).click();
  await expect(page).toHaveURL("/profile/edit");

  // Wait for page to pull down existing profile.
  await apiUserGet;

  let apiUserAvatarResponse = page.waitForResponse("**/api/user/avatar");
  await page
    .locator("#upload-profile-photo")
    .setInputFiles(["e2e/testdata/kittyface.jpg"]);
  await apiUserAvatarResponse;

  await expect(page.locator(".avatar-wrapper img")).toHaveScreenshot(
    "kittyface-80px-circle.png"
  );

  await page.locator("#save-profile").click();

  await expect(page).toHaveURL("/staging_jimmy");

  // Workaround for timing.
  await page.reload();
  await expect(page.locator(".profile img")).toHaveScreenshot(
    "kittyface-150px-circle.png"
  );

  // Delete the avatar image.
  await page.getByRole("link", { name: "Edit" }).click();
  await expect(page).toHaveURL("/profile/edit");

  apiUserAvatarResponse = page.waitForResponse("**/api/user/avatar");
  await page.locator("#delete-profile-photo").click();
  await apiUserAvatarResponse;

  await page.locator("#save-profile").click();

  await expect(page).toHaveURL("/staging_jimmy");
  // Workaround for timing.
  await page.reload();
  await expect(page.locator(".profile img")).toHaveScreenshot(
    "unknown-person.png"
  );
});
