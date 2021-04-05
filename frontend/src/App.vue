<template>
  <div id="app">
      <b-navbar type="dark" variant="info">
          <b-navbar-brand>Djavue</b-navbar-brand>
          <b-navbar-nav class="ml-auto">
        <b-nav-form>
          <b-button v-if="loggedIn" size="sm" class="my-2 my-sm-0" @click="logOut()">Logout</b-button>
        </b-nav-form>
          </b-navbar-nav>
      </b-navbar>
      <LoginForm  v-if="!loggedIn" />
      <SubmitJobForm v-if="loggedIn" />

 </div>
</template>

<script>
import { mapState } from 'vuex'
import LoginForm from './components/LoginForm.vue'
import SubmitJobForm from './components/SubmitJobForm.vue'

export default {
  name: 'Djavule',
  computed: {
      ...mapState(['loggedIn'])
  },
  methods: {
      logOut() {
          this.$store.dispatch('logOut')
      }
  },
  components: {
      LoginForm,
      SubmitJobForm
  }
}
</script>

