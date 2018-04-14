<template>
  <div>
    <div v-if="challenge.id">
      <vue-headful :title="`${challenge.name} | CPCTF2018`" />
      <h1>{{ challenge.name }}</h1>
      <div class="row">
        <div class="col-md-4">
          <dl class="row">
            <dt class="col-xs-4">Genre</dt>
            <dd class="col-xs-8">{{ challenge.genre }}</dd>
          </dl>
          <dl class="row">
            <dt class="col-xs-4">Author</dt>
            <dd class="col-xs-8"><router-link :to="{name: 'user', params: {id: challenge.author.id}}"><img :src="challenge.author.icon_url" class="icon">{{ challenge.author.name }}<small v-if="challenge.author.twitter_screen_name">(@{{ challenge.author.twitter_screen_name }})</small></router-link></dd>
          </dl>
          <dl class="row">
            <dt class="col-xs-4">Score</dt>
            <dd class="col-xs-8">{{ challenge.score }}</dd>
          </dl>
          <dl class="row">
            <dt class="col-xs-4">Solved By</dt>
            <dd class="col-xs-8">
              <ul class="list-unstyled">
                <li v-for="who_solved in challenge.who_solved"><router-link :to="{name: 'user', params: {id: who_solved.id}}"><img :src="who_solved.icon_url" class="icon">{{ who_solved.name }}<small v-if="who_solved.twitter_screen_name">(@{{ who_solved.twitter_screen_name }})</small></router-link></li>
              </ul>
            </dd>
          </dl>
        </div>
        <div class="col-md-8">
          <h2 style="margin-top: 0;">Description</h2>
          <div class="row">
            <div class="col-md-10">
              <p class="well">{{ challenge.caption }}</p>
            </div>
          </div>
          <div v-for="hint in challenge.hints">
            <h2>Hint {{ parseInt(hint.id.substr(challenge.id.length + 1)) + 1 }}</h2>
            <div class="row">
              <div class="col-md-10">
                <p class="well">{{ hint.caption || `Not opened. This hint's penalty is ${hint.penalty}.` }}</p>
              </div>
              <div class="col-md-2">
                <button v-if="!hint.caption" class="btn btn-primary" style="width: 100%;" @click="openHint(hint.id);">Open Hint {{ hint.id.substr(challenge.id.length + 1) }}</button>
              </div>
            </div>
          </div>
          <div v-if="challenge.answer">
            <h2>Answer</h2>
            <div class="row">
              <div class="col-md-10">
                <p class="well">{{ challenge.answer }}</p>
              </div>
            </div>
          </div>
          <div class="row" style="margin-top: 20px;">
            <div class="col-md-10">
              <input class="form-control" v-model="flag" placeholder="CPCTF{FL4G_1S_H3RE}">
            </div>
            <div class="col-md-2">
              <button v-if="!challenge.answer" @click="checkFlag" class="btn btn-primary" style="width: 100%;">Check</button>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div v-else>
      <p>Loading...</p>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
const api = axios.create({
  withCredentials: true
})

export default {
  props: [
    'id'
  ],
  data () {
    return {
      challenge: {},
      flag: ''
    }
  },
  created () {
    this.fetchChallenge()
  },
  watch: {
    id (val) {
     this.fetchChallenge()
    }
  },
  methods: {
    fetchChallenge () {
      api.get(`${process.env.API_URL_PREFIX}/challenges/${this.id}`)
      .then(res => res.data)
      .then((data) => {
        for (const key in data) {
          this.$set(this.challenge, key, data[key])
        }
        this.flag = this.challenge.flag
        this.postURL = `${process.env.API_URL_PREFIX}/challenges/${this.id}`
      })
    },
    checkFlag () {
      api.post(`${process.env.API_URL_PREFIX}/challenges/${this.id}`, {
        flag: this.flag
      })
      .then(() => {
        fetchChallenge()
      })
    },
    openHint (id) {
      api.post(`${process.env.API_URL_PREFIX}/users/me`, {
        code: `hint:${id}`
      })
      .then(() => {
        this.fetchChallenge()
      })
    }
  }
}
</script>

<style>
</style>
