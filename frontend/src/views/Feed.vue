<template>
  <div class="feed">
    <h1>Feed</h1>
    <h2>TODO</h2>

    <p>
      Make a view like <router-link to="/recent">Recent</router-link>, except
      with only followed users.
    </p>

    <h2>Following</h2>
    <p>{{ following }}</p>
  </div>
</template>

<script>
export default {
  name: 'Feed',
  computed: {
    following: function() {
      return this.$store.state.following.join(', ');
    },
  },
};
</script>
