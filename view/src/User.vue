<template>
  <div>
    <div v-if="user.id">
      <vue-headful :title="`${user.name} | CPCTF2018`" />
      <h1>{{ user.name }}</h1>
      <div class="row">
        <div class="col-md-4">
          <img :src="user.icon_url" style="width: 80%; height: 80%; margin-bottom: 5px">
          <dl class="row">
            <dt class="col-xs-4">Twitter</dt>
            <dd class="col-xs-8"><a v-if="user.twitter_screen_name" :href="`https://twitter.com/${user.twitter_screen_name}`">@{{ user.twitter_screen_name }}</a></dd>
          </dl>
          <div v-if="!user.is_author">
            <dl class="row">
              <dt class="col-xs-4">Score</dt>
              <dd class="col-xs-8">{{ user.score }}</dd>
            </dl>
            <dl class="row">
              <dt class="col-xs-4">Solved</dt>
              <dd class="col-xs-8">
                <ul class="list-unstyled">
                  <li v-for="_solved in solved">{{ _solved.name }}</li>
                </ul>
              </dd>
            </dl>
          </div>
        </div>
        <div v-if="user.id === myID" class="col-md-8">
          <h2 style="margin-top: 0;">Settings</h2>
          <h3>WebShell</h3>
          <button v-show="!recreatingContainer" class="btn btn-danger" style="width: 100%;" @click="recreateContainerWarn = true">Recreate WebShell Container</button>
        </div>
      </div>
    </div>
    <div v-else>
      <p>Loading...</p>
    </div>
    <modal
      :show="recreateContainerWarn"
      @close="recreateContainerWarn = false"
      :callback="recreateContainer"
      :modal="{ title: 'Are you sure?', body: 'Your container will be fully initialized and your data will be lost.', showCancel: true, btnClass: 'btn-danger', btnBody: 'Sure' }"
    />
    <error-modal :errors="errors" />
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
      user: {},
      solved: [],
      myID: '',
      recreateContainerWarn: false,
      recreatingContainer: false,
      errors: []
    }
  },
  created () {
    this.fetchUser()
  },
  watch: {
    id (to, from) {
      this.fetchUser()
    }
  },
  methods: {
    fetchUser () {
      api.get(`${process.env.API_URL_PREFIX}/users/${this.id}`)
      .then(res => res.data)
      .then((data) => {
        for (const key in data) {
          this.$set(this.user, key, data[key])
        }
      })
      .catch((err) => {
        this.errors.push(err.response.data)
      })
      api.get(`${process.env.API_URL_PREFIX}/users/${this.id}/solved`)
      .then(res => res.data)
      .then((data) => {
        this.solved.push(...data)
      })
      api.get(`${process.env.API_URL_PREFIX}/users/me`)
      .then(res => res.data)
      .then((data) => {
        this.myID = data.id
      })
    },
    recreateContainer () {
      this.recreatingContainer = true
      api.post(`${process.env.API_URL_PREFIX}/users/me`, { code: 'rwsc' })
      .catch((err) => {
        this.errors.push(err.response.data)
      })
      .then(() => {
        this.recreatingContainer = false
      })
    }
  }
}
</script>

<style>
</style>
