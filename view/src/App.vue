<template>
  <div id="app" class="container" @click="showDropdown=false">
    <ul class="nav nav-tabs">
      <router-link tag="li" :to="{name: 'index'}" exact><a>Home</a></router-link>
      <router-link tag="li" :to="{name: 'challenges'}" exact><a>Challenges</a></router-link>
      <router-link tag="li" :to="{name: 'questions'}" exact><a>Questions</a></router-link>
      <router-link tag="li" :to="{name: 'ranking'}" exact><a>Ranking</a></router-link>
      <li class="dropdown">
        <a :aria-expanded="showDropdown ? 'true' : 'false'" class="dropdown-toggle" @click.stop="showDropdown=!showDropdown" href="#">Me <span class="caret"></span></a>
        <ul class="dropdown-menu" v-bind:style="{ display: showDropdown ? 'block' : '' }">
          <router-link tag="li" v-if="me.id" :to="{name: 'user', params: {id: me.id}}"><a><img :src="me.icon_url" class="icon big">{{ me.name }}<small v-if="me.twitter_screen_name">(@{{ me.twitter_screen_name }})</small></a></router-link>
          <li v-else><a @click.prevent="showDropdown = false" href="#"><img :src="me.icon_url" class="icon big">{{ me.name }}</a></li>
          <li class="divider" v-if="me.web_shell_pass"></li>
          <li v-if="me.web_shell_pass"><a target="_blank" :href="`https://${me.id}:${me.web_shell_pass}@client.cpctf.site/`">Open webshell</a></li>
          <li v-if="me.web_shell_pass"><a target="_blank" :href="`https://${me.id}:${me.web_shell_pass}@client.cpctf.site/_/`">Open file browser</a></li>
          <li class="divider"></li>
          <li :class="{disabled: me.id}"><a @click="showDropdown = false" :href="!me.id && twitterLoginURL">Login with Twitter</a></li>
          <li :class="{disabled: !me.id}"><a @click="showDropdown = false" :href="me.id && logoutURL">Logout</a></li>
        </ul>
      </li>
      <li class="pull-right"><a>{{ timer }}</a></li>
    </ul>
    <div v-show="newAnswer" class="alert alert-dismissible alert-success">
      <button type="button" class="close" @click="newAnswer = false">&times;</button>
      A new question has been answered! <router-link :to="{name: 'questions'}" @click.native="newAnswer = false">Check</router-link>
    </div>
    <router-view
      :me="me"
      :questions="questions"
      @reloadMe="fetchMe"
      @reloadQuestions="fetchQuestions"
      @error="(error) => { errors.push(error); }"
      @success="(success) => { successes.push(success); }"
    />
    <ez-modals title="Error"  bodyClass="text-danger" :bodies="errors" />
    <ez-modals title="Success" :bodies="successes" />
  </div>
</template>

<script>
import axios from 'axios'
const api = axios.create({
  withCredentials: true
})

import nobodyIcon from './assets/nobody.svg'

const START_TIME = Date.parse(process.env.START_TIME)
const FINISH_TIME = Date.parse(process.env.FINISH_TIME)

export default {
  data () {
    return {
      timer: "",
      showDropdown: false,
      twitterLoginURL: `${process.env.API_URL_PREFIX}/auth/twitter`,
      logoutURL: `${process.env.API_URL_PREFIX}/logout`,
      me: {
        name: 'Guest',
        icon_url: nobodyIcon
      },
      questions: [],
      newAnswer: false,
      errors: [],
      successes: []
    }
  },
  created () {
    this.fetchMe()
    this.fetchQuestions(() => {}, () => {}, true)
    setInterval(this.fetchQuestions, 1<<15)
    setInterval(this.tickCount, 1000)
  },
  methods: {
    fetchQuestions (resolve, reject, first) {
      return api.get(`${process.env.API_URL_PREFIX}/questions`)
      .then(res => res.data).then((data) => {
        if (!first) {
          this.newAnswer = this.questions.filter(question => question.answer).length !== data.filter(datum => datum.answer).length
        }
        this.questions.splice(0, data.length, ...data)
      })
      .then(() => {
        (resolve || (() => {}))()
      })
      .catch((err) => {
        (reject || ((err) => { throw err }))(err)
      })
    },
    fetchMe (resolve, reject) {
      return api.get(`${process.env.API_URL_PREFIX}/users/me`)
      .then(res => res.data).then((data) => {
        for (const key in data) {
          this.$set(this.me, key, data[key])
        }
      })
      .then(() => {
        (resolve || (() => {}))()
      })
      .catch((err) => {
        (reject || ((err) => { throw err }))(err)
      })
    },
    tickCount () {
      const fill = i => `0${parseInt(i)}`.substr(-2)
      const format = t => `${fill(t/3600)}:${fill((t/60)%60)}:${fill(t%60)}`

      const now = new Date().getTime()
      const start = parseInt((START_TIME - now) / 1000)
      const finish = parseInt((FINISH_TIME - now) / 1000)

      if(start > 0){
        this.timer = format(start) + " until CTF starts"
      }else if(finish > 0){
        this.timer = format(finish) + " left"
      }else{
        this.timer = "CTF is over!"
      }
    }
  }
}
</script>

<style>
.container {
  padding-bottom: 72px;
}
</style>
