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
              <b-button size="lg" type="submit" variant="primary">Login</b-button>
            </div>
         </b-form>
        </b-contaier>
    </div>
</template>

<script>
export default {
    name: 'LoginForm',
    data() {
        return {
            "username": "",
            "password": ""
        }
    },
    methods: {
        onSubmit(event) {
            event.preventDefault()
            alert(this.username)
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
