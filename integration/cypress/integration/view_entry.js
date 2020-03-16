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
