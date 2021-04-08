<template>
  <div id="app">
      <b-navbar type="dark" variant="info">
          <b-navbar-brand>Djavue</b-navbar-brand>
          <b-navbar-nav class="ml-auto">
            <b-nav-form>
                <b-button v-if="loggedIn" size="sm" class="my-2 my-sm-0" @click="logOut()">Logout</b-button>
                <b-button to="/register" v-if="showRegisterBtn" size="sm" class="my-2 my-sm-0">Register</b-button>
                <b-button to="/login" v-if="showLoginBtn" size="sm" class="my-2 my-sm-0">Login</b-button>
            </b-nav-form>
          </b-navbar-nav>
      </b-navbar>
      <router-view />
 </div>
</template>

<script>
import { mapState } from 'vuex'

export default {
  name: 'Djavule',
  computed: {
      ...mapState(['loggedIn']),
      showRegisterBtn(){
          if (this.$route.name == 'register')
              return false
          if (this.loggedIn)
              return false
          return true
      },
      showLoginBtn() {
          if (this.loggedIn)
            return false
          if (this.$route.name == 'login')
              return false
          return true
      }
  },
  beforCreate() {
      if (!this.loggedIn) {
          this.$router.push('/login')
      }
  },
  methods: {
      logOut() {
          this.$store.dispatch('logOut')
          this.$router.push('/Login')
      }
  },
}
</script>

