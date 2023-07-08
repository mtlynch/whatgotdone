import { expect, test } from "@playwright/test";
import {
  mockLoginAsUser,
  populateDummyData,
  wipeDB,
} from "./helpers/test_apis.js";

test.beforeEach(async ({ page }) => {
  await wipeDB(page);
});

test("logs in and posts an update", async ({ page }) => {
  const apiDraftGet = page.waitForRequest(
    (request) =>
      request.url().includes("/api/draft") && request.method() === "GET"
  );
  await mockLoginAsUser(page, "staging_jimmy");

  await expect(page).toHaveURL(/\/entry\/edit\/.+/g);

  // Wait for page to pull down any previous entry.
  await apiDraftGet;

  const entryText = "Posted an update at " + new Date().toISOString();

  await page.locator(".editor-content .ProseMirror").clear();
  await page.locator(".editor-content .ProseMirror").fill(entryText);

  await page.locator("form button[type='submit']").click();

  await expect(page).toHaveURL(/\/staging_jimmy\/.+/g);

  // Reload the page to fetch the new HTML rather than using what the front-end
  // generated client-side.
  await page.reload();

  await expect(page).toHaveTitle(
    /^staging_jimmy's What Got Done for the week of/
  );
  await expect(page.locator("meta[property='og:type']")).toHaveAttribute(
    "content",
    "article"
  );

  await expect(page.locator(".journal-body")).toHaveText(entryText);
  await expect(page.locator(".missing-entry")).toHaveCount(0);
});

test("logs in and backdates an update from a previous week", async ({
  page,
}) => {
  const apiDraftGet = page.waitForRequest(
    (request) =>
      request.url().includes("/api/draft") && request.method() === "GET"
  );
  await mockLoginAsUser(page, "staging_jimmy");

  // Wait for page to pull down any previous entry.
  await apiDraftGet;

  await page.goto("/entry/edit/2019-12-13");

  const entryText = "Posted a backdated update at " + new Date().toISOString();

  await page.locator(".editor-content .ProseMirror").clear();
  await page.locator(".editor-content .ProseMirror").fill(entryText);

  await page.locator("form button[type='submit']").click();

  await expect(page).toHaveURL("/staging_jimmy/2019-12-13");

  // Reload the page to fetch the new HTML rather than using what the front-end
  // generated client-side.
  await page.reload();

  await expect(page).toHaveTitle(
    "staging_jimmy's What Got Done for the week of 2019-12-13"
  );
  await expect(page.locator("meta[name='description']")).toHaveAttribute(
    "content",
    "Find out what staging_jimmy accomplished for the week ending on December 13, 2019"
  );
  await expect(page.locator("meta[property='og:type']")).toHaveAttribute(
    "content",
    "article"
  );
  await expect(page.locator("meta[property='og:title']")).toHaveAttribute(
    "content",
    "staging_jimmy's What Got Done for the week of Dec. 13, 2019"
  );
  await expect(page.locator("meta[property='og:description']")).toHaveAttribute(
    "content",
    "Find out what staging_jimmy accomplished for the week ending on December 13, 2019"
  );

  await expect(page.locator(".journal-body")).toHaveText(entryText);
  await expect(page.locator(".missing-entry")).toHaveCount(0);
});

test("posts an update and then unpublishes it", async ({ page }) => {
  // Populate dummy data or else we won't get the right "missing entry" message
  // for staging_jimmy.
  await populateDummyData(page);

  const apiDraftGet = page.waitForRequest(
    (request) =>
      request.url().includes("/api/draft") && request.method() === "GET"
  );
  await mockLoginAsUser(page, "staging_jimmy");

  // Wait for page to pull down any previous entry.
  await apiDraftGet;

  await page.goto("/entry/edit/2019-06-28");

  await page.locator(".editor-content .ProseMirror").clear();
  await page
    .locator(".editor-content .ProseMirror")
    .fill("felt cute, might unpublish later");

  await page.locator("form button[type='submit']").click();

  await expect(page).toHaveURL("/staging_jimmy/2019-06-28");

  await expect(page.locator(".missing-entry")).toHaveCount(0);

  await page.getByRole("button", { name: "Unpublish" }).click();

  // Unpublishing takes the user back to the edit entry page.
  await expect(page).toHaveURL("/entry/edit/2019-06-28");

  // Go back and make sure the published entry is gone.
  await page.goBack();
  await expect(page).toHaveURL("/staging_jimmy/2019-06-28");
  await expect(page.locator(".missing-entry")).toBeVisible();
});
