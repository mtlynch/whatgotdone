<template>
  <div>
    <h1><Username :username="username" />'s updates for {{ project }}</h1>

    <template v-if="entriesLoaded && projectBodies.length > 0">
      <ProjectEntry
        class="project-entry"
        v-bind:entry="item"
        v-for="item in projectBodies"
        :key="item.key"
      />
    </template>

    <p
      v-if="entriesLoaded && projectBodies.length === 0"
      class="no-entries-message"
    >
      This user has not submitted any recent updates.
    </p>
  </div>
</template>

<script>
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
      projectBodies: [],
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
      this.projectBodies = [];
      this.entriesLoaded = false;
    },
    loadProjectBodies: function() {
      this.projectBodies = [];
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/entries/${this.username}/project/${this.project}`;
      this.$http.get(url).then(result => {
        for (const entry of result.data) {
          this.projectBodies.push({
            key: `${entry.date}`,
            date: new Date(entry.date),
            markdown: entry.markdown,
            sourceUrl: `/${this.username}/${entry.date}`,
          });
        }
        // Sort newest to oldest.
        this.projectBodies.sort((a, b) => b.date - a.date);
        this.entriesLoaded = true;
      });
    },
  },
  created() {
    this.loadProjectBodies();
  },
  watch: {
    $route(to, from) {
      if (to.params.username != from.params.username) {
        this.clear();
        this.loadprojectBodies();
      }
    },
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
