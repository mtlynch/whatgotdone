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

it('reacting to an entry before authenticating prompts login', () => {
  cy.visit('/staging.jimmy/2019-06-28')

  cy.get('.reaction-buttons .btn:first-of-type').click();

  cy.url().should('include', '/login')
})

it('renders the date correctly', () => {
  cy.visit('/staging.jimmy/2019-06-28')

  cy.get('.journalHeader')
    .then((element) => {
      expect(element.text().replace(/\s+/g, ' ')).to.equal('staging.jimmy\'s update for the week ending on Friday, Jun 28, 2019');
    });
})

it('reaction buttons should not appear when the user has no posts', () => {
  cy.visit('/dummyUserWithZeroPosts')

  cy.get('.reaction-buttons').should('not.exist');
})

it('reaction buttons should not appear when the post is missing', () => {
  cy.visit('/staging.jimmy/2000-01-07')

  cy.get('.reaction-buttons').should('not.exist');
})

it('logs in and posts an update', () => {
  cy.server()
  cy.route('/api/draft/*').as('getDraft')

  cy.login('staging.jimmy', 'just4st@ginG!')

  cy.url().should('include', '/submit')

  // Wait for page to pull down any previous entry.
  cy.wait('@getDraft')

  const entryText = 'Posted an update at ' + new Date().toISOString();

  cy.get('.journal-markdown')
    .clear()
    .type(entryText)
  cy.get('form').submit()

  cy.url().should('include', '/staging.jimmy/')
  cy.get('.journal-body')
    .should('contain', entryText)
})

it('logs in and saves a draft', () => {
  cy.server()
  cy.route('GET', '/api/draft/*').as('getDraft')
  cy.route('POST', '/api/draft/*').as('postDraft')

  cy.login('staging.jimmy', 'just4st@ginG!')

  cy.url().should('include', '/submit')

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
  cy.url().should('include', '/submit')

  cy.visit('/recent')

  // Private drafts should not appear on the recent page
  cy.get('#app').should('not.contain', entryText)
})

it('logs in and views profile', () => {
  cy.login('staging.jimmy', 'just4st@ginG!')

  cy.get('.account-dropdown').click()
  cy.get('.profile-link a').click()

  cy.url().should('include', '/staging.jimmy')
})

it('logs in and reacts to an entry', () => {
  cy.server()
  cy.route('POST', 'https://api.userkit.io/v1/widget/login').as('postUserKitLogin')
  cy.login('staging.jimmy', 'just4st@ginG!')
  cy.wait('@postUserKitLogin')

  // Clear any existing reaction on the entry.
  cy.request('POST', '/api/reactions/entry/staging.jimmy/2019-06-28', { reactionSymbol: "" }).then(() => {

    // Visit a post and react to the entry.
    cy.visit('/staging.jimmy/2019-06-28')

    cy.get('.reaction-buttons .btn:first-of-type').click();

    // TODO(mtlynch): We should really be selecting the *first* div.reaction element.
    cy.get('.reaction')
      .then((element) => {
        expect(element.text().replace(/\s+/g, ' ')).to.equal('staging.jimmy reacted with a 👍');
      });
  });

})

it('logs in and signs out', () => {
  cy.login('staging.jimmy', 'just4st@ginG!')

  cy.url().should('include', '/submit')

  cy.visit('/logout')
  cy.url().should('include', '/login')
  cy.get('#userkit_username')
})