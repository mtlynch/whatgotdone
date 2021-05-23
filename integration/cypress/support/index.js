require("cypress-file-upload");

import "cypress-file-upload";

import "./commands";

// Reset the datastore before each test.
beforeEach(function () {
  cy.request("POST", Cypress.env("testDataManagerUrl") + "/reset");
});
