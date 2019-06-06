import axios from "axios";
import store from "../store.js";

export default function refreshRecent() {
  const url = `${process.env.VUE_APP_BACKEND_URL}/api/recentEntries`;
  axios
    .get(url)
    .then(result => {
      const recentEntries = [];
      for (const entry of result.data) {
        const formattedDate = new Date(entry.date).toISOString().slice(0, 10);
        recentEntries.push({
          key: `/${entry.author}/${formattedDate}`,
          author: entry.author,
          date: entry.date,
          markdown: entry.markdown
        });
      }
      store.commit("setRecent", recentEntries);
    });
}