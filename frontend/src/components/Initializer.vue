<template>
  <div>
    <!-- This is a dummy component. We just want the created() event. -->
  </div>
</template>

<script>
import store from '@/store.js';

import {getRecent} from '@/controllers/Recent.js';
import updateLoginState from '@/controllers/LoginState.js';
import {loadUserKit} from '@/controllers/UserKit.js';

export default {
  name: 'Initializer',
  created() {
    loadUserKit(process.env.VUE_APP_USERKIT_APP_ID).then(userKit => {
      if (userKit.isLoggedIn() === true) {
        updateLoginState();
      }
    });
    getRecent(/*start=*/ 0).then(recentEntries => {
      store.commit('setRecent', recentEntries);
    });
  },
};
</script>
