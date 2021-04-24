<template>
  <div class="footer container">
    <hr />
    <div class="links">
      <router-link to="/privacy-policy">Privacy Policy</router-link>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Footer',
};
</script>

<style scoped>
.links {
  display: flex;
  align-content: center;
  justify-content: center;
}

.footer {
  margin-top: 2rem;
}
</style>
