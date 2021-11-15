import {getCsrfToken, processJsonResult} from '@/controllers/Common.js';

export function getPreferences() {
  return fetch(`${process.env.VUE_APP_BACKEND_URL}/api/preferences`, {
    credentials: 'include',
  }).then(processJsonResult);
}

export function savePreferences(preferences) {
  return fetch(`${process.env.VUE_APP_BACKEND_URL}/api/preferences`, {
    method: 'PUT',
    headers: {'X-CSRF-Token': getCsrfToken()},
    credentials: 'include',
    body: JSON.stringify(preferences),
  });
}
