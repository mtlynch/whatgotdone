it('loads the homepage', () => {
  cy.visit('/')

  cy.get('h1')
    .should('contain', 'What did you get done this week?')
})

it('views recent posts', () => {
  cy.visit('/recent')

  cy.get('div.journal').should('contain', 'staging.jimmy\'s update')
})

it('clicking "Post Update" before authenticating prompts login', () => {
  cy.visit('/')

  cy.get('nav .account-dropdown').should('not.exist');

  cy.get('nav .post-update').click()

  cy.url().should('include', '/login')
})

it('logs in and posts an update', () => {
  cy.server()
  cy.route('/api/entry/staging.jimmy/*').as('getEntry')

  cy.visit('/login')

  cy.get('#userkit_username')
    .type('staging.jimmy')
  cy.get('#userkit_password')
    .type('just4st@ginG!')
  cy.get('form').submit()

  cy.url().should('include', '/submit')

  // Wait for page to pull down any previous entry.
  cy.wait('@getEntry')

  const entryText = 'Posted an update at ' + new Date().toISOString();

  cy.get('.journal-markdown')
    .clear()
    .type(entryText)
  cy.get('form').submit()

  cy.url().should('include', '/staging.jimmy/')
  cy.get('.journal-body')
    .should('contain', entryText)
})

it('logs in and views profile', () => {
  cy.visit('/login')
  cy.get('#userkit_username')
    .type('staging.jimmy')
  cy.get('#userkit_password')
    .type('just4st@ginG!')
  cy.get('form').submit()

  cy.get('.account-dropdown').click()
  cy.get('.profile-link a').click()

  cy.url().should('include', '/staging.jimmy')
})

it('logs in and signs out', () => {
  cy.visit('/login')

  cy.get('#userkit_username')
    .type('staging.jimmy')
  cy.get('#userkit_password')
    .type('just4st@ginG!')
  cy.get('form').submit()

  cy.url().should('include', '/submit')

  cy.visit('/logout')
  cy.url().should('include', '/login')
  cy.get('#userkit_username')
})