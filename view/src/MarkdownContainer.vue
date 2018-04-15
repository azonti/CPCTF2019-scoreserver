<template>
  <article class="markdown well" v-html="rendered"></article>
</template>

<script>
import MarkdownIt from "markdown-it";
import MarkdownItKatex from "markdown-it-katex";

const md = new MarkdownIt({
  html: true,
  breaks: true,
  linkify: true,
});
md.use(MarkdownItKatex);
md.linkify.set({fuzzyLink: false});

export default {
  props: [
    'body'
  ],
  computed: {
    rendered () {
      return md.render(this.body)
    }
  }
}
</script>

<style>
p:last-child {
  margin: 0;
}
</style>
