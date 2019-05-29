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

      <JournalHeader :username="$route.params.username" :date="$route.params.date"/>
      <Journal v-bind:entry="currentEntry" v-if="currentEntry"/>
      <p v-else>
        No journal entry found for
        <b>{{ $route.params.date }}</b>
      </p>
    </template>
    <template v-else>
      <p>
        <span class="username">{{ $route.params.username }}</span> has not posted any What Got Done updates.
      </p>
    </template>
    <p v-if="backendError" class="error">Failed to connect to backend: {{ backendError }}</p>
    <b-button
      v-if="canEdit"
      :to="'/entry/edit/' + this.$route.params.date"
      variant="primary"
      class="float-right"
    >Edit</b-button>
  </div>
</template>

<script>
import Vue from "vue";
import moment from "moment";
import Journal from "./Journal.vue";
import JournalHeader from "./JournalHeader.vue";
import Pagination from "bootstrap-vue/es/components/pagination";

Vue.use(Pagination);

export default {
  name: "ViewEntry",
  components: {
    Journal,
    JournalHeader
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
      return new moment(dateRaw).format("MMM. D");
    },
    goToLatestEntry() {
      const lastEntry = this.journalEntries[this.journalEntries.length - 1];
      this.$router.replace(`/${this.$route.params.username}/${lastEntry.key}`);
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
    username: function() {
      return this.$store.state.username;
    },
    canEdit: function() {
      return this.username && this.username === this.$route.params.username;
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
          this.goToLatestEntry();
          return;
        }
      })
      .catch(error => {
        this.backendError = error;
      });
  },
  watch: {
    $route() {
      if (!this.$route.params.date) {
        this.goToLatestEntry();
      }
    }
  }
};
</script>

<style scoped>
span.username {
  color: rgb(255, 208, 56);
  font-weight: bold;
}
</style>
