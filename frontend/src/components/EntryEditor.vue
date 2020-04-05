<template>
  <textarea-autosize
    v-model="contents"
    ref="editor"
    :min-height="250"
    :max-height="650"
    @input="onInput"
    @change.native="onChange"
    @paste.native="onPaste"
  ></textarea-autosize>
</template>

<script>
import Vue from 'vue';
import VueTextareaAutosize from 'vue-textarea-autosize';

import {uploadImage} from '@/controllers/Images';

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
      for (const item of evt.clipboardData.items) {
        if (item.type.indexOf('image') < 0) continue;
        let selectedText = this.getSelectedText();
        if (!selectedText) {
          selectedText = 'image';
        }
        this.insertTextAtCursorPosition(
          `[${selectedText}](uploading...)`,
          true
        );
        uploadImage(item.getAsFile())
          .then(url => {
            this.insertTextAtCursorPosition(`[${selectedText}](${url})`, false);
          })
          .catch(err => {
            this.insertTextAtCursorPosition(
              `[${selectedText}](upload failed: ${err})`,
              false
            );
          });
      }
    },
    getSelectedText: function() {
      const textarea = this.$refs.editor.$el;
      return textarea.value.slice(
        textarea.selectionStart,
        textarea.selectionEnd
      );
    },
    insertTextAtCursorPosition: function(text, highlight) {
      const textarea = this.$refs.editor.$el;
      if (!text) {
        return;
      }

      let cursorPos = textarea.selectionStart;
      textarea.value =
        textarea.value.substring(0, textarea.selectionStart) +
        text +
        textarea.value.substring(textarea.selectionEnd, this.contents.length);

      this.$emit('input', textarea.value);

      this.$nextTick(function() {
        if (highlight) {
          textarea.selectionStart = cursorPos;
          textarea.selectionEnd = cursorPos + text.length;
        } else {
          cursorPos += text.length;
          textarea.selectionStart = cursorPos;
          textarea.selectionEnd = cursorPos;
        }
      });
    },
  },
  watch: {
    value: function(newValue) {
      this.contents = newValue;
    },
  },
};
</script>
