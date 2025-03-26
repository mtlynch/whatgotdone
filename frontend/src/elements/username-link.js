class UsernameLink extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {
    if (!this.hasAttribute('username')) {
      console.error('username-link element requires username attribute');
      return;
    }

    const username = this.getAttribute('username');

    // Create the link
    const a = document.createElement('a');
    a.href = `/${username}`;
    a.textContent = username;
    a.style.fontWeight = 'bold';
    a.style.textDecoration = 'none';

    // Add any classes from the original element
    if (this.className) {
      a.className = this.className;
    }

    this.appendChild(a);
  }
}

// Register the custom element if it's not already defined
if (!customElements.get('username-link')) {
  customElements.define('username-link', UsernameLink);
}

// Export the class for testing or other uses
export default UsernameLink;
