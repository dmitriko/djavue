<template>
    <div>
        <b-container>
            <div id="submit-job-form" class="mx-auto" >
             <b-form-file
                accept="image/*"
                placeholder="Choose a file or drop it here..."
                drop-placeholder="Drop file here..."
                v-model="file"
                @change="clearErrMsg"
                ></b-form-file>
             <b-form-select @change="clearErrMsg" id="kind-select" v-model="kind" :options="kind_options" />
                  <div class="col-md-12 text-center">
                      <b-button @click="submitJob" :disabled="loading">
                          <b-spinner small label="Loading..." v-if="loading" ></b-spinner>
                          Submit
                      </b-button>
                 </div>
                 <div v-if="error_msg">{{error_msg}}</div>
            </div>
                 <AuthImage v-for="image in result_images" :key="image.pk" :pk="image.pk" />
        </b-container>
    </div>
</template>

<script>
import AuthImage from './AuthImage.vue'

import axios from 'axios'

const JOB_URL = process.env.VUE_APP_API_URL + 'job/'

export default {
    name: 'SubmitJobForm',
    components: {
        AuthImage
    },
    beforeCreate() {
        if (!this.$store.state.loggedIn) {
            this.$router.push('/login')
        }
    },
    data() {
        return {
            'error_msg': '',
            'result_images': [],
            'file': null,
            'kind': null,
            'loading': false,
            'kind_options': [
                 { value: null, text: 'Please, select a kind of a job' },
                 { value: 'original', text: 'Save original' },
                 { value: 'square_original', text: 'Make a big square' },
                 { value: 'square_small', text: 'Make a small square' },
                 { value: 'all_three', text: 'Save all three' }
        ]
        }
    },
    methods: {
        clearErrMsg() {
            this.error_msg = ''
        },
        fetchResults(job_id) {
            axios.get(JOB_URL + job_id + '/',
                {headers: {'Authorization': 'Token ' + this.$store.state.token}}
            ).then(response => {
                this.loading = false
                this.result_images = response.data.images
            }).catch(error => {
                this.loading = false
                this.error_msg = error.response.data})
        },
        submitJob() {
            if (!this.kind) {
                this.error_msg = 'Please, select a job kind'
                return
            }
            if (!this.file) {
                this.error_msg = 'Please, choose an image from your computer'
                return
            }
            let auth = 'Token ' + this.$store.state.token

            let formData = new FormData()
            formData.append('kind', this.kind)
            formData.append('file', this.file)
            this.loading = true
            axios.post(JOB_URL,
                formData,
                  {
                    headers: {
                    'Content-Type': 'multipart/form-data',
                    'Authorization': auth
                    }
                }).then(response => {
                    this.fetchResults(response.data.job_id)
                }).catch(error => {
                    this.loading = false
                    if (error.response.data && error.response.data.detail == 'Invalid token.'){
                        this.$store.dispatch('logOut')
                        this.$router.push('/login')
                        return
                    }
                    this.error_msg = error.response.data
                })

        }
    }
}
</script>

<style scoped >

  #submit-job-form {
     width: 50%;
     margin-top: 10px;
     margin-bottom: 50px;
 }
 #kind-select {
     margin-top: 10px;
     margin-bottom: 20px;
 }

</style>
