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
  border: 1px solid rgb(8, 24, 8);
  background: rgb(21, 255, 177);
  color: rgb(5, 70, 9) !important;
  margin-right: 20px;
}

.btn-primary:hover {
  border: 1px solid rgb(4, 61, 23);
  background: rgb(180, 250, 226);
  color: rgb(5, 70, 9) !important;
}

a.btn-secondary {
  background: rgb(115, 116, 199);
  border: 1px solid rgb(0, 0, 0);
  color: rgb(41, 2, 49) !important;
}

a.btn-secondary:hover {
  background: rgb(163, 164, 233);
  border: 1px solid rgb(0, 0, 0);
  color: rgb(41, 2, 49) !important;
}
</style>
