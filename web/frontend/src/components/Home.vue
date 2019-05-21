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

    <h2>"Where did the whole week go?"</h2>

    <p>Have you ever ended a week wondering where all the time went? Take a step back and remember what got done.</p>

    <p>Make What Got Done your weekly ritual to remind yourself of what you accomplished.</p>

    <figure>
      <a href="/images/typing-update.gif">
        <img id="typing-update" src="/images/typing-update.gif">
      </a>
      <figcaption>Summarize your accomplishments at the end of each week.</figcaption>
    </figure>

    <h2>Share progress with your teammates</h2>

    <p>When you're done, share a clean summary of your week of work with your teammates.</p>

    <figure>
      <a href="/images/example-rendered.jpg">
        <img id="example-rendered" src="/images/example-rendered.jpg">
      </a>
      <figcaption>Share a summary of your work with your teammates.</figcaption>
    </figure>

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
  font-size: 48px;
  margin-top: 50px;
}

@media screen and (min-width: 768px) {
  h1 {
    font-size: 72px;
  }
}

h2 {
  margin-top: 100px;
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
  margin-bottom: 30px;
}

@media screen and (min-width: 768px) {
  .btn-primary {
    margin-top: 0px;
    margin-bottom: 0px;
  }
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

figure {
  margin-top: 60px;
}

@media screen and (min-width: 768px) {
  figure {
    margin-left: 45px;
  }
}

figure img {
  width: 100%;
}

figcaption {
  font-size: small;
}

#typing-update {
  max-width: 685px;
}

#example-rendered {
  max-width: 694px;
  border: 1px solid white;
}
</style>
