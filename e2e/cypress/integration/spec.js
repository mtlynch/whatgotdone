it('loads the homepage', () => {
  cy.visit('/')

  cy.get('h1')
    .should('contain', 'What did you get done this week?')
})