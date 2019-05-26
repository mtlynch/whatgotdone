it('loads the homepage', () => {
  cy.visit('/')

  cy.get('h1')
    .should('contain', 'What did you get done this week?')
})

it('views recent posts', () => {
  cy.visit('/recent')

  cy.get('div.journal').should('contain', 'staging.jimmy\'s update')
})

it('logs in and posts an update', () => {
  cy.server()
  cy.route('/api/entries/*').as('getEntries')

  cy.visit('/login')

  cy.get('#userkit_username')
    .type('staging.jimmy')
  cy.get('#userkit_password')
    .type('just4st@ginG!')
  cy.get('form').submit()

  cy.url().should('include', '/submit')

  // Wait for page to pull down any previous entry.
  cy.wait('@getEntries')

  cy.get('.journal-markdown')
    .clear()
    .type('Posted an update at ' + new Date().toISOString())
})

it('logs in and signs out', () => {
  cy.server()
  cy.route('/api/entries/*').as('getEntries')

  cy.visit('/login')

  cy.get('#userkit_username')
    .type('staging.jimmy')
  cy.get('#userkit_password')
    .type('just4st@ginG!')
  cy.get('form').submit()

  cy.url().should('include', '/submit')

  cy.visit('/logout')
  cy.url().should('include', '/login')
})