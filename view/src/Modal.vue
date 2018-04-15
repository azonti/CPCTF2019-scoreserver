<template>
  <transition name="fade">
    <div class="modal" v-show="show" style="display: block;">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <button type="button" class="close" :aria-hidden="show ? 'false' : 'true'" @click="$emit('close');">&times;</button>
            <h3 class="modal-title" :class="modal.bodyClass">{{ modal.title }}</h3>
          </div>
          <div class="modal-body">
            <p :class="modal.bodyClass">{{ modal.body }}</p>
          </div>
          <div class="modal-footer">
            <button v-if="modal.showCancel" type="button" class="btn btn-default" @click="$emit('close');">Cancel</button>
            <button type="button" class="btn" :class="modal.btnClass || 'btn-primary'" @click="callback && callback(); $emit('close');">{{ modal.btnBody || 'OK' }}</button>
          </div>
        </div>
      </div>
    </div>
  </transition>
</template>

<script>
import axios from 'axios'
const api = axios.create({
  withCredentials: true
})

export default {
  props: [
    'show',
    'modal',
    'callback'
  ]
}
</script>

<style>
.fade-enter-active, .fade-leave-active {
  display: block;
  transition: opacity .5s;
}
.fade-enter, .fade-leave-to {
  opacity: 0;
}

</style>
