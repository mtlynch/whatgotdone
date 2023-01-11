import store from '@/store.js';

import {getFollowing} from '@/controllers/Follow.js';
import {getUserSelfMetadata} from '@/controllers/User.js';
import {logoutUserKit} from '@/controllers/UserKit.js';

function clearCachedAuthInformation() {
  store.commit('clearUserState');
  logoutUserKit();
}

function updateFollowingUsers(username) {
  getFollowing(username).then((following) => {
    for (const followedUser of following) {
      store.commit('addFollowedUser', followedUser);
    }
  });
}

export default function initializeUserState() {
  return new Promise(function (resolve, reject) {
    getUserSelfMetadata()
      .then((metadata) => {
        store.commit('setUsername', metadata.username);
        updateFollowingUsers(metadata.username);
        resolve();
      })
      .catch((error) => {
        clearCachedAuthInformation();
        reject(error);
      });
  });
}
