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
            <dd class="col-xs-8">{{ challenge.score - challenge.hints.filter(hint => hint.caption).map(hint => hint.penalty).reduce((a, c) => (a + c)) }}</dd>
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
              <markdown-container :body="challenge.caption" />
            </div>
          </div>
          <div v-for="(hint, index) in challenge.hints">
            <h2>Hint #{{ index + 1 }}</h2>
            <div class="row">
              <div class="col-md-10">
                <markdown-container v-if="hint.caption" :body="hint.caption" />
                <p class="well" v-else-if="index == 0 || challenge.hints[index-1].caption">This hint's penalty is <strong>{{ hint.penalty }}</strong>.</p>
                <p class="well" v-else>Open hint #{{ index }} first.</p>
              </div>
              <div class="col-md-2">
                <button v-if="!openingHint && !hint.caption && (index == 0 || challenge.hints[index-1].caption)" class="btn btn-primary" style="width: 100%;" @click="hintToOpen = hint; openHintWarn = true;">Open Hint #{{ index+1 }}</button>
              </div>
            </div>
          </div>
          <div v-if="challenge.answer">
            <h2>Answer</h2>
            <div class="row">
              <div class="col-md-10">
                <markdown-container :body="challenge.answer" />
              </div>
            </div>
          </div>
          <div class="row" style="margin-top: 20px;">
            <div class="col-md-10">
              <input class="form-control" v-model="flag" placeholder="FLAG{FL4G_1S_H3RE}">
            </div>
            <div class="col-md-2">
              <button v-if="!checkingFlag && !challenge.answer" @click="checkFlag" class="btn btn-primary" style="width: 100%;">Check</button>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div v-else>
      <p>Loading...</p>
    </div>
    <modal
      :show="openHintWarn"
      @close="openHintWarn = false"
      :callback="openHint"
      :modal="{ title: 'Are you sure?', body: `This hint's penalty is ${hintToOpen.penalty}.`, showCancel: true, btnBody: 'Sure' }"
    />
  </div>
</template>

<script>
import axios from 'axios'
const api = axios.create({
  withCredentials: true
})

export default {
  props: [
    'id',
    'me'
  ],
  data () {
    return {
      challenge: {},
      flag: '',
      hintToOpen: {},
      openHintWarn: false,
      openingHint: false,
      checkingFlag: false
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
      .catch((err) => {
        this.$emit('error', `Message: ${err.response.data.message}`)
      })
    },
    checkFlag () {
      this.checkingFlag = true
      api.post(`${process.env.API_URL_PREFIX}/challenges/${this.id}`, {
        flag: this.flag
      })
      .then(() => {
        this.$emit('success', 'Congrats!!')
      })
      .then(() => this.fetchChallenge())
      .catch((err) => {
        this.$emit('error', `Message: ${err.response.data.message}`)
      })
      .then(() => {
        this.checkingFlag = false
      })
    },
    openHint () {
      this.openingHint = true
      api.post(`${process.env.API_URL_PREFIX}/users/me`, {
        code: `hint:${this.hintToOpen.id}`
      })
      .then(() => this.fetchChallenge())
      .catch((err) => {
        this.$emit('error', `Message: ${err.response.data.message}`)
      })
      .then(() => {
        this.openingHint = false
      })
    }
  }
}
</script>

<style>
</style>
