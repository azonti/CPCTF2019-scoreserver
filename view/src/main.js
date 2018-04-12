import Vue from 'vue'

import Router from 'vue-router'
import headful from 'vue-headful'

import App from './App'
import Index from './Index'
import Challenges from './Challenges'
import Challenge from './Challenge'

import './assets/css/hacker.css'
import './assets/css/my.css'

Vue.use(Router)

const router = new Router({
  mode: 'history',
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
    },
    {
      path: '/challenges/:id',
      name: 'challenge',
      props: true,
      component: Challenge
    }
  ]
})

Vue.component('vue-headful', headful)

new Vue({
  el: '#app',
  router,
  render: h => h(App)
})
