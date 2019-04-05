<template>
  <div class="journal">
    <div class="journalEntry">
      <div class="journalBody">
        <vue-markdown
          :linkify="false"
          :html="false"
          :anchorAttributes="{rel: 'nofollow' }"
        >{{ entry.markdown }}</vue-markdown>
      </div>
      <div class="date">Last modified at {{ entry.lastModified | dateFormat('h:mm a MM/D/YYYY') }}</div>
    </div>
  </div>
</template>

<script>
import Vue from "vue";
import VueFilterDateFormat from "vue-filter-date-format";
import VueMarkdown from "vue-markdown";

Vue.use(VueMarkdown);
Vue.use(VueFilterDateFormat);

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
