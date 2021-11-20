import {processJsonResponse} from '@/controllers/Common.js';

export function exportGet() {
  return fetch(`${process.env.VUE_APP_BACKEND_URL}/api/export`, {
    credentials: 'include',
  }).then(processJsonResponse);
}
