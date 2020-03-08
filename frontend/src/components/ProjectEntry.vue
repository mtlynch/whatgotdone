<template>
  <div class="project-entry">
    <p class="date">
      For the week ending
      <strong>
        <router-link :to="entry.permalink">{{
          entry.date | moment('utc', 'dddd, ll')
        }}</router-link></strong
      >
    </p>
    <vue-markdown
      class="body"
      :linkify="false"
      :html="false"
      :anchorAttributes="{rel: 'ugc'}"
      :source="entry.markdown"
    ></vue-markdown>
  </div>
</template>

<script>
import Vue from 'vue';
import VueMarkdown from 'vue-markdown';

Vue.use(VueMarkdown);

export default {
  name: 'ProjectEntry',
  props: {
    entry: Object,
  },
  components: {
    VueMarkdown,
  },
};
</script>

<style scoped>
* {
  text-align: left;
}

.body {
  border: 1px solid rgb(26, 0, 68);
  padding: 10px 15px 5px 2px;
  background-color: #4e5d6c;
  overflow: auto;
}

@media screen and (min-width: 768px) {
  .body {
    overflow: visible;
  }
}
</style>
