<template>
  <div>
    <div v-if="user.id">
      <vue-headful :title="`${user.name} | CPCTF2018`" />
      <h1>{{ user.name }}</h1>
      <div class="row">
        <div class="col-md-4">
          <img :src="user.icon_url" style="width: 80%; height: 80%; margin-bottom: 5px">
          <p>
            <span class="badge" v-if="user.is_author">Author</span>
            <span class="badge" v-if="user.is_onsite">Onsite</span>
          </p>
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
                  <li class="solved-chals" v-for="_solved in solved"><router-link :to="{ name: 'challenge', params: { id: _solved.id } }">{{ _solved.name }}</router-link></li>
                </ul>
              </dd>
            </dl>
          </div>
        </div>
        <div v-if="user.id === me.id" class="col-md-8">
          <h2 style="margin-top: 0;">Settings</h2>
          <h3>Code</h3>
          <div class="row">
            <div class="col-md-10">
              <input class="form-control" v-model="code" placeholder="4WESOME_C0DE_I5_H3RE">
            </div>
            <div class="col-md-2">
              <button v-if="!sendingCode" @click="sendCode" class="btn btn-primary" style="width: 100%;">Check</button>
            </div>
          </div>
          <h3 v-if="user.is_onsite">WebShell</h3>
          <button v-if="user.is_onsite && !recreatingContainer" class="btn btn-danger" style="width: 100%;" @click="recreateContainerWarn = true">Recreate WebShell Container</button>
        </div>
      </div>
    </div>
    <div v-else>
      <p class="loading">Loading ...</p>
    </div>
    <modal
      :show="recreateContainerWarn"
      @close="recreateContainerWarn = false"
      :callback="recreateContainer"
      :modal="{ title: 'Are you sure?', body: 'Your container will be fully initialized and your data will be lost.', showCancel: true, btnClass: 'btn-danger', btnBody: 'Sure' }"
    />
  </div>
</template>

<style>
.badge {
  border: 1px solid;
  border-radius: 2px;
}
</style>

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
      user: {},
      solved: [],
      code: '',
      recreateContainerWarn: false,
      sendingCode: false,
      recreatingContainer: false,
    }
  },
  created () {
    this.fetchUser()
  },
  watch: {
    id (val) {
      this.fetchUser()
    }
  },
  methods: {
    fetchUser () {
      return Promise.all([
        api.get(`${process.env.API_URL_PREFIX}/users/${this.id}`)
        .then(res => res.data).then((data) => {
          for (const key in data) {
            this.$set(this.user, key, data[key])
          }
        }),
        api.get(`${process.env.API_URL_PREFIX}/users/${this.id}/solved`)
        .then(res => res.data).then((data) => {
          this.solved.splice(0, data.length, ...data)
        })
      ])
      .catch((err) => {
        this.$emit('error', err.response ? `Message: ${err.response.data.message}` : err)
      })
    },
    sendCode () {
      this.sendingCode = true
      return api.post(`${process.env.API_URL_PREFIX}/users/me`, { code: this.code })
      .then(() => {
        this.$emit('success', 'Your code has been activated.')
      })
      .then(() => Promise.all([
        new Promise((resolve) => { this.$emit('reloadMe', resolve, resolve) }),
        this.fetchUser()
      ]))
      .catch((err) => {
        this.$emit('error', err.response ? `Message: ${err.response.data.message}` : err)
      })
      .finally(() => {
        this.sendingCode = false
      })
    },
    recreateContainer () {
      this.recreatingContainer = true
      return api.post(`${process.env.API_URL_PREFIX}/users/me`, { code: 'rwsc' })
      .then(() => {
        this.$emit('success', 'Your container has been recreated.')
      })
      .then(() => Promise.all([
        new Promise((resolve) => { this.$emit('reloadMe', resolve, resolve) }),
        this.fetchUser()
      ]))
      .catch((err) => {
        this.$emit('error', err.response ? `Message: ${err.response.data.message}` : err)
      })
      .finally(() => {
        this.recreatingContainer = false
      })
    }
  }
}
</script>

<style>
</style>
