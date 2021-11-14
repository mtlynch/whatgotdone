it("logs in and saves a draft", () => {
  cy.server();
  cy.route("GET", "/api/draft/*").as("getDraft");
  cy.route("PUT", "/api/draft/*").as("putDraft");

  cy.login("staging_jimmy");

  cy.location("pathname").should("include", "/entry/edit");

  // Wait for page to pull down any previous entry.
  cy.wait("@getDraft");

  const entryText = "Saved a private draft at " + new Date().toISOString();

  cy.get(".editor-content .ProseMirror").clear().type(entryText);

  // Wait for auto-save to complete.
  cy.get(".save-draft").should("contain", "Changes Saved");
  cy.wait("@putDraft");

  // User should stay on the same page after saving a draft.
  cy.location("pathname").should("include", "/entry/edit");

  cy.visit("/recent");

  // Private drafts should not appear on the recent page
  cy.get("#app").should("not.contain", entryText);
});

it("don't overwrite draft until we successfully sync the latest draft from the server", () => {
  cy.server();
  cy.route({
    method: "GET",
    url: "/api/draft/*",
    response: {},
    status: 500,
  }).as("getDraft");
  cy.route("PUT", "/api/draft/*").as("putDraft");

  cy.login("staging_jimmy");

  cy.location("pathname").should("include", "/entry/edit");

  // Wait for page to fail on its request to pull down the previous draft.
  cy.wait("@getDraft");

  cy.get(".journal-markdown").should("not.exist");
  cy.get(".save-draft").should("not.exist");

  cy.routeShouldBeCalled("putDraft", 0);
  cy.get(".entry-form").should("not.exist");
});

it("uses the entry template for new drafts", () => {
  cy.server();
  cy.route("GET", "/api/draft/*").as("getDraft");
  cy.route("PUT", "/api/draft/*").as("putDraft");
  cy.route("POST", "/api/prefrences").as("postPreferences");

  cy.login("staging_jimmy");

  cy.location("pathname").should("include", "/entry/edit");
  cy.visit("/entry/edit/2020-03-06");

  // Wait for page to pull down any previous entry.
  cy.wait("@getDraft");

  cy.get(".editor-content .ProseMirror").should("have.value", "");

  // Set an entry template on the preferences page.
  cy.visit("/preferences");

  cy.get("#entry-template-input").type(
    "# Example project\n\n* Item A\n* Item B"
  );
  cy.get(".btn-primary").click();

  cy.get(".alert-success").should("be.visible");

  // Verify new entries start with the template.
  cy.visit("/entry/edit/2020-03-06");
  cy.location("pathname").should("include", "/entry/edit");

  cy.wait("@getDraft");
  cy.get(".switch-mode .btn").click();
  cy.get(".markdown-editor .editor-textarea").should(
    "have.value",
    "# Example project\n\n* Item A\n* Item B"
  );
});
