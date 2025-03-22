class UserAvatar extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {
    if (!this.hasAttribute('username')) {
      console.error('user-avatar element requires username attribute');
      return;
    }

    const username = this.getAttribute('username');
    const size = this.getAttribute('size') || '40px';

    // Create the avatar image
    const img = document.createElement('img');
    img.src = `/${username}/avatar`;
    img.alt = `${username}'s avatar`;
    img.style.width = size;
    img.style.height = size;
    img.style.borderRadius = '50%';
    img.style.objectFit = 'cover';

    // Add any classes from the original element
    if (this.className) {
      img.className = this.className;
    }

    this.appendChild(img);
  }
}

// Register the custom element if it's not already defined
if (!customElements.get('user-avatar')) {
  customElements.define('user-avatar', UserAvatar);
}

// Export the class for testing or other uses
export default UserAvatar;
