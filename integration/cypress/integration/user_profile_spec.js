it("logs in and views profile", () => {
  cy.login("staging_jimmy");

  cy.get(".account-dropdown").click();
  cy.get(".profile-link a").click();

  cy.location("pathname").should("eq", "/staging_jimmy");
});

it("gets 404 for non-existent user's profile page", () => {
  cy.request({
    url: "/nonExistentUser",
    failOnStatusCode: false,
  }).then((response) => {
    expect(response.status).to.eq(404);
  });
});

it("logs in and updates profile", () => {
  cy.intercept("/api/user/staging_jimmy").as("getUserProfile");

  cy.login("staging_jimmy");

  cy.location("pathname").should("include", "/entry/edit");

  cy.visit("/staging_jimmy");
  cy.get('[data-test-id="edit-btn"]').click();

  // Wait for page to pull down existing profile.
  cy.wait("@getUserProfile");

  cy.get("#user-bio").clear().type("Hello, my name is staging_jimmy!");

  cy.get("#email-address").clear().type("staging_jimmy@example.com");

  cy.get("#twitter-handle").clear().type("jack");

  cy.get("#save-profile").click();
  cy.location("pathname").should("eq", "/staging_jimmy");

  cy.get(".user-bio").should("contain", "Hello, my name is staging_jimmy!");
  cy.get(".email-address").should("contain", "staging_jimmy@example.com");
  cy.get(".twitter-handle").should("contain", "jack");
});

it("logs in and sets profile photo", () => {
  cy.intercept("/api/user/staging_jimmy").as("getUserProfile");
  cy.intercept("PUT", "/api/user/avatar").as("setProfilePhoto");
  cy.intercept("DELETE", "/api/user/avatar").as("deleteProfilePhoto");

  cy.login("staging_jimmy");

  cy.location("pathname").should("include", "/entry/edit");

  cy.visit("/staging_jimmy");
  cy.get('[data-test-id="edit-btn"]').click();

  // Wait for page to pull down existing profile.
  cy.wait("@getUserProfile");

  cy.get("#upload-profile-photo").attachFile("kittyface.jpg");
  cy.wait("@setProfilePhoto");

  cy.visit("/staging_jimmy");

  cy.get('[data-test-id="edit-btn"]').click();
  // Wait for page to pull down existing profile.
  cy.wait("@getUserProfile");

  cy.get("#delete-profile-photo").click();
  cy.wait("@deleteProfilePhoto");
});
