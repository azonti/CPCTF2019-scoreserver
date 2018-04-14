<template>
  <div id="app" class="container">
    <ul class="nav nav-tabs">
      <li :class="{active: isActive('/challenges')}"><router-link :to="{name: 'challenges'}">Challenges</router-link></li>
      <li :class="{active: isActive('/challenges')}"><router-link :to="{name: 'challenges'}">Questions</router-link></li>
      <li :class="{active: isActive('/ranking')}"><router-link :to="{name: 'ranking'}">Ranking</router-link></li>
      <li class="dropdown">
        <a :aria-expanded="showDropdown ? 'true' : 'false'" class="dropdown-toggle" @click.prevent="showDropdown = !showDropdown" href="#">Me <span class="caret"></span></a>
        <ul class="dropdown-menu" v-bind:style="{ display: showDropdown ? 'block' : 'none' }">
          <li><router-link @click.native="showDropdown = !showDropdown" :to="me.id ? {name: 'user', params: {id: me.id}} : {}"><img :src="me.icon_url" class="icon big">{{ me.name }}<small v-if="me.twitter_screen_name">(@{{ me.twitter_screen_name }})</small></router-link></li>
          <li class="divider"></li>
          <li :class="{disabled: me.id}"><a @click="showDropdown = !showDropdown" :href="!me.id ? twitterLoginURL : '#'">Login with Twitter</a></li>
          <li :class="{disabled: !me.id}"><a @click="showDropdown = !showDropdown" :href="me.id ? logoutURL : '#'">Logout</a></li>
        </ul>
      </li>
    </ul>
    <router-view>
    </router-view>
  </div>
</template>

<script>
import axios from 'axios'
const api = axios.create({
  withCredentials: true
})

import nobodyIcon from './assets/nobody.svg'

export default {
  data () {
    return {
      showDropdown: false,
      twitterLoginURL: `${process.env.API_URL_PREFIX}/auth/twitter`,
      logoutURL: `${process.env.API_URL_PREFIX}/logout`,
      me: {
        name: 'Guest',
        icon_url: nobodyIcon
      }
    }
  },
  created () {
    this.fetchMe()
  },
  methods: {
    isActive (path) {
      return this.$route.path === path
    },
    fetchMe () {
      api.get(`${process.env.API_URL_PREFIX}/users/me`)
      .then(res => res.data)
      .then((data) => {
        for (const key in data) {
          this.$set(this.me, key, data[key])
        }
      })
    }
  }
}
</script>

<style>
</style>
