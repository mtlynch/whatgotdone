<template>
  <b-navbar toggleable="md" class="navbar navbar-expand-lg navbar-light bg-light">
    <b-navbar-brand v-b-toggle.nav-collapse class="navbar-brand" to="/">What Got Done</b-navbar-brand>

    <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

    <b-collapse id="nav-collapse" is-nav>
      <b-navbar-nav>
        <b-nav-item to="/">About</b-nav-item>
        <b-nav-item to="/recent">Recent</b-nav-item>
        <b-nav-item to="/pro">Pro</b-nav-item>
      </b-navbar-nav>

      <!-- Right aligned nav items -->
      <b-navbar-nav class="ml-auto">
        <b-nav-item-dropdown text="Account" right>
          <!-- This must be href instead of to so that we get the proper server-side redirect -->
          <b-dropdown-item href="/me">Profile</b-dropdown-item>
          <b-dropdown-item to="/logout">Sign Out</b-dropdown-item>
        </b-nav-item-dropdown>
        <b-button variant="success" v-b-toggle.nav-collapse to="/submit">Post Update</b-button>
      </b-navbar-nav>
    </b-collapse>
  </b-navbar>
</template>

<script>
export default {
  name: "NavigationBar"
};
</script>

<style scoped>
.navbar-brand {
  margin: 10px 20px 10px 0px;
}
</style>
