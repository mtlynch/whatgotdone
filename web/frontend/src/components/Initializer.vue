<template>
  <div>
    <!-- This is a dummy component. We just want the created() event. -->
  </div>
</template>

<script>
import refreshRecent from "../controllers/Recent.js";
import updateLoginState from "../controllers/LoginState.js";

export default {
  name: "Initializer",
  created() {
    updateLoginState(5);
    refreshRecent();
  }
};
</script>