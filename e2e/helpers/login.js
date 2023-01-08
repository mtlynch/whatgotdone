export async function loginAsUser(page, username) {
  await page.goto("/login");

  await page.locator("id=userkit_username").fill(username);
  await page.locator("id=userkit_password").fill("password");
  await page.locator("form button[type='submit']").click();
}
