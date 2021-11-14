// Retrieve the CSRF token from the meta tag in the DOM.
export function getCsrfToken() {
  return document
    .querySelector("meta[name='csrf-token']")
    .getAttribute('content');
}
