<template>
  <b-navbar
    toggleable="md"
    id="navbar"
    class="navbar navbar-expand-lg navbar-dark bg-dark"
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
        <b-nav-item v-if="isLoggedIn" to="/feed">Feed</b-nav-item>
        <b-nav-item to="/recent">Recent</b-nav-item>
        <b-nav-item to="/shutdown-notice" class="warning-link ml-3"
          >Shutdown Notice</b-nav-item
        >
      </b-navbar-nav>

      <!-- Right aligned nav items -->
      <b-navbar-nav class="ml-auto">
        <b-nav-item-dropdown v-if="isLoggedIn" text="Account" right>
          <b-dropdown-item :to="'/' + username">Profile</b-dropdown-item>
          <b-dropdown-item to="/preferences">Preferences</b-dropdown-item>
          <b-dropdown-item to="/export">Export</b-dropdown-item>
          <b-dropdown-item to="/forwarding-address"
            >Forwarding Address</b-dropdown-item
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
import {thisFriday} from '@/controllers/EntryDates.js';

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
    isLoggedIn() {
      if (this.username) {
        return true;
      }
      return false;
    },
    isOnEntryEditPage() {
      return this.$route.path.startsWith('/entry/edit/');
    },
    logoLinkTarget() {
      if (!this.username) {
        return '/';
      }
      return this.editCurrentWeekUrl;
    },
  },
};
</script>

<style scoped>
#navbar .navbar-brand {
  margin: 10px 20px 10px 0px;
}

#navbar .navbar {
  background: #2b3e50;
}

#navbar .navbar,
#navbar .navbar-brand,
#navbar .nav-link {
  color: white;
}

#navbar .navbar-toggler {
  border: none;
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.4);
}

#navbar .nav-link {
  padding: 8px;
  font-size: 16px;
}

@media screen and (min-width: 768px) {
  #navbar .nav-link {
    font-size: 14px;
  }
}

#navbar .nav-link:hover {
  cursor: pointer;
  background: rgba(255, 255, 255, 0.4);
  border-radius: 6px;
}

#navbar .dropdown-menu {
  border-radius: 6px;
  border: none;
  box-shadow: 0px 3px 6px rgba(0, 0, 0, 0.4);
}

#navbar .dropdown-item {
  font-size: 16px;
}

@media screen and (min-width: 768px) {
  #navbar .dropdown-item {
    font-size: 14px;
  }
}

#navbar .warning-link .nav-link {
  background-color: #f39c12 !important;
  color: #2c3e50 !important;
  font-weight: bold;
  border-radius: 6px;
}

#navbar .warning-link .nav-link:hover {
  background-color: #e67e22 !important;
  color: #2c3e50 !important;
}
</style>
