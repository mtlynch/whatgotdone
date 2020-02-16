<template>
  <div class="userkit">
    <p>Please wait, logging in...</p>
  </div>
</template>

<script>
import updateLoginState from '../controllers/LoginState.js';
import loadUserKit from '../controllers/UserKit.js';

export default {
  name: 'Login',
  data() {
    return {
      previousRoute: null,
    };
  },
  methods: {
    goBackOrGoHome: function() {
      if (this.previousRoute && this.previousRoute.path) {
        this.$router.replace(this.previousRoute.path);
      } else {
        this.$router.replace('/');
      }
    },
  },
  beforeRouteEnter(to, from, next) {
    next(vm => {
      vm.previousRoute = from;
    });
  },
  mounted() {
    loadUserKit(
      process.env.VUE_APP_USERKIT_APP_ID,
      (userKit, userKitWidget) => {
        if (userKit.isLoggedIn() === true) {
          this.goBackOrGoHome();
        } else {
          userKitWidget.open('login');
        }
      },
      () => {
        updateLoginState(/*attempts=*/ 5, () => {
          this.goBackOrGoHome();
        });
      }
    );
  },
};
</script>
