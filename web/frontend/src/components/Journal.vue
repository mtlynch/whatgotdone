<template>
  <div class="journal">
    <div class="journalEntry">
      <div class="journalBody">
        <vue-markdown
          :linkify="false"
          :html="false"
          :anchorAttributes="{rel: 'nofollow' }"
          :source="entry.markdown"
        ></vue-markdown>
      </div>
      <div class="date">Last modified {{ entry.lastModified | moment("lll") }}</div>
    </div>
  </div>
</template>

<script>
import Vue from "vue";
import VueMarkdown from "vue-markdown";

Vue.use(VueMarkdown);

export default {
  name: "Journal",
  props: {
    entry: Object
  },
  components: {
    VueMarkdown
  }
};
</script>

<style scoped>
.journalBody {
  text-align: left;
}
.date {
  font-style: italic;
}
</style>
