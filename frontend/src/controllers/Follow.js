import axios from 'axios';

import getCsrfToken from '@/controllers/CsrfToken.js';

export function getFollowing(username) {
  return new Promise(function (resolve, reject) {
    const url = `${process.env.VUE_APP_BACKEND_URL}/api/user/${username}/following`;
    axios
      .get(url)
      .then((result) => {
        resolve(result.data.following);
      })
      .catch((err) => reject(err));
  });
}

export function follow(username) {
  return new Promise(function (resolve, reject) {
    const url = `${process.env.VUE_APP_BACKEND_URL}/api/follow/${username}`;
    axios
      .put(
        url,
        {},
        {
          withCredentials: true,
          headers: {'X-CSRF-Token': getCsrfToken()},
        }
      )
      .then(() => {
        resolve();
      })
      .catch((err) => reject(err));
  });
}

export function unfollow(username) {
  return new Promise(function (resolve, reject) {
    const url = `${process.env.VUE_APP_BACKEND_URL}/api/follow/${username}`;
    axios
      .delete(url, {
        withCredentials: true,
        headers: {'X-CSRF-Token': getCsrfToken()},
      })
      .then(() => {
        resolve();
      })
      .catch((err) => reject(err));
  });
}
