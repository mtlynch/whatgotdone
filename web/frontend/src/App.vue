<template>
  <div id="app" class="container">
    <router-view></router-view>
  </div>
</template>

<script>
export default {
  name: "app"
};
</script>

<style>
#app {
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  margin-top: 60px;
}

#app a {
  color: #4b7ccc;
}

#app a.page-link {
  color: white;
}

#app .page-link {
  border: 1px solid rgb(124, 133, 145);
}
</style>
