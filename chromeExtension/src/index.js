// @ts-check
/// <reference types="chrome" />
console.log("hi form the index.js");

try {
  let element = document.getElementById("status");
  if (element !== null) {
    element.addEventListener("click", async function a() {
      let [key, error] = await getKeyFromStorageOrBackend();
      if (error) {
        console.log(
          "error is there in getting the key and it is -->",
          error,
          "\n\n the key is",
          key,
        );
        return;
      }
      console.log("the key is -->", key);
    });
  }
} catch (e) {
  console.log("error ++--++", e);
}