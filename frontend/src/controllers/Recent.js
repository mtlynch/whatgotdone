import axios from 'axios';
import store from '../store.js';

// Number of entries to request from the server on each request;
const updateSize = 15;

export function refreshRecent() {
  getRecent(0, updateSize, recentEntries => {
    store.commit('setRecent', recentEntries);
  });
}

export function extendRecent() {
  let recentEntries = store.state.recentEntries;
  getRecent(
    recentEntries.length,
    recentEntries.length + updateSize,
    newEntries => {
      recentEntries.push(...newEntries);
      store.commit('setRecent', recentEntries);
    }
  );
}

export function getRecent(start, limit, callback) {
  const url = `${process.env.VUE_APP_BACKEND_URL}/api/recentEntries?start=${start}&limit=${limit}`;
  axios.get(url).then(result => {
    const recentEntries = [];
    for (const entry of result.data) {
      recentEntries.push(processEntry(entry));
    }
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
