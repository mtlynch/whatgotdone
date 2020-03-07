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

import {getRecent, mergeEntryArrays} from '@/controllers/Recent.js';

export default {
  name: 'Recent',
  components: {
    EntryFeed,
  },
  data() {
    return {
      loadMoreInProgress: false,
    };
  },
  computed: {
    recentEntries() {
      return this.$store.state.recentEntries;
    },
  },
  methods: {
    readEntriesFromServer(start) {
      return getRecent(start);
    },
    readEntriesFromStore() {
      return this.$store.state.recentEntries;
    },
    writeEntriesToStore(newEntries) {
      this.$store.commit(
        'setRecent',
        mergeEntryArrays(this.$store.state.recentEntries, newEntries)
      );
    },
  },
  created() {
    getRecent(/*start=*/ 0).then(recentEntries => {
      this.$store.commit(
        'setRecent',
        mergeEntryArrays(this.recentEntries, recentEntries)
      );
    });
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
