<template>
  <div class="preferences">
    <h1>Preferences</h1>
    <b-form @submit.prevent="handleSubmit">
      <textarea-autosize
        class="form-control"
        v-model="preferences.entryTemplate"
        name="entry-template"
        :min-height="250"
        :max-height="650"
        placeholder="Enter boilerplate to structure each of your weekly updates"
      ></textarea-autosize>
      <div class="d-flex justify-content-end">
        <b-button :to="'/' + username" class="btn btn-cancel">
          Cancel
        </b-button>
        <button type="submit" class="btn btn-primary">Save</button>
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
    };
  },
  computed: {
    username: function() {
      return this.$store.state.username;
    },
  },
  methods: {
    loadPreferences() {
      getPreferences.then(preferences => {
        this.preferences = preferences;
      });
    },
    handleSubmit() {
      console.log('saving preferences');
      console.log(this.preferences);
      savePreferences(this.preferences).then(() => {
        this.$router.push(`/${this.username}`);
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
