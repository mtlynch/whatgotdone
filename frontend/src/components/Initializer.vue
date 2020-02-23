<template>
  <div>
    <!-- This is a dummy component. We just want the created() event. -->
  </div>
</template>

<script>
import store from '@/store.js';

import {getRecent} from '@/controllers/Recent.js';
import updateLoginState from '@/controllers/LoginState.js';

export default {
  name: 'Initializer',
  methods: {
    refreshFollowing: function() {},
  },
  created() {
    updateLoginState(5);
    getRecent(/*start=*/ 0).then(recentEntries => {
      store.commit('setRecent', recentEntries);
    });
  },
};
</script>
