<template>
  <div class="editor rich-text-editor form-control">
    <editor-menu-bar :editor="editor" v-slot="{commands, isActive}">
      <div class="menubar">
        <EditorMenuButton
          class="btn-bold"
          :isActive="isActive.bold()"
          tooltip="Bold"
          @click="commands.bold"
        >
          <b-icon-type-bold></b-icon-type-bold>
        </EditorMenuButton>

        <EditorMenuButton
          class="btn-italic"
          :isActive="isActive.italic()"
          tooltip="Italic"
          @click="commands.italic"
        >
          <b-icon-type-italic></b-icon-type-italic>
        </EditorMenuButton>

        <EditorMenuButton
          class="btn-strikethrough"
          :isActive="isActive.strike()"
          tooltip="Strikethrough"
          @click="commands.strike"
        >
          <b-icon-type-strikethrough></b-icon-type-strikethrough>
        </EditorMenuButton>

        <EditorMenuButton
          class="btn-link"
          :class="{'is-active': isActive.link()}"
          tooltip="Link"
          @click="onClickLink"
        >
          <b-icon-link></b-icon-link>
        </EditorMenuButton>

        <EditorMenuButton
          class="btn-inline-code"
          :isActive="isActive.code()"
          tooltip="Inline code"
          @click="commands.code"
        >
          <b-icon-code></b-icon-code>
        </EditorMenuButton>

        <EditorMenuButton
          class="btn-code-block"
          :isActive="isActive.code_block()"
          tooltip="Code block"
          @click="commands.code_block"
        >
          <b-icon-file-code></b-icon-file-code>
        </EditorMenuButton>

        <EditorMenuButton
          class="btn-h1"
          :isActive="isActive.heading({level: 1})"
          tooltip="Main heading"
          @click="commands.heading({level: 1})"
        >
          H1
        </EditorMenuButton>

        <EditorMenuButton
          class="btn-h2"
          :isActive="isActive.heading({level: 2})"
          tooltip="Subheading"
          @click="commands.heading({level: 2})"
        >
          H2
        </EditorMenuButton>

        <EditorMenuButton
          class="btn-h3"
          :isActive="isActive.heading({level: 3})"
          tooltip="Subsection heading"
          @click="commands.heading({level: 3})"
        >
          H3
        </EditorMenuButton>

        <EditorMenuButton
          class="btn-bulleted-list"
          :isActive="isActive.bullet_list()"
          tooltip="Bulleted list"
          @click="commands.bullet_list"
        >
          <b-icon-list-ul></b-icon-list-ul>
        </EditorMenuButton>

        <EditorMenuButton
          class="btn-numbered-list"
          :isActive="isActive.ordered_list()"
          tooltip="Numbered list"
          @click="commands.ordered_list"
        >
          <b-icon-list-ol></b-icon-list-ol>
        </EditorMenuButton>

        <EditorMenuButton
          class="btn-blockquote"
          :isActive="isActive.blockquote()"
          tooltip="Quote"
          @click="commands.blockquote"
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
    </editor-menu-bar>

    <hr />

    <editor-content class="editor-content" :editor="editor" />
  </div>
</template>

<script>
import Vue from 'vue';
import VueMarkdown from 'vue-markdown';

Vue.use(VueMarkdown);

import {Editor, EditorContent, EditorMenuBar} from 'tiptap';
import {
  Blockquote,
  CodeBlock,
  Heading,
  OrderedList,
  BulletList,
  ListItem,
  Bold,
  Code,
  Italic,
  Link,
  Strike,
} from 'tiptap-extensions';
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
    EditorMenuBar,
    EditorMenuButton,
  },
  props: {
    value: String,
  },
  data() {
    return {
      inRichTextMode: true,
      editor: new Editor({
        extensions: [
          new Blockquote(),
          new BulletList(),
          new CodeBlock(),
          new Heading({levels: [1, 2, 3]}),
          new ListItem(),
          new OrderedList(),
          new Link(),
          new Bold(),
          new Code(),
          new Italic(),
          new Strike(),
        ],
        content: showdownService.makeHtml(this.value),
        onUpdate: ({getHTML}) => {
          this.$emit('input', this.htmlToMarkdown(getHTML()));
        },
      }),
    };
  },
  methods: {
    onClickLink() {
      console.log('onClickLink');
      // TODO: Implement a dialog for entering links
    },
    htmlToMarkdown(html) {
      console.log('html=', html);

      turndownService.use(gfm);
      let markdown = turndownService.turndown(html);
      // Trim trailing whitespace
      markdown = markdown.replace(/\s+\n/g, '\n');
      console.log('markdown=', {x: markdown});
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
