<template>
    <b-img class="auth-image" :src="image_src_data" fluid />
</template>

<script>

import axios from 'axios'

const IMAGE_URL = 'http://127.0.0.1:8000/api/image/'

export default {
    name: 'AuthImage',
    props: [
        'pk'
    ],
    data() {
        return {
            image_src_data:""
        }
    },
    methods: {
        fetchSrcData() {
            let url = IMAGE_URL + this.pk + '/'
            let auth = 'Token ' + this.$store.state.token
            axios({
                method: 'get',
                url: url,
                headers: {
                    'Authorization': auth
                },
                responseType: 'arraybuffer'
            }).then(resp => {
                let mimeType = resp.headers['content-type'].toLowerCase();
                let imgBase64 = new Buffer(resp.data, 'binary').toString('base64');
                this.image_src_data = 'data:' + mimeType + ';base64,' + imgBase64;
            }).catch(error => {
                console.log(error.response.data)
                })
        }
    },
    created() {
        this.fetchSrcData()
    }
}
</script>

<style scoped>
 .auth-image {
     border-style: dotted;
 }
</style>

