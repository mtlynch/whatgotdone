import axios from 'axios';

import getCsrfToken from '@/controllers/CsrfToken.js';

const url = `${process.env.VUE_APP_BACKEND_URL}/api/preferences`;

export function getPreferences() {
  return new Promise(function(resolve, reject) {
    axios
      .get(url, {
        withCredentials: true,
      })
      .then(result => {
        resolve(result.data);
      })
      .catch(error => {
        reject(error);
      });
  });
}

export function savePreferences(preferences) {
  return new Promise(function(resolve, reject) {
    axios
      .post(url, preferences, {
        withCredentials: true,
        headers: {'X-CSRF-Token': getCsrfToken()},
      })
      .then(() => {
        resolve();
      })
      .catch(error => {
        reject(error);
      });
  });
}
