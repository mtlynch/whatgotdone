import axios from 'axios';

import getCsrfToken from '@/controllers/CsrfToken.js';

export function getReactions(entryAuthor, entryDate) {
  return new Promise(function (resolve, reject) {
    const url = `${process.env.VUE_APP_BACKEND_URL}/api/reactions/entry/${entryAuthor}/${entryDate}`;
    axios
      .get(url)
      .then((result) => {
        resolve(result.data);
      })
      .catch((error) => {
        reject(error);
      });
  });
}

export function setReaction(entryAuthor, entryDate, reaction) {
  return new Promise(function (resolve, reject) {
    const url = `${process.env.VUE_APP_BACKEND_URL}/api/reactions/entry/${entryAuthor}/${entryDate}`;
    axios
      .post(
        url,
        {
          reactionSymbol: reaction,
        },
        {withCredentials: true, headers: {'X-CSRF-Token': getCsrfToken()}}
      )
      .then(() => {
        resolve();
      })
      .catch((error) => {
        reject(error);
      });
  });
}
