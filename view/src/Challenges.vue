<template>
  <div>
    <vue-headful title="Challenges | CPCTF2018" />
    <div v-for="challenge in challenges">
      <p>{{ challenge.name }}</p>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      loading: true,
      challenges: []
    }
  },
  created() {
    this.fetchChallenges()
  },
  methods: {
    fetchChallenges() {
      axios.get(process.env.API_URL_PREFIX + "/challenges")
      .then((res) => {
        this.challenges = res.data
      })
    }
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}

h1, h2 {
  font-weight: normal;
}

ul {
  list-style-type: none;
  padding: 0;
}

li {
  display: inline-block;
  margin: 0 10px;
}

a {
  color: #42b983;
}
</style>
