<template>
  <div>
    <h1>Edit Profile</h1>

    <div class="form-group profile-photo-form">
      <label for="upload-profile-photo">Profile Photo</label>

      <a class="avatar-wrapper" href="#" @click="openFileUploadDialog">
        <Avatar :username="loggedInUsername" size="80px" ref="avatar" />
      </a>

      <div class="profile-photo-btns">
        <b-form-file
          id="upload-profile-photo"
          @input="onUploadProfilePhoto"
          v-model="profilePhoto"
          :plain="true"
          ref="fileUpload"
          input="profile-photo-input"
          accept="image/jpeg"
        ></b-form-file>
        <b-button
          variant="danger"
          id="delete-profile-photo"
          @click="onRemoveProfilePhoto"
          class="profile-photo-btn"
          >Remove</b-button
        >
      </div>
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
        placeholder="@handle"
      />
    </div>

    <div class="form-group">
      <label for="mastodon-address">Mastodon address</label>
      <input
        type="text"
        v-model="profile.mastodonAddress"
        class="form-control"
        id="mastodon-address"
        :disabled="!profileLoaded"
        placeholder="handle@example.com"
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
import {deleteAvatar, uploadAvatar} from '@/controllers/Avatar.js';
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
        mastodonAddress: '',
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
    openFileUploadDialog: function () {
      this.$refs.fileUpload.reset();
      this.$refs.fileUpload.$el.click();
    },
    onUploadProfilePhoto: function (file) {
      uploadAvatar(file)
        .then(() => {
          this.$refs.avatar.refresh();
        })
        .catch((error) => {
          this.formError = error;
        });
    },
    onRemoveProfilePhoto: function () {
      deleteAvatar()
        .then(() => {
          this.$refs.avatar.refresh();
        })
        .catch((error) => {
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

.profile-photo-form {
  display: flex;
  flex-direction: column;
}

.avatar-wrapper {
  margin: 1rem 0;
  display: inline-block;
  align-self: flex-start;
}

.profile-photo-btns {
  display: flex;
  flex-direction: column;
  max-width: 400px;
}

.profile-photo-btns {
  display: flex;
  flex-direction: column;
}

.profile-photo-btns > * {
  margin: 0.5rem 0;
}

#delete-profile-photo {
  max-width: 20ch;
}
</style>
