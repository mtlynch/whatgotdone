<template>
  <div class="missing-page">
    <p>You typed {{ $route.path }}</p>
    <p>We don't have a page for that. Shrug!</p>
  </div>
</template>

<script>
export default {
  name: "MissingPage"
};
</script>
