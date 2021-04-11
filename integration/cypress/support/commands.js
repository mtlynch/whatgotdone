import "./selectionCommand";

Cypress.Commands.add("completeLoginForm", (username, options = {}) => {
  cy.get("#userkit_username").type(username);
  cy.get("#userkit_password").type("password"); // Test-mode password is 'password'.
  cy.get("form").submit();
});

Cypress.Commands.add("login", (username, options = {}) => {
  cy.visit("/login");
  cy.completeLoginForm(username, options);
});

Cypress.Commands.add("routeShouldBeCalled", (alias, timesCalled) => {
  expect(
    cy.state("requests").filter((call) => call.alias === alias),
    `${alias} should have been called ${timesCalled} times`
  ).to.have.length(timesCalled);
});
