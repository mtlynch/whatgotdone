<template>
  <div class="preferences">
    <h1>Preferences</h1>
    <b-form @submit.prevent="handleSubmit">
      <h2>Entry template</h2>
      <b-form-group
        label="Create a template for your weekly What Got Done entries. Every new update you create will start with this text."
        label-for="entry-template-input"
      >
        <textarea-autosize
          class="form-control"
          id="entry-template-input"
          v-model="entryTemplate"
          :disabled="entryTemplate === undefined"
          name="entry-template"
          @input="onEntryTemplateChanged"
          :min-height="250"
          :max-height="650"
          :placeholder="'# Project 1\n\n*Update 1\n*Update 2\n\n# Project 2\n\n* Update 1'"
        ></textarea-autosize>
      </b-form-group>

      <div class="d-flex justify-content-end">
        <button
          type="submit"
          class="btn btn-primary"
          :disabled="!templateChanged"
        >
          Save
        </button>
      </div>
      <div class="d-flex justify-content-end">
        <b-alert
          variant="success"
          class="mt-2"
          dismissible
          fade
          :show="preferencesSaved"
          >Preferences saved</b-alert
        >
      </div>
    </b-form>
  </div>
</template>

<script>
import Vue from 'vue';
import VueTextareaAutosize from 'vue-textarea-autosize';

import {getPreferences, savePreferences} from '@/controllers/Preferences.js';

Vue.use(VueTextareaAutosize);

export default {
  name: 'UserPreferencesPage',
  data() {
    return {
      entryTemplate: undefined,
      entryTemplateFromServer: undefined,
      savingPreferences: false,
      templateChanged: false,
      preferencesSaved: false,
    };
  },
  computed: {
    username: function () {
      return this.$store.state.username;
    },
  },
  methods: {
    loadPreferences() {
      getPreferences()
        .then((preferences) => {
          this.entryTemplateFromServer = preferences.entryTemplate;
          this.entryTemplate = preferences.entryTemplate;
        })
        // Ignore errors on pulling down preferences.
        .catch(() => {
          this.entryTemplate = '';
        });
    },
    onEntryTemplateChanged() {
      // Template is considered changed if either:
      //  - The server version is empty or undefined, and the local version is
      //      non-empty.
      //  - The server version is non-empty and it doesn't match the local
      //      version.
      this.templateChanged =
        (!this.entryTemplateFromServer && this.entryTemplate) ||
        (this.entryTemplateFromServer &&
          this.entryTemplateFromServer !== this.entryTemplate);
    },
    handleSubmit() {
      this.savingPreferences = true;
      savePreferences({
        entryTemplate: this.entryTemplate,
      })
        .then(() => {
          this.preferencesSaved = true;
          this.templateChanged = false;
        })
        .finally(() => {
          this.savingPreferences = false;
        });
    },
  },
  created() {
    if (!this.$store.state.username) {
      this.$router.replace('/login');
      return;
    }
    this.loadPreferences();
  },
};
</script>

<style scoped>
* {
  text-align: left;
}

.preferences h2 {
  font-size: 1.4em;
}
</style>
