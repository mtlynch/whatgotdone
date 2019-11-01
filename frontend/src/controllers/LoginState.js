import axios from "axios";
import store from "../store.js";
import { logoutUserKit }  from "../controllers/UserKit.js";

function clearCachedAuthInformation() {
  store.commit("clearUsername");
  logoutUserKit();
}

export default function updateLoginState(attempts, callback) {
  if (attempts <= 0) {
    return;
  }
  const url = `${process.env.VUE_APP_BACKEND_URL}/api/user/me`;
  axios
    .get(url, { withCredentials: true })
    .then(result => {
      store.commit("setUsername", result.data.username);
      if(typeof callback === 'function') {
        callback();
      }
    })
    .catch((error) => {
      // If checking user information fails, the cached authentication information
      // is no longer correct, so we need to clear it.
      if (error.response && error.response.status === 403) {
        clearCachedAuthInformation();
        return;
      }
      updateLoginState(attempts - 1, callback);
    });
}