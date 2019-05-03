<template>
  <div class="submit">
    <p>What did you do this week?</p>
    <template v-if="!submitSucceeded">
      <form @submit.prevent="handleSubmit">
        <b-form-select v-model="date" class="mb-3">
          <option :value="lastFriday">Week ending {{ lastFriday | moment("dddd, LL") }}</option>
          <option :value="thisFriday" selected>Week ending {{ thisFriday | moment("dddd, LL") }}</option>
        </b-form-select>
        <textarea class="form-control" v-model="entryContent" name="markdown" rows="5"></textarea>
        <button type="submit" class="btn btn-primary">Submit</button>
      </form>
    </template>
  </div>
</template>

<script>
import moment from "moment";

export default {
  name: "Submit",
  data() {
    return {
      username: "",
      date: "",
      entryContent: "",
      submitSucceeded: null
    };
  },
  created() {
    const url = `${process.env.VUE_APP_BACKEND_URL}/api/user/me`;
    this.$http.get(url).then(result => {
      this.username = result.data.username;
    });
    this.date = this.thisFriday;
  },
  watch: {
    date: function(newDate) {
      if (
        !newDate ||
        newDate.length != 10 ||
        this.username.length == 0 ||
        this.entryContent.length > 0
      ) {
        return;
      }
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/entry/${
        this.username
      }/${newDate}`;
      this.$http.get(url).then(result => {
        this.entryContent = result.data.markdown;
      });
    }
  },
  computed: {
    thisFriday: function() {
      const daysInWeek = 7;
      const daysToAdd =
        (moment().isoWeekday("Friday") - moment().isoWeekday()) % daysInWeek;
      return moment()
        .add(daysToAdd, "days")
        .format("YYYY-MM-DD");
    },
    lastFriday: function() {
      const daysInWeek = 7;
      return moment(this.thisFriday)
        .subtract(daysInWeek, "days")
        .format("YYYY-MM-DD");
    }
  },
  methods: {
    handleSubmit() {
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/submit`;
      this.$http
        .post(url, {
          date: this.date,
          entryContent: this.entryContent
        })
        .then(result => {
          if (result.ok) {
            this.submitSucceeded = true;
          }
        });
    }
  }
};
</script>

<style scoped>
.submit {
  font-size: 11pt;
}
</style>
