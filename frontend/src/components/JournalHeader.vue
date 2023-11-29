<template>
  <p class="journal-header text-left">
    <user-avatar :username="entryAuthor" class="mr-1" />
    <username-link :username="entryAuthor" />'s update for the week ending on
    <b>{{ entryDate | moment('utc', 'dddd, ll') }}</b>
  </p>
</template>

<script>
export default {
  name: 'JournalHeader',
  props: {
    entryAuthor: String,
    entryDate: Date,
  },
};
</script>
