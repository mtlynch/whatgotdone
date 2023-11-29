<template>
  <div>
    <div class="profile">
      <user-avatar :username="username" />

      <div class="profile-body">
        <h1>{{ username }}</h1>

        <vue-markdown
          :linkify="false"
          :html="false"
          :anchorAttributes="{rel: 'ugc'}"
          :source="aboutMarkdown"
          data-testid="user-bio"
        ></vue-markdown>

        <p v-if="profileLoaded && !aboutMarkdown" class="font-italic">
          This user has not yet created a public bio.
        </p>

        <template v-if="twitterHandle || emailAddress">
          <ul>
            <li v-if="emailAddress">
              <a :href="'mailto:' + emailAddress" data-testid="email-address">{{
                emailAddress
              }}</a>
              (Email)
            </li>
            <li v-if="twitterHandle">
              <a
                :href="'https://twitter.com/' + twitterHandle"
                data-testid="twitter-handle"
                >@{{ twitterHandle }}</a
              >
              (Twitter)
            </li>
            <li v-if="mastodonAddress">
              <a data-testid="mastodon-address" :href="mastodonUrl">{{
                mastodonAddress
              }}</a>
              (Mastodon)
            </li>
          </ul>
        </template>
      </div>
    </div>

    <div class="d-flex justify-content-end">
      <b-button v-if="canEdit" to="/profile/edit" variant="primary"
        >Edit</b-button
      >
      <b-button
        data-testid="follow-btn"
        v-if="canFollow"
        variant="primary"
        v-on:click="onFollow"
        >Follow</b-button
      >
      <b-button
        data-testid="unfollow-btn"
        v-if="canUnfollow"
        variant="primary"
        v-on:click="onUnfollow"
        >Unfollow</b-button
      >
    </div>

    <h2>Recent entries</h2>

    <PartialJournal
      v-bind:key="entry.permalink"
      v-bind:entry="entry"
      v-for="entry in recentEntries"
    />

    <p v-if="entriesLoaded && recentEntries.length === 0" class="font-italic">
      This user has not submitted any recent updates.
    </p>
    <template v-if="canEdit">
      <h2>Export</h2>
      <p>
        Download a backup of your account data, including all What Got Done
        drafts, entries, and preferences.
      </p>

      <b-button variant="primary" v-on:click="onExport">Download</b-button>
      <a class="d-none" ref="file-download-helper"
        ><!-- Dummy element to allow file downloads --></a
      >
    </template>
  </div>
</template>

<script>
import Vue from 'vue';
import VueMarkdown from 'vue-markdown';

import {getEntriesFromUser} from '@/controllers/Entries.js';
import {exportGet} from '@/controllers/Export.js';
import {follow, unfollow} from '@/controllers/Follow.js';
import {getUserMetadata} from '@/controllers/User.js';

import PartialJournal from '@/components/PartialJournal.vue';

Vue.use(VueMarkdown);

export default {
  name: 'UserProfilePage',
  components: {
    VueMarkdown,
    PartialJournal,
  },
  data() {
    return {
      aboutMarkdown: '',
      twitterHandle: null,
      emailAddress: null,
      mastodonAddress: null,
      recentEntries: [],
      profileLoaded: false,
      entriesLoaded: false,
    };
  },
  computed: {
    username: function () {
      return this.$route.params.username;
    },
    loggedInUsername: function () {
      return this.$store.state.username;
    },
    canEdit: function () {
      return this.loggedInUsername && this.loggedInUsername === this.username;
    },
    isFollowing: function () {
      if (!this.$store.state.following) {
        return false;
      }
      return this.$store.state.following.includes(this.username);
    },
    avatarLink: function () {
      if (this.canEdit) {
        return '/profile/edit';
      }
      return null;
    },
    mastodonUrl: function () {
      if (!this.mastodonAddress) {
        return null;
      }
      const addressParts = this.mastodonAddress.split('@');
      const username = addressParts[0];
      const domain = addressParts[1];
      return `https://${domain}/@${username}`;
    },
    isSelf: function () {
      if (!this.loggedInUsername) {
        return false;
      }
      return this.loggedInUsername == this.username;
    },
    canFollow: function () {
      return this.loggedInUsername && !this.isFollowing && !this.isSelf;
    },
    canUnfollow: function () {
      return this.loggedInUsername && this.isFollowing && !this.isSelf;
    },
  },
  methods: {
    clear: function () {
      this.aboutMarkdown = '';
      this.twitterHandle = null;
      this.emailAddress = null;
      this.mastodonAddress = null;
      this.recentEntries = [];
      this.profileLoaded = false;
      this.entriesLoaded = false;
    },
    loadProfile: function () {
      getUserMetadata(this.username).then((metadata) => {
        this.aboutMarkdown = metadata.aboutMarkdown;
        this.twitterHandle = metadata.twitterHandle;
        this.emailAddress = metadata.emailAddress;
        this.mastodonAddress = metadata.mastodonAddress;
        this.profileLoaded = true;
      });
    },
    loadRecentEntries: function () {
      this.recentEntries = [];
      getEntriesFromUser(this.username).then((entries) => {
        this.recentEntries = entries;
        this.entriesLoaded = true;
      });
    },
    onFollow: function () {
      this.$store.commit('addFollowedUser', this.username);
      follow(this.username);
    },
    onUnfollow: function () {
      this.$store.commit('removeFollowedUser', this.username);
      unfollow(this.username);
    },
    onExport: function () {
      exportGet().then((exportedData) => {
        const blobURL = window.URL.createObjectURL(
          new Blob([JSON.stringify(exportedData, null, 2)], {
            type: 'text/json',
          })
        );
        const downloadHelper = this.$refs['file-download-helper'];
        downloadHelper.style = 'display: none';
        downloadHelper.href = blobURL;
        const timestamp = new Date()
          .toISOString()
          .replace(/:/g, '')
          .replace(/\.\d+/g, '');
        downloadHelper.download = `whatgotdone-${this.username}-${timestamp}.json`;
        downloadHelper.click();
      });
    },
  },
  created() {
    this.loadProfile();
    this.loadRecentEntries();
  },
  watch: {
    $route(to, from) {
      if (to.params.username != from.params.username) {
        this.clear();
        this.loadProfile();
        this.loadRecentEntries();
      }
    },
  },
};
</script>

<style scoped>
* {
  text-align: left;
}

.profile {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

@media screen and (min-width: 768px) {
  .profile {
    flex-direction: row;
  }
}

h2 {
  clear: both;
  margin-top: 40px;
  margin-bottom: 30px;
}
</style>
