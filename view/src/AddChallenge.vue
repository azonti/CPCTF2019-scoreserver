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
          <h3>Group ID</h3>
          <div class="row">
            <div class="col-md-10">
              <input class="form-control" v-model="group_id" placeholder="Group ID">
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
          <h3>部分点問題の場合はチェック</h3>
          <div class="row">
            <div class="col-md-10">
              <input type="checkbox" v-model="is_not_complete">
            </div>
          </div>
          <h3>Flag</h3>
          <div class="row">
            <div class="col-md-10">
              <input class="form-control" v-model="flag" placeholder="FLAG_X00{THIS_IS_FLAG}">
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
      <p class="loading">Loading ...</p>
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
      group_id:   '',
      genre:      '',
      score:      '',
    	is_not_complete:false,
      flag:       '',
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
      console.log("run api.post")
      return api.post(`${process.env.API_URL_PREFIX}/challenges`, { 
          name:       this.name,
          group_id:   this.group_id,
          author:     {id: this.me.id},
          genre:      this.genre,
          score:      score,
          is_complete:!this.is_complete,
          flag:       this.flag,
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
