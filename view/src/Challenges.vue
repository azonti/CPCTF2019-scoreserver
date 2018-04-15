<template>
  <div>
    <vue-headful title="Challenges | CPCTF2018" />
    <div v-if="!loading">
      <section v-for="challenges in genre2Challenges">
        <h1>{{ challenges[0].genre }}</h1>
        <div class="row">
          <section v-for="challenge in challenges" class="col-md-4">
            <div class="panel" :class="challenge.who_solved.map(user => user.id).includes(me.id) ? 'panel-success' : 'panel-primary'">
              <div class="panel-heading">
                <h2 class="panel-title"><router-link :to="{name: 'challenge', params: {id: challenge.id}}">{{ challenge.name }}</router-link></h2>
              </div>
              <div class="panel-body">
                <dt class="row">
                  <dt class="col-xs-4">Author</dt>
                  <dd class="col-xs-8"><router-link :to="{name: 'user', params: {id: challenge.author.id}}"><img :src="challenge.author.icon_url" class="icon">{{ challenge.author.name }}<small v-if="challenge.author.twitter_screen_name">(@{{ challenge.author.twitter_screen_name }})</small></router-link></dd>
                </dt>
                <dt class="row">
                  <dt class="col-xs-4">Score</dt>
                  <dd class="col-xs-8">{{ challenge.score - challenge.hints.filter(hint => hint.caption).map(hint => hint.penalty).reduce((a, c) => (a + c)) }}</dd>
                </dt>
                <dt class="row">
                  <dt class="col-xs-4">Solved By</dt>
                  <dd class="col-xs-8">{{ challenge.who_solved.length }}</dd>
                </dt>
              </div>
            </div>
          </section>
        </div>
      </section>
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
    }
  }
}
</script>

<style>
</style>
