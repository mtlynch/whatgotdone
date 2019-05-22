<template>
  <div class="home">
    <div class="jumbotron">
      <h1 class="display-3">What did you get done this week?</h1>
      <p class="lead">What Got Done is the simple, easy way to share progress with your teammates.</p>
      <hr class="my-4">
      <p>End each week by reflecting on what you accomplished.</p>
      <ul>
        <li>Keep your teammates informed</li>
        <li>Strengthen your motivation</li>
        <li>Put an end to boring status meetings</li>
      </ul>

      <div class="jumbotron-footer">
        <b-button
          class="btn btn-primary btn-lg"
          variant="outline-primary"
          to="/submit"
        >Start&nbsp;now</b-button>
        <b-button class="btn btn-lg btn-secondary" to="/recent">View&nbsp;recent&nbsp;entries</b-button>
      </div>
    </div>

    <h2>"Where did the whole week go?"</h2>

    <p>Have you ever ended a week wondering where all the time went? Take a step back and remember what got done.</p>

    <p>What Got Done is a weekly ritual to remind yourself of your accomplishments.</p>

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

    <h2>Share Updates Your Way</h2>

    <p>To-do lists are great for managing what you will do, but they do a poor job of communicating what you accomplished.</p>

    <p>What Got Done lets you write updates your way, highlighting what's important instead of forcing your teammates to dig through a list of checkboxes.</p>

    <figure>
      <a href="/images/example-freeform.jpg">
        <img id="example-freeform" src="/images/example-freeform.jpg">
      </a>
      <figcaption>Write updates any way that you choose.</figcaption>
    </figure>

    <h2>Contact Us</h2>
    <p>
      If you have questions or feature requests for What Got Done, we'd love to speak to you. Send an email to
      <a
        href="mailto:contact@whatgotdone.com"
      >contact@whatgotdone.com</a>.
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

.jumbotron {
  padding-bottom: 40px;
}

.jumbotron-footer {
  margin-top: 40px;
  text-align: center;
}

@media screen and (min-width: 768px) {
  .jumbotron-footer {
    text-align: right;
  }
}

.btn {
  margin-top: 20px;
}

@media screen and (min-width: 768px) {
  .btn {
    margin-top: 0px;
  }
}

.btn-primary {
  border: 1px solid rgba(0, 0, 0, 0.45);
  background: rgb(21, 255, 177);
  color: rgb(5, 70, 9) !important;
}

@media screen and (min-width: 768px) {
  .btn {
    margin-right: 20px;
  }
}

.btn-primary:hover {
  border: 1px solid rgb(4, 61, 23);
  background: rgb(180, 250, 226);
  color: rgb(5, 70, 9) !important;
}

a.btn-secondary {
  background: rgb(115, 116, 199);
  border: 1px solid rgba(0, 0, 0, 0.45);
  color: rgb(41, 2, 49) !important;
  margin-right: 0px;
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

#example-freeform {
  max-width: 470px;
  border: 1px solid black;
}
</style>
