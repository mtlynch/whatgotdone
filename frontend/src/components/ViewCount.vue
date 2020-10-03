<template>
  <span class="view-count" v-if="viewCount">
    Viewed {{ viewCount }} times
  </span>
</template>

<script>
import {getPageViews} from '@/controllers/PageViews.js';

export default {
  name: 'ViewCount',
  data() {
    return {
      path: this.$route.path,
      viewCount: null,
    };
  },
  methods: {
    loadViewCount: function () {
      if (!this.path) {
        return;
      }
      getPageViews(this.path)
        .then((viewCount) => {
          if (viewCount) {
            this.viewCount = viewCount;
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
