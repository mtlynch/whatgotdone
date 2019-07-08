<template>
  <div class="reactions">
    <div class="reaction" v-for="reaction in reactions" v-bind:key="reaction.key">
      <p>
        <Username :username="reaction.username" />
        reacted with a {{ reaction.reaction }}
      </p>
    </div>
  </div>
</template>

<script>
import Username from "./Username.vue";
export default {
  name: "Reactions",
  props: {
    username: String,
    date: String
  },
  components: {
    Username
  },
  data() {
    return {
      reactions: [],
      backendError: null
    };
  },
  methods: {
    loadReactions: function() {
      this.reactions = [];
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/reactions/entry/${this.username}/${this.date}`;
      this.$http
        .get(url)
        .then(result => {
          for (const reaction of result.data) {
            this.reactions.push({
              key: reaction.username,
              username: reaction.username,
              timestamp: new Date(reaction.timestamp),
              reaction: reaction.reaction
            });
          }
          this.reactions.sort((a, b) => a.timestamp - b.timestamp);
        })
        .catch(error => {
          this.backendError = error;
        });
    }
  },
  created() {
    this.loadReactions();
  }
};
</script>

<style scoped>
</style>