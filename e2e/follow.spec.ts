import { expect, test } from "@playwright/test";
import {
  mockLoginAsUser,
  wipeDB,
  populateDummyData,
} from "./helpers/test_apis.js";

test.beforeEach(async ({ page }) => {
  await wipeDB(page);
  await populateDummyData(page);
});

test("follows a user", async ({ page }) => {
  // Log in as follower user.
  await mockLoginAsUser(page, "follower_frank");
  await expect(page).toHaveURL(/\/entry\/edit\/.+/g);

  // Verify the personalized feed is empty.
  await page.getByRole("link", { name: "Feed" }).click();
  await expect(page.locator(".alert")).toContainText(
    "You're not following anyone yet."
  );
  await expect(page.locator(".journal")).not.toBeVisible();

  // Follow leader_lenny
  let apiFollowingGet = page.waitForRequest(
    (request) =>
      request.url().endsWith("/api/user/follower_frank/following") &&
      request.method() === "GET"
  );

  await page.goto("/leader_lenny");
  await apiFollowingGet;
  await page.getByTestId("follow-btn").click();
  await expect(page.getByTestId("unfollow-btn")).toBeVisible();
  await expect(page.getByTestId("follow-btn")).not.toBeVisible();

  // Verify personalized feed is non-empty
  await page.getByRole("link", { name: "Feed" }).click();
  await expect(page).toHaveURL("/feed");
  await expect(page.locator(".journal")).toHaveCount(2);

  // Unfollow leader_lenny
  await page.goBack();
  await expect(page).toHaveURL("/leader_lenny");
  await page.getByTestId("unfollow-btn").click();

  // Verify the personalized feed is empty again.
  await page.getByRole("link", { name: "Feed" }).click();
  await expect(page).toHaveURL("/feed");
  await expect(page.locator(".alert")).toContainText(
    "You're not following anyone yet."
  );
  await expect(page.locator(".journal")).not.toBeVisible();
});
