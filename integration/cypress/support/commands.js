Cypress.Commands.add("login", (username, options = {}) => {
  cy.visit("/login");

  cy.get("#userkit_username").type(username);
  cy.get("#userkit_password").type("password");
  cy.get("form").submit();
});
