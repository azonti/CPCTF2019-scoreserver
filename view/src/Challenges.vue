<template>
  <div>
    <vue-headful title="Challenges | CPCTF2019" />
    <div class="toggle row">
      <div class="col-sm-4"><button :class="['btn', group == 'genre' ? 'btn-info' : 'btn-primary']" @click="group='genre'">Group by genre</button></div>
      <div class="col-sm-4"><button :class="['btn', group == 'real_score' ? 'btn-info' : 'btn-primary']" @click="group='real_score'">Group by score</button></div>
      <div class="col-sm-4"><button :class="['btn', group == 'solveCount' ? 'btn-info' : 'btn-primary']" @click="group='solveCount'">Group by solve count</button></div>
    </div>
    <div class="toggle row">
      <div class="col-xs-12"><button :class="['btn', hide ? 'btn-info' : 'btn-primary']" @click="hide=!hide">Hide solved challenges</button></div>
    </div>
    <div v-if="!loading">
      <div v-for="(chals, key) in chalSet">
        <h1 :class="'chal-name-' + key">{{ key }}</h1>
        <div class="row">
          <div v-for="challenge in chals" class="col-xl-3 col-md-4 col-sm-6" v-if="!challenge.solved || !hide">
            <router-link :to="{name: 'challenge', params: {id: challenge.id}, query: { hide: contestFinished ? 'true' : ''}}" :class="['panel', 'panel-' + challenge.genre, challenge.solved ? 'panel-success' : 'panel-primary']">
              <div class="panel-heading">
                <h2 class="panel-title">{{ challenge.name }}</h2>
              </div>
              <div class="panel-body panel-challenge">
                <dl class="row">
                  <dt class="col-xs-4 col-a-left">Author</dt>
                  <dd class="col-xs-8 col-a-right author"><router-link :to="{name: 'user', params: {id: challenge.author.id}}"><img :src="challenge.author.icon_url" class="icon">{{ challenge.author.name }} <small v-if="challenge.author.twitter_screen_name">(@{{ challenge.author.twitter_screen_name }})</small></router-link></dd>
                </dl>
                <dl class="row">
                  <dt class="col-xs-4 col-a-left">Score</dt>
                  <dd class="col-xs-8 col-a-right chal-score">{{ challenge.score }} <small class="level">({{ "★".repeat(challenge.real_score/100) }})</small></dd>
                  <div v-for="i in challenge.flags.length - 1">
                    <dt class="col-xs-4 col-a-left"></dt><dd class="col-xs-8 col-a-right chal-score">{{"   "}}<small class="level">({{ "★".repeat(challenge.flags[challenge.flags.length-i-1].score/100) }})</small></dd>
                  </div>
                </dl>
                <dl class="row">
                  <dt class="col-xs-4 col-a-left">Solved</dt>
                  <dd class="col-xs-8 col-a-right">{{ challenge.solveCount }} time{{ challenge.solveCount != 1 ? "s" : "" }}</dd>
                </dl>
              </div>
            </router-link>
          </div>
        </div>
      </div>
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

a.panel {
  display: block;
}
a.panel:hover .panel-title {
  text-decoration: underline;
}

.author {
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
}
.level {
  font-size: .6em;
}
</style>

<script>
import axios from 'axios'
const api = axios.create({
  withCredentials: true
})

const seasoner = {
  genre (obj) {
    return obj
  },
  real_score (obj) {
    const nobj = {}
    Object.keys(obj).sort().forEach(key => nobj[`${key} (${"★".repeat(key / 100)})`] = obj[key])
    return nobj
  },
  solveCount (obj) {
    const nobj = {}
    Object.keys(obj).map(e => parseInt(e)).sort((x, y) => y - x).forEach(key => nobj[`Solved ${key} time${key != 1 ? "s" : ""}`] = obj[key])
    return nobj
  }
}

export default {
  props: [
    "me"
  ],
  data () {
    return {
      loading: true,
      hide: false,
      contestFinished: false,
      grouped: {genre: {}, real_score: {}, solveCount: {}},
      group: "genre"
    }
  },
  created () {
    this.contestFinished = Date.parse(process.env.FINISH_TIME) <= Date.now()
    setInterval(() => {this.contestFinished = Date.parse(process.env.FINISH_TIME) <= Date.now()}, 1000)
    return api.get(`${process.env.API_URL_PREFIX}/challenges`)
    .then((res) => res.data)
    .then((data) => {
      data.forEach((chal) => {
        chal.solveCount = chal.who_solved.length
        chal.solved = chal.who_solved.some(user => user.id == this.me.id)

        for(const [key, chals] of Object.entries(this.grouped)){
          if(!(chal[key] in chals)){
            chals[chal[key]] = []
          }
          chals[chal[key]].push(chal)
        }
      })
    })
    .catch((err) => {
      this.$emit('error', err.response ? `Message: ${err.response.data.message}` : err)
    })
    .then(() => {
      this.loading = false
    })
  },
  computed: {
    chalSet () {
      return seasoner[this.group](this.grouped[this.group])
    }
  }
}
</script>

<style>
</style>
