<template>
  <div class="submit">
    <h1>What got done this week?</h1>
    <form @submit.prevent="handleSubmit">
      <p>
        Enter your update for the week ending
        <span class="endDate">{{ date | moment("dddd, ll") }}</span>.
      </p>
      <textarea-autosize
        class="form-control journal-markdown"
        v-model="entryContent"
        name="markdown"
        :min-height="250"
        :max-height="650"
      ></textarea-autosize>
      <p>
        (You can use
        <a href="https://www.markdownguide.org/cheat-sheet/">Markdown</a>)
      </p>
      <div class="d-flex justify-content-end">
        <button
          @click.prevent="handleSaveDraft"
          class="btn btn-primary save-draft"
          :disabled="changesSaved"
        >{{ saveLabel }}</button>
        <button type="submit" class="btn btn-primary">Publish</button>
      </div>
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
      date: "",
      entryContent: "",
      changesSaved: true,
      saveLabel: "Save Draft"
    };
  },
  computed: {
    username() {
      return this.$store.state.username;
    }
  },
  methods: {
    loadEntryContent() {
      if (this.date.length == 0 || !this.username) {
        return;
      }
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/draft/${this.date}`;
      this.$http
        .get(url, { withCredentials: true })
        .then(result => {
          this.entryContent = result.data.markdown;
        })
        .catch(error => {
          if (error.response.status == 404) {
            this.entryContent = "";
          }
        });
    },
    handleSaveDraft() {
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/draft/${this.date}`;
      this.$http
        .post(
          url,
          {
            entryContent: this.entryContent
          },
          { withCredentials: true }
        )
        .then(result => {
          if (result.data.ok) {
            this.changesSaved = true;
            this.saveLabel = "Changes Saved";
          }
        })
        .catch(() => {
          this.changesSaved = false;
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
    },
    thisFriday() {
      const today = moment().isoWeekday();
      const friday = 5;

      if (today <= friday) {
        return moment().isoWeekday(friday);
      } else {
        return moment()
          .add(1, "weeks")
          .isoWeekday(friday);
      }
    },
    validateDate(d) {
      const m = moment(d);
      if (!m.isValid()) {
        return false;
      }
      const whatGotDoneCreationYear = 2019;
      if (m.year() < whatGotDoneCreationYear) {
        return false;
      }
      if (m > this.thisFriday()) {
        return false;
      }
      const friday = 5;
      if (m.isoWeekday() != friday) {
        return false;
      }
      return true;
    }
  },
  created() {
    if (!this.username) {
      this.$router.push("/login");
      return;
    }
    if (this.$route.params.date && this.validateDate(this.$route.params.date)) {
      this.date = this.$route.params.date;
    } else {
      this.date = this.thisFriday().format("YYYY-MM-DD");
    }
  },
  watch: {
    date: function() {
      this.loadEntryContent();
    },
    username: function() {
      this.loadEntryContent();
    },
    entryContent: function() {
      this.changesSaved = false;
      this.saveLabel = "Save Draft";
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

.save-draft {
  margin-right: 20px;
}
</style>
