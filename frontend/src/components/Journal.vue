<template>
  <div class="journal">
    <JournalHeader :entryAuthor="entry.author" :entryDate="entry.date" />
    <div class="journal-body">
      <vue-markdown
        :linkify="false"
        :html="false"
        :anchorAttributes="{rel: 'ugc'}"
        :source="entry.markdown"
      ></vue-markdown>
      <ViewCount class="view-count" />
      <div class="last-modified-date">
        Last modified {{ entry.lastModified | moment('lll') }}
      </div>
    </div>
  </div>
</template>

<script>
import Vue from 'vue';
import VueMarkdown from 'vue-markdown';
import JournalHeader from './JournalHeader.vue';
import ViewCount from '../components/ViewCount.vue';

Vue.use(VueMarkdown);

export default {
  name: 'Journal',
  props: {
    entry: Object,
  },
  components: {
    JournalHeader,
    ViewCount,
    VueMarkdown,
  },
};
</script>

<style scoped>
.journal-body {
  text-align: left;
  overflow: auto;
}

@media screen and (min-width: 768px) {
  overflow: visible;
}

.view-count {
  display: block;
  font-style: italic;
  margin-top: 40px;
}

@media screen and (min-width: 768px) {
  .view-count {
    text-align: right;
    margin-top: 5px;
  }
}

.last-modified-date {
  font-style: italic;
  margin-top: 40px;
}

@media screen and (min-width: 768px) {
  .last-modified-date {
    text-align: right;
    margin-top: 5px;
  }
}
</style>
