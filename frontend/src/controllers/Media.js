import {getCsrfToken, processJsonResponse} from '@/controllers/Common.js';

export function uploadMedia(media) {
  let formData = new FormData();
  formData.append('file', media);
  return fetch(`${process.env.VUE_APP_BACKEND_URL}/api/media`, {
    method: 'PUT',
    credentials: 'include',
    headers: {
      'X-CSRF-Token': getCsrfToken(),
    },
    body: formData,
  })
    .then(processJsonResponse)
    .then((result) => {
      Promise.resolve(result.url);
    });
}
