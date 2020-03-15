Cypress.Commands.add("login", (username, options = {}) => {
  cy.visit("/login");

  cy.get("#userkit_username").type(username);
  cy.get("#userkit_password").type("password");
  cy.get("form").submit();
});

it("views recent posts", () => {
  cy.visit("/recent");

  cy.get("div.journal").should("contain", "staging_jimmy's update");
});

it("reacting to an entry before authenticating prompts login", () => {
  cy.visit("/staging_jimmy/2019-06-28");

  cy.get(".reaction-buttons .btn:first-of-type").click();

  cy.location("pathname").should("eq", "/login");
});

it("renders the date correctly", () => {
  cy.visit("/staging_jimmy/2019-06-28");

  cy.title().should(
    "eq",
    "staging_jimmy's What Got Done for the week of 2019-06-28"
  );

  cy.get(".journal-header").then(element => {
    expect(element.text().replace(/\s+/g, " ")).to.equal(
      "staging_jimmy's update for the week ending on Friday, Jun 28, 2019"
    );
  });
});

it("logs in and posts an empty update (deleting the update)", () => {
  cy.server();
  cy.route("/api/draft/*").as("getDraft");

  cy.login("staging_jimmy");

  cy.location("pathname").should("include", "/entry/edit");

  // Wait for page to pull down any previous entry.
  cy.wait("@getDraft");

  cy.get(".journal-markdown").clear();
  cy.get("form").submit();

  cy.location("pathname").should("include", "/staging_jimmy/");
  cy.get(".missing-entry").should("be.visible");
});

it("logs in and reacts to an entry", () => {
  cy.server();
  cy.route("POST", "https://api.userkit.io/v1/widget/login").as(
    "postUserKitLogin"
  );
  cy.login("reacting_tommy");
  cy.wait("@postUserKitLogin");

  cy.visit("/staging_jimmy/2019-06-28")
    .its("document")
    .then(document => {
      const csrfToken = document
        .querySelector("meta[name='csrf-token']")
        .getAttribute("content");

      // Clear any existing reaction on the entry.
      cy.request({
        url: "/api/reactions/entry/staging_jimmy/2019-06-28",
        method: "POST",
        headers: {
          "X-Csrf-Token": csrfToken
        },
        body: { reactionSymbol: "" }
      }).then(() => {
        cy.visit("/staging_jimmy/2019-06-28").then(() => {
          cy.get(".reaction-buttons .btn:first-of-type").click();

          // TODO(mtlynch): We should really be selecting the *first* div.reaction element.
          cy.get(".reaction").then(element => {
            expect(element.text().replace(/\s+/g, " ")).to.equal(
              "reacting_tommy reacted with a 👍"
            );
          });
        });
      });
    });
});

it("reacting to an entry prompts login and redirects back to the entry", () => {
  cy.server();
  cy.route("POST", "https://api.userkit.io/v1/widget/login").as(
    "postUserKitLogin"
  );

  cy.visit("/staging_jimmy/2019-06-28");

  cy.get(".reaction-buttons .btn:first-of-type").click();

  // Do login (can't use login method because it performs a visit())
  cy.get("#userkit_username").type("reacting_tommy");
  cy.get("#userkit_password").type("password");
  cy.get("form").submit();
  cy.wait("@postUserKitLogin");

  cy.location("pathname").should("eq", "/staging_jimmy/2019-06-28");
});
