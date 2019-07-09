<template>
  <p>
    <Username :username="entryAuthor" />'s update for the week ending on
    <b>{{ entryDate | moment("dddd, ll") }}</b>
  </p>
</template>

<script>
import Username from "./Username.vue";
export default {
  name: "JournalHeader",
  props: {
    entryAuthor: String,
    entryDate: String
  },
  components: {
    Username
  }
};
</script>

<style scoped>
</style>