import axios from 'axios';

import getCsrfToken from '@/controllers/CsrfToken.js';

export function getReactions(entryAuthor, entryDate) {
  const url = `${process.env.VUE_APP_BACKEND_URL}/api/reactions/entry/${entryAuthor}/${entryDate}`;
  axios.get(url).then((result) => {
    return result.data;
  });
}

export function setReaction(entryAuthor, entryDate, reaction) {
  const url = `${process.env.VUE_APP_BACKEND_URL}/api/reactions/entry/${entryAuthor}/${entryDate}`;
  axios.post(
    url,
    {
      reactionSymbol: reaction,
    },
    {withCredentials: true, headers: {'X-CSRF-Token': getCsrfToken()}}
  );
}
