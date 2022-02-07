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

export function processBlankResponse(response) {
  if (!response.ok) {
    return response.text().then((error) => {
      return Promise.reject(error);
    });
  }
  return Promise.resolve();
}
