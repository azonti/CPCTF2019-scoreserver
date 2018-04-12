<template>
  <div>
    <div v-if="challenge">
      <vue-headful :title="`${challenge.name} | CPCTF2018`" />
      <h1>{{ challenge.name }}</h1>
      <div class="row">
        <div class="col-md-3">
          <dl class="row">
            <dt class="col-xs-6">Author</dt>
            <dd class="col-xs-6"><img :src="challenge.author.icon_url" class="icon">{{ challenge.author.name }}</dd>
          </dl>
          <dl class="row">
            <dt class="col-xs-6">Score</dt>
            <dd class="col-xs-6">{{ challenge.score }}</dd>
          </dl>
          <dl class="row">
            <dt class="col-xs-6">Solved</dt>
            <dd class="col-xs-6">
              <ul class="list-unstyled">
                <li v-for="who_solved in challenge.who_solved"><img :src="who_solved.icon_url" class="icon">{{ who_solved.name }}</li>
              </ul>
            </dd>
          </dl>
        </div>
        <div class="col-md-9">
          <h2>Description</h2>
          <p class="well">{{ challenge.caption }}</p>
          <div v-for="hint in challenge.hints">
            <h2>Hint {{ hint.id.substr(challenge.id.length + 1) }}</h2>
            <div class="row">
              <div class="col-md-10">
                <p class="well">{{ hint.caption ? hint.caption : `Not opened. This hint's penalty is ${hint.penalty}.` }}</p>
              </div>
              <div class="col-md-2">
                <button :class="hint.caption ? 'btn-disabled ' : 'btn-primary'" class="btn">Open Hint {{ hint.id.substr(challenge.id.length + 1) }}</button>
              </div>
            </div>
          </div>
          <div v-if="challenge.answer">
            <h2>Answer</h2>
            <p class="well">{{ challenge.answer }}</p>
          </div>
          <hr>
          <div class="row">
            <div class="col-md-10">
              <input class="form-control" v-model="flag" placeholder="CPCTF{FL4G_1S_H3RE}">
            </div>
            <div class="col-md-2">
              <button v-on:click="checkFlag" class="btn btn-primary">Check</button>
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
    "id"
  ],
  data() {
    return {
      challenge: {},
      flag: ''
    }
  },
  created() {
    this.fetchChallenge()
  },
  methods: {
    fetchChallenge() {
      api.get(`${process.env.API_URL_PREFIX}/challenges/${this.id}`)
      .then(res => res.data)
      .then((data) => {
        this.challenge = data
        this.postURL = `${process.env.API_URL_PREFIX}/challenges/${this.id}`
      })
    },
    checkFlag() {
      api.post(`${process.env.API_URL_PREFIX}/challenges/${this.id}`, {
        flag: this.flag
      })
      .then(() => {
        fetchChallenge()
      })
    }
  }
}
</script>

<style>
</style>
