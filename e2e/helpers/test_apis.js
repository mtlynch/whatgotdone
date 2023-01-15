export async function wipeDB(page) {
  await page.goto("/api/testing/db/wipe");
}

export async function populateDummyData(page) {
  await page.goto("/api/testing/db/populate-dummy-data");
}

export async function mockLoginAsUser(page, username) {
  let apiUserMetGet = page.waitForRequest(
    (request) =>
      request.url().indexOf("/api/user/me") >= 0 && request.method() === "GET"
  );

  await page.goto(`/api/testing/auth/login/${username}`);

  // Reload the page now that it's authenticated.
  await apiUserMetGet;
  await page.reload();

  // Click the brand so that the frontend places us in the page where we'd be
  // after a real login.
  await page.locator(".navbar .navbar-brand").click();
}
