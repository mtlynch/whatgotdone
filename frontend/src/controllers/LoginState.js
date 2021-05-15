import store from '@/store.js';

import {getUserSelfMetadata} from '@/controllers/User.js';
import {logoutUserKit} from '@/controllers/UserKit.js';

function clearCachedAuthInformation() {
  store.commit('clearLoginState');
  logoutUserKit();
}

export default function updateLoginState() {
  getUserSelfMetadata()
    .then((metadata) => {
      store.commit('setUsername', metadata.username);
      return metadata;
    })
    .catch((error) => {
      // If checking user information fails, the cached authentication information
      // is no longer correct, so we need to clear it.
      if (error.response && error.response.status === 403) {
        clearCachedAuthInformation();
      }
      throw error;
    });
}
