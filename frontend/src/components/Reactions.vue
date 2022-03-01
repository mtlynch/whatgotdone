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
        <Avatar :username="reaction.username" class="avatar" />
        <UsernameLink :username="reaction.username" />&nbsp;reacted with a
        <span class="reaction-symbol">{{ reaction.reaction }}</span>
      </p>
    </div>
  </div>
</template>

<script>
import Avatar from '@/components/Avatar.vue';
import UsernameLink from '@/components/UsernameLink.vue';

import {
  getReactions,
  deleteReaction,
  setReaction,
} from '@/controllers/Reactions.js';

export default {
  name: 'Reactions',
  props: {
    entryAuthor: String,
    entryDate: String,
  },
  components: {
    Avatar,
    UsernameLink,
  },
  data() {
    return {
      reactions: [],
      reactionSymbols: ['👍', '🎉', '🙁'],
      selectedReaction: '',
    };
  },
  methods: {
    clear: function () {
      this.reactions = [];
      this.selectedReaction = '';
    },
    loadReactions: function () {
      if (!this.entryAuthor || !this.entryDate) {
        return;
      }
      const reactions = [];
      let newSelectedReaction = '';
      getReactions(this.entryAuthor, this.entryDate)
        .then((fetchedReactions) => {
          for (const reaction of fetchedReactions) {
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
    reloadReactions: function () {
      this.clear();
      this.loadReactions();
    },
    handleReaction: function (reactionSymbol) {
      if (!this.loggedInUsername) {
        this.$router.push('/login');
        return;
      }
      if (this.selectedReaction == reactionSymbol) {
        this.selectedReaction = '';
        deleteReaction(this.entryAuthor, this.entryDate);
      } else {
        this.selectedReaction = reactionSymbol;
        setReaction(this.entryAuthor, this.entryDate, this.selectedReaction);
      }
      this.updateReactions();
    },
    updateReactions: function () {
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
    buttonVariant: function (reactionSymbol) {
      if (this.selectedReaction == reactionSymbol) {
        return 'light';
      } else {
        return 'secondary';
      }
    },
  },
  computed: {
    loggedInUsername: function () {
      return this.$store.state.username;
    },
  },
  created() {
    this.loadReactions();
  },
  watch: {
    entryAuthor: function () {
      this.reloadReactions();
    },
    entryDate: function () {
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
  margin: 1rem 0;
  clear: both;
}

@media screen and (min-width: 768px) {
  .reaction-buttons {
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

.reaction .avatar {
  margin-right: 0.5rem;
}

.reaction .username {
  background: #2b3e50;
  border-radius: 6px;
  padding: 2px 8px 5px;
  color: white;
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
