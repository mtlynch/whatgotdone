<template>
  <div class="reactions">
    <div class="reaction-buttons">
      <b-button
        v-for="symbol in reactionSymbols"
        v-bind:key="symbol"
        :variant="buttonVariant(symbol)"
        @click="handleReaction(symbol)"
        >{{ symbol }}</b-button
      >
    </div>
    <div
      class="reaction"
      v-for="reaction in reactions"
      v-bind:key="reaction.key"
    >
      <p>
        <Username :username="reaction.username" />&nbsp;reacted with a
        <span class="reaction-symbol">{{ reaction.reaction }}</span>
      </p>
    </div>
  </div>
</template>

<script>
import Username from './Username.vue';
import getCsrfToken from '../controllers/CsrfToken.js';

export default {
  name: 'Reactions',
  props: {
    entryAuthor: String,
    entryDate: String,
  },
  components: {
    Username,
  },
  data() {
    return {
      reactions: [],
      reactionSymbols: ['👍', '🎉', '🙁'],
      selectedReaction: '',
    };
  },
  methods: {
    clear: function() {
      this.reactions = [];
      this.selectedReaction = '';
    },
    loadReactions: function() {
      if (!this.entryAuthor || !this.entryDate) {
        return;
      }
      const reactions = [];
      let newSelectedReaction = '';
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/reactions/entry/${this.entryAuthor}/${this.entryDate}`;
      this.$http
        .get(url)
        .then(result => {
          for (const reaction of result.data) {
            if (reaction.username == this.loggedInUsername) {
              if (this.selectedReaction == '') {
                newSelectedReaction = reaction.symbol;
              } else {
                // Don't overwrite the local reaction symbol if the user
                // clicked a reaction before the request finished.
                continue;
              }
            }

            reactions.push({
              key: reaction.username,
              username: reaction.username,
              timestamp: new Date(reaction.timestamp),
              reaction: reaction.symbol,
            });
          }
          if (this.selectedReaction != '') {
            reactions.push({
              key: this.loggedInUsername,
              username: this.loggedInUsername,
              timestamp: new Date(),
              reaction: this.selectedReaction,
            });
          }
          // Sort from oldest to newest.
          reactions.sort((a, b) => a.timestamp - b.timestamp);
          this.reactions = reactions;
          if (this.selectedReaction == '') {
            this.selectedReaction = newSelectedReaction;
          }
        })
        .catch(() => {
          // Ignore error for reactions, as they're non-essential.
        });
    },
    reloadReactions: function() {
      this.clear();
      this.loadReactions();
    },
    handleReaction: function(reactionSymbol) {
      if (!this.loggedInUsername) {
        this.$router.push('/login');
        return;
      }
      if (this.selectedReaction == reactionSymbol) {
        this.selectedReaction = '';
      } else {
        this.selectedReaction = reactionSymbol;
      }
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/reactions/entry/${this.entryAuthor}/${this.entryDate}`;
      this.$http.post(
        url,
        {
          reactionSymbol: this.selectedReaction,
        },
        {withCredentials: true, headers: {'X-CSRF-Token': getCsrfToken()}}
      );
      this.updateReactions();
    },
    updateReactions: function() {
      const newReactions = [];
      for (const reaction of this.reactions) {
        if (reaction.username == this.loggedInUsername) {
          continue;
        }
        newReactions.push(reaction);
      }
      if (this.selectedReaction != '') {
        newReactions.push({
          key: this.loggedInUsername,
          username: this.loggedInUsername,
          timestamp: new Date(),
          reaction: this.selectedReaction,
        });
      }
      this.reactions = newReactions;
    },
    buttonVariant: function(reactionSymbol) {
      if (this.selectedReaction == reactionSymbol) {
        return 'light';
      } else {
        return 'secondary';
      }
    },
  },
  computed: {
    loggedInUsername: function() {
      return this.$store.state.username;
    },
  },
  created() {
    this.loadReactions();
  },
  watch: {
    entryAuthor: function() {
      this.reloadReactions();
    },
    entryDate: function() {
      this.reloadReactions();
    },
  },
};
</script>

<style scoped>
@media screen and (min-width: 768px) {
  .reactions {
    text-align: left;
  }
}

.reaction-buttons {
  margin-top: 40px;
  margin-bottom: 25px;
  clear: both;
}

@media screen and (min-width: 768px) {
  .reaction-buttons {
    margin-top: 0px;
    clear: none;
  }
}

.btn {
  margin-right: 12px;
  font-size: 24pt;
}

@media screen and (min-width: 768px) {
  .btn {
    margin-right: 8px;
    font-size: 12pt;
  }
}

.reaction-symbol {
  font-size: 18pt;
}

@media screen and (min-width: 768px) {
  .reaction-symbol {
    font-size: 14pt;
  }
}
</style>
