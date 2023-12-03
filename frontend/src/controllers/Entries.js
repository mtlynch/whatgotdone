import {getCsrfToken, processJsonResponse} from '@/controllers/Common.js';

export function getEntriesFromUser(username, project = null) {
  let url = `${process.env.VUE_APP_BACKEND_URL}/api/entries/${username}`;
  if (project) {
    url += `/project/${project}`;
  }
  return fetch(url)
    .then(processJsonResponse)
    .then((entriesRaw) => {
      let entries = [];
      for (const entry of entriesRaw) {
        entries.push({
          permalink: `/${username}/${entry.date}`,
          author: username,
          date: entry.date,
          lastModified: new Date(entry.lastModified),
          markdown: entry.markdown,
        });
      }
      // Sort newest to oldest.
      entries.sort((a, b) => b.date - a.date);
      return Promise.resolve(entries);
    })
    .catch((error) => {
      return Promise.reject(error);
    });
}

export function saveEntry(entryDate, entryContent) {
  return fetch(`${process.env.VUE_APP_BACKEND_URL}/api/entry/${entryDate}`, {
    method: 'PUT',
    credentials: 'include',
    headers: {'X-CSRF-Token': getCsrfToken()},
    body: JSON.stringify({entryContent: entryContent}),
  }).then(processJsonResponse);
}

export function entryDelete(entryDate) {
  return fetch(`${process.env.VUE_APP_BACKEND_URL}/api/entry/${entryDate}`, {
    method: 'DELETE',
    credentials: 'include',
    headers: {'X-CSRF-Token': getCsrfToken()},
  }).then((response) => {
    if (response.ok) {
      return Promise.resolve();
    } else {
      return response.text().then((error) => {
        return Promise.reject(error);
      });
    }
  });
}
