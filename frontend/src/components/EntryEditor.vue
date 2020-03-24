<template>
  <textarea-autosize
    v-model="contents"
    :min-height="250"
    :max-height="650"
    @input="onInput"
  ></textarea-autosize>
</template>

<script>
import Vue from 'vue';
import VueTextareaAutosize from 'vue-textarea-autosize';

Vue.use(VueTextareaAutosize);

export default {
  name: 'EntryEditor',
  props: {
    value: String,
  },
  data() {
    return {
      contents: this.value,
    };
  },
  methods: {
    onInput(newValue) {
      this.$emit('input', newValue);
    },
  },
};
</script>
