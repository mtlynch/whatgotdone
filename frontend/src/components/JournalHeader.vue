<template>
  <p class="journal-header text-left">
    <Avatar
      class="avatar"
      :username="entryAuthor"
      :to="'/' + entryAuthor"
      size="md"
    />
    <Username class="username" :username="entryAuthor" />'s update for the week
    ending on
    <b>{{ entryDate | moment('utc', 'dddd, ll') }}</b>
  </p>
</template>

<script>
import Avatar from '@/components/Avatar.vue';
import Username from '@/components/Username.vue';

export default {
  name: 'JournalHeader',
  props: {
    entryAuthor: String,
    entryDate: Date,
  },
  components: {
    Avatar,
    Username,
  },
};
</script>

<style scoped>
.avatar {
  margin-right: 0.5rem;
}
</style>
