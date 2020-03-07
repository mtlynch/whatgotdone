<template>
  <div>
    <!-- This is a dummy component. We just want the created() event. -->
  </div>
</template>

<script>
import store from '@/store.js';

import {getFollowing} from '@/controllers/Follow.js';
import {getRecent} from '@/controllers/Recent.js';
import updateLoginState from '@/controllers/LoginState.js';

export default {
  name: 'Initializer',
  created() {
    updateLoginState(5);
    getFollowing().then(following => {
      store.commit('setFollowing', following);
    });
    getRecent(/*start=*/ 0).then(recentEntries => {
      store.commit('setRecent', recentEntries);
    });
  },
};
</script>
