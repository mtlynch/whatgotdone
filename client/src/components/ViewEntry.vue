<template>
  <div class="view-entry container">
    <p v-if="previousPath">
      <router-link v-bind:to="previousPath">Previous</router-link>
    </p>
    <p v-if="nextPath">
      <router-link v-bind:to="nextPath">Next</router-link>
    </p>

    <p>
      Journal entry for
      <b>{{ $route.params.username }}</b> on
      <b>{{ $route.params.date }}</b>
    </p>
    <Journal v-bind:entry="currentEntry" v-if="currentEntry"/>
    <p v-else>
      No journal entry found for
      <b>{{ $route.params.date }}</b>
    </p>

    <p v-if="previousPath">
      <router-link v-bind:to="previousPath">Previous</router-link>
    </p>
    <p v-if="nextPath">
      <router-link v-bind:to="nextPath">Next</router-link>
    </p>
    <p v-if="backendError" class="error">Failed to connect to backend: {{ backendError }}</p>
  </div>
</template>

<script>
import Journal from "./Journal.vue";

export default {
  name: "ViewEntry",
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
    currentEntryIndex: function() {
      if (!this.journalEntries || !this.$route.params.date) {
        return null;
      }
      for (const [index, entry] of this.journalEntries.entries()) {
        if (this.$route.params.date === entry.key) {
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
      return `/${this.$route.params.username}/${this.nextEntryKey}`;
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
      return `/${this.$route.params.username}/${this.previousEntryKey}`;
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
            lastModified: new Date(entry.lastModified),
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
.view-entry {
  color: #2c3e50;
}
</style>
