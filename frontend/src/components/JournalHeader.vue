<template>
  <p class="journal-header text-left">
    <Avatar
      class="avatar"
      :username="entryAuthor"
      :to="'/' + entryAuthor"
      size="md"
    />
    <username-link :data-username="entryAuthor" />'s update for the week ending
    on
    <b>{{ entryDate | moment('utc', 'dddd, ll') }}</b>
  </p>
</template>

<script>
import Avatar from '@/components/Avatar.vue';

export default {
  name: 'JournalHeader',
  props: {
    entryAuthor: String,
    entryDate: Date,
  },
  components: {
    Avatar,
  },
};
</script>

<style scoped>
.avatar {
  margin-right: 0.5rem;
}
</style>
