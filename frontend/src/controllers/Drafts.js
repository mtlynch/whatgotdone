import axios from 'axios';

import {getCsrfToken} from '@/controllers/Common.js';

export function getDraft(entryDate) {
  return new Promise(function (resolve, reject) {
    const url = `${process.env.VUE_APP_BACKEND_URL}/api/draft/${entryDate}`;
    axios
      .get(url, {withCredentials: true})
      .then((result) => {
        if (!Object.prototype.hasOwnProperty.call(result, 'data')) {
          return Promise.reject(
            new Error('response is missing expected field: data')
          );
        }
        if (!Object.prototype.hasOwnProperty.call(result.data, 'markdown')) {
          return Promise.reject(
            new Error('response is missing expected field: data.markdown')
          );
        }
        resolve(result.data.markdown);
      })
      .catch((error) => {
        // A 404 is not an error.
        if (error?.response?.status === 404) {
          resolve('');
        } else {
          reject(error);
        }
      });
  });
}

export function saveDraft(entryDate, entryContent) {
  return fetch(`${process.env.VUE_APP_BACKEND_URL}/api/draft/${entryDate}`, {
    method: 'PUT',
    body: JSON.stringify({entryContent: entryContent}),
    credentials: 'include',
    headers: {'X-CSRF-Token': getCsrfToken()},
  }).then((response) => {
    if (response.ok) {
      return Promise.resolve();
    }
    return response.text().then((error) => {
      return Promise.reject(error);
    });
  });
}
