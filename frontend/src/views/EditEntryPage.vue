<template>
  <div class="submit">
    <h1>What got done this week?</h1>
    <template v-if="entryContent !== null">
      <form class="entry-form" @submit.prevent="handleSubmit">
        <p>
          Enter your update for the week ending
          <span class="end-date">{{ date | moment('dddd, ll') }}</span
          >.
        </p>
        <MarkdownEditor
          ref="entryText"
          v-model="entryContent"
          @input="debouncedSaveDraft"
        />
        <div class="d-flex justify-content-end">
          <button
            @click.prevent="handleSaveDraft"
            class="btn btn-primary save-draft"
            :disabled="changesSaved"
          >
            {{ saveLabel }}
          </button>
          <button
            type="submit"
            class="btn btn-primary"
            :disabled="!entryHasContent"
          >
            Publish
          </button>
        </div>
      </form>
      <JournalPreview :markdown="entryContent" v-if="entryContent !== null" />
    </template>
    <template v-else-if="errorMessage">
      <b-alert show variant="warning"
        ><h3>Failed to load entry draft</h3>
        <p>{{ errorMessage }}</p>
        <p>Please reload the page to try again.</p>
      </b-alert>
    </template>
    <template v-else>
      <b-spinner type="grow" label="Spinning"></b-spinner>
      <p>Loading draft...</p>
    </template>
  </div>
</template>

<script>
import Vue from 'vue';
import VueTextareaAutosize from 'vue-textarea-autosize';
import _ from 'lodash';

import {draftDelete, getDraft, saveDraft} from '@/controllers/Drafts.js';
import {saveEntry} from '@/controllers/Entries.js';
import {isValidEntryDate, thisFriday} from '@/controllers/EntryDates.js';

import JournalPreview from '@/components/JournalPreview.vue';
import MarkdownEditor from '@/components/MarkdownEditor.vue';

Vue.use(VueTextareaAutosize);

export default {
  name: 'EditEntryPage',
  components: {
    JournalPreview,
    MarkdownEditor,
  },
  data() {
    return {
      date: '',
      entryContent: null,
      errorMessage: null,
      changesSaved: true,
      saveLabel: 'Save Draft',
    };
  },
  computed: {
    username() {
      return this.$store.state.username;
    },
    entryHasContent() {
      return this.entryContent.trim() !== '';
    },
  },
  methods: {
    loadEntryContent() {
      if (this.date.length == 0 || !this.username) {
        return;
      }
      getDraft(this.date)
        .then((content) => {
          this.entryContent = content;
        })
        .catch((error) => {
          this.errorMessage = error;
        });
    },
    onContentChanged() {
      this.changesSaved = false;
      this.debouncedSaveDraft();
    },
    handleSaveDraft() {
      this.changesSaved = false;
      this.saveLabel = 'Saving';

      let saveFn = null;
      if (this.entryHasContent) {
        saveFn = () => {
          return saveDraft(this.date, this.entryContent);
        };
      } else {
        saveFn = () => {
          return draftDelete(this.date);
        };
      }
      saveFn()
        .then(() => {
          this.changesSaved = true;
          this.saveLabel = 'Changes Saved';
        })
        .catch(() => {
          this.changesSaved = false;
        });
    },
    debouncedSaveDraft: _.debounce(function () {
      this.handleSaveDraft();
    }, 2500),
    handleSubmit() {
      if (!this.entryContent) {
        return;
      }
      saveEntry(this.date, this.entryContent).then((result) => {
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
    date: function () {
      this.loadEntryContent();
    },
    username: function () {
      this.loadEntryContent();
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
  color: #2b3e50 !important;
  font-weight: bold;
}

.save-draft {
  width: 150px;
  margin-right: 20px;
}
</style>
