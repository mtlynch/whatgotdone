<template>
  <div class="userkit">
    <p>Please wait, logging in...</p>
  </div>
</template>

<script>
import { thisFriday } from "../controllers/EntryDates.js";
import updateLoginState from "../controllers/LoginState.js";
import loadUserKit from "../controllers/UserKit.js";

export default {
  name: "Login",
  mounted() {
    loadUserKit(
      process.env.VUE_APP_USERKIT_APP_ID,
      (userKit, userKitWidget) => {
        if (userKit.isLoggedIn() === true) {
          this.$router.replace("/entry/edit/" + thisFriday());
        } else {
          userKitWidget.open("login");
        }
      },
      // eslint-disable-next-line no-unused-vars
      (userKit, userKitWidget) => {
        updateLoginState(/*attempts=*/ 5, () => {
          this.$router.replace("/entry/edit/" + thisFriday());
        });
      }
    );
  }
};
</script>