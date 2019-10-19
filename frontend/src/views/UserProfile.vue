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

    <b-button v-if="canEdit" to="/profile/edit/" variant="primary" class="edit-btn float-right">Edit</b-button>

    <h2>Recent entries</h2>

    <p>TODO: Show user's recent partial journals</p>
  </div>
</template>

<script>
import Vue from "vue";
import VueMarkdown from "vue-markdown";

Vue.use(VueMarkdown);

export default {
  name: "UserProfile",
  data() {
    return {
      aboutMarkdown: "",
      twitterHandle: null,
      emailAddress: null
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
        })
        .catch(() => {
          // TODO: Handle errors.
        });
    }
  },
  created() {
    this.loadProfile();
  },
  components: {
    VueMarkdown
  }
};
</script>

<style scoped>
* {
  text-align: left;
}

h2 {
  clear: both;
}
</style>