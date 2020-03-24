<template>
  <textarea-autosize
    v-model="contents"
    ref="editor"
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
        if (item.type.indexOf('image') < 0) continue;
        const imgFile = item.getAsFile();
        console.log(imgFile);
        // TODO: upload file to What Got Done.
        this.insertTextAtCursorPosition('[](/uploads/todo.jpg)');
      }
    },
    insertTextAtCursorPosition: function(text) {
      const element = this.$refs.editor.$el;
      if (!text) {
        return;
      }

      // get cursor's position
      const startPos = element.selectionStart;
      const endPos = element.selectionEnd;
      let cursorPos = startPos;
      const tmpStr = element.value;

      // insert
      this.contents =
        tmpStr.substring(0, startPos) +
        text +
        tmpStr.substring(endPos, tmpStr.length);

      // move cursor
      setTimeout(() => {
        cursorPos += text.length;
        element.selectionStart = element.selectionEnd = cursorPos;
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
