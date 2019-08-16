// Retrieve the CSRF token from the meta tag in the DOM.
export default function getCsrfToken() {
  return document
    .querySelector("meta[name='csrf-token']")
    .getAttribute("content");
}