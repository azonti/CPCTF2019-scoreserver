<template>
  <div>
    <div v-if="me.is_author">
      <vue-headful :title="`${me.name} | CPCTF2019`" />
      <h1>{{ me.name }}</h1>
      <div class="row">
        <div class="col-md-4">
          <img :src="me.icon_url" style="width: 80%; height: 80%; margin-bottom: 5px">
          <p>
            <span class="badge" v-if="me.is_author">Organizer</span>
            <span class="badge" v-if="!me.is_author && me.is_onsite">Onsite Participant</span>
          </p>
          <dl class="row">
            <dt class="col-xs-4">Twitter</dt>
            <dd class="col-xs-8"><a v-if="me.twitter_screen_name" :href="`https://twitter.com/${me.twitter_screen_name}`">@{{ me.twitter_screen_name }}</a></dd>
          </dl>
        </div>
        <div class="col-md-8">
          <h2 style="margin-top: 0;">Post challenge</h2>
          <div class="row">
            <h3>Challenge name</h3>
            <div class="col-md-10">
              <input class="form-control" v-model="name" placeholder="Challenge name">
            </div>
          </div>
          <h3>Genre</h3>
          <div class="row">
            <div class="col-md-10">
              <input class="form-control" v-model="genre" placeholder="Genre">
            </div>
          </div>
          <h3>Score</h3>
          <div class="row">
            <div class="col-md-10">
              <input class="form-control" v-model="score" placeholder="X00">
            </div>
          </div>
          <h3>Flags</h3>
          <div class="row">
            <div class="col-md-10">
              <textarea v-model="flags_text" placeholder="FLAG_X00{FLAG1}\nFLAG_X00{FLAG2}" cols="60" rows="5"></textarea>
            </div>
          </div>
          <h3>問題文</h3>
          <div class="row">
            <div class="col-md-10">
              <textarea v-model="caption" placeholder="問題文" cols="60" rows="20"></textarea>
            </div>
          </div>
          <h3>ヒント1</h3>
          <div class="row">
            <div class="col-md-10">
              <textarea v-model="hints[0]" placeholder="ヒント1" cols="60" rows="20"></textarea>
            </div>
          </div>
          <h3>ヒント2</h3>
          <div class="row">
            <div class="col-md-10">
              <textarea v-model="hints[1]" placeholder="ヒント2" cols="60" rows="20"></textarea>
            </div>
          </div>
          <h3>解説</h3>
          <div class="row">
            <div class="col-md-10">
              <textarea v-model="hints[2]" placeholder="解説" cols="60" rows="20"></textarea>
            </div>
          </div>
          <div class="row">
            <div class="col-md-2">
              <button v-if="!sendingCode" @click="postChallenge" class="btn btn-primary" style="width: 100%;">Post</button>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div v-else>
      <p class="loading">ERROR ...</p>
    </div>
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
    'me'
  ],
  data () {
    return {
      name:       '',
      genre:      '',
      score:      '',
    	is_not_complete:false,
      flags_text: '',
      caption:    '',
      hints:      [],
      sendingCode: false,
    }
  },
  created () {
  },
  watch: {
  },
  methods: {
    postChallenge () {
      this.sendingCode = true
      var hints_tmp = []
      var penalty_percent = [10,30-10,99-30]
      var score = parseInt(this.score,10)
      for (var i = 0; i < this.hints.length; i++) {
        hints_tmp[i] = {
          id:this.name+":"+i.toString(),
          caption:this.hints[i],
          penalty:score * penalty_percent[i] / 100
        }
      }
      var flags = this.flags_text.split(/\r\n|\n/);
      var flags_tmp = []
      for (var i = 0; i < flags.length; i++) {
        if(flags[i]===""){continue;}
        var flagscore =  parseInt(flags[i].substring(5,8),10)
        if (isNaN(flagscore) || flags[i].substring(0,5) !== "FLAG_"){
          this.$emit('error', 'invalid flags')
          return;
        }
        flags_tmp[i] = {
          id:this.name+":"+i.toString(),
          score: flagscore,
          flag:flags[i]
        }
      }

      console.log("run api.post")
      return api.post(`${process.env.API_URL_PREFIX}/challenges`, { 
          name:       this.name,
          author:     {id: this.me.id},
          genre:      this.genre,
          score:      score,
          is_complete:!this.is_complete,
          flags:       flags_tmp,
          caption:    this.caption,
          hints:      hints_tmp,
          answer:     this.flag
       })
      .then(() => {
        this.$emit('success', 'New challenge added.')
      })
      .catch((err) => {
        this.$emit('error', err.response ? `Message: ${err.response.data.message}` : err)
      })
      .then(() => {
        this.sendingCode = false
      })
    }
  }
}
</script>

<style>
</style>
