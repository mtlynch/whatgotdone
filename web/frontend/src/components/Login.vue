<template>
  <div class="userkit">
    <p>Please wait, logging in...</p>
  </div>
</template>

<script>
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
        const url = `${process.env.VUE_APP_BACKEND_URL}/api/user/me`;
        this.$http
          .get(url)
          .then(result => {
            this.$router.push(`/${result.data.username}`);
          })
          .catch(function() {
            // Do nothing, wait for request to succeed.
          });
      }, 100);
    }
  }
};
</script>
