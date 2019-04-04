<template>
  <div class="journal">
    <div class="journalEntry">
      <div class="date">{{ entry.date | dateFormat('YYYY-MM-DD') }}</div>
      <div class="journalBody">
        <vue-markdown>{{ entry.markdown }}</vue-markdown>
      </div>
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
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
