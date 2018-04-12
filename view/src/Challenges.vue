<template>
  <div>
    <div v-if="!loading">
      <vue-headful title="Challenges | CPCTF2018" />
      <div v-for="challenges in genre2Challenges">
        <h1>{{ challenges[0].genre }}</h1>
        <div class="row">
          <div v-for="challenge in challenges" class="col-md-3">
            <div class="panel panel-primary">
              <div class="panel-heading">
                <h2 class="panel-title"><router-link :to="{name: 'challenge', params: {id: challenge.id}}">{{ challenge.name }}</router-link></h2>
              </div>
              <div class="panel-body">
                <dt class="row">
                  <dt class="col-xs-6">Author</dt>
                  <dd class="col-xs-6"><img :src="challenge.author.icon_url" class="icon">{{ challenge.author.name }}</dd>
                </dt>
                <dt class="row">
                  <dt class="col-xs-6">Score</dt>
                  <dd class="col-xs-6">{{ challenge.score }}</dd>
                </dt>
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
  data() {
    return {
      loading: true,
      genre2Challenges: {}
    }
  },
  created() {
    this.fetchChallenges()
  },
  methods: {
    fetchChallenges() {
      api.get(`${process.env.API_URL_PREFIX}/challenges`)
      .then(res => res.data)
      .then((data) => {
        for (const challenge of data) {
          this.$set(this.genre2Challenges, challenge.genre, (this.genre2Challenges[challenge.genre] || []).concat(challenge))
          this.loading = false
        }
      })
    }
  }
}
</script>

<style>
</style>
