it("writes an entry without formatting", () => {
  cy.login("staging_jimmy");

  cy.location("pathname").should("include", "/entry/edit");

  const entryText = "Posted an update at " + new Date().toISOString();

  cy.get(".editor-content .ProseMirror").type(entryText);
  cy.get("form").submit();

  cy.location("pathname").should("include", "/staging_jimmy/");

  cy.get(".journal-body").should("contain", entryText);
});

it("canonicalizes newlines after bulleted list", () => {
  cy.login("staging_jimmy");

  cy.location("pathname").should("include", "/entry/edit");

  cy.get(".btn-bulleted-list .btn").click();
  cy.get(".editor-content .ProseMirror").type("a{enter}{enter}b");

  cy.get(".switch-mode .btn").click();

  cy.get(".editor-textarea").should("have.value", "*   a\n\nb");
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

  cy.get(".editor-content .ProseMirror").type("Full report is on the wiki");
  cy.get(".editor-content .ProseMirror").setSelection("wiki");
  cy.get(".btn-link .btn").click();
  cy.get(".modal-dialog input").type("https://example.com/wiki{enter}");
  cy.get(".editor-content .ProseMirror").type("{rightarrow}.{enter}{enter}");

  cy.get(".editor-content .ProseMirror").type("Most were in the ");
  cy.get(".btn-inline-code .btn").click();
  cy.get(".editor-content .ProseMirror").type("Frombobulator");
  cy.get(".btn-inline-code .btn").click();
  cy.get(".editor-content .ProseMirror").type(
    " component. The typical bad code looks like this:{enter}"
  );

  cy.get(".btn-code-block .btn").click();
  cy.get(".editor-content .ProseMirror").type("f = new Frombobulator(){enter}");
  cy.get(".editor-content .ProseMirror").type("f.frombobulate(){ctrl}{enter}");
  cy.get(".editor-content .ProseMirror").type("Yuck!{enter}");

  cy.get(".editor-content .ProseMirror").type(
    "These are the things I picked up from the supermarket:{enter}"
  );
  cy.get(".btn-bulleted-list .btn").click();
  cy.get(".editor-content .ProseMirror").type("eggs{enter}");
  cy.get(".editor-content .ProseMirror").type("egg cartons{enter}");
  cy.get(".editor-content .ProseMirror").type(
    "egg carton holders{enter}{enter}"
  );

  cy.get(".editor-content .ProseMirror").type(
    "And I prepared these meals (in order):{enter}"
  );
  cy.get(".btn-numbered-list .btn").click();
  cy.get(".editor-content .ProseMirror").type("scrambled eggs{enter}");
  cy.get(".editor-content .ProseMirror").type("omelettes{enter}");
  cy.get(".editor-content .ProseMirror").type("fritattas{enter}{enter}");

  cy.get(".editor-content .ProseMirror").type(
    "When my manager sees me eat lunch every day, he says:{enter}"
  );
  cy.get(".btn-blockquote .btn").click();
  cy.get(".editor-content .ProseMirror").type(
    "hey jimmy why u eat so much eggs?{enter}{enter}"
  );
  cy.get(".editor-content .ProseMirror").type(
    "And I just tell him I love the protein!"
  );

  cy.get(".switch-mode .btn").click();

  cy.get(".markdown-editor textarea").should(
    "have.value",
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

  cy.get("form").submit();

  cy.location("pathname").should("include", "/staging_jimmy/");
});

it("supports undo", () => {
  cy.login("staging_jimmy");

  cy.location("pathname").should("include", "/entry/edit");

  cy.get(".editor-content .ProseMirror").type("Hello, world!");

  // Wait for auto-save to complete
  cy.get(".save-draft").should("contain", "Changes Saved");
  cy.reload();

  // Delete content
  cy.get(".editor-content .ProseMirror").type("{selectall}{del}");

  // Undo
  cy.get(".editor-content .ProseMirror").type("{ctrl}z");
  cy.get(".switch-mode .btn").click();
  cy.get(".editor-textarea").should("have.value", "Hello, world!");
});

it("canonicalizes newlines after bulleted list", () => {
  cy.login("staging_jimmy");

  cy.location("pathname").should("include", "/entry/edit");

  cy.get(".btn-bulleted-list .btn").click();
  cy.get(".editor-content .ProseMirror").type("a{enter}{enter}b");

  cy.get(".switch-mode .btn").click();

  cy.get(".editor-textarea").should("have.value", "*   a\n\nb");
});
