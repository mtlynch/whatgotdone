<template>
  <div class="view-entry container">
    <template v-if="journalEntries.length > 0">
      <b-pagination-nav :pages="pages" v-if="pages.length > 0" align="center" use-router></b-pagination-nav>

      <JournalHeader :entryAuthor="entryAuthor" :entryDate="entryDate" />
      <Journal v-bind:entry="currentEntry" v-if="currentEntry" />
      <p v-else>
        No journal entry found for
        <b>{{ entryDate }}</b>
      </p>
    </template>
    <template v-else>
      <p>
        <span class="username">{{ entryAuthor }}</span> has not posted any What Got Done updates.
      </p>
    </template>
    <p v-if="backendError" class="error">Failed to connect to backend: {{ backendError }}</p>
    <b-button
      v-if="canEdit"
      :to="'/entry/edit/' + this.entryDate"
      variant="primary"
      class="float-right edit-btn"
    >Edit</b-button>
    <Reactions :entryAuthor="entryAuthor" :entryDate="entryDate" />
  </div>
</template>

<script>
import Vue from "vue";
import moment from "moment";
import Journal from "../components/Journal.vue";
import JournalHeader from "../components/JournalHeader.vue";
import Reactions from "../components/Reactions.vue";
import Pagination from "bootstrap-vue/es/components/pagination";

Vue.use(Pagination);

export default {
  name: "ViewEntry",
  components: {
    Journal,
    JournalHeader,
    Reactions
  },
  data() {
    return {
      journalEntries: [],
      backendError: null
    };
  },
  methods: {
    goToLatestEntry() {
      // I don't understand how this can happen, but sometimes I'm seeing the
      // e2e test try to redirect the client to /undefined/[somedate] and it
      // seems to be caused when goToLatestEntry is called when
      // this.entryAuthor is undefined, even though that seems like it should
      // never happen.
      if (!this.entryAuthor) {
        return;
      }
      const lastEntry = this.journalEntries[this.journalEntries.length - 1];
      this.$router.replace(`/${this.entryAuthor}/${lastEntry.key}`);
    },
    loadJournalEntries: function() {
      this.journalEntries = [];
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/entries/${this.entryAuthor}`;
      this.$http
        .get(url)
        .then(result => {
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

          if (!this.entryDate) {
            this.goToLatestEntry();
            return;
          }
        })
        .catch(error => {
          this.backendError = error;
        });
    }
  },
  computed: {
    pages: function() {
      let pages = [];
      for (const entry of this.journalEntries) {
        pages.push({
          link: `/${this.entryAuthor}/${entry.key}`,
          text: new moment(entry.key).format("MMM. D").replace("May.", "May")
        });
      }
      return pages;
    },
    loggedInUsername: function() {
      return this.$store.state.username;
    },
    canEdit: function() {
      return (
        this.loggedInUsername && this.loggedInUsername === this.entryAuthor
      );
    },
    entryAuthor: function() {
      return this.$route.params.username;
    },
    entryDate: function() {
      return this.$route.params.date;
    },
    currentEntry: function() {
      if (!this.entryDate) {
        return null;
      }
      for (const entry of this.journalEntries) {
        if (this.entryDate === entry.key) {
          return entry;
        }
      }
      return null;
    }
  },
  created() {
    this.loadJournalEntries();
  },
  watch: {
    $route(to, from) {
      if (to.params.username != from.params.username) {
        this.loadJournalEntries();
      }
      if (!to.params.date) {
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

.edit-btn {
  margin: 25px 0px;
}
</style>
