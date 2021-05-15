import axios from 'axios';

import getCsrfToken from '@/controllers/CsrfToken.js';

export function getDraft(entryDate) {
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
      return result.data.markdown;
    })
    .catch((error) => {
      // A 404 is not an error.
      if (error?.response?.status === 404) {
        return;
      } else {
        throw error;
      }
    });
}

export function saveDraft(entryDate, entryContent) {
  const url = `${process.env.VUE_APP_BACKEND_URL}/api/draft/${entryDate}`;
  axios
    .post(
      url,
      {
        entryContent: entryContent,
      },
      {withCredentials: true, headers: {'X-CSRF-Token': getCsrfToken()}}
    )
    .then((result) => {
      return result.data;
    })
    .catch((error) => {
      throw error;
    });
}
