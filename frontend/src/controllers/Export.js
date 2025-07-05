import {processJsonResponse} from '@/controllers/Common.js';

export function exportGet() {
  return fetch(`${process.env.VUE_APP_BACKEND_URL}/api/export`, {
    credentials: 'include',
  }).then(processJsonResponse);
}

export function exportMarkdown() {
  return fetch(`${process.env.VUE_APP_BACKEND_URL}/api/export/markdown`, {
    credentials: 'include',
  }).then((response) => {
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return response.blob().then((blob) => ({
      blob,
      filename: extractFilenameFromContentDisposition(
        response.headers.get('Content-Disposition')
      ),
    }));
  });
}

function extractFilenameFromContentDisposition(contentDisposition) {
  if (!contentDisposition) {
    return 'export.zip';
  }
  const filenameMatch = contentDisposition.match(/filename="([^"]+)"/);
  return filenameMatch ? filenameMatch[1] : 'export.zip';
}
