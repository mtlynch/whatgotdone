<template>
  <div class="journal">
    <div class="header">
      <router-link :to="entry.author">{{ entry.author }}</router-link>
      's update for the week of {{ entry.date.toLocaleDateString() }}
    </div>
    <div class="journalEntry">
      <div class="journalBody">
        <vue-markdown
          :linkify="false"
          :html="false"
          :anchorAttributes="{rel: 'nofollow' }"
          :source="entrySnippet"
        ></vue-markdown>
      </div>
      <router-link :to="entry.key">Read more</router-link>
    </div>
  </div>
</template>

<script>
import Vue from "vue";
import VueMarkdown from "vue-markdown";

Vue.use(VueMarkdown);

export default {
  name: "PartialJournal",
  props: {
    entry: Object
  },
  components: {
    VueMarkdown
  },
  computed: {
    entrySnippet: function() {
      const maxLines = 15;
      const entryLines = this.entry.markdown.split("\n");
      if (entryLines.length < maxLines) {
        return entryLines.join("\n");
      }
      return entryLines.slice(0, maxLines).join("\n");
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
.journalBody {
  text-align: left;
}
.header {
  font-style: italic;
}
</style>
