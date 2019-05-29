<template>
  <div>
    <NavigationBar/>
    <div id="app" class="container">
      <router-view></router-view>
    </div>
    <Footer/>
    <Initializer/>
  </div>
</template>

<script>
import Initializer from "./components/Initializer";
import Footer from "./components/Footer";
import NavigationBar from "./components/NavigationBar";
export default {
  name: "app",
  components: {
    Initializer,
    Footer,
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

/* TODO: Remove the #app part of the selector */
#app a {
  color: rgb(101, 168, 255);
}

#app a.btn {
  color: white;
}

/* TODO: Move these to the view entry component */
#app a.page-link {
  color: white;
}

#app .page-link {
  border: 1px solid rgb(124, 133, 145);
}
</style>
