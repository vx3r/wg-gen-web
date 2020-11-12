import Vue from "vue";
import TokenService from "./token.service";

const ApiService = {

  setHeader() {
    Vue.axios.defaults.headers['x-wg-gen-web-auth'] = `${TokenService.getToken()}`;
  },

  get(resource) {
    return Vue.axios.get(resource)
      .then(response => response.data)
        .catch(error => {
          if(typeof error.response !== 'undefined') {
            throw new Error(`${error.response.status} - ${error.response.statusText}: ${error.response.data}`)
          } else {
            throw new Error(`ApiService: ${error}`)
          }
        });
  },

  post(resource, params) {
    return Vue.axios.post(resource, params)
      .then(response => response.data)
      .catch(error => {
        throw new Error(`ApiService: ${error}`)
      });
  },

  put(resource, params) {
    return Vue.axios.put(resource, params)
      .then(response => response.data)
      .catch(error => {
        throw new Error(`ApiService: ${error}`)
      });
  },

  patch(resource, params) {
    return Vue.axios.patch(resource, params)
      .then(response => response.data)
      .catch(error => {
        throw new Error(`ApiService: ${error}`)
      });
  },

  delete(resource) {
    return Vue.axios.delete(resource)
      .then(response => response.data)
      .catch(error => {
        throw new Error(`ApiService: ${error}`)
      });
  },

  getWithConfig(resource, config) {
    return Vue.axios.get(resource, config)
      .then(response => response.data)
      .catch(error => {
        throw new Error(`ApiService: ${error}`)
      });
  },
};

export default ApiService;
