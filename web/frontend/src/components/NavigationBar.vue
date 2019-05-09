<template>
  <b-navbar class="navbar navbar-expand-lg navbar-light bg-light">
    <b-navbar-brand class="navbar-brand" to="/">What Got Done</b-navbar-brand>

    <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

    <b-collapse id="nav-collapse" is-nav>
      <b-navbar-nav>
        <b-nav-item to="/">About</b-nav-item>
        <b-nav-item to="/pro">Pro</b-nav-item>
      </b-navbar-nav>

      <!-- Right aligned nav items -->
      <b-navbar-nav class="ml-auto">
        <b-nav-item right to="/submit">
          <button class="btn btn-success my-2 my-sm-0" type="submit">Post Update</button>
        </b-nav-item>
        <b-nav-item-dropdown right>
          <template slot="button-content">
            <em>Account</em>
          </template>
          <b-dropdown-item href="/me">Profile</b-dropdown-item>
          <b-dropdown-item href="/logout">Sign Out</b-dropdown-item>
        </b-nav-item-dropdown>
      </b-navbar-nav>
    </b-collapse>
  </b-navbar>
</template>

<script>
export default {
  name: "NavigationBar"
};
</script>
