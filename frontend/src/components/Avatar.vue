<template>
  <b-avatar
    variant="info"
    :to="to"
    :src="src"
    :size="size"
    class="avatar"
    ref="avatar"
  ></b-avatar>
</template>

<script>
export default {
  props: {
    username: String,
    to: String,
    size: {
      type: String,
      default: 'md',
    },
  },
  data() {
    return {
      // This just exists for cache-busting.
      version: 1,
    };
  },
  computed: {
    src: function () {
      let size = '40px';
      if (this.size !== 'md') {
        size = '300px';
      }
      return `${process.env.VUE_APP_GCS_PUBLIC_BASE_URL}/avatars/${this.username}/${this.username}-avatar-${size}.jpg?v=${this.version}`;
    },
  },
  methods: {
    // Increment the version number on the avatar image, which causes the browser to request a new version from the server.
    refresh: function () {
      this.version += 1;
    },
  },
};
</script>
