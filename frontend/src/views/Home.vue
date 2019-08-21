<template>
  <div class="home">
    <div class="jumbotron">
      <h1 class="display-3">What did you get done this week?</h1>
      <p class="lead">What Got Done is the simple, easy way to share progress with your teammates.</p>
      <hr class="my-4" />
      <p>End each week by reflecting on what you accomplished.</p>
      <ul>
        <li>Keep your teammates informed</li>
        <li>Strengthen your motivation</li>
        <li>Put an end to boring status meetings</li>
      </ul>

      <div class="jumbotron-footer">
        <b-button class="btn btn-lg btn-start-now" variant="primary" to="/submit">Start&nbsp;now</b-button>
        <b-button
          class="btn btn-lg btn-success btn-view-recent"
          to="/recent"
        >View&nbsp;recent&nbsp;entries</b-button>
      </div>
    </div>

    <h2>"Where did the whole week go?"</h2>

    <div class="flex-grid">
      <div class="col">
        <p>Have you ever ended a week wondering where all the time went? Take a step back and remember what got done.</p>

        <p>What Got Done is a weekly ritual to remind yourself of your accomplishments.</p>
      </div>
      <div class="col">
        <img id="typing-update" src="/images/typing-update.gif" />
      </div>
    </div>

    <h2>Share progress with your teammates</h2>

    <div class="flex-grid">
      <div class="col">
        <p>When you're done, share a clean summary of your week of work with your teammates.</p>
      </div>
      <div class="col">
        <img id="example-rendered" src="/images/example-rendered.jpg" />
      </div>
    </div>

    <h2>Share Updates Your Way</h2>

    <div class="flex-grid">
      <div class="col">
        <p>To-do lists are great for managing what you will do, but they do a poor job of communicating what you accomplished.</p>

        <p>What Got Done lets you write updates your way, highlighting what's important instead of forcing your teammates to dig through a list of checkboxes.</p>
      </div>
      <div class="col">
        <img id="example-freeform" src="/images/example-freeform.jpg" />
      </div>
    </div>

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
  margin-top: 80px;
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
  margin-top: 60px;
  text-align: center;
}

@media screen and (min-width: 768px) {
  .jumbotron-footer {
    margin-top: 40px;
    text-align: right;
  }
}

@media screen and (min-width: 768px) {
  .btn {
    margin-right: 20px;
  }
}

.btn-start-now {
  margin-right: 15px;
}

.btn-view-recent {
  margin-right: 0px;
}

.flex-grid {
  display: flex;
}

@media (max-width: 768px) {
  .flex-grid {
    display: block;
  }
}

.col {
  flex: 1;
}

.col img {
  margin-top: 28px;
  max-width: 100%;
}

#example-rendered {
  border: 1px solid white;
}

#example-freeform {
  border: 1px solid black;
}
</style>
