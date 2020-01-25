<template>
  <b-navbar
    toggleable="md"
    class="navbar navbar-expand-lg navbar-light bg-light"
  >
    <b-navbar-brand
      v-b-toggle.nav-collapse
      class="navbar-brand"
      :to="logoLinkTarget"
      >What Got Done</b-navbar-brand
    >

    <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

    <b-collapse id="nav-collapse" is-nav>
      <b-navbar-nav>
        <b-nav-item to="/about">About</b-nav-item>
        <b-nav-item to="/recent">Recent</b-nav-item>
        <b-nav-item href="https://github.com/mtlynch/whatgotdone"
          >Contribute</b-nav-item
        >
      </b-navbar-nav>

      <!-- Right aligned nav items -->
      <b-navbar-nav class="ml-auto">
        <b-nav-item-dropdown
          v-if="username"
          text="Account"
          class="account-dropdown"
          right
        >
          <b-dropdown-item :to="'/' + username" class="profile-link"
            >Profile</b-dropdown-item
          >
          <b-dropdown-item to="/logout">Sign Out</b-dropdown-item>
        </b-nav-item-dropdown>
        <b-button
          class="post-update"
          variant="success"
          v-b-toggle.nav-collapse
          :to="editCurrentWeekUrl"
          v-if="!isOnEntryEditPage"
          >Post Update</b-button
        >
      </b-navbar-nav>
    </b-collapse>
  </b-navbar>
</template>

<script>
import {thisFriday} from '../controllers/EntryDates.js';

export default {
  name: 'NavigationBar',
  data() {
    return {
      editCurrentWeekUrl: '/entry/edit/' + thisFriday(),
    };
  },
  computed: {
    username() {
      return this.$store.state.username;
    },
    isOnEntryEditPage() {
      return this.$route.path.startsWith('/entry/edit/');
    },
    logoLinkTarget() {
      if (!this.username) {
        return '/about';
      }
      return this.editCurrentWeekUrl;
    },
  },
};
</script>

<style scoped>
.navbar-brand {
  margin: 10px 20px 10px 0px;
}
</style>
