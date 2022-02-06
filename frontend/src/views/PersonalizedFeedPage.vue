<template>
  <div>
    <h1>Feed</h1>
    <template v-if="isFollowingAnyone">
      <p>Here are the latest updates from users you're following:</p>

      <EntryFeed
        :readEntriesFromServer="readEntriesFromServer"
        :readEntriesFromStore="readEntriesFromStore"
        :writeEntriesToStore="writeEntriesToStore"
      />
    </template>
    <template v-else>
      <b-alert show variant="secondary"
        >You're not following anyone yet. Check the updates in the
        <router-link to="/recent">Recent feed</router-link> to find users to
        follow.</b-alert
      >
    </template>
  </div>
</template>

<script>
import EntryFeed from '@/components/EntryFeed.vue';

import {getRecentFollowing} from '@/controllers/Recent.js';

export default {
  name: 'PersonalizedFeedPage',
  components: {
    EntryFeed,
  },
  computed: {
    isFollowingAnyone: function () {
      if (!this.$store.state.following) {
        return false;
      }
      return this.$store.state.following.length > 0;
    },
  },
  methods: {
    readEntriesFromServer(start) {
      return getRecentFollowing(start);
    },
    readEntriesFromStore() {
      return this.$store.state.recentFollowingEntries;
    },
    writeEntriesToStore(newEntries) {
      this.$store.commit('setRecentFollowing', newEntries);
    },
  },
};
</script>

<style scoped>
h1 {
  text-align: left;
}

p {
  text-align: left;
}
</style>
