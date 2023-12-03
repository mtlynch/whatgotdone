<template>
  <div class="view-entry">
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
        <username-link :username="entryAuthor" />&nbsp;has not posted a journal
        entry for&nbsp;
        <entry-date :date="entryDate" />
      </p>
    </template>
    <template v-else-if="isLoading">
      <b-spinner type="grow" label="Spinning"></b-spinner>
      <p>Loading &nbsp;<username-link :username="entryAuthor" />'s update...</p>
    </template>
    <template v-else>
      <p>
        <username-link :username="entryAuthor" /> has not posted any What Got
        Done updates.
      </p>
    </template>
    <p v-if="backendError" class="error">
      Failed to connect to backend: {{ backendError }}
    </p>
    <div class="author-controls mb-4" v-if="canEdit">
      <div class="ml-auto">
        <b-button variant="danger" @click="onDelete">Unpublish</b-button>
        <b-button
          :to="'/entry/edit/' + this.entryDate"
          variant="primary"
          class="ml-2"
          >Edit</b-button
        >
      </div>
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

import {entryDelete, getEntriesFromUser} from '@/controllers/Entries.js';
import {thisFriday} from '@/controllers/EntryDates.js';

import Journal from '@/components/Journal.vue';
import Reactions from '@/components/Reactions.vue';

export default {
  name: 'ViewEntryPage',
  components: {
    Journal,
    Reactions,
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
    onDelete: function () {
      entryDelete(this.entryDate).then(() => {
        this.$router.push('/entry/edit/' + this.entryDate);
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
          dates.push(new Date(entry.date).toISOString());
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
