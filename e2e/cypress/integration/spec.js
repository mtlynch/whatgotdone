Cypress.Commands.add('login', (username, password, options = {}) => {
  cy.visit('/login')

  // Temporary workaround: Wait a few ms before logging in, otherwise UserKit
  // fails with "Too many requests." This should go away in a few days when we
  // switch to the upcoming dummy app ID with no limits.
  cy.wait(250)

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

  cy.get('div.journal').should('contain', 'staging_jimmy\'s update')
})

it('clicking "Post Update" before authenticating prompts login', () => {
  cy.visit('/')

  cy.get('nav .account-dropdown').should('not.exist');

  cy.get('nav .post-update').click()

  cy.url().should('include', '/login')
})

it('reacting to an entry before authenticating prompts login', () => {
  cy.visit('/staging_jimmy/2019-06-28')

  cy.get('.reaction-buttons .btn:first-of-type').click();

  cy.url().should('include', '/login')
})

it('renders the date correctly', () => {
  cy.visit('/staging_jimmy/2019-06-28')

  cy.title().should('eq', 'staging_jimmy\'s What Got Done for the week of 2019-06-28')

  cy.get('.journal-header')
    .then((element) => {
      expect(element.text().replace(/\s+/g, ' ')).to.equal('staging_jimmy\'s update for the week ending on Friday, Jun 28, 2019');
    });
})

it('reaction buttons should not appear when the post is missing', () => {
  cy.visit('/staging_jimmy/2000-01-07')

  cy.get('.reaction-buttons').should('not.exist');
})

it('logs in and posts an update', () => {
  cy.server()
  cy.route('/api/draft/*').as('getDraft')

  cy.login('staging_jimmy', 'just4st@ginG!')

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

it('logs in and posts an empty update (deleting the update)', () => {
  cy.server()
  cy.route('/api/draft/*').as('getDraft')

  cy.login('staging_jimmy', 'just4st@ginG!')

  cy.url().should('include', '/entry/edit')

  // Wait for page to pull down any previous entry.
  cy.wait('@getDraft')

  cy.get('.journal-markdown').clear()
  cy.get('form').submit()

  cy.url().should('include', '/staging_jimmy/')
  cy.get('.missing-entry')
    .should('be.visible')
})

it('logs in and saves a draft', () => {
  cy.server()
  cy.route('GET', '/api/draft/*').as('getDraft')
  cy.route('POST', '/api/draft/*').as('postDraft')

  cy.login('staging_jimmy', 'just4st@ginG!')

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
  cy.login('staging_jimmy', 'just4st@ginG!')

  cy.get('.account-dropdown').click()
  cy.get('.profile-link a').click()

  cy.url().should('include', '/staging_jimmy')
})

it('logs in and reacts to an entry', () => {
  cy.server()
  cy.route('POST', 'https://api.userkit.io/v1/widget/login').as('postUserKitLogin')
  cy.login('staging_jimmy', 'just4st@ginG!')
  cy.wait('@postUserKitLogin')

  cy.visit('/staging_jimmy/2019-06-28')
    .its('document')
    .then((document) => {
      const csrfToken = document
        .querySelector("meta[name='csrf-token']")
        .getAttribute("content");

      // Clear any existing reaction on the entry.
      cy.request({
        url: '/api/reactions/entry/staging_jimmy/2019-06-28',
        method: 'POST',
        headers: {
          'X-Csrf-Token': csrfToken,
        },
        body: { reactionSymbol: "" },
      }).then(() => {
        cy.visit('/staging_jimmy/2019-06-28').then(() => {
          cy.get('.reaction-buttons .btn:first-of-type').click();

          // TODO(mtlynch): We should really be selecting the *first* div.reaction element.
          cy.get('.reaction')
            .then((element) => {
              expect(element.text().replace(/\s+/g, ' ')).to.equal('staging_jimmy reacted with a 👍');
            });
        })
      });
    });
})

it('logs in and signs out', () => {
  cy.login('staging_jimmy', 'just4st@ginG!')

  cy.url().should('include', '/entry/edit')

  cy.visit('/logout')
  cy.location('pathname').should('eq', '/')
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

  cy.login('staging_jimmy', 'just4st@ginG!')

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