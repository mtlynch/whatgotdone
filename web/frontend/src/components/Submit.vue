<template>
  <div class="submit">
    <h1>What got done this week?</h1>
    <form @submit.prevent="handleSubmit">
      <p>
        Enter your update for the week ending
        <span class="endDate">{{ date | moment("dddd, ll") }}</span>.
      </p>
      <textarea-autosize
        class="form-control"
        v-model="entryContent"
        name="markdown"
        :min-height="250"
        :max-height="650"
      ></textarea-autosize>
      <p>
        (You can use
        <a href="https://www.markdownguide.org/cheat-sheet/">Markdown</a>)
      </p>
      <button type="submit" class="btn btn-primary float-right">Submit</button>
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
    this.date = moment()
      .isoWeekday("Friday")
      .format("YYYY-MM-DD");
  },
  watch: {
    date: function() {
      this.loadEntryContent();
    },
    username: function() {
      this.loadEntryContent();
    }
  }
};
</script>

<style scoped>
.submit {
  text-align: left;
  font-size: 11pt;
}

span.endDate {
  color: rgb(255, 208, 56);
  font-weight: bold;
}
</style>
