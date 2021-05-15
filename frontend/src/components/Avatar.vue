<template>
  <b-avatar variant="info" :src="src" :to="'/' + username" size="md"></b-avatar>
</template>

<script>
export default {
  props: {
    username: String,
    src: String,
  },
};
</script>
