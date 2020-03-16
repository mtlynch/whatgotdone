Cypress.Commands.add("completeLoginForm", (username, options = {}) => {
  cy.get("#userkit_username").type(username);
  cy.get("#userkit_password").type("password"); // Test-mode password is 'password'.
  cy.get("form").submit();
});

Cypress.Commands.add("login", (username, options = {}) => {
  cy.visit("/login");
  cy.completeLoginForm(username, options);
});
