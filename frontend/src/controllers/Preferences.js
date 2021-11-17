import {getCsrfToken} from '@/controllers/Common.js';

export function getPreferences() {
  return fetch(`${process.env.VUE_APP_BACKEND_URL}/api/preferences`, {
    credentials: 'include',
  }).then((response) => {
    // A 404 is not an error.
    if (response.status === 404) {
      return Promise.resolve({});
    } else if (!response.ok) {
      return response.text().then((error) => {
        return Promise.reject(error);
      });
    }
    return response.json();
  });
}

export function savePreferences(preferences) {
  return fetch(`${process.env.VUE_APP_BACKEND_URL}/api/preferences`, {
    method: 'PUT',
    headers: {'X-CSRF-Token': getCsrfToken()},
    credentials: 'include',
    body: JSON.stringify(preferences),
  });
}
