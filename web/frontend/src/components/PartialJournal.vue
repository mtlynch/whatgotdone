<template>
  <div class="journal">
    <JournalHeader :username="entry.author" :date="entry.date"/>
    <div class="journal-entry">
      <div class="journal-body">
        <vue-markdown
          :linkify="false"
          :html="false"
          :anchorAttributes="{rel: 'nofollow' }"
          :source="entrySnippet"
        ></vue-markdown>
      </div>
      <b-button pill variant="primary" class="read-more" :to="entry.key">More</b-button>
    </div>
  </div>
</template>

<script>
import Vue from "vue";
import VueMarkdown from "vue-markdown";
import JournalHeader from "./JournalHeader.vue";

Vue.use(VueMarkdown);

export default {
  name: "PartialJournal",
  props: {
    entry: Object
  },
  components: {
    VueMarkdown,
    JournalHeader
  },
  computed: {
    entrySnippet: function() {
      const maxLines = 12;
      const entryLines = this.entry.markdown.split("\n");
      if (entryLines.length < maxLines) {
        return entryLines.join("\n");
      }
      return entryLines.slice(0, maxLines).join("\n") + "\n\n...";
    }
  }
};
</script>

<style scoped>
div.journal {
  border: 1px solid rgb(26, 0, 68);
  padding: 15px;
  margin-bottom: 20px;
  background: rgb(79, 87, 161);
}
.header {
  font-style: italic;
  margin-bottom: 15px;
}
.journal-entry {
  padding: 20px 20px 0px 20px;
}
.journal-body {
  text-align: left;
}

.read-more {
  padding: 8px 25px;
}
</style>
