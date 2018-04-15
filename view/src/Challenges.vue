<template>
  <div>
    <vue-headful title="Challenges | CPCTF2018" />
    <div v-if="!loading">
      <div v-for="challenges in genre2Challenges">
        <h1 :class="'chal-name-' + challenges[0].genre">{{ challenges[0].genre }}</h1>
        <div class="row">
          <div v-for="challenge in sort(challenges)" class="col-md-4">
            <div :class="['panel', 'panel-' + challenge.genre, challenge.who_solved.map(user => user.id).includes(me.id) ? 'panel-success' : 'panel-primary']">
              <div class="panel-heading">
                <h2 class="panel-title"><router-link :to="{name: 'challenge', params: {id: challenge.id}}">{{ challenge.name }}</router-link></h2>
              </div>
              <div class="panel-body panel-challenge">
                <dl class="row">
                  <dt class="col-xs-4 col-a-left">Author</dt>
                  <dd class="col-xs-8 col-a-right"><router-link :to="{name: 'user', params: {id: challenge.author.id}}"><img :src="challenge.author.icon_url" class="icon">{{ challenge.author.name }}<small v-if="challenge.author.twitter_screen_name">(@{{ challenge.author.twitter_screen_name }})</small></router-link></dd>
                </dl>
                <dl class="row">
                  <dt class="col-xs-4 col-a-left">Score</dt>
                  <dd class="col-xs-8 col-a-right chal-score">{{ challenge.score }}</dd>
                </dl>
                <dl class="row">
                  <dt class="col-xs-4 col-a-left">Solved By</dt>
                  <dd class="col-xs-8 col-a-right">{{ challenge.who_solved.length }}</dd>
                </dl>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div v-else>
      <p>Loading ...</p>
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
    "me"
  ],
  data () {
    return {
      loading: true,
      genre2Challenges: {}
    }
  },
  created () {
    this.fetchChallenges()
  },
  methods: {
    fetchChallenges () {
      api.get(`${process.env.API_URL_PREFIX}/challenges`)
      .then(res => res.data)
      .then((data) => {
        for (const challenge of data) {
          this.$set(this.genre2Challenges, challenge.genre, (this.genre2Challenges[challenge.genre] || []).concat(challenge))
        }
      })
      .catch((err) => {
        this.$emit('error', `Message: ${err.response.data.message}`)
      })
      .then(() => {
        this.loading = false
      })
    },
    sort (challenges) {
      return [].concat(challenges).sort((a, b) => a.score > b.score)
    }
  }
}
</script>

<style>
</style>
