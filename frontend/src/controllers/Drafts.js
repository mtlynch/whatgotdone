import axios from 'axios';

import getCsrfToken from '@/controllers/CsrfToken.js';

export function getDraft(entryDate) {
  return new Promise(function(resolve, reject) {
    const url = `${process.env.VUE_APP_BACKEND_URL}/api/draft/${entryDate}`;
    axios
      .get(url, {withCredentials: true})
      .then(result => {
        resolve(result.data.markdown);
      })
      .catch(error => {
        if (error.response.status == 404) {
          resolve('');
        } else {
          reject(error);
        }
      });
  });
}

export function saveDraft(entryDate, entryContent) {
  return new Promise(function(resolve, reject) {
    const url = `${process.env.VUE_APP_BACKEND_URL}/api/entry/${entryDate}`;
    axios
      .post(
        url,
        {
          entryContent: entryContent,
        },
        {withCredentials: true, headers: {'X-CSRF-Token': getCsrfToken()}}
      )
      .then(result => {
        resolve(result.data);
      })
      .catch(error => {
        reject(error);
      });
  });
}
