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

import {getIssueMetadata} from '@/controllers/Github';
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
    onChange(newValue) {
      console.log('onChange', newValue);
    },
    onPaste(evt) {
      console.log('paste', evt);
      const clipboardText = evt.clipboardData.getData('text');
      evt.preventDefault();
      if (
        clipboardText &&
        (this.isGithubIssueUrl(clipboardText) ||
          this.isGithubPrUrl(clipboardText))
      ) {
        console.log('getting issue metadata');
        getIssueMetadata(clipboardText)
          .then(meta => {
            console.log('in then');
            this.insertTextAtCursorPosition(
              `[#${meta.number}: ${meta.title}](${clipboardText})`
            );
          })
          .catch(err => {
            console.log('caught error', err);
            this.insertTextAtCursorPosition(clipboardText);
          });
        return;
      }
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
        uploadImage(item.getAsFile()).then(url => {
          this.insertTextAtCursorPosition(`[${selectedText}](${url})`, false);
        });
      }
    },
    isGithubIssueUrl: function(url) {
      return url.match(/github.com\/.+\/.+\/issues\/[0-9]+/);
    },
    isGithubPrUrl: function(url) {
      return url.match(/github.com\/.+\/.+\/pull\/[0-9]+/);
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
