it("writes an entry without formatting", () => {
  cy.login("staging_jimmy");

  cy.location("pathname").should("include", "/entry/edit");

  const entryText = "Posted an update at " + new Date().toISOString();

  cy.get(".editor-content .ProseMirror").type(entryText);
  cy.get("form").submit();

  cy.location("pathname").should("include", "/staging_jimmy/");

  cy.get(".journal-body").should("contain", entryText);
});

it("writes an entry with every type of formatting", () => {
  cy.login("staging_jimmy");

  cy.location("pathname").should("include", "/entry/edit");

  cy.get(".btn-h1 .btn").click();
  cy.get(".editor-content .ProseMirror").type("Project A{enter}{enter}");

  cy.get(".btn-h2 .btn").click();
  cy.get(".editor-content .ProseMirror").type("Subproject B{enter}{enter}");

  cy.get(".btn-h3 .btn").click();
  cy.get(".editor-content .ProseMirror").type("Topic 1{enter}{enter}");

  cy.get(".editor-content .ProseMirror").type("This week was ");
  cy.get(".btn-bold .btn").click();
  cy.get(".btn-bold .btn").should("have.class", "is-active");
  cy.get(".editor-content .ProseMirror").type("very difficult");
  cy.get(".btn-bold .btn").click();
  cy.get(".btn-bold .btn").should("not.have.class", "is-active");
  cy.get(".editor-content .ProseMirror").type("!");
  cy.get(".editor-content .ProseMirror").type("{enter}{enter}");

  cy.get(".editor-content .ProseMirror").type("I ");
  cy.get(".btn-italic .btn").click();
  cy.get(".editor-content .ProseMirror").type("discovered ");
  cy.get(".btn-italic .btn").click();
  cy.get(".btn-strikethrough .btn").click();
  cy.get(".editor-content .ProseMirror").type("11");
  cy.get(".btn-strikethrough .btn").click();
  cy.get(".editor-content .ProseMirror").type(" 22 new bugs.{enter}{enter}");

  // TODO: use link

  cy.get(".editor-content .ProseMirror").type("Most were in the ");
  cy.get(".btn-inline-code .btn").click();
  cy.get(".editor-content .ProseMirror").type("Frombobulator");
  cy.get(".btn-inline-code .btn").click();
  cy.get(".editor-content .ProseMirror").type(
    "component. The typical bad code looks like this:{enter}"
  );

  cy.get(".btn-code-block .btn").click();
  cy.get(".editor-content .ProseMirror").type("f = new Frombobulator(){enter}");
  cy.get(".editor-content .ProseMirror").type("f.frombobulate(){ctrl}{enter}");
  cy.get(".editor-content .ProseMirror").type("Yuck!");
  // TODO: use bulleted list
  // TODO: use ordered list
  // TODO: use blockquote

  cy.get(".switch-mode .btn").click();

  cy.get(".markdown-editor textarea").should(
    "have.value",
    `# Project A
## Subproject B
### Topic 1
This week was **very difficult**!

I _discovered_ ~11~ 22 new bugs.

Most were in the \`Frombobulator\` component. The typical bad code looks like this:
    f = new Frombobulator()
    f.frombobulate()
Yuck!`
  );

  cy.get("form").submit();

  cy.location("pathname").should("include", "/staging_jimmy/");

  // TODO: Check rendered text
});

it("does not inject HTML comments", () => {
  cy.login("staging_jimmy");

  cy.location("pathname").should("include", "/entry/edit");

  cy.get(".btn-bulleted-list .btn").click();
  cy.get(".editor-content .ProseMirror").type("a{enter}{enter}b");

  cy.get(".switch-mode .btn").click();

  cy.get(".editor-textarea").should("have.value", "*    a\n\nb");
});
