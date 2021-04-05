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
        </b-container>
    </div>
</template>

<script>

import axios from 'axios'

const JOB_URL = 'http://127.0.0.1:8000/api/job/'

export default {
    name: 'SubmitJobForm',
    data() {
        return {
            'error_msg': '',
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
                    console.log(response.data)
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
 }
 #kind-select {
     margin-top: 10px;
     margin-bottom: 20px;
 }

</style>
