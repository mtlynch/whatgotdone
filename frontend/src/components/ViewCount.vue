<template>
  <span class="view-count" v-if="viewCount">
    Viewed {{ viewCount }} times
  </span>
</template>

<script>
export default {
  name: 'ViewCount',
  data() {
    return {
      path: this.$route.path,
      viewCount: null,
    };
  },
  methods: {
    loadViewCount: function() {
      if (!this.path) {
        return;
      }
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/pageViews`;
      this.$http
        .get(url, {
          params: {
            path: this.path,
          },
        })
        .then(result => {
          if (result.data) {
            this.viewCount = result.data.views;
          }
        })
        .catch(() => {
          // Ignore error for view count, as it's non-essential.
        });
    },
  },
  created() {
    this.loadViewCount();
  },
  watch: {
    $route(to) {
      this.viewCount = null;
      this.path = to.path;
      this.loadViewCount();
    },
  },
};
</script>
