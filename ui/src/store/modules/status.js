import ApiService from "../../services/api.service";

const state = {
  error: null,
  enabled: false,
  interfaceStatus: null,
  clientStatus: [],
  version: '_ci_build_not_run_properly_',
}

const getters = {
  error(state) {
    return state.error;
  },

  enabled(state) {
    return state.enabled;
  },

  interfaceStatus(state) {
    return state.interfaceStatus;
  },

  clientStatus(state) {
    return state.clientStatus;
  },

  version(state) {
    return state.version;
  },
}

const actions = {
  error({ commit }, error){
    commit('error', error)
  },

  read({ commit }){
    ApiService.get("/status/interface")
      .then(resp => {
        commit('interfaceStatus', resp)
      })
      .catch(err => {
        commit('interfaceStatus', null);
        commit('error', err)
      });
    ApiService.get("/status/clients")
      .then(resp => {
        commit('clientStatus', resp)
      })
      .catch(err => {
        commit('clientStatus', []);
        commit('error', err)
      });
  },

  isEnabled({ commit }){
    ApiService.get("/status/enabled")
      .then(resp => {
        commit('enabled', resp)
      })
      .catch(err => {
        commit('enabled', false);
        commit('error', err.response.data)
      });
  },
}

const mutations = {
  error(state, error) {
    state.error = error;
  },

  enabled(state, enabled) {
    state.enabled = enabled;
  },

  interfaceStatus(state, interfaceStatus){
    state.interfaceStatus = interfaceStatus
  },

  clientStatus(state, clientStatus){
    state.clientStatus = clientStatus
  },

  version(state, version){
    state.version = version
  },
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
