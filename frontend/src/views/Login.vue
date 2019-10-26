<template>
  <div class="userkit">
    <p>Please wait, logging in...</p>
  </div>
</template>

<script>
import { thisFriday } from "../controllers/EntryDates.js";
import updateLoginState from "../controllers/LoginState.js";

export default {
  name: "Login",
  mounted() {
    if(this.$store.runOnce === true) {
      // ESLint can't detect that these global variables are defined in a
      // dynamically-loaded script, so make local copies.
      const userKit = UserKit; // eslint-disable-line no-undef
      const userKitWidget = UserKitWidget; // eslint-disable-line no-undef

      if(userKit.isLoggedIn() === true) {
        this.$router.replace("/entry/edit/" + thisFriday());
      } else {
        userKitWidget.open("login");
      }
      return;
    }

    let userKitScript = document.createElement("script");
    userKitScript.setAttribute("src", "https://widget.userkit.io/widget.js");
    userKitScript.setAttribute(
      "data-app-id",
      process.env.VUE_APP_USERKIT_APP_ID
    );
    userKitScript.setAttribute("data-show-on-load", "login");
    userKitScript.setAttribute("data-login-dismiss", "false");
    document.head.appendChild(userKitScript);

    document.addEventListener("UserKitSignIn", () => {
      updateLoginState(/*attempts=*/ 5, () => {
        this.$router.replace("/entry/edit/" + thisFriday());
      });
    });

    this.$store.runOnce = true;
  },
};
</script>