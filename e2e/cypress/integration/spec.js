it('loads the homepage', () => {
  cy.visit('/')

  cy.get('h1')
    .should('contain', 'What did you get done this week?')
})

it('views recent posts', () => {
  cy.visit('/recent')

  cy.get('div.journal').should('contain', 'staging.jimmy\'s update')
})

it('gets recent entries by API', () => {
  cy.request('/api/recentEntries')
})

it('logs in', () => {
  cy.visit('/login')

  cy.get('#userkit_username')
    .type('staging.jimmy')
  cy.get('#userkit_password')
    .type('just4st@ginG!')
  cy.get('form').submit()

  cy.url().should('include', '/staging.jimmy')
})