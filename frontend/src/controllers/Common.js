// Retrieve the CSRF token from the meta tag in the DOM.
export function getCsrfToken() {
  return document
    .querySelector("meta[name='csrf-token']")
    .getAttribute('content');
}

export function processJsonResponse(response) {
  if (response.ok) {
    return response.json();
  }
  return response.text().then((error) => {
    return Promise.reject(error);
  });
}
