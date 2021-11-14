import {getCsrfToken, processJsonResponse} from '@/controllers/Common.js';

export function getDraft(entryDate) {
  return fetch(`${process.env.VUE_APP_BACKEND_URL}/api/draft/${entryDate}`, {
    credentials: 'include',
  })
    .then(processJsonResponse)
    .then((draft) => {
      if (!Object.prototype.hasOwnProperty.call(draft, 'markdown')) {
        return Promise.reject(
          new Error('response is missing expected field: data.markdown')
        );
      }
      Promise.resolve(draft.markdown);
    })
    .catch((error) => {
      // A 404 is not an error.
      if (error?.response?.status === 404) {
        Promise.resolve('');
      } else {
        Promise.reject(error);
      }
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
