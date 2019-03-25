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
      journalEntries: [
        { id: 0, date: "2019-03-25", markdown: "Ate some crackers" },
        { id: 1, date: "2019-03-24", markdown: "Took a nap" },
        { id: 2, date: "2019-03-23", markdown: "Watched a movie" }
      ]
    };
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
