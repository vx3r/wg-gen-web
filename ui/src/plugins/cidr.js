import Vue from 'vue'
const isCidr = require('is-cidr');

const plugin = {
  install () {
    Vue.isCidr = isCidr;
    Vue.prototype.$isCidr = isCidr
  }
};

Vue.use(plugin);
