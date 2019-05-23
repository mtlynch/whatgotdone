<template>
  <div class="footer container">
    <hr>
    <p class="float-center">
      <a href="mailto:contact@whatgotdone.com">Email Us</a>
      <span class="separator">|</span>
      <router-link to="/privacy-policy">Privacy Policy</router-link>
    </p>
  </div>
</template>

<script>
export default {
  name: "Footer"
};
</script>

<style scoped>
* {
  text-align: center;
}

.footer {
  margin-top: 80px;
}

.separator {
  margin: 0px 15px;
}
</style>
