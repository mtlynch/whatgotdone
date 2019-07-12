<template>
  <div class="journal">
    <JournalHeader :entryAuthor="entry.author" :entryDate="entry.date" />
    <div class="journalEntry">
      <div class="journal-body">
        <vue-markdown
          :linkify="false"
          :html="false"
          :anchorAttributes="{rel: 'nofollow' }"
          :source="entry.markdown"
        ></vue-markdown>
      </div>
      <div class="last-modified-date">Last modified {{ entry.lastModified | moment("lll") }}</div>
    </div>
  </div>
</template>

<script>
import Vue from "vue";
import VueMarkdown from "vue-markdown";
import JournalHeader from "./JournalHeader.vue";

Vue.use(VueMarkdown);

export default {
  name: "Journal",
  props: {
    entry: Object
  },
  components: {
    JournalHeader,
    VueMarkdown
  }
};
</script>

<style scoped>
.journal-body {
  text-align: left;
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
