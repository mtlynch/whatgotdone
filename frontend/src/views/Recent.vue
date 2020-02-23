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
      v-bind:disabled="loadMoreInProgress"
      v-on:click="onLoadMore"
      >More Entries</b-button
    >
  </div>
</template>

<script>
import PartialJournal from '@/components/PartialJournal.vue';

import {getRecent, mergeEntryArrays} from '@/controllers/Recent.js';

export default {
  name: 'Recent',
  components: {
    PartialJournal,
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
    onLoadMore() {
      this.loadMoreInProgress = true;
      return getRecent(this.recentEntries.length)
        .then(newEntries => {
          this.$store.commit(
            'setRecent',
            mergeEntryArrays(this.recentEntries, newEntries)
          );
        })
        .finally(() => {
          this.loadMoreInProgress = false;
        });
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
