import axios from 'axios';

import getCsrfToken from '@/controllers/CsrfToken.js';

export function uploadMedia(media) {
  return new Promise(function (resolve, reject) {
    let formData = new FormData();
    formData.append('file', media);
    let url = `${process.env.VUE_APP_BACKEND_URL}/api/media`;
    axios
      .put(url, formData, {
        withCredentials: true,
        headers: {
          'X-CSRF-Token': getCsrfToken(),
          'Content-Type': 'multipart/form-data',
        },
      })
      .then((response) => {
        resolve(response.data.url);
      })
      .catch((err) => {
        reject(err);
      });
  });
}
