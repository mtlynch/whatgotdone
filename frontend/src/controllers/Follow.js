import axios from 'axios';

export function getFollowing(username) {
  return new Promise(function(resolve, reject) {
    const url = `${process.env.VUE_APP_BACKEND_URL}/api/user/${username}/following`;
    axios
      .get(url)
      .then(result => {
        resolve(result.data.following);
      })
      .catch(err => reject(err));
  });
}
