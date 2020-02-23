import axios from 'axios';

// Number of entries to request from the server on each request;
const updateSize = 15;

export function getRecent(start) {
  return new Promise(function(resolve, reject) {
    const url = `${process.env.VUE_APP_BACKEND_URL}/api/recentEntries?start=${start}&limit=${updateSize}`;
    axios
      .get(url)
      .then(result => {
        // Transform each response data into entry object
        const recentEntries = result.data.map(rawEntry => {
          return processEntry(rawEntry);
        });
        resolve(recentEntries);
      })
      .catch(err => reject(err));
  });
}

// Merges two arrays of entries so that there are no duplicate keys.
export function mergeEntryArrays(a, b) {
  let merged = Array.from(a);
  let keySet = new Set();
  a.forEach(entry => {
    keySet.add(entry.key);
  });

  b.forEach(newEntry => {
    if (!keySet.has(newEntry.key)) {
      merged.push(newEntry);
    }
  });

  return merged;
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
