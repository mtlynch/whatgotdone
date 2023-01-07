require("cypress-file-upload");

import "cypress-file-upload";

import "./commands";

// Reset the datastore before each test.
beforeEach(function () {
  cy.request("GET", "/api/testing/db/wipe");

  // TODO: Move this to per-test, as not all tests need initial data.
  cy.request("GET", "/api/testing/db/populate-dummy-data");
});
