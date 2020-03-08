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
