import store from '@/store.js';

import {getUserSelfMetadata} from '@/controllers/User.js';
import {logoutUserKit} from '@/controllers/UserKit.js';

function clearCachedAuthInformation() {
  store.commit('clearUsername');
  logoutUserKit();
}

export default function updateLoginState(attempts) {
  return new Promise(function(resolve, reject) {
    if (attempts <= 0) {
      return;
    }
    getUserSelfMetadata()
      .then(metadata => {
<<<<<<< HEAD
        store.commit('setUsername', metadata.username);
        resolve();
=======
        resolve(metadata.username);
>>>>>>> Adjusting to resolve the full metadata
      })
      .catch(error => {
        // If checking user information fails, the cached authentication information
        // is no longer correct, so we need to clear it.
        if (error.response && error.response.status === 403) {
          clearCachedAuthInformation();
          return;
        }
        updateLoginState(attempts - 1).catch(error => {
          if (attempts - 1 <= 0) {
            reject(error);
          }
        });
      });
  });
}
