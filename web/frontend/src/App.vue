<template>
  <div>
    <NavigationBar/>
    <div id="app" class="container">
      <router-view></router-view>
    </div>
  </div>
</template>

<script>
import NavigationBar from "./components/NavigationBar";
export default {
  name: "app",
  components: {
    NavigationBar
  }
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
  color: rgb(101, 168, 255);
}

#app a.btn {
  color: white;
}
</style>
