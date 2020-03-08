it("logs in and posts an update", () => {
  cy.server();
  cy.route("/api/draft/*").as("getDraft");

  cy.login("staging_jimmy");

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

  cy.login("staging_jimmy");

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
