require("cypress-file-upload");

import "cypress-file-upload";

import "./commands";

// Reset the datastore before each test.
beforeEach(function () {
  cy.request("GET", Cypress.env("baseUrl") + "/api/testing/db/wipe");
  cy.request(
    "GET",
    Cypress.env("baseUrl") + "/api/testing/db/populate-dummy-data"
  );
});
