import {getCsrfToken, processBlankResponse} from '@/controllers/Common.js';

export function deleteAvatar() {
  return fetch(`${process.env.VUE_APP_BACKEND_URL}/api/user/avatar`, {
    method: 'DELETE',
    credentials: 'include',
    headers: {
      'X-CSRF-Token': getCsrfToken(),
    },
  }).then(processBlankResponse);
}

export function uploadAvatar(image) {
  const formData = new FormData();
  formData.append('file', image);
  return fetch(`${process.env.VUE_APP_BACKEND_URL}/api/user/avatar`, {
    method: 'PUT',
    credentials: 'include',
    headers: {
      'X-CSRF-Token': getCsrfToken(),
    },
    body: formData,
  }).then(processBlankResponse);
}
