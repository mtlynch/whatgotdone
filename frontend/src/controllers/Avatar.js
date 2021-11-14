import axios from 'axios';

import {getCsrfToken} from '@/controllers/Common.js';

export function deleteAvatar() {
  return new Promise(function (resolve, reject) {
    axios
      .delete(`${process.env.VUE_APP_BACKEND_URL}/api/user/avatar`, {
        withCredentials: true,
        headers: {
          'X-CSRF-Token': getCsrfToken(),
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
