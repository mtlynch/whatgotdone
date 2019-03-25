<template>
  <div class="journal">
    <div class="journalEntry">
      <div class="date">{{ entry.date }}</div>
      <div class="journalBody">{{ entry.markdown }}</div>
    </div>
  </div>
</template>

<script>
export default {
  name: "Journal",
  props: {
    entry: Object
  }
};
</script>

<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
