import { expect, test } from "@playwright/test";

test("loads the homepage", async ({ page }) => {
  await page.goto("/");

  // Verify <head> metadata.
  expect(page).toHaveTitle("What Got Done");
  await expect(page.locator("meta[name='description']")).toHaveAttribute(
    "content",
    "The simple, easy way to share progress with your teammates."
  );
  await expect(page.locator("meta[property='og:type']")).toHaveAttribute(
    "content",
    "website"
  );
  await expect(page.locator("meta[property='og:title']")).toHaveAttribute(
    "content",
    "What Got Done"
  );
  await expect(page.locator("meta[property='og:description']")).toHaveAttribute(
    "content",
    "The simple, easy way to share progress with your teammates."
  );
  await expect(page.locator("h1")).toHaveText(
    "What did you get done this week?"
  );
});
