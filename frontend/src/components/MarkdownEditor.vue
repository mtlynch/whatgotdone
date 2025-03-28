<template>
  <div class="editor markdown-editor">
    <textarea-autosize
      v-model="contents"
      ref="editor"
      class="editor-textarea form-control"
      :min-height="250"
      :max-height="650"
      @input="onInput"
      @paste.native="onPaste"
    ></textarea-autosize>
    <p>
      (You can use
      <a href="https://www.markdownguide.org/cheat-sheet/" target="_blank"
        >Markdown</a
      >)
    </p>
  </div>
</template>

<script>
import Vue from 'vue';
import VueTextareaAutosize from 'vue-textarea-autosize';

import {uploadMedia} from '@/controllers/Media';

Vue.use(VueTextareaAutosize);

export default {
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
        const pastedFile = item.getAsFile();
        if (!pastedFile) {
          return;
        }
        let selectedText = this.getSelectedText();
        if (!selectedText) {
          selectedText = 'image';
        }
        this.insertTextAtCursorPosition(
          `[${selectedText}](uploading...)`,
          true
        );
        uploadMedia(pastedFile)
          .then((url) => {
            this.insertTextAtCursorPosition(`[${selectedText}](${url})`, false);
          })
          .catch((err) => {
            this.insertTextAtCursorPosition(
              `[${selectedText}](upload failed: ${err})`,
              false
            );
          });
      }
    },
    getSelectedText: function () {
      const textarea = this.$refs.editor.$el;
      return textarea.value.slice(
        textarea.selectionStart,
        textarea.selectionEnd
      );
    },
    insertTextAtCursorPosition: function (text, highlight) {
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

      this.$nextTick(function () {
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
    value: function (newValue) {
      this.contents = newValue;
    },
  },
};
</script>

<style scoped>
.editor-textarea {
  border: 0;
}
</style>
