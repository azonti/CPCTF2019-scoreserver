<template>
  <div>
    <vue-headful title="Ranking | CPCTF2018" />
    <div class="toggle row">
      <div class="col-xs-6"><button :class="['btn', all ? 'btn-info' : 'btn-primary']" @click="all=true">Overall Ranking</button></div>
      <div class="col-xs-6"><button :class="['btn', !all ? 'btn-info' : 'btn-primary']" @click="all=false">Onsite Ranking</button></div>
    </div>
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
          <tr v-for="(user, index) in filteredUsers">
            <td>{{ index + 1 }}</td>
            <td><router-link :to="{name: 'user', params: {id: user.id}}"><img :src="user.icon_url" class="icon">{{ user.name }}<small v-if="user.twitter_screen_name">(@{{ user.twitter_screen_name }})</small></router-link></td>
            <td><span v-if="user.is_onsite" class="glyphicon glyphicon-star"></span></th>
            <td>{{ user.score }}</td>
          </tr>
        </tbody>
      </table>
    </div>
    <div v-else>
      <p class="loading">Loading ...</p>
    </div>
  </div>
</template>

<style>
button {
  width: 100%;
}
.toggle {
  margin: 1em 0;
}
</style>

<script>
import axios from 'axios'
const api = axios.create({
  withCredentials: true
})

export default {
  props: [
    'me'
  ],
  data () {
    return {
      all: !this.me.is_onsite,
      loading: true,
      users: []
    }
  },
  created () {
    api.get(`${process.env.API_URL_PREFIX}/users`)
    .then(res => res.data).then((data) => {
      this.users = data.filter(a => !a.is_author).sort((a, b) => b.score - a.score)
    })
    .catch((err) => {
      this.$emit('error', err.response ? `Message: ${err.response.data.message}` : err)
    })
    .finally(() => {
      this.loading = false
    })
  },
  computed: {
    filteredUsers () {
      return this.all ? this.users : this.users.filter(u => u.is_onsite)
    }
  }
}
</script>
