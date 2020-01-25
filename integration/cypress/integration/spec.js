Cypress.Commands.add("login", (username, password, options = {}) => {
  cy.visit("/login");

  cy.get("#userkit_username").type(username);
  cy.get("#userkit_password").type(password);
  cy.get("form").submit();
});

it("loads the homepage", () => {
  cy.visit("/");

  // Verify <head> metadata.
  cy.title().should("include", "What Got Done");
  cy.get('meta[name="description"]').should(
    "have.attr",
    "content",
    "The simple, easy way to share progress with your teammates."
  );
  cy.get('meta[property="og:type"]').should("have.attr", "content", "website");
  cy.get('meta[property="og:title"]').should(
    "have.attr",
    "content",
    "What Got Done"
  );
  cy.get('meta[property="og:description"]').should(
    "have.attr",
    "content",
    "The simple, easy way to share progress with your teammates."
  );

  cy.get("h1").should("contain", "What did you get done this week?");
});

it('clicking "Post Update" before authenticating prompts login', () => {
  cy.visit("/");

  cy.get("nav .account-dropdown").should("not.exist");

  cy.get("nav .post-update").click();

  cy.location("pathname").should("eq", "/login");
});

it("reaction buttons should not appear when the post is missing", () => {
  cy.visit("/staging_jimmy/2000-01-07");

  cy.get(".reaction-buttons").should("not.exist");
});

it("logs in and posts an update", () => {
  cy.server();
  cy.route("/api/draft/*").as("getDraft");

  cy.login("staging_jimmy", "password");

  cy.location("pathname").should("include", "/entry/edit");

  // Wait for page to pull down any previous entry.
  cy.wait("@getDraft");

  const entryText = "Posted an update at " + new Date().toISOString();

  cy.get(".journal-markdown")
    .clear()
    .type(entryText);
  cy.get("form").submit();

  cy.location("pathname").should("include", "/staging_jimmy/");
  // Reload the page to fetch the new HTML rather than using what the front-end
  // generated client-side.
  cy.reload();

  // Verify <head> metadata.
  cy.title().should("include", "staging_jimmy's What Got Done for the week of");
  cy.get('meta[property="og:type"]').should("have.attr", "content", "article");
  cy.get(".journal-body").should("contain", entryText);
});

it("logs in and backdates an update from a previous week", () => {
  cy.server();
  cy.route("/api/draft/*").as("getDraft");

  cy.login("staging_jimmy", "password");

  // Wait for page to pull down any previous entry.
  cy.wait("@getDraft");

  cy.visit("/entry/edit/2019-12-13");

  const entryText = "Posted an update at " + new Date().toISOString();

  cy.get(".journal-markdown")
    .clear()
    .type(entryText);
  cy.get("form").submit();

  cy.location("pathname").should("eq", "/staging_jimmy/2019-12-13");
  // Reload the page to fetch the new HTML rather than using what the front-end
  // generated client-side.
  cy.reload();

  // Verify <head> metadata.
  cy.title().should(
    "include",
    "staging_jimmy's What Got Done for the week of 2019-12-13"
  );
  cy.get('meta[name="description"]').should(
    "have.attr",
    "content",
    "Find out what staging_jimmy accomplished for the week ending on December 13, 2019"
  );
  cy.get('meta[property="og:type"]').should("have.attr", "content", "article");
  cy.get('meta[property="og:title"]').should(
    "have.attr",
    "content",
    "staging_jimmy's What Got Done for the week of Dec. 13, 2019"
  );
  cy.get('meta[property="og:description"]').should(
    "have.attr",
    "content",
    "Find out what staging_jimmy accomplished for the week ending on December 13, 2019"
  );

  cy.get(".journal-body").should("contain", entryText);
});

it("logs in and saves a draft", () => {
  cy.server();
  cy.route("GET", "/api/draft/*").as("getDraft");
  cy.route("POST", "/api/draft/*").as("postDraft");

  cy.login("staging_jimmy", "password");

  cy.location("pathname").should("include", "/entry/edit");

  // Wait for page to pull down any previous entry.
  cy.wait("@getDraft");

  const entryText = "Saved a private draft at " + new Date().toISOString();

  cy.get(".journal-markdown")
    .clear()
    .type(entryText);
  cy.get(".save-draft").click();

  // Wait for "save draft" operation to complete.
  cy.wait("@postDraft");

  // User should stay on the same page after saving a draft.
  cy.location("pathname").should("include", "/entry/edit");

  cy.visit("/recent");

  // Private drafts should not appear on the recent page
  cy.get("#app").should("not.contain", entryText);
});

it("logs in and views profile", () => {
  cy.login("staging_jimmy", "password");

  cy.get(".account-dropdown").click();
  cy.get(".profile-link a").click();

  cy.location("pathname").should("eq", "/staging_jimmy");
});

it("logs in and signs out", () => {
  cy.login("staging_jimmy", "password");

  cy.location("pathname").should("include", "/entry/edit");

  cy.visit("/logout");
  cy.location("pathname").should("eq", "/");

  cy.get("nav .account-dropdown").should("not.exist");
});

it("views a non-existing user profile with empty information", () => {
  cy.visit("/nonExistentUser");
  cy.get(".no-bio-message").should("be.visible");
  cy.get(".no-entries-message").should("be.visible");
});

it("logs in updates profile", () => {
  cy.server();
  cy.route("/api/user/staging_jimmy").as("getUserProfile");

  cy.login("staging_jimmy", "password");

  cy.location("pathname").should("include", "/entry/edit");

  cy.visit("/staging_jimmy");
  cy.get(".edit-btn").click();

  // Wait for page to pull down existing profile.
  cy.wait("@getUserProfile");

  cy.get("#user-bio")
    .clear()
    .type("Hello, my name is staging_jimmy!");

  cy.get("#email-address")
    .clear()
    .type("staging_jimmy@example.com");

  cy.get("#twitter-handle")
    .clear()
    .type("jack");

  cy.get("#save-profile").click();
  cy.location("pathname").should("eq", "/staging_jimmy");

  cy.get(".user-bio").should("contain", "Hello, my name is staging_jimmy!");
  cy.get(".email-address").should("contain", "staging_jimmy@example.com");
  cy.get(".twitter-handle").should("contain", "jack");
});
