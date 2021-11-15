it("pre-populates new entries with the user's draft template", () => {
  cy.intercept("PUT", "/api/preferences").as("putPreferences");
  cy.intercept("GET", "/api/preferences").as("getPreferences");

  cy.visit("/preferences");
  cy.completeLoginForm("staging_jimmy");

  cy.location("pathname").should("equal", "/preferences");

  cy.get("#entry-template-input").type("# Project Falcon\n\n* Item 1");
  cy.get("form").submit();

  cy.wait("@putPreferences");

  cy.get(".alert").contains("Preferences saved");

  cy.get(".post-update").click();

  cy.location("pathname").should("include", "/entry/edit");

  cy.wait("@getPreferences");
  cy.get(".switch-mode .btn").click();
  cy.get(".markdown-editor textarea").should(
    "have.value",
    "# Project Falcon\n\n* Item 1"
  );
});
