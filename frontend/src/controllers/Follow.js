import axios from 'axios';

import getCsrfToken from '@/controllers/CsrfToken.js';

export function getFollowing(username) {
  const url = `${process.env.VUE_APP_BACKEND_URL}/api/user/${username}/following`;
  axios.get(url).then((result) => {
    return result.data.following;
  });
}

export function follow(username) {
  const url = `${process.env.VUE_APP_BACKEND_URL}/api/follow/${username}`;
  axios.put(
    url,
    {},
    {
      withCredentials: true,
      headers: {'X-CSRF-Token': getCsrfToken()},
    }
  );
}

export function unfollow(username) {
  const url = `${process.env.VUE_APP_BACKEND_URL}/api/follow/${username}`;
  axios.delete(url, {
    withCredentials: true,
    headers: {'X-CSRF-Token': getCsrfToken()},
  });
}
