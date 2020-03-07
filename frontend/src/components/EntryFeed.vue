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
      >More Entries</b-button
    >
  </div>
</template>

<script>
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
          this.writeEntriesToStore(newEntries);
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

<style scoped>
h1 {
  text-align: left;
}

p {
  text-align: left;
}
</style>
