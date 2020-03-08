import axios from 'axios';

export function getUserSelfMetadata() {
  return new Promise(function(resolve, reject) {
    const url = `${process.env.VUE_APP_BACKEND_URL}/api/user/me`;
    axios
      .get(url, {withCredentials: true})
      .then(result => {
        resolve(result.data);
      })
      .catch(error => {
        reject(error);
      });
  });
}

export function getUserMetadata(username) {
  return new Promise(function(resolve, reject) {
    const url = `${process.env.VUE_APP_BACKEND_URL}/api/user/${username}`;
    axios
      .get(url)
      .then(result => {
        resolve(result.data);
      })
      .catch(error => {
        if (error.response && error.response.status == 404) {
          // A 404 for a user profile is equivalent to retrieving an empty profile.
          resolve({});
        } else {
          reject(error);
        }
      });
  });
}
