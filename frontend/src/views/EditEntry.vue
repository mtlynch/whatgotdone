<template>
  <div class="submit">
    <h1>What got done this week?</h1>
    <form @submit.prevent="handleSubmit">
      <p>
        Enter your update for the week ending
        <span class="end-date">{{ date | moment('dddd, ll') }}</span
        >.
      </p>
      <textarea-autosize
        class="form-control journal-markdown"
        v-model="entryContent"
        name="markdown"
        @input="debouncedSaveDraft"
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
        >
          {{ saveLabel }}
        </button>
        <button type="submit" class="btn btn-primary">Publish</button>
      </div>
    </form>
    <JournalPreview :markdown="entryContent" />
  </div>
</template>

<script>
import Vue from 'vue';
import VueTextareaAutosize from 'vue-textarea-autosize';
import _ from 'lodash';
import JournalPreview from '../components/JournalPreview.vue';
import getCsrfToken from '../controllers/CsrfToken.js';
import {isValidEntryDate, thisFriday} from '../controllers/EntryDates.js';

Vue.use(VueTextareaAutosize);

export default {
  name: 'EditEntry',
  components: {
    JournalPreview,
  },
  data() {
    return {
      date: '',
      entryContent: '',
      changesSaved: true,
      saveLabel: 'Save Draft',
    };
  },
  computed: {
    username() {
      return this.$store.state.username;
    },
  },
  methods: {
    loadEntryContent() {
      if (this.date.length == 0 || !this.username) {
        return;
      }
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/draft/${this.date}`;
      this.$http
        .get(url, {withCredentials: true})
        .then(result => {
          this.entryContent = result.data.markdown;
        })
        .catch(error => {
          if (error.response.status == 404) {
            this.entryContent = '';
          }
        });
    },
    handleSaveDraft() {
      this.saveLabel = 'Saving';
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/draft/${this.date}`;
      this.$http
        .post(
          url,
          {
            entryContent: this.entryContent,
          },
          {withCredentials: true, headers: {'X-CSRF-Token': getCsrfToken()}}
        )
        .then(result => {
          if (result.data.ok) {
            this.changesSaved = true;
            this.saveLabel = 'Changes Saved';
          }
        })
        .catch(() => {
          this.changesSaved = false;
        });
    },
    debouncedSaveDraft: _.debounce(function() {
      this.handleSaveDraft();
    }, 2500),
    handleSubmit() {
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/entry/${this.date}`;
      this.$http
        .post(
          url,
          {
            entryContent: this.entryContent,
          },
          {withCredentials: true, headers: {'X-CSRF-Token': getCsrfToken()}}
        )
        .then(result => {
          if (result.data.ok) {
            this.$router.push(result.data.path);
          }
        });
    },
  },
  created() {
    if (!this.username) {
      this.$router.replace('/login');
      return;
    }
    if (this.$route.params.date && isValidEntryDate(this.$route.params.date)) {
      this.date = this.$route.params.date;
    } else {
      this.date = thisFriday();
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
      this.saveLabel = 'Save Draft';
    },
    $route(to, from) {
      if (to.params.date != from.params.date) {
        if (isValidEntryDate(to.params.date)) {
          this.date = to.params.date;
        }
      }
    },
  },
};
</script>

<style scoped>
.submit {
  text-align: left;
  font-size: 11pt;
}

span.end-date {
  color: rgb(255, 208, 56);
  font-weight: bold;
}

.save-draft {
  width: 150px;
  margin-right: 20px;
}
</style>
