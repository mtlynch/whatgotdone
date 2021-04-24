<template>
  <div class="home">
    <div class="jumbotron">
      <h1 class="display-3">What did you get done this week?</h1>
      <p class="lead">
        What Got Done is the simple, easy way to share progress with your
        teammates.
      </p>
      <hr class="my-4" />
      <p>End each week by reflecting on what you accomplished.</p>
      <ul>
        <li>Keep your teammates informed</li>
        <li>Strengthen your motivation</li>
        <li>Put an end to boring status meetings</li>
      </ul>

      <div class="jumbotron-footer">
        <b-button
          class="btn btn-lg btn-start-now"
          variant="primary"
          :to="'/entry/edit/' + this.thisFriday"
          >Start&nbsp;now</b-button
        >
        <b-button class="btn btn-lg btn-success btn-view-recent" to="/recent"
          >View&nbsp;recent&nbsp;entries</b-button
        >
      </div>
    </div>

    <div class="feature">
      <h2>"Where did the whole week go?"</h2>

      <div class="text">
        <p>
          Have you ever ended a week wondering where all the time went? Take a
          step back and remember what got done.
        </p>

        <p>
          What Got Done is a weekly ritual to remind yourself of your
          accomplishments.
        </p>
      </div>
      <img src="/images/typing-update.gif" />
    </div>

    <div class="feature">
      <h2>Share progress with your teammates</h2>

      <div class="text">
        <p>
          When you're done, share a clean summary of your week of work with your
          teammates.
        </p>
      </div>
      <img class="border-white" src="/images/example-rendered.jpg" />
    </div>

    <div class="feature">
      <h2>Share updates your way</h2>

      <div class="text">
        <p>
          To-do lists are great for managing what you will do, but they do a
          poor job of communicating what you accomplished.
        </p>

        <p>
          What Got Done lets you write updates your way, highlighting what's
          important instead of forcing your teammates to dig through a list of
          checkboxes.
        </p>
      </div>
      <img class="border-black" src="/images/example-freeform.jpg" />
    </div>

    <div>
      <h2>Who made this?</h2>
      <p>
        What Got Done is a free, open-source project by
        <a href="https://mtlynch.io">Michael Lynch</a>.
      </p>
      <p>
        Anyone can contribute features or bugfix through What Got Done's
        <a href="https://github.com/mtlynch/whatgotdone">public Github</a>.
      </p>
    </div>
  </div>
</template>

<script>
import {thisFriday} from '@/controllers/EntryDates.js';

export default {
  name: 'Home',
  data() {
    return {
      thisFriday: thisFriday(),
    };
  },
  computed: {
    username() {
      return this.$store.state.username;
    },
  },
  created() {
    if (this.$route.path === '/' && this.username) {
      this.$router.push('/entry/edit/' + this.thisFriday);
      return;
    }
  },
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
  margin-bottom: 100px;
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

.feature {
  display: grid;
  grid-template-columns: 1fr;
  grid-template-rows: auto 1fr auto;
  grid-gap: 1rem;
  margin-bottom: 100px;
}

.feature h2 {
  grid-row: 1;
  grid-column: 1;
  margin-top: 0;
}

.feature img {
  grid-row: 2;
  grid-column: 1;
  max-width: 100%;
}

.feature .text {
  grid-row: 3 / 4;
  grid-column: 1 / 2;
}

@media screen and (min-width: 768px) {
  .feature {
    grid-template-columns: 1fr 500px;
    grid-template-rows: auto 1fr;
  }

  .feature h2 {
    grid-row: 1;
    grid-column: 1;
  }

  .feature img {
    grid-row: 1 / 3;
    grid-column: 2;
    max-width: 100%;
  }

  .feature .text {
    grid-row: 2;
    grid-column: 1;
    padding-right: 50px;
  }
}

.border-white {
  border: 1px solid white;
}

.border-black {
  border: 1px solid black;
}
</style>
