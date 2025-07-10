<template>
  <div>
    <h1>Forwarding Address</h1>
    <p>
      Due to the
      <router-link to="/shutdown-notice"
        >shutdown of whatgotdone.com</router-link
      >, you can redirect your profile and posts to a new URL.
    </p>
    <p>
      Forwarding might be useful if you're hosting What Got Done or an
      equivalent service at another URL.
    </p>
    <p>
      When you set a forwarding address, visitors who access your profile will
      be redirected to your new location. For example:
    </p>
    <ul>
      <li>
        <code>{{ currentHostname }}/{{ loggedInUsername }}</code> will redirect
        to <code>https://new.example.com/{{ loggedInUsername }}</code>
      </li>
      <li>
        <code>{{ currentHostname }}/{{ loggedInUsername }}/2025-07-11</code>
        will redirect to
        <code>https://new.example.com/{{ loggedInUsername }}/2025-07-11</code>
      </li>
    </ul>

    <div v-if="currentForwardingUrl" class="mt-4">
      <div class="alert alert-info">
        <strong>Current forwarding address:</strong>
        <br />
        <code>{{ currentForwardingUrl }}</code>
      </div>
      <p>
        Visitors to your What Got Done profile are currently being redirected to
        this address.
      </p>
      <b-button variant="danger" v-on:click="onDelete" :disabled="isLoading">
        {{ isLoading ? 'Deleting...' : 'Delete Forwarding Address' }}
      </b-button>
    </div>

    <div v-else class="mt-4">
      <b-form @submit.prevent="onSave">
        <b-form-group
          label="Base URL of your new server:"
          label-for="forwarding-url"
          description="Enter the base URL where your content is hosted (e.g., https://new-location.example.com)"
        >
          <b-form-input
            id="forwarding-url"
            v-model="newForwardingUrl"
            type="url"
            placeholder="https://new.example.com"
            required
            :disabled="isLoading"
          ></b-form-input>
        </b-form-group>
        <b-button
          type="submit"
          variant="primary"
          :disabled="isLoading || !isValidUrl"
        >
          {{ isLoading ? 'Saving...' : 'Save Forwarding Address' }}
        </b-button>
      </b-form>
    </div>

    <div v-if="errorMessage" class="alert alert-danger mt-3">
      {{ errorMessage }}
    </div>
  </div>
</template>

<script>
import {
  getForwardingAddress,
  saveForwardingAddress,
  deleteForwardingAddress,
} from '@/controllers/ForwardingAddress.js';

export default {
  name: 'ForwardingAddressPage',
  data() {
    return {
      currentForwardingUrl: null,
      newForwardingUrl: '',
      isLoading: false,
      errorMessage: '',
    };
  },
  computed: {
    loggedInUsername: function () {
      return this.$store.state.username;
    },
    currentHostname: function () {
      return `${window.location.protocol}//${window.location.host}`;
    },
    isValidUrl: function () {
      if (!this.newForwardingUrl) return false;
      try {
        const url = new URL(this.newForwardingUrl);
        return url.protocol === 'http:' || url.protocol === 'https:';
      } catch {
        return false;
      }
    },
  },
  methods: {
    loadForwardingAddress: function () {
      this.isLoading = true;
      this.errorMessage = '';
      getForwardingAddress()
        .then((data) => {
          this.currentForwardingUrl = data.forwardingUrl || null;
        })
        .catch((error) => {
          this.errorMessage = `Failed to load forwarding address: ${error}`;
        })
        .finally(() => {
          this.isLoading = false;
        });
    },
    onSave: function () {
      if (!this.isValidUrl) return;

      this.isLoading = true;
      this.errorMessage = '';
      saveForwardingAddress(this.newForwardingUrl)
        .then(() => {
          this.currentForwardingUrl = this.newForwardingUrl;
          this.newForwardingUrl = '';
        })
        .catch((error) => {
          this.errorMessage = `Failed to save forwarding address: ${error}`;
        })
        .finally(() => {
          this.isLoading = false;
        });
    },
    onDelete: function () {
      this.isLoading = true;
      this.errorMessage = '';
      deleteForwardingAddress()
        .then(() => {
          this.currentForwardingUrl = null;
        })
        .catch((error) => {
          this.errorMessage = `Failed to delete forwarding address: ${error}`;
        })
        .finally(() => {
          this.isLoading = false;
        });
    },
  },
  created() {
    // Redirect to login if not authenticated
    if (!this.loggedInUsername) {
      this.$router.push('/login');
      return;
    }
    this.loadForwardingAddress();
  },
};
</script>

<style scoped>
* {
  text-align: left;
}

h1 {
  margin-bottom: 30px;
}

code {
  background-color: #f8f9fa;
  padding: 2px 4px;
  border-radius: 3px;
  font-family: 'Courier New', monospace;
}

.alert {
  border-radius: 6px;
}
</style>
