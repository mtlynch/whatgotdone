import {getCsrfToken, processJsonResponse} from '@/controllers/Common.js';

export function getReactions(entryAuthor, entryDate) {
  return fetch(
    `${process.env.VUE_APP_BACKEND_URL}/api/reactions/entry/${entryAuthor}/${entryDate}`
  ).then(processJsonResponse);
}

export function setReaction(entryAuthor, entryDate, reaction) {
  return fetch(
    `${process.env.VUE_APP_BACKEND_URL}/api/reactions/entry/${entryAuthor}/${entryDate}`,
    {
      method: 'POST',
      credentials: 'include',
      headers: {'X-CSRF-Token': getCsrfToken()},
      body: JSON.stringify({
        reactionSymbol: reaction,
      }),
    }
  );
}

export function deleteReaction(entryAuthor, entryDate) {
  return fetch(
    `${process.env.VUE_APP_BACKEND_URL}/api/reactions/entry/${entryAuthor}/${entryDate}`,
    {
      method: 'DELETE',
      credentials: 'include',
      headers: {'X-CSRF-Token': getCsrfToken()},
    }
  );
}
