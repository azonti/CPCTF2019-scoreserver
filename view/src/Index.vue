<template>
  <div>
    <vue-headful title="CPCTF2018" />
    <markdown-container :body="body" />
  </div>
</template>

<script>
import axios from 'axios'
import md from './index.md'

export default {
  data () {
    return {
      body: "loading ..."
    }
  },
  created () {
    axios.get(md).then(r => this.body = r.data)
  }
}
</script>
