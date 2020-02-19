<template>
  <div>
    <h1>Recent Entries</h1>
    <p>Check out what other users got done this week:</p>

    <PartialJournal
      v-bind:key="item.key"
      v-bind:entry="item"
      v-for="item in recentEntries"
    />

    <b-button
      variant="secondary"
      v-bind:disabled="requestInProgress"
      v-on:click="onLoadMore"
      >More Entries</b-button
    >
  </div>
</template>

<script>
import PartialJournal from '../components/PartialJournal.vue';

import {refreshRecent, extendRecent} from '../controllers/Recent.js';

export default {
  name: 'Recent',
  components: {
    PartialJournal,
  },
  data() {
    return {
      requestInProgress: false,
    };
  },
  computed: {
    recentEntries() {
      return this.$store.state.recentEntries;
    },
  },
  methods: {
    onLoadMore() {
      /**
       * To prevent muliple click, disable button while api call is in progress
       */
      this.requestInProgress = true; // Set the flag true before fetch call
      return extendRecent(() => {
        this.requestInProgress = false; // Reset the flag in callback
      });
    },
  },
  created() {
    refreshRecent();
  },
};
</script>

<style scoped>
h1 {
  text-align: left;
}

p {
  text-align: left;
}
</style>
