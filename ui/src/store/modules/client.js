import ApiService from "../../services/api.service";

const state = {
  error: null,
  clients: [],
  clientQrcodes: [],
  clientConfigs: []
}

const getters = {
  error(state) {
    return state.error;
  },
  clients(state) {
    return state.clients;
  },
  getClientQrcode: (state) => (id) => {
    let item = state.clientQrcodes.find(item => item.id === id)
    // initial load fails, must wait promise and stuff...
    return item ? item.qrcode : "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg=="
  },
  getClientConfig: (state) => (id) => {
    let item = state.clientConfigs.find(item => item.id === id)
    return item ? item.config : null
  }
}

const actions = {
  error({ commit }, error){
    commit('error', error)
  },

  readAll({ commit, dispatch }){
    ApiService.get("/client")
      .then(resp => {
        commit('clients', resp)
        dispatch('readQrcodes')
        dispatch('readConfigs')
      })
      .catch(err => {
        commit('error', err)
      })
  },

  create({ commit, dispatch }, client){
    ApiService.post("/client", client)
      .then(resp => {
        dispatch('readQrcode', resp)
        dispatch('readConfig', resp)
        commit('create', resp)
      })
      .catch(err => {
        commit('error', err)
      })
  },

  update({ commit, dispatch }, client){
    ApiService.patch(`/client/${client.id}`, client)
      .then(resp => {
        dispatch('readQrcode', resp)
        dispatch('readConfig', resp)
        commit('update', resp)
      })
      .catch(err => {
        commit('error', err)
      })
  },

  delete({ commit }, client){
    ApiService.delete(`/client/${client.id}`)
      .then(() => {
        commit('delete', client)
      })
      .catch(err => {
        commit('error', err)
      })
  },

  email({ commit }, client){
    ApiService.get(`/client/${client.id}/email`)
      .then(() => {
      })
      .catch(err => {
        commit('error', err)
      })
  },

  readQrcode({ state, commit }, client){
    ApiService.getWithConfig(`/client/${client.id}/config?qrcode=true`, {responseType: 'arraybuffer'})
      .then(resp => {
        let image = Buffer.from(resp, 'binary').toString('base64')
        commit('clientQrcodes', { client, image })
      })
      .catch(err => {
        commit('error', err)
      })
  },

  readConfig({ state, commit }, client){
    ApiService.getWithConfig(`/client/${client.id}/config?qrcode=false`, {responseType: 'arraybuffer'})
      .then(resp => {
        commit('clientConfigs', { client: client, config: resp })
      })
      .catch(err => {
        commit('error', err)
      })
  },

  readQrcodes({ state, dispatch }){
    state.clients.forEach(client => {
      dispatch('readQrcode', client)
    })
  },

  readConfigs({ state, dispatch }){
    state.clients.forEach(client => {
      dispatch('readConfig', client)
    })
  },
}

const mutations = {
  error(state, error) {
    state.error = error;
  },
  clients(state, clients){
    state.clients = clients
  },
  create(state, client){
    state.clients.push(client)
  },
  update(state, client){
    let index = state.clients.findIndex(x => x.id === client.id);
    if (index !== -1) {
      state.clients.splice(index, 1);
      state.clients.push(client);
    } else {
      state.error = "update client failed, not in list"
    }
  },
  delete(state, client){
    let index = state.clients.findIndex(x => x.id === client.id);
    if (index !== -1) {
      state.clients.splice(index, 1);
    } else {
      state.error = "delete client failed, not in list"
    }
  },
  clientQrcodes(state, { client, image }){
    let index = state.clientQrcodes.findIndex(x => x.id === client.id);
    if (index !== -1) {
      state.clientQrcodes.splice(index, 1);
    }
    state.clientQrcodes.push({
      id: client.id,
      qrcode: image
    })
  },
  clientConfigs(state, { client, config }){
    let index = state.clientConfigs.findIndex(x => x.id === client.id);
    if (index !== -1) {
      state.clientConfigs.splice(index, 1);
    }
    state.clientConfigs.push({
      id: client.id,
      config: config
    })
  },
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
