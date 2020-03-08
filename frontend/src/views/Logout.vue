<template>
  <div>
    <p>Signing out...</p>
  </div>
</template>

<script>
import {logout} from '@/controllers/Logout.js';

export default {
  name: 'Logout',
  created() {
    logout().then(() => {
      window.location.href = '/';
    });
  },
};
</script>
