import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import vuetify from './plugins/vuetify';
import './plugins/moment';
import './plugins/cidr'
import './plugins/axios'

// Don't warn about using the dev version of Vue in development.
Vue.config.productionTip = process.env.NODE_ENV === 'production'

new Vue({
  router,
  store,
  vuetify,
  render: function (h) { return h(App) }
}).$mount('#app')
