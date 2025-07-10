import {
  getCsrfToken,
  processJsonResponse,
  processBlankResponse,
} from '@/controllers/Common.js';

export function getForwardingAddress() {
  return fetch(`${process.env.VUE_APP_BACKEND_URL}/api/forwarding-address`, {
    credentials: 'include',
  }).then(processJsonResponse);
}

export function saveForwardingAddress(url) {
  return fetch(`${process.env.VUE_APP_BACKEND_URL}/api/forwarding-address`, {
    method: 'PUT',
    credentials: 'include',
    headers: {
      'X-CSRF-Token': getCsrfToken(),
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({forwardingUrl: url}),
  }).then(processBlankResponse);
}

export function deleteForwardingAddress() {
  return fetch(`${process.env.VUE_APP_BACKEND_URL}/api/forwarding-address`, {
    method: 'DELETE',
    credentials: 'include',
    headers: {
      'X-CSRF-Token': getCsrfToken(),
    },
  }).then(processBlankResponse);
}
