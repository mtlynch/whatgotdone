<template>
  <div class="view-entry container">
    <template v-if="journalEntries.length > 0">
      <b-pagination-nav
        :pages="pages"
        :number-of-pages="totalPages"
        v-if="pages.length > 0"
        v-model="currentEntryIndex"
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
      journalEntries: [],
      showEmptyEntries: false,
      backendError: null
    };
  },
  methods: {
    goToLatestEntry() {
      const lastEntry = this.journalEntries[this.journalEntries.length - 1];
      this.$router.replace(`/${this.$route.params.username}/${lastEntry.key}`);
    },
    thisFriday() {
      const today = moment().isoWeekday();
      const friday = 5;

      if (today <= friday) {
        return moment().isoWeekday(friday);
      } else {
        return moment()
          .add(1, "weeks")
          .isoWeekday(friday);
      }
    },
    entryDates: function() {
      const dates = [];
      if (this.showEmptyEntries) {
        let d = moment(this.journalEntries[0].key);
        while (d <= this.thisFriday()) {
          dates.push(d.format("YYYY-MM-DD"));
          d = d.add(1, "weeks");
        }
      } else {
        for (const entry of this.journalEntries) {
          dates.push(entry.key);
        }
      }
      return dates;
    },
    getCurrentEntryIndex: function() {
      console.log("in getCurrentEntryIndex");
      if (!this.$route.params.date) {
        console.log("route date is null");
        return null;
      }
      const dates = this.entryDates();
      for (let i in dates) {
        const entryDate = dates[i];
        console.log(`${i} -> ${entryDate}`);
        if (this.$route.params.date === entryDate) {
          let ans = +i + 1;
          console.log(`returning ${ans} from getCurrentEntryIndex`);
          return +i + 1;
        }
      }
      console.log("returning null from getCurrentEntryIndex");
      return null;
    },
    loadJournalEntries: function() {
      this.journalEntries = [];
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/entries/${
        this.$route.params.username
      }`;
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
    totalPages: function() {
      let total = this.entryDates().length;
      console.log(`totalPages: ${total}`);
      return this.entryDates().length;
    },
    pages: function() {
      const pages = [];
      for (const d of this.entryDates()) {
        pages.push({
          link: {
            path: `/${this.$route.params.username}/${d}`
          },
          text: new moment(d).format("MMM. D").replace("May.", "May")
        });
      }
      console.log("generating pages");
      console.log(pages);
      return pages;
    },
    username: function() {
      return this.$store.state.username;
    },
    canEdit: function() {
      return this.username && this.username === this.$route.params.username;
    },
    currentEntryIndex: {
      get: function() {
        return this.getCurrentEntryIndex();
      },
      set: function(newValue) {
        console.log(`Setting currentEntryIndex to ${newValue}`);
        return false;
      }
    },
    currentEntry: function() {
      const index = this.getCurrentEntryIndex();
      console.log(`currentEntry: ${index}`);
      if (index === null) {
        return null;
      }
      return this.journalEntries[index];
    }
  },
  created() {
    this.loadJournalEntries();
  },
  watch: {
    $route(to, from) {
      if (!to.params.date) {
        this.goToLatestEntry();
      }
      if (to.params.username != from.params.username) {
        this.loadJournalEntries();
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
