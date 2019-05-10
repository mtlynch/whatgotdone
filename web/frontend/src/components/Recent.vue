<template>
  <div>
    <p>Here are some recent What Got Done entries:</p>
    <p v-if="backendError" class="error">Failed to connect to backend: {{ backendError }}</p>

    <PartialJournal v-bind:key="item.key" v-bind:entry="item" v-for="item in recentEntries"/>
  </div>
</template>

<script>
import Vue from "vue";
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
            date: new Date(entry.date),
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
