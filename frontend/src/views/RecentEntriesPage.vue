<template>
  <div>
    <h1>Recent Entries</h1>
    <p>Check out what other users got done this week:</p>

    <EntryFeed
      :readEntriesFromServer="readEntriesFromServer"
      :readEntriesFromStore="readEntriesFromStore"
      :writeEntriesToStore="writeEntriesToStore"
    />
  </div>
</template>

<script>
import EntryFeed from '@/components/EntryFeed.vue';

import {getRecent} from '@/controllers/Recent.js';

export default {
  name: 'RecentEntriesPage',
  components: {
    EntryFeed,
  },
  data() {
    return {
      loadMoreInProgress: false,
    };
  },
  methods: {
    readEntriesFromServer(start) {
      return getRecent(start);
    },
    readEntriesFromStore() {
      return this.$store.state.recentEntries;
    },
    writeEntriesToStore(newEntries) {
      this.$store.commit('setRecent', newEntries);
    },
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
