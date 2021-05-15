import axios from 'axios';

export function getPageViews(path) {
  const url = `${process.env.VUE_APP_BACKEND_URL}/api/pageViews`;
  axios
    .get(url, {
      params: {
        path: path,
      },
    })
    .then((result) => {
      return result?.data?.views;
    });
}
