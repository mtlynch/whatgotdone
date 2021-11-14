import {getCsrfToken} from '@/controllers/Common.js';

export function getDraft(entryDate) {
  return fetch(`${process.env.VUE_APP_BACKEND_URL}/api/draft/${entryDate}`, {
    credentials: 'include',
  })
    .then((response) => {
      if (response.ok) {
        return response.json();
      }
      // A 404 is not an error.
      if (response.status === 404) {
        return Promise.resolve({markdown: ''});
      }
      return response.text().then((error) => {
        return Promise.reject(error);
      });
    })
    .then((draft) => {
      if (!Object.prototype.hasOwnProperty.call(draft, 'markdown')) {
        return Promise.reject(
          new Error('response is missing expected field: data.markdown')
        );
      }
      Promise.resolve(draft.markdown);
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
