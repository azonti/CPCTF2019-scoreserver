<template>
  <div id="app" class="container">
    <ul class="nav nav-tabs">
      <router-link tag="li" :to="{name: 'challenges'}" exact @click.native="showDropdown = false"><a>Challenges</a></router-link>
      <router-link tag="li" :to="{name: 'challenges'}" exact @click.native="showDropdown = false"><a>Questions</a></router-link>
      <router-link tag="li" :to="{name: 'ranking'}" exact @click.native="showDropdown = false"><a>Ranking</a></router-link>
      <li class="dropdown">
        <a :aria-expanded="showDropdown ? 'true' : 'false'" class="dropdown-toggle" @click.prevent="showDropdown = !showDropdown" href="#">Me <span class="caret"></span></a>
        <ul class="dropdown-menu" v-bind:style="{ display: showDropdown ? 'block' : '' }">
          <router-link tag="li" v-if="me.id" :to="{name: 'user', params: {id: me.id}}" @click.native="showDropdown = false"><a><img :src="me.icon_url" class="icon big">{{ me.name }}<small v-if="me.twitter_screen_name">(@{{ me.twitter_screen_name }})</small></a></router-link>
          <li v-else><a @click.prevent="showDropdown = false" href="#"><img :src="me.icon_url" class="icon big">{{ me.name }}</a></li>
          <li class="divider" v-if="me.web_shell_pass"></li>
          <li v-if="me.web_shell_pass"><a target="_blank" :href="`https://${me.id}:${me.web_shell_pass}@client.cpctf.site/`">Open webshell</a></li>
          <li v-if="me.web_shell_pass"><a target="_blank" :href="`https://${me.id}:${me.web_shell_pass}@client.cpctf.site/_/`">Open file browser</a></li>
          <li class="divider"></li>
          <li :class="{disabled: me.id}"><a @click="showDropdown = false" :href="!me.id && twitterLoginURL">Login with Twitter</a></li>
          <li :class="{disabled: !me.id}"><a @click="showDropdown = false" :href="me.id && logoutURL">Logout</a></li>
        </ul>
      </li>
    </ul>
    <router-view
      :me="me"
      @reloadMe="fetchMe"
      @error="(error) => { errors.push(error); }"
      @success="(success) => { successes.push(success); }"
    >
    </router-view>
    <error-modal :errors="errors" />
    <success-modal :successes="successes" />
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
      },
      errors: [],
      successes: []
    }
  },
  created () {
    this.fetchMe()
  },
  methods: {
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
