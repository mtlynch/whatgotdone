<template>
  <div>
    <h1>Edit Profile</h1>

    <h2>About</h2>

    <textarea v-model="aboutMarkdown" />

    <h2>Contact</h2>

    <ul>
      <li>
        Twitter:
        <input type="text" v-model="twitterHandle" maxlength="15" />
      </li>
      <li>
        Email:
        <input type="text" v-model="emailAddress" />
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
      aboutMarkdown: "",
      twitterHandle: "",
      emailAddress: ""
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
      this.$http
        .post(
          url,
          {
            aboutMarkdown: this.aboutMarkdown,
            twitterHandle: this.twitterHandle,
            emailAddress: this.emailAddress
          },
          { withCredentials: true, headers: { "X-CSRF-Token": getCsrfToken() } }
        )
        .then(result => {
          if (result.data.ok) {
            this.$router.push(`/${this.loggedInUsername}`);
          }
        })
        .catch(() => {
          // TODO: Handle error.
        });
    }
  }
};
</script>

<style scoped>
* {
  text-align: left;
}
</style>