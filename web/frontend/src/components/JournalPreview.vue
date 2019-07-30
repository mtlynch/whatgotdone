<template>
  <div class="journal-preview">
    <p class="preview-bar">preview</p>
    <div class="content">
      <vue-markdown
        :linkify="false"
        :html="false"
        :anchorAttributes="{rel: 'nofollow' }"
        :source="markdown"
      ></vue-markdown>
    </div>
  </div>
</template>

<script>
import Vue from "vue";
import VueMarkdown from "vue-markdown";

Vue.use(VueMarkdown);

export default {
  name: "JournalPreview",
  props: {
    markdown: String
  },
  components: {
    VueMarkdown
  }
};
</script>

<style scoped>
.journal-preview {
  margin-top: 40px;
  background-color: #4e5d6c;
  border: 1px solid rgb(90, 90, 90);
  text-align: left;
}

.preview-bar {
  background-color: #3f464d;
  font-size: 9pt;
  padding: 4px 2px 4px 14px;
}

.content {
  padding: 4px 14px;
}
</style>
