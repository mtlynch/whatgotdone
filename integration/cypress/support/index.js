import "./commands";

// Reset the datastore before each test.
beforeEach(function () {
  cy.request("POST", Cypress.env("testDataManagerUrl") + "/reset");
});
