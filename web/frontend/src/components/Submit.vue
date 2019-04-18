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
      date: "",
      entryContent: "",
      submitSucceeded: null
    };
  },
  watch: {
    date: function(newDate) {
      if (newDate.length != 10) {
        return;
      }
      // Hard-code the username until we support real accounts.
      const username = "michael";
      const url = `${
        process.env.VUE_APP_BACKEND_URL
      }/api/entry/${username}/${newDate}`;
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
