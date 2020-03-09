<template>
  <div class="preferences">
    <h1>Preferences</h1>
    <b-form @submit.prevent="handleSubmit">
      <textarea-autosize
        class="form-control"
        v-model="preferences.entryTemplate"
        name="entry-template"
        @input="onEntryTemplateChanged"
        :min-height="250"
        :max-height="650"
        placeholder="Enter boilerplate to structure each of your weekly updates"
      ></textarea-autosize>
      <div class="d-flex justify-content-end">
        <button
          type="submit"
          class="btn btn-primary"
          :disabled="preferencesSaved"
        >
          Save
        </button>
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
  name: 'UserPreferences',
  data() {
    return {
      preferences: {
        entryTemplate: '',
      },
      savingPreferences: false,
      preferencesSaved: true,
    };
  },
  computed: {
    username: function() {
      return this.$store.state.username;
    },
  },
  methods: {
    loadPreferences() {
      getPreferences().then(preferences => {
        this.preferences = preferences;
      });
    },
    onEntryTemplateChanged() {
      this.preferencesSaved = false;
    },
    handleSubmit() {
      this.savingPreferences = true;
      savePreferences(this.preferences)
        .then(() => {
          this.preferencesSaved = true;
        })
        .finally(() => {
          this.savingPreferences = false;
        });
    },
  },
  created() {
    this.loadPreferences();
  },
};
</script>

<style scoped>
* {
  text-align: left;
}
</style>
