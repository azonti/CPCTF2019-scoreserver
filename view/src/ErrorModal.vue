<template>
  <div>
    <modal
      v-for="(error, i) in errors"
      :key="i"
      :show="show[i]"
      @close="$set(show, i, false)"
      :modal="{ title: 'Error', body: error }"
    />
  </div>
</template>

<script>
export default {
  props: [
    'errors'
  ],
  data () {
    return {
      show: []
    }
  },
  watch: {
    errors (val) {
      while (this.show.length < val.length) {
        this.show.push(true)
      }
    }
  }
}
</script>

<style>
</style>
