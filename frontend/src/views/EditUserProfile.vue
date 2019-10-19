<template>
  <div>
    <h1>Edit Profile</h1>

    <h2>About</h2>

    <textarea v-model="aboutMarkdown" />

    <h2>Contact</h2>

    <ul>
      <li>
        Twitter:
        <input type="text" :value="twitterUsername" />
      </li>
      <li>
        Email:
        <input type="text" :value="email" />
      </li>
    </ul>

    <b-button to="/michael" variant="primary" class="float-right">Save</b-button>
  </div>
</template>

<script>
export default {
  name: "EditUserProfile",
  data() {
    return {
      // TODO(mtlynch): Retrieve this from the server.
      aboutMarkdown:
        "Hi, I'm Michael, creator of What Got Done.\n\nI also blog at [mtlynch.io](https://mtlynch.io).",
      twitterUsername: "deliberatecoder",
      email: "michael@mtlynch.io",
      canEdit: true
    };
  }
};
</script>

<style scoped>
* {
  text-align: left;
}
</style>