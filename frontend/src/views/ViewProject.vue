<template>
  <div>
    <h1><Username :username="username" />'s updates for {{ project }}</h1>

    <template v-if="entriesLoaded && entries.length > 0">
      <ProjectEntry
        class="project-entry"
        v-bind:entry="entry"
        v-bind:key="entry.key"
        v-for="entry in entries"
      />
    </template>

    <p v-if="entriesLoaded && entries.length === 0" class="no-entries-message">
      This user has not submitted any recent updates.
    </p>
  </div>
</template>

<script>
import {getEntriesFromUser} from '@/controllers/Entries.js';

import Username from '@/components/Username.vue';
import ProjectEntry from '@/components/ProjectEntry.vue';

export default {
  name: 'ViewProject',
  components: {
    ProjectEntry,
    Username,
  },
  data() {
    return {
      entries: [],
      entriesLoaded: false,
    };
  },
  computed: {
    username: function() {
      return this.$route.params.username;
    },
    project: function() {
      return this.$route.params.project;
    },
  },
  methods: {
    clear: function() {
      this.entries = [];
      this.entriesLoaded = false;
    },
    loadEntries: function() {
      this.entries = [];
      getEntriesFromUser(this.username, this.project).then(entries => {
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
  margin-bottom: 50px;
}

.project-entry {
  margin-bottom: 50px;
}
</style>
