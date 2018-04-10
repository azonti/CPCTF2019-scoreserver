import Vue from 'vue'

import Router from 'vue-router'
import headful from 'vue-headful'

import Index from './Index'
import Challenges from './Challenges'

Vue.use(Router)

const router = new Router({
  routes: [
    {
      path: '/',
      name: 'index',
      component: Index
    },
    {
      path: '/challenges',
      name: 'challenges',
      component: Challenges
    }
  ]
})

Vue.component('vue-headful', headful)

new Vue({
  el: '#app',
  router
})
