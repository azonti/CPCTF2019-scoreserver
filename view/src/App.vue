<template>
  <div id="app" class="container">
    <ul class="nav nav-tabs">
      <li :class="isActive('/challenges') ? 'active' : ''"><router-link :to="{name: 'challenges'}">Challenges</router-link></li>
      <li :class="isActive('/challenges') ? 'active' : ''"><router-link :to="{name: 'challenges'}">Questions</router-link></li>
      <li :class="isActive('/challenges') ? 'active' : ''"><router-link :to="{name: 'challenges'}">Ranking</router-link></li>
      <li class="dropdown">
        <a :aria-expanded="showDropdown ? 'true' : 'false'" class="dropdown-toggle" v-on:click.prevent="showDropdown = !showDropdown" href="#">Me <span class="caret"></span></a>
        <ul class="dropdown-menu" v-bind:style="{ display: showDropdown ? 'block' : 'none' }">
          <li class="divider"></li>
          <li :class="loading || me.id ? 'disabled' : ''"><a v-on:click="showDropdown = !showDropdown" href="#">Login with Twitter</a></li>
          <li :class="loading || !me.id ? 'disabled' : ''"><a v-on:click="showDropdown = !showDropdown" href="#">Logout</a></li>
        </ul>
      </li>
    </ul>
    <router-view>
    </router-view>
  </div>
</template>

<script>
export default {
  data() {
    return {
      loading: true,
      showDropdown: false,
      me: {}
    }
  },
  methods: {
    isActive(path) {
      return this.$route.path === path
    }
  }
}
</script>

<style>
</style>
