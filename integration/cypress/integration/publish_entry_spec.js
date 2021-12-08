it("logs in and posts an update", () => {
  cy.intercept("/api/draft/*").as("getDraft");

  cy.login("staging_jimmy");

  cy.location("pathname").should("include", "/entry/edit");

  // Wait for page to pull down any previous entry.
  cy.wait("@getDraft");

  const entryText = "Posted an update at " + new Date().toISOString();

  cy.get(".editor-content .ProseMirror").clear().type(entryText);
  cy.get("form").submit();

  cy.location("pathname").should("include", "/staging_jimmy/");
  // Reload the page to fetch the new HTML rather than using what the front-end
  // generated client-side.
  cy.reload();

  // Verify <head> metadata.
  cy.title().should("include", "staging_jimmy's What Got Done for the week of");
  cy.get('meta[property="og:type"]').should("have.attr", "content", "article");

  cy.get(".journal-body").should("contain", entryText);
  cy.get(".view-count").should("contain", "Viewed 1 times");
  cy.get(".missing-entry").should("not.exist");
});

it("logs in and backdates an update from a previous week", () => {
  cy.intercept("/api/draft/*").as("getDraft");

  cy.login("staging_jimmy");

  // Wait for page to pull down any previous entry.
  cy.wait("@getDraft");

  cy.visit("/entry/edit/2019-12-13");

  const entryText = "Posted an update at " + new Date().toISOString();

  cy.get(".editor-content .ProseMirror").clear().type(entryText);
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
  cy.get(".view-count").should("contain", "Viewed 1 times");
  cy.get(".missing-entry").should("not.exist");
});

it("posts an update and then unpublishes it", () => {
  cy.intercept("GET", "/api/draft/2019-06-28").as("getDraft");

  cy.login("staging_jimmy");

  cy.location("pathname").should("include", "/entry/edit");
  cy.visit("/entry/edit/2019-06-28");

  // Wait for page to pull down the previous entry.
  cy.wait("@getDraft");

  cy.get(".editor-content .ProseMirror")
    .clear()
    .type("felt cute, might unpublish later");
  cy.get("form").submit();

  cy.location("pathname").should("eq", "/staging_jimmy/2019-06-28");
  cy.get(".missing-entry").should("not.exist");

  cy.get("[data-test-id='unpublish-btn']").click();

  // Unpublishing takes the user back to the edit entry page.
  cy.location("pathname").should("eq", "/entry/edit/2019-06-28");

  // Go back and make sure the published entry is gone.
  cy.go("back");
  cy.location("pathname").should("eq", "/staging_jimmy/2019-06-28");
  cy.get(".missing-entry").should("be.visible");
});
