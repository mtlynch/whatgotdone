<template>
  <div>
    <h1>{{ username }} : {{ topic }}</h1>

    <template v-if="entriesLoaded && topicBodies.length > 0">
      <div v-for="item in topicBodies" :key="item.key">
        <p>
          <b>{{ item.date | moment('utc', 'dddd, ll') }}</b>
        </p>
        <div class="journal">
          <vue-markdown
            :linkify="false"
            :html="false"
            :anchorAttributes="{rel: 'ugc'}"
            :source="item.markdown"
          ></vue-markdown>
        </div>
      </div>
    </template>

    <p
      v-if="entriesLoaded && topicBodies.length === 0"
      class="no-entries-message"
    >
      This user has not submitted any recent updates.
    </p>
  </div>
</template>

<script>
import Vue from 'vue';
import VueMarkdown from 'vue-markdown';

Vue.use(VueMarkdown);

export default {
  name: 'ViewTopic',
  components: {
    VueMarkdown,
  },
  data() {
    return {
      topicBodies: [],
      entriesLoaded: false,
    };
  },
  computed: {
    username: function() {
      return this.$route.params.username;
    },
    topic: function() {
      return this.$route.params.topic;
    },
  },
  methods: {
    clear: function() {
      this.topicBodies = [];
      this.entriesLoaded = false;
    },
    loadTopicBodies: function() {
      this.topicBodies = [];
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/entries/${this.username}/topic/${this.topic}`;
      this.$http.get(url).then(result => {
        for (const entry of result.data) {
          this.topicBodies.push({
            key: `${entry.date}`,
            date: new Date(entry.date),
            markdown: entry.markdown,
          });
        }
        // Sort newest to oldest.
        this.topicBodies.sort((a, b) => b.date - a.date);
        this.entriesLoaded = true;
      });
    },
  },
  created() {
    this.loadTopicBodies();
  },
  watch: {
    $route(to, from) {
      if (to.params.username != from.params.username) {
        this.clear();
        this.loadtopicBodies();
      }
    },
  },
};
</script>

<style scoped>
* {
  text-align: left;
}

h1 {
  margin-bottom: 50px;
}

h2 {
  clear: both;
  margin-top: 40px;
  margin-bottom: 30px;
}

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
