import Vue from 'vue'
import Router from 'vue-router'
import SubmitJobForm from './../components/SubmitJobForm.vue'
import LoginForm from './../components/LoginForm.vue'
import RegisterForm from './../components/RegisterForm.vue'


Vue.use(Router)

export default new Router({
    routes: [
        { path: '/', component: SubmitJobForm, name: 'home' },
        { path: '/login', component: LoginForm, name: 'login' },
        { path: '/register', component: RegisterForm, name: 'register' }
    ]
})


