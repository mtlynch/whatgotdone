it("/preferences redirects unauthenticated users to login page", () => {
  cy.visit("/preferences");

  cy.location("pathname").should("eq", "/login");
});

it("/preferences allows authenticated users to stay on page", () => {
  cy.login("staging_jimmy");

  cy.location("pathname").should("include", "/entry/edit");

  cy.visit("/preferences");

  cy.get("h1").should("contain", "Preferences");
});
