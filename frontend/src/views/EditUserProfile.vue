<template>
  <div>
    <h1>Edit Profile</h1>

    <p>Fill out your public profile below:</p>

    <div class="form-group">
      <label for="userBio">Bio</label>
      <textarea
        id="userBio"
        v-model="aboutMarkdown"
        :disabled="!profileLoaded"
        class="form-control"
        maxlength="300"
        placeholder="Tell others who you are and what your current projects are"
      />
    </div>

    <div class="form-group">
      <label for="emailAddress">Public email address</label>
      <input
        type="email"
        v-model="emailAddress"
        class="form-control"
        id="emailAddress"
        :disabled="!profileLoaded"
        placeholder="name@example.com"
      />
    </div>

    <div class="form-group">
      <label for="twitterHandle">Twitter handle</label>
      <input
        type="text"
        v-model="twitterHandle"
        class="form-control"
        id="twitterHandle"
        :disabled="!profileLoaded"
      />
    </div>

    <div class="alert alert-primary" role="alert" v-if="formError">{{ formError }}</div>

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
      emailAddress: "",
      profileLoaded: false,
      formError: null
    };
  },
  computed: {
    loggedInUsername: function() {
      return this.$store.state.username;
    }
  },
  methods: {
    loadProfile: function() {
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/user/${this.loggedInUsername}`;
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
        .catch(error => {
          if (error.response && error.response.data) {
            this.formError = error.response.data;
          } else {
            this.formError = error;
          }

          // TODO: Handle error.
        });
    }
  },
  created() {
    this.loadProfile();
  }
};
</script>

<style scoped>
* {
  text-align: left;
}
</style>