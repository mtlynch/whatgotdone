<template>
  <div class="home">
    <div class="jumbotron">
      <h1 class="display-3">What did you get done this week?</h1>
      <p class="lead">What Got Done is a lightweight, weekly task reflections app.</p>
      <hr class="my-4">
      <p>End each week by highlighting what you accomplished.</p>
      <p>It's a great way to keep your teammates informed and to strengthen your motivation.</p>
      <p class="lead float-right">
        <b-button class="btn btn-primary btn-lg" variant="outline-primary" to="/submit">Start now</b-button>
        <b-button class="btn btn-lg btn-secondary" to="/recent">View recent entries</b-button>
      </p>
      <p>&nbsp;</p>
    </div>

    <h2>Interested?</h2>
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

a {
  color: red;
}

.btn-primary {
  margin-right: 20px;
}

a.btn-secondary {
  background: rgba(85, 58, 39, 0.685);
}
</style>
