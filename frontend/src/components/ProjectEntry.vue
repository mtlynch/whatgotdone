<template>
  <div class="project-entry">
    <p class="date">
      For the week ending&nbsp;
      <router-link :to="entry.permalink"
        ><entry-date :date="entry.date"
      /></router-link>
    </p>
    <JournalBody :source="entry.markdown" />
  </div>
</template>

<script>
import Vue from 'vue';
import VueMarkdown from 'vue-markdown';

import JournalBody from '@/components/JournalBody';

Vue.use(VueMarkdown);

export default {
  name: 'ProjectEntry',
  props: {
    entry: Object,
  },
  components: {
    JournalBody,
  },
};
</script>
