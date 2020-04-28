const TOKEN_KEY = "token";
const CLIENT_ID_KEY = "client_id";

export const getToken = () => {
  return window.localStorage.getItem(TOKEN_KEY);
};

export const saveToken = token => {
  window.localStorage.setItem(TOKEN_KEY, token);
};

export const destroyToken = () => {
  window.localStorage.removeItem(TOKEN_KEY);
};

export const getClientId = () => {
  return window.localStorage.getItem(CLIENT_ID_KEY);
};

export const saveClientId = token => {
  window.localStorage.setItem(CLIENT_ID_KEY, token);
};

export const destroyClientId = () => {
  window.localStorage.removeItem(CLIENT_ID_KEY);
};

export default {
  getToken,
  saveToken,
  destroyToken,
  getClientId,
  saveClientId,
  destroyClientId
};
