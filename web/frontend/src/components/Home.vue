<template>
  <div class="home">
    <h1>What did you get done this week?</h1>
    <p>What Got Done is a weekly task journaling app.</p>
    <p>End each week by reflecting on what you accomplished. It's a great way to reflect on the week and share your progress with others.</p>
    <p>
      Check out some
      <router-link to="/recent">recent entries</router-link>.
    </p>
    <h1>Interested?</h1>
    <p>What Got Done is still a work in progress, and new user registration is not yet open.</p>
    <p>
      If you're interested in beta testing the site, please send an email to
      <a
        href="mailto:michael@mtlynch.io"
      >michael@mtlynch.io</a>.
    </p>
  </div>
</template>

<script>
export default {
  name: "Home"
};
</script>

<style scoped>
.home {
  text-align: left;
}

h1 {
  margin-top: 50px;
}

p {
  margin-top: 20px;
}
</style>
