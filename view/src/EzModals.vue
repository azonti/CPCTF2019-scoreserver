<template>
  <div>
    <modal
      v-for="(body, i) in bodies"
      :key="i"
      :show="show[i]"
      @close="$set(show, i, false)"
      :modal="{ title: title, bodyClass: bodyClass, body: body === 'Message: undefined' ? 'network error' : body }"
    />
  </div>
</template>

<script>
export default {
  props: [
    'title',
    'bodyClass',
    'bodies'
  ],
  data () {
    return {
      show: []
    }
  },
  watch: {
    bodies (val) {
      while (this.show.length < val.length) {
        this.show.push(true)
      }
    }
  }
}
</script>

<style>
</style>
