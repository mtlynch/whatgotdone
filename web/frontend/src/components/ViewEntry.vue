<template>
  <div class="view-entry container">
    <template v-if="entries.length > 0">
      <b-pagination-nav
        :pages="pages"
        :number-of-pages="numberOfPages"
        v-model="currentEntryIndex"
        v-if="pages.length > 0"
        align="center"
        no-page-detect
        use-router
      ></b-pagination-nav>
      <b-form-checkbox v-model="showEmptyEntries" v-if="canEdit">Show empty entries</b-form-checkbox>

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
      journalEntriesRaw: [],
      showEmptyEntries: false,
      backendError: null
    };
  },
  methods: {
    goToLatestEntry() {
      const lastEntry = this.entries[this.entries.length - 1];
      this.$router.replace(`/${this.$route.params.username}/${lastEntry.key}`);
    },
    loadjournalEntries: function() {
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/entries/${
        this.$route.params.username
      }`;
      this.$http
        .get(url)
        .then(result => {
          for (const entry of result.data) {
            this.journalEntriesRaw = result.data;
          }

          if (!this.$route.params.date) {
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
      for (const entry of this.entries) {
        pages.push({
          link: `/${this.$route.params.username}/${entry.key}`,
          text: new moment(entry.key).format("MMM. D").replace("May.", "May")
        });
      }
      return pages;
    },
    numberOfPages: function() {
      console.log(`numberOfPages: ${this.pages.length}`);
      return this.pages.length;
    },
    username: function() {
      return this.$store.state.username;
    },
    canEdit: function() {
      return this.username && this.username === this.$route.params.username;
    },
    currentEntryIndex: {
      get: function() {
        console.log("in getCurrentEntryIndex");
        if (!this.$route.params.date) {
          console.log("route date is null");
          return null;
        }
        const entries = this.entries;
        for (let i in entries) {
          const entryDate = entries[i].key;
          console.log(`${i} -> ${entryDate}`);
          if (this.$route.params.date === entryDate) {
            let ans = +i + 1;
            console.log(`returning ${ans} from getCurrentEntryIndex`);
            return +i + 1;
          }
        }
        console.log("returning 1 from getCurrentEntryIndex");
        return 1;
      },
      set: function(newValue) {
        console.log(`Setting currentEntryIndex to ${newValue}`);
      }
    },
    currentEntry: function() {
      if (!this.$route.params.date) {
        return null;
      }
      for (const entry of this.entries) {
        if (this.$route.params.date === entry.key && entry.markdown) {
          return entry;
        }
      }
      return null;
    },
    entries: function() {
      const entries = [];
      for (const entry of this.journalEntriesRaw) {
        entries.push({
          key: entry.date,
          date: new moment(entry.date),
          lastModified: new moment(entry.lastModified),
          markdown: entry.markdown
        });
      }
      entries.push({
        key: "2019-06-07",
        date: new moment("2019-06-07"),
        lastModified: null,
        markdown: null
      });
      entries.push({
        key: "2019-06-14",
        date: new moment("2019-06-14"),
        lastModified: null,
        markdown: null
      });
      console.log(entries);
      return entries;
    }
  },
  created() {
    this.loadjournalEntries();
  },
  watch: {
    $route(to, from) {
      if (!to.params.date) {
        this.goToLatestEntry();
      }
      if (to.params.username != from.params.username) {
        this.loadjournalEntries();
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
