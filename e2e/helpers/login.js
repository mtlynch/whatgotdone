import { expect } from "@playwright/test";

async function loginAsUser(page, username, password) {
  await page.goto("/");

  await page.locator("data-test-id=sign-in-btn").click();

  await expect(page).toHaveURL("/login");
  await page.locator("id=username").fill(username);
  await page.locator("id=password").fill(password);
  await page.locator("form input[type='submit']").click();

  await expect(page).toHaveURL("/reviews");
}

export async function loginAsAdmin(page) {
  await loginAsUser(page, "dummyadmin", "dummypass");
}

export async function loginAsUserA(page) {
  await loginAsUser(page, "userA", "password123");
}

export async function loginAsUserB(page) {
  await loginAsUser(page, "userB", "password456");
}
