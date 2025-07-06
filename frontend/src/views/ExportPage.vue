<template>
  <div>
    <h1>Export Data</h1>
    <p>
      Download a backup of your account data, including all What Got Done
      drafts, entries, and preferences.
    </p>

    <b-button variant="primary" v-on:click="onExport">Download (JSON)</b-button>

    <p>
      Download only your public posts in Markdown format. (optimized for
      migrating to a static site generator)
    </p>
    <b-button variant="secondary" v-on:click="onExportMarkdown" class="ml-2"
      >Download (Markdown)</b-button
    >
    <a class="d-none" ref="file-download-helper"
      ><!-- Dummy element to allow file downloads --></a
    >
  </div>
</template>

<script>
import {exportGet, exportMarkdown} from '@/controllers/Export.js';

export default {
  name: 'ExportPage',
  computed: {
    loggedInUsername: function () {
      return this.$store.state.username;
    },
  },
  methods: {
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
        downloadHelper.download = `whatgotdone-${this.loggedInUsername}-${timestamp}.json`;
        downloadHelper.click();
      });
    },
    onExportMarkdown: function () {
      exportMarkdown().then(({blob, filename}) => {
        const blobURL = window.URL.createObjectURL(blob);
        const downloadHelper = this.$refs['file-download-helper'];
        downloadHelper.style = 'display: none';
        downloadHelper.href = blobURL;
        downloadHelper.download = filename;
        downloadHelper.click();
      });
    },
  },
  created() {
    // Redirect to login if not authenticated
    if (!this.loggedInUsername) {
      this.$router.push('/login');
    }
  },
};
</script>

<style scoped>
* {
  text-align: left;
}

h1 {
  margin-bottom: 30px;
}
</style>
