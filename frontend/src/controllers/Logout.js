import axios from 'axios';

import store from '@/store.js';

import getCsrfToken from '@/controllers/CsrfToken.js';
import {logoutUserKit} from '@/controllers/UserKit.js';

function deleteCookie(name) {
  document.cookie = name + '=;expires=Thu, 01 Jan 1970 00:00:01 GMT';
}

export function logout() {
  store.commit('clearLoginState');
  return new Promise(function (resolve, reject) {
    const url = `${process.env.VUE_APP_BACKEND_URL}/api/logout`;
    axios
      .post(url, {}, {headers: {'X-CSRF-Token': getCsrfToken()}})
      .then(() => {
        logoutUserKit();
        resolve();
      })
      .catch((error) => {
        reject(error);
      })
      .finally(() => {
        // Logout can fail if CSRF goes out of state. In this case, still
        // delete the CSRF cookie.
        deleteCookie('csrf_base_v3');
      });
  });
}
