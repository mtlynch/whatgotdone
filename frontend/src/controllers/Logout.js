import store from '@/store.js';

import {getCsrfToken} from '@/controllers/Common.js';
import {logoutUserKit} from '@/controllers/UserKit.js';

function deleteCookie(name) {
  document.cookie = name + '=;expires=Thu, 01 Jan 1970 00:00:01 GMT';
}

export function logout() {
  store.commit('clearUserState');
  return fetch(`${process.env.VUE_APP_BACKEND_URL}/api/logout`, {
    method: 'POST',
    credentials: 'include',
    headers: {'X-CSRF-Token': getCsrfToken()},
  })
    .then(() => {
      logoutUserKit();
      return Promise.resolve();
    })
    .finally(() => {
      // Logout can fail if CSRF goes out of state. In this case, still
      // delete the CSRF cookie.
      deleteCookie('csrf_base_v4');
    });
}
