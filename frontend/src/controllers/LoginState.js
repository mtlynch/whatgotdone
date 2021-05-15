import store from '@/store.js';

import {getUserSelfMetadata} from '@/controllers/User.js';
import {logoutUserKit} from '@/controllers/UserKit.js';

function clearCachedAuthInformation() {
  store.commit('clearLoginState');
  logoutUserKit();
}

export default function updateLoginState() {
  return new Promise(function (resolve, reject) {
    getUserSelfMetadata()
      .then((metadata) => {
        store.commit('setUsername', metadata.username);
        resolve(metadata);
      })
      .catch((error) => {
        // If checking user information fails, the cached authentication information
        // is no longer correct, so we need to clear it.
        if (error.response && error.response.status === 403) {
          clearCachedAuthInformation();
        }
        reject(error);
      });
  });
}
