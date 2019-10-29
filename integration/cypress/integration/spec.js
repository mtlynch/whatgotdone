Cypress.Commands.add('login', (username, password, options = {}) => {
  cy.visit('/login')

  cy.get('#userkit_username')
    .type(username)
  cy.get('#userkit_password')
    .type(password)
  cy.get('form').submit()
})

it('loads the homepage', () => {
  cy.visit('/')

  cy.get('h1')
    .should('contain', 'What did you get done this week?')
})


it('clicking "Post Update" before authenticating prompts login', () => {
  cy.visit('/')

  cy.get('nav .account-dropdown').should('not.exist');

  cy.get('nav .post-update').click()

  cy.url().should('include', '/login')
})

it('reaction buttons should not appear when the post is missing', () => {
  cy.visit('/staging_jimmy/2000-01-07')

  cy.get('.reaction-buttons').should('not.exist');
})

it('logs in and posts an update', () => {
  cy.server()
  cy.route('/api/draft/*').as('getDraft')

  cy.login('staging_jimmy', 'password')

  cy.url().should('include', '/entry/edit')

  // Wait for page to pull down any previous entry.
  cy.wait('@getDraft')

  const entryText = 'Posted an update at ' + new Date().toISOString();

  cy.get('.journal-markdown')
    .clear()
    .type(entryText)
  cy.get('form').submit()

  cy.url().should('include', '/staging_jimmy/')
  cy.get('.journal-body')
    .should('contain', entryText)
})

it('logs in and saves a draft', () => {
  cy.server()
  cy.route('GET', '/api/draft/*').as('getDraft')
  cy.route('POST', '/api/draft/*').as('postDraft')

  cy.login('staging_jimmy', 'password')

  cy.url().should('include', '/entry/edit')

  // Wait for page to pull down any previous entry.
  cy.wait('@getDraft')

  const entryText = 'Saved a private draft at ' + new Date().toISOString();

  cy.get('.journal-markdown')
    .clear()
    .type(entryText)
  cy.get('.save-draft').click()

  // Wait for "save draft" operation to complete.
  cy.wait('@postDraft')

  // User should stay on the same page after saving a draft.
  cy.url().should('include', '/entry/edit')

  cy.visit('/recent')

  // Private drafts should not appear on the recent page
  cy.get('#app').should('not.contain', entryText)
})

it('logs in and views profile', () => {
  cy.login('staging_jimmy', 'password')

  cy.get('.account-dropdown').click()
  cy.get('.profile-link a').click()

  cy.url().should('include', '/staging_jimmy')
})

it('logs in and signs out', () => {
  cy.login('staging_jimmy', 'password')

  cy.url().should('include', '/entry/edit')

  cy.visit('/logout')
  cy.location('pathname').should('eq', '/')

  cy.get('nav .account-dropdown').should('not.exist');
})

it('views a non-existing user profile with empty information', () => {
  cy.visit('/nonExistentUser')
  cy.get('.no-bio-message')
    .should('be.visible')
  cy.get('.no-entries-message')
    .should('be.visible')
})

it('logs in updates profile', () => {
  cy.server()
  cy.route('/api/user/staging_jimmy').as('getUserProfile')

  cy.login('staging_jimmy', 'password')

  cy.url().should('include', '/entry/edit')

  cy.visit('/staging_jimmy')
  cy.get('.edit-btn').click()

  // Wait for page to pull down existing profile.
  cy.wait('@getUserProfile')

  cy.get('#user-bio')
    .clear()
    .type("Hello, my name is staging_jimmy!")

  cy.get('#email-address')
    .clear()
    .type("staging_jimmy@example.com")

  cy.get('#twitter-handle')
    .clear()
    .type("jack")

  cy.get('#save-profile').click()
  cy.url().should('include', '/staging_jimmy')

  cy.get('.user-bio')
    .should('contain', "Hello, my name is staging_jimmy!")
  cy.get('.email-address')
    .should('contain', "staging_jimmy@example.com")
  cy.get('.twitter-handle')
    .should('contain', "jack")
})