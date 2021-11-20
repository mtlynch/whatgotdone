it("can export the logged-in user's data", () => {
  cy.login("staging_jimmy");

  cy.get(".account-dropdown").click();
  cy.get(".profile-link a").click();

  cy.location("pathname").should("eq", "/staging_jimmy");

  cy.get('[data-test-id="export-data-btn"]').click();
});

it("can't see an export data button for other users", () => {
  cy.login("staging_jimmy");

  cy.visit("/leader_lenny");

  cy.get('[data-test-id="export-data-btn"]').should("not.exist");
});
