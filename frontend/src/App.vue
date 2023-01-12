<template>
  <div>
    <NavigationBar />
    <b-container id="app">
      <router-view></router-view>
    </b-container>
    <PageFooter />
  </div>
</template>

<script>
import {getRecent} from '@/controllers/Recent.js';
import initializeUserState from '@/controllers/UserState.js';
import {loadUserKit} from '@/controllers/UserKit.js';

import PageFooter from '@/components/PageFooter';
import NavigationBar from '@/components/NavigationBar';

export default {
  name: 'app',
  components: {
    PageFooter,
    NavigationBar,
  },
  created() {
    loadUserKit(process.env.VUE_APP_USERKIT_APP_ID).then((userKit) => {
      if (userKit.isLoggedIn() === true) {
        initializeUserState().catch(() => {
          this.$store.commit('clearUserState');
        });
      } else {
        this.$store.commit('clearUserState');
      }
    });
    getRecent(/*start=*/ 0).then((recentEntries) => {
      this.$store.commit('setRecent', recentEntries);
    });
  },
};
</script>

<style>
@import '~@fortawesome/fontawesome-svg-core/styles.css';

#app {
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  margin-top: 60px;
}

#app a.btn {
  color: white;
}

/* TODO: Move these to the view entry component */
#app a.page-link {
  color: white;
}

#app .page-link {
  border: 1px solid rgb(124, 133, 145);
}
</style>
