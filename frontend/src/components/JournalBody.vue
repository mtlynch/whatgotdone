<template>
  <div class="journal-body">
    <vue-markdown
      :linkify="false"
      :html="false"
      :anchorAttributes="{rel: 'ugc'}"
      :source="source"
    ></vue-markdown>
  </div>
</template>

<script>
import Vue from 'vue';
import VueMarkdown from 'vue-markdown';

Vue.use(VueMarkdown);

export default {
  name: 'JournalBody',
  props: {
    source: String,
  },
  components: {
    VueMarkdown,
  },
};
</script>

<style scoped>
.journal-body {
  text-align: left;
  overflow: auto;
}

.journal-body >>> .contains-task-list {
  list-style-type: none;
}

.journal-body >>> blockquote {
  background: #4c5b68;
  border-left: 5px solid #ccc;
  margin: 1rem 0.5rem;
  padding: 0.5em 10px;
  quotes: '\201C''\201D''\2018''\2019';
  display: inline-block;
  font-style: italic;
}

.journal-body >>> blockquote:before {
  color: #ccc;
  content: open-quote;
  font-size: 4em;
  line-height: 0.1em;
  margin-right: 0.25em;
  vertical-align: -0.4em;
}

.journal-body >>> blockquote p {
  display: inline;
}
</style>
