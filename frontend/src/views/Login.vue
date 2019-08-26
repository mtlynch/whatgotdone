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
    let userKitScript = document.createElement("script");
    userKitScript.setAttribute("src", "https://widget.userkit.io/widget.js");
    userKitScript.setAttribute(
      "data-app-id",
      process.env.VUE_APP_USERKIT_APP_ID
    );
    userKitScript.setAttribute("data-show-on-load", "login");
    userKitScript.setAttribute("data-login-dismiss", "false");
    document.head.appendChild(userKitScript);
  },
  created() {
    document.addEventListener("UserKitSignIn", () => {
      updateLoginState(/*attempts=*/ 5);
      this.$router.push("/entry/edit/" + thisFriday());
    });
  }
};
</script>
