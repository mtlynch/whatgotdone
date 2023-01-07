import { test } from "@playwright/test";

test("gets the sitemap", async ({ page }) => {
  await page.goto("/sitemap.xml");
});

test("gets the robots.txt file", async ({ page }) => {
  await page.goto("/robots.txt");
});
