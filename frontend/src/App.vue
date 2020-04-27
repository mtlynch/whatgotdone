<template>
  <div>
    <NavigationBar />
    <div id="app" class="container">
      <router-view></router-view>
    </div>
    <Footer />
  </div>
</template>

<script>
import {getRecent} from '@/controllers/Recent.js';
import updateLoginState from '@/controllers/LoginState.js';
import {loadUserKit} from '@/controllers/UserKit.js';

import Footer from '@/components/Footer';
import NavigationBar from '@/components/NavigationBar';

export default {
  name: 'app',
  components: {
    Footer,
    NavigationBar,
  },
  created() {
    loadUserKit(process.env.VUE_APP_USERKIT_APP_ID).then(userKit => {
      if (userKit.isLoggedIn() === true) {
        updateLoginState();
      } else {
        this.$store.commit('clearLoginState');
        if (this.routeRequiresLogin) {
          this.$router.push('/login');
        }
      }
    });
    getRecent(/*start=*/ 0).then(recentEntries => {
      this.$store.commit('setRecent', recentEntries);
    });
  },
  computed: {
    routeRequiresLogin: function() {
      const routeName = this.$router.currentRoute.name;
      if (!routeName) {
        return false;
      }
      if (routeName === 'Preferences') {
        return true;
      }
      if (routeName.indexOf('Edit') === 0) {
        return true;
      }
      return false;
    },
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

/* TODO: Remove the #app part of the selector */
#app a {
  color: rgb(101, 168, 255);
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
