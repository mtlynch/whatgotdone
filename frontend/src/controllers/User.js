import {
  getCsrfToken,
  processBlankResponse,
  processJsonResponse,
} from '@/controllers/Common.js';

export function getUserSelfMetadata() {
  return fetch(`${process.env.VUE_APP_BACKEND_URL}/api/user/me`, {
    credentials: 'include',
  }).then(processJsonResponse);
}

export function getUserMetadata(username) {
  return fetch(`${process.env.VUE_APP_BACKEND_URL}/api/user/${username}`).then(
    (response) => {
      if (response.ok) {
        return response.json();
      } else if (response.status === 404) {
        // A 404 for a user profile is equivalent to retrieving an empty profile.
        return Promise.resolve({});
      }
      return response.text().then((error) => {
        return Promise.reject(error);
      });
    }
  );
}

export function setUserMetadata(metadata) {
  return fetch(`${process.env.VUE_APP_BACKEND_URL}/api/user`, {
    method: 'POST',
    headers: {'X-CSRF-Token': getCsrfToken()},
    credentials: 'include',
    body: JSON.stringify(metadata),
  }).then(processBlankResponse);
}
