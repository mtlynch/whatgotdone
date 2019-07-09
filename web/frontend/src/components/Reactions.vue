<template>
  <div class="reactions">
    <div class="reactionButtons">
      <b-button
        v-for="symbol in reactionSymbols"
        v-bind:key="symbol"
        :variant="buttonVariant(symbol)"
        @click="sendReaction(symbol)"
      >{{ symbol }}</b-button>
    </div>
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
      reactionSymbols: ["👍", "🎉", "🙁"],
      selectedReaction: ""
    };
  },
  methods: {
    loadReactions: function() {
      const reactions = [];
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/reactions/entry/${this.username}/${this.date}`;
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
          reactions.sort((a, b) => b.timestamp - a.timestamp);
          this.reactions = reactions;
        })
        .catch(error => {
          // Ignore error for reactions, as they're non-essential.
        });
    },
    sendReaction: function(reactionSymbol) {
      if (this.selectedReaction == reactionSymbol) {
        this.selectedReaction = "";
      } else {
        this.selectedReaction = reactionSymbol;
      }
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/reactions/entry/${this.$route.params.username}/${this.$route.params.date}`;
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
        console.log("Adding reaction to end");
        newReactions.push({
          key: this.loggedInUsername,
          username: this.loggedInUsername,
          timestamp: new Date(),
          reaction: this.selectedReaction
        });
        console.log(this.reactions);
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
  }
};
</script>

<style scoped>
@media screen and (min-width: 768px) {
  .reactions {
    text-align: left;
  }
}

.reactionButtons {
  margin-top: 40px;
  margin-bottom: 25px;
  clear: both;
}

@media screen and (min-width: 768px) {
  .reactionButtons {
    margin-top: 0px;
    clear: none;
  }
}

.btn {
  margin-right: 5px;
  font-size: 24pt;
}

@media screen and (min-width: 768px) {
  .btn {
    font-size: 12pt;
  }
}
</style>