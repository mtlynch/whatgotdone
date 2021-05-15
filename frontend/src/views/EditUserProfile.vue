<template>
  <div>
    <h1>Edit Profile</h1>

    <div class="form-group">
      <label for="upload-profile-photo">Profile Photo</label>

      <Avatar :username="loggedInUsername" size="80px" />

      <b-button variant="secondary" id="edit-profile-photo">Edit</b-button>
      <b-button variant="danger" id="delete-profile-photo">Remove</b-button>

      <b-form-file
        id="upload-profile-photo"
        v-on:input="onUploadProfilePhoto"
        v-model="profilePhoto"
        :state="Boolean(profilePhoto)"
        accept="image/jpeg"
      ></b-form-file>
      <b-button variant="secondary" @click="onUploadProfilePhoto"
        >Upload</b-button
      >
    </div>

    <div class="form-group">
      <label for="user-bio">Bio</label>
      <textarea
        id="user-bio"
        v-model="profile.aboutMarkdown"
        :disabled="!profileLoaded"
        class="form-control"
        maxlength="300"
        placeholder="Tell others who you are and what your current projects are"
      />
    </div>

    <div class="form-group">
      <label for="email-address">Public email address</label>
      <input
        type="email"
        v-model="profile.emailAddress"
        class="form-control"
        id="email-address"
        :disabled="!profileLoaded"
        placeholder="name@example.com"
      />
    </div>

    <div class="form-group">
      <label for="twitter-handle">Twitter handle</label>
      <input
        type="text"
        v-model="profile.twitterHandle"
        class="form-control"
        id="twitter-handle"
        :disabled="!profileLoaded"
      />
    </div>

    <div class="alert alert-primary" role="alert" v-if="formError">
      {{ formError }}
    </div>

    <b-button
      variant="primary"
      class="d-block ml-auto"
      @click.prevent="handleSave()"
      id="save-profile"
      >Save</b-button
    >
  </div>
</template>

<script>
import {uploadAvatar} from '@/controllers/Avatar.js';
import {getUserMetadata, setUserMetadata} from '@/controllers/User.js';

import Avatar from '@/components/Avatar.vue';

export default {
  name: 'EditUserProfile',
  components: {
    Avatar,
  },
  data() {
    return {
      profilePhoto: null,
      profile: {
        aboutMarkdown: '',
        twitterHandle: '',
        emailAddress: '',
      },
      profileLoaded: false,
      formError: null,
    };
  },
  computed: {
    loggedInUsername: function () {
      return this.$store.state.username;
    },
  },
  methods: {
    loadProfile: function () {
      getUserMetadata(this.loggedInUsername).then((metadata) => {
        this.profile = metadata;
        this.profileLoaded = true;
      });
    },
    handleSave: function () {
      setUserMetadata(this.profile)
        .then(() => {
          this.$router.push(`/${this.loggedInUsername}`);
        })
        .catch((error) => {
          if (error.response && error.response.data) {
            this.formError = error.response.data;
          } else {
            this.formError = error;
          }
        });
    },
    onUploadProfilePhoto: function (file) {
      uploadAvatar(file).catch((error) => {
        this.formError = error;
      });
    },
  },
  created() {
    this.loadProfile();
  },
};
</script>

<style scoped>
* {
  text-align: left;
}
</style>
