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
    this.$http.post(url).then(result => {
      window.location.href = "/login";
    });
  }
};
</script>
