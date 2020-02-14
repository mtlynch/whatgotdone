<template>
  <div class="view-entry container">
    <template v-if="journalEntries.length > 0">
      <b-pagination-nav
        :pages="pages"
        v-if="pages.length > 0"
        align="center"
        use-router
      ></b-pagination-nav>
      <b-form-checkbox
        v-model="showEmptyEntries"
        v-if="canEdit"
        class="show-empty"
        >Show empty entries</b-form-checkbox
      >

      <Journal v-bind:entry="currentEntry" v-if="currentEntry" />
      <p class="missing-entry" v-else>
        <Username :username="entryAuthor" />&nbsp;has not posted a journal entry
        for
        <b>{{ entryDate | moment('utc', 'dddd, ll') }}</b>
      </p>
    </template>
    <template v-else>
      <p>
        <span class="username">{{ entryAuthor }}</span> has not posted any What
        Got Done updates.
      </p>
    </template>
    <p v-if="backendError" class="error">
      Failed to connect to backend: {{ backendError }}
    </p>
    <div class="author-controls mb-4" v-if="canEdit">
      <b-button
        :href="twitterShareUrl"
        title="Share on Twitter"
        class="twitter"
        variant="info"
        ><font-awesome-icon :icon="['fab', 'twitter']" class="mr-3" /> Share on
        Twitter</b-button
      >
      <b-button
        :to="'/entry/edit/' + this.entryDate"
        variant="primary"
        class="float-right edit-btn"
        >Edit</b-button
      >
    </div>
    <Reactions
      :entryAuthor="entryAuthor"
      :entryDate="entryDate"
      v-if="currentEntry"
    />
  </div>
</template>

<script>
import moment from 'moment';
import Journal from '../components/Journal.vue';
import Reactions from '../components/Reactions.vue';
import Username from '../components/Username.vue';
import {thisFriday} from '../controllers/EntryDates.js';

export default {
  name: 'ViewEntry',
  components: {
    Journal,
    Reactions,
    Username,
  },
  data() {
    return {
      journalEntries: [],
      showEmptyEntries: false,
      backendError: null,
    };
  },
  methods: {
    loadJournalEntries: function() {
      this.journalEntries = [];
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/entries/${this.entryAuthor}`;
      this.$http
        .get(url)
        .then(result => {
          for (const entry of result.data) {
            this.journalEntries.push({
              key: entry.date,
              author: this.entryAuthor,
              date: new Date(entry.date),
              lastModified: new Date(entry.lastModified),
              markdown: entry.markdown,
            });
          }
          if (this.journalEntries.length == 0) {
            return;
          }
          this.journalEntries.sort((a, b) => a.date - b.date);
        })
        .catch(error => {
          this.backendError = error;
        });
    },
  },
  computed: {
    pages: function() {
      let pages = [];
      for (const date of this.entryDates) {
        pages.push({
          link: `/${this.entryAuthor}/${date}`,
          text: new moment(date).format('MMM. D').replace('May.', 'May'),
        });
      }
      return pages;
    },
    entryDates: function() {
      let dates = [];
      if (this.showEmptyEntries) {
        let d = moment(this.journalEntries[0].key);
        while (d <= moment(thisFriday())) {
          dates.push(d.format('YYYY-MM-DD'));
          d = d.add(1, 'weeks');
        }
      } else {
        for (const entry of this.journalEntries) {
          dates.push(entry.key);
        }
      }
      return dates;
    },
    loggedInUsername: function() {
      return this.$store.state.username;
    },
    canEdit: function() {
      return (
        this.loggedInUsername && this.loggedInUsername === this.entryAuthor
      );
    },
    twitterShareUrl: function() {
      const permalink =
        location.protocol + '//' + location.host + this.$route.fullPath;
      const text =
        encodeURIComponent("Here's what I got done this week ") + permalink;
      return `https://twitter.com/intent/tweet?text=${text}`;
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
    },
  },
  created() {
    this.loadJournalEntries();
  },
  watch: {
    $route(to, from) {
      if (to.params.username != from.params.username) {
        this.loadJournalEntries();
      }
    },
  },
};
</script>

<style scoped>
.show-empty {
  margin: 10px 0px 25px 0px;
}

.author-controls {
  margin: 25px 0px;
  display: flex;
  flex-direction: column-reverse;
}

@media screen and (min-width: 768px) {
  .author-controls {
    flex-direction: row;
    justify-content: space-between;
  }
}

.author-controls .btn:first-of-type {
  margin-top: 25px;
}

@media screen and (min-width: 768px) {
  .author-controls .btn:first-of-type {
    margin: 0px;
  }
}
</style>
