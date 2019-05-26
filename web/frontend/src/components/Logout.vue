<template>
  <div>
    <p>Signing out...</p>
  </div>
</template>

<script>
export default {
  name: "Logout",
  created() {
    const url = `${process.env.VUE_APP_BACKEND_URL}/api/logout`;
    this.$http.post(url).then(() => {
      this.deleteCookie("userkit_auth_token");
      window.location.href = "/login";
    });
  },
  methods: {
    deleteCookie(name) {
      document.cookie = name + "=;expires=Thu, 01 Jan 1970 00:00:01 GMT;";
    }
  }
};
</script>
