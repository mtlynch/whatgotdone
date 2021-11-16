import {processJsonResponse} from '@/controllers/Common.js';

// Number of entries to request from the server on each request;
const updateSize = 15;

export function getRecent(start) {
  return fetch(
    `${process.env.VUE_APP_BACKEND_URL}/api/recentEntries?start=${start}&limit=${updateSize}`
  ).then((result) => {
    // Transform each response data into entry object
    const recentEntries = result.data.map((rawEntry) => {
      return processEntry(rawEntry);
    });
    return Promise.resolve(recentEntries);
  });
}

// Retrieve recent entries from users that the logged-in user is following.
export function getRecentFollowing(start) {
  return fetch(
    `${process.env.VUE_APP_BACKEND_URL}/api/entries/following?start=${start}&limit=${updateSize}`,
    {credentials: 'include'}
  )
    .then(processJsonResponse)
    .then((recentEntries) => {
      if (!recentEntries.entries) {
        Promise.resolve([]);
        return;
      }
      // Transform each response data into entry object
      const entries = recentEntries.entries.map((rawEntry) => {
        return processEntry(rawEntry);
      });
      Promise.resolve(entries);
    });
}

// Merges two arrays of entries so that there are no duplicate entries.
export function mergeEntryArrays(a, b) {
  let merged = Array.from(a);
  let keySet = new Set();
  a.forEach((entry) => {
    keySet.add(entry.permalink);
  });

  b.forEach((newEntry) => {
    if (!keySet.has(newEntry.permalink)) {
      merged.push(newEntry);
    }
  });

  return merged;
}

function processEntry(entry) {
  const formattedDate = new Date(entry.date).toISOString().slice(0, 10);
  return {
    permalink: `/${entry.author}/${formattedDate}`,
    author: entry.author,
    date: new Date(entry.date),
    markdown: entry.markdown,
  };
}
