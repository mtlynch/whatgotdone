import axios from 'axios';

import {getCsrfToken} from '@/controllers/Common.js';

export function getEntriesFromUser(username, project = null) {
  return new Promise(function (resolve, reject) {
    let url = `${process.env.VUE_APP_BACKEND_URL}/api/entries/${username}`;
    if (project) {
      url += `/project/${project}`;
    }
    axios
      .get(url)
      .then((result) => {
        let entries = [];
        for (const entry of result.data) {
          entries.push({
            permalink: `/${username}/${entry.date}`,
            author: username,
            date: new Date(entry.date),
            lastModified: new Date(entry.lastModified),
            markdown: entry.markdown,
          });
        }
        // Sort newest to oldest.
        entries.sort((a, b) => b.date - a.date);
        resolve(entries);
      })
      .catch((error) => {
        reject(error);
      });
  });
}

export function saveEntry(entryDate, entryContent) {
  return new Promise(function (resolve, reject) {
    const url = `${process.env.VUE_APP_BACKEND_URL}/api/entry/${entryDate}`;
    axios
      .post(
        url,
        {
          entryContent: entryContent,
        },
        {withCredentials: true, headers: {'X-CSRF-Token': getCsrfToken()}}
      )
      .then((result) => {
        resolve(result.data);
      })
      .catch((error) => {
        reject(error);
      });
  });
}
