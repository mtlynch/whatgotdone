it("follows a user", () => {
  cy.intercept("GET", "/api/draft/*").as("getDraft");
  cy.intercept("GET", "/api/user/*/following").as("getFollowing");
  cy.intercept("POST", "/api/logout").as("logout");

  // Log in as a leader user and create an entry.
  cy.login("leader_lenny");
  cy.location("pathname").should("include", "/entry/edit");
  cy.wait("@getDraft");

  cy.get(".editor-content .ProseMirror")
    .clear()
    .type("It's good to be the leader, as other users love to follow me!");
  cy.get("form").submit();
  cy.location("pathname").should("include", "/leader_lenny/");

  cy.visit("/logout");
  cy.wait("@logout");

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
  cy.get('[data-test-id="unfollow-btn"]').should("exist");
  cy.get('[data-test-id="follow-btn"]').should("not.exist");

  cy.visit("/feed");
  cy.wait("@getFollowing");
  cy.get(".journal").should("exist");

  cy.visit("/leader_lenny");
  cy.wait("@getFollowing");
  cy.get('[data-test-id="unfollow-btn"]').click();
  cy.get('[data-test-id="follow-btn"]').should("exist");
  cy.get('[data-test-id="unfollow-btn"]').should("not.exist");
});
