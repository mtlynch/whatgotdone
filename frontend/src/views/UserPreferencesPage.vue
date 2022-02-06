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
          v-model="preferences.entryTemplate"
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
          :disabled="!preferencesChanged"
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
      preferences: {
        entryTemplate: '',
      },
      preferencesFromServer: null,
      savingPreferences: false,
      preferencesChanged: false,
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
          this.preferencesFromServer = Object.assign({}, preferences);
          this.preferences = Object.assign({}, preferences);
        })
        // Ignore errors on pulling down preferences.
        .catch(() => {});
    },
    onEntryTemplateChanged() {
      if (
        this.preferencesFromServer &&
        this.preferencesFromServer.entryTemplate !=
          this.preferences.entryTemplate
      ) {
        this.preferencesChanged = true;
      }
    },
    handleSubmit() {
      this.savingPreferences = true;
      savePreferences(this.preferences)
        .then(() => {
          this.preferencesSaved = true;
          this.preferencesChanged = false;
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
