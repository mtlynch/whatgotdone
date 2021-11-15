import {processJsonResponse} from '@/controllers/Common.js';

export function getPageViews(path) {
  return fetch(
    `${process.env.VUE_APP_BACKEND_URL}/api/pageViews` +
      new URLSearchParams({
        path,
      })
  )
    .then(processJsonResponse)
    .then((viewData) => {
      return Promise.resolve(viewData.views);
    });
}
