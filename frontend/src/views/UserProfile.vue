<template>
  <div>
    <h1>{{ username }}</h1>

    <h2>About</h2>

    <vue-markdown
      :linkify="false"
      :html="false"
      :anchorAttributes="{rel: 'ugc' }"
      :source="aboutMarkdown"
    ></vue-markdown>

    <template v-if="twitterHandle || emailAddress">
      <h2>Contact</h2>

      <ul>
        <li v-if="twitterHandle">
          <a :href="'https://twitter.com/' + twitterHandle">@{{ twitterHandle }}</a> (Twitter)
        </li>
        <li v-if="emailAddress">
          <a :href="'mailto:' + emailAddress">{{ emailAddress }}</a> (Email)
        </li>
      </ul>
    </template>

    <b-button
      v-if="canEdit"
      to="/profile/edit"
      variant="primary"
      class="edit-btn float-right"
      :disabled="!profileLoaded"
    >Edit</b-button>

    <h2>Recent entries</h2>

    <PartialJournal v-bind:key="item.key" v-bind:entry="item" v-for="item in recentEntries" />
  </div>
</template>

<script>
import Vue from "vue";
import VueMarkdown from "vue-markdown";
import PartialJournal from "../components/PartialJournal.vue";

Vue.use(VueMarkdown);

export default {
  name: "UserProfile",
  components: {
    VueMarkdown,
    PartialJournal
  },
  data() {
    return {
      aboutMarkdown: "",
      twitterHandle: null,
      emailAddress: null,
      recentEntries: [],
      profileLoaded: false
    };
  },
  computed: {
    username: function() {
      return this.$route.params.username;
    },
    loggedInUsername: function() {
      return this.$store.state.username;
    },
    canEdit: function() {
      // TODO: Delete the next line.
      return true;
      //return this.loggedInUsername && this.loggedInUsername === this.username;
    }
  },
  methods: {
    loadProfile: function() {
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/user/${this.username}`;
      this.$http
        .get(url)
        .then(result => {
          this.aboutMarkdown = result.data.aboutMarkdown;
          this.twitterHandle = result.data.twitterHandle;
          this.emailAddress = result.data.emailAddress;
          this.profileLoaded = true;
        })
        .catch(() => {
          // TODO: Handle errors.
        });
    },
    loadrecentEntries: function() {
      this.recentEntries = [];
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/entries/${this.username}`;
      this.$http
        .get(url)
        .then(result => {
          for (const entry of result.data) {
            this.recentEntries.push({
              key: `${this.username}/${entry.date}`,
              author: this.username,
              date: new Date(entry.date),
              lastModified: new Date(entry.lastModified),
              markdown: entry.markdown
            });
          }
          if (this.recentEntries.length == 0) {
            return;
          }
          this.recentEntries.sort((a, b) => a.date - b.date);
        })
        .catch(() => {
          // TODO: Handle errors.
        });
    }
  },
  created() {
    this.loadProfile();
    this.loadrecentEntries();
  }
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
</style>