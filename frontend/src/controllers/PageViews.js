import axios from 'axios';

export function getPageViews(path) {
  return new Promise(function(resolve, reject) {
    const url = `${process.env.VUE_APP_BACKEND_URL}/api/pageViews`;
    axios
      .get(url, {
        params: {
          path: path,
        },
      })
      .then(result => {
        if (result.data) {
          resolve(result.data.views);
        } else {
          resolve(null);
        }
      })
      .catch(error => {
        reject(error);
      });
  });
}
