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
  mounted() {
    loadUserKit(
      process.env.VUE_APP_USERKIT_APP_ID,
      (userKit, userKitWidget) => {
        if (userKit.isLoggedIn() === true) {
          this.$router.back();
        } else {
          userKitWidget.open('login');
        }
      },
      () => {
        updateLoginState(/*attempts=*/ 5, () => {
          this.$router.back();
        });
      }
    );
  },
};
</script>
