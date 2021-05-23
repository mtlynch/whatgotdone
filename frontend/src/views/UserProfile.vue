<template>
  <div>
    <div class="profile">
      <Avatar class="avatar" :username="username" size="150px" />

      <div class="profile-body">
        <h1>{{ username }}</h1>

        <vue-markdown
          :linkify="false"
          :html="false"
          :anchorAttributes="{rel: 'ugc'}"
          :source="aboutMarkdown"
          class="user-bio"
        ></vue-markdown>

        <p v-if="profileLoaded && !aboutMarkdown" class="no-bio-message">
          This user has not yet created a public bio.
        </p>

        <template v-if="twitterHandle || emailAddress">
          <ul>
            <li v-if="emailAddress">
              <a :href="'mailto:' + emailAddress" class="email-address">{{
                emailAddress
              }}</a>
              (Email)
            </li>
            <li v-if="twitterHandle">
              <a
                :href="'https://twitter.com/' + twitterHandle"
                class="twitter-handle"
                >@{{ twitterHandle }}</a
              >
              (Twitter)
            </li>
          </ul>
        </template>
      </div>
    </div>

    <div class="d-flex justify-content-end">
      <b-button
        class="edit-btn"
        v-if="canEdit"
        to="/profile/edit"
        variant="primary"
        >Edit</b-button
      >
      <b-button
        class="follow-btn"
        v-if="canFollow"
        variant="primary"
        v-on:click="onFollow"
        >Follow</b-button
      >
      <b-button
        class="unfollow-btn"
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

    <p
      v-if="entriesLoaded && recentEntries.length === 0"
      class="no-entries-message"
    >
      This user has not submitted any recent updates.
    </p>
  </div>
</template>

<script>
import Vue from 'vue';
import VueMarkdown from 'vue-markdown';

import {getEntriesFromUser} from '@/controllers/Entries.js';
import {follow, unfollow} from '@/controllers/Follow.js';
import {getUserMetadata} from '@/controllers/User.js';

import Avatar from '@/components/Avatar.vue';
import PartialJournal from '@/components/PartialJournal.vue';

Vue.use(VueMarkdown);

export default {
  name: 'UserProfile',
  components: {
    Avatar,
    VueMarkdown,
    PartialJournal,
  },
  data() {
    return {
      aboutMarkdown: '',
      twitterHandle: null,
      emailAddress: null,
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
      this.recentEntries = [];
      this.profileLoaded = false;
      this.entriesLoaded = false;
    },
    loadProfile: function () {
      getUserMetadata(this.username).then((metadata) => {
        this.aboutMarkdown = metadata.aboutMarkdown;
        this.twitterHandle = metadata.twitterHandle;
        this.emailAddress = metadata.emailAddress;
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
      follow(this.username).then(() => {
        this.$store.commit('addFollowedUser', this.username);
      });
    },
    onUnfollow: function () {
      unfollow(this.username).then(() => {
        this.$store.commit('removeFollowedUser', this.username);
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

.profile .avatar {
  margin-bottom: 1rem;
}

@media screen and (min-width: 768px) {
  .profile .avatar {
    margin-right: 2.5rem;
  }
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
