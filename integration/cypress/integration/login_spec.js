it('clicking "Post Update" before authenticating prompts login', () => {
  cy.visit("/");

  cy.get("nav .account-dropdown").should("not.exist");

  cy.get("nav .post-update").click();

  cy.location("pathname").should("eq", "/login");
});

it("back button should work if the user decides not to login/sign up", () => {
  cy.visit("/");
  cy.get("nav .post-update").click();

  cy.location("pathname").should("eq", "/login");

  cy.go(-1);

  cy.location("pathname").should("eq", "/");
});

it("logs in and signs out", () => {
  cy.login("staging_jimmy");

  cy.location("pathname").should("include", "/entry/edit");

  cy.get(".account-dropdown").click();
  cy.get(".sign-out-link a").click();
  cy.location("pathname").should("eq", "/");

  cy.get("nav .account-dropdown").should("not.exist");
});

it("bare route should redirect authenticated user to their edit entry page", () => {
  cy.server();
  cy.route("/api/user/me").as("getUsername");

  cy.visit("/");
  cy.location("pathname").should("eq", "/");

  // Clicking the navbar brand should point to /about page.
  cy.get(".navbar .navbar-brand").click();
  cy.location("pathname").should("eq", "/");

  cy.login("staging_jimmy");
  cy.location("pathname").should("include", "/entry/edit");
  cy.wait("@getUsername");

  // Navigating back to bare route should redirect to edit entry page.
  cy.visit("/");
  cy.reload();
  cy.location("pathname").should("contain", "/entry/edit/");

  // Clicking navbar brand should point to edit entry page.
  cy.get(".navbar .navbar-brand").click();
  cy.location("pathname").should("contain", "/entry/edit/");

  // Log out
  cy.get(".account-dropdown").click();
  cy.get(".sign-out-link a").click();

  // Post-logout, user should be on bare route.
  cy.location("pathname").should("eq", "/");

  // Clicking navbar brand should point to bare route.
  cy.get(".navbar .navbar-brand").click();
  cy.location("pathname").should("eq", "/");
});
