<template>
  <div>
    <h1>{{ username }} : {{ topic }}</h1>

    <template v-if="entriesLoaded && topicBodies.length > 0">
      <div v-for="item in topicBodies" :key="item.key">
        <p>
          For the week ending on
          <b>{{ item.date | moment('utc', 'dddd, ll') }}</b>
        </p>
        <vue-markdown
          :linkify="false"
          :html="false"
          :anchorAttributes="{rel: 'ugc'}"
          :source="item.markdown"
        ></vue-markdown>
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

.no-bio-message {
  font-style: italic;
}

.no-entries-message {
  font-style: italic;
}
</style>
