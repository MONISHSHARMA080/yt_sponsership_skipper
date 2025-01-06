import { getKeyFromStorageOrBackend as o } from "./helper.js";
console.log("hi from the contentscript in ts");
async function n() {
  console.log("in the main");
  let [t, e] = await o();
  if (e) {
    console.log(
      "error is there in getting the key and it is -->",
      e,
      `

 the key is`,
      t
    );
    return;
  }
  console.log(
    "the key is  -->",
    t,
    "from the  content script and the error is->",
    e
  );
}
n();
