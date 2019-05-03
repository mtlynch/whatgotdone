<template>
  <div class="submit">
    <p>What did you do this week?</p>
    <template v-if="!submitSucceeded">
      <form @submit.prevent="handleSubmit">
        <div class="form-group">
          <label for="date">Date</label>
          <input
            type="text"
            class="form-control"
            v-model="date"
            name="date"
            placeholder="2019-01-11"
          >
        </div>
        <textarea class="form-control" v-model="entryContent" name="markdown" rows="5"></textarea>
        <button type="submit" class="btn btn-primary">Submit</button>
      </form>
    </template>
  </div>
</template>

<script>
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
  },
  watch: {
    date: function(newDate) {
      if (
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
