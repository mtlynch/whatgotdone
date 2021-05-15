import axios from 'axios';

import getCsrfToken from '@/controllers/CsrfToken.js';

export function getUserSelfMetadata() {
  const url = `${process.env.VUE_APP_BACKEND_URL}/api/user/me`;
  axios.get(url, {withCredentials: true}).then((result) => {
    return result.data;
  });
}

export function getUserMetadata(username) {
  const url = `${process.env.VUE_APP_BACKEND_URL}/api/user/${username}`;
  axios
    .get(url)
    .then((result) => {
      return result.data;
    })
    .catch((error) => {
      if (error.response && error.response.status == 404) {
        // A 404 for a user profile is equivalent to retrieving an empty profile.
        return;
      } else {
        throw error;
      }
    });
}

export function setUserMetadata(metadata) {
  const url = `${process.env.VUE_APP_BACKEND_URL}/api/user`;
  axios.post(url, metadata, {
    withCredentials: true,
    headers: {'X-CSRF-Token': getCsrfToken()},
  });
}
