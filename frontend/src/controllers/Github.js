import axios from 'axios';

function parseUrl(url) {
  const parts = url.split('/');
  console.log(parts);
  return {owner: parts[3], repo: parts[4], issueNumber: parts[6]};
}

export function getIssueMetadata(issueUrl) {
  return new Promise(function(resolve, reject) {
    const parsed = parseUrl(issueUrl);
    const url = `https://api.github.com/repos/${parsed.owner}/${parsed.repo}/issues/${parsed.issueNumber}`;
    axios
      .get(url, {withCredentials: true})
      .then(result => {
        console.log('resolving');
        resolve(result.data);
      })
      .catch(err => {
        console.log('rejecting');
        reject(err);
      });
  });
}
