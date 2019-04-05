<template>
  <div id="app" class="container">
    <p>
      Journal entry for
      <b>{{ username }}</b> on
      <b>{{ selectedDate }}</b>
    </p>
    <Journal v-bind:entry="currentEntry" v-if="currentEntry"/>
    <p v-else>
      No journal entry found for
      <b>{{ selectedDate }}</b>
    </p>

    <p v-if="previousPath">
      <a v-bind:href="previousPath">Previous</a>
    </p>
    <p v-if="nextPath">
      <a v-bind:href="nextPath">Next</a>
    </p>
    <p v-if="backendError" class="error">Failed to connect to backend: {{ backendError }}</p>
  </div>
</template>

<script>
import Journal from "./components/Journal.vue";

export default {
  name: "app",
  components: {
    Journal
  },
  data() {
    return {
      journalEntries: [],
      backendError: null
    };
  },
  computed: {
    username: function() {
      if (!window.location.pathname) {
        return null;
      }
      return window.location.pathname.split("/")[1];
    },
    selectedDate: function() {
      if (!window.location.pathname) {
        return null;
      }
      return window.location.pathname.split("/")[2];
    },
    currentEntryIndex: function() {
      if (!this.journalEntries || !this.selectedDate) {
        return null;
      }
      for (const [index, entry] of this.journalEntries.entries()) {
        if (this.selectedDate === entry.key) {
          return index;
        }
      }
      return null;
    },
    currentEntry: function() {
      if (this.currentEntryIndex === null) {
        return null;
      }
      return this.journalEntries[this.currentEntryIndex];
    },
    nextEntryKey: function() {
      if (
        this.currentEntryIndex === null ||
        this.currentEntryIndex === this.journalEntries.length - 1
      ) {
        return null;
      }
      return this.journalEntries[this.currentEntryIndex + 1].key;
    },
    nextPath: function() {
      if (!this.nextEntryKey) {
        return null;
      }
      return `/${this.username}/${this.nextEntryKey}`;
    },
    previousEntryKey: function() {
      if (this.currentEntryIndex === null || this.currentEntryIndex === 0) {
        return null;
      }
      return this.journalEntries[this.currentEntryIndex - 1].key;
    },
    previousPath: function() {
      if (!this.previousEntryKey) {
        return null;
      }
      return `/${this.username}/${this.previousEntryKey}`;
    }
  },
  created() {
    const url = `${process.env.VUE_APP_BACKEND_URL}/entries`;
    this.$http
      .get(url)
      .then(result => {
        this.journalEntries = [];
        for (const entry of result.data) {
          this.journalEntries.push({
            key: entry.date,
            date: new Date(entry.date),
            markdown: entry.markdown
          });
        }
        this.journalEntries.sort((a, b) => a.date - b.date);
      })
      .catch(function(error) {
        this.backendError = error;
      });
  }
};
</script>

<style>
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
