<template>
  <div>
    <p>Signing out...</p>
  </div>
</template>

<script>
import {logout} from '@/controllers/Logout.js';

export default {
  name: 'LogoutPage',
  created() {
    logout().then(() => {
      window.location.href = '/';
    });
  },
};
</script>
