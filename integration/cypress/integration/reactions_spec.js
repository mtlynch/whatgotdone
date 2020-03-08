it("reaction buttons should not appear when the post is missing", () => {
  cy.visit("/staging_jimmy/2000-01-07");

  cy.get(".reaction-buttons").should("not.exist");
});
