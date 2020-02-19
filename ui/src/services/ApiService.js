import axios from 'axios'

let baseUrl = "/api/v1.0";
if (process.env.NODE_ENV === "development"){
  baseUrl = process.env.VUE_APP_API_BASE_URL
}

export const API_BASE_URL = baseUrl;

export class ApiService {
  get(resource) {
    return axios
      .get(`${API_BASE_URL}${resource}`)
      .then(response => response.data)
  };

  post(resource, data) {
    return axios
      .post(`${API_BASE_URL}${resource}`, data)
      .then(response => response.data)
  };

  put(resource, data) {
    return axios
      .put(`${API_BASE_URL}${resource}`, data)
      .then(response => response.data)
  };

  patch(resource, data) {
    return axios
      .patch(`${API_BASE_URL}${resource}`, data)
      .then(response => response.data)
  };

  delete(resource) {
    return axios
      .delete(`${API_BASE_URL}${resource}`)
      .then(response => response.data)
  };
}
