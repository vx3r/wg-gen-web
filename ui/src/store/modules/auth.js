import ApiService from "../../services/api.service";
import TokenService from "../../services/token.service";

const state = {
  error: null,
  user: null,
  authStatus: '',
  authRedirectUrl: '',
};

const getters = {
  error(state) {
    return state.error;
  },
  user(state) {
    return state.user;
  },
  isAuthenticated(state) {
    return state.user !== null;
  },
  authRedirectUrl(state) {
    return state.authRedirectUrl
  },
  authStatus(state) {
    return state.authStatus
  },
};

const actions = {
  user({ commit }){
    ApiService.get("/auth/user")
      .then( resp => {
        commit('user', resp)
      })
      .catch(err => {
        commit('error', err);
        commit('logout')
      });
  },

  oauth2_url({ commit, dispatch }){
    if (TokenService.getToken()) {
      ApiService.setHeader();
      dispatch('user');
      return
    }
    ApiService.get("/auth/oauth2_url")
      .then(resp => {
        if (resp.codeUrl === '_magic_string_fake_auth_no_redirect_'){
          console.log("server report oauth2 is disabled, fake exchange")
          commit('authStatus', 'disabled')
          TokenService.saveClientId(resp.clientId)
          dispatch('oauth2_exchange', {code: "", state: resp.state})
        } else {
          commit('authStatus', 'redirect')
          commit('authRedirectUrl', resp)
        }
      })
      .catch(err => {
        commit('authStatus', 'error')
        commit('error', err);
        commit('logout')
      })
  },

  oauth2_exchange({ commit, dispatch }, data){
    data.clientId = TokenService.getClientId()
    ApiService.post("/auth/oauth2_exchange", data)
      .then(resp => {
        commit('authStatus', 'success')
        commit('token', resp)
        dispatch('user');
      })
      .catch(err => {
        commit('authStatus', 'error')
        commit('error', err);
        commit('logout')
      })
  },

  logout({ commit }){
    ApiService.get("/auth/logout")
      .then(resp => {
        commit('logout')
      })
      .catch(err => {
        commit('authStatus', '')
        commit('error', err);
        commit('logout')
      })
  },
}

const mutations = {
  error(state, error) {
    state.error = error;
  },
  authStatus(state, authStatus) {
    state.authStatus = authStatus;
  },
  authRedirectUrl(state, resp) {
    state.authRedirectUrl = resp.codeUrl;
    TokenService.saveClientId(resp.clientId);
  },
  token(state, token) {
    TokenService.saveToken(token);
    ApiService.setHeader();
    TokenService.destroyClientId();
  },
  user(state, user) {
    state.user = user;
  },
  logout(state) {
    state.user = null;
    TokenService.destroyToken();
    TokenService.destroyClientId();
  }
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
