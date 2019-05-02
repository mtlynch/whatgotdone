<template>
  <div class="userkit">
    <p>Please wait, logging in...</p>
  </div>
</template>

<script>
export default {
  name: "Login",
  mounted() {
    let userKitScript = document.createElement("script");
    userKitScript.setAttribute("src", "https://widget.userkit.io/widget.js");
    userKitScript.setAttribute("data-app-id", "app_zXjb2TF8");
    userKitScript.setAttribute("data-show-on-load", "login");
    userKitScript.setAttribute("data-login-dismiss", "false");
    document.head.appendChild(userKitScript);
  }
};
</script>
