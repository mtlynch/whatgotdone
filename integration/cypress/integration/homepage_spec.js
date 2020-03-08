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
