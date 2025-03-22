class UserReactions extends HTMLElement {
  constructor() {
    super();
    this.reactionSymbols = ['👍', '🎉', '🙁'];
  }

  connectedCallback() {
    if (
      !this.hasAttribute('entry-author') ||
      !this.hasAttribute('entry-date')
    ) {
      console.error(
        'user-reactions element requires entry-author and entry-date attributes'
      );
      return;
    }

    // Initial load of reactions if the element doesn't have content yet
    if (this.children.length === 0) {
      this.loadReactions();
    }

    // Set up event delegation for reaction buttons if not using htmx
    this.addEventListener('click', (e) => {
      if (e.target.matches('button.btn') && !e.target.hasAttribute('hx-post')) {
        this.handleReactionClick(e.target.textContent.trim());
      }
    });
  }

  get entryAuthor() {
    return this.getAttribute('entry-author');
  }

  get entryDate() {
    return this.getAttribute('entry-date');
  }

  get loggedInUsername() {
    return this.getAttribute('logged-in-username') || '';
  }

  // Fetch reactions and update the element (fallback for non-htmx scenarios)
  async loadReactions() {
    if (!this.entryAuthor || !this.entryDate) return;

    try {
      const response = await fetch(
        `/api/reactions/entry/${this.entryAuthor}/${this.entryDate}`,
        {
          headers: {
            Accept: 'text/html',
          },
        }
      );
      if (!response.ok) throw new Error('Failed to load reactions');

      const html = await response.text();
      this.innerHTML = html;
    } catch (err) {
      console.error('Error loading reactions:', err);
    }
  }

  // Handle reaction button clicks (fallback for non-htmx scenarios)
  async handleReactionClick(reactionSymbol) {
    if (!this.loggedInUsername) {
      window.location.href = '/login';
      return;
    }

    // Determine if we're adding or removing a reaction
    const isSelected = Array.from(
      this.querySelectorAll('button.btn-light')
    ).some((btn) => btn.textContent.trim() === reactionSymbol);

    try {
      if (isSelected) {
        // Delete reaction
        await fetch(
          `/api/reactions/entry/${this.entryAuthor}/${this.entryDate}`,
          {
            method: 'DELETE',
            credentials: 'include',
            headers: {
              Accept: 'text/html',
              'X-CSRF-Token': this.getCsrfToken(),
            },
          }
        );
      } else {
        // Add reaction
        await fetch(
          `/api/reactions/entry/${this.entryAuthor}/${this.entryDate}`,
          {
            method: 'POST',
            credentials: 'include',
            headers: {
              'Content-Type': 'application/json',
              Accept: 'text/html',
              'X-CSRF-Token': this.getCsrfToken(),
            },
            body: JSON.stringify({
              reactionSymbol: reactionSymbol,
            }),
          }
        );
      }
      // Reload reactions
      this.loadReactions();
    } catch (err) {
      console.error('Error updating reaction:', err);
    }
  }

  // Helper to get CSRF token
  getCsrfToken() {
    const metaTag = document.querySelector('meta[name="csrf-token"]');
    return metaTag ? metaTag.getAttribute('content') : '';
  }
}

// Register the custom element if it's not already defined
if (!customElements.get('user-reactions')) {
  customElements.define('user-reactions', UserReactions);
}

// Export the class for testing or other uses
export default UserReactions;
