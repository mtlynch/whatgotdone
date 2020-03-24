<template>
  <textarea-autosize
    v-model="contents"
    :min-height="250"
    :max-height="650"
    @input="onInput"
    @paste.native="onPaste"
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
    onPaste(evt) {
      console.log('on paste', evt.clipboardData.items);
      for (const item of evt.clipboardData.items) {
        console.log(item);
        if (item.type.indexOf('image') < 0) continue;
        const imgFile = item.getAsFile();
        console.log(imgFile);
        this.insertSomething('[](/uploads/adfang.jpg)');
      }
    },
    insertSomething: function(insert) {
      const self = this;
      var tArea = this.$refs.entryText;
      // filter:
      if (0 == insert) {
        return;
      }
      if (0 == cursorPos) {
        return;
      }

      // get cursor's position:
      var startPos = tArea.selectionStart,
        endPos = tArea.selectionEnd,
        cursorPos = startPos,
        tmpStr = tArea.value;

      // insert:
      self.entryContent =
        tmpStr.substring(0, startPos) +
        insert +
        tmpStr.substring(endPos, tmpStr.length);

      // move cursor:
      setTimeout(() => {
        cursorPos += insert.length;
        tArea.selectionStart = tArea.selectionEnd = cursorPos;
      }, 10);
    },
  },
  watch: {
    value: function(newValue) {
      this.contents = newValue;
    },
  },
};
</script>
