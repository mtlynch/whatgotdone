it("gets the sitemap", () => {
  cy.request("/sitemap.xml");
});

it("gets the robots.txt file", () => {
  cy.request("/robots.txt");
});
