<template>
  <div id="app">
      <b-navbar type="dark" variant="info">
          <b-navbar-brand>Djavue</b-navbar-brand>
      </b-navbar>
      <LoginForm  v-if="!isLogged" />
      <div v-if="isLogged">SubmitJob Form</div>

 </div>
</template>

<script>

import LoginForm from './components/LoginForm.vue'

export default {
  name: 'Djavule',
  computed: {
      isLogged() {
          return false
      }
  },
  components: {
      LoginForm
  }
}
</script>

