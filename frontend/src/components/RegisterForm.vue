<template>
    <div>
        <b-container>
         <b-form @submit="onSubmit" class="mx-auto" id="register-form">
            <b-form-group label="Your Username:">
                <b-form-input v-model="username" required />
            </b-form-group>
            <b-form-group label="Password:">
                <b-form-input v-model="password0" required type="password" />
                <b-form-input placeholder="Repeat password..." v-model="password1" required type="password" />
            </b-form-group>
            <b-form-group label="Invitation Code:">
                <b-form-input v-model="invite_code" required />
            </b-form-group>
            <div class="col-md-12 text-center">
              <b-button size="lg" type="submit" variant="outline-primary">Register</b-button>
            </div>
             <b-alert
                :show="dismissCountDown"
                dismissible
                variant="warning"
                @dismissed="dismissCountDown=0"
                @dismiss-count-down="countDownChanged"
                > {{errorMsg}}</b-alert>
         </b-form>
        </b-container>
    </div>
</template>

<script>

import axios from 'axios'
import { mapActions } from 'vuex'

const USER_URL = process.env.VUE_APP_API_URL + 'user/'

export default {
    name: 'RegisterForm',
    data() {
        return {
            dismissCountDown: 0,
            dismissSecs: 10,
            errorMsg: "",
            username: "",
            password0: "",
            password1: "",
            invite_code: "",
        }
    },
    methods: {
        ...mapActions(['setToken']),
        countDownChanged(dismissCountDown) {
            this.dismissCountDown = dismissCountDown
        },
        showAlert(err) {
            this.errorMsg = err
            this.dismissCountDown = this.dismissSecs
        },
        onSubmit(event) {
            event.preventDefault()
            if (this.password0 != this.password1) {
                this.showAlert("Password do not match.")
                return
            }
            axios.post(USER_URL, {
                username: this.username,
                password: this.password0,
                code: this.invite_code
            }).then(
                response => {
                    this.setToken(response.data.token)
                    this.$router.push('/')
                },
                error => {
                    this.showAlert(error.response.data.error)
                })
        }
    }
}
</script>

<style scoped >

  #register-form {
     width: 60%;
     margin-top: 10px
 }

</style>
