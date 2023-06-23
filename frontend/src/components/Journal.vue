<template>
  <div class="journal">
    <JournalHeader :entryAuthor="entry.author" :entryDate="entry.date" />
    <JournalBody :source="entry.markdown" />
    <div class="metadata">
      <div class="last-modified-date">
        Last modified {{ entry.lastModified | moment('lll') }}
      </div>
    </div>
  </div>
</template>

<script>
import Vue from 'vue';
import VueMarkdown from 'vue-markdown';
import JournalBody from '@/components/JournalBody.vue';
import JournalHeader from '@/components/JournalHeader.vue';

Vue.use(VueMarkdown);

export default {
  name: 'Journal',
  props: {
    entry: Object,
  },
  components: {
    JournalBody,
    JournalHeader,
  },
};
</script>

<style scoped>
.metadata {
  font-style: italic;
  margin-top: 3rem;
}

@media screen and (min-width: 768px) {
  .metadata {
    text-align: right;
    margin-top: 1rem;
  }
}

.view-count {
  display: block;
}
</style>
