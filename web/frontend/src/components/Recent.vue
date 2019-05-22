<template>
  <div>
    <h1>Recent Entries</h1>
    <p>Check out what other users got done this week:</p>
    <p v-if="backendError" class="error">Failed to connect to backend: {{ backendError }}</p>

    <PartialJournal v-bind:key="item.key" v-bind:entry="item" v-for="item in recentEntries"/>
  </div>
</template>

<script>
import PartialJournal from "./PartialJournal.vue";

export default {
  name: "Recent",
  components: {
    PartialJournal
  },
  data() {
    return {
      recentEntries: [],
      backendError: null
    };
  },
  created() {
    const url = `${process.env.VUE_APP_BACKEND_URL}/api/recentEntries`;
    this.$http
      .get(url)
      .then(result => {
        this.recentEntries = [];
        for (const entry of result.data) {
          const formattedDate = new Date(entry.date).toISOString().slice(0, 10);
          this.recentEntries.push({
            key: `/${entry.author}/${formattedDate}`,
            author: entry.author,
            date: entry.date,
            markdown: entry.markdown
          });
        }
      })
      .catch(error => {
        this.backendError = error;
      });
  }
};
</script>

<style scoped>
h1 {
  text-align: left;
}

p {
  text-align: left;
}
</style>