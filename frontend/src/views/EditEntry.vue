<template>
  <div class="submit">
    <h1>What got done this week?</h1>
    <p>
      Enter your update for the week ending
      <span class="end-date">{{ date | moment('dddd, ll') }}</span
      >.
    </p>

    <template v-if="entryContent === null">
      <b-spinner type="grow" label="Spinning"></b-spinner>
      <p>Loading draft...</p>
    </template>
    <template v-else>
      <form @submit.prevent="handleSubmit">
        <EntryEditor
          class="form-control journal-markdown"
          ref="entryText"
          v-model="entryContent"
          @input="debouncedSaveDraft"
        />
        <p>
          (You can use
          <a href="https://www.markdownguide.org/cheat-sheet/" target="_blank"
            >Markdown</a
          >)
        </p>
        <div class="d-flex justify-content-end">
          <button
            @click.prevent="handleSaveDraft"
            class="btn btn-primary save-draft"
            :disabled="changesSaved"
          >
            {{ saveLabel }}
          </button>
          <button type="submit" class="btn btn-primary">Publish</button>
        </div>
      </form>
      <JournalPreview :markdown="entryContent" />
    </template>
  </div>
</template>

<script>
import Vue from 'vue';
import VueTextareaAutosize from 'vue-textarea-autosize';
import _ from 'lodash';

import {getDraft, saveDraft} from '@/controllers/Drafts.js';
import {saveEntry} from '@/controllers/Entries.js';
import {isValidEntryDate, thisFriday} from '@/controllers/EntryDates.js';

import EntryEditor from '@/components/EntryEditor.vue';
import JournalPreview from '@/components/JournalPreview.vue';

Vue.use(VueTextareaAutosize);

export default {
  name: 'EditEntry',
  components: {
    EntryEditor,
    JournalPreview,
  },
  data() {
    return {
      date: '',
      entryContent: null,
      changesSaved: true,
      saveLabel: 'Save Draft',
    };
  },
  computed: {
    username() {
      return this.$store.state.username;
    },
  },
  methods: {
    loadEntryContent() {
      if (this.date.length == 0 || !this.username) {
        return;
      }
      getDraft(this.date).then(content => {
        this.entryContent = content;
      });
    },
    handleSaveDraft() {
      if (this.entryContent === null) {
        return;
      }
      this.saveLabel = 'Saving';
      saveDraft(this.date, this.entryContent)
        .then(() => {
          this.changesSaved = true;
          this.saveLabel = 'Changes Saved';
        })
        .catch(() => {
          this.changesSaved = false;
        });
    },
    debouncedSaveDraft: _.debounce(function() {
      this.handleSaveDraft();
    }, 2500),
    handleSubmit() {
      saveEntry(this.date, this.entryContent).then(result => {
        this.$router.push(result.path);
      });
    },
  },
  created() {
    if (!this.username) {
      this.$router.replace('/login');
      return;
    }
    if (this.$route.params.date && isValidEntryDate(this.$route.params.date)) {
      this.date = this.$route.params.date;
    } else {
      this.date = thisFriday();
    }
  },
  watch: {
    date: function() {
      this.loadEntryContent();
    },
    username: function() {
      this.loadEntryContent();
    },
    entryContent: function() {
      this.changesSaved = false;
      this.saveLabel = 'Save Draft';
    },
    $route(to, from) {
      if (to.params.date != from.params.date) {
        if (isValidEntryDate(to.params.date)) {
          this.date = to.params.date;
        }
      }
    },
  },
};
</script>

<style scoped>
.submit {
  text-align: left;
  font-size: 11pt;
}

span.end-date {
  color: rgb(255, 208, 56);
  font-weight: bold;
}

.save-draft {
  width: 150px;
  margin-right: 20px;
}
</style>
