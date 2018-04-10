import Vue from 'vue'
import router from './router'
import headful from 'vue-headful'

Vue.component('vue-headful', headful)

new Vue({
  el: '#app',
  router
})
