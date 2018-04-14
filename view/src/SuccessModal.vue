<template>
  <div>
    <modal
      v-for="(success, i) in successes"
      :key="i"
      :show="show[i]"
      @close="$set(show, i, false)"
      :modal="{ title: 'Success', body: success }"
    />
  </div>
</template>

<script>
export default {
  props: [
    'successes'
  ],
  data () {
    return {
      show: []
    }
  },
  watch: {
    successes (val) {
      while (this.show.length < val.length) {
        this.show.push(true)
      }
    }
  }
}
</script>

<style>
</style>
