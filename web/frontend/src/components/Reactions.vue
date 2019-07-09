<template>
  <div class="reactions">
    <div class="reaction-buttons">
      <b-button
        v-for="symbol in reactionSymbols"
        v-bind:key="symbol"
        :variant="buttonVariant(symbol)"
        @click="handleReaction(symbol)"
      >{{ symbol }}</b-button>
    </div>
    <div class="reaction" v-for="reaction in reactions" v-bind:key="reaction.key">
      <p>
        <Username :username="reaction.username" />&nbsp;reacted with a
        <span class="reaction-symbol">{{ reaction.reaction }}</span>
      </p>
    </div>
  </div>
</template>

<script>
import Username from "./Username.vue";
export default {
  name: "Reactions",
  props: {
    entryAuthor: String,
    entryDate: String
  },
  components: {
    Username
  },
  data() {
    return {
      reactions: [],
      reactionSymbols: ["👍", "🎉", "🙁"],
      selectedReaction: ""
    };
  },
  methods: {
    clear: function() {
      this.reactions = [];
      this.selectedReaction = "";
    },
    loadReactions: function() {
      const reactions = [];
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/reactions/entry/${this.entryAuthor}/${this.entryDate}`;
      this.$http
        .get(url)
        .then(result => {
          for (const reaction of result.data) {
            reactions.push({
              key: reaction.username,
              username: reaction.username,
              timestamp: new Date(reaction.timestamp),
              reaction: reaction.symbol
            });
            if (reaction.username == this.loggedInUsername) {
              this.selectedReaction = reaction.symbol;
            }
          }
          // Sort from newest to oldest.
          reactions.sort((a, b) => b.timestamp - a.timestamp);
          this.reactions = reactions;
        })
        .catch(error => {
          // Ignore error for reactions, as they're non-essential.
        });
    },
    reloadReactions: function() {
      this.clear();
      this.loadReactions();
    },
    handleReaction: function(reactionSymbol) {
      if (!this.loggedInUsername) {
        this.$router.push("/login");
        return;
      }
      if (this.selectedReaction == reactionSymbol) {
        this.selectedReaction = "";
      } else {
        this.selectedReaction = reactionSymbol;
      }
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/reactions/entry/${this.entryAuthor}/${this.entryDate}`;
      this.$http.post(
        url,
        {
          reactionSymbol: this.selectedReaction
        },
        { withCredentials: true }
      );
      this.updateReactions();
    },
    updateReactions: function() {
      const newReactions = [];
      if (this.selectedReaction) {
        newReactions.push({
          key: this.loggedInUsername,
          username: this.loggedInUsername,
          timestamp: new Date(),
          reaction: this.selectedReaction
        });
      }
      for (const reaction of this.reactions) {
        if (reaction.username == this.loggedInUsername) {
          continue;
        }
        newReactions.push(reaction);
      }
      this.reactions = newReactions;
    },
    buttonVariant: function(reactionSymbol) {
      if (this.selectedReaction == reactionSymbol) {
        return "light";
      } else {
        return "secondary";
      }
    }
  },
  computed: {
    loggedInUsername: function() {
      return this.$store.state.username;
    }
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
    }
  }
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