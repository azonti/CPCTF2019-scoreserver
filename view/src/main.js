import Vue from 'vue'

import Router from 'vue-router'
import headful from 'vue-headful'

import App from './App.vue'
import Index from './Index.vue'
import Challenges from './Challenges.vue'
import Challenge from './Challenge.vue'
import Ranking from './Ranking.vue'
import User from './User.vue'

import Modal from './Modal.vue'
import ErrorModal from './ErrorModal.vue'
import SuccessModal from './SuccessModal.vue'

import './assets/css/hacker.css'
import './assets/css/my.css'

Vue.use(Router)

Vue.component('modal', Modal)
Vue.component('error-modal', ErrorModal)
Vue.component('success-modal', SuccessModal)

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
    },
    {
      path: '/ranking',
      name: 'ranking',
      component: Ranking
    },
    {
      path: '/users/:id',
      name: 'user',
      props: true,
      component: User
    }

  ]
})

Vue.component('vue-headful', headful)

new Vue({
  el: '#app',
  router,
  render: h => h(App)
})
