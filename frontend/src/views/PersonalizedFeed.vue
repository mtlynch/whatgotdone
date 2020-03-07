<template>
  <div>
    <h1>Feed</h1>
    <p>Here are the latest updates from users you're following:</p>

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
  name: 'PersonalizedFeed',
  components: {
    EntryFeed,
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
