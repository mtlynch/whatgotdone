it("pre-populates new entries with the user's draft template", () => {
  cy.server();
  cy.route("POST", "/api/preferences").as("postPreferences");
  cy.route("GET", "/api/preferences").as("getPreferences");

  cy.visit("/preferences");
  cy.completeLoginForm("staging_jimmy");

  cy.location("pathname").should("equal", "/preferences");

  cy.get("#entry-template-input").type("# Project Falcon\n\n* Item 1");
  cy.get("form").submit();

  cy.wait("@postPreferences");

  cy.get(".alert").contains("Preferences saved");

  cy.get(".post-update").click();

  cy.location("pathname").should("include", "/entry/edit");

  cy.wait("@getPreferences");
  cy.get(".journal-markdown").should(
    "have.value",
    "# Project Falcon\n\n* Item 1"
  );
});
