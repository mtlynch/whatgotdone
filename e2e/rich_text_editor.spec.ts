import { expect, test } from "@playwright/test";
import { mockLoginAsUser, wipeDB } from "./helpers/test_apis.js";

test.beforeEach(async ({ page }) => {
  await wipeDB(page);
});

test("writes an entry without formatting", async ({ page }) => {
  await mockLoginAsUser(page, "staging_jimmy");

  await expect(page).toHaveURL(/\/entry\/edit\/.+/g);

  await page.locator(".editor-content .ProseMirror").clear();
  await page
    .locator(".editor-content .ProseMirror")
    .fill("Posted a new update without formatting!");

  await page.locator("form button[type='submit']").click();

  await expect(page).toHaveURL(/\/staging_jimmy\/.+/g);

  await expect(page.locator(".journal-body")).toHaveText(
    "Posted a new update without formatting!"
  );
});

test("canonicalizes newlines after bulleted list", async ({ page }) => {
  await mockLoginAsUser(page, "staging_jimmy");

  await expect(page).toHaveURL(/\/entry\/edit\/.+/g);

  await page.locator(".editor-content .ProseMirror").clear();
  await page.locator(".btn-bulleted-list .btn").click();
  await page.keyboard.press("a");
  await page.keyboard.press("Enter");
  await page.keyboard.press("Enter");
  await page.keyboard.press("b");

  await page.locator(".switch-mode .btn").click();
  await expect(page.locator(".editor-textarea")).toHaveValue("*   a\n\nb");
});

test("writes an entry with every type of formatting", async ({ page }) => {
  await mockLoginAsUser(page, "staging_jimmy");

  await expect(page).toHaveURL(/\/entry\/edit\/.+/g);

  await page.locator(".btn-h1 .btn").click();
  await page.keyboard.type("Project A");
  await page.keyboard.press("Enter");

  await page.locator(".btn-h2 .btn").click();
  await page.keyboard.type("Subproject B");
  await page.keyboard.press("Enter");

  await page.locator(".btn-h3 .btn").click();
  await page.keyboard.type("Topic 1");
  await page.keyboard.press("Enter");

  await page.keyboard.type("This week was ");
  await page.locator(".btn-bold .btn").click();
  await expect(page.locator(".btn-bold .btn")).toHaveClass(/is-active/);
  await page.keyboard.type("very difficult");
  await page.locator(".btn-bold .btn").click();
  await expect(page.locator(".btn-bold .btn")).not.toHaveClass(/is-active/);
  await page.keyboard.type("!");
  await page.keyboard.press("Enter");
  await page.keyboard.press("Enter");

  await page.keyboard.type("I ");
  await page.locator(".btn-italic .btn").click();
  await page.keyboard.type("discovered ");
  await page.locator(".btn-italic .btn").click();
  await page.locator(".btn-strikethrough .btn").click();
  await page.keyboard.type("11");
  await page.locator(".btn-strikethrough .btn").click();
  await page.keyboard.type(" 22 new bugs.");
  await page.keyboard.press("Enter");
  await page.keyboard.press("Enter");

  await page
    .locator(".editor-content .ProseMirror")
    .type("Full report is on the wiki.");
  await page.keyboard.press("ArrowLeft");
  await page.keyboard.down("Shift");
  for (let i = 0; i < "wiki".length; i++) {
    await page.keyboard.press("ArrowLeft");
  }
  await page.keyboard.up("Shift");
  await page.locator(".btn-link .btn").click();
  await page.locator(".modal-dialog input").click();
  await page.keyboard.type("example.com/wiki");
  await page.keyboard.press("Enter");
  await page.keyboard.press("End");
  await page.keyboard.press("Enter");
  await page.keyboard.press("Enter");

  await page.keyboard.type("Most were in the ");
  await page.locator(".btn-inline-code .btn").click();
  await page.keyboard.type("Frombobulator");
  await page.locator(".btn-inline-code .btn").click();
  await page.keyboard.type(" component. The typical bad code looks like this:");
  await page.keyboard.press("Enter");

  await page.locator(".btn-code-block .btn").click();
  await page.keyboard.type("f = new Frombobulator()");
  await page.keyboard.press("Enter");
  await page.keyboard.type("f.frombobulate()");
  await page.keyboard.press("Control+Enter");

  await page.keyboard.type("Yuck!");
  await page.keyboard.press("Enter");

  await page
    .locator(".editor-content .ProseMirror")
    .type("These are the things I picked up from the supermarket:");
  await page.keyboard.press("Enter");
  await page.locator(".btn-bulleted-list .btn").click();
  await page.keyboard.type("eggs");
  await page.keyboard.press("Enter");
  await page.keyboard.type("egg cartons");
  await page.keyboard.press("Enter");
  await page.keyboard.type("egg carton holders");
  await page.keyboard.press("Enter");
  await page.keyboard.press("Enter");

  await page
    .locator(".editor-content .ProseMirror")
    .type("And I prepared these meals (in order):");
  await page.keyboard.press("Enter");
  await page.locator(".btn-numbered-list .btn").click();
  await page.keyboard.type("scrambled eggs");
  await page.keyboard.press("Enter");
  await page.keyboard.type("omelettes");
  await page.keyboard.press("Enter");
  await page.keyboard.type("fritattas");
  await page.keyboard.press("Enter");
  await page.keyboard.press("Enter");

  await page
    .locator(".editor-content .ProseMirror")
    .type("When my manager sees me eat lunch every day, he says:");
  await page.keyboard.press("Enter");
  await page.locator(".btn-blockquote .btn").click();
  await page
    .locator(".editor-content .ProseMirror")
    .type("hey jimmy why u eat so much eggs?");
  await page.keyboard.press("Enter");
  await page.keyboard.press("Enter");
  await page
    .locator(".editor-content .ProseMirror")
    .type("And I just tell him I love the protein!");

  await page.locator(".switch-mode .btn").click();

  await expect(page.locator(".markdown-editor textarea")).toHaveValue(
    `# Project A

## Subproject B

### Topic 1

This week was **very difficult**!

I _discovered_ ~11~ 22 new bugs.

Full report is on the [wiki](https://example.com/wiki).

Most were in the \`Frombobulator\` component. The typical bad code looks like this:

    f = new Frombobulator()
    f.frombobulate()

Yuck!

These are the things I picked up from the supermarket:

*   eggs

*   egg cartons

*   egg carton holders

And I prepared these meals (in order):

1.  scrambled eggs

2.  omelettes

3.  fritattas

When my manager sees me eat lunch every day, he says:

> hey jimmy why u eat so much eggs?

And I just tell him I love the protein!`
  );

  await page.locator("form button[type='submit']").click();

  await expect(page).toHaveURL(/\/staging_jimmy\/.+/g);
});

test("supports undo", async ({ page }) => {
  await mockLoginAsUser(page, "staging_jimmy");

  await expect(page).toHaveURL(/\/entry\/edit\/.+/g);

  await page.locator(".editor-content .ProseMirror").click();

  await page.keyboard.type("Hello, world!");

  // Wait for auto-save to complete
  await expect(page.locator(".save-draft")).toHaveText("Changes Saved");
  page.reload();

  // Delete content
  await page.keyboard.press("Control+A");
  await page.keyboard.press("Delete");

  // Undo
  await page.keyboard.press("Control+z");
  await page.locator(".switch-mode .btn").click();
  await expect(page.locator(".editor-textarea")).toHaveValue("Hello, world!");
});
