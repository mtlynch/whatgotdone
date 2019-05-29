<template>
  <div>
    <!-- This is a dummy component. We just want the created() event. -->
  </div>
</template>

<script>
export default {
  name: "Initializer",
  methods: {
    checkLoginState(attempts) {
      if (attempts <= 0) {
        return;
      }
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/user/me`;
      this.$http
        .get(url, { withCredentials: true })
        .then(result => {
          this.$store.commit("setUsername", result.data.username);
        })
        .catch(error => {
          // If checking user information fails, the cached authentication information
          // is no longer correct, so we need to clear it.
          if (error.response && error.response.status === 403) {
            this.clearCachedAuthInformation();
            return;
          }
          this.checkLoginState(attempts - 1);
        });
    },
    clearCachedAuthInformation() {
      this.$store.commit("clearUsername");
      this.deleteCookie("userkit_auth_token");
    },
    deleteCookie(name) {
      document.cookie = name + "=;expires=Thu, 01 Jan 1970 00:00:01 GMT;";
    }
  },
  created() {
    this.checkLoginState(5);
  }
};
</script>