import { expect, test } from "@playwright/test";

test("gets the sitemap", async ({ page }) => {
  const response = await page.goto("/sitemap.xml");

  expect(response?.ok()).toBe(true);
});

test("gets the robots.txt file", async ({ page }) => {
  const response = await page.goto("/robots.txt");

  expect(response?.ok()).toBe(true);
});
