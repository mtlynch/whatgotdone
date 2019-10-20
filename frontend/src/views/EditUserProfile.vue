<template>
  <div>
    <h1>Edit Profile</h1>

    <h2>About</h2>

    <textarea v-model="aboutMarkdown" />

    <h2>Contact</h2>

    <ul>
      <li>
        Twitter:
        <input type="text" :value="twitterUsername" />
      </li>
      <li>
        Email:
        <input type="text" :value="email" />
      </li>
    </ul>

    <b-button variant="primary" class="float-right" @click.prevent="handleSave()">Save</b-button>
  </div>
</template>

<script>
import getCsrfToken from "../controllers/CsrfToken.js";

export default {
  name: "EditUserProfile",
  data() {
    return {
      // TODO(mtlynch): Retrieve this from the server.
      aboutMarkdown:
        "Hi, I'm Michael, creator of What Got Done.\n\nI also blog at [mtlynch.io](https://mtlynch.io).",
      twitterUsername: "deliberatecoder",
      email: "michael@mtlynch.io"
    };
  },
  computed: {
    loggedInUsername: function() {
      return this.$store.state.username;
    }
  },
  methods: {
    handleSave: function() {
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/user`;
      this.$http.post(
        url,
        {
          aboutMarkdown: this.aboutMarkdown,
          twitterUsername: this.twitterUsername,
          emailAddress: this.emailAddress
        },
        { withCredentials: true, headers: { "X-CSRF-Token": getCsrfToken() } }
      );
    }
  }
};
</script>

<style scoped>
* {
  text-align: left;
}
</style>