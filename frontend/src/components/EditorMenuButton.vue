<template>
  <div class="button-wrapper">
    <b-button
      :variant="isActive ? 'dark' : 'outline-dark'"
      class="menubar__button"
      ref="editor-button"
      :class="{'is-active': isActive}"
      :disabled="disabled"
      @click="onClick"
    >
      <slot></slot>
    </b-button>
    <b-tooltip
      :target="() => $refs['editor-button']"
      triggers="hover"
      v-if="tooltip"
      :title="tooltip"
      placement="bottom"
    >
    </b-tooltip>
  </div>
</template>

<script>
export default {
  props: {
    isActive: Boolean,
    disabled: {
      type: Boolean,
      default: false,
    },
    tooltip: String,
  },
  methods: {
    onClick(evt) {
      this.$emit('click', evt);
    },
  },
};
</script>

<style scoped>
.button-wrapper {
  display: inline-block;
}
</style>
