import {getCsrfToken, processJsonResponse} from '@/controllers/Common.js';

export function getFollowing(username) {
  return fetch(
    `${process.env.VUE_APP_BACKEND_URL}/api/user/${username}/following`
  )
    .then(processJsonResponse)
    .then((followData) => {
      return Promise.resolve(followData.following);
    });
}

export function follow(username) {
  return fetch(`${process.env.VUE_APP_BACKEND_URL}/api/follow/${username}`, {
    method: 'PUT',
    credentials: 'include',
    headers: {'X-CSRF-Token': getCsrfToken()},
  });
}

export function unfollow(username) {
  return fetch(`${process.env.VUE_APP_BACKEND_URL}/api/follow/${username}`, {
    method: 'DELETE',
    credentials: 'include',
    headers: {'X-CSRF-Token': getCsrfToken()},
  });
}
