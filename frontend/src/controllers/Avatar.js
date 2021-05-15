import axios from 'axios';

import getCsrfToken from '@/controllers/CsrfToken.js';

export function uploadAvatar(image) {
  return new Promise(function (resolve, reject) {
    const formData = new FormData();
    formData.append('file', image);
    axios
      .put(`${process.env.VUE_APP_BACKEND_URL}/api/user/avatar`, formData, {
        withCredentials: true,
        headers: {
          'X-CSRF-Token': getCsrfToken(),
          'Content-Type': 'multipart/form-data',
        },
      })
      .then(() => {
        resolve();
      })
      .catch((error) => {
        reject(error);
      });
  });
}
