<template>
  <div class="journal">
    <JournalHeader
      :entryAuthor="entry.author"
      :entryDate="new Date(entry.date)"
    />
    <div class="journal-entry">
      <div class="journal-body">
        <vue-markdown
          :linkify="false"
          :html="false"
          :anchorAttributes="{rel: 'ugc'}"
          :source="entrySnippet"
        ></vue-markdown>
      </div>
      <div class="text-center">
        <b-button variant="primary" :to="entry.key">More</b-button>
      </div>
    </div>
  </div>
</template>

<script>
import Vue from 'vue';
import VueMarkdown from 'vue-markdown';
import JournalHeader from '@/components/JournalHeader.vue';

Vue.use(VueMarkdown);

export default {
  name: 'PartialJournal',
  props: {
    entry: Object,
  },
  components: {
    VueMarkdown,
    JournalHeader,
  },
  computed: {
    entrySnippet: function() {
      const maxLines = 12;
      const entryLines = this.entry.markdown.split('\n');
      if (entryLines.length < maxLines) {
        return entryLines.join('\n');
      }
      return entryLines.slice(0, maxLines).join('\n') + '\n\n...';
    },
  },
};
</script>

<style scoped>
div.journal {
  border: 1px solid rgb(26, 0, 68);
  padding: 10px;
  margin-bottom: 60px;
  background-color: #4e5d6c;
  overflow: auto;
}

@media screen and (min-width: 768px) {
  div.journal {
    padding: 15px;
    margin-bottom: 40px;
    overflow: visible;
  }
}

.header {
  font-style: italic;
  margin-bottom: 15px;
}

.journal-entry {
  padding: 4px 5px;
}

@media screen and (min-width: 768px) {
  .journal-entry {
    padding: 20px 20px 0px 20px;
  }
}

.journal-body {
  text-align: left;
  margin-bottom: 50px;
}
</style>
