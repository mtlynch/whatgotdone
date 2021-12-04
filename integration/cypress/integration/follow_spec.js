it("follows a user", () => {
  cy.intercept("GET", "/api/draft/*").as("getDraft");
  cy.intercept("GET", "/api/user/follower_frank/following").as("getFollowing");
  cy.intercept("PUT", "/api/follow/leader_lenny").as("putFollow");
  cy.intercept("DELETE", "/api/follow/leader_lenny").as("deleteFollow");
  cy.intercept("POST", "/api/logout").as("logout");

  // Log in as a follow user to follow.
  cy.login("follower_frank");
  cy.location("pathname").should("include", "/entry/edit");

  cy.visit("/feed");
  cy.wait("@getFollowing");

  // Verify the personalized feed is empty.
  cy.get(".alert").should("contain", "You're not following anyone yet.");
  cy.get(".journal").should("not.exist");

  cy.visit("/leader_lenny");
  cy.wait("@getFollowing");
  cy.get('[data-test-id="follow-btn"]').click();
  cy.wait("@putFollow");
  cy.get('[data-test-id="unfollow-btn"]').should("exist");
  cy.get('[data-test-id="follow-btn"]').should("not.exist");

  cy.visit("/feed");
  cy.wait("@getFollowing");
  cy.get(".journal").should("exist");

  cy.visit("/leader_lenny");
  cy.wait("@getFollowing");
  cy.get('[data-test-id="unfollow-btn"]').click();
  cy.wait("@deleteFollow");
  cy.get('[data-test-id="follow-btn"]').should("exist");
  cy.get('[data-test-id="unfollow-btn"]').should("not.exist");
});
