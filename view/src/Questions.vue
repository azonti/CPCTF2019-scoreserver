<template>
  <div>
    <vue-headful title="Questions | CPCTF2018" />
    <div v-if="!loading">
      <div class="row">
        <div v-for="question in questions" class="col-md-6">
          <div class="panel" :class="question.answer ? 'panel-success' : 'panel-primary'">
            <div class="panel-body">
              <dt class="row">
                <dt class="col-xs-4">Questioner</dt>
                <dd class="col-xs-8"><router-link :to="{name: 'user', params: {id: question.questioner.id}}"><img :src="question.questioner.icon_url" class="icon">{{ question.questioner.name }}<small v-if="question.questioner.twitter_screen_name">(@{{ question.questioner.twitter_screen_name }})</small></router-link></dd>
              </dt>
              <dt class="row">
                <dt class="col-xs-4">Question</dt>
                <dd class="col-xs-8">{{ question.query }}</dd>
              </dt>
              <dt class="row">
                <dt class="col-xs-4">Answerer</dt>
                <dd v-if="question.answerer" class="col-xs-8"><router-link :to="{name: 'user', params: {id: question.answerer.id}}"><img :src="question.answerer.icon_url" class="icon">{{ question.answerer.name }}<small v-if="question.answerer.twitter_screen_name">(@{{ question.answerer.twitter_screen_name }})</small></router-link></dd>
              </dt>
              <dt class="row">
                <dt class="col-xs-4">Answer</dt>
                <dd class="col-xs-8">{{ question.answer }}</dd>
              </dt>
            </div>
          </div>
        </div>
        <div class="col-md-6">
          <div class="panel panel-primary">
            <div class="panel-body">
              <div class="col-xs-10">
                <textarea class="form-control" v-model="query" placeholder="New question..."></textarea>
              </div>
              <div class="col-xs-2">
                <button v-if="!sendingQuestion" @click="sendQuestion" class="btn btn-primary" style="width: 100%;">Send</button>
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

<style>
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
      loading: true,
      sendingQuestion: false,
      questions: [],
      query: ""
    }
  },
  created () {
    this.fetchQuestions()
  },
  methods: {
    fetchQuestions () {
      api.get(`${process.env.API_URL_PREFIX}/questions`)
      .then(res => res.data)
      .then((data) => {
        this.questions.splice(0, data.length, ...data)
        this.loading = false
      })
    },
    sendQuestion () {
      this.sendingQuestion = true
      api.post(`${process.env.API_URL_PREFIX}/questions`, {
        questioner: this.me,
        query: this.query
      })
      .then(() => {
        this.query = ""
        this.$emit('success', 'Your question has been sent.')
      })
      .then(() => this.fetchQuestions())
      .catch((err) => {
        this.$emit('error', `Message: ${err.response.data.message}`)
      })
      .then(() => {
        this.sendingQuestion = false
      })
    }
  }
}
</script>

<style>
</style>
