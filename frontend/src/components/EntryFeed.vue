<template>
  <div>
    <PartialJournal
      v-bind:key="entry.key"
      v-bind:entry="entry"
      v-for="entry in entries"
    />

    <b-button
      variant="secondary"
      v-bind:disabled="loadMoreInProgress"
      v-on:click="onLoadMore"
      v-if="serverHasMore"
      >More Entries</b-button
    >
  </div>
</template>

<script>
import {mergeEntryArrays} from '@/controllers/Recent.js';

import PartialJournal from '@/components/PartialJournal.vue';

export default {
  name: 'EntryFeed',
  components: {
    PartialJournal,
  },
  props: {
    readEntriesFromServer: Function,
    readEntriesFromStore: Function,
    writeEntriesToStore: Function,
  },
  data() {
    return {
      loadMoreInProgress: false,
      serverHasMore: true,
    };
  },
  computed: {
    entries() {
      return this.readEntriesFromStore();
    },
  },
  methods: {
    onLoadMore() {
      this.loadMoreInProgress = true;
      return this.readEntriesFromServer(this.entries.length)
        .then(newEntries => {
          if (newEntries.length === 0) {
            this.serverHasMore = false;
            return;
          }
          this.writeEntriesToStore(
            mergeEntryArrays(this.readEntriesFromStore(), newEntries)
          );
        })
        .finally(() => {
          this.loadMoreInProgress = false;
        });
    },
  },
  created() {
    this.readEntriesFromServer(/*start=*/ 0).then(entries => {
      this.writeEntriesToStore(entries);
    });
  },
};
</script>
