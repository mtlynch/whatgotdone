<template>
  <b-avatar variant="info" :src="src" :size="size" class="avatar"></b-avatar>
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
  computed: {
    src: function () {
      let size = 'thumb-30px';
      if (this.size !== 'md') {
        size = 'full-300px';
      }
      return `https://storage.googleapis.com/${this.username}/avatars/%s/${this.username}-avatar-${size}.jpg`;
    },
  },
};
</script>

<style scoped>
/* Prevent us from using the standard link color, which distorts the avatar image. */
.avatar a {
  color: initial !important;
}
</style>
