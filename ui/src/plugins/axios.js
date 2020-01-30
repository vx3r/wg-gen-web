import Vue from 'vue'
import VueAxios from 'vue-plugin-axios'
import axios from 'axios'

export const myVar = 'This is my variable'

// https://www.npmjs.com/package/vue-cli-plugin-vuetify
Vue.use(VueAxios, {
  axios,
  config: {
    baseURL: process.env.VUE_APP_API_BASE_URL || '/api/v1.0',
  },
});
