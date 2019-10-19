<template>
  <div>
    <img :src="avatarUrl" />

    <h1>{{ username }}</h1>

    <p>Member since: {{ joinDate }}</p>
    <p>Updates: {{ entryCount }}</p>

    <h2>About</h2>

    <vue-markdown
      :linkify="false"
      :html="false"
      :anchorAttributes="{rel: 'ugc' }"
      :source="aboutMarkdown"
    ></vue-markdown>

    <h2>Contact</h2>

    <ul>
      <li><a :href="'https://twitter.com/' + twitterUsername">@{{ twitterUsername }}</a> (Twitter)</li>
      <li><a :href="'mailto:' + email">{{ email }}</a> (Email)</li>
    </ul>

    <b-button
      v-if="canEdit"
      to="/profile/edit/"
      variant="primary"
      class="edit-btn float-right"
    >Edit</b-button>
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
      // TODO(mtlynch): Retrieve this from the server.
      avatarUrl: "https://i.stack.imgur.com/rdgMY.jpg",
      username: "michael",
      joinDate: "2019-03-01",
      entryCount: 26,
      aboutMarkdown: "Hi, I'm Michael, creator of What Got Done.\n\nI also blog at [mtlynch.io](https://mtlynch.io).",
      twitterUsername: "deliberatecoder",
      email: "michael@mtlynch.io",
      canEdit: true,
    };
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
</style>