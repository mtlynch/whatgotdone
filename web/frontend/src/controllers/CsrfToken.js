export default function getCsrfToken() {
  console.log("Calling getCsrfToken");
  console.log(document
    .querySelector("meta[name='csrf-token']")
    .getAttribute("content"));
  return document
    .querySelector("meta[name='csrf-token']")
    .getAttribute("content");
}