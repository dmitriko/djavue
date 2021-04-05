import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

let initialState = {
    "token": "",
    "loggedIn": false
}
let token = localStorage.getItem("token")
if (token) {
    initialState.token = token
    initialState.loggedIn = true
}

export default new Vuex.Store({
    state: initialState,

    mutations: {
        SET_TOKEN(state, token) {
            state.token = token
        },
        SET_LOGGED_IN(state, value) {
            state.loggedIn = value
        }
    },

    actions: {
        setToken({ commit }, token) {
            commit('SET_TOKEN', token)
            commit('SET_LOGGED_IN', true)
            localStorage.setItem('token', token)
        },
        logOut({ commit }) {
            commit('SET_TOKEN', '')
            commit('SET_LOGGED_IN', false)
            localStorage.setItem('token', '')
        }
    }

})

