<template>
  <div class="submit">
    <form @submit.prevent="handleSubmit">
      <h2>Select a week</h2>
      <b-form-select v-model="date" class="mb-3">
        <option :value="lastFriday">Week ending {{ lastFriday | moment("dddd, LL") }}</option>
        <option :value="thisFriday" selected>Week ending {{ thisFriday | moment("dddd, LL") }}</option>
      </b-form-select>
      <h2>What did you do this week?</h2>
      <p>
        (You can use
        <a href="https://www.markdownguide.org/cheat-sheet/">Markdown</a>)
      </p>
      <textarea-autosize
        class="form-control"
        v-model="entryContent"
        name="markdown"
        :min-height="250"
        :max-height="650"
      ></textarea-autosize>
      <button type="submit" class="btn btn-primary">Submit</button>
    </form>
  </div>
</template>

<script>
import Vue from "vue";
import VueTextareaAutosize from "vue-textarea-autosize";
import moment from "moment";

Vue.use(VueTextareaAutosize);

export default {
  name: "Submit",
  data() {
    return {
      username: "",
      date: "",
      entryContent: ""
    };
  },
  methods: {
    loadUsername() {
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/user/me`;
      this.$http
        .get(url, { withCredentials: true })
        .then(result => {
          this.username = result.data.username;
        })
        .catch(error => {
          if (error.response.status == 403) {
            this.$router.push("/login");
          }
        });
    },
    loadEntryContent() {
      if (this.date.length == 0 || this.username.length == 0) {
        return;
      }
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/entry/${
        this.username
      }/${this.date}`;
      this.$http
        .get(url)
        .then(result => {
          this.entryContent = result.data.markdown;
        })
        .catch(error => {
          if (error.response.status == 404) {
            this.entryContent = "";
          }
        });
    },
    handleSubmit() {
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/submit`;
      this.$http
        .post(
          url,
          {
            date: this.date,
            entryContent: this.entryContent
          },
          { withCredentials: true }
        )
        .then(result => {
          if (result.data.ok) {
            this.$router.push(result.data.path);
          }
        });
    }
  },
  created() {
    this.loadUsername();
    this.date = this.thisFriday;
  },
  watch: {
    date: function() {
      this.loadEntryContent();
    },
    username: function() {
      this.loadEntryContent();
    }
  },
  computed: {
    thisFriday: function() {
      return moment()
        .isoWeekday("Friday")
        .format("YYYY-MM-DD");
    },
    lastFriday: function() {
      const daysInWeek = 7;
      return moment(this.thisFriday)
        .subtract(daysInWeek, "days")
        .format("YYYY-MM-DD");
    }
  }
};
</script>

<style scoped>
.submit {
  font-size: 11pt;
}

button {
  margin-top: 25px;
}
</style>
