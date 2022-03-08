<template>
  <div class="footer">
    <hr />
    <div class="links">
      <router-link to="/privacy-policy">Privacy Policy</router-link>
    </div>
  </div>
</template>

<style scoped>
.footer {
  margin: 2rem auto;
  max-width: 95%;
}

.links {
  display: flex;
  align-content: center;
  justify-content: center;
}
</style>

<script>
export default {
  name: 'PageFooter',
};
</script>
