<template>
  <p class="journal-header text-center">
    <Username :username="entryAuthor" />'s update for the week ending on
    <b>{{ entryDate | moment("utc", "dddd, ll") }}</b>
  </p>
</template>

<script>
import Username from "./Username.vue";
export default {
  name: "JournalHeader",
  props: {
    entryAuthor: String,
    entryDate: Date
  },
  components: {
    Username
  }
};
</script>

<style scoped>
* {
  text-align: center;
}
</style>