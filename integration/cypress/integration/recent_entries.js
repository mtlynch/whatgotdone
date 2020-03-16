it("views recent posts", () => {
  cy.visit("/recent");

  cy.get("div.journal").should("contain", "staging_jimmy's update");
});
