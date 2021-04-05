<template>
    <div>
        <b-contaier>
         <b-form @submit="onSubmit" class="mx-auto" id="login-form">
            <b-form-group label="Your Name:">
                <b-form-input v-model="username" required />
            </b-form-group>
            <b-form-group label="Password:">
                <b-form-input v-model="password" required type="password" />
            </b-form-group>
            <div class="col-md-12 text-center">
              <b-button size="lg" type="submit" variant="outline-primary">Login</b-button>
            </div>
            <div v-if="error_msg">{{error_msg}}</div>
         </b-form>
        </b-contaier>
    </div>
</template>

<script>

import axios from 'axios'

const TOKEN_URL = 'http://127.0.0.1:8000/api/token/'

export default {
    name: 'LoginForm',
    data() {
        return {
            "username": "",
            "password": "",
            "error_msg": ""
        }
    },
    methods: {
        onSubmit(event) {
            event.preventDefault()
            axios.post(TOKEN_URL, {
                username: this.username,
                password: this.password
            }).then(
                response => {
                    this.error_msg = response.data
                },
                error => {
                    this.error_msg = error
                })
        }
    }
}
</script>

<style scoped >

  #login-form {
     width: 50%;
     margin-top: 10px
 }

</style>
