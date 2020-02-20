import axios from 'axios';
import store from '../store.js';

// Number of entries to request from the server on each request;
const updateSize = 15;

export function refreshRecent() {
  getRecent(0, updateSize, recentEntries => {
    store.commit('setRecent', recentEntries);
  });
}

export function extendRecent(callback) {
  let recentEntries = store.state.recentEntries;
  getRecent(recentEntries.length, updateSize, newEntries => {
    // Extract keys from recentEntries array
    const recentEntriesKeySet = getRecentEntriesKey(recentEntries);
    // Loop through new entries and add to recentEntries store only if key doesn't exist
    newEntries.forEach(newEntry => {
      if (!recentEntriesKeySet.has(newEntry.key)) {
        recentEntries.push(newEntry);
      }
    });
    store.commit('setRecent', recentEntries);
    callback();
  });
}

export function getRecent(start, limit, callback) {
  const url = `${process.env.VUE_APP_BACKEND_URL}/api/recentEntries?start=${start}&limit=${limit}`;
  axios.get(url).then(result => {
    // Transform each response data into entry object
    const recentEntries = result.data.map(rawEntry => {
      return processEntry(rawEntry);
    });
    callback(recentEntries);
  });
}

function processEntry(entry) {
  const formattedDate = new Date(entry.date).toISOString().slice(0, 10);
  return {
    key: `/${entry.author}/${formattedDate}`,
    author: entry.author,
    date: new Date(entry.date),
    markdown: entry.markdown,
  };
}

//Loop through all entries and put entry key in Set
function getRecentEntriesKey(recentEntries) {
  const entriesKeySet = new Set();

  recentEntries.forEach(recentEntry => {
    entriesKeySet.add(recentEntry.key);
  });

  return entriesKeySet;
}
