import { expect, test } from "@playwright/test";

async function loginAsUser(page, username) {
  await page.goto("/login");

  await page.locator("id=userkit_username").fill(username);
  await page.locator("id=userkit_password").fill("password");
  await page.locator("form button[type='submit']").click();
}

test('clicking "Post Update" before authenticating prompts login', async ({
  page,
}) => {
  await page.goto("/");

  await expect(page.getByTestId("account-dropdown")).toHaveCount(0);

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

test("logs in and signs out", async ({ page }) => {
  await loginAsUser(page, "staging_jimmy");

  await expect(page).toHaveURL(/\/entry\/edit\/.+/g);

  await page.getByTestId("account-dropdown").click();
  await page.getByTestId("sign-out-link").click();

  await expect(page).toHaveURL("/");

  await expect(page.getByTestId("account-dropdown")).toHaveCount(0);

  // Try signing in again.
  await loginAsUser(page, "staging_jimmy");

  await page.waitForTimeout(5 * 1000);

  await expect(page).toHaveURL(/\/entry\/edit\/.+/g, { timeout: 0 * 1000 });

  await page.getByTestId("account-dropdown").click();
  await page.getByTestId("sign-out-link").click();

  await expect(page).toHaveURL("/");
});

test("bare route should redirect authenticated user to their edit entry page", async ({
  page,
}) => {
  await page.goto("/");

  // Clicking the navbar brand should point to homepage.
  await page.locator(".navbar .navbar-brand").click();
  await expect(page).toHaveURL("/");

  await loginAsUser(page, "staging_jimmy");
  await expect(page).toHaveURL(/\/entry\/edit\/.+/g);

  // Navigating back to bare route should redirect to edit entry page.
  await page.goto("/");
  await expect(page).toHaveURL(/\/entry\/edit\/.+/g);

  // Clicking navbar brand should point to edit entry page.
  await page.locator(".navbar .navbar-brand").click();
  await expect(page).toHaveURL(/\/entry\/edit\/.+/g);

  // Log out
  await page.getByTestId("account-dropdown").click();
  await page.getByTestId("sign-out-link").click();
  await expect(page).toHaveURL("/");

  // Clicking the navbar brand should point to homepage.
  await page.locator(".navbar .navbar-brand").click();
  await expect(page).toHaveURL("/");
});

test("visiting authenticated page after UserKit token expires should redirect to login", async ({
  browser,
}) => {
  const browserContext = await browser.newContext();
  const page = await browserContext.newPage();
  await loginAsUser(page, "joe123");
  await expect(page).toHaveURL(/\/entry\/edit\/.+/g);

  await page.getByTestId("account-dropdown").click();
  await page.getByTestId("preferences-link").click();
  await expect(page).toHaveURL("/preferences");

  // Simulate a UserKit cookie going stale.
  browserContext.addCookies([
    {
      name: "userkit_auth_token",
      value: "some-invalid-value",
      domain: "localhost",
      path: "/",
    },
  ]);

  await page.reload();
  await expect(page).toHaveURL("/login");
  await loginAsUser(page, "joe123");

  await expect(page).toHaveURL(/\/entry\/edit\/.+/g);
});
