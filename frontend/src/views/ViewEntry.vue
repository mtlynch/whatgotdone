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
    <template v-else-if="isLoading">
      <b-spinner type="grow" label="Spinning"></b-spinner>
      <p>Loading &nbsp;<Username :username="entryAuthor" />'s update...</p>
    </template>
    <template v-else>
      <p>
        <Username :username="entryAuthor" /> has not posted any What Got Done
        updates.
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
        class="d-inline-block ml-auto edit-btn"
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

import {getEntriesFromUser} from '@/controllers/Entries.js';
import {thisFriday} from '@/controllers/EntryDates.js';

import Journal from '@/components/Journal.vue';
import Reactions from '@/components/Reactions.vue';
import Username from '@/components/Username.vue';

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
      isLoading: false,
      backendError: null,
    };
  },
  methods: {
    loadJournalEntries: function () {
      this.isLoading = true;
      this.journalEntries = [];
      getEntriesFromUser(this.entryAuthor)
        .then((entries) => {
          // Sort oldest to newest.
          entries.sort((a, b) => a.date - b.date);
          this.journalEntries = entries;
        })
        .catch((error) => {
          this.backendError = error;
        })
        .finally(() => {
          this.isLoading = false;
        });
    },
  },
  computed: {
    pages: function () {
      let pages = [];
      for (const date of this.entryDates) {
        const m = moment(date).utc();
        const linkFormattedDate = m.format('YYYY-MM-DD');
        pages.push({
          link: `/${this.entryAuthor}/${linkFormattedDate}`,
          text: m.format('MMM. D').replace('May.', 'May'),
        });
      }
      return pages;
    },
    entryDates: function () {
      let dates = [];
      if (this.showEmptyEntries) {
        let d = moment(this.journalEntries[0].date.toISOString());
        while (d <= moment(thisFriday())) {
          dates.push(d.toISOString());
          d = d.add(1, 'weeks');
        }
      } else {
        for (const entry of this.journalEntries) {
          dates.push(entry.date.toISOString());
        }
      }
      return dates;
    },
    loggedInUsername: function () {
      return this.$store.state.username;
    },
    canEdit: function () {
      return (
        this.currentEntry &&
        this.loggedInUsername &&
        this.loggedInUsername === this.entryAuthor
      );
    },
    twitterShareUrl: function () {
      const permalink =
        location.protocol + '//' + location.host + this.$route.fullPath;
      const text =
        encodeURIComponent("Here's what I got done this week ") + permalink;
      return `https://twitter.com/intent/tweet?text=${text}`;
    },
    entryAuthor: function () {
      return this.$route.params.username;
    },
    entryDate: function () {
      return this.$route.params.date;
    },
    currentEntry: function () {
      if (!this.entryDate) {
        return null;
      }
      for (const entry of this.journalEntries) {
        if (this.$route.path === entry.permalink) {
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
