<template>
  <div class="userkit">
    <p>Please wait, logging in...</p>
  </div>
</template>

<script>
import updateLoginState from "../controllers/LoginState.js";

export default {
  name: "Login",
  data() {
    return {
      polling: null
    };
  },
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
  beforeDestroy() {
    clearInterval(this.polling);
  },
  created() {
    this.pollLoginStatus();
  },
  methods: {
    pollLoginStatus() {
      this.polling = setInterval(() => {
        if (this.isLoggedIn()) {
          updateLoginState(5);
          this.$router.push("/submit");
        }
      }, 100);
    },
    isLoggedIn() {
      // eslint-disable-next-line
      return typeof UserKit !== "undefined" && UserKit.isLoggedIn();
    }
  }
};
</script>
