<template>
  <div class="userkit">
    <p>Please wait, logging in...</p>
  </div>
</template>

<script>
import initializeUserState from '@/controllers/UserState.js';
import {loadUserKit} from '@/controllers/UserKit.js';

export default {
  name: 'Login',
  data() {
    return {
      previousRoute: null,
    };
  },
  methods: {
    goBackOrGoHome: function () {
      if (this.previousRoute) {
        this.$router.replace(this.previousRoute);
      } else {
        this.$router.replace('/');
      }
    },
  },
  beforeRouteEnter(to, from, next) {
    next((vm) => {
      if (from.path) {
        vm.previousRoute = from.path;
      }
    });
  },
  mounted() {
    loadUserKit(process.env.VUE_APP_USERKIT_APP_ID).then((userKit) => {
      userKit.authenticate().then(() => {
        initializeUserState().then(() => {
          this.goBackOrGoHome();
        });
      });
    });
  },
};
</script>
