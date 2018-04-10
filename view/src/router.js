import Vue from 'vue'
import Router from 'vue-router'
import Index from './Index'
import Challenges from './Challenges'

Vue.use(Router)

export default new Router({
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
