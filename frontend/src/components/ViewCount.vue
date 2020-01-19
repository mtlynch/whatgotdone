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
      console.log('in loadViewCount');
      if (!this.path) {
        console.log('path is empty, returning');
        return;
      }
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/pageViews`;
      console.log(`About to make request to ${url}`);
      this.$http
        .get(url, {
          params: {
            path: this.path,
          },
        })
        .then(result => {
          console.log('Got view count');
          this.viewCount = result.data.views;
        })
        .catch(e => {
          console.log('Failed to get view count');
          console.log(e);
          // Ignore error for view count, as it's non-essential.
        });
    },
  },
  created() {
    console.log('created ViewCount');
    this.loadViewCount();
  },
};
</script>
