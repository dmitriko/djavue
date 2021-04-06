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
                      <b-button @click="submitJob">Submit</b-button>
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

const JOB_URL = process.env.VUE_APP_JOB_URL

export default {
    name: 'SubmitJobForm',
    components: {
        AuthImage
    },
    data() {
        return {
            'error_msg': '',
            'result_images': [],
            'file': null,
            'kind': null,
            'kind_options': [
                 { value: null, text: 'Please select a job kind' },
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
                this.result_images = response.data.data.images
            }).catch(error => {this.error_msg = error.response.data})
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
                    if (error.response.data && error.response.data.detail == 'Invalid token.'){
                        this.$store.dispatch('logOut')
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
