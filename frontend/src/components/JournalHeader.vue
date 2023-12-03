<template>
  <div class="d-flex align-items-center mb-4">
    <user-avatar :username="entryAuthor" class="mr-2 d-block avatar" />
    <username-link :username="entryAuthor" />'s update for the week ending
    on&nbsp;
    <entry-date :date="entryDate" />
  </div>
</template>

<script>
export default {
  name: 'JournalHeader',
  props: {
    entryAuthor: String,
    entryDate: String,
  },
};
</script>

<style scoped>
.avatar {
  max-width: 45px;
  max-height: 45px;
}
</style>
