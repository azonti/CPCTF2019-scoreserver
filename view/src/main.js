import Vue from 'vue'

import Router from 'vue-router'
import headful from 'vue-headful'

import App from './App.vue'
import Index from './Index.vue'
import Challenges from './Challenges.vue'
import Challenge from './Challenge.vue'
import Ranking from './Ranking.vue'
import User from './User.vue'
import Questions from './Questions.vue'
import AddChallenge from './AddChallenge.vue'
import UpdateChallenge from './UpdateChallenge.vue'

import Modal from './Modal.vue'
import EzModals from './EzModals.vue'
import MarkdownContainer from './MarkdownContainer.vue'

import 'katex/dist/katex.min.css'
import './assets/css/hacker.css'
import './assets/css/my.css'

Vue.use(Router)

Vue.component('modal', Modal)
Vue.component('ez-modals', EzModals)
Vue.component('markdown-container', MarkdownContainer)

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
    },
    {
      path: '/questions',
      name: 'questions',
      props: true,
      component: Questions
    },
    {
      path: '/add_challenge-e29375e2-046e-4bb6-b953-e65ea498824a',
      name: 'add_challenge',
      component: AddChallenge
    },
    {
      path: '/challenges/:id/update_challenge',
      name: 'update_challenge',
      props: true,
      component: UpdateChallenge
    },
  ],
  linkActiveClass: 'active'
})

Vue.component('vue-headful', headful)

new Vue({
  el: '#app',
  router,
  render: h => h(App)
})
