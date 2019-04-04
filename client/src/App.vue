<template>
  <div id="app">
    <Journal v-for="item in journalEntries" v-bind:entry="item" v-bind:key="item.id"/>
  </div>
</template>

<script>
import Journal from "./components/Journal.vue";

export default {
  name: "app",
  components: {
    Journal
  },
  data() {
    return {
      journalEntries: []
    };
  },
  created() {
    const url = `${process.env.VUE_APP_BACKEND_URL}/entries`;
    this.$http.get(url).then(result => {
      this.journalEntries = [];
      for (const entry of result.data) {
        this.journalEntries.push({
          id: entry.id,
          date: new Date(entry.date),
          markdown: entry.markdown
        });
      }
    });
  }
};
</script>

<style>
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
