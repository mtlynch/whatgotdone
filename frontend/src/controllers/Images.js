import axios from 'axios';

import getCsrfToken from '@/controllers/CsrfToken.js';

export function uploadImage(image) {
  return new Promise(function(resolve, reject) {
    let formData = new FormData();
    formData.append('file', image);
    let url = `${process.env.VUE_APP_BACKEND_URL}/api/images`;
    axios
      .put(url, formData, {
        withCredentials: true,
        headers: {
          'X-CSRF-Token': getCsrfToken(),
          'Content-Type': 'multipart/form-data',
        },
      })
      .then(response => {
        resolve(response.data.url);
      })
      .catch(err => {
        reject(err);
      });
  });
}
