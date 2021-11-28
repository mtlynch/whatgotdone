it("logs in and reacts to an entry", () => {
  cy.intercept("POST", "https://api.userkit.io/v1/widget/login").as(
    "postUserKitLogin"
  );

  // Try reacting to an entry before logging in to ensure What Got Done prompts
  // a login.
  cy.visit("/staging_jimmy/2019-06-28");
  cy.get(".reaction-buttons .btn:first-of-type").click();
  cy.location("pathname").should("eq", "/login");
  cy.completeLoginForm("reacting_tommy");
  cy.wait("@postUserKitLogin");
  cy.location("pathname").should("eq", "/staging_jimmy/2019-06-28");

  // React to the entry
  cy.get(".reaction-buttons .btn:first-of-type").click();

  // TODO(mtlynch): We should really be selecting the *first* div.reaction element.
  cy.get(".reaction").then((element) => {
    expect(element.text().replace(/\s+/g, " ")).to.equal(
      "reacting_tommy reacted with a 👍"
    );
  });
});
