<template>
  <span class="username">
    <router-link :to="'/' + username">{{ username }}</router-link>
  </span>
</template>

<script>
export default {
  name: 'UsernameLink',
  props: {
    username: String,
  },
};
</script>

<style scoped>
.username a {
  font-weight: bold;
  color: rgb(255, 208, 56) !important;
  text-decoration: none;
}
</style>
