it("logs in and saves a draft", () => {
  cy.server();
  cy.route("GET", "/api/draft/*").as("getDraft");
  cy.route("POST", "/api/draft/*").as("postDraft");

  cy.login("staging_jimmy");

  cy.location("pathname").should("include", "/entry/edit");

  // Wait for page to pull down any previous entry.
  cy.wait("@getDraft");

  const entryText = "Saved a private draft at " + new Date().toISOString();

  cy.get(".journal-markdown")
    .clear()
    .type(entryText);
  cy.get(".save-draft").click();

  // Wait for "save draft" operation to complete.
  cy.wait("@postDraft");

  // User should stay on the same page after saving a draft.
  cy.location("pathname").should("include", "/entry/edit");

  cy.visit("/recent");

  // Private drafts should not appear on the recent page
  cy.get("#app").should("not.contain", entryText);
});

it("uses the entry template for new drafts", () => {
  cy.server();
  cy.route("GET", "/api/draft/*").as("getDraft");
  cy.route("POST", "/api/draft/*").as("postDraft");
  cy.route("POST", "/api/prefrences").as("postPreferences");

  cy.login("staging_jimmy");

  cy.location("pathname").should("include", "/entry/edit");
  cy.visit("/entry/edit/2020-03-06");

  // Wait for page to pull down any previous entry.
  cy.wait("@getDraft");

  cy.get(".journal-markdown").should("have.value", "");

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
  cy.get(".journal-markdown").should(
    "have.value",
    "# Example project\n\n* Item A\n* Item B"
  );
});
