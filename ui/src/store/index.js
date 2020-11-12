import Vue from 'vue'
import Vuex from 'vuex'
import auth from "./modules/auth";
import client from "./modules/client";
import server from "./modules/server";
import status from "./modules/status";

Vue.use(Vuex)

export default new Vuex.Store({
  state: {},
  getters : {},
  mutations: {},
  actions:{},
  modules: {
    auth,
    client,
    server,
    status,
  }
})
