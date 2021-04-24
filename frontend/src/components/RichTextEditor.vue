<template>
  <div
    class="editor rich-text-editor form-control"
    v-on:foo="this.editor.chain().focus().setLink({href: linkUrl}).run()"
  >
    <div class="menubar">
      <EditorMenuButton
        class="btn-bold"
        :isActive="editor.isActive('bold')"
        tooltip="Bold"
        @click="editor.chain().focus().toggleBold().run()"
      >
        <b-icon-type-bold></b-icon-type-bold>
      </EditorMenuButton>

      <EditorMenuButton
        class="btn-italic"
        :isActive="editor.isActive('italic')"
        tooltip="Italic"
        @click="editor.chain().focus().toggleItalic().run()"
      >
        <b-icon-type-italic></b-icon-type-italic>
      </EditorMenuButton>

      <EditorMenuButton
        class="btn-strikethrough"
        :isActive="editor.isActive('strike')"
        tooltip="Strikethrough"
        @click="editor.chain().focus().toggleStrike().run()"
      >
        <b-icon-type-strikethrough></b-icon-type-strikethrough>
      </EditorMenuButton>

      <EditorMenuButton
        class="btn-link"
        :isActive="editor.isActive('link')"
        tooltip="Link"
        @click="
          editor.chain().focus().setLink({href: 'https://google.com'}).run()
        "
      >
        <b-icon-link></b-icon-link>
      </EditorMenuButton>

      <EditorMenuButton
        class="btn-inline-code"
        :isActive="editor.isActive('code')"
        tooltip="Inline code"
        @click="editor.chain().focus().toggleCode().run()"
      >
        <b-icon-code></b-icon-code>
      </EditorMenuButton>

      <EditorMenuButton
        class="btn-code-block"
        :isActive="editor.isActive('codeBlock')"
        tooltip="Code block"
        @click="editor.chain().focus().toggleCodeBlock().run()"
      >
        <b-icon-file-code></b-icon-file-code>
      </EditorMenuButton>

      <EditorMenuButton
        class="btn-h1"
        :isActive="editor.isActive('heading', {level: 1})"
        tooltip="Main heading"
        @click="editor.chain().focus().toggleHeading({level: 1}).run()"
      >
        H1
      </EditorMenuButton>

      <EditorMenuButton
        class="btn-h2"
        :isActive="editor.isActive('heading', {level: 2})"
        tooltip="Subheading"
        @click="editor.chain().focus().toggleHeading({level: 2}).run()"
      >
        H2
      </EditorMenuButton>

      <EditorMenuButton
        class="btn-h3"
        :isActive="editor.isActive('heading', {level: 3})"
        tooltip="Subsection heading"
        @click="editor.chain().focus().toggleHeading({level: 3}).run()"
      >
        H3
      </EditorMenuButton>

      <EditorMenuButton
        class="btn-bulleted-list"
        :isActive="editor.isActive('bulletList')"
        tooltip="Bulleted list"
        @click="editor.chain().focus().toggleBulletList().run()"
      >
        <b-icon-list-ul></b-icon-list-ul>
      </EditorMenuButton>

      <EditorMenuButton
        class="btn-numbered-list"
        :isActive="editor.isActive('orderedList')"
        tooltip="Numbered list"
        @click="editor.chain().focus().toggleOrderedList().run()"
      >
        <b-icon-list-ol></b-icon-list-ol>
      </EditorMenuButton>

      <EditorMenuButton
        class="btn-blockquote"
        :isActive="editor.isActive('blockquote')"
        tooltip="Quote"
        @click="editor.chain().focus().toggleBlockquote().run()"
      >
        <b-icon-blockquote-left></b-icon-blockquote-left>
      </EditorMenuButton>

      <EditorMenuButton
        class="switch-mode"
        tooltip="For markdown enthusiasts"
        @click="$emit('change-mode')"
      >
        Plaintext Editor
      </EditorMenuButton>
    </div>

    <hr />

    <editor-content class="editor-content" :editor="editor" />
    <b-modal
      ref="insert-link"
      title="Insert link"
      @shown="$refs['url-input'].focus()"
    >
      <b-form @submit.stop.prevent>
        <b-form-input ref="url-input" v-model="linkUrl"></b-form-input>
      </b-form>
    </b-modal>
  </div>
</template>

<script>
import Vue from 'vue';
import VueMarkdown from 'vue-markdown';

Vue.use(VueMarkdown);

import {Editor, EditorContent} from '@tiptap/vue-2';
import Link from '@tiptap/extension-link';
import {defaultExtensions} from '@tiptap/starter-kit';
import showdown from 'showdown';
import TurndownService from 'turndown';
import {gfm} from 'turndown-plugin-gfm';

import EditorMenuButton from '@/components/EditorMenuButton';

const showdownService = new showdown.Converter({
  omitExtraWLInCodeBlocks: true,
  noHeaderId: true,
  simplifiedAutoLink: true,
  excludeTrailingPunctuationFromURLs: true,
  strikethrough: true,
  tables: true,
  openLinksInNewWindow: true,
  emoji: true,
  simpleLineBreaks: true,
  encodeEmails: false,
  requireSpaceBeforeHeadingText: false,
});

const turndownService = new TurndownService({headingStyle: 'atx'});

export default {
  components: {
    EditorContent,
    EditorMenuButton,
  },
  props: {
    value: String,
  },
  data() {
    return {
      linkUrl: '',
      editor: new Editor({
        autofocus: true,
        extensions: [Link, ...defaultExtensions()],
        content: showdownService.makeHtml(this.value),
        onUpdate: () => {
          this.$emit('input', this.htmlToMarkdown(this.editor.getHTML()));
        },
      }),
    };
  },
  methods: {
    onClickLink() {
      this.linkUrl = 'https://';
      this.$refs['insert-link'].show();
    },
    htmlToMarkdown(html) {
      turndownService.use(gfm);
      let markdown = turndownService.turndown(html);
      // Trim trailing spaces.
      markdown = markdown.replace(/[ \t\r]+\n/g, '\n');
      // Consolidate multiple whitespace/newlines to just the double newlines.
      markdown = markdown.replace(/\s+\n\n/g, '\n\n');
      return markdown;
    },
  },
  beforeDestroy() {
    this.editor.destroy();
  },
};
</script>

<style scoped>
.menubar {
  margin: 0.5em 0em;
}

.menubar >>> .btn {
  margin-right: 0.5rem;
}

.editor {
  margin: 1rem 0;
  height: inherit;
}

.editor-content >>> .ProseMirror {
  outline: none;
  min-height: 250px;
}

.editor-content >>> blockquote {
  font-style: italic;
  margin-left: 1rem;
  background: rgb(224, 224, 224);
}
</style>
