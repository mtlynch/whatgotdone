<template>
  <p>
    <Username :username="username"/>'s update for the week ending on
    <b>{{ date | moment("dddd, ll") }}</b>
  </p>
</template>

<script>
import Username from "./Username.vue";
export default {
  name: "JournalHeader",
  props: {
    username: String,
    date: String
  },
  components: {
    Username
  }
};
</script>

<style scoped>
</style>