import axios from 'axios';

import getCsrfToken from '@/controllers/CsrfToken.js';

const url = `${process.env.VUE_APP_BACKEND_URL}/api/preferences`;

export function getPreferences() {
  axios
    .get(url, {
      withCredentials: true,
    })
    .then((result) => {
      return result.data;
    });
}

export function savePreferences(preferences) {
  axios.post(url, preferences, {
    withCredentials: true,
    headers: {'X-CSRF-Token': getCsrfToken()},
  });
}
