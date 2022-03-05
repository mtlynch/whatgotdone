<template>
  <div>
    <h1><UsernameLink :username="username" />'s updates for {{ project }}</h1>

    <template v-if="entriesLoaded && entries.length > 0">
      <div
        class="journal project-entry"
        v-for="entry in entries"
        v-bind:key="entry.key"
      >
        <ProjectEntry v-bind:entry="entry" />
      </div>
    </template>

    <p v-if="entriesLoaded && entries.length === 0" class="no-entries-message">
      This user has not submitted any recent updates.
    </p>
  </div>
</template>

<script>
import {getEntriesFromUser} from '@/controllers/Entries.js';

import UsernameLink from '@/components/UsernameLink.vue';
import ProjectEntry from '@/components/ProjectEntry.vue';

export default {
  name: 'ViewProjectPage',
  components: {
    ProjectEntry,
    UsernameLink,
  },
  data() {
    return {
      entries: [],
      entriesLoaded: false,
    };
  },
  computed: {
    username: function () {
      return this.$route.params.username;
    },
    project: function () {
      return this.$route.params.project;
    },
  },
  methods: {
    clear: function () {
      this.entries = [];
      this.entriesLoaded = false;
    },
    loadEntries: function () {
      this.entries = [];
      getEntriesFromUser(this.username, this.project).then((entries) => {
        this.entries = entries;
        this.entriesLoaded = true;
      });
    },
  },
  created() {
    this.loadEntries();
  },
};
</script>

<style scoped>
h1 {
  margin-bottom: 2rem;
}

.project-entry + .project-entry {
  margin-top: 5rem;
}
</style>
