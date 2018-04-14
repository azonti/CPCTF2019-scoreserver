<template>
  <div>
    <vue-headful title="Ranking | CPCTF2018" />
    <div v-if="!loading">
      <table class="table table-hover">
        <thead>
          <tr>
            <th>#</th>
            <th>Name</th>
            <th>Onsite</th>
            <th>Score</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(user, index) in users">
            <td>{{ index + 1 }}</td>
            <td><router-link :to="{name: 'user', params: {id: user.id}}"><img :src="user.icon_url" class="icon">{{ user.name }}<small v-if="user.twitter_screen_name">(@{{ user.twitter_screen_name }})</small></router-link></td>
            <td><span v-if="user.is_onsite" class="glyphicon glyphicon-star"></span></th>
            <td>{{ user.score }}</td>
          </tr>
        </tbody>
      </table>
    </div>
    <div v-else>
      <p>Loading...</p>
    </div>
    <error-modal :errors="errors" />
    <success-modal :successes="successes" />
  </div>
</template>

<script>
import axios from 'axios'
const api = axios.create({
  withCredentials: true
})

export default {
  data () {
    return {
      loading: true,
      users: [],
      errors: []
    }
  },
  created () {
    this.fetchUsers()
  },
  methods: {
    fetchUsers () {
      api.get(`${process.env.API_URL_PREFIX}/users`)
      .then(res => res.data)
      .then((data) => {
        this.users.push(...data.filter(a => !a.is_author).sort((a, b) => (b.score - a.score)))
      })
      .catch((err) => {
        this.errors.push(`Message: ${err.response.data.message}`)
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
