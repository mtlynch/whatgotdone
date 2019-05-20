<template>
  <div class="view-entry container">
    <template v-if="journalEntries.length > 0">
      <div v-if="journalEntries" class="overflow-auto">
        <b-pagination-nav
          :link-gen="linkGen"
          :page-gen="pageGen"
          :number-of-pages="links.length"
          v-if="links.length > 0"
          align="center"
          use-router
        ></b-pagination-nav>
      </div>

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
    </template>
    <template v-else>
      <p>This user has not posted any What Got Done updates.</p>
    </template>
    <p v-if="backendError" class="error">Failed to connect to backend: {{ backendError }}</p>
  </div>
</template>

<script>
import Vue from "vue";
import Journal from "./Journal.vue";
import Pagination from "bootstrap-vue/es/components/pagination";

Vue.use(Pagination);

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
  methods: {
    linkGen(pageNum) {
      return this.links[pageNum - 1];
    },
    pageGen(pageNum) {
      const dateRaw = this.links[pageNum - 1].split("/")[2];
      const dateUtc = new Date(`${dateRaw}T00:00:00Z`);
      const userTimezoneOffset = dateUtc.getTimezoneOffset() * 60000;
      const date = new Date(dateUtc.getTime() + userTimezoneOffset);
      const friendlyMonth = date.toLocaleString("en-us", {
        month: "short"
      });
      return `${friendlyMonth}. ${date.getDate()}`;
    }
  },
  computed: {
    links: function() {
      let links = [];
      for (const entry of this.journalEntries) {
        links.push(`/${this.$route.params.username}/${entry.key}`);
      }
      return links;
    },
    currentEntry: function() {
      if (!this.$route.params.date) {
        return null;
      }
      for (const entry of this.journalEntries) {
        if (this.$route.params.date === entry.key) {
          return entry;
        }
      }
      return null;
    }
  },
  created() {
    const url = `${process.env.VUE_APP_BACKEND_URL}/api/entries/${
      this.$route.params.username
    }`;
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
        if (this.journalEntries.length == 0) {
          return;
        }
        this.journalEntries.sort((a, b) => a.date - b.date);

        if (!this.$route.params.date) {
          const lastEntry = this.journalEntries[this.journalEntries.length - 1];
          this.$router.push(`/${this.$route.params.username}/${lastEntry.key}`);
          return;
        }
      })
      .catch(error => {
        this.backendError = error;
      });
  }
};
</script>
