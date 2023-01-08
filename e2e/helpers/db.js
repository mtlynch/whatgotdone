export async function wipeDB(page) {
  await page.goto("/api/testing/db/wipe");
}

export async function populateDummyData(page) {
  await page.goto("/api/testing/db/populate-dummy-data");
}
