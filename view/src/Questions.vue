<template>
  <div>
    <vue-headful title="Questions | CPCTF2018" />
    <div class="row">
      <div v-for="question in questions" class="col-md-6">
        <div class="panel" :class="question.answer ? 'panel-success' : 'panel-primary'">
          <div class="panel-body">
            <div v-if="question.questioner.id !== 'nobody'">
              <dl class="row">
                <dt class="col-xs-2">Questioner</dt>
                <dd class="col-xs-10"><router-link :to="{name: 'user', params: {id: question.questioner.id}}"><img :src="question.questioner.icon_url" class="icon">{{ question.questioner.name }}<small v-if="question.questioner.twitter_screen_name">(@{{ question.questioner.twitter_screen_name }})</small></router-link></dd>
              </dl>
              <dl class="row">
                <dt class="col-xs-2">Question</dt>
                <dd class="col-xs-10">{{ question.query }}</dd>
              </dl>
            </div>
            <div v-if="question.answer || !me.is_author">
              <dl class="row">
                <dt class="col-xs-2">Answerer</dt>
                <dd v-if="question.answerer" class="col-xs-10"><router-link :to="{name: 'user', params: {id: question.answerer.id}}"><img :src="question.answerer.icon_url" class="icon">{{ question.answerer.name }}<small v-if="question.answerer.twitter_screen_name">(@{{ question.answerer.twitter_screen_name }})</small></router-link></dd>
              </dl>
              <dl class="row">
                <dt class="col-xs-2">Answer</dt>
                <dd class="col-xs-10">{{ question.answer }}</dd>
              </dl>
            </div>
            <div v-else>
              <div class="col-xs-10">
                <textarea class="form-control" v-model="question._answer" placeholder="Answer..."></textarea>
              </div>
              <div class="col-xs-2">
                <button v-if="!answeringQuestion" @click="questionToAnswer = question; answerQuestion()" class="btn btn-primary" style="width: 100%;">Send</button>
              </div>
            </div>
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
    'me',
    'questions'
  ],
  data () {
    return {
      sendingQuestion: false,
      answeringQuestion: false,
      questionToAnswer: {},
      query: ""
    }
  },
  methods: {
    sendQuestion () {
      this.sendingQuestion = true
      api.post(`${process.env.API_URL_PREFIX}/questions`, {
        questioner: this.me,
        query: this.query
      })
      .then(() => {
        this.query = ""
        this.$emit('success', 'Your question has been sent.')
        this.$emit('reloadQuestions')
      })
      .catch((err) => {
        this.$emit('error', `Message: ${err.response.data.message}`)
      })
      .then(() => {
        this.sendingQuestion = false
      })
    },
    answerQuestion () {
      this.answeringQuestion = true
      api.put(`${process.env.API_URL_PREFIX}/questions/${this.questionToAnswer.id}`, {
        questioner: this.questionToAnswer.questioner,
        answerer: this.me,
        query: this.questionToAnswer.query,
        answer: this.questionToAnswer._answer
      })
      .then(() => {
        this.$emit('success', 'Your answer has been sent.')
        this.$emit('reloadQuestions')
      })
      .catch((err) => {
        this.$emit('error', `Message: ${err.response.data.message}`)
      })
      .then(() => {
        this.answeringQuestion = false
      })
    }
  }
}
</script>

<style>
</style>
